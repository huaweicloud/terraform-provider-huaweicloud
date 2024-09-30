package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/instances/auditlogs
func DataSourceSecmasterPlaybookAuditLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterPlaybookAuditLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workspace ID.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start time.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end time.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status. The value can be **RUNNING**, **FINISHED**, **FAILED**, **RETRYING**, or **TERMINATED**.`,
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the instance type. The value can be **AOP_WORKFLOW**, **SCRIPT**, or **PLAYBOOK**.`,
			},
			"action_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the workflow name.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the instance ID.`,
			},
			"log_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the log level. The value can be **DEBUG**, **INFO**, **WARN** or **ERROR**.`,
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the triggering type.`,
			},
			"action_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the workflow ID.`,
			},
			"parent_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the instance ID of the parent node.`,
			},
			"input": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the input information.`,
			},
			"output": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the output information.`,
			},
			"error_msg": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the error message.`,
			},
			"audit_logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The audit log list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status.`,
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance type.`,
						},
						"action_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The workflow name.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID.`,
						},
						"log_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The log level.`,
						},
						"trigger_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The triggering type.`,
						},
						"action_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The workflow ID.`,
						},
						"parent_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID of the parent node.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The end time.`,
						},
						"input": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The input information.`,
						},
						"output": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The output information.`,
						},
						"error_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The error message.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceSecmasterPlaybookAuditLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// getPlaybookAuditLogs: Query the SecMaster playbook audit logs.
	var (
		listPlaybookAuditLogsHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/instances/auditlogs?limit=50"
		listPlaybookAuditLogsProduct = "secmaster"
	)
	client, err := cfg.NewServiceClient(listPlaybookAuditLogsProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	listPlaybookAuditLogsPath := client.Endpoint + listPlaybookAuditLogsHttpUrl
	listPlaybookAuditLogsPath = strings.ReplaceAll(listPlaybookAuditLogsPath, "{project_id}", client.ProjectID)
	listPlaybookAuditLogsPath = strings.ReplaceAll(listPlaybookAuditLogsPath, "{workspace_id}", d.Get("workspace_id").(string))

	listPlaybookAuditLogsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listPlaybookAuditLogsOpt.JSONBody = utils.RemoveNil(buildPlaybookAuditLogsBodyParams(d))

	playbookAuditLogs := make([]interface{}, 0)
	offset := 0
	for {
		currentPath := listPlaybookAuditLogsPath + fmt.Sprintf("&offset=%d", offset)
		listPlaybookAuditLogsResp, err := client.Request("POST", currentPath, &listPlaybookAuditLogsOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listPlaybookAuditLogsRespBody, err := utils.FlattenResponse(listPlaybookAuditLogsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		auditLogs := utils.PathSearch("audit_logs", listPlaybookAuditLogsRespBody, make([]interface{}, 0)).([]interface{})
		playbookAuditLogs = append(playbookAuditLogs, auditLogs...)

		if len(auditLogs) == 0 {
			break
		}
		offset += len(auditLogs)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("audit_logs", flattenAuditLogsResponseBody(playbookAuditLogs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildPlaybookAuditLogsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action_id":          utils.ValueIgnoreEmpty(d.Get("action_id")),
		"action_name":        utils.ValueIgnoreEmpty(d.Get("action_name")),
		"end_time":           utils.ValueIgnoreEmpty(d.Get("end_time")),
		"error_msg":          utils.ValueIgnoreEmpty(d.Get("error_msg")),
		"input":              utils.ValueIgnoreEmpty(d.Get("input")),
		"instance_id":        utils.ValueIgnoreEmpty(d.Get("instance_id")),
		"instance_type":      utils.ValueIgnoreEmpty(d.Get("instance_type")),
		"log_level":          utils.ValueIgnoreEmpty(d.Get("log_level")),
		"output":             utils.ValueIgnoreEmpty(d.Get("output")),
		"parent_instance_id": utils.ValueIgnoreEmpty(d.Get("parent_instance_id")),
		"start_time":         utils.ValueIgnoreEmpty(d.Get("start_time")),
		"status":             utils.ValueIgnoreEmpty(d.Get("status")),
		"trigger_type":       utils.ValueIgnoreEmpty(d.Get("trigger_type")),
	}

	return bodyParams
}

func flattenAuditLogsResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	auditLogs := make([]interface{}, len(resp))
	for i, v := range resp {
		auditLogs[i] = map[string]interface{}{
			"input":              utils.PathSearch("input", v, nil),
			"start_time":         utils.PathSearch("start_time", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"instance_type":      utils.PathSearch("instance_type", v, nil),
			"action_name":        utils.PathSearch("action_name", v, nil),
			"instance_id":        utils.PathSearch("instance_id", v, nil),
			"log_level":          utils.PathSearch("log_level", v, nil),
			"end_time":           utils.PathSearch("last_report_time", v, nil),
			"trigger_type":       utils.PathSearch("trigger_type", v, nil),
			"action_id":          utils.PathSearch("action_id", v, nil),
			"parent_instance_id": utils.PathSearch("parent_instance_id", v, nil),
			"output":             utils.PathSearch("output", v, nil),
			"error_msg":          utils.PathSearch("error_msg", v, nil),
		}
	}
	return auditLogs
}
