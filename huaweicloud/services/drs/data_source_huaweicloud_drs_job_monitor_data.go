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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/monitor-data
func DataSourceDrsJobMonitorData() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsJobMonitorDataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_src_normal": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_dst_normal": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"src_offset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_offset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dst_offset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"src_delay": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dst_delay": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"src_rps": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"src_io": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dst_rps": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dst_io": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trans_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trans_lines": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_volumes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_memory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_cpu_percent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"node_memory_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"apply_rate": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDrsJobMonitorDataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/jobs/{job_id}/monitor-data"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", d.Get("job_id").(string))

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS job monitor data: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("bandwidth", utils.PathSearch("bandwidth", respBody, nil)),
		d.Set("is_src_normal", utils.PathSearch("is_src_normal", respBody, nil)),
		d.Set("is_dst_normal", utils.PathSearch("is_dst_normal", respBody, nil)),
		d.Set("src_offset", utils.PathSearch("src_offset", respBody, nil)),
		d.Set("node_offset", utils.PathSearch("node_offset", respBody, nil)),
		d.Set("dst_offset", utils.PathSearch("dst_offset", respBody, nil)),
		d.Set("src_delay", utils.PathSearch("src_delay", respBody, nil)),
		d.Set("dst_delay", utils.PathSearch("dst_delay", respBody, nil)),
		d.Set("src_rps", utils.PathSearch("src_rps", respBody, nil)),
		d.Set("src_io", utils.PathSearch("src_io", respBody, nil)),
		d.Set("dst_rps", utils.PathSearch("dst_rps", respBody, nil)),
		d.Set("dst_io", utils.PathSearch("dst_io", respBody, nil)),
		d.Set("trans_data", utils.PathSearch("trans_data", respBody, nil)),
		d.Set("trans_lines", utils.PathSearch("trans_lines", respBody, nil)),
		d.Set("used_volumes", utils.PathSearch("used_volumes", respBody, nil)),
		d.Set("used_memory", utils.PathSearch("used_memory", respBody, nil)),
		d.Set("used_cpu_percent", utils.PathSearch("used_cpu_percent", respBody, nil)),
		d.Set("node_volume_size", utils.PathSearch("node_volume_size", respBody, nil)),
		d.Set("node_memory_size", utils.PathSearch("node_memory_size", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
		d.Set("apply_rate", utils.PathSearch("apply_rate", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
