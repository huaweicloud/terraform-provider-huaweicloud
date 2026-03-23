package cfw

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

// @API CFW GET /v1/{project_id}/cfw/logs/trend-attack
func DataSourceAttackLogTrend() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAttackLogTrendRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"range": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vgw_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deny_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"permit_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAttackLogTrendQueryParams(d *schema.ResourceData) string {
	rst := fmt.Sprintf(
		"?fw_instance_id=%s&log_type=%s",
		d.Get("fw_instance_id").(string),
		d.Get("log_type").(string),
	)

	if v, ok := d.GetOk("range"); ok {
		rst += fmt.Sprintf("&range=%s", v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		rst += fmt.Sprintf("&start_time=%s", v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		rst += fmt.Sprintf("&end_time=%s", v.(string))
	}

	rawArray, ok := d.Get("vgw_id").([]interface{})
	if ok {
		for _, v := range rawArray {
			rst += fmt.Sprintf("&vgw_id=%s", v.(string))
		}
	}

	return rst
}

func dataSourceAttackLogTrendRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/logs/trend-attack"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAttackLogTrendQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW attack log trend: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("data", flattenLogTrendData(utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLogTrendData(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"deny_count":   utils.PathSearch("deny_count", v, nil),
			"permit_count": utils.PathSearch("permit_count", v, nil),
			"time":         utils.PathSearch("time", v, nil),
		})
	}

	return rst
}
