package dli

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDliFlinkSQLJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDliFlinkSQLJobsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"queue_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"manager_cu_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cu_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"parallel_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tm_cu_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tm_slot_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flink_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"run_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"queue_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cu_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"parallel_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"checkpoint_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"checkpoint_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"checkpoint_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"obs_bucket": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"smn_topic": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"restart_when_exception": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"idle_state_retention": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"edge_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"dirty_data_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"udf_jar_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"manager_cu_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tm_cu_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tm_slot_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resume_checkpoint": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"resume_max_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						// The runtime_config is a json string.
						"runtime_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"static_estimator_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type FlinkSQLJobsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newFlinkSQLJobsDSWrapper(d *schema.ResourceData, meta interface{}) *FlinkSQLJobsDSWrapper {
	return &FlinkSQLJobsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDliFlinkSQLJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newFlinkSQLJobsDSWrapper(d, meta)
	lisFliJobRst, err := wrapper.ListFlinkJobs()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listFlinkJobsToSchema(lisFliJobRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DLI GET /v1.0/{project_id}/streaming/jobs
func (w *FlinkSQLJobsDSWrapper) ListFlinkJobs() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dli")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/{project_id}/streaming/jobs"
	params := map[string]any{
		"queue_name":  w.Get("queue_name"),
		"tags":        w.getTags(),
		"show_detail": true,
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("job_list.jobs", "offset", "limit", 100).
		Filter(
			filters.New().From("job_list.jobs").
				Where("job_id", "=", w.GetToInt("job_id")).
				Where("job_config.cu_number", "=", w.Get("cu_num")).
				Where("job_config.parallel_number", "=", w.Get("parallel_num")).
				Where("job_config.manager_cu_number", "=", w.Get("manager_cu_num")).
				Where("job_config.tm_cus", "=", w.Get("tm_cu_num")).
				Where("job_config.tm_slot_num", "=", w.Get("tm_slot_num")),
		).
		Request().
		Result()
}

func (w *FlinkSQLJobsDSWrapper) listFlinkJobsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("jobs", schemas.SliceToList(body.Get("job_list.jobs"),
			func(job gjson.Result) any {
				return map[string]any{
					"id":                      w.setJobLisJobJobId(job),
					"name":                    job.Get("name").Value(),
					"flink_version":           job.Get("job_config.flink_version").Value(),
					"type":                    job.Get("job_type").Value(),
					"status":                  job.Get("status").Value(),
					"run_mode":                job.Get("run_mode").Value(),
					"description":             job.Get("desc").Value(),
					"queue_name":              job.Get("queue_name").Value(),
					"sql":                     job.Get("sql_body").Value(),
					"cu_num":                  job.Get("job_config.cu_number").Value(),
					"parallel_num":            job.Get("job_config.parallel_number").Value(),
					"checkpoint_enabled":      job.Get("job_config.checkpoint_enabled").Value(),
					"checkpoint_mode":         job.Get("job_config.checkpoint_mode").Value(),
					"checkpoint_interval":     job.Get("job_config.checkpoint_interval").Value(),
					"obs_bucket":              job.Get("job_config.obs_bucket").Value(),
					"log_enabled":             job.Get("job_config.log_enabled").Value(),
					"smn_topic":               job.Get("job_config.smn_topic").Value(),
					"restart_when_exception":  job.Get("job_config.restart_when_exception").Value(),
					"idle_state_retention":    job.Get("job_config.idle_state_retention").Value(),
					"edge_group_ids":          job.Get("job_config.edge_group_ids").Value(),
					"dirty_data_strategy":     job.Get("job_config.dirty_data_strategy").Value(),
					"udf_jar_url":             job.Get("job_config.udf_jar_url").Value(),
					"manager_cu_num":          job.Get("job_config.manager_cu_number").Value(),
					"tm_cu_num":               job.Get("job_config.tm_cus").Value(),
					"tm_slot_num":             job.Get("job_config.tm_slot_num").Value(),
					"resume_checkpoint":       job.Get("job_config.resume_checkpoint").Value(),
					"resume_max_num":          job.Get("job_config.resume_max_num").Value(),
					"runtime_config":          job.Get("job_config.runtime_config").Value(),
					"operator_config":         job.Get("operator_config").Value(),
					"static_estimator_config": job.Get("static_estimator_config").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (w *FlinkSQLJobsDSWrapper) getTags() string {
	raw := w.Get("tags")
	if raw == nil {
		return ""
	}

	tags := raw.(map[string]interface{})
	tagsList := make([]string, 0, len(tags))
	for k, v := range tags {
		tagsList = append(tagsList, k+"="+v.(string))
	}
	return strings.Join(tagsList, ",")
}

func (*FlinkSQLJobsDSWrapper) setJobLisJobJobId(data gjson.Result) string {
	return data.Get("job_id").String()
}
