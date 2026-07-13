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
	assetDomainLabelNonUpdatableParams = []string{
		"name",
		"parent_id",
	}

	assetDomainLabelNotFoundErrCodes = []string{
		// dsc.10000009: Current instance does not exist.
		"dsc.10000009",
		// dsc.40000007: The label does not exist.
		"dsc.40000007",
	}
)

// @API DSC POST /v1/{project_id}/metadata/asset-domain-labels
// @API DSC GET /v1/{project_id}/metadata/asset-domain-labels
// @API DSC DELETE /v1/{project_id}/metadata/asset-domain-labels/{id}
func ResourceAssetDomainLabel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetDomainLabelCreate,
		ReadContext:   resourceAssetDomainLabelRead,
		UpdateContext: resourceAssetDomainLabelUpdate,
		DeleteContext: resourceAssetDomainLabelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAssetDomainLabelImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(assetDomainLabelNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the asset domain label is located.`,
			},

			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the label.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the parent to which the label belongs.`,
			},

			// Attributes.
			"depth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The depth of the label.`,
			},
			"is_leaf": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Whether the label is a leaf node.`,
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

func buildCreateAssetDomainLabelBodyParams(name, parentId string) map[string]interface{} {
	return map[string]interface{}{
		"name":      name,
		"parent_id": parentId,
	}
}

func createAssetDomainLabel(client *golangsdk.ServiceClient, params map[string]interface{}) error {
	httpUrl := "v1/{project_id}/metadata/asset-domain-labels"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         params,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func resourceAssetDomainLabelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		name     = d.Get("name").(string)
		parentId = d.Get("parent_id").(string)
	)
	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = createAssetDomainLabel(client, buildCreateAssetDomainLabelBodyParams(name, parentId))
	if err != nil {
		return diag.Errorf("error creating asset domain label: %s", err)
	}

	label, err := GetAssetDomainLabelByName(client, name, parentId)
	if err != nil {
		return diag.Errorf("error getting asset domain label (%s) under parent ID (%s): %s", name, parentId, err)
	}

	labelId, _ := utils.PathSearch("id", label, "").(string)
	if labelId == "" {
		return diag.Errorf("unable to find the ID of the asset domain label from the API response")
	}

	d.SetId(labelId)

	return resourceAssetDomainLabelRead(ctx, d, meta)
}

func listAssetDomainLabels(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v1/{project_id}/metadata/asset-domain-labels"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func findLabelByNameAndParentId(labels []interface{}, name, parentId string) interface{} {
	if len(labels) == 0 {
		return nil
	}

	for _, item := range labels {
		if utils.PathSearch("name", item, "").(string) == name && utils.PathSearch("parent_id", item, "").(string) == parentId {
			return item
		}

		label := findLabelByNameAndParentId(utils.PathSearch("sons", item, make([]interface{}, 0)).([]interface{}), name, parentId)
		if label != nil {
			return label
		}
	}

	return nil
}

func GetAssetDomainLabelByName(client *golangsdk.ServiceClient, name, parentId string) (interface{}, error) {
	respBody, err := listAssetDomainLabels(client)
	if err != nil {
		return nil, common.ConvertExpected401ErrInto404Err(err, "error_code", assetDomainLabelNotFoundErrCodes[0])
	}

	label := findLabelByNameAndParentId(utils.PathSearch("sons", respBody, make([]interface{}, 0)).([]interface{}), name, parentId)
	if label == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/metadata/asset-domain-labels",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the asset domain label with name (%s) and parent ID (%s) does not exist", name, parentId)),
			},
		}
	}

	return label, nil
}

func resourceAssetDomainLabelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		name     = d.Get("name").(string)
		parentId = d.Get("parent_id").(string)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	label, err := GetAssetDomainLabelByName(client, name, parentId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving asset domain label (%s) under parent ID (%s)", name, parentId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", label, nil)),
		d.Set("parent_id", utils.PathSearch("parent_id", label, nil)),
		// Attributes.
		d.Set("depth", utils.PathSearch("depth", label, nil)),
		d.Set("is_leaf", utils.PathSearch("is_leaf", label, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAssetDomainLabelUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func deleteAssetDomainLabel(client *golangsdk.ServiceClient, labelId string) error {
	httpUrl := "v1/{project_id}/metadata/asset-domain-labels/{id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", labelId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceAssetDomainLabelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = deleteAssetDomainLabel(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(common.ConvertExpected401ErrInto404Err(err, "error_code",
			assetDomainLabelNotFoundErrCodes...), "error_code", assetDomainLabelNotFoundErrCodes...),
			fmt.Sprintf("error deleting asset domain label (%s)", d.Get("name").(string)))
	}

	return nil
}

func resourceAssetDomainLabelImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<name>/<parent_id>', but got '%s'", importedId)
	}

	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	name := parts[0]
	parentId := parts[1]
	label, err := GetAssetDomainLabelByName(client, name, parentId)
	if err != nil {
		return nil, err
	}

	d.SetId(utils.PathSearch("id", label, "").(string))

	mErr := multierror.Append(nil,
		d.Set("name", name),
		d.Set("parent_id", parentId),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
