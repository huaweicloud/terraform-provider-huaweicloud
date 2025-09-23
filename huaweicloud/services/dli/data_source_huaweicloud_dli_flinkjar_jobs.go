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

func DataSourceDliFlinkjarJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDliFlinkjarJobsRead,

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
						"status": {
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
						"main_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entrypoint_args": {
							Type:     schema.TypeString,
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
						"manager_cu_num": {
							Type:     schema.TypeInt,
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
						"restart_when_exception": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"entrypoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dependency_jars": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dependency_files": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"resume_checkpoint": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						// The runtime_config is a json string.
						"runtime_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resume_max_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"checkpoint_path": {
							Type:     schema.TypeString,
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
						"image": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"feature": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flink_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type FlinkjarJobsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newFlinkjarJobsDSWrapper(d *schema.ResourceData, meta interface{}) *FlinkjarJobsDSWrapper {
	return &FlinkjarJobsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDliFlinkjarJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newFlinkjarJobsDSWrapper(d, meta)
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
func (w *FlinkjarJobsDSWrapper) ListFlinkJobs() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dli")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/{project_id}/streaming/jobs"
	params := map[string]any{
		"queue_name":  w.Get("queue_name"),
		"tags":        w.getTags(),
		"job_type":    "flink_jar_job",
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

func (w *FlinkjarJobsDSWrapper) listFlinkJobsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("jobs", schemas.SliceToList(body.Get("job_list.jobs"),
			func(job gjson.Result) any {
				return map[string]any{
					"id":                     w.setJobLisJobJobId(job),
					"name":                   job.Get("name").Value(),
					"status":                 job.Get("status").Value(),
					"description":            job.Get("desc").Value(),
					"queue_name":             job.Get("queue_name").Value(),
					"main_class":             job.Get("main_class").Value(),
					"entrypoint_args":        job.Get("entrypoint_args").Value(),
					"obs_bucket":             job.Get("job_config.obs_bucket").Value(),
					"log_enabled":            job.Get("job_config.log_enabled").Value(),
					"smn_topic":              job.Get("job_config.smn_topic").Value(),
					"manager_cu_num":         job.Get("job_config.manager_cu_number").Value(),
					"cu_num":                 job.Get("job_config.cu_number").Value(),
					"parallel_num":           job.Get("job_config.parallel_number").Value(),
					"restart_when_exception": job.Get("job_config.restart_when_exception").Value(),
					"entrypoint":             job.Get("job_config.entrypoint").Value(),
					"dependency_jars":        schemas.SliceToStrList(job.Get("job_config.dependency_jars")),
					"dependency_files":       schemas.SliceToStrList(job.Get("job_config.dependency_files")),
					"resume_checkpoint":      job.Get("job_config.resume_checkpoint").Value(),
					"runtime_config":         job.Get("job_config.runtime_config").Value(),
					"resume_max_num":         job.Get("job_config.resume_max_num").Value(),
					"checkpoint_path":        job.Get("job_config.checkpoint_path").Value(),
					"tm_cu_num":              job.Get("job_config.tm_cus").Value(),
					"tm_slot_num":            job.Get("job_config.tm_slot_num").Value(),
					"image":                  job.Get("job_config.image").Value(),
					"feature":                job.Get("job_config.feature").Value(),
					"flink_version":          job.Get("job_config.flink_version").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (w *FlinkjarJobsDSWrapper) getTags() string {
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

func (*FlinkjarJobsDSWrapper) setJobLisJobJobId(data gjson.Result) string {
	return data.Get("job_id").String()
}
