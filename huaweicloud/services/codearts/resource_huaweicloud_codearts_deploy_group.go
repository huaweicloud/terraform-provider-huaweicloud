// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CodeArts
// ---------------------------------------------------------------

package codearts

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	groupNotFound = "Deploy.00021423"
)

// @API CodeArtsDeploy POST /v1/resources/host-groups
// @API CodeArtsDeploy PUT /v2/host-groups/{group_id}
// @API CodeArtsDeploy GET /v1/resources/host-groups/{group_id}
// @API CodeArtsDeploy DELETE /v2/host-groups/{group_id}
func ResourceDeployGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeployGroupCreate,
		UpdateContext: resourceDeployGroupUpdate,
		ReadContext:   resourceDeployGroupRead,
		DeleteContext: resourceDeployGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDeployGroupImportState,
		},

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
				Description: `Specifies the group name.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the project ID.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the operating system.`,
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource pool ID.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description.`,
			},
			"is_proxy_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				ForceNew:    true,
				Description: `Specifies whether the host is in agent access mode.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"created_by": {
				Type:     schema.TypeList,
				Elem:     deployGroupUserSchema(),
				Computed: true,
			},
			"permission": {
				Type:     schema.TypeList,
				Elem:     deployGroupPermissionSchema(),
				Computed: true,
			},
		},
	}
}

func deployGroupUserSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user ID.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user name.`,
			},
		},
	}
	return &sc
}

func deployGroupPermissionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"can_view": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the view permission.`,
			},
			"can_edit": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the edit permission.`,
			},
			"can_delete": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the deletion permission.`,
			},
			"can_add_host": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the permission to add hosts.`,
			},
			"can_manage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the management permission.`,
			},
		},
	}
	return &sc
}

func resourceDeployGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/resources/host-groups"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: buildCreateDeployGroupBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy group: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody, groupNotFound); err != nil {
		return diag.Errorf("error creating CodeArts deploy group: %s", err)
	}

	id, err := jmespath.Search("id", createRespBody)
	if err != nil || id == nil {
		return diag.Errorf("error creating CodeArts deploy group: ID is not found in API response")
	}

	d.SetId(id.(string))

	return resourceDeployGroupRead(ctx, d, meta)
}

func buildCreateDeployGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":             d.Get("name"),
		"project_id":       d.Get("project_id"),
		"os":               d.Get("os_type"),
		"is_proxy_mode":    d.Get("is_proxy_mode"),
		"slave_cluster_id": d.Get("resource_pool_id"),
		"description":      d.Get("description"),
	}
}

func resourceDeployGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/resources/host-groups/{group_id}"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{group_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving CodeArts deploy group: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(getRespBody, groupNotFound); err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts deploy group")
	}

	resultRespBody := utils.PathSearch("result", getRespBody, nil)
	if resultRespBody == nil {
		return diag.Errorf("error retrieving CodeArts deploy group: result is not found in API response")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", resultRespBody, nil)),
		d.Set("os_type", utils.PathSearch("os", resultRespBody, nil)),
		d.Set("resource_pool_id", utils.PathSearch("slave_cluster_id", resultRespBody, nil)),
		d.Set("description", utils.PathSearch("description", resultRespBody, nil)),
		d.Set("is_proxy_mode", utils.PathSearch("is_proxy_mode", resultRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_time", resultRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_time", resultRespBody, nil)),
		d.Set("created_by", flattenDeployGroupCreatedBy(resultRespBody)),
		d.Set("permission", flattenDeployGroupPermission(resultRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDeployGroupCreatedBy(resp interface{}) []interface{} {
	curJson, err := jmespath.Search("created_by", resp)
	if err != nil {
		log.Printf("[ERROR] error flatten created_by, cause this field is not found in API response")
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"user_id":   utils.PathSearch("user_id", curJson, nil),
			"user_name": utils.PathSearch("user_name", curJson, nil),
		},
	}
}

func flattenDeployGroupPermission(resp interface{}) []interface{} {
	curJson, err := jmespath.Search("permission", resp)
	if err != nil {
		log.Printf("[ERROR] error flatten permission, cause this field is not found in API response")
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"can_view":     utils.PathSearch("can_view", curJson, nil),
			"can_edit":     utils.PathSearch("can_edit", curJson, nil),
			"can_delete":   utils.PathSearch("can_delete", curJson, nil),
			"can_add_host": utils.PathSearch("can_add_host", curJson, nil),
			"can_manage":   utils.PathSearch("can_manage", curJson, nil),
		},
	}
}

func resourceDeployGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/host-groups/{group_id}"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{group_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: buildUpdateDeployGroupBodyParams(d),
	}
	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating CodeArts deploy group: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(updateRespBody, groupNotFound); err != nil {
		return diag.Errorf("error updating CodeArts deploy group: %s", err)
	}

	return resourceDeployGroupRead(ctx, d, meta)
}

func buildUpdateDeployGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":             d.Get("name"),
		"description":      d.Get("description"),
		"slave_cluster_id": d.Get("resource_pool_id"),
	}
}

func resourceDeployGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/host-groups/{group_id}"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{group_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting CodeArts deploy group: %s", err)
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(deleteRespBody, groupNotFound); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts deploy group")
	}

	return nil
}

// checkResponseError use to check whether the CodeArts deploy API response body contains error code.
// Parameter 'respBody' is the response body and 'notFoundCode' is the error code when the resource is not found.
// An example of an error response body is as follows: {"error_code": "XXX", "error_msg": "XXX", "status": "XXX"}
func checkResponseError(respBody interface{}, notFoundCode string) error {
	errorCode := utils.PathSearch("error_code", respBody, "")
	if errorCode == "" {
		return nil
	}

	errorMsg := utils.PathSearch("error_msg", respBody, "")
	err := fmt.Errorf("error code: %s, error message: %s", errorCode, errorMsg)
	if errorCode != notFoundCode {
		return err
	}

	return golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte(err.Error()),
		},
	}
}

func resourceDeployGroupImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<project_id>/<id>', but got '%s'",
			d.Id())
	}

	d.SetId(parts[1])
	if err := d.Set("project_id", parts[0]); err != nil {
		return nil, fmt.Errorf("error saving project ID: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
