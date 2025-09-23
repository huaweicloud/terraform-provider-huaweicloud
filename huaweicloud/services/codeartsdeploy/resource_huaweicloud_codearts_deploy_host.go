package codeartsdeploy

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

// @API CodeArtsDeploy POST /v1/resources/host-groups/{group_id}/hosts
// @API CodeArtsDeploy PUT /v1/resources/host-groups/{group_id}/hosts/{host_id}
// @API CodeArtsDeploy GET /v1/resources/host-groups/{group_id}/hosts/{host_id}
// @API CodeArtsDeploy DELETE /v1/resources/host-groups/{group_id}/hosts/{host_id}
func ResourceDeployHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeployHostCreate,
		UpdateContext: resourceDeployHostUpdate,
		ReadContext:   resourceDeployHostRead,
		DeleteContext: resourceDeployHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDeployHostImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the CodeArts deploy group ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the host name.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the IP address.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the SSH port.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the operating system.`,
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the username.`,
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"password", "private_key"},
				Description:  `Specifies the password.`,
			},
			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `Specifies the private key.`,
			},
			"as_proxy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether the host is an agent host.`,
			},
			"proxy_host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the agent ID.`,
			},
			"install_icagent": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `Specifies whether to enable Application Operations Management (AOM) for free to provide
metric monitoring, log query and alarm functions.`,
			},
			"sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `Specifies whether to synchronize the password of the current host to the hosts with the
same IP address, username and port number in other group in the same project.`,
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
			"lastest_connection_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last connection time.`,
			},
			"connection_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection status.`,
			},
			"permission": {
				Type:     schema.TypeList,
				Elem:     deployHostPermissionSchema(),
				Computed: true,
			},
		},
	}
}

func deployHostPermissionSchema() *schema.Resource {
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
			"can_copy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the permission to copy hosts.`,
			},
		},
	}
	return &sc
}

func resourceDeployHostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/resources/host-groups/{group_id}/hosts"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{group_id}", d.Get("group_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: buildCreateDeployHostBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy host: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	hostId := utils.PathSearch("id", createRespBody, "").(string)
	if hostId == "" {
		return diag.Errorf("unable to find the CodeArts deploy host ID from the API response")
	}
	d.SetId(hostId)

	if d.Get("sync").(bool) {
		if err := updateDeployHost(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDeployHostRead(ctx, d, meta)
}

func buildCreateDeployHostBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"host_name":       d.Get("name"),
		"ip":              d.Get("ip_address"),
		"port":            d.Get("port"),
		"os":              d.Get("os_type"),
		"as_proxy":        d.Get("as_proxy"),
		"install_icagent": d.Get("install_icagent"),
		"authorization":   buildAuthorizationBodyParam(d),
	}

	if v, ok := d.GetOk("proxy_host_id"); ok {
		bodyParams["proxy_host_id"] = v
	}

	return bodyParams
}

func buildAuthorizationBodyParam(d *schema.ResourceData) map[string]interface{} {
	// Use password authentication
	if v, ok := d.GetOk("password"); ok {
		return map[string]interface{}{
			"username":     d.Get("username"),
			"password":     v,
			"trusted_type": 0,
		}
	}
	// Use key authentication
	return map[string]interface{}{
		"username":     d.Get("username"),
		"private_key":  d.Get("private_key"),
		"trusted_type": 1,
	}
}

func resourceDeployHostRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/resources/host-groups/{group_id}/hosts/{host_id}"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{group_id}", d.Get("group_id").(string))
	getPath = strings.ReplaceAll(getPath, "{host_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts deploy host")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resultRespBody := utils.PathSearch("result", getRespBody, nil)
	if resultRespBody == nil {
		return diag.Errorf("error retrieving CodeArts deploy host: result is not found in API response")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("group_id", utils.PathSearch("group_id", resultRespBody, nil)),
		d.Set("name", utils.PathSearch("host_name", resultRespBody, nil)),
		d.Set("ip_address", utils.PathSearch("ip", resultRespBody, nil)),
		d.Set("port", utils.PathSearch("port", resultRespBody, nil)),
		d.Set("os_type", utils.PathSearch("os", resultRespBody, nil)),
		d.Set("as_proxy", utils.PathSearch("as_proxy", resultRespBody, nil)),
		d.Set("proxy_host_id", utils.PathSearch("proxy_host_id", resultRespBody, nil)),
		d.Set("username", utils.PathSearch("authorization.username", resultRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", resultRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", resultRespBody, nil)),
		d.Set("lastest_connection_at", utils.PathSearch("lastest_connection_time", resultRespBody,
			nil)),
		d.Set("connection_status", utils.PathSearch("connection_status", resultRespBody,
			nil)),
		d.Set("permission", flattenDeployHostPermission(resultRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDeployHostPermission(resp interface{}) []interface{} {
	curJson := utils.PathSearch("permission", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"can_view":     utils.PathSearch("can_view", curJson, nil),
			"can_edit":     utils.PathSearch("can_edit", curJson, nil),
			"can_delete":   utils.PathSearch("can_delete", curJson, nil),
			"can_add_host": utils.PathSearch("can_add_host", curJson, nil),
			"can_copy":     utils.PathSearch("can_copy", curJson, nil),
		},
	}
}

func resourceDeployHostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("codearts_deploy", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	if err := updateDeployHost(client, d); err != nil {
		return diag.FromErr(err)
	}
	return resourceDeployHostRead(ctx, d, meta)
}

func updateDeployHost(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updatePath := client.Endpoint + "v1/resources/host-groups/{group_id}/hosts/{host_id}"
	updatePath = strings.ReplaceAll(updatePath, "{group_id}", d.Get("group_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{host_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: buildUpdateDeployHostBodyParams(d),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating CodeArts deploy host: %s", err)
	}

	return nil
}

func buildUpdateDeployHostBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"host_name":       d.Get("name"),
		"ip":              d.Get("ip_address"),
		"port":            d.Get("port"),
		"as_proxy":        d.Get("as_proxy"),
		"install_icagent": d.Get("install_icagent"),
		"sync":            d.Get("sync"),
		"authorization":   buildAuthorizationBodyParam(d),
	}

	if v, ok := d.GetOk("proxy_host_id"); ok {
		bodyParams["proxy_host_id"] = v
	}

	return bodyParams
}

func resourceDeployHostDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/resources/host-groups/{group_id}/hosts/{host_id}"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{group_id}", d.Get("group_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{host_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts deploy host")
	}

	return nil
}

func resourceDeployHostImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID,"+
			" must be <group_id>/<id> but got %s", d.Id())
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("group_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import CodeArts deploy host, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
