package gaussdb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/limit-task
// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/limit-task/{task_id}
// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/limit-task-list
// @API GaussDB DELETE /v3/{project_id}/instances/{instance_id}/limit-task/{task_id}
func ResourceOpenGaussSqlThrottlingTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussSqlThrottlingTaskCreate,
		UpdateContext: resourceOpenGaussSqlThrottlingTaskUpdate,
		ReadContext:   resourceOpenGaussSqlThrottlingTaskRead,
		DeleteContext: resourceOpenGaussSqlThrottlingTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOpenGaussSqlThrottlingTaskImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"task_scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"limit_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"limit_type_value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parallel_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_words": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sql_model": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"cpu_utilization": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"memory_utilization": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"databases": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_infos": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     sqlThrottlingTaskNodeInfosSchema(),
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func sqlThrottlingTaskNodeInfosSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sql_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	return &sc
}

func resourceOpenGaussSqlThrottlingTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/limit-task"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateOpenGaussSqlThrottlingTaskBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB OpenGauss SQL throttling task: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskId := utils.PathSearch("task_id", createRespBody, "").(string)
	if taskId == "" {
		return diag.Errorf("error creating GaussDB OpenGauss SQL throttling task: task_id is not found in API response")
	}
	d.SetId(taskId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"WAIT_EXCUTE", "EXCUTING"},
		Refresh:      sqlThrottlingTaskStateRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for SQL throttling task(%s) to completed: %s", d.Id(), err)
	}
	return resourceOpenGaussSqlThrottlingTaskRead(ctx, d, meta)
}

func buildCreateOpenGaussSqlThrottlingTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"task_scope":         d.Get("task_scope"),
		"limit_type":         d.Get("limit_type"),
		"limit_type_value":   d.Get("limit_type_value"),
		"task_name":          d.Get("task_name"),
		"parallel_size":      d.Get("parallel_size"),
		"start_time":         utils.ValueIgnoreEmpty(d.Get("start_time")),
		"end_time":           utils.ValueIgnoreEmpty(d.Get("end_time")),
		"key_words":          utils.ValueIgnoreEmpty(d.Get("key_words")),
		"sql_model":          utils.ValueIgnoreEmpty(d.Get("sql_model")),
		"cpu_utilization":    utils.ValueIgnoreEmpty(d.Get("cpu_utilization")),
		"memory_utilization": utils.ValueIgnoreEmpty(d.Get("memory_utilization")),
		"databases":          utils.ValueIgnoreEmpty(d.Get("databases")),
		"node_infos":         buildCreateSqlThrottlingTaskBodyParam(d),
	}
	return bodyParams
}

func buildCreateSqlThrottlingTaskBodyParam(d *schema.ResourceData) []interface{} {
	nodeInfos := d.Get("node_infos").(*schema.Set)
	if nodeInfos.Len() == 0 {
		return nil
	}
	rst := make([]interface{}, 0, nodeInfos.Len())
	for _, v := range nodeInfos.List() {
		if raw, ok := v.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"node_id": raw["node_id"].(string),
				"sql_id":  raw["sql_id"].(string),
			})
		}
	}
	return rst
}

func resourceOpenGaussSqlThrottlingTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/limit-task/{task_id}"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{task_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateOpenGaussSqlThrottlingTaskBodyParams(d))

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating GaussDB OpenGauss SQL throttling task: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"UPDATING"},
		Target:       []string{"WAIT_EXCUTE", "EXCUTING"},
		Refresh:      sqlThrottlingTaskStateRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for SQL throttling task(%s) to completed: %s", d.Id(), err)
	}

	return resourceOpenGaussSqlThrottlingTaskRead(ctx, d, meta)
}

func buildUpdateOpenGaussSqlThrottlingTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time":         utils.ValueIgnoreEmpty(d.Get("start_time")),
		"end_time":           utils.ValueIgnoreEmpty(d.Get("end_time")),
		"key_words":          utils.ValueIgnoreEmpty(d.Get("key_words")),
		"task_name":          utils.ValueIgnoreEmpty(d.Get("task_name")),
		"parallel_size":      utils.ValueIgnoreEmpty(d.Get("parallel_size")),
		"cpu_utilization":    utils.ValueIgnoreEmpty(d.Get("cpu_utilization")),
		"memory_utilization": utils.ValueIgnoreEmpty(d.Get("memory_utilization")),
		"databases":          utils.ValueIgnoreEmpty(d.Get("databases")),
	}
	return bodyParams
}

func resourceOpenGaussSqlThrottlingTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	sqlThrottlingTask, err := getSqlThrottlingTask(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB OpenGauss SQL throttling task")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", sqlThrottlingTask, nil)),
		d.Set("task_scope", utils.PathSearch("task_scope", sqlThrottlingTask, nil)),
		d.Set("limit_type", utils.PathSearch("limit_type", sqlThrottlingTask, nil)),
		d.Set("limit_type_value", utils.PathSearch("limit_type_value", sqlThrottlingTask, nil)),
		d.Set("task_name", utils.PathSearch("task_name", sqlThrottlingTask, nil)),
		d.Set("parallel_size", utils.PathSearch("parallel_size", sqlThrottlingTask, nil)),
		d.Set("key_words", utils.PathSearch("key_words", sqlThrottlingTask, nil)),
		d.Set("sql_model", utils.PathSearch("sql_model", sqlThrottlingTask, nil)),
		d.Set("cpu_utilization", utils.PathSearch("cpu_utilization", sqlThrottlingTask, nil)),
		d.Set("memory_utilization", utils.PathSearch("memory_utilization", sqlThrottlingTask, nil)),
		d.Set("databases", utils.PathSearch("databases", sqlThrottlingTask, nil)),
		d.Set("node_infos", flattenSqlThrottlingTaskResponseBodyNodeInfos(sqlThrottlingTask)),
		d.Set("created_at", utils.PathSearch("created", sqlThrottlingTask, nil)),
		d.Set("updated_at", utils.PathSearch("updated", sqlThrottlingTask, nil)),
		d.Set("creator", utils.PathSearch("creator", sqlThrottlingTask, nil)),
		d.Set("modifier", utils.PathSearch("modifier", sqlThrottlingTask, nil)),
		d.Set("rule_name", utils.PathSearch("rule_name", sqlThrottlingTask, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSqlThrottlingTaskResponseBodyNodeInfos(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("node_infos", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"node_id": utils.PathSearch("node_id", v, nil),
			"sql_id":  utils.PathSearch("sql_id", v, nil),
		})
	}
	return rst
}

func getSqlThrottlingTask(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/limit-task-list"
	)

	listBasePath := client.Endpoint + httpUrl
	listBasePath = strings.ReplaceAll(listBasePath, "{project_id}", client.ProjectID)
	listBasePath = strings.ReplaceAll(listBasePath, "{instance_id}", d.Get("instance_id").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var sqlThrottlingTask interface{}
	offset := 0
	for {
		listPath := listBasePath + buildPageQueryParam(100, offset)
		requestResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return nil, common.ConvertUndefinedErrInto404Err(err, 500, "errCode", "DBS.280005")
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		sqlThrottlingTasks := utils.PathSearch("limit_task_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(sqlThrottlingTasks) == 0 {
			break
		}
		searchExpression := fmt.Sprintf("[?task_id=='%s']|[0]", d.Id())
		sqlThrottlingTask = utils.PathSearch(searchExpression, sqlThrottlingTasks, nil)
		if sqlThrottlingTask != nil {
			break
		}

		offset += 100
		totalCount := utils.PathSearch("total_count", respBody, float64(0)).(float64)
		if int(totalCount) <= offset {
			break
		}
	}
	if sqlThrottlingTask == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return sqlThrottlingTask, nil
}

func buildPageQueryParam(limit, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func sqlThrottlingTaskStateRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := getSqlThrottlingTask(client, d)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		status := utils.PathSearch("status", res, "").(string)
		return res, status, nil
	}
}

func resourceOpenGaussSqlThrottlingTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/limit-task/{task_id}"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{task_id}", d.Id())

	deleteGOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteGOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errCode", "DBS.06010013"),
			"error deleting GaussDB OpenGauss SQL throttling task")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING"},
		Target:       []string{"DELETED"},
		Refresh:      sqlThrottlingTaskStateRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		PollInterval: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for SQL throttling task(%s) to completed: %s", d.Id(), err)
	}

	return nil
}

func resourceOpenGaussSqlThrottlingTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import id, must be <instance_id>/<id>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
