package css

import (
	"context"
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

const (
	ClusterScanTaskRunning   = "RUNNING"
	ClusterScanTaskCompleted = "COMPLETED"
	ClusterScanTaskFailed    = "FAILED"
)

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/ai-ops
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/ai-ops
// @API CSS DELETE /v1.0/{project_id}/clusters/{cluster_id}/ai-ops/{aiops_id}
func ResourceScanTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScanTaskCreate,
		ReadContext:   resourceScanTaskRead,
		DeleteContext: resourceScanTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceScanTaskImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"alarm": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"level": {
							Type:     schema.TypeString,
							Required: true,
						},
						"smn_topic": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"smn_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"smn_fail_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"summary": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"high_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"medium_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"suggestion_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"task_risks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"risk": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"suggestion": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceScanTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)

	// createScanTask: create a scan task.
	var (
		createScanTaskUrl     = "v1.0/{project_id}/clusters/{cluster_id}/ai-ops"
		createScanTaskProduct = "css"
	)

	createScanTaskClient, err := conf.NewServiceClient(createScanTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}
	createScanTaskPath := createScanTaskClient.Endpoint + createScanTaskUrl
	createScanTaskPath = strings.ReplaceAll(createScanTaskPath, "{project_id}", createScanTaskClient.ProjectID)
	createScanTaskPath = strings.ReplaceAll(createScanTaskPath, "{cluster_id}", clusterID)

	createScanTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createScanTaskOpt.JSONBody = utils.RemoveNil(buildcreateScanTaskBodyParams(d))
	_, err = createScanTaskClient.Request("POST", createScanTaskPath, &createScanTaskOpt)
	if err != nil {
		return diag.Errorf("error creating CSS cluster scan task: %s", err)
	}

	id, err := refreshScanTaskID(createScanTaskClient, clusterID, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{ClusterScanTaskRunning},
		Target:       []string{ClusterScanTaskCompleted, ClusterScanTaskFailed},
		Refresh:      scanTaskStateRefreshFunc(createScanTaskClient, clusterID, d.Get("name").(string)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        60 * time.Second,
		PollInterval: 60 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the cluster scan task to become available: %s", err)
	}

	return resourceScanTaskRead(ctx, d, meta)
}

func resourceScanTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	scanTask, err := getScanTaskByName(client, d.Get("cluster_id").(string), d.Get("name").(string))
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, "error querying CSS cluster scan task")
	}

	createdAt := utils.PathSearch("create_time", scanTask, float64(0)).(float64) / 1000
	status := int(utils.PathSearch("status", scanTask, float64(0)).(float64))
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", scanTask, nil)),
		d.Set("description", utils.PathSearch("desc", scanTask, nil)),
		d.Set("status", convertClusterScanTaskStatus(status)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(createdAt), false)),
		d.Set("smn_status", utils.PathSearch("smn_status", scanTask, nil)),
		d.Set("smn_fail_reason", utils.PathSearch("smn_fail_reason", scanTask, nil)),
		d.Set("summary", flattenScanTaskSummaryResponse(scanTask)),
		d.Set("task_risks", flattenScanTaskRisksResponse(scanTask)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceScanTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)

	// deleteScanTask: delete a scan task.
	var (
		deleteScanTaskUrl     = "v1.0/{project_id}/clusters/{cluster_id}/ai-ops/{aiops_id}"
		deleteScanTaskProduct = "css"
	)

	deleteScanTaskClient, err := conf.NewServiceClient(deleteScanTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}
	deleteScanTaskPath := deleteScanTaskClient.Endpoint + deleteScanTaskUrl
	deleteScanTaskPath = strings.ReplaceAll(deleteScanTaskPath, "{project_id}", deleteScanTaskClient.ProjectID)
	deleteScanTaskPath = strings.ReplaceAll(deleteScanTaskPath, "{cluster_id}", clusterID)
	deleteScanTaskPath = strings.ReplaceAll(deleteScanTaskPath, "{aiops_id}", d.Id())

	deleteScanTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteScanTaskClient.Request("DELETE", deleteScanTaskPath, &deleteScanTaskOpt)
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, "error deleting CSS cluster scan task")
	}

	return nil
}

func buildcreateScanTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"alarm":       buildcreateScanTaskAlarmBodyParams(d.Get("alarm")),
	}
	return bodyParams
}

func buildcreateScanTaskAlarmBodyParams(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok && len(rawArray) == 1 {
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"level":     raw["level"],
			"smn_topic": raw["smn_topic"],
		}
		return params
	}
	return nil
}

func refreshScanTaskID(client *golangsdk.ServiceClient, clusterId, taskName string) (string, error) {
	scanTask, err := getScanTaskByName(client, clusterId, taskName)
	if err != nil {
		return "", err
	}
	id := utils.PathSearch("id", scanTask, "").(string)

	return id, nil
}

func scanTaskStateRefreshFunc(client *golangsdk.ServiceClient, clusterId, taskName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		scanTask, err := getScanTaskByName(client, clusterId, taskName)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				return "RESOURCE_NOT_FOUND", ClusterScanTaskCompleted, nil
			}
			return scanTask, "ERROR", err
		}
		status := utils.PathSearch("status", scanTask, float64(0)).(float64)
		return scanTask, convertClusterScanTaskStatus(int(status)), nil
	}
}

func getScanTaskByName(client *golangsdk.ServiceClient, clusterId, taskName string) (interface{}, error) {
	getScanTaskHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/ai-ops"

	getScanTaskPath := client.Endpoint + getScanTaskHttpUrl
	getScanTaskPath = strings.ReplaceAll(getScanTaskPath, "{project_id}", client.ProjectID)
	getScanTaskPath = strings.ReplaceAll(getScanTaskPath, "{cluster_id}", clusterId)

	getScanTaskPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	currentTotal := 0
	for {
		currentPath := fmt.Sprintf("%s?limit=10&start=%d", getScanTaskPath, currentTotal)
		getScanTaskResp, err := client.Request("GET", currentPath, &getScanTaskPathOpt)
		if err != nil {
			return getScanTaskResp, err
		}
		getScanTaskRespBody, err := utils.FlattenResponse(getScanTaskResp)
		if err != nil {
			return nil, err
		}
		scanTasks := utils.PathSearch("aiops_list", getScanTaskRespBody, make([]interface{}, 0)).([]interface{})
		findAiopsList := fmt.Sprintf("aiops_list | [?name=='%s'] | [0]", taskName)
		scanTask := utils.PathSearch(findAiopsList, getScanTaskRespBody, nil)
		if scanTask != nil {
			return scanTask, nil
		}
		total := utils.PathSearch("total_size", getScanTaskRespBody, float64(0)).(float64)
		currentTotal += len(scanTasks)
		if float64(currentTotal) == total {
			break
		}
	}
	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v1.0/{project_id}/clusters/{cluster_id}/ai-ops",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the scan task (%s) does not exist", taskName)),
		},
	}
}

func flattenScanTaskSummaryResponse(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("summary", resp, nil)
	if curJson == nil {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"high_num":       int(utils.PathSearch("high", curJson, float64(0)).(float64)),
			"medium_num":     int(utils.PathSearch("medium", curJson, float64(0)).(float64)),
			"suggestion_num": int(utils.PathSearch("suggestion", curJson, float64(0)).(float64)),
		},
	}

	return rst
}

func flattenScanTaskRisksResponse(resp interface{}) []interface{} {
	curJson := utils.PathSearch("task_risks", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"risk":        utils.PathSearch("riskType", v, nil),
			"level":       utils.PathSearch("level", v, nil),
			"description": utils.PathSearch("desc", v, nil),
			"suggestion":  utils.PathSearch("suggestion", v, nil),
		})
	}
	return rst
}

func convertClusterScanTaskStatus(status int) string {
	state := ClusterScanTaskRunning
	if status == 200 {
		state = ClusterScanTaskCompleted
	}
	if status == 300 {
		state = ClusterScanTaskFailed
	}

	return state
}

func resourceScanTaskImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <cluster_id>/<name>")
	}

	d.Set("cluster_id", parts[0])
	d.Set("name", parts[1])

	cssClient, err := conf.NewServiceClient("css", region)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS client: %s", err)
	}
	scanTask, err := getScanTaskByName(cssClient, parts[0], parts[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	id := utils.PathSearch("id", scanTask, "")
	d.SetId(id.(string))

	return []*schema.ResourceData{d}, nil
}
