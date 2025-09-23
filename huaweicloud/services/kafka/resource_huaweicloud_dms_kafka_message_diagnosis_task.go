package kafka

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

// @API Kafka DELETE /v2/{project_id}/kafka/instances/{instance_id}/message-diagnosis-tasks
// @API Kafka POST /v2/{project_id}/kafka/instances/{instance_id}/message-diagnosis-tasks
// @API Kafka GET /v2/{project_id}/kafka/instances/{instance_id}/message-diagnosis/{report_id}
// @API Kafka GET /v2/{project_id}/kafka/instances/{instance_id}/diagnosis-check
// @API Kafka GET /v2/{project_id}/kafka/instances/{instance_id}/message-diagnosis-tasks
func ResourceDmsKafkaMessageDiagnosisTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaMessageDiagnosisTaskCreate,
		ReadContext:   resourceDmsKafkaMessageDiagnosisTaskRead,
		DeleteContext: resourceDmsKafkaMessageDiagnosisTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMessageDiagnosisTaskImportState,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"accumulated_partitions": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"diagnosis_dimension_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"abnormal_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"failed_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"diagnosis_item_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"result": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cause_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     resourceSchemeDiagnosisConclusionEntity(),
									},
									"advice_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     resourceSchemeDiagnosisConclusionEntity(),
									},
									"partitions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"failed_partitions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"broker_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSchemeDiagnosisConclusionEntity() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"params": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func precheckMessageDiagnosisTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	precheckHttpUrl := "v2/{project_id}/kafka/instances/{instance_id}/diagnosis-check?group={group}&topic={topic}"
	precheckPath := client.Endpoint + precheckHttpUrl
	precheckPath = strings.ReplaceAll(precheckPath, "{project_id}", client.ProjectID)
	precheckPath = strings.ReplaceAll(precheckPath, "{instance_id}", d.Get("instance_id").(string))
	precheckPath = strings.ReplaceAll(precheckPath, "{group}", d.Get("group_name").(string))
	precheckPath = strings.ReplaceAll(precheckPath, "{topic}", d.Get("topic_name").(string))
	precheckOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("GET", precheckPath, &precheckOpt)
	if err != nil {
		return fmt.Errorf("error prechecking Kafka message diagnosis task: %s", err)
	}

	return nil
}

func resourceDmsKafkaMessageDiagnosisTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	// precheck
	err = precheckMessageDiagnosisTask(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	createHttpUrl := "v2/{project_id}/kafka/instances/{instance_id}/message-diagnosis-tasks"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"group_name": d.Get("group_name").(string),
			"topic_name": d.Get("topic_name").(string),
		},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Kafka message diagnosis task: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}

	reportID := utils.PathSearch("report_id", createRespBody, "").(string)
	if reportID == "" {
		return diag.Errorf("unable to find report ID from API response")
	}

	d.SetId(reportID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      kafkaMessageDiagnosisTaskStateRefreshFunc(client, instanceID, reportID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the message diagnosis task(%s) to be done: %s", d.Id(), err)
	}

	return resourceDmsKafkaMessageDiagnosisTaskRead(ctx, d, cfg)
}

func kafkaMessageDiagnosisTaskStateRefreshFunc(client *golangsdk.ServiceClient, instanceID, reportID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		task, err := filterKafkaMessageDiagnosisTaskFromList(client, instanceID, reportID)
		if err != nil {
			return nil, "QUERY ERROR", err
		}

		status := utils.PathSearch("status", task, "").(string)
		switch status {
		case "diagnosing":
			return task, "PENDING", nil
		case "finished":
			return task, "SUCCESS", nil
		}

		return task, status, nil
	}
}

func filterKafkaMessageDiagnosisTaskFromList(client *golangsdk.ServiceClient, instanceID, reportID string) (interface{}, error) {
	listHttpUrl := "v2/{project_id}/kafka/instances/{instance_id}/message-diagnosis-tasks"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pageLimit is `10`
	listPath += fmt.Sprintf("?limit=%d", pageLimit)
	offset := 0
	for {
		currentPath := listPath + fmt.Sprintf("&offset=%d", offset)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving the message diagnosis tasks list: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, fmt.Errorf("error flattening the message diagnosis tasks list: %s", err)
		}

		searchPath := fmt.Sprintf("report_list[?report_id=='%s']|[0]", reportID)
		result := utils.PathSearch(searchPath, listRespBody, nil)
		if result != nil {
			return result, nil
		}

		// `total_num` means the number of all reports, and type is float64.
		offset += pageLimit
		total := utils.PathSearch("total_num", listRespBody, float64(0))
		if int(total.(float64)) <= offset {
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Body: []byte(fmt.Sprintf("unable to find task(%s) from API response", reportID)),
				}}
		}
	}
}

func resourceDmsKafkaMessageDiagnosisTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	reportID := d.Id()

	reportInfo, err := filterKafkaMessageDiagnosisTaskFromList(client, instanceID, reportID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the message diagnosis task report")
	}

	reportDetail, err := GetKafkaMessageDiagnosisTaskReport(client, instanceID, reportID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the message diagnosis task report")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("group_name", utils.PathSearch("group_name", reportInfo, nil)),
		d.Set("topic_name", utils.PathSearch("topic_name", reportInfo, nil)),
		d.Set("status", utils.PathSearch("status", reportInfo, nil)),
		d.Set("begin_time", utils.PathSearch("begin_time", reportInfo, nil)),
		d.Set("end_time", utils.PathSearch("end_time", reportInfo, nil)),
		d.Set("accumulated_partitions", utils.PathSearch("accumulated_partitions", reportInfo, nil)),
		d.Set("diagnosis_dimension_list", flattenDiagnosisDimensionList(
			utils.PathSearch("diagnosis_dimension_list", reportDetail, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetKafkaMessageDiagnosisTaskReport(client *golangsdk.ServiceClient, instanceID, reportID string) (interface{}, error) {
	getHttpUrl := "v2/{project_id}/kafka/instances/{instance_id}/message-diagnosis/{report_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)
	getPath = strings.ReplaceAll(getPath, "{report_id}", reportID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the message diagnosis task report: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening the message diagnosis task report: %s", err)
	}

	return getRespBody, nil
}

func flattenDiagnosisDimensionList(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"name":         utils.PathSearch("name", params, nil),
			"abnormal_num": utils.PathSearch("abnormal_num", params, nil),
			"failed_num":   utils.PathSearch("failed_num", params, nil),
			"diagnosis_item_list": flattenDiagnosisDimensionListDiagnosisItemList(
				utils.PathSearch("diagnosis_item_list", params, make([]interface{}, 0)).([]interface{})),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenDiagnosisDimensionListDiagnosisItemList(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"name":   utils.PathSearch("name", params, nil),
			"result": utils.PathSearch("result", params, nil),
			"cause_ids": flattenDiagnosisItemListDiagnosisConclusionEntity(
				utils.PathSearch("cause_ids", params, make([]interface{}, 0)).([]interface{})),
			"advice_ids": flattenDiagnosisItemListDiagnosisConclusionEntity(
				utils.PathSearch("advice_ids", params, make([]interface{}, 0)).([]interface{})),
			"partitions":        utils.PathSearch("partitions", params, nil),
			"failed_partitions": utils.PathSearch("failed_partitions", params, nil),
			"broker_ids":        utils.PathSearch("broker_ids", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenDiagnosisItemListDiagnosisConclusionEntity(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"id":     utils.PathSearch("id", params, nil),
			"params": utils.PathSearch("params", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func resourceDmsKafkaMessageDiagnosisTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/kafka/instances/{instance_id}/message-diagnosis-tasks"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"report_id_list": []string{d.Id()},
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting message diagnosis task report")
	}

	return nil
}

func resourceMessageDiagnosisTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<report_id>")
	}

	d.Set("instance_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
