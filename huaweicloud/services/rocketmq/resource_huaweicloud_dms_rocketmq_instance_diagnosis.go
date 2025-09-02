package rocketmq

import (
	"context"
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

// @API RocketMQ POST /v2/{engine}/{project_id}/instances/{instance_id}/diagnosis
// @API RocketMQ GET /v2/{engine}/{project_id}/diagnosis/{report_id}
// @API RocketMQ DELETE /v2/{engine}/{project_id}/instances/{instance_id}/diagnosis
func ResourceInstanceDiagnosis() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceDiagnosisCreate,
		ReadContext:   resourceInstanceDiagnosisRead,
		DeleteContext: resourceInstanceDiagnosisDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which to create the resource.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the RocketMQ instance to be diagnosed.`,
			},
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the consumer group to be diagnosed.`,
			},
			"node_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of node IDs to be diagnosed.`,
			},

			// Attributes
			"report_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the diagnosis report.`,
			},
			"consumer_nums": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of consumers.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the diagnosis task.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the diagnosis report.`,
			},
			"abnormal_item_sum": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of abnormal items.`,
			},
			"faulted_node_sum": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of abnormal nodes.`,
			},
			"online": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the consumer group is online.`,
			},
			"message_accumulation": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of accumulated messages.`,
			},
			"subscription_consistency": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the subscription is consistent.`,
			},
			"subscriptions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of subscribers.`,
			},
			"diagnosis_node_reports": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of diagnosis node report.`,
			},
		},
	}
}

func buildCreateDiagnosisBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{}

	if groupName, ok := d.GetOk("group_name"); ok {
		bodyParams["group_name"] = groupName
	}

	if nodeIdList, ok := d.GetOk("node_ids"); ok {
		bodyParams["node_id_list"] = utils.ExpandToStringList(nodeIdList.([]interface{}))
	}

	return bodyParams
}

func resourceInstanceDiagnosisCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	httpUrl := "v2/{engine}/{project_id}/instances/{instance_id}/diagnosis"
	httpUrl = strings.ReplaceAll(httpUrl, "{engine}", "rocketmq")
	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{instance_id}", d.Get("instance_id").(string))

	createPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateDiagnosisBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating RocketMQ instance diagnosis task: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	reportId := utils.PathSearch("report_id", respBody, "").(string)
	if reportId == "" {
		return diag.Errorf("unable to find report ID from response")
	}
	d.SetId(reportId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"diagnosing"},
		Target:  []string{"finished"},
		Refresh: diagnosisTaskStateRefreshFunc(client, reportId),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for RocketMQ diagnosis task (%s) to complete: %s", reportId, err)
	}

	return resourceInstanceDiagnosisRead(ctx, d, meta)
}

func diagnosisTaskStateRefreshFunc(client *golangsdk.ServiceClient, reportId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		httpUrl := "v2/{engine}/{project_id}/diagnosis/{report_id}"
		httpUrl = strings.ReplaceAll(httpUrl, "{engine}", "rocketmq")
		httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
		httpUrl = strings.ReplaceAll(httpUrl, "{report_id}", reportId)

		queryPath := client.Endpoint + httpUrl

		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}

		requestResp, err := client.Request("GET", queryPath, &opt)
		if err != nil {
			return nil, "", err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("status", respBody, "").(string)

		return respBody, status, nil
	}
}

func resourceInstanceDiagnosisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	reportId := d.Id()

	httpUrl := "v2/{engine}/{project_id}/diagnosis/{report_id}"
	httpUrl = strings.ReplaceAll(httpUrl, "{engine}", "rocketmq")
	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{report_id}", reportId)

	queryPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", queryPath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RocketMQ diagnosis report")
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	createdAt := utils.PathSearch("creat_at", respBody, float64(0)).(float64)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("report_id", utils.PathSearch("report_id", respBody, nil)),
		d.Set("group_name", utils.PathSearch("group_name", respBody, nil)),
		d.Set("consumer_nums", utils.PathSearch("consumer_nums", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(createdAt)/1000, false)),
		d.Set("abnormal_item_sum", utils.PathSearch("abnormal_item_sum", respBody, nil)),
		d.Set("faulted_node_sum", utils.PathSearch("faulted_node_sum", respBody, nil)),
		d.Set("online", utils.PathSearch("online", respBody, nil)),
		d.Set("message_accumulation", utils.PathSearch("message_accumulation", respBody, nil)),
		d.Set("subscription_consistency", utils.PathSearch("subscription_consistency", respBody, nil)),
		d.Set("subscriptions", utils.ExpandToStringList(
			utils.PathSearch("subscriptions", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("diagnosis_node_reports", utils.ExpandToStringList(
			utils.PathSearch("diagnosis_node_report_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceInstanceDiagnosisDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	reportId := d.Id()

	httpUrl := "v2/{engine}/{project_id}/instances/{instance_id}/diagnosis"
	httpUrl = strings.ReplaceAll(httpUrl, "{engine}", "rocketmq")
	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{instance_id}", instanceId)

	deletePath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"report_id_list": []string{reportId},
		},
	}

	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return diag.Errorf("error deleting RocketMQ diagnosis report: %s", err)
	}

	return nil
}
