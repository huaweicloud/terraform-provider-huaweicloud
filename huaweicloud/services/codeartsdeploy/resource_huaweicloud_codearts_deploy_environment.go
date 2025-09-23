package codeartsdeploy

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsDeploy POST /v1/applications/{application_id}/environments
// @API CodeArtsDeploy PUT /v1/applications/{application_id}/environments/{environment_id}
// @API CodeArtsDeploy GET /v1/applications/{application_id}/environments/{environment_id}
// @API CodeArtsDeploy DELETE /v1/applications/{application_id}/environments/{environment_id}
// @API CodeArtsDeploy POST /v1/applications/{application_id}/environments/{environment_id}/hosts/import
// @API CodeArtsDeploy DELETE /v1/applications/{application_id}/environments/{environment_id}/{host_id}
// @API CodeArtsDeploy GET /v1/applications/{application_id}/environments/{environment_id}/hosts
// @API CodeArtsDeploy GET /v2/applications/{application_id}/environments/{environment_id}/permissions
func ResourceDeployEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeployEnvironmentCreate,
		UpdateContext: resourceDeployEnvironmentUpdate,
		ReadContext:   resourceDeployEnvironmentRead,
		DeleteContext: resourceDeployEnvironmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDeployEnvironmentImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the project ID for CodeArts service.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the application ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the environment name.`,
			},
			"deploy_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the deployment type.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the operating system.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description.`,
			},
			"hosts": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the hosts list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the cluster group ID.`,
						},
						"host_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the host ID.`,
						},
						"host_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the host name.`,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the IP address.`,
						},
						"connection_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the connection status.`,
						},
					},
				},
			},
			"proxies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the proxy hosts list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the cluster group ID.`,
						},
						"host_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the host ID.`,
						},
						"host_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the host name.`,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the IP address.`,
						},
						"connection_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the connection status.`,
						},
					},
				},
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of host instances in the environment.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the create time.`,
			},
			"created_by": {
				Type:        schema.TypeList,
				Elem:        deployEnvironmentUserSchema(),
				Computed:    true,
				Description: `Indicates the creator information.`,
			},
			"permission": {
				Type:        schema.TypeList,
				Elem:        deployEnvironmentPermissionSchema(),
				Computed:    true,
				Description: `Indicates the user permission.`,
			},
			"permission_matrix": {
				Type:        schema.TypeList,
				Elem:        deployEnvironmentPermissionMatrixSchema(),
				Computed:    true,
				Description: `Indicates the permission matrix.`,
			},
		},
	}
}

func deployEnvironmentUserSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the user ID.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the user name.`,
			},
		},
	}
	return &sc
}

func deployEnvironmentPermissionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"can_view": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the view environment.`,
			},
			"can_edit": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the permission to edit environments.`,
			},
			"can_delete": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the permission to delete environments.`,
			},
			"can_manage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the permission to edit the environment permission matrix.`,
			},
			"can_deploy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the deploy permission.`,
			},
		},
	}
	return &sc
}

func deployEnvironmentPermissionMatrixSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"permission_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the permission ID.`,
			},
			"role_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the role ID.`,
			},
			"role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the role name.`,
			},
			"role_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the role type.`,
			},
			"can_view": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the view environment.`,
			},
			"can_edit": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the permission to edit environments.`,
			},
			"can_delete": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the permission to delete environments.`,
			},
			"can_manage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the permission to edit the environment permission matrix.`,
			},
			"can_deploy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the deploy permission.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the create time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the update time.`,
			},
		},
	}
	return &sc
}

func resourceDeployEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	httpUrl := "v1/applications/{application_id}/environments"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{application_id}", d.Get("application_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: utils.RemoveNil(buildCreateDeployEnvironmentBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy environment: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts deploy environment ID from the API response")
	}

	d.SetId(id)

	if v, ok := d.GetOk("hosts"); ok {
		if err := importEnvironmentHosts(client, d, v.(*schema.Set).List()); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDeployEnvironmentRead(ctx, d, meta)
}

func buildCreateDeployEnvironmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"project_id":  d.Get("project_id"),
		"os":          d.Get("os_type"),
		"deploy_type": d.Get("deploy_type"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func importEnvironmentHosts(client *golangsdk.ServiceClient, d *schema.ResourceData, importList []interface{}) error {
	for _, v := range importList {
		params := v.(map[string]interface{})
		hostId := params["host_id"]

		httpUrl := "v1/applications/{application_id}/environments/{environment_id}/hosts/import"
		importPath := client.Endpoint + httpUrl
		importPath = strings.ReplaceAll(importPath, "{application_id}", d.Get("application_id").(string))
		importPath = strings.ReplaceAll(importPath, "{environment_id}", d.Id())
		importOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf-8",
			},
			// there is no difference between importing one by one and importing all in a time
			JSONBody: map[string]interface{}{
				"group_id": params["group_id"],
				"host_ids": []interface{}{hostId},
			},
		}

		_, err := client.Request("POST", importPath, &importOpt)
		if err != nil {
			return fmt.Errorf("error importing host(%s) to environment: %s", hostId, err)
		}
	}

	return nil
}

func removeEnvironmentHosts(client *golangsdk.ServiceClient, d *schema.ResourceData, removeList []interface{}) error {
	for _, v := range removeList {
		params := v.(map[string]interface{})
		hostId := params["host_id"].(string)

		httpUrl := "v1/applications/{application_id}/environments/{environment_id}/{host_id}"
		removePath := client.Endpoint + httpUrl
		removePath = strings.ReplaceAll(removePath, "{application_id}", d.Get("application_id").(string))
		removePath = strings.ReplaceAll(removePath, "{environment_id}", d.Id())
		removePath = strings.ReplaceAll(removePath, "{host_id}", hostId)
		removeOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf-8",
			},
		}

		_, err := client.Request("DELETE", removePath, &removeOpt)
		if err != nil {
			return fmt.Errorf("error removing environment host(%s): %s", hostId, err)
		}
	}

	return nil
}

func resourceDeployEnvironmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	httpUrl := "v1/applications/{application_id}/environments/{environment_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{application_id}", d.Get("application_id").(string))
	getPath = strings.ReplaceAll(getPath, "{environment_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts deploy environment")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	environment := utils.PathSearch("result", getRespBody, nil)
	if environment == nil {
		return diag.Errorf("error retrieving CodeArts deploy environment: result is not found in API response")
	}

	hosts, proxies, err := getDeployEnvironmentHosts(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	permissionMatrix, err := getDeployEnvironmentPermissionMatrix(client, d.Get("application_id").(string), d.Id())
	if err != nil {
		log.Printf("[WARN] failed to retrieve environment permission matrix: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", environment, nil)),
		d.Set("os_type", utils.PathSearch("os", environment, nil)),
		d.Set("description", utils.PathSearch("description", environment, nil)),
		d.Set("deploy_type", utils.PathSearch("deploy_type", environment, nil)),
		d.Set("created_at", utils.PathSearch("created_time", environment, nil)),
		d.Set("instance_count", utils.PathSearch("instance_count", environment, nil)),
		d.Set("created_by", flattenDeployEnvironmentCreatedBy(environment)),
		d.Set("permission", flattenDeployEnvironmentPermission(environment)),
		d.Set("permission_matrix", flattenDeployEnvironmentPermissionMatrix(permissionMatrix)),
		d.Set("hosts", hosts),
		d.Set("proxies", proxies),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getDeployEnvironmentPermissionMatrix(client *golangsdk.ServiceClient, appId, id string) (interface{}, error) {
	httpUrl := "v2/applications/{application_id}/environments/{environment_id}/permissions"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{application_id}", appId)
	getPath = strings.ReplaceAll(getPath, "{environment_id}", id)
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

	permissionMatrix := getRespBody.([]interface{})
	if len(permissionMatrix) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte("error retrieving CodeArts deploy environment permission matrix, empty list"),
			},
		}
	}
	return permissionMatrix, nil
}

func getDeployEnvironmentHosts(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, interface{}, error) {
	getHttpUrl := "v1/applications/{application_id}/environments/{environment_id}/hosts"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{application_id}", d.Get("application_id").(string))
	getPath = strings.ReplaceAll(getPath, "{environment_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pageSize is `10`
	getPath += fmt.Sprintf("?page_size=%v", pageSize)
	pageIndex := 1

	hosts := make([]interface{}, 0)
	proxies := make([]interface{}, 0)

	for {
		currentPath := getPath + fmt.Sprintf("&page_index=%d", pageIndex)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return nil, nil, fmt.Errorf("error retrieving hosts: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, nil, fmt.Errorf("error flatten response: %s", err)
		}

		results := utils.PathSearch("result", getRespBody, make([]interface{}, 0)).([]interface{})
		for _, result := range results {
			asProxy := utils.PathSearch("as_proxy", result, false).(bool)
			if asProxy {
				proxies = append(proxies, flattenDeployEnvironmentHosts(result))
				continue
			}
			hosts = append(hosts, flattenDeployEnvironmentHosts(result))
		}

		total := utils.PathSearch("total", getRespBody, float64(0)).(float64)
		if pageSize*(pageIndex-1)+len(results) >= int(total) {
			break
		}
		pageIndex++
	}

	return hosts, proxies, nil
}

func flattenDeployEnvironmentHosts(resp interface{}) interface{} {
	return map[string]interface{}{
		"group_id":          utils.PathSearch("group_id", resp, nil),
		"host_id":           utils.PathSearch("host_id", resp, nil),
		"host_name":         utils.PathSearch("host_name", resp, nil),
		"ip_address":        utils.PathSearch("ip", resp, nil),
		"connection_status": utils.PathSearch("connection_status", resp, nil),
	}
}

func flattenDeployEnvironmentCreatedBy(resp interface{}) []interface{} {
	curJson := utils.PathSearch("created_by", resp, nil)
	if curJson == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"user_id":   utils.PathSearch("user_id", curJson, nil),
			"user_name": utils.PathSearch("user_name", curJson, nil),
		},
	}
}

func flattenDeployEnvironmentPermission(resp interface{}) []interface{} {
	curJson := utils.PathSearch("permission", resp, nil)
	if curJson == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"can_view":   utils.PathSearch("can_view", curJson, nil),
			"can_edit":   utils.PathSearch("can_edit", curJson, nil),
			"can_delete": utils.PathSearch("can_delete", curJson, nil),
			"can_manage": utils.PathSearch("can_manage", curJson, nil),
			"can_deploy": utils.PathSearch("can_deploy", curJson, nil),
		},
	}
}

func flattenDeployEnvironmentPermissionMatrix(respBody interface{}) []interface{} {
	if resp, isList := respBody.([]interface{}); isList {
		rst := make([]interface{}, 0, len(resp))
		for _, v := range resp {
			rst = append(rst, map[string]interface{}{
				"permission_id": utils.PathSearch("id", v, nil),
				"role_id":       utils.PathSearch("role_id", v, nil),
				"role_name":     utils.PathSearch("name", v, nil),
				"role_type":     utils.PathSearch("role_type", v, nil),
				"can_view":      utils.PathSearch("can_view", v, nil),
				"can_edit":      utils.PathSearch("can_edit", v, nil),
				"can_delete":    utils.PathSearch("can_delete", v, nil),
				"can_manage":    utils.PathSearch("can_manage", v, nil),
				"can_deploy":    utils.PathSearch("can_deploy", v, nil),
				"created_at":    utils.PathSearch("create_time", v, nil),
				"updated_at":    utils.PathSearch("update_time", v, nil),
			})
		}
		return rst
	}

	return nil
}

func resourceDeployEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	if d.HasChanges("name", "description") {
		httpUrl := "v1/applications/{application_id}/environments/{environment_id}"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{application_id}", d.Get("application_id").(string))
		updatePath = strings.ReplaceAll(updatePath, "{environment_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf-8",
			},
			JSONBody: buildUpdateDeployEnvironmentBodyParams(d),
		}
		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating CodeArts deploy environment: %s", err)
		}
	}

	if d.HasChange("hosts") {
		o, n := d.GetChange("hosts")
		os, ns := o.(*schema.Set), n.(*schema.Set)
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()

		// remove hosts first
		if len(remove) > 0 {
			if err := removeEnvironmentHosts(client, d, remove); err != nil {
				return diag.FromErr(err)
			}
		}

		// add hosts
		if len(add) > 0 {
			if err := importEnvironmentHosts(client, d, add); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceDeployEnvironmentRead(ctx, d, meta)
}

func buildUpdateDeployEnvironmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}
}

func resourceDeployEnvironmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	httpUrl := "v1/applications/{application_id}/environments/{environment_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{application_id}", d.Get("application_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{environment_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts deploy environment")
	}

	return nil
}

func resourceDeployEnvironmentImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<project_id>/<application_id>/<id>', but got '%s'",
			d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("project_id", parts[0]),
		d.Set("application_id", parts[1]),
	)
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
