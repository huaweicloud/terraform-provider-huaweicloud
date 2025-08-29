package rds

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceRdsInstantTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsInstantTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tasks": {
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
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"order_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fail_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

type InstantTasksDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newInstantTasksDSWrapper(d *schema.ResourceData, meta interface{}) *InstantTasksDSWrapper {
	return &InstantTasksDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceRdsInstantTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newInstantTasksDSWrapper(d, meta)
	listTasksRst, err := wrapper.ListTasks()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listTasksToSchema(listTasksRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API RDS GET /v3/{project_id}/tasklist
func (w *InstantTasksDSWrapper) ListTasks() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rds")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/tasklist"
	params := map[string]any{
		"id":          w.Get("task_id"),
		"instance_id": w.Get("instance_id"),
		"order_id":    w.Get("order_id"),
		"name":        w.Get("name"),
		"status":      w.Get("status"),
		"start_time":  w.Get("start_time"),
		"end_time":    w.Get("end_time"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("tasks", "offset", "limit", 50).
		Request().
		Result()
}

func (w *InstantTasksDSWrapper) listTasksToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("tasks", schemas.SliceToList(body.Get("tasks"),
			func(tasks gjson.Result) any {
				return map[string]any{
					"instance_name":   tasks.Get("instance_name").Value(),
					"instance_status": tasks.Get("instance_status").Value(),
					"process":         tasks.Get("process").Value(),
					"fail_reason":     tasks.Get("fail_reason").Value(),
					"status":          tasks.Get("status").Value(),
					"id":              tasks.Get("id").Value(),
					"name":            tasks.Get("name").Value(),
					"create_time":     tasks.Get("create_time").Value(),
					"end_time":        tasks.Get("end_time").Value(),
					"instance_id":     tasks.Get("instance_id").Value(),
					"order_id":        tasks.Get("order_id").Value(),
				}
			},
		)),
		d.Set("actions", body.Get("actions").Value()),
	)
	return mErr.ErrorOrNil()
}
