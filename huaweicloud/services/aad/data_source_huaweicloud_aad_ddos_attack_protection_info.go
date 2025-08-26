package aad

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

// Due to limited testing conditions, this data source cannot be tested and the API was not successfully called.

// @API AAD GET /v2/aad/instances/{instance_id}/ddos-info/flow
func DataSourceDdosAttackProtectionInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdosAttackProtectionInfoRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the instance ID.",
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the high defense IP address.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the request type. The valid values are **pps** and **bps**.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the start time.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the end time.",
			},
			"flow_bps": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The BPS flow information list. This field is returned when type is **bps**.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"utime": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The data time.",
						},
						"attack_bps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The attack traffic (bps).",
						},
						"normal_bps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The normal traffic (bps).",
						},
					},
				},
			},
			"flow_pps": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The PPS flow information list. This field is returned when type is **pps**.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"utime": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The data time.",
						},
						"attack_pps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The attack packet rate (pps).",
						},
						"normal_pps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The normal packet rate (pps).",
						},
					},
				},
			},
		},
	}
}

func buildDdosAttackProtectionInfoQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?ip=%v", d.Get("ip"))
	queryParams += fmt.Sprintf("&type=%v", d.Get("type"))
	queryParams += fmt.Sprintf("&start_time=%v", d.Get("start_time"))
	queryParams += fmt.Sprintf("&end_time=%v", d.Get("end_time"))

	return queryParams
}

func dataSourceDdosAttackProtectionInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/instances/{instance_id}/ddos-info/flow"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Get("instance_id").(string))
	requestPath += buildDdosAttackProtectionInfoQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD DDoS attack protection info: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("flow_bps", flattenDdosAttackProtectionInfoFlowBps(utils.PathSearch("flow_bps", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("flow_pps", flattenDdosAttackProtectionInfoFlowPps(utils.PathSearch("flow_pps", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDdosAttackProtectionInfoFlowBps(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"utime":      utils.PathSearch("utime", v, nil),
			"attack_bps": utils.PathSearch("attack_bps", v, nil),
			"normal_bps": utils.PathSearch("normal_bps", v, nil),
		})
	}

	return rst
}

func flattenDdosAttackProtectionInfoFlowPps(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"utime":      utils.PathSearch("utime", v, nil),
			"attack_pps": utils.PathSearch("attack_pps", v, nil),
			"normal_pps": utils.PathSearch("normal_pps", v, nil),
		})
	}

	return rst
}
