package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	scanTemplateClassificationNonUpdatableParams = []string{
		"template_id",
		"parent_id",
	}

	scanTemplateClassificationNotFoundErrCodes = []string{
		"dsc.600000025", // The template does not exist.
		"dsc.10000009",  // The DSC instance does not exist.
	}
)

// @API DSC POST /v1/{project_id}/scan-templates/{template_id}/classifications
// @API DSC GET /v1/{project_id}/scan-templates/{template_id}/classifications
// @API DSC PUT /v1/{project_id}/scan-templates/{template_id}/classifications/{classification_id}
// @API DSC POST /v1/{project_id}/scan-templates/{template_id}/classifications/batch-delete
func ResourceScanTemplateClassification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScanTemplateClassificationCreate,
		ReadContext:   resourceScanTemplateClassificationRead,
		UpdateContext: resourceScanTemplateClassificationUpdate,
		DeleteContext: resourceScanTemplateClassificationDelete,

		CustomizeDiff: config.FlexibleForceNew(scanTemplateClassificationNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceScanTemplateClassificationImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the scan template classification is located.`,
			},

			// Required parameters.
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the scan template.`,
			},
			"classification_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the classification.`,
			},

			// Optional parameters.
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The parent classification ID.`,
			},

			// Attributes.
			"depth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The depth of the classification.`,
			},
			"root_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the root to which the classification belongs.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the classification is created, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the classification, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildCreateScanTemplateClassificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"classification_name": d.Get("classification_name"),
		"parent_id":           utils.ValueIgnoreEmpty(d.Get("parent_id")),
	}
}

func createScanTemplateClassification(client *golangsdk.ServiceClient, templateId string, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/scan-templates/{template_id}/classifications"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_id}", templateId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateScanTemplateClassificationBodyParams(d)),
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceScanTemplateClassificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		templateId = d.Get("template_id").(string)
		name       = d.Get("classification_name").(string)
	)
	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = createScanTemplateClassification(client, templateId, d)
	if err != nil {
		return diag.Errorf("error creating scan classification under template (%s): %s", templateId, err)
	}

	// After creation, query the classification list to get corresponding classification ID.
	respBody, err := listScanTemplateClassifications(client, templateId)
	if err != nil {
		return diag.Errorf("error getting classifications under template (%s): %s", templateId, err)
	}

	classification := findScanTemplateClassificationByNameAndParentId(utils.PathSearch("classification_trees",
		respBody, make([]interface{}, 0)).([]interface{}), name, d.Get("parent_id").(string))
	classificationId := utils.PathSearch("id", classification, "").(string)
	if classificationId == "" {
		return diag.Errorf("unable to find the classification ID from API response")
	}

	d.SetId(classificationId)

	return resourceScanTemplateClassificationRead(ctx, d, meta)
}

func findScanTemplateClassificationByNameAndParentId(classification []interface{}, name, parentId string) interface{} {
	if len(classification) == 0 {
		return nil
	}

	for _, item := range classification {
		isMatch := utils.PathSearch("name", item, "").(string) == name
		if parentId != "" {
			isMatch = isMatch && utils.PathSearch("parent_id", item, "").(string) == parentId
		}

		if isMatch {
			return item
		}

		children := findScanTemplateClassificationByNameAndParentId(utils.PathSearch("children", item, make([]interface{}, 0)).([]interface{}),
			name, parentId)
		if children != nil {
			return children
		}
	}

	return nil
}

func listScanTemplateClassifications(client *golangsdk.ServiceClient, templateId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/scan-templates/{template_id}/classifications"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_id}", templateId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func GetScanTemplateClassificationById(client *golangsdk.ServiceClient, templateId, classificationId string) (interface{}, error) {
	respBody, err := listScanTemplateClassifications(client, templateId)
	if err != nil {
		return nil, err
	}

	classification := findScanTemplateClassificationById(utils.PathSearch("classification_trees",
		respBody, make([]interface{}, 0)).([]interface{}), classificationId)
	if classification == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/scan-templates/{template_id}/classifications",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the classification (%s) does not exist", classificationId)),
			},
		}
	}

	return classification, nil
}

func resourceScanTemplateClassificationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		templateId       = d.Get("template_id").(string)
		classificationId = d.Id()
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	classification, err := GetScanTemplateClassificationById(client, templateId, classificationId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(
				common.ConvertExpected400ErrInto404Err(err, "error_code", scanTemplateClassificationNotFoundErrCodes...),
				"error_code",
				scanTemplateClassificationNotFoundErrCodes...,
			),
			fmt.Sprintf("error retrieving scan template classifications (%s)", classificationId),
		)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Required parameters.
		d.Set("template_id", templateId),
		d.Set("classification_name", utils.PathSearch("name", classification, nil)),
		// Optional parameters.
		d.Set("parent_id", utils.PathSearch("parent_id", classification, nil)),
		// Attributes.
		d.Set("depth", utils.PathSearch("depth", classification, nil)),
		d.Set("root_id", utils.PathSearch("root_id", classification, nil)),
		d.Set("create_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
			classification, float64(0)).(float64))/1000, false)),
		d.Set("update_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time",
			classification, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func findScanTemplateClassificationById(classifications []interface{}, classificationId string) interface{} {
	for _, classification := range classifications {
		if utils.PathSearch("id", classification, "").(string) == classificationId {
			return classification
		}

		children := utils.PathSearch("children", classification, make([]interface{}, 0)).([]interface{})
		if result := findScanTemplateClassificationById(children, classificationId); result != nil {
			return result
		}
	}

	return nil
}

func resourceScanTemplateClassificationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		templateId       = d.Get("template_id").(string)
		classificationId = d.Id()
		name             = d.Get("classification_name").(string)
	)

	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	if d.HasChangeExcept("enable_force_new") {
		err = updateScanTemplateClassificationName(client, templateId, classificationId, name)
		if err != nil {
			return diag.Errorf("error updating classification (%s) under the template (%s): %s", classificationId, templateId, err)
		}
	}

	return resourceScanTemplateClassificationRead(ctx, d, meta)
}

func updateScanTemplateClassificationName(client *golangsdk.ServiceClient, templateId, classificationId, name string) error {
	httpUrl := "v1/{project_id}/scan-templates/{template_id}/classifications/{classification_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{template_id}", templateId)
	updatePath = strings.ReplaceAll(updatePath, "{classification_id}", classificationId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"classification_name": name,
		},
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func deleteScanTemplateClassification(client *golangsdk.ServiceClient, templateId string, classificationId string) error {
	httpUrl := "v1/{project_id}/scan-templates/{template_id}/classifications/batch-delete"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_id}", templateId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"classification_ids": []string{classificationId},
		},
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceScanTemplateClassificationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		templateId = d.Get("template_id").(string)
	)
	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = deleteScanTemplateClassification(client, templateId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", scanTemplateClassificationNotFoundErrCodes...),
			fmt.Sprintf("error deleting classification (%v) under template (%s)", d.Get("classification_name"), templateId),
		)
	}

	return nil
}

func resourceScanTemplateClassificationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<template_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("template_id", parts[0])
}
