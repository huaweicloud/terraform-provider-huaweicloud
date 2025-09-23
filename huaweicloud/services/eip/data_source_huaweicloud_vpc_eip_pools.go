package eip

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP GET /v3/{project_id}/eip/publicip-pools
func DataSourceVpcEipPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcEipPoolsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pool name.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the pool size.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pool status.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pool type.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pool description.`,
			},
			"public_border_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether the pool is at the center or at the edge.`,
			},
			"pools": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the public network pools.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the pool ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the pool name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the pool type.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the pool size.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the pool status.`,
						},
						"used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `indicates the number of used IP addresses.`,
						},
						"shared": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to share the pool.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the description.`,
						},
						"public_border_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates whether the pool is at the center or at the edge.`,
						},
						"allow_share_bandwidth_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Indicates the list of shared bandwidth types to which the public IP address can be added.`,
						},
						"billing_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the order information. If an order is available, it indicates a yearly/monthly pool.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"order_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the order ID.`,
									},
									"product_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the product ID`,
									},
								},
							},
						},
						"tags": common.TagsComputedSchema(),
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of an enterprise project.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the create time of the pool.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the update time of the pool.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcEipPoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	getHttpUrl := "v3/{project_id}/eip/publicip-pools"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPath += fmt.Sprintf("?limit=%v", pageLimit)
	getPath = buildQueryPoolsListPath(d, getPath)

	marker := ""
	results := make([]map[string]interface{}, 0)
	currentPath := getPath
	for {
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		pools := utils.PathSearch("publicip_pools", getRespBody, make([]interface{}, 0)).([]interface{})
		for _, pool := range pools {
			results = append(results, map[string]interface{}{
				"id":                          utils.PathSearch("id", pool, nil),
				"name":                        utils.PathSearch("name", pool, nil),
				"type":                        utils.PathSearch("type", pool, nil),
				"status":                      utils.PathSearch("status", pool, nil),
				"allow_share_bandwidth_types": utils.PathSearch("allow_share_bandwidth_types", pool, nil),
				"public_border_group":         utils.PathSearch("public_border_group", pool, nil),
				"used":                        utils.PathSearch("used", pool, nil),
				"size":                        utils.PathSearch("size", pool, nil),
				"shared":                      utils.PathSearch("shared", pool, nil),
				"description":                 utils.PathSearch("description", pool, nil),
				"enterprise_project_id":       utils.PathSearch("enterprise_project_id", pool, nil),
				"created_at":                  utils.PathSearch("created_at", pool, nil),
				"updated_at":                  utils.PathSearch("updated_at", pool, nil),
				"billing_info":                flattenEipPoolsBillingInfo(utils.PathSearch("billing_info", pool, nil)),
				"tags":                        utils.FlattenTagsToMap(utils.PathSearch("tags", pool, nil)),
			})
		}

		if len(pools) < pageLimit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", getRespBody, "").(string)
		if marker == "" {
			break
		}
		currentPath = getPath + fmt.Sprintf("&marker=%s", marker)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("pools", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// some attrs do not return by default and need to be specified in `field` key to return them
var getPoolsKeys = []string{"id", "name", "type", "size", "status", "allow_share_bandwidth_types", "public_border_group",
	"used", "shared", "description", "billing_info", "tags", "enterprise_project_id", "created_at", "updated_at"}

func buildQueryPoolsListPath(d *schema.ResourceData, getPath string) string {
	for _, k := range getPoolsKeys {
		getPath += fmt.Sprintf("&fields=%s", k)
	}
	if name, ok := d.GetOk("name"); ok {
		getPath += fmt.Sprintf("&name=%s", name)
	}
	if size, ok := d.GetOk("size"); ok {
		getPath += fmt.Sprintf("&size=%s", size)
	}
	if status, ok := d.GetOk("status"); ok {
		getPath += fmt.Sprintf("&status=%s", status)
	}
	if poolType, ok := d.GetOk("type"); ok {
		getPath += fmt.Sprintf("&type=%s", poolType)
	}
	if description, ok := d.GetOk("description"); ok {
		getPath += fmt.Sprintf("&description=%s", description)
	}
	if publicBorderGroup, ok := d.GetOk("public_border_group"); ok {
		getPath += fmt.Sprintf("&public_border_group=%s", publicBorderGroup)
	}

	return getPath
}

func flattenEipPoolsBillingInfo(params interface{}) interface{} {
	return []map[string]interface{}{
		{
			"order_id":   utils.PathSearch("order_id", params, nil),
			"product_id": utils.PathSearch("product_id", params, nil),
		},
	}
}
