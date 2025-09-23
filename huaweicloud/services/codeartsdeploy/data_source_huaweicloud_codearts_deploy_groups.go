package codeartsdeploy

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const pageSize = 10

// @API CodeArtsDeploy GET /v1/resources/host-groups
func DataSourceCodeartsDeployGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeartsDeployGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the project ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of host cluster.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the operating system. Valid values are **windows**, **linux**.`,
			},
			"is_proxy_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether the host is an agent host.`,
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the customized resource pool ID.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the host cluster list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the host cluster ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the host cluster name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the description of host cluster.`,
						},
						"env_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of environments.`,
						},
						"host_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the the number of hosts in a cluster.`,
						},
						"is_proxy_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates whether the host is an agent host.`,
						},
						"os_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the operating system.`,
						},
						"resource_pool_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the slave cluster ID.`,
						},
						"permission": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the permission list.`,
							Elem:        deployGroupPermissionSchema(),
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator name.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeartsDeployGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	getHttpUrl := "v1/resources/host-groups"
	getPath := client.Endpoint + getHttpUrl
	getPath += buildCodeartsDeployGroupsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pageSize is `10`
	getPath += fmt.Sprintf("&page_size=%v", pageSize)
	pageIndex := 1

	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&page_index=%d", pageIndex)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving groups: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		groups := utils.PathSearch("result", getRespBody, make([]interface{}, 0)).([]interface{})
		for _, group := range groups {
			rst = append(rst, map[string]interface{}{
				"id":               utils.PathSearch("id", group, nil),
				"name":             utils.PathSearch("name", group, nil),
				"description":      utils.PathSearch("description", group, nil),
				"env_count":        utils.PathSearch("env_count", group, nil),
				"host_count":       utils.PathSearch("host_count", group, nil),
				"is_proxy_mode":    utils.PathSearch("is_proxy_mode", group, nil),
				"os_type":          utils.PathSearch("os", group, nil),
				"resource_pool_id": utils.PathSearch("slave_cluster_id", group, nil),
				"created_by":       utils.PathSearch("nick_name", group, nil),
				"permission":       flattenDeployGroupPermission(group),
			})
		}

		total := utils.PathSearch("total", getRespBody, float64(0)).(float64)
		if pageSize*(pageIndex-1)+len(groups) >= int(total) {
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
		d.Set("groups", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCodeartsDeployGroupsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?project_id=%v", d.Get("project_id"))

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("os_type"); ok {
		res = fmt.Sprintf("%s&os=%v", res, v)
	}
	if v, ok := d.GetOk("is_proxy_mode"); ok {
		res = fmt.Sprintf("%s&is_proxy_mode=%v", res, v)
	}
	if v, ok := d.GetOk("resource_pool_id"); ok {
		res = fmt.Sprintf("%s&slave_cluster_id=%v", res, v)
	}

	return res
}
