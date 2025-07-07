// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CES
// ---------------------------------------------------------------

package ces

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CES POST /v2/{project_id}/resource-groups
// @API CES GET /v2/{project_id}/resource-groups/{id}
// @API CES PUT /v2/{project_id}/resource-groups/{id}
// @API CES POST /v2/{project_id}/resource-groups/{id}/resources/batch-create
// @API CES POST /v2/{project_id}/resource-groups/{id}/resources/batch-delete
// @API CES POST /v2/{project_id}/resource-groups/batch-delete
func ResourceResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceGroupCreate,
		UpdateContext: resourceResourceGroupUpdate,
		ReadContext:   resourceResourceGroupRead,
		DeleteContext: resourceResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource group name.`,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  `Specifies the resource group type.`,
				ExactlyOneOf: []string{"resources"},
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project ID of the resource group.`,
			},
			"tags": common.TagsSchema(),
			"associated_eps_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project IDs where the resources from.`,
			},
			"resources": {
				Type:         schema.TypeList,
				Elem:         ResourceGroupResourcesOptsSchema(),
				Optional:     true,
				ExactlyOneOf: []string{"associated_eps_ids", "tags"},
				Description:  `Specifies the list of resources to add into the group.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
		},
	}
}

func ResourceGroupResourcesOptsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace in **service.item** format.`,
			},
			"dimensions": {
				Type:        schema.TypeList,
				Elem:        ResourceGroupResourcesOptsDimensionOptsSchema(),
				Required:    true,
				Description: `Specifies the list of dimensions.`,
			},
		},
	}
	return &sc
}

func ResourceGroupResourcesOptsDimensionOptsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the dimension name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the dimension value.`,
			},
		},
	}
	return &sc
}

func resourceResourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createResourceGroup: Create a CES resource group.
	var (
		createResourceGroupHttpUrl = "v2/{project_id}/resource-groups"
		createResourceGroupProduct = "ces"
	)
	createResourceGroupClient, err := cfg.NewServiceClient(createResourceGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	createResourceGroupPath := createResourceGroupClient.Endpoint + createResourceGroupHttpUrl
	createResourceGroupPath = strings.ReplaceAll(createResourceGroupPath, "{project_id}", createResourceGroupClient.ProjectID)

	createResourceGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createResourceGroupOpt.JSONBody = utils.RemoveNil(buildCreateResourceGroupBodyParams(d, cfg))
	createResourceGroupResp, err := createResourceGroupClient.Request("POST", createResourceGroupPath, &createResourceGroupOpt)
	if err != nil {
		return diag.Errorf("error creating CES resource group: %s", err)
	}

	createResourceGroupRespBody, err := utils.FlattenResponse(createResourceGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("group_id", createResourceGroupRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CES resource group: ID is not found in API response")
	}
	d.SetId(id)

	// add resources to CES resource group if resources is specified.
	if _, ok := d.GetOk("resources"); ok {
		var (
			batchCreateResourcesHttpUrl = "v2/{project_id}/resource-groups/{id}/resources/batch-create"
		)

		batchCreateResourcesPath := createResourceGroupClient.Endpoint + batchCreateResourcesHttpUrl
		batchCreateResourcesPath = strings.ReplaceAll(batchCreateResourcesPath, "{project_id}", createResourceGroupClient.ProjectID)
		batchCreateResourcesPath = strings.ReplaceAll(batchCreateResourcesPath, "{id}", d.Id())

		batchCreateResourcesOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		batchCreateResourcesOpt.JSONBody = utils.RemoveNil(buildBatchCreateDeleteResourcesBodyParams(d.Get("resources")))
		_, err = createResourceGroupClient.Request("POST", batchCreateResourcesPath, &batchCreateResourcesOpt)
		if err != nil {
			return diag.Errorf("error adding resources to resource group: %s", err)
		}
	}

	return resourceResourceGroupRead(ctx, d, meta)
}

func buildCreateResourceGroupBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_name":            utils.ValueIgnoreEmpty(d.Get("name")),
		"type":                  utils.ValueIgnoreEmpty(d.Get("type")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"tags":                  utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
		"association_ep_ids":    utils.ValueIgnoreEmpty(d.Get("associated_eps_ids")),
	}
	return bodyParams
}

func buildBatchCreateDeleteResourcesBodyParams(resourcesRaw interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"resources": buildBatchCreateResourcesRequestBodyResourcesOpts(resourcesRaw),
	}
	return bodyParams
}

func buildBatchCreateResourcesRequestBodyResourcesOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"namespace":  utils.ValueIgnoreEmpty(raw["namespace"]),
				"dimensions": buildResourcesOptsDimensionOpts(raw["dimensions"]),
			}
		}
		return rst
	}
	return nil
}

func buildResourcesOptsDimensionOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":  utils.ValueIgnoreEmpty(raw["name"]),
				"value": utils.ValueIgnoreEmpty(raw["value"]),
			}
		}
		return rst
	}
	return nil
}

func resourceResourceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getResourceGroup: Query the CES resource group detail
	var (
		getResourceGroupHttpUrl = "v2/{project_id}/resource-groups/{id}"
		getResourceGroupProduct = "ces"
	)
	getResourceGroupClient, err := cfg.NewServiceClient(getResourceGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	getResourceGroupPath := getResourceGroupClient.Endpoint + getResourceGroupHttpUrl
	getResourceGroupPath = strings.ReplaceAll(getResourceGroupPath, "{project_id}", getResourceGroupClient.ProjectID)
	getResourceGroupPath = strings.ReplaceAll(getResourceGroupPath, "{id}", d.Id())

	getResourceGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getResourceGroupResp, err := getResourceGroupClient.Request("GET", getResourceGroupPath, &getResourceGroupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ResourceGroup")
	}

	getResourceGroupRespBody, err := utils.FlattenResponse(getResourceGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// the API for getting resource is inconvenient, so it's not set
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("group_name", getResourceGroupRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getResourceGroupRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getResourceGroupRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", getResourceGroupRespBody, nil))),
		d.Set("associated_eps_ids", utils.PathSearch("association_ep_ids", getResourceGroupRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", getResourceGroupRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceResourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	resourceGroupId := d.Id()
	var (
		updateResourceGroupHttpUrl  = "v2/{project_id}/resource-groups/{id}"
		batchCreateResourcesHttpUrl = "v2/{project_id}/resource-groups/{id}/resources/batch-create"
		batchDeleteResourcesHttpUrl = "v2/{project_id}/resource-groups/{id}/resources/batch-delete"
		updateResourceGroupProduct  = "ces"
	)
	updateResourceGroupClient, err := cfg.NewServiceClient(updateResourceGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	updateResourceGroupChanges := []string{
		"name",
		"tags",
	}

	if d.HasChanges(updateResourceGroupChanges...) {
		// updateResourceGroup: Update the configuration of CES resource group

		updateResourceGroupPath := updateResourceGroupClient.Endpoint + updateResourceGroupHttpUrl
		updateResourceGroupPath = strings.ReplaceAll(updateResourceGroupPath, "{project_id}", updateResourceGroupClient.ProjectID)
		updateResourceGroupPath = strings.ReplaceAll(updateResourceGroupPath, "{id}", resourceGroupId)

		updateResourceGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateResourceGroupOpt.JSONBody = utils.RemoveNil(buildUpdateResourceGroupBodyParams(d))
		_, err = updateResourceGroupClient.Request("PUT", updateResourceGroupPath, &updateResourceGroupOpt)
		if err != nil {
			return diag.Errorf("error updating CES resource group: %s", err)
		}
	}

	// update resources
	if d.HasChange("resources") {
		oldResources, newResources := d.GetChange("resources")

		if len(oldResources.([]interface{})) > 0 {
			batchDeleteResourcesHttpPath := updateResourceGroupClient.Endpoint + batchDeleteResourcesHttpUrl
			batchDeleteResourcesHttpPath = strings.ReplaceAll(batchDeleteResourcesHttpPath, "{project_id}", updateResourceGroupClient.ProjectID)
			batchDeleteResourcesHttpPath = strings.ReplaceAll(batchDeleteResourcesHttpPath, "{id}", resourceGroupId)

			batchDeleteResourcesOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			batchDeleteResourcesOpt.JSONBody = utils.RemoveNil(buildBatchCreateDeleteResourcesBodyParams(oldResources))
			_, err = updateResourceGroupClient.Request("POST", batchDeleteResourcesHttpPath, &batchDeleteResourcesOpt)
			if err != nil {
				return diag.Errorf("error deleting resources to resource group: %s", err)
			}
		}

		if len(newResources.([]interface{})) > 0 {
			batchCreateResourcesPath := updateResourceGroupClient.Endpoint + batchCreateResourcesHttpUrl
			batchCreateResourcesPath = strings.ReplaceAll(batchCreateResourcesPath, "{project_id}", updateResourceGroupClient.ProjectID)
			batchCreateResourcesPath = strings.ReplaceAll(batchCreateResourcesPath, "{id}", resourceGroupId)

			batchCreateResourcesOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			batchCreateResourcesOpt.JSONBody = utils.RemoveNil(buildBatchCreateDeleteResourcesBodyParams(newResources))
			_, err = updateResourceGroupClient.Request("POST", batchCreateResourcesPath, &batchCreateResourcesOpt)
			if err != nil {
				return diag.Errorf("error adding resources to resource group: %s", err)
			}
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   resourceGroupId,
			ResourceType: "CES-resourceGroup",
			RegionId:     region,
			ProjectId:    updateResourceGroupClient.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceResourceGroupRead(ctx, d, meta)
}

func buildUpdateResourceGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_name": utils.ValueIgnoreEmpty(d.Get("name")),
		"tags":       utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
	}
	return bodyParams
}

func resourceResourceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteResourceGroup: Delete an existing CES resource group
	var (
		deleteResourceGroupHttpUrl = "v2/{project_id}/resource-groups/batch-delete"
		deleteResourceGroupProduct = "ces"
	)
	deleteResourceGroupClient, err := cfg.NewServiceClient(deleteResourceGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	deleteResourceGroupPath := deleteResourceGroupClient.Endpoint + deleteResourceGroupHttpUrl
	deleteResourceGroupPath = strings.ReplaceAll(deleteResourceGroupPath, "{project_id}", deleteResourceGroupClient.ProjectID)

	deleteResourceGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	deleteResourceGroupOpt.JSONBody = utils.RemoveNil(buildDeleteResourceGroupBodyParams(d))
	_, err = deleteResourceGroupClient.Request("POST", deleteResourceGroupPath, &deleteResourceGroupOpt)
	if err != nil {
		return diag.Errorf("error deleting CES resource group: %s", err)
	}

	return nil
}

func buildDeleteResourceGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_ids": []string{d.Id()},
	}
	return bodyParams
}
