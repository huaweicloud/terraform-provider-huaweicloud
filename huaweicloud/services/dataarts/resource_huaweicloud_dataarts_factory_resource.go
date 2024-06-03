package dataarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio POST /v1/{project_id}/resources
// @API DataArtsStudio DELETE /v1/{project_id}/resources/{resource_id}
// @API DataArtsStudio GET /v1/{project_id}/resources/{resource_id}
// @API DataArtsStudio PUT /v1/{project_id}/resources/{resource_id}
func ResourceFactoryResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFactoryResourceCreate,
		ReadContext:   resourceFactoryResourceRead,
		UpdateContext: resourceFactoryResourceUpdate,
		DeleteContext: resourceFactoryResourceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceFactoryResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"directory": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"depend_packages": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"location": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceFactoryResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	createResourceHttpUrl := "v1/{project_id}/resources"
	createResourceProduct := "dataarts-dlf"

	createResourceClient, err := cfg.NewServiceClient(createResourceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	createResourcePath := createResourceClient.Endpoint + createResourceHttpUrl
	createResourcePath = strings.ReplaceAll(createResourcePath, "{project_id}", createResourceClient.ProjectID)

	createResourceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}
	createResourceOpt.JSONBody = utils.RemoveNil(buildCreateOrUpdateResourceBodyParams(d))
	createResourceResp, err := createResourceClient.Request("POST", createResourcePath, &createResourceOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts Factory resource: %s", err)
	}

	createResourceRespBody, err := utils.FlattenResponse(createResourceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("resourceId", createResourceRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating DataArts Factory resource: %s is not found in API response", "id")
	}
	d.SetId(id.(string))
	return resourceFactoryResourceRead(ctx, d, meta)
}

func buildCreateOrUpdateResourceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":            d.Get("name"),
		"type":            d.Get("type"),
		"location":        d.Get("location"),
		"directory":       d.Get("directory"),
		"depend_packages": utils.ValueIgnoreEmpty(buildDependPackagesParams(d)),
		"desc":            utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func buildDependPackagesParams(d *schema.ResourceData) []map[string]string {
	rawDependPackages := d.Get("depend_packages").([]interface{})
	if len(rawDependPackages) == 0 {
		return nil
	}

	dependPackages := make([]map[string]string, len(rawDependPackages))

	for i := range rawDependPackages {
		item := rawDependPackages[i].(map[string]interface{})
		dependPackage := make(map[string]string)
		dependPackage["type"] = item["type"].(string)
		dependPackage["location"] = item["location"].(string)
		dependPackages[i] = dependPackage
	}
	return dependPackages
}

func resourceFactoryResourceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)

	var mErr *multierror.Error

	getResourceHttpUrl := "v1/{project_id}/resources/{resource_id}"
	getResourceProduct := "dataarts-dlf"

	getResourceClient, err := cfg.NewServiceClient(getResourceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	getResourcePath := getResourceClient.Endpoint + getResourceHttpUrl
	getResourcePath = strings.ReplaceAll(getResourcePath, "{project_id}", getResourceClient.ProjectID)
	getResourcePath = strings.ReplaceAll(getResourcePath, "{resource_id}", d.Id())

	getResourceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	getResourceResp, err := getResourceClient.Request("GET", getResourcePath, &getResourceOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DataArts Factory resource")
	}

	getResourceRespBody, err := utils.FlattenResponse(getResourceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("workspace_id", workspaceID),
		d.Set("name", utils.PathSearch("name", getResourceRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getResourceRespBody, nil)),
		d.Set("location", utils.PathSearch("location", getResourceRespBody, nil)),
		d.Set("directory", utils.PathSearch("directory", getResourceRespBody, nil)),
		d.Set("depend_packages", flattenDependPackagesInGetResourceResponseBody(getResourceRespBody)),
		d.Set("description", utils.PathSearch("desc", getResourceRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDependPackagesInGetResourceResponseBody(getResourceRespBody interface{}) []map[string]interface{} {
	if getResourceRespBody == nil {
		return nil
	}
	curJson := utils.PathSearch("dependPackages", getResourceRespBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	dependPackages := make([]map[string]interface{}, len(curArray))

	for i, val := range curArray {
		dependPackage := make(map[string]interface{})
		dependPackage["type"] = utils.PathSearch("type", val, "")
		dependPackage["location"] = utils.PathSearch("location", val, "")
		dependPackages[i] = dependPackage
	}
	return dependPackages
}

func resourceFactoryResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateResourceHttpUrl := "v1/{project_id}/resources/{resource_id}"
	updateResourceProduct := "dataarts-dlf"

	updateResourceClient, err := cfg.NewServiceClient(updateResourceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	updateResourcePath := updateResourceClient.Endpoint + updateResourceHttpUrl
	updateResourcePath = strings.ReplaceAll(updateResourcePath, "{project_id}", updateResourceClient.ProjectID)
	updateResourcePath = strings.ReplaceAll(updateResourcePath, "{resource_id}", d.Id())

	updateResourceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		OkCodes: []int{
			204,
		},
	}
	updateResourceOpt.JSONBody = utils.RemoveNil(buildCreateOrUpdateResourceBodyParams(d))
	_, err = updateResourceClient.Request("PUT", updateResourcePath, &updateResourceOpt)
	if err != nil {
		return diag.Errorf("error updating DataArts Factory resource: %s", err)
	}

	return resourceFactoryResourceRead(ctx, d, meta)
}

func resourceFactoryResourceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	deleteResourceHttpUrl := "v1/{project_id}/resources/{resource_id}"
	deleteResourceProduct := "dataarts-dlf"

	deleteResourceClient, err := cfg.NewServiceClient(deleteResourceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	deleteResourcePath := deleteResourceClient.Endpoint + deleteResourceHttpUrl
	deleteResourcePath = strings.ReplaceAll(deleteResourcePath, "{project_id}", deleteResourceClient.ProjectID)
	deleteResourcePath = strings.ReplaceAll(deleteResourcePath, "{resource_id}", d.Id())

	deleteResourceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}
	_, err = deleteResourceClient.Request("DELETE", deleteResourcePath, &deleteResourceOpt)
	if err != nil {
		return diag.Errorf("error deleting DataArts Factory resource: %s", err)
	}

	return nil
}

func resourceFactoryResourceImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <workspace_id>/<id>")
	}

	d.Set("workspace_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
