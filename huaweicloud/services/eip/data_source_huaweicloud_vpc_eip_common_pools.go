package eip

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP GET /v3/{project_id}/eip/publicip-pools/common-pools
// @API EIP POST /v3/{project_id}/eip/resources/available
func DataSourceVpcEipCommonPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEIPVpcEipCommonPoolsRead,

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
				Description: `Specifies the common pool name.`,
			},
			"public_border_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether the common pool is at the center or at the edge.`,
			},
			"common_pools": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the common pools.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the common pool ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the common pool name。`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the common pool type, such as **bgp** and **sbgp**.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the common pool status.`,
						},
						"allow_share_bandwidth_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Indicates the list of shared bandwidth types that the public IP address can be added to.`,
						},
						"public_border_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the whether common pool is at the central site or edge site.`,
						},
						"used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of used IP addresses.`,
						},
						"available": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of available IP addresses.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceEIPVpcEipCommonPoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	getHttpUrl := "v3/{project_id}/eip/publicip-pools/common-pools"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPath += fmt.Sprintf("?limit=%v", pageLimit)
	getPath = buildQueryCommonPoolsListPath(d, getPath)

	currentTotal := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", currentTotal)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		commonPools := utils.PathSearch("common_pools", getRespBody, make([]interface{}, 0)).([]interface{})
		for _, commonPool := range commonPools {
			name := utils.PathSearch("name", commonPool, "").(string)
			available := 0
			if name != "" {
				available, err = getCommonPoolAvailableResources(client, name)
				if err != nil {
					log.Printf("[WARN] failed to fetch the available num for common pool (%s): %s", name, err)
				}
			}
			results = append(results, map[string]interface{}{
				"id":                          utils.PathSearch("id", commonPool, nil),
				"name":                        name,
				"type":                        utils.PathSearch("type", commonPool, nil),
				"status":                      utils.PathSearch("status", commonPool, nil),
				"allow_share_bandwidth_types": utils.PathSearch("allow_share_bandwidth_types", commonPool, nil),
				"public_border_group":         utils.PathSearch("public_border_group", commonPool, nil),
				"used":                        utils.PathSearch("used", commonPool, nil),
				"available":                   available,
			})
		}

		length := len(commonPools)
		if length < pageLimit {
			break
		}

		currentTotal += length
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("common_pools", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// some attrs do not return by default and need to be specified in `field` key to return them
var getCommonPoolsKeys = []string{"id", "name", "type", "status", "allow_share_bandwidth_types", "public_border_group", "used"}

func buildQueryCommonPoolsListPath(d *schema.ResourceData, getPath string) string {
	for _, k := range getCommonPoolsKeys {
		getPath += fmt.Sprintf("&fields=%s", k)
	}
	if name, ok := d.GetOk("name"); ok {
		getPath += fmt.Sprintf("&name=%s", name)
	}
	if publicBorderGroup, ok := d.GetOk("public_border_group"); ok {
		getPath += fmt.Sprintf("&public_border_group=%s", publicBorderGroup)
	}

	return getPath
}

func getCommonPoolAvailableResources(client *golangsdk.ServiceClient, name string) (int, error) {
	getHttpUrl := "v3/{project_id}/eip/resources/available"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"type": name,
		},
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return 0, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return 0, err
	}

	return int(utils.PathSearch("result", getRespBody, float64(0)).(float64)), nil
}
