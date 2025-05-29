// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product COC
// ---------------------------------------------------------------

package coc

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

var scriptOrderNotFoundErrCodes = []string{
	"COC.00040709", // Script not found
}

var scriptExecuteNonUpdatableParams = []string{
	"script_id", "instance_id", "timeout", "execute_user", "parameters", "parameters.*.name", "parameters.*.value", "is_sync",
}

// @API COC POST /v1/job/scripts/{script_uuid}
// @API COC GET /v1/job/script/orders/{execute_uuid}
// @API COC PUT /v1/job/script/orders/{execute_uuid}/operation
// @API COC POST /v1/external/resources/sync
// @API COC GET /v1/external/resources
func ResourceScriptExecute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScriptExecuteCreate,
		ReadContext:   resourceScriptExecuteRead,
		UpdateContext: resourceScriptExecuteUpdate,
		DeleteContext: resourceScriptExecuteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(scriptExecuteNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"script_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"execute_user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
								return oldValue == defaultSensitiveValue
							},
						},
					},
				},
			},
			"is_sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			// attributes
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"script_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finished_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildExecuteParamsBody(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"param_order": i + 1, // param_order starts counting from 1
				"param_name":  raw["name"],
				"param_value": raw["value"],
			}
		}
		return params
	}

	return nil
}

func buildExecuteParamReqBody(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"resourceful":   true,
		"success_rate":  100,
		"timeout":       d.Get("timeout"),
		"execute_user":  d.Get("execute_user"),
		"script_params": buildExecuteParamsBody(d.Get("parameters")),
	}
}

func buildExecuteBatchesReqBody(info interface{}) []map[string]interface{} {
	targetInstances := []map[string]interface{}{
		{
			"resource_id": utils.PathSearch("resource_id", info, nil),
			"agent_sn":    utils.PathSearch("agent_id", info, nil),
			"project_id":  utils.PathSearch("project_id", info, nil),
			"region_id":   utils.PathSearch("region_id", info, nil),
		},
	}

	return []map[string]interface{}{
		{
			"batch_index":       1,
			"rotation_strategy": "PAUSE",
			"target_instances":  targetInstances,
		},
	}
}

func buildCreateExecuteBodyParams(d *schema.ResourceData, extra interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"execute_param":   buildExecuteParamReqBody(d),
		"execute_batches": buildExecuteBatchesReqBody(extra),
	}
	return bodyParams
}

func resourceScriptExecuteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "coc"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	// sync the ECS instance to get information about UniAgent
	if d.Get("is_sync").(bool) {
		if err = syncResourceInfo(client); err != nil {
			return diag.FromErr(err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"online"},
		Refresh:      doGetResources(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 15 * time.Second,
	}
	info, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for ECS instance(%s) to be online in COC: %s", instanceID, err)
	}

	createExecuteHttpUrl := fmt.Sprintf("v1/job/scripts/%s", d.Get("script_id"))
	createExecutePath := client.Endpoint + createExecuteHttpUrl

	createExecuteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createExecuteOpt.JSONBody = utils.RemoveNil(buildCreateExecuteBodyParams(d, info))
	createExecuteResp, err := client.Request("POST", createExecutePath, &createExecuteOpt)
	if err != nil {
		return diag.Errorf("error executing COC script: %s", err)
	}

	createExecuteRespBody, err := utils.FlattenResponse(createExecuteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	ticketID := utils.PathSearch("data", createExecuteRespBody, "").(string)
	if ticketID == "" {
		return diag.Errorf("unable to find the executing COC script ID from the API response")
	}

	d.SetId(ticketID)

	// waiting the execution status of COC script
	stateConf = &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"exited"},
		Refresh:      refreshGetExecutionTicketDetail(client, ticketID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        15 * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		// skip the timeout error, and return error when the result is ABNORMAL
		if _, ok := err.(*resource.UnexpectedStateError); ok {
			return diag.Errorf("error executing COC script: %s", err)
		}
	}

	return resourceScriptExecuteRead(ctx, d, meta)
}

func refreshGetExecutionTicketDetail(client *golangsdk.ServiceClient, ticketID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ticketDetail, err := getExecutionTicketDetail(client, ticketID)
		if err != nil {
			return nil, "error", err
		}

		status := utils.PathSearch("data.status", ticketDetail, "")
		ticketStatus := status.(string)

		if ticketStatus == "PROCESSING" {
			return ticketDetail, "pending", nil
		} else if ticketStatus == "ABNORMAL" {
			return ticketDetail, "error", nil
		}
		return ticketDetail, "exited", nil
	}
}

func getExecutionTicketDetail(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	getExecutionTicketHttpUrl := "v1/job/script/orders/{id}"
	getExecutionTicketPath := client.Endpoint + getExecutionTicketHttpUrl
	getExecutionTicketPath = strings.ReplaceAll(getExecutionTicketPath, "{id}", id)

	getExecutionTicketOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getExecutionTicketResp, err := client.Request("GET", getExecutionTicketPath, &getExecutionTicketOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getExecutionTicketResp)
}

func resourceScriptExecuteRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "coc"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	ticketID := d.Id()
	ticketDetail, err := getExecutionTicketDetail(client, ticketID)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			scriptOrderNotFoundErrCodes...), "COC script execute")
	}

	mErr := multierror.Append(nil,
		d.Set("status", utils.PathSearch("data.status", ticketDetail, nil)),
		d.Set("script_id", utils.PathSearch("data.properties.script_uuid", ticketDetail, nil)),
		d.Set("script_name", utils.PathSearch("data.properties.script_name", ticketDetail, nil)),
		d.Set("timeout", utils.PathSearch("data.properties.execute_param.timeout", ticketDetail, nil)),
		d.Set("execute_user", utils.PathSearch("data.properties.execute_param.execute_user", ticketDetail, nil)),
		d.Set("created_at", flattenScriptTimeStamp(ticketDetail, "data.gmt_created")),
		d.Set("finished_at", flattenScriptTimeStamp(ticketDetail, "data.gmt_finished")),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting COC script execute fields: %s", err)
	}

	return nil
}

// now, the value of script_params in API response is the default value, not the value when executing
// so we donot set `parameters` until COC service fixed this bug.
// nolint: unused
func flattenScriptExecuteParams(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("data.properties.execute_param.script_params", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"name":  utils.PathSearch("param_name", v, nil),
			"value": utils.PathSearch("param_value", v, nil),
		}
	}
	return rst
}

func resourceScriptExecuteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceScriptExecuteDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "coc"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	status := d.Get("status").(string)
	exitList := []string{
		"FINISHED", "ROLLBACKED", "CANCELED", "ERROR", "ABNORMAL",
	}
	if utils.StrSliceContains(exitList, status) {
		return nil
	}

	// cancel the ticket when it in PROCESSING and PAUSED status
	operationExecuteHttpUrl := fmt.Sprintf("v1/job/script/orders/%s/operation", d.Id())
	operationExecutePath := client.Endpoint + operationExecuteHttpUrl

	operationExecuteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         map[string]string{"operation_type": "CANCEL_ORDER"},
	}

	_, err = client.Request("PUT", operationExecutePath, &operationExecuteOpt)
	if err != nil {
		return diag.Errorf("error canceling COC script execution: %s", err)
	}
	return nil
}

func syncResourceInfo(client *golangsdk.ServiceClient) error {
	syncResourceHttpUrl := "v1/external/resources/sync"
	syncResourcePath := client.Endpoint + syncResourceHttpUrl

	syncResourceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]string{
			"provider": "ecs",
			"type":     "cloudservers",
		},
	}

	_, err := client.Request("POST", syncResourcePath, &syncResourceOpt)
	if err != nil {
		return fmt.Errorf("error syncing ECS instances: %s", err)
	}

	return nil
}

func doGetResources(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		listResourceInfoHttpUrl := fmt.Sprintf("v1/external/resources?provider=%s&type=%s&limit=10&resource_id_list=%s",
			"ecs", "cloudservers", instanceID)
		listResourceInfoPath := client.Endpoint + listResourceInfoHttpUrl

		listResourceInfoOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}

		respRaw, err := client.Request("GET", listResourceInfoPath, &listResourceInfoOpt)
		if err != nil {
			return nil, "error", fmt.Errorf("error syncing ECS instances: %s", err)
		}

		respBody, err := utils.FlattenResponse(respRaw)
		if err != nil {
			return nil, "error", err
		}

		item := utils.PathSearch(fmt.Sprintf("data[?resource_id=='%s']|[0]", instanceID), respBody, nil)
		if item == nil {
			return "", "pending", nil
		}

		agentID := utils.PathSearch("agent_id", item, nil)
		agentStatus := utils.PathSearch("agent_state", item, "uninstalled")
		if agentID == nil || agentStatus != "ONLINE" {
			return item, "pending", nil
		}

		return item, "online", nil
	}
}
