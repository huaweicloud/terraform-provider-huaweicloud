package aad

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

// Due to limited testing conditions, this data source cannot be tested and the API was not successfully called.

// @API AAD GET /v2/aad/policies/ddos/flow-block
func DataSourceFlowBlock() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFlowBlockRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_center": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"foreign_switch_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"udp_switch_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildFlowBlockQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?instance_id=%v", d.Get("instance_id"))
}

func dataSourceFlowBlockRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/policies/ddos/flow-block"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildFlowBlockQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD flow block information: %s", err)
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

	mErr := multierror.Append(nil,
		d.Set("ips", flattenFlowBlockIps(utils.PathSearch("ips", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFlowBlockIps(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"ip_id":                 utils.PathSearch("ip_id", v, nil),
			"ip":                    utils.PathSearch("ip", v, nil),
			"isp":                   utils.PathSearch("isp", v, nil),
			"data_center":           utils.PathSearch("data_center", v, nil),
			"foreign_switch_status": utils.PathSearch("foreign_switch_status", v, nil),
			"udp_switch_status":     utils.PathSearch("udp_switch_status", v, nil),
		})
	}

	return rst
}
