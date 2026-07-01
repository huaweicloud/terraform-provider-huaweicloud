package gaussdb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/full-sql-switches
func DataSourceInstanceSqlExplorerStatusRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceFullSqlSwitchesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"full_sql_switches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbInstanceFullSqlSwitchesSchema(),
			},
			"allowed_sql_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbSqlTypeRangeSchema(),
			},
		},
	}
}

func gaussDbInstanceFullSqlSwitchesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"is_open": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"save_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"storage_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_exclude_sys_user": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"lts_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbLtsConfigSchema(),
			},
			"sql_type_range": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbSqlTypeRangeSchema(),
			},
		},
	}
}

func gaussDbLtsConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_ttl_in_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"group_log_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_stream_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stream_log_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stream_ttl_in_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"stream_structure_config_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stream_index_config_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func gaussDbSqlTypeRangeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prefixes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_preset": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func listFullSqlSwitches(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, []interface{}, error) {
	var (
		httpUrl         = "v3/{project_id}/instances/{instance_id}/full-sql-switches?limit={limit}"
		result          = make([]interface{}, 0)
		allowedSqlTypes = make([]interface{}, 0)
		limit           = 100
		offset          = 0
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, nil, err
		}

		if len(allowedSqlTypes) == 0 {
			allowedSqlTypes = flattenSqlTypeRange(respBody, "allowed_sql_types")
		}

		sqlSwitches := utils.PathSearch("full_sql_switches", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, sqlSwitches...)
		if len(sqlSwitches) < limit {
			break
		}

		offset += len(sqlSwitches)
	}

	return result, allowedSqlTypes, nil
}

func dataSourceInstanceFullSqlSwitchesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	switches, allowedSqlTypes, err := listFullSqlSwitches(client, d)
	if err != nil {
		return diag.Errorf("error querying GaussDB sql explorer status records: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("full_sql_switches", flattenFullSqlSwitches(switches)),
		d.Set("allowed_sql_types", allowedSqlTypes),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFullSqlSwitches(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}
	res := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		res = append(res, map[string]interface{}{
			"is_open":             utils.PathSearch("is_open", v, nil),
			"begin_time":          utils.PathSearch("begin_time", v, nil),
			"end_time":            utils.PathSearch("end_time", v, nil),
			"save_days":           utils.PathSearch("save_days", v, nil),
			"storage_mode":        utils.PathSearch("storage_mode", v, nil),
			"is_exclude_sys_user": utils.PathSearch("is_exclude_sys_user", v, nil),
			"lts_config": []interface{}{
				flattenLtsConfig(utils.PathSearch("lts_config", v, nil)),
			},
			"sql_type_range": flattenSqlTypeRange(v, "sql_type_range"),
		})
	}
	return res
}

func flattenSqlTypeRange(resp interface{}, key string) []interface{} {
	if resp == nil {
		return nil
	}

	curArray := utils.PathSearch(key, resp, make([]interface{}, 0)).([]interface{})
	res := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"category":  utils.PathSearch("category", v, nil),
			"prefixes":  utils.PathSearch("prefixes", v, nil),
			"is_preset": utils.PathSearch("is_preset", v, nil),
		})
	}
	return res
}

func flattenLtsConfig(ltsConfig interface{}) map[string]interface{} {
	if ltsConfig == nil {
		return nil
	}

	return map[string]interface{}{
		"group_ttl_in_days":          utils.PathSearch("group_ttl_in_days", ltsConfig, nil),
		"group_log_type":             utils.PathSearch("group_log_type", ltsConfig, nil),
		"log_group_name":             utils.PathSearch("log_group_name", ltsConfig, nil),
		"log_group_id":               utils.PathSearch("log_group_id", ltsConfig, nil),
		"log_stream_name":            utils.PathSearch("log_stream_name", ltsConfig, nil),
		"log_stream_id":              utils.PathSearch("log_stream_id", ltsConfig, nil),
		"stream_log_type":            utils.PathSearch("stream_log_type", ltsConfig, nil),
		"stream_ttl_in_days":         utils.PathSearch("stream_ttl_in_days", ltsConfig, nil),
		"stream_structure_config_id": utils.PathSearch("stream_structure_config_id", ltsConfig, nil),
		"stream_index_config_id":     utils.PathSearch("stream_index_config_id", ltsConfig, nil),
	}
}
