package ga

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA GET /v1/ip-groups
func DataSourceAddressGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAddressGroupsRead,
		Schema: map[string]*schema.Schema{
			"address_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the IP address group.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the IP address group.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the IP address group.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the listener associated with the IP address group.",
			},
			"address_groups": {
				Type:        schema.TypeList,
				Elem:        addressGroupsSchema(),
				Computed:    true,
				Description: "The list of the IP address groups.",
			},
		},
	}
}

func addressGroupsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the IP address group.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the IP address group.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the IP address group.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the IP address group.",
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of CIDR block configurations of the IP address group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CIDR block included in the IP address group.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the CIDR block.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the CIDR block.",
						},
					},
				},
			},
			"associated_listeners": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of the listeners associated with the IP address group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the listener associated with the IP address group.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The listener type associated with the IP address group.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the IP address group.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the IP address group.",
			},
		},
	}
	return &sc
}

func dataSourceAddressGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/ip-groups"
		product = "ga"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildListAddressGroupsQueryParams(d)
	resp, err := pagination.ListAllItems(
		client,
		"marker",
		requestPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving GA IP address groups: %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("address_groups", filterListAddressGroupsResponseBody(flattenListAddressGroupsResponseBody(respBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListAddressGroupsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("ip_groups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                   utils.PathSearch("id", v, nil),
			"name":                 utils.PathSearch("name", v, nil),
			"status":               utils.PathSearch("status", v, nil),
			"description":          utils.PathSearch("description", v, nil),
			"ip_addresses":         flattenIpList(utils.PathSearch("ip_list", v, make([]interface{}, 0))),
			"associated_listeners": flattenListeners(utils.PathSearch("associated_listeners", v, make([]interface{}, 0))),
			"created_at":           utils.PathSearch("created_at", v, nil),
			"updated_at":           utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func flattenIpList(raw interface{}) []map[string]interface{} {
	curArray := raw.([]interface{})
	result := make([]map[string]interface{}, len(curArray))
	for i, ipList := range curArray {
		result[i] = map[string]interface{}{
			"cidr":        utils.PathSearch("cidr", ipList, nil),
			"description": utils.PathSearch("description", ipList, nil),
			"created_at":  utils.PathSearch("created_at", ipList, nil),
		}
	}
	return result
}

func flattenListeners(raw interface{}) []map[string]interface{} {
	curArray := raw.([]interface{})
	result := make([]map[string]interface{}, len(curArray))
	for i, listeners := range curArray {
		result[i] = map[string]interface{}{
			"id":   utils.PathSearch("listener_id", listeners, nil),
			"type": utils.PathSearch("type", listeners, nil),
		}
	}
	return result
}

func filterListAddressGroupsResponseBody(all []interface{}, d *schema.ResourceData) []interface{} {
	var (
		addressGroupID = d.Get("address_group_id").(string)
		name           = d.Get("name").(string)
		status         = d.Get("status").(string)
		rst            = make([]interface{}, 0, len(all))
	)

	for _, v := range all {
		if addressGroupID != "" && addressGroupID != utils.PathSearch("id", v, "").(string) {
			continue
		}

		if name != "" && name != utils.PathSearch("name", v, "").(string) {
			continue
		}

		if status != "" && status != utils.PathSearch("status", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListAddressGroupsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("listener_id"); ok {
		res = fmt.Sprintf("%s&listener_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
