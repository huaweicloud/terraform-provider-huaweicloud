package ces

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var agentMaintenanceTaskNonUpdatableParams = []string{
	"invocation_type", "instance_id", "invocation_target",
	"invocation_id", "version_type", "version",
}

// @API CES POST /v3/{project_id}/agent-invocations/batch-create
// @API CES GET /v3/{project_id}/agent-invocations
func ResourceAgentMaintenanceTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentMaintenanceTaskCreate,
		ReadContext:   resourceAgentMaintenanceTaskRead,
		UpdateContext: resourceAgentMaintenanceTaskUpdate,
		DeleteContext: resourceAgentMaintenanceTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(agentMaintenanceTaskNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"invocation_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the task type.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the server ID.`,
			},
			"invocation_target": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "telescope",
				Description: `Specifies the task object. Only **telescope** is supported.`,
			},
			"invocation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task ID.`,
			},
			"version_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the version the agent will be upgraded to.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the version number.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"invocations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The result of the agent maintenance task.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"invocation_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The task ID.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The server ID.`,
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The server name`,
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The server type.`,
						},
						"intranet_ips": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `The private IP address list.`,
						},
						"elastic_ips": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `The EIP list.`,
						},
						"invocation_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The task type.`,
						},
						"invocation_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The task status.`,
						},
						"invocation_target": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The task object.`,
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `When the task was created.`,
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `When the task was updated.`,
						},
						"current_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current version of the agent.`,
						},
						"target_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The target version.`,
						},
					},
				},
			},
		},
	}
}

func resourceAgentMaintenanceTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		createAgentMaintenanceTaskHttpUrl = "v3/{project_id}/agent-invocations/batch-create"
		createAgentMaintenanceTaskProduct = "ces"
	)
	client, err := conf.NewServiceClient(createAgentMaintenanceTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	createAgentMaintenanceTaskPath := client.Endpoint + createAgentMaintenanceTaskHttpUrl
	createAgentMaintenanceTaskPath = strings.ReplaceAll(createAgentMaintenanceTaskPath, "{project_id}", client.ProjectID)

	createAgentMaintenanceTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createAgentMaintenanceTaskOpt.JSONBody = utils.RemoveNil(buildCreateAgentMaintenanceTaskBodyParams(d))
	resp, err := client.Request("POST", createAgentMaintenanceTaskPath, &createAgentMaintenanceTaskOpt)
	if err != nil {
		return diag.Errorf("error creating CES agent maintenance task: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	status := utils.PathSearch("invocations|[0].ret_status", respBody, "").(string)
	if status != "successful" {
		return diag.Errorf("error creating CES agent maintenance task: %s", respBody)
	}

	id := utils.PathSearch("invocations|[0].invocation_id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CES agent maintenance task: ID is not found in API response")
	}
	d.SetId(id)

	err = waitingForAgentMaintenanceTaskCompleted(ctx, client, d)
	if err != nil {
		return diag.Errorf("error waiting for agent maintenance task creation completed: %s", err)
	}

	return resourceAgentMaintenanceTaskRead(ctx, d, meta)
}

func buildCreateAgentMaintenanceTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"invocation_type":   d.Get("invocation_type"),
		"invocation_target": utils.ValueIgnoreEmpty(d.Get("invocation_target")),
		"version_type":      utils.ValueIgnoreEmpty(d.Get("version_type")),
		"version":           utils.ValueIgnoreEmpty(d.Get("version")),
	}

	if v, ok := d.GetOk("instance_id"); ok {
		bodyParam["instance_ids"] = []string{v.(string)}
	}

	if v, ok := d.GetOk("invocation_id"); ok {
		bodyParam["invocation_ids"] = []string{v.(string)}
	}
	return bodyParam
}

func waitingForAgentMaintenanceTaskCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			result, err := getAgentMaintenanceTaskExecutionResult(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("invocation_status", result, "").(string)
			if status == "SUCCEEDED" {
				return result, "COMPLETED", nil
			}

			if status == "PENDING" || status == "RUNNING" {
				return result, "PENDING", nil
			}
			return result, status, fmt.Errorf("agent maintenance task abnormal status: %s", status)
		},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func resourceAgentMaintenanceTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	var mErr *multierror.Error

	getAgentMaintenanceTaskProduct := "ces"
	client, err := conf.NewServiceClient(getAgentMaintenanceTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	invocation, err := getAgentMaintenanceTaskExecutionResult(client, d.Id())
	if invocation == nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES agent maintenance task")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("invocations", flattenAgentMaintenanceTask(invocation)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func getAgentMaintenanceTaskExecutionResult(client *golangsdk.ServiceClient, taskID string) (interface{}, error) {
	httpUrl := "v3/{project_id}/agent-invocations"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += fmt.Sprintf("?invocation_id=%s", taskID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", path, &opts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	task := utils.PathSearch("invocations|[0]", respBody, nil)
	if task == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return task, nil
}

func flattenAgentMaintenanceTask(param interface{}) []interface{} {
	if param == nil {
		return nil
	}
	rst := map[string]interface{}{
		"invocation_id":     utils.PathSearch("invocation_id", param, nil),
		"instance_id":       utils.PathSearch("instance_id", param, nil),
		"instance_name":     utils.PathSearch("instance_name", param, nil),
		"instance_type":     utils.PathSearch("instance_type", param, nil),
		"intranet_ips":      utils.PathSearch("intranet_ips", param, nil),
		"elastic_ips":       utils.PathSearch("elastic_ips", param, nil),
		"invocation_type":   utils.PathSearch("invocation_type", param, nil),
		"invocation_status": utils.PathSearch("invocation_status", param, nil),
		"invocation_target": utils.PathSearch("invocation_target", param, nil),
		"create_time":       utils.PathSearch("create_time", param, nil),
		"update_time":       utils.PathSearch("update_time", param, nil),
		"current_version":   utils.PathSearch("current_version", param, nil),
		"target_version":    utils.PathSearch("target_version", param, nil),
	}
	return []interface{}{rst}
}

func resourceAgentMaintenanceTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAgentMaintenanceTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource. Deleting this resource will not change
		the status of the current resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
