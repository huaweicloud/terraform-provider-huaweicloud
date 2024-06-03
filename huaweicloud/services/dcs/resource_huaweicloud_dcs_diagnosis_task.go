package dcs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS POST /v2/{project_id}/instances/{instance_id}/diagnosis
// @API DCS GET /v2/{project_id}/diagnosis/{report_id}
// @API DCS GET /v2/{project_id}/instances/{instance_id}/diagnosis
// @API DCS GET /v2/{project_id}/instances/{instance_id}
// @API DCS DELETE /v2/{project_id}/instances/{instance_id}/diagnosis
func ResourceDiagnosisTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDiagnosisTaskCreate,
		ReadContext:   resourceDiagnosisTaskRead,
		DeleteContext: resourceDiagnosisTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDianosisReportImportState,
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
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the DCS instance.",
			},
			"begin_time": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressTimeDiffs,
				Description:      `Specifies the start time of the diagnosis task, in RFC3339 format.`,
			},
			"end_time": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressTimeDiffs,
				Description:      `Specifies the end time of the diagnosis task, in RFC3339 format.`,
			},
			"node_ip_list": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the IP addresses of diagnosed nodes. By default, all nodes are diagnosed.`,
			},
			"abnormal_item_sum": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total number of abnormal diagnosis items.`,
			},
			"failed_item_sum": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total number of failed diagnosis items.`,
			},
			"diagnosis_node_report_list": {
				Type:        schema.TypeList,
				Elem:        diagnosisNodeReportSchema(),
				Computed:    true,
				Description: `Indicates the list of node diagnosis report`,
			},
		},
	}
}

func diagnosisNodeReportSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the node IP address.`,
			},
			"az_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the code of the AZ where the node is.`,
			},
			"group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the shard where the node is.`,
			},
			"abnormal_sum": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total number of abnormal diagnosis items.`,
			},
			"failed_sum": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total number of failed diagnosis items.`,
			},
			"is_faulted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the node is faulted.`,
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the node role.`,
			},
			"diagnosis_dimension_list": {
				Type:        schema.TypeList,
				Elem:        diagnosisDimensionSchema(),
				Computed:    true,
				Description: `Indicates the diagnosis dimension list.`,
			},
			"command_time_taken_list": {
				Type:        schema.TypeList,
				Elem:        commandTimeTakenSchema(),
				Computed:    true,
				Description: `Indicates the command execution duration list.`,
			},
		},
	}
	return &sc
}

func diagnosisDimensionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the diagnosis dimension name.`,
			},
			"abnormal_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total number of abnormal diagnosis items.`,
			},
			"failed_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total number of failed diagnosis items.`,
			},
			"diagnosis_item_list": {
				Type:        schema.TypeList,
				Elem:        diagnosisItemSchema(),
				Computed:    true,
				Description: `Indicates the diagnosis items.`,
			},
		},
	}
	return &sc
}

func commandTimeTakenSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"total_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total number of times that commands are executed.`,
			},
			"total_usec_sum": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Indicates the total duration of command execution.`,
			},
			"result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the command execution latency result.`,
			},
			"error_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the error code for the command time taken.`,
			},
			"command_list": {
				Type:        schema.TypeList,
				Elem:        commandSchema(),
				Computed:    true,
				Description: `Indicates the command execution latency statistics.`,
			},
		},
	}
	return &sc
}

func diagnosisItemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the diagnosis item name.`,
			},
			"cause_ids": {
				Type:        schema.TypeList,
				Elem:        conclusionItemSchema(),
				Computed:    true,
				Description: `Indicates the list of cause IDs.`,
			},
			"impact_ids": {
				Type:        schema.TypeList,
				Elem:        conclusionItemSchema(),
				Computed:    true,
				Description: `Indicates the list of impact IDs.`,
			},
			"advice_ids": {
				Type:        schema.TypeList,
				Elem:        conclusionItemSchema(),
				Computed:    true,
				Description: `Indicates the list of suggestion IDs.`,
			},
			"result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the diagnosis result.`,
			},
			"error_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the error code for the diagnosis item.`,
			},
		},
	}
	return &sc
}

func conclusionItemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the conclusion ID.`,
			},
			"params": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the conclusion parameters.`,
			},
		},
	}
	return &sc
}

func commandSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"calls_sum": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of calls.`,
			},
			"usec_sum": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Indicates the total time consumed.`,
			},
			"command_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the command name.`,
			},
			"per_usec": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the duration percentage.`,
			},
			"average_usec": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Indicates the average duration of calls.`,
			},
		},
	}
	return &sc
}

func buildCreateDiagnosisTaskBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	beginTime, err := time.Parse(time.RFC3339, d.Get("begin_time").(string))
	if err != nil {
		return nil, err
	}
	beginTimeUnix := beginTime.UnixMicro() / 1000

	endTime, err := time.Parse(time.RFC3339, d.Get("end_time").(string))
	if err != nil {
		return nil, err
	}
	endTimeUnix := endTime.UnixMicro() / 1000

	bodyParams := map[string]interface{}{
		"begin_time":   beginTimeUnix,
		"end_time":     endTimeUnix,
		"node_ip_list": utils.ValueIgnoreEmpty(d.Get("node_ip_list").(*schema.Set).List()),
	}
	return bodyParams, nil
}

func resourceDiagnosisTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// createDiagnosisTask: create diagnosis task
	var (
		createDiagnosisTaskHttpUrl = "v2/{project_id}/instances/{instance_id}/diagnosis"
		createDiagnosisTaskProduct = "dcs"
	)
	createDiagnosisTaskClient, err := cfg.NewServiceClient(createDiagnosisTaskProduct, region)

	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createDiagnosisTaskPath := createDiagnosisTaskClient.Endpoint + createDiagnosisTaskHttpUrl
	createDiagnosisTaskPath = strings.ReplaceAll(createDiagnosisTaskPath, "{project_id}", createDiagnosisTaskClient.ProjectID)
	createDiagnosisTaskPath = strings.ReplaceAll(createDiagnosisTaskPath, "{instance_id}", instanceID)

	createDiagnosisTaskOpt := golangsdk.RequestOpts{KeepResponseBody: true}

	bodyParam, err := buildCreateDiagnosisTaskBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}
	createDiagnosisTaskOpt.JSONBody = utils.RemoveNil(bodyParam)

	retryFunc := func() (interface{}, bool, error) {
		createDiagnosisTaskResp, createErr := createDiagnosisTaskClient.Request("POST", createDiagnosisTaskPath, &createDiagnosisTaskOpt)
		retry, err := handleOperationError(createErr)
		return createDiagnosisTaskResp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(createDiagnosisTaskClient, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error creating diagnosis task: %v", err)
	}

	createDiagnosisTaskRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	reportId := utils.PathSearch("report_id", createDiagnosisTaskRespBody, nil)
	if reportId == nil {
		return diag.Errorf("error creating diagnosis task: report_id is not found in API response")
	}
	d.SetId(reportId.(string))

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"diagnosing"},
		Target:       []string{"finished"},
		Refresh:      diagnosisTaskRefreshFunc(instanceID, d.Id(), createDiagnosisTaskClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for creating diagnosis task to complete: %s", err)
	}

	return resourceDiagnosisTaskRead(ctx, d, meta)
}

func resourceDiagnosisTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// getDiagnosisTask: get diagnosis task
	var (
		mErr                    *multierror.Error
		getDiagnosisTaskHttpUrl = "v2/{project_id}/diagnosis/{report_id}"
		getDiagnosisTaskProduct = "dcs"
	)

	getDiagnosisTaskClient, err := cfg.NewServiceClient(getDiagnosisTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	getDiagnosisTaskPath := getDiagnosisTaskClient.Endpoint + getDiagnosisTaskHttpUrl
	getDiagnosisTaskPath = strings.ReplaceAll(getDiagnosisTaskPath, "{project_id}", getDiagnosisTaskClient.ProjectID)
	getDiagnosisTaskPath = strings.ReplaceAll(getDiagnosisTaskPath, "{report_id}", d.Id())
	getDiagnosisTaskOpt := golangsdk.RequestOpts{KeepResponseBody: true}

	getDiagnosisTaskResp, err := getDiagnosisTaskClient.Request("GET", getDiagnosisTaskPath, &getDiagnosisTaskOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving diagnosis report")
	}

	getDiagnosisTaskRespBody, err := utils.FlattenResponse(getDiagnosisTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// get begin_time and end_time
	report, err := getDiagnosisReport(instanceID, d.Id(), getDiagnosisTaskClient)
	if err != nil {
		return diag.FromErr(err)
	}

	beginTime := utils.PathSearch("begin_time", report, nil)
	endTime := utils.PathSearch("end_time", report, nil)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("begin_time", beginTime),
		d.Set("end_time", endTime),
		d.Set("node_ip_list", getNodeIpAddressList(getDiagnosisTaskRespBody)),
		d.Set("abnormal_item_sum", utils.PathSearch("abnormal_item_sum", getDiagnosisTaskRespBody, nil)),
		d.Set("failed_item_sum", utils.PathSearch("failed_item_sum", getDiagnosisTaskRespBody, nil)),
		d.Set("diagnosis_node_report_list", getDiagnosisNodeReportList(getDiagnosisTaskRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getDiagnosisNodeReportList(resp interface{}) []map[string]interface{} {
	diagnosisNodeReportList := utils.PathSearch("diagnosis_node_report_list", resp, make([]interface{}, 0)).([]interface{})
	reportList := make([]map[string]interface{}, 0, len(diagnosisNodeReportList))
	for _, v := range diagnosisNodeReportList {
		reportList = append(reportList, map[string]interface{}{
			"abnormal_sum":             utils.PathSearch("abnormal_sum", v, nil),
			"az_code":                  utils.PathSearch("az_code", v, nil),
			"failed_sum":               utils.PathSearch("failed_sum", v, nil),
			"group_name":               utils.PathSearch("group_name", v, nil),
			"is_faulted":               utils.PathSearch("is_faulted", v, nil),
			"node_ip":                  utils.PathSearch("node_ip", v, nil),
			"role":                     utils.PathSearch("role", v, nil),
			"command_time_taken_list":  getCommandTimeTakenList(v),
			"diagnosis_dimension_list": utils.PathSearch("diagnosis_dimension_list", v, nil),
		})
	}
	return reportList
}

func getNodeIpAddressList(resp interface{}) []string {
	nodeIpList := utils.PathSearch("diagnosis_node_report_list[*].node_ip", resp, make([]interface{}, 0)).([]interface{})
	nodeIpAddressList := make([]string, len(nodeIpList))
	for i, addressPort := range nodeIpList {
		nodeIpAddressList[i] = strings.Split(addressPort.(string), ":")[0]
	}
	return nodeIpAddressList
}

func getCommandTimeTakenList(nodeReport interface{}) []map[string]interface{} {
	commandTimeTakenStruct := utils.PathSearch("command_time_taken_list", nodeReport, make(map[string]interface{})).(map[string]interface{})
	commandTimeTakenList := make([]map[string]interface{}, 1)
	commandTimeTakenList[0] = commandTimeTakenStruct
	return commandTimeTakenList
}

func resourceDiagnosisTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDiagnosisTask: delete diagnosis task
	var (
		deleteDiagnosisTaskHttpUrl = "v2/{project_id}/instances/{instance_id}/diagnosis"
		deleteDiagnosisTaskProduct = "dcs"
	)

	deleteDiagnosisTaskClient, err := cfg.NewServiceClient(deleteDiagnosisTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deleteDiagnosisTaskPath := deleteDiagnosisTaskClient.Endpoint + deleteDiagnosisTaskHttpUrl
	deleteDiagnosisTaskPath = strings.ReplaceAll(deleteDiagnosisTaskPath, "{project_id}", deleteDiagnosisTaskClient.ProjectID)
	deleteDiagnosisTaskPath = strings.ReplaceAll(deleteDiagnosisTaskPath, "{instance_id}", instanceID)

	deleteDiagnosisTaskOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	deleteDiagnosisTaskOpt.JSONBody = map[string]interface{}{
		"report_id_list": []string{d.Id()},
	}

	_, err = deleteDiagnosisTaskClient.Request("DELETE", deleteDiagnosisTaskPath, &deleteDiagnosisTaskOpt)

	if err != nil {
		return diag.Errorf("error deleting the diagnosis report (%s): %v", d.Id(), err)
	}

	return resourceDiagnosisTaskRead(ctx, d, meta)
}

func diagnosisTaskRefreshFunc(instanceID, reportID string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		report, err := getDiagnosisReport(instanceID, reportID, client)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", report, "").(string)
		return report, status, nil
	}
}

func getDiagnosisReport(instanceID, reportID string, client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getDiagnosisTaskHttpUrl = "v2/{project_id}/instances/{instance_id}/diagnosis"
	)

	getDiagnosisTaskPath := client.Endpoint + getDiagnosisTaskHttpUrl
	getDiagnosisTaskPath = strings.ReplaceAll(getDiagnosisTaskPath, "{project_id}", client.ProjectID)
	getDiagnosisTaskPath = strings.ReplaceAll(getDiagnosisTaskPath, "{instance_id}", instanceID)

	getDiagnosisTaskResp, err := pagination.ListAllItems(
		client,
		"offset",
		getDiagnosisTaskPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, err
	}

	getDiagnosisTaskRespJson, err := json.Marshal(getDiagnosisTaskResp)
	if err != nil {
		return nil, err
	}
	var getDiagnosisTaskRespBody interface{}
	err = json.Unmarshal(getDiagnosisTaskRespJson, &getDiagnosisTaskRespBody)
	if err != nil {
		return nil, err
	}

	report := utils.PathSearch(fmt.Sprintf("diagnosis_report_list[?report_id=='%s']|[0]", reportID), getDiagnosisTaskRespBody, nil)
	if report == nil {
		return nil, fmt.Errorf("unable to find the report %s", reportID)
	}
	return report, nil
}

func resourceDianosisReportImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func suppressTimeDiffs(_, old, new string, _ *schema.ResourceData) bool {
	oldTime, err := time.Parse(time.RFC3339, old)
	if err != nil {
		return false
	}

	newTime, err := time.Parse(time.RFC3339, new)
	if err != nil {
		return false
	}
	layout := "2006-01-02 15:04" // only validate year-month-day hour:minute
	oldTimeWithoutSecond := oldTime.Format(layout)
	newTimeWithoutSecond := newTime.Format(layout)

	return oldTimeWithoutSecond == newTimeWithoutSecond
}
