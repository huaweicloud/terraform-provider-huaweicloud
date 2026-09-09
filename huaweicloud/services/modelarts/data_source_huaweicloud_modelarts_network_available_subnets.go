package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/networks/{network_name}/network-ip-availabilities
func DataSourceNetworkAvailableSubnets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkAvailableSubnetsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"network_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the network ID.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the subnet ID.`,
			},

			// Attributes.
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subnet name.`,
			},
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subnet ID.`,
			},
			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The number of available network IP addresses of the subnet.`,
				Elem:        networkAvailableSubnetsSchema(),
			},
		},
	}
}

func networkAvailableSubnetsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cidr of the subnet.`,
			},
			"ip_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The IP address type.`,
			},
			"used_ips": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of used IP addresses.`,
			},
			"total_ips": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of IP addresses in the subnet.`,
			},
		},
	}
}

func buildNetworkAvailableSubnetsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("subnet_id"); ok {
		res = fmt.Sprintf("%s?network_id=%v", res, v)
	}

	return res
}

func listNetworkAvailableSubnets(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, string, string, error) {
	httpUrl := "v1/{project_id}/networks/{network_name}/network-ip-availabilities"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{network_name}", d.Get("network_name").(string))
	listPath += buildNetworkAvailableSubnetsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, "", "", err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, "", "", err
	}

	subnetName := utils.PathSearch("name", respBody, "").(string)
	subnetId := utils.PathSearch("networkId", respBody, "").(string)
	ipAddresses := utils.PathSearch("subnetIpAvailability", respBody, make([]interface{}, 0)).([]interface{})

	return ipAddresses, subnetName, subnetId, nil
}

func dataSourceNetworkAvailableSubnetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	ipAddresses, subnetName, subnetId, err := listNetworkAvailableSubnets(client, d)
	if err != nil {
		return diag.Errorf("error querying network available subnets: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", subnetName),
		d.Set("network_id", subnetId),
		d.Set("subnets", flattenNetworkAvailableSubnets(ipAddresses)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNetworkAvailableSubnets(ipAddresses []interface{}) []map[string]interface{} {
	if len(ipAddresses) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(ipAddresses))
	for _, v := range ipAddresses {
		result = append(result, map[string]interface{}{
			"cidr":       utils.PathSearch("cidr", v, nil),
			"ip_version": utils.PathSearch("ipVersion", v, nil),
			"used_ips":   utils.PathSearch("usedIps", v, nil),
			"total_ips":  utils.PathSearch("totalIps", v, nil),
		})
	}

	return result
}
