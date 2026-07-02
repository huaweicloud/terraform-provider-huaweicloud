package secmaster

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

var classifierNonUpdatableParams = []string{"workspace_id", "dataclass_id", "data_source"}

// There are many issues with creating APIs currently, and the actual use of many parameters does not match the API
// documentation. Currently, only some mandatory parameters are supported, and other optional parameters are not
// currently supported.

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/mappings/classifiers
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/mappings/classifiers/{classifier_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/mappings/classifiers/{classifier_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/mappings/classifiers/{classifier_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/mappings/{mapping_id}
func ResourceClassifier() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClassifierCreate,
		ReadContext:   resourceClassifierRead,
		UpdateContext: resourceClassifierUpdate,
		DeleteContext: resourceClassifierDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceClassifierImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(classifierNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region where the resource is located.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the workspace ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the name.",
			},
			"dataclass_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the data class ID.",
			},
			// Query API no return.
			"data_source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the data source.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the description.",
			},
			"classifier": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"direct_classifier": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies whether to classify directly.",
						},
					},
				},
				Description: "Specifies the classifier information.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"mapping_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mapping ID.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project ID.",
			},
			"dataclass_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The data class name.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status.",
			},
			"complete_degree": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The completion degree.",
			},
			"instance_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The number of associated instances.",
			},
			"built_in": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the data is built-in.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},
			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator ID.",
			},
			"creator_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator name.",
			},
			"modifier_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The modifier ID.",
			},
			"modifier_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The modifier name.",
			},
		},
	}
}

func buildClassifierBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":         d.Get("name"),
		"dataclass_id": d.Get("dataclass_id"),
		"data_source":  d.Get("data_source"),
		"description":  d.Get("description"),
		"classifier":   buildResourceClassifierBodyParams(d.Get("classifier").([]interface{})),
	}
}

func buildResourceClassifierBodyParams(rawParams []interface{}) map[string]interface{} {
	raw, ok := rawParams[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"direct_classifier": raw["direct_classifier"],
	}
}

func resourceClassifierCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/classifiers"
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildClassifierBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster classifier: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.classifier.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster classifier: ID is not found in API response")
	}

	d.SetId(id)

	return resourceClassifierRead(ctx, d, meta)
}

func resourceClassifierRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/classifiers/{classifier_id}"
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{classifier_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster classifier: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The query API always returns `200`, and when the resource does not exist, the value of `data.mapping_info`
	// is nil. So this is used as a condition to determine whether resources exist.
	mappingInfo := utils.PathSearch("data.mapping_info", respBody, nil)
	if mappingInfo == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving SecMaster classifier")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("workspace_id", utils.PathSearch("workspace_id", mappingInfo, nil)),
		d.Set("name", utils.PathSearch("name", mappingInfo, nil)),
		d.Set("dataclass_id", utils.PathSearch("dataclass_id", mappingInfo, nil)),
		d.Set("description", utils.PathSearch("description", mappingInfo, nil)),
		d.Set("classifier", flattenResourceClassifier(
			utils.PathSearch("data.direct_classifier", respBody, nil))),
		d.Set("mapping_id", utils.PathSearch("id", mappingInfo, nil)),
		d.Set("project_id", utils.PathSearch("project_id", mappingInfo, nil)),
		d.Set("dataclass_name", utils.PathSearch("dataclass_name", mappingInfo, nil)),
		d.Set("status", utils.PathSearch("status", mappingInfo, nil)),
		d.Set("complete_degree", utils.PathSearch("complete_degree", mappingInfo, nil)),
		d.Set("instance_num", utils.PathSearch("instance_num", mappingInfo, nil)),
		d.Set("built_in", utils.PathSearch("built_in", mappingInfo, nil)),
		d.Set("update_time", utils.PathSearch("update_time", mappingInfo, nil)),
		d.Set("create_time", utils.PathSearch("create_time", mappingInfo, nil)),
		d.Set("creator_id", utils.PathSearch("creator_id", mappingInfo, nil)),
		d.Set("creator_name", utils.PathSearch("creator_name", mappingInfo, nil)),
		d.Set("modifier_id", utils.PathSearch("modifier_id", mappingInfo, nil)),
		d.Set("modifier_name", utils.PathSearch("modifier_name", mappingInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResourceClassifier(directClassifier interface{}) []interface{} {
	if directClassifier == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"direct_classifier": directClassifier,
		},
	}
}

func buildUpdateClassifierBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"id":           d.Get("mapping_id"),
		"name":         d.Get("name"),
		"dataclass_id": d.Get("dataclass_id"),
		"data_source":  d.Get("data_source"),
		"description":  d.Get("description"),
		"classifier":   buildResourceClassifierBodyParams(d.Get("classifier").([]interface{})),
	}
}

func resourceClassifierUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/classifiers/{classifier_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{classifier_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildUpdateClassifierBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster classifier: %s", err)
	}

	return resourceClassifierRead(ctx, d, meta)
}

func resourceClassifierDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
		mappingId   = d.Get("mapping_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	classifierUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/classifiers/{classifier_id}"
	classifierPath := client.Endpoint + classifierUrl
	classifierPath = strings.ReplaceAll(classifierPath, "{project_id}", client.ProjectID)
	classifierPath = strings.ReplaceAll(classifierPath, "{workspace_id}", workspaceId)
	classifierPath = strings.ReplaceAll(classifierPath, "{classifier_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	_, err = client.Request("DELETE", classifierPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster classifier: %s", err)
	}

	// After deleting the classifier, there will still be mapping relationships present.
	// Need to delete the remaining mapping again.
	if mappingId == "" {
		return nil
	}

	mappingPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/{mapping_id}"
	mappingPath = strings.ReplaceAll(mappingPath, "{project_id}", client.ProjectID)
	mappingPath = strings.ReplaceAll(mappingPath, "{workspace_id}", workspaceId)
	mappingPath = strings.ReplaceAll(mappingPath, "{mapping_id}", mappingId)

	_, err = client.Request("DELETE", mappingPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster mapping info: %s", err)
	}

	return nil
}

func resourceClassifierImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	importIdParts := strings.Split(importId, "/")
	if len(importIdParts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want '<workspace_id>/<id>', but got '%s'", importId)
	}

	d.SetId(importIdParts[1])

	return []*schema.ResourceData{d}, d.Set("workspace_id", importIdParts[0])
}
