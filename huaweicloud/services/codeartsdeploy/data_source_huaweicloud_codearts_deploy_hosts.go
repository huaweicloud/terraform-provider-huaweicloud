package codeartsdeploy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsDeploy GET /v1/resources/host-groups/{group_id}/hosts
func DataSourceCodeartsDeployHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeartsDeployHostsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the group ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of host.`,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the environment ID.`,
			},
			"as_proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether the host is proxy or not.`,
			},
			"hosts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the host list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the host ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the host name.`,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the IP address.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the SSH port.`,
						},
						"os_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the operating system.`,
						},
						"username": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the username.`,
						},
						"trusted_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the trusted type.`,
						},
						"as_proxy": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the host is an agent host.`,
						},
						"proxy_host_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the agent ID.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the create time.`,
						},
						"lastest_connection_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last connection time.`,
						},
						"connection_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the connection status.`,
						},
						"permission": {
							Type:        schema.TypeList,
							Elem:        deployHostPermissionSchema(),
							Computed:    true,
							Description: `Indicates the permission.`,
						},
						"owner_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the owner ID.`,
						},
						"owner_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the owner name.`,
						},
						"env_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of environments.`,
						},
						"import_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the import status.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeartsDeployHostsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	getHttpUrl := "v1/resources/host-groups/{group_id}/hosts"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{group_id}", d.Get("group_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pageSize is `10`
	getPath += fmt.Sprintf("?page_size=%v", pageSize)
	getPath += buildCodeartsDeployHostsQueryParams(d)
	pageIndex := 1

	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&page_index=%d", pageIndex)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving hosts: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		hosts := utils.PathSearch("result", getRespBody, make([]interface{}, 0)).([]interface{})
		for _, host := range hosts {
			rst = append(rst, map[string]interface{}{
				"id":                    utils.PathSearch("uuid", host, nil),
				"name":                  utils.PathSearch("host_name", host, nil),
				"ip_address":            utils.PathSearch("ip", host, nil),
				"port":                  utils.PathSearch("port", host, nil),
				"os_type":               utils.PathSearch("os", host, nil),
				"username":              utils.PathSearch("authorization.username", host, nil),
				"trusted_type":          utils.PathSearch("authorization.trusted_type", host, nil),
				"as_proxy":              utils.PathSearch("as_proxy", host, nil),
				"proxy_host_id":         utils.PathSearch("proxy_host_id", host, nil),
				"created_at":            utils.PathSearch("create_time", host, nil),
				"lastest_connection_at": utils.PathSearch("lastest_connection_time", host, nil),
				"connection_status":     utils.PathSearch("connection_status", host, nil),
				"owner_id":              utils.PathSearch("owner_id", host, nil),
				"owner_name":            utils.PathSearch("owner_name", host, nil),
				"permission":            flattenDeployHostPermission(host),
				"env_count":             utils.PathSearch("env_count", host, nil),
				"import_status":         utils.PathSearch("import_status", host, nil),
			})
		}

		total := utils.PathSearch("total", getRespBody, float64(0)).(float64)
		if pageSize*(pageIndex-1)+len(hosts) >= int(total) {
			break
		}
		pageIndex++
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("hosts", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCodeartsDeployHostsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&key_field=%v", res, v)
	}
	if v, ok := d.GetOk("environment_id"); ok {
		res = fmt.Sprintf("%s&environment_id=%v", res, v)
	}
	if v, ok := d.GetOk("as_proxy"); ok {
		res = fmt.Sprintf("%s&as_proxy=%v", res, v)
	}

	return res
}
