package swrenterprise

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var enterpriseNamespaceNonUpdatableParams = []string{
	"instance_id", "namespace_name",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/namespaces
// @API SWR GET /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}
// @API SWR PUT /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}
// @API SWR DELETE /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}
// @API SWR POST /v2/{project_id}/{resource_type}/{resource_id}/{sub_resource_type}/{sub_resource_id}/tags/create
// @API SWR DELETE /v2/{project_id}/{resource_type}/{resource_id}/{sub_resource_type}/{sub_resource_id}/tags/delete
// @API SWR GET /v2/{project_id}/{resource_type}/{resource_id}/{sub_resource_type}/{sub_resource_id}/tags
func ResourceSwrEnterpriseNamespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseNamespaceCreate,
		UpdateContext: resourceSwrEnterpriseNamespaceUpdate,
		ReadContext:   resourceSwrEnterpriseNamespaceRead,
		DeleteContext: resourceSwrEnterpriseNamespaceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSwrEnterpriseNamespaceImportStateFunc,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(enterpriseNamespaceNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace name.`,
			},
			"metadata": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the metadata.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
							Description:  `Specifies whether the namespace is public.`,
						},
					},
				},
			},
			"tags": common.TagsSchema(`Specifies the key/value pairs to associate with the namespace.`),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"namespace_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the namespace ID.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last update time.`,
			},
			"repo_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the repo count of the namespace.`,
			},
		},
	}
}

func resourceSwrEnterpriseNamespaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateSwrEnterpriseNamespaceBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR namespace: %s", err)
	}

	d.SetId(d.Get("instance_id").(string) + "/" + d.Get("name").(string))

	if v, ok := d.GetOk("tags"); ok {
		err := addSwrEnterpriseInstanceNamespaceTags(client, d, v.(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSwrEnterpriseNamespaceRead(ctx, d, meta)
}

func buildCreateSwrEnterpriseNamespaceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"namespace_name": d.Get("name"),
		"metadata":       buildSwrEnterpriseNamespaceMetadataBodyParams(d),
	}

	return bodyParams
}

func buildSwrEnterpriseNamespaceMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	paramsList := d.Get("metadata").([]interface{})
	if params, ok := paramsList[0].(map[string]interface{}); ok {
		return map[string]interface{}{
			"public": params["public"],
		}
	}

	return nil
}

func resourceSwrEnterpriseNamespaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{namespace_name}", d.Get("name").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR namespace")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("metadata", flattenSwrEnterpriseNamespaceMetadata(getRespBody)),
		d.Set("namespace_id", utils.PathSearch("namespace_id", getRespBody, nil)),
		d.Set("repo_count", utils.PathSearch("repo_count", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getRespBody, nil)),
	)

	if resourceTags, err := getSwrEnterpriseNamespaceTags(client, d); err == nil {
		mErr = multierror.Append(mErr,
			d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", resourceTags, make([]interface{}, 0)))),
		)
	} else {
		log.Printf("[WARN] error fetching tags of namespace: %s", err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSwrEnterpriseNamespaceMetadata(resp interface{}) []interface{} {
	params := utils.PathSearch("metadata", resp, nil)
	if params == nil {
		return nil
	}

	rst := map[string]interface{}{
		"public": utils.PathSearch("public", params, nil),
	}

	return []interface{}{rst}
}

func getSwrEnterpriseNamespaceTags(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/{resource_type}/{resource_id}/{sub_resource_type}/{sub_resource_id}/tags"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{resource_type}", "instances")
	getPath = strings.ReplaceAll(getPath, "{resource_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{sub_resource_type}", "namespaces")
	getPath = strings.ReplaceAll(getPath, "{sub_resource_id}", d.Get("name").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func resourceSwrEnterpriseNamespaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	if d.HasChanges("metadata") {
		updateHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
		updatePath = strings.ReplaceAll(updatePath, "{namespace_name}", d.Get("name").(string))
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateSwrEnterpriseNamespaceBodyParams(d),
		}

		_, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating SWR instance namespace: %s", err)
		}
	}

	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})
		// remove old tags
		if len(oMap) > 0 {
			err := deleteSwrEnterpriseInstanceNamespaceTags(client, d, oMap)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// add new tags
		if len(nMap) > 0 {
			err := addSwrEnterpriseInstanceNamespaceTags(client, d, nMap)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceSwrEnterpriseNamespaceRead(ctx, d, meta)
}

func buildUpdateSwrEnterpriseNamespaceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": buildSwrEnterpriseNamespaceMetadataBodyParams(d),
	}

	return bodyParams
}

func addSwrEnterpriseInstanceNamespaceTags(client *golangsdk.ServiceClient, d *schema.ResourceData, tags map[string]interface{}) error {
	httpUrl := "v2/{project_id}/{resource_type}/{resource_id}/{sub_resource_type}/{sub_resource_id}/tags/create"
	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{project_id}", client.ProjectID)
	addPath = strings.ReplaceAll(addPath, "{resource_type}", "instances")
	addPath = strings.ReplaceAll(addPath, "{resource_id}", d.Get("instance_id").(string))
	addPath = strings.ReplaceAll(addPath, "{sub_resource_type}", "namespaces")
	addPath = strings.ReplaceAll(addPath, "{sub_resource_id}", d.Get("name").(string))

	addOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: buildUpdateSwrEnterpriseInstanceNamespaceTagsBodyParams(tags),
	}

	_, err := client.Request("POST", addPath, &addOpt)
	if err != nil {
		return fmt.Errorf("error adding SWR enterprise instance namespace tags: %s", err)
	}

	return nil
}

func deleteSwrEnterpriseInstanceNamespaceTags(client *golangsdk.ServiceClient, d *schema.ResourceData, tags map[string]interface{}) error {
	httpUrl := "v2/{project_id}/{resource_type}/{resource_id}/{sub_resource_type}/{sub_resource_id}/tags/delete"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{resource_type}", "instances")
	deletePath = strings.ReplaceAll(deletePath, "{resource_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{sub_resource_type}", "namespaces")
	deletePath = strings.ReplaceAll(deletePath, "{sub_resource_id}", d.Get("name").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: buildUpdateSwrEnterpriseInstanceNamespaceTagsBodyParams(tags),
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return fmt.Errorf("error deleting SWR enterprise instance namespace tags: %s", err)
	}

	return nil
}

func buildUpdateSwrEnterpriseInstanceNamespaceTagsBodyParams(tags map[string]interface{}) interface{} {
	bodyParams := map[string]interface{}{
		"tags": utils.ExpandResourceTagsMap(tags),
	}
	return bodyParams
}

func resourceSwrEnterpriseNamespaceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{namespace_name}", d.Get("name").(string))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR namespace")
	}

	return nil
}

func resourceSwrEnterpriseNamespaceImportStateFunc(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<name>', but got '%s'", d.Id())
	}

	if err := d.Set("instance_id", parts[0]); err != nil {
		return nil, fmt.Errorf("error saving instance ID: %s", err)
	}

	if err := d.Set("name", parts[1]); err != nil {
		return nil, fmt.Errorf("error saving namespace name: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
