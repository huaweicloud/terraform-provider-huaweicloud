package modelarts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/workflows/{workflow_id}/schedules
func DataSourceV2WorkflowSchedules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2WorkflowSchedulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the workflow schedules are located.`,
			},

			// Required parameters.
			"workflow_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workflow to which the schedule configurations belong.`,
			},

			// Attributes.
			"schedules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2WorkflowSchedulesElemSchema(),
				Description: `The list of the workflow schedules.`,
			},
		},
	}
}

func dataSourceV2WorkflowSchedulePoliciesElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"on_failure": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The policy action when the workflow execution fails.`,
			},
			"on_running": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The policy action when the workflow is already running.`,
			},
		},
	}
}

func dataSourceV2WorkflowSchedulesElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the workflow schedule.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the workflow schedule.`,
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The content of the workflow schedule, in JSON format.`,
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The action of the workflow schedule.`,
			},
			"workflow_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the workflow to which the schedule belongs.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user ID that created the workflow schedule.`,
			},
			"enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the workflow schedule is enabled.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2WorkflowSchedulePoliciesElemSchema(),
				Description: `The scheduling policies of the workflow schedule.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the workflow schedule, in RFC3339 format.`,
			},
		},
	}
}

func listV2WorkflowSchedules(client *golangsdk.ServiceClient, workflowId string) ([]interface{}, error) {
	var httpUrl = "v2/{project_id}/workflows/{workflow_id}/schedules"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workflow_id}", workflowId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("schedules", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenV2WorkflowSchedules(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("uuid", item, nil),
			"type":        utils.PathSearch("type", item, nil),
			"content":     utils.JsonToString(utils.PathSearch("content", item, nil)),
			"action":      utils.PathSearch("action", item, nil),
			"workflow_id": utils.PathSearch("workflow_id", item, nil),
			"user_id":     utils.PathSearch("user_id", item, nil),
			"enable":      utils.PathSearch("enable", item, nil),
			"policies": flattenV2WorkflowSchedulePolicies(utils.PathSearch("policies", item,
				make(map[string]interface{})).(map[string]interface{})),
			"created_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at", item, "").(string))/1000, false),
		})
	}

	return result
}

func dataSourceV2WorkflowSchedulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		workflowId = d.Get("workflow_id").(string)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	schedules, err := listV2WorkflowSchedules(client, workflowId)
	if err != nil {
		return diag.Errorf("error querying ModelArts workflow schedules: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("schedules", flattenV2WorkflowSchedules(schedules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
