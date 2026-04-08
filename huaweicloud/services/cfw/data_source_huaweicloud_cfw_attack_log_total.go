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

// @API CFW GET /v1/{project_id}/cfw/logs/total-attack
func DataSourceAttackLogTotal() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAttackLogTotalRead,

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
			// The API documentation is of type int, but here it is changed to string to support scenarios set to `0`.
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
					Schema: attackLogTotalDataSchema(),
				},
			},
		},
	}
}

func attackLogTotalDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attack_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"deny_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"permit_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"risk_ports": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func buildAttackLogTotalQueryParams(d *schema.ResourceData) string {
	req := fmt.Sprintf("?fw_instance_id=%v&log_type=%v", d.Get("fw_instance_id"), d.Get("log_type"))

	if v, ok := d.GetOk("range"); ok {
		req = fmt.Sprintf("%s&range=%v", req, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		req = fmt.Sprintf("%s&start_time=%v", req, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		req = fmt.Sprintf("%s&end_time=%v", req, v)
	}
	if rawArray, ok := d.Get("vgw_id").([]interface{}); ok {
		for _, v := range rawArray {
			req = fmt.Sprintf("%s&vgw_id=%s", req, v.(string))
		}
	}

	return req
}

func dataSourceAttackLogTotalRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/logs/total-attack"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAttackLogTotalQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW attack log total: %s", err)
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

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenAttackTotalData(utils.PathSearch("data", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAttackTotalData(dataResp interface{}) []interface{} {
	if dataResp == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"attack_count": utils.PathSearch("attack_count", dataResp, nil),
		"deny_count":   utils.PathSearch("deny_count", dataResp, nil),
		"permit_count": utils.PathSearch("permit_count", dataResp, nil),
		"risk_ports":   utils.PathSearch("risk_ports", dataResp, nil),
	}}
}
