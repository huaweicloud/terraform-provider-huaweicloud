package antiddos

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

// @API ANTI-DDOS GET /v1/{project_id}/antiddos/{floating_ip_id}/daily
func DataSourceEipProtectionTraffic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEipProtectionTrafficRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"floating_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the EIP.`,
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the EIP address.`,
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of daily protection statistics.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period_start": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The start time of the statistics period.`,
						},
						"bps_in": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The inbound traffic rate in bit/s.`,
						},
						"bps_attack": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The attack traffic rate in bit/s.`,
						},
						"total_bps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The total traffic rate in bit/s.`,
						},
						"pps_in": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The inbound packet rate in pps.`,
						},
						"pps_attack": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The attack packet rate in pps.`,
						},
						"total_pps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The total packet rate in pps.`,
						},
					},
				},
			},
		},
	}
}

func buildEipProtectionTrafficQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("ip"); ok {
		return fmt.Sprintf("?ip=%s", v.(string))
	}

	return ""
}

func dataSourceEipProtectionTrafficRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/antiddos/{floating_ip_id}/daily"
		product = "anti-ddos"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{floating_ip_id}", d.Get("floating_ip_id").(string))
	requestPath += buildEipProtectionTrafficQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving Anti-DDoS daily protection traffic: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("data", flattenEipProtectionTraffic(utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEipProtectionTraffic(dataArray []interface{}) []map[string]interface{} {
	if len(dataArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataArray))
	for _, v := range dataArray {
		result = append(result, map[string]interface{}{
			"period_start": utils.PathSearch("period_start", v, nil),
			"bps_in":       utils.PathSearch("bps_in", v, 0),
			"bps_attack":   utils.PathSearch("bps_attack", v, 0),
			"total_bps":    utils.PathSearch("total_bps", v, 0),
			"pps_in":       utils.PathSearch("pps_in", v, 0),
			"pps_attack":   utils.PathSearch("pps_attack", v, 0),
			"total_pps":    utils.PathSearch("total_pps", v, 0),
		})
	}

	return result
}
