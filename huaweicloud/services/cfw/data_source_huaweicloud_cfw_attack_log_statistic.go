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

// @API CFW GET /v1/{project_id}/cfw/logs/attack-statistic
func DataSourceAttackLogStatistic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAttackLogStatisticRead,

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
			"item": {
				Type:     schema.TypeString,
				Required: true,
			},
			// The API documentation is of type int, but here it is changed to string to support scenarios set to `0`.
			"range": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": {
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
			// The API documentation is of type int, but here it is changed to string to support scenarios set to `0`.
			"size": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: attackLogStatisticDataSchema(),
				},
			},
		},
	}
}

func attackLogStatisticDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"apps": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     attackStatisticTopInfoSchema(),
		},
		"associated_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"associated_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"attack_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"attack_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deny_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"dst_ports": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     attackStatisticTopInfoSchema(),
		},
		"ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"latest_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"region_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"region_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"src_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vgw_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func attackStatisticTopInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"item": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"item_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildAttackLogStatisticQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?fw_instance_id=%v&log_type=%v&item=%v",
		d.Get("fw_instance_id"), d.Get("log_type"), d.Get("item"))

	if v, ok := d.GetOk("range"); ok {
		res = fmt.Sprintf("%s&range=%v", res, v)
	}
	if v, ok := d.GetOk("direction"); ok {
		res = fmt.Sprintf("%s&direction=%v", res, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}
	if rawArray, ok := d.Get("vgw_id").([]interface{}); ok {
		for _, v := range rawArray {
			res = fmt.Sprintf("%s&vgw_id=%s", res, v.(string))
		}
	}
	if v, ok := d.GetOk("size"); ok {
		res = fmt.Sprintf("%s&size=%v", res, v)
	}

	return res
}

func dataSourceAttackLogStatisticRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/logs/attack-statistic"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAttackLogStatisticQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW attack log statistic: %s", err)
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
		d.Set("data", flattenAttackLogStatisticData(
			utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAttackLogStatisticData(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataResp))
	for _, raw := range dataResp {
		result = append(result, map[string]interface{}{
			"apps": flattenAttackLogStatisticTopInfos(
				utils.PathSearch("apps", raw, make([]interface{}, 0)).([]interface{})),
			"associated_name": utils.PathSearch("associated_name", raw, nil),
			"associated_type": utils.PathSearch("associated_type", raw, nil),
			"attack_count":    utils.PathSearch("attack_count", raw, nil),
			"attack_type":     utils.PathSearch("attack_type", raw, nil),
			"deny_count":      utils.PathSearch("deny_count", raw, nil),
			"dst_ports": flattenAttackLogStatisticTopInfos(
				utils.PathSearch("dst_ports", raw, make([]interface{}, 0)).([]interface{})),
			"ip":          utils.PathSearch("ip", raw, nil),
			"latest_time": utils.PathSearch("latest_time", raw, nil),
			"region_id":   utils.PathSearch("region_id", raw, nil),
			"region_name": utils.PathSearch("region_name", raw, nil),
			"src_type":    utils.PathSearch("src_type", raw, nil),
			"vgw_id":      utils.PathSearch("vgw_id", raw, nil),
		})
	}

	return result
}

func flattenAttackLogStatisticTopInfos(rawResp []interface{}) []interface{} {
	if len(rawResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rawResp))
	for _, raw := range rawResp {
		result = append(result, map[string]interface{}{
			"count":   utils.PathSearch("count", raw, nil),
			"item":    utils.PathSearch("item", raw, nil),
			"item_id": utils.PathSearch("item_id", raw, nil),
		})
	}

	return result
}
