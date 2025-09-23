// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product NAT
// ---------------------------------------------------------------

package nat

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API NAT GET /v3/{project_id}/private-nat/transit-ips
func DataSourcePrivateTransitIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateTransitIpsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region where the transit IPs are located.",
			},
			"transit_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the transit IP.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address of the transit IP.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the private NAT gateway to which the transit IP belongs.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the subnet to which the transit IPs belong.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The key/value pairs to associate with the transit IPs.",
			},
			"network_interface_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The network interface ID of the transit IP for private NAT.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project to which the transit IPs belong.",
			},
			// This field is not tested due to insufficient testing conditions.
			"transit_subnet_id": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The ID of the the transit subnet.",
			},
			"transit_ips": {
				Type:        schema.TypeList,
				Elem:        transitIpSchema(),
				Computed:    true,
				Description: "The list of the transit IPs.",
			},
		},
	}
}

func transitIpSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the transit IP.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the transit IP",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the private NAT gateway to which the transit IP belongs.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the subnet to which the transit IP belongs.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "The key/value pairs to associate the transit IPs used for filter.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the transit IP.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the transit IP.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the transit IP.",
			},
			"network_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network interface ID of the transit IP for private NAT.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project to which the transit IP belongs.",
			},
		},
	}
	return &sc
}

func dataSourcePrivateTransitIpsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/transit-ips"
		product = "nat"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildListTransitIpsQueryParams(d, cfg)

	resp, err := pagination.ListAllItems(
		client,
		"marker",
		requestPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving transit IPs %s", err)
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

	curJson := utils.PathSearch("transit_ips", respBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("transit_ips", flattenListTransitIpsResponseBody(filterListTransitIpsResponseBody(curArray, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListTransitIpsResponseBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"ip_address":            utils.PathSearch("ip_address", v, nil),
			"subnet_id":             utils.PathSearch("virsubnet_id", v, nil),
			"gateway_id":            utils.PathSearch("gateway_id", v, nil),
			"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"status":                utils.PathSearch("status", v, nil),
			"created_at":            utils.PathSearch("created_at", v, nil),
			"updated_at":            utils.PathSearch("updated_at", v, nil),
			"network_interface_id":  utils.PathSearch("network_interface_id", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
		})
	}
	return rst
}

func filterListTransitIpsResponseBody(all []interface{}, d *schema.ResourceData) []interface{} {
	tagFilter := d.Get("tags").(map[string]interface{})
	if len(tagFilter) == 0 {
		return all
	}

	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		tags := utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil))
		tagMap := utils.ExpandToStringMap(tags)
		if !utils.HasMapContains(tagMap, tagFilter) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListTransitIpsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsID := cfg.GetEnterpriseProjectID(d)
	transitSubnetId := d.Get("transit_subnet_id").([]interface{})

	if v, ok := d.GetOk("transit_ip_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("ip_address"); ok {
		res = fmt.Sprintf("%s&ip_address=%v", res, v)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		res = fmt.Sprintf("%s&virsubnet_id=%v", res, v)
	}
	if v, ok := d.GetOk("gateway_id"); ok {
		res = fmt.Sprintf("%s&gateway_id=%v", res, v)
	}
	if v, ok := d.GetOk("network_interface_id"); ok {
		res = fmt.Sprintf("%s&network_interface_id=%v", res, v)
	}
	if len(transitSubnetId) > 0 {
		for _, v := range transitSubnetId {
			res = fmt.Sprintf("%s&transit_subnet_id=%v", res, v)
		}
	}
	if epsID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsID)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
