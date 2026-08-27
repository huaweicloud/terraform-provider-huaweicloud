package coc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var scriptBatchExecuteNonUpdatableParams = []string{
	"script_id",
	"execute_batches",
	"execute_batches.*.batch_index",
	"execute_batches.*.instance_ids",
	"timeout",
	"execute_user",
	"parameters",
	"parameters.*.name",
	"parameters.*.value",
	"is_sync",
	"resource_provider",
	"type",
}

// @API COC POST /v1/external/resources/sync
// @API COC GET /v1/resources
// @API COC POST /v1/job/scripts/{script_uuid}
// @API COC GET /v1/job/script/orders/{execute_uuid}
// @API COC PUT /v1/job/script/orders/{execute_uuid}/operation
func ResourceScriptBatchExecute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScriptBatchExecuteCreate,
		ReadContext:   resourceScriptBatchExecuteRead,
		UpdateContext: resourceScriptBatchExecuteUpdate,
		DeleteContext: resourceScriptBatchExecuteDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(scriptBatchExecuteNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"script_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the COC script.`,
			},
			"execute_batches": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        scriptBatchExecuteBatchesSchema(),
				Description: `The batch information of the target instances.`,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The maximum time to execute the script, in seconds.`,
			},
			"execute_user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user to execute the script.`,
			},

			// Optional parameters.
			"parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        scriptBatchExecuteParametersSchema(),
				Description: `The input parameters of the script.`,
			},
			"is_sync": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to sync data before executing the script.`,
			},
			"resource_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ecs",
				Description: `The resource provider.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "cloudservers",
				Description: `The resource type of the resource provider.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the script execution.`,
			},
			"script_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the script.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The start time of the script execution, in RFC3339 format.`,
			},
			"finished_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end time of the script execution, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func scriptBatchExecuteBatchesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"batch_index": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The batch index.`,
			},
			"instance_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The ID list of the specified resource instances in this batch.`,
			},
		},
	}
}

func scriptBatchExecuteParametersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the parameter.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The value of the parameter.`,
			},
		},
	}
}

func buildScriptBatchExecuteParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"success_rate":  100,
		"timeout":       d.Get("timeout"),
		"execute_user":  d.Get("execute_user"),
		"script_params": buildScriptBatchExecuteScriptParams(d.Get("parameters").([]interface{})),
	}
}

func buildScriptBatchExecuteScriptParams(scriptParams []interface{}) []map[string]interface{} {
	if len(scriptParams) < 1 {
		return nil
	}

	params := make([]map[string]interface{}, 0, len(scriptParams))
	for _, v := range scriptParams {
		params = append(params, map[string]interface{}{
			"param_name":  utils.PathSearch("name", v, ""),
			"param_value": utils.PathSearch("value", v, ""),
		})
	}

	return params
}

func buildScriptBatchExecuteBodyParams(d *schema.ResourceData, instances []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"execute_param":   buildScriptBatchExecuteParam(d),
		"execute_batches": buildScriptBatchExecuteBatches(d.Get("execute_batches").([]interface{}), instances),
	}
}

func buildScriptBatchExecuteBatches(executeBatches, instances []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(executeBatches))
	for _, v := range executeBatches {
		instanceIds := utils.ExpandToStringList(utils.PathSearch("instance_ids", v, make([]interface{}, 0)).([]interface{}))
		result = append(result, map[string]interface{}{
			"batch_index":       utils.PathSearch("batch_index", v, nil),
			"rotation_strategy": "CONTINUE",
			"target_instances":  buildScriptBatchExecuteBatchesTargetInstances(instanceIds, instances),
		})
	}

	return result
}

func buildScriptBatchExecuteBatchesTargetInstances(instanceIds []string, instances []interface{}) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(instanceIds))
	for _, instanceId := range instanceIds {
		targetInstance := utils.PathSearch(fmt.Sprintf("[?resource_id== '%s']|[0]", instanceId), instances, nil)
		results = append(results, map[string]interface{}{
			"resource_id": utils.PathSearch("resource_id", targetInstance, nil),
			"region_id":   utils.PathSearch("region_id", targetInstance, nil),
		})
	}

	return results
}

func listResources(client *golangsdk.ServiceClient, resourceProvider, resourceType, instanceId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/resources"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	pathParams := fmt.Sprintf("%s?provider=%s&type=%s&limit=%v&resource_id_list=%s", httpUrl, resourceProvider, resourceType, limit, instanceId)
	listPath := client.Endpoint + pathParams
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%v", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		resources := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, resources...)
		if len(resources) < limit {
			break
		}

		offset += limit
	}

	return result, nil
}

func refreshResourcesAgentStatus(client *golangsdk.ServiceClient, resourceProvider, resourceType, instanceId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resources, err := listResources(client, resourceProvider, resourceType, instanceId)
		if err != nil {
			return nil, "ERROR", err
		}

		if len(resources) < 1 {
			return nil, "ERROR", fmt.Errorf("unable to find any resources")
		}

		for _, v := range resources {
			agentId := utils.PathSearch("agent_id", v, "").(string)
			agentStatus := utils.PathSearch("agent_state", v, "").(string)
			if agentId != "" && agentStatus == "ONLINE" {
				continue
			}

			if agentId == "" && agentStatus == "" {
				return nil, "ERROR", fmt.Errorf("UniAgent is not installed on the instance (%v)", utils.PathSearch("resource_id", v, ""))
			}

			if utils.StrSliceContains([]string{"FAILED", "OFFLINE", "UNINSTALLED"}, agentStatus) {
				return nil, "ERROR", fmt.Errorf("unexpected UniAgent status (%s)", agentStatus)
			}

			return resources, "PENDING", nil
		}

		return resources, "COMPLETE", nil
	}
}

func batchExecuteScript(client *golangsdk.ServiceClient, d *schema.ResourceData, instances []interface{}) (interface{}, error) {
	httpUrl := "v1/job/scripts/{script_uuid}"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{script_uuid}", d.Get("script_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildScriptBatchExecuteBodyParams(d, instances)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(createResp)
}

func resourceScriptBatchExecuteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	resourceProvider := d.Get("resource_provider").(string)
	resourceType := d.Get("type").(string)
	// Synchronize the ECS instances to get information about UniAgent.
	if d.Get("is_sync").(bool) {
		if err = syncResourceInfo(client, resourceProvider, resourceType); err != nil {
			return diag.FromErr(err)
		}
	}

	instanceIds := utils.PathSearch("[*].instance_ids[]", d.Get("execute_batches").([]interface{}), make([]interface{}, 0))
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETE"},
		Refresh: refreshResourcesAgentStatus(client, resourceProvider, resourceType,
			strings.Join(utils.ExpandToStringList(instanceIds.([]interface{})), ",")),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 15 * time.Second,
	}

	instances, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for ECS instance agent to be online: %s", err)
	}

	respBody, err := batchExecuteScript(client, d, instances.([]interface{}))
	if err != nil {
		return diag.Errorf("error executing script on instances: %s", err)
	}

	ticketId := utils.PathSearch("data", respBody, "").(string)
	if ticketId == "" {
		return diag.Errorf("unable to find the script execution ticket ID from the API response")
	}

	d.SetId(ticketId)

	stateConf = &retry.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"exited"},
		Refresh:      refreshGetExecutionTicketDetail(client, ticketId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        15 * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for script execution to be completed: %s", err)
	}

	return resourceScriptBatchExecuteRead(ctx, d, meta)
}

func resourceScriptBatchExecuteRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		ticketId = d.Id()
	)
	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	ticketDetail, err := GetExecutionTicketDetail(client, ticketId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			scriptOrderNotFoundErrCodes...), fmt.Sprintf("error retrieving COC script execute ticket (%s)", ticketId))
	}

	mErr := multierror.Append(nil,
		// Required parameters.
		d.Set("script_id", utils.PathSearch("data.properties.script_uuid", ticketDetail, nil)),
		d.Set("timeout", utils.PathSearch("data.properties.execute_param.timeout", ticketDetail, nil)),
		d.Set("execute_user", utils.PathSearch("data.properties.execute_param.execute_user", ticketDetail, nil)),
		// Attributes.
		d.Set("status", utils.PathSearch("data.status", ticketDetail, nil)),
		d.Set("script_name", utils.PathSearch("data.properties.script_name", ticketDetail, nil)),
		d.Set("created_at", flattenScriptTimeStamp(ticketDetail, "data.gmt_created")),
		d.Set("finished_at", flattenScriptTimeStamp(ticketDetail, "data.gmt_finished")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceScriptBatchExecuteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func cancelScriptExecute(client *golangsdk.ServiceClient, ticketId string, params map[string]interface{}) error {
	httpUrl := "v1/job/script/orders/{execute_uuid}/operation"
	operationPath := client.Endpoint + httpUrl
	operationPath = strings.ReplaceAll(operationPath, "{execute_uuid}", ticketId)

	operationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         params,
	}

	_, err := client.Request("PUT", operationPath, &operationOpt)
	return err
}

func resourceScriptBatchExecuteDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	if utils.StrSliceContains([]string{"FINISHED", "ROLLBACKED", "CANCELED", "ERROR", "ABNORMAL"}, d.Get("status").(string)) {
		return nil
	}

	// Cancel the ticket when it is in PROCESSING or PAUSED status.
	params := map[string]interface{}{
		"operation_type": "CANCEL_ORDER",
	}
	err = cancelScriptExecute(client, d.Id(), params)
	if err != nil {
		return diag.Errorf("error canceling COC script batch execute (%s): %s", d.Id(), err)
	}

	return nil
}
