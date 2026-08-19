package dsc

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/data-watermark-embed-tasks
func DataSourceDatabaseWatermarkEmbedTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseWatermarkEmbedTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the database watermark embed tasks are located.`,
			},

			// Optional parameters.
			"task_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the database watermark embed task.`,
			},
			"start": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"end"},
				Description:  `The start time of the task running time interval, in RFC3339 format.`,
			},
			"end": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"start"},
				Description:  `The end time of the task running time interval, in RFC3339 format.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the database watermark embed task.`,
			},

			// Attributes.
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        databaseWatermarkEmbedTasksSchema(),
				Description: `The list of the database watermark embed tasks.`,
			},
		},
	}
}

func databaseWatermarkEmbedTasksSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the database watermark embed task.`,
			},
			"task_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the database watermark embed task.`,
			},
			"water_mark": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The watermark content.`,
			},
			"db_water_param": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        databaseWatermarkEmbedTasksDbWaterParamSchema(),
				Description: `The database watermark embedding configuration.`,
			},
			"selected_fields": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the selected column.`,
						},
						"column_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data type of the selected column.`,
						},
					},
				},
				Description: `The selected field list used for watermark embedding.`,
			},
			"source_db_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        databaseWatermarkEmbedTasksDbInfoSchema(),
				Description: `The source database information.`,
			},
			"target_db_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        databaseWatermarkEmbedTasksDbInfoSchema(),
				Description: `The target database information.`,
			},
			"watermark_describe": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The watermark description.`,
			},
			"watermark_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The watermark algorithm version.`,
			},
			"start_now": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the task starts immediately.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The scheduled start time of the task, in RFC3339 format.`,
			},
			"schedule_switch": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether task scheduling is enabled.`,
			},
			"schedule_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The schedule type of the task.`,
			},
			"task_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The running state of the watermark embed task.`,
			},
			"task_create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the task, in RFC3339 format.`,
			},
			"task_end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end time of the task, in RFC3339 format.`,
			},
			"water_extract_result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The watermark extract result.`,
			},
		},
	}
}

func databaseWatermarkEmbedTasksDbWaterParamSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"embed_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The watermark embed mode.`,
			},
			"row_spacing": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The row spacing used by fake-row watermark.`,
			},
			"watermark_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The watermark key used to embed and extract the watermark.`,
			},
			"params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        databaseWatermarkEmbedTasksEmbedParamSchema(),
				Description: `The fake-column embed parameter list.`,
			},
		},
	}
}

func databaseWatermarkEmbedTasksEmbedParamSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"new_column_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the new fake column.`,
			},
			"new_column_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data type of the new fake column.`,
			},
			"fake_strategy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The strategy used to generate fake data.`,
			},
			"fake_param": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        databaseWatermarkEmbedTasksFakeParamSchema(),
				Description: `The configuration of fake data generation.`,
			},
		},
	}
}

func databaseWatermarkEmbedTasksFakeParamSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"address_accuracy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The accuracy of generated address data.`,
			},
			"date_begin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The start date of the generated date range.`,
			},
			"date_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end date of the generated date range.`,
			},
			"random_accuracy": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The precision of generated random numbers.`,
			},
			"random_begin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The lower bound of the generated random value range.`,
			},
			"random_distribute": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The distribution mode of generated random values.`,
			},
			"random_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The upper bound of the generated random value range.`,
			},
			"string_distribute": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The distribution mode of generated string values.`,
			},
		},
	}
}

func databaseWatermarkEmbedTasksDbInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the authorized database.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the database.`,
			},
			"db_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The database type.`,
			},
			"ins_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the database instance.`,
			},
			"ins_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the database instance.`,
			},
			"schema_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The schema name of the database.`,
			},
			"table_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the database table.`,
			},
		},
	}
}

func buildDatabaseWatermarkEmbedTasksQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("task_id"); ok {
		queryParams = fmt.Sprintf("%s&id=%v", queryParams, v)
	}

	if v, ok := d.GetOk("start"); ok {
		queryParams = fmt.Sprintf("%s&start=%v", queryParams, utils.ConvertTimeStrToNanoTimestamp(v.(string)))
	}

	if v, ok := d.GetOk("end"); ok {
		queryParams = fmt.Sprintf("%s&end=%v", queryParams, utils.ConvertTimeStrToNanoTimestamp(v.(string)))
	}

	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceDatabaseWatermarkEmbedTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	tasks, err := listDatabaseWatermarkEmbedTasks(client, buildDatabaseWatermarkEmbedTasksQueryParams(d))
	if err != nil {
		return diag.Errorf("error retrieving database watermark embed tasks: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", flattenDatabaseWatermarkEmbedTasks(tasks)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDatabaseWatermarkEmbedTasks(tasks []interface{}) []map[string]interface{} {
	if len(tasks) < 1 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(tasks))
	for _, v := range tasks {
		rst = append(rst, map[string]interface{}{
			"id":         strconv.Itoa(int(utils.PathSearch("id", v, float64(0)).(float64))),
			"task_name":  utils.PathSearch("task_name", v, nil),
			"water_mark": utils.PathSearch("water_mark", v, nil),
			"db_water_param": flattenDatabaseWatermarkEmbedTasksDbWaterParam(utils.PathSearch("db_water_param",
				v, make(map[string]interface{})).(map[string]interface{})),
			"selected_fields": flattenDatabaseWatermarkEmbedTasksColumnInfos(utils.PathSearch("selected_fields",
				v, make([]interface{}, 0)).([]interface{})),
			"source_db_info": flattenDatabaseWatermarkEmbedTasksDbInfo(utils.PathSearch("source_db_info",
				v, make(map[string]interface{})).(map[string]interface{})),
			"target_db_info": flattenDatabaseWatermarkEmbedTasksDbInfo(utils.PathSearch("target_db_info",
				v, make(map[string]interface{})).(map[string]interface{})),
			"watermark_describe": utils.PathSearch("watermark_describe", v, nil),
			"watermark_version":  utils.PathSearch("watermark_version", v, nil),
			"start_now":          utils.PathSearch("start_now", v, nil),
			"start_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("start_time",
				v, float64(0)).(float64))/1000, false),
			"schedule_switch": utils.PathSearch("schedule_switch", v, nil),
			"schedule_type":   utils.PathSearch("schedule_type", v, nil),
			"task_state":      utils.PathSearch("task_state", v, nil),
			"task_create_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("task_create_time",
				v, float64(0)).(float64))/1000, false),
			"task_end_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("task_end_time",
				v, float64(0)).(float64))/1000, false),
			"water_extract_result": utils.PathSearch("water_extract_result", v, nil),
		})
	}

	return rst
}

func flattenDatabaseWatermarkEmbedTasksDbWaterParam(dbWaterParam map[string]interface{}) []map[string]interface{} {
	if len(dbWaterParam) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"embed_mode":    utils.PathSearch("embed_mode", dbWaterParam, nil),
			"row_spacing":   utils.PathSearch("row_spacing", dbWaterParam, nil),
			"watermark_key": utils.PathSearch("watermark_key", dbWaterParam, nil),
			"params": flattenDatabaseWatermarkEmbedTasksEmbedParams(utils.PathSearch("params",
				dbWaterParam, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenDatabaseWatermarkEmbedTasksEmbedParams(params []interface{}) []map[string]interface{} {
	if len(params) < 1 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(params))
	for _, v := range params {
		rst = append(rst, map[string]interface{}{
			"new_column_name": utils.PathSearch("new_column_name", v, nil),
			"new_column_type": utils.PathSearch("new_column_type", v, nil),
			"fake_strategy":   utils.PathSearch("fake_strategy", v, nil),
			"fake_param": flattenDatabaseWatermarkEmbedTasksFakeParam(utils.PathSearch("fake_param",
				v, make(map[string]interface{})).(map[string]interface{})),
		})
	}

	return rst
}

func flattenDatabaseWatermarkEmbedTasksFakeParam(fakeParam map[string]interface{}) []map[string]interface{} {
	if len(fakeParam) == 0 {
		return nil
	}

	beginTime, err := strconv.ParseInt(utils.PathSearch("date_begin", fakeParam, "").(string), 10, 64)
	if err != nil {
		log.Printf("[WARN] error parsing 'date_begin' field to Integer: %s", err)
	}

	endTime, err := strconv.ParseInt(utils.PathSearch("date_end", fakeParam, "").(string), 10, 64)
	if err != nil {
		log.Printf("[WARN] error parsing 'date_end' field to Integer: %s", err)
	}

	return []map[string]interface{}{
		{
			"address_accuracy":  utils.PathSearch("address_accuracy", fakeParam, nil),
			"date_begin":        utils.FormatTimeStampRFC3339(beginTime/1000, false),
			"date_end":          utils.FormatTimeStampRFC3339(endTime/1000, false),
			"random_accuracy":   utils.PathSearch("random_accuracy", fakeParam, nil),
			"random_begin":      utils.PathSearch("random_begin", fakeParam, nil),
			"random_distribute": utils.PathSearch("random_distribute", fakeParam, nil),
			"random_end":        utils.PathSearch("random_end", fakeParam, nil),
			"string_distribute": utils.PathSearch("string_distribute", fakeParam, nil),
		},
	}
}

func flattenDatabaseWatermarkEmbedTasksColumnInfos(fields []interface{}) []map[string]interface{} {
	if len(fields) < 1 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(fields))
	for _, v := range fields {
		rst = append(rst, map[string]interface{}{
			"column_name": utils.PathSearch("column_name", v, nil),
			"column_type": utils.PathSearch("column_type", v, nil),
		})
	}

	return rst
}

func flattenDatabaseWatermarkEmbedTasksDbInfo(dbInfo map[string]interface{}) []map[string]interface{} {
	if len(dbInfo) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"db_id":       utils.PathSearch("db_id", dbInfo, nil),
			"db_name":     utils.PathSearch("db_name", dbInfo, nil),
			"db_type":     utils.PathSearch("db_type", dbInfo, nil),
			"ins_id":      utils.PathSearch("ins_id", dbInfo, nil),
			"ins_name":    utils.PathSearch("ins_name", dbInfo, nil),
			"schema_name": utils.PathSearch("schema_name", dbInfo, nil),
			"table_name":  utils.PathSearch("table_name", dbInfo, nil),
		},
	}
}
