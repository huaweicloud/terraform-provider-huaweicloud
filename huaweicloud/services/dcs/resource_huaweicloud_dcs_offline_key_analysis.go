package dcs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var offlineKeyAnalysisNonUpdatableParams = []string{
	"instance_id",
	"node_id",
	"backup_id",
}

// @API DCS POST /v2/{project_id}/instances/{instance_id}/offline/key-analysis
// @API DCS GET /v2/{project_id}/instances/{instance_id}
// @API DCS GET /v2/{project_id}/instances/{instance_id}/offline/key-analysis/{task_id}
// @API DCS GET /v2/{project_id}/instances/{instance_id}/offline/key-analysis
// @API DCS DELETE /v2/{project_id}/instances/{instance_id}/offline/key-analysis/{task_id}
func ResourceOfflineKeyAnalysis() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOfflineKeyAnalysisCreate,
		UpdateContext: resourceOfflineKeyAnalysisUpdate,
		ReadContext:   resourceOfflineKeyAnalysisRead,
		DeleteContext: resourceOfflineKeyAnalysisDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOfflineKeyAnalysisImport,
		},

		CustomizeDiff: config.FlexibleForceNew(offlineKeyAnalysisNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

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
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"analysis_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"started_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finished_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"largest_key_prefixes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     offlineKeyAnalysisLargestKeyPrefixSchema(),
			},
			"largest_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     offlineKeyAnalysisLargestKeySchema(),
			},
			"total_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type_bytes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     offlineKeyAnalysisKeyTypeByteSchema(),
			},
			"type_num": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     offlineKeyAnalysisKeyTypeNumSchema(),
			},
		},
	}
}

func offlineKeyAnalysisLargestKeyPrefixSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func offlineKeyAnalysisLargestKeySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"num_of_elem": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func offlineKeyAnalysisKeyTypeByteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func offlineKeyAnalysisKeyTypeNumSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceOfflineKeyAnalysisCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/offline/key-analysis"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateOfflineKeyAnalysisBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		r, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleOperationError(err)
		return r, retry, err
	}
	createResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"RUNNING"},
		WaitPending:  []string{"PENDING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating DCS offline key analysis: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	taskId := utils.PathSearch("task_id", createRespBody, "").(string)
	if taskId == "" {
		return diag.Errorf("error creating DCS offline key analysis: task_id is not found in API response")
	}

	d.SetId(taskId)

	if err = checkOfflineKeyAnalysisJobFinish(ctx, client, d, "success", d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for DCS offline key analysis task (%s) to complete: %s", taskId, err)
	}

	return resourceOfflineKeyAnalysisRead(ctx, d, meta)
}

func buildCreateOfflineKeyAnalysisBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_id":   d.Get("node_id"),
		"backup_id": utils.ValueIgnoreEmpty(d.Get("backup_id")),
	}
	return bodyParams
}

func resourceOfflineKeyAnalysisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/offline/key-analysis/{task_id}"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{task_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DCS.4232"),
			"error retrieving DCS offline key analysis")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		nil,
		d.Set("instance_id", utils.PathSearch("instance_id", getRespBody, nil)),
		d.Set("node_id", utils.PathSearch("node_id", getRespBody, nil)),
		d.Set("backup_id", utils.PathSearch("backup_id", getRespBody, nil)),
		d.Set("group_name", utils.PathSearch("group_name", getRespBody, nil)),
		d.Set("node_ip", utils.PathSearch("node_ip", getRespBody, nil)),
		d.Set("node_ipv6", utils.PathSearch("node_ipv6", getRespBody, nil)),
		d.Set("node_type", utils.PathSearch("node_type", getRespBody, nil)),
		d.Set("analysis_type", utils.PathSearch("analysis_type", getRespBody, nil)),
		d.Set("started_at", utils.PathSearch("started_at", getRespBody, nil)),
		d.Set("finished_at", utils.PathSearch("finished_at", getRespBody, nil)),
		d.Set("total_bytes", utils.PathSearch("total_bytes", getRespBody, nil)),
		d.Set("total_num", utils.PathSearch("total_num", getRespBody, nil)),
		d.Set("largest_key_prefixes", flattenOfflineKeyAnalysisLargestKeyPrefixes(getRespBody)),
		d.Set("largest_keys", flattenOfflineKeyAnalysisLargestKeys(getRespBody)),
		d.Set("type_bytes", flattenOfflineKeyAnalysisTypeBytes(getRespBody)),
		d.Set("type_num", flattenOfflineKeyAnalysisTypeNum(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOfflineKeyAnalysisLargestKeyPrefixes(resp interface{}) []map[string]interface{} {
	curArray := utils.PathSearch("largest_key_prefixes", resp, make([]interface{}, 0)).([]interface{})
	result := make([]map[string]interface{}, len(curArray))
	for i, item := range curArray {
		result[i] = map[string]interface{}{
			"key_prefix": utils.PathSearch("key_prefix", item, nil),
			"type":       utils.PathSearch("type", item, nil),
			"bytes":      utils.PathSearch("bytes", item, nil),
			"num":        utils.PathSearch("num", item, nil),
		}
	}

	return result
}

func flattenOfflineKeyAnalysisLargestKeys(resp interface{}) []map[string]interface{} {
	curArray := utils.PathSearch("largest_keys", resp, make([]interface{}, 0)).([]interface{})
	result := make([]map[string]interface{}, 0, len(curArray))
	for i, item := range curArray {
		result[i] = map[string]interface{}{
			"key":         utils.PathSearch("key", item, nil),
			"type":        utils.PathSearch("type", item, nil),
			"bytes":       utils.PathSearch("bytes", item, nil),
			"num_of_elem": utils.PathSearch("num_of_elem", item, nil),
		}
	}

	return result
}

func flattenOfflineKeyAnalysisTypeBytes(resp interface{}) []map[string]interface{} {
	curArray := utils.PathSearch("type_bytes", resp, make([]interface{}, 0)).([]interface{})
	result := make([]map[string]interface{}, 0, len(curArray))
	for i, item := range curArray {
		result[i] = map[string]interface{}{
			"type":  utils.PathSearch("type", item, nil),
			"bytes": utils.PathSearch("bytes", item, nil),
		}
	}

	return result
}

func flattenOfflineKeyAnalysisTypeNum(resp interface{}) []map[string]interface{} {
	curArray := utils.PathSearch("type_num", resp, make([]interface{}, 0)).([]interface{})
	result := make([]map[string]interface{}, 0, len(curArray))
	for i, item := range curArray {
		result[i] = map[string]interface{}{
			"type": utils.PathSearch("type", item, nil),
			"num":  utils.PathSearch("num", item, nil),
		}
	}

	return result
}

func resourceOfflineKeyAnalysisUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOfflineKeyAnalysisDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/offline/key-analysis/{task_id}"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{task_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	retryFunc := func() (interface{}, bool, error) {
		r, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleOperationError(err)
		return r, retry, err
	}
	deleteResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"RUNNING"},
		WaitPending:  []string{"PENDING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DCS.4232"),
			"error deleting DCS offline key analysis")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	taskId := utils.PathSearch("task_id", deleteRespBody, "").(string)
	if taskId == "" {
		return diag.Errorf("error deleting DCS offline key analysis: task_id is not found in API response")
	}

	if err = checkOfflineKeyAnalysisJobFinish(ctx, client, d, "deleted", d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for DCS offline key analysis task (%s) to be deleted: %s", taskId, err)
	}

	return nil
}

func checkOfflineKeyAnalysisJobFinish(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	target string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{target},
		Refresh:      offlineKeyAnalysisJobStatusRefreshFunc(client, d.Get("instance_id").(string), d.Id()),
		Timeout:      timeout,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for DCS offline key analysis task (%s) to be completed: %s", d.Id(), err)
	}
	return nil
}

func offlineKeyAnalysisJobStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl = "v2/{project_id}/instances/{instance_id}/offline/key-analysis"
		)

		getPath := client.Endpoint + httpUrl
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

		getResp, err := pagination.ListAllItems(
			client,
			"offset",
			getPath,
			&pagination.QueryOpts{MarkerField: ""})
		if err != nil {
			return nil, "failed", err
		}

		getRespJson, err := json.Marshal(getResp)
		if err != nil {
			return nil, "failed", err
		}
		var getRespBody interface{}
		err = json.Unmarshal(getRespJson, &getRespBody)
		if err != nil {
			return nil, "failed", err
		}

		status := utils.PathSearch(fmt.Sprintf("records[?id=='%s']|[0].status", taskId), getRespBody, "").(string)
		if status == "" {
			return getRespBody, "deleted", nil
		}
		if status == "success" {
			return getRespBody, status, nil
		}
		if status == "failed" {
			return getRespBody, status, fmt.Errorf("error getting the status of the offline key analysis task(%s)", taskId)
		}

		return getRespBody, "pending", nil
	}
}

func resourceOfflineKeyAnalysisImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<task_id>")
	}

	instanceId := parts[0]
	taskId := parts[1]

	d.SetId(taskId)
	if err := d.Set("instance_id", instanceId); err != nil {
		return nil, fmt.Errorf("error setting instance_id: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
