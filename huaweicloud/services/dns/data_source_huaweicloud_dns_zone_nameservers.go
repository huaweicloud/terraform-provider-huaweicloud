package dns

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DNS GET /v2/zones/{zone_id}/nameservers
func DataSourceZoneNameservers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZoneNameserversRead,

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the public zone to be queried.`,
			},

			// Attributes.
			"nameservers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the name servers of the public zone.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The host name of the name server.`,
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The priority of the name server.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceZoneNameserversRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("dns", "")
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	respBody, err := getZoneNameservers(client, d.Get("zone_id").(string))
	if err != nil {
		return diag.Errorf("error querying DNS zone nameservers: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("nameservers", flattenZoneNameservers(utils.PathSearch("nameservers",
			respBody, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func getZoneNameservers(client *golangsdk.ServiceClient, zoneId string) (interface{}, error) {
	var (
		httpUrl = "v2/zones/{zone_id}/nameservers"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{zone_id}", zoneId)

	reqOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	requestResp, err := client.Request("GET", getPath, &reqOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func flattenZoneNameservers(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"hostname": utils.PathSearch("hostname", item, nil),
			"priority": utils.PathSearch("priority", item, nil),
		})
	}
	return result
}
