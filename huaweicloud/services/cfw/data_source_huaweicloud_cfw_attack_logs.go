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

// @API CFW GET /v1/{project_id}/cfw/logs/attack
func DataSourceCfwAttackLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCfwAttackLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the firewall instance ID.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the start time.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the end time.`,
			},
			"src_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the source IP address.`,
			},
			"src_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the source port.`,
			},
			"dst_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the destination IP address.`,
			},
			"dst_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the destination port.`,
			},
			"app": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the application protocol.`,
			},
			"attack_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the intrusion event type.`,
			},
			"attack_rule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the intrusion event rule.`,
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the threat level.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"log_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the log type.`,
			},
			"attack_rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the attack rule ID.`,
			},
			"src_region_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the source region name.`,
			},
			"dst_region_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the destination region name.`,
			},
			"src_province_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the source province name.`,
			},
			"dst_province_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the destination province name.`,
			},
			"src_city_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the source city name.`,
			},
			"dst_city_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the destination city name.`,
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The attack log records.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"packet": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The attack log packet.`,
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source.`,
						},
						"src_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source IP address.`,
						},
						"src_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The source port.`,
						},
						"direction": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The direction.`,
						},
						"dst_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The destination port.`,
						},
						"app": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application protocol.`,
						},
						"attack_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The attack rule ID.`,
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The protocol.`,
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Then action.`,
						},
						"event_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The event time.`,
						},
						"attack_rule": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The attack rule.`,
						},
						"log_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The log ID.`,
						},
						"dst_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The destination IP address.`,
						},
						"packet_messages": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The packet messages.`,
							Elem:        dataRecPacMesElem(),
						},
						"attack_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The attack type.`,
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The threat level.`,
						},
						"packet_length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The packet length.`,
						},
						"src_region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source region ID.`,
						},
						"src_region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source region name.`,
						},
						"dst_region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The destination region ID.`,
						},
						"dst_region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The destination region name.`,
						},
						"src_province_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source province ID.`,
						},
						"src_province_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source province name.`,
						},
						"src_city_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source city ID.`,
						},
						"src_city_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source city name.`,
						},
						"dst_province_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The distination province ID.`,
						},
						"dst_province_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The distination province name.`,
						},
						"dst_city_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The distination city ID.`,
						},
						"dst_city_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The distination city name.`,
						},
					},
				},
			},
		},
	}
}

func dataRecPacMesElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"utf8_string": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The utf-8 string.`,
			},
			"hex_index": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The hexadecimal index.`,
			},
			"hexs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The hexadecimal series.`,
			},
		},
	}
}

func dataSourceCfwAttackLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cfw", region)

	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	result, err := getLogs(client, cfg, d)
	if err != nil {
		return diag.Errorf("error retrieving CFW attack logs: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("records", flattenListAttackLogs(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getLogs(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/cfw/logs/attack"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)

	params, err := buildCfwAttackLogQueryParams(d, cfg)
	if err != nil {
		return nil, fmt.Errorf("error building CFW attack log query params: %s", err)
	}
	path += params

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", path, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	curJson := utils.PathSearch("data.records", respBody, make([]interface{}, 0))

	return curJson.([]interface{}), nil
}

func buildCfwAttackLogQueryParams(d *schema.ResourceData, cfg *config.Config) (string, error) {
	epsID := cfg.GetEnterpriseProjectID(d)

	startTime, err := utils.FormatUTCTimeStamp(d.Get("start_time").(string))
	if err != nil {
		return "", err
	}
	endTime, err := utils.FormatUTCTimeStamp(d.Get("end_time").(string))
	if err != nil {
		return "", err
	}
	res := fmt.Sprintf("?fw_instance_id=%v&start_time=%v&end_time=%v&limit=1000", d.Get("fw_instance_id").(string), startTime*1000, endTime*1000)

	if v, ok := d.GetOk("src_ip"); ok {
		res = fmt.Sprintf("%s&src_ip=%v", res, v)
	}
	if v, ok := d.GetOk("src_port"); ok {
		res = fmt.Sprintf("%s&src_port=%v", res, v)
	}
	if v, ok := d.GetOk("dst_ip"); ok {
		res = fmt.Sprintf("%s&dst_ip=%v", res, v)
	}
	if v, ok := d.GetOk("dst_port"); ok {
		res = fmt.Sprintf("%s&dst_port=%v", res, v)
	}
	if v, ok := d.GetOk("app"); ok {
		res = fmt.Sprintf("%s&app=%v", res, v)
	}
	if v, ok := d.GetOk("attack_type"); ok {
		res = fmt.Sprintf("%s&attack_type=%v", res, v)
	}
	if v, ok := d.GetOk("attack_rule"); ok {
		res = fmt.Sprintf("%s&attack_rule=%v", res, v)
	}
	if v, ok := d.GetOk("level"); ok {
		res = fmt.Sprintf("%s&level=%v", res, v)
	}
	if epsID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsID)
	}
	if v, ok := d.GetOk("log_type"); ok {
		res = fmt.Sprintf("%s&log_type=%v", res, v)
	}
	if v, ok := d.GetOk("attack_rule_id"); ok {
		res = fmt.Sprintf("%s&attack_rule_id=%v", res, v)
	}
	if v, ok := d.GetOk("src_region_name"); ok {
		res = fmt.Sprintf("%s&src_region_name=%v", res, v)
	}
	if v, ok := d.GetOk("dst_region_name"); ok {
		res = fmt.Sprintf("%s&dst_region_name=%v", res, v)
	}
	if v, ok := d.GetOk("src_province_name"); ok {
		res = fmt.Sprintf("%s&src_province_name=%v", res, v)
	}
	if v, ok := d.GetOk("dst_province_name"); ok {
		res = fmt.Sprintf("%s&dst_province_name=%v", res, v)
	}
	if v, ok := d.GetOk("src_city_name"); ok {
		res = fmt.Sprintf("%s&src_city_name=%v", res, v)
	}
	if v, ok := d.GetOk("dst_city_name"); ok {
		res = fmt.Sprintf("%s&dst_city_name=%v", res, v)
	}

	return res, nil
}

func flattenListAttackLogs(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		eventTime := int64(utils.PathSearch("event_time", v, float64(0)).(float64) / 1000)
		rst = append(rst, map[string]interface{}{
			"packet":            utils.PathSearch("packet", v, nil),
			"source":            utils.PathSearch("source", v, nil),
			"src_ip":            utils.PathSearch("src_ip", v, nil),
			"src_port":          utils.PathSearch("src_port", v, nil),
			"direction":         utils.PathSearch("direction", v, nil),
			"dst_port":          utils.PathSearch("dst_port", v, nil),
			"app":               utils.PathSearch("app", v, nil),
			"attack_rule_id":    utils.PathSearch("attack_rule_id", v, nil),
			"protocol":          utils.PathSearch("protocol", v, nil),
			"action":            utils.PathSearch("action", v, nil),
			"event_time":        utils.FormatTimeStampUTC(eventTime),
			"attack_rule":       utils.PathSearch("attack_rule", v, nil),
			"log_id":            utils.PathSearch("log_id", v, nil),
			"dst_ip":            utils.PathSearch("dst_ip", v, nil),
			"packet_messages":   flattenPacketMessages(v),
			"attack_type":       utils.PathSearch("attack_type", v, nil),
			"level":             utils.PathSearch("level", v, nil),
			"packet_length":     utils.PathSearch("packet_length", v, nil),
			"src_region_id":     utils.PathSearch("src_region_id", v, nil),
			"src_region_name":   utils.PathSearch("src_region_name", v, nil),
			"dst_region_id":     utils.PathSearch("dst_region_id", v, nil),
			"dst_region_name":   utils.PathSearch("dst_region_name", v, nil),
			"src_province_id":   utils.PathSearch("src_province_id", v, nil),
			"src_province_name": utils.PathSearch("src_province_name", v, nil),
			"src_city_id":       utils.PathSearch("src_city_id", v, nil),
			"src_city_name":     utils.PathSearch("src_city_name", v, nil),
			"dst_province_id":   utils.PathSearch("dst_province_id", v, nil),
			"dst_province_name": utils.PathSearch("dst_province_name", v, nil),
			"dst_city_id":       utils.PathSearch("dst_city_id", v, nil),
			"dst_city_name":     utils.PathSearch("dst_city_name", v, nil),
		})
	}

	return rst
}

func flattenPacketMessages(resp interface{}) []interface{} {
	curJson := utils.PathSearch("packetMessages", resp, nil)

	if curJson == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	if curArray, ok := curJson.([]interface{}); ok {
		for _, item := range curArray {
			rst = append(rst, map[string]interface{}{
				"hex_index":   utils.PathSearch("hex_index", item, nil),
				"hexs":        utils.PathSearch("hexs", item, make([]interface{}, 0)),
				"utf8_string": utils.PathSearch("utf8_String", item, nil),
			})
		}
	}
	return rst
}
