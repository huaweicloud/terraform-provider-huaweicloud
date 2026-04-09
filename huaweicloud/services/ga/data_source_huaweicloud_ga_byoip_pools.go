package ga

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA GET /v1/byoip-pools
func DataSourceGaByoipPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaByoipPoolsRead,

		Schema: map[string]*schema.Schema{
			"byoip_pools": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of BYOIP pools.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the BYOIP pool.`,
						},
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CIDR block of the BYOIP pool.`,
						},
						"ip_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address version.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the BYOIP pool.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the BYOIP pool.`,
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The acceleration area.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain ID.`,
						},
					},
				},
			},
		},
	}
}

func listByoipPools(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/byoip-pools"
		result  = make([]interface{}, 0)
		limit   = 500
		marker  = ""
	)

	listPath := client.Endpoint + httpUrl
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	reqOpt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, reqOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		byoipPools := utils.PathSearch("byoip_pools", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, byoipPools...)
		if len(byoipPools) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func dataSourceGaByoipPoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("ga", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	byoipPools, err := listByoipPools(client)
	if err != nil {
		return diag.Errorf("error listing GA BYOIP pools: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}

	d.SetId(generateUUID)

	return diag.FromErr(d.Set("byoip_pools", flattenByoipPools(byoipPools)))
}

func flattenByoipPools(byoipPools []interface{}) []interface{} {
	if len(byoipPools) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(byoipPools))
	for _, pool := range byoipPools {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", pool, nil),
			"cidr":       utils.PathSearch("cidr", pool, nil),
			"ip_type":    utils.PathSearch("ip_type", pool, nil),
			"created_at": utils.PathSearch("created_at", pool, nil),
			"updated_at": utils.PathSearch("updated_at", pool, nil),
			"area":       utils.PathSearch("area", pool, nil),
			"domain_id":  utils.PathSearch("domain_id", pool, nil),
		})
	}

	return result
}
