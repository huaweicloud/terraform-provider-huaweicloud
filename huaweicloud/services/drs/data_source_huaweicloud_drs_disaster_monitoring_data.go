package drs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS POST /v3/{project_id}/jobs/disaster-recovery-monitoring-data
func DataSourceDrsDisasterMonitoringData() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsDisasterMonitoringDataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     disasterRecoveryMonitoringResultSchema(),
			},
		},
	}
}

func disasterRecoveryMonitoringResultSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_guard_monitor": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataGuardMonitorSchema(),
			},
		},
	}
}

func dataGuardMonitorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_used_percent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dst_delay": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dst_io": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dst_normal": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dst_offset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dst_rps": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mem_used_in_mb": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_mem_in_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"node_offset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_volume_in_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sr_delay": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sr_offset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"src_io": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"src_normal": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"src_rps": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trans_in_mb": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trans_lines": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_used_in_gb": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"migration_bytes_per_second": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildDisasterMonitoringDataBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"jobs": utils.ExpandToStringList(d.Get("job_ids").([]interface{})),
	}
}

func dataSourceDrsDisasterMonitoringDataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v3/{project_id}/jobs/disaster-recovery-monitoring-data"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildDisasterMonitoringDataBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS disaster recovery monitoring data: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	resultsRaw := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("results", flattenDisasterMonitoringResults(resultsRaw)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDisasterMonitoringResults(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	resultList := make([]interface{}, 0, len(results))
	for _, result := range results {
		resultList = append(resultList, map[string]interface{}{
			"id":                 utils.PathSearch("id", result, nil),
			"data_guard_monitor": flattenDataGuardMonitor(utils.PathSearch("data_guard_minitor", result, nil)),
		})
	}

	return resultList
}

func flattenDataGuardMonitor(monitor interface{}) []interface{} {
	if monitor == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"bandwidth":                  utils.PathSearch("bandwidth", monitor, nil),
			"cpu_used_percent":           utils.PathSearch("cpuUsed_percent", monitor, nil),
			"dst_delay":                  utils.PathSearch("dst_delay", monitor, nil),
			"dst_io":                     utils.PathSearch("dst_io", monitor, nil),
			"dst_normal":                 utils.PathSearch("dst_normal", monitor, nil),
			"dst_offset":                 utils.PathSearch("dst_offset", monitor, nil),
			"dst_rps":                    utils.PathSearch("dst_rps", monitor, nil),
			"mem_used_in_mb":             utils.PathSearch("mem_used_inMB", monitor, nil),
			"node_mem_in_mb":             utils.PathSearch("node_mem_inMB", monitor, nil),
			"node_offset":                utils.PathSearch("node_offset", monitor, nil),
			"node_volume_in_gb":          utils.PathSearch("node_volume_inGB", monitor, nil),
			"sr_delay":                   utils.PathSearch("sr_delay", monitor, nil),
			"sr_offset":                  utils.PathSearch("sr_offset", monitor, nil),
			"src_io":                     utils.PathSearch("src_io", monitor, nil),
			"src_normal":                 utils.PathSearch("src_normal", monitor, nil),
			"src_rps":                    utils.PathSearch("src_rps", monitor, nil),
			"trans_in_mb":                utils.PathSearch("trans_inMB", monitor, nil),
			"trans_lines":                utils.PathSearch("trans_lines", monitor, nil),
			"volume_used_in_gb":          utils.PathSearch("volume_used_inGB", monitor, nil),
			"migration_bytes_per_second": utils.PathSearch("migration_bytes_per_second", monitor, nil),
		},
	}
}
