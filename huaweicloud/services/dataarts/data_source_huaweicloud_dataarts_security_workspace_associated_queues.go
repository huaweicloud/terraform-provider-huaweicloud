package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/security/permission/queue/assigned-source
func DataSourceSecurityWorkspaceAssociatedQueues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityWorkspaceAssociatedQueuesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the workspace associated queues are located.`,
			},

			// Required
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace.`,
			},

			// Optional
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the queue.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the cluster.`,
			},

			// Attributes.
			"queues": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of workspace associated queues that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the queue.`,
						},
						"source_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The service type of the queue.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the queue.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the queue.`,
						},
						"attribute": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The attribute of the queue.`,
						},
						"connection_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the data connection.`,
						},
						"connection_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the data connection.`,
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the cluster.`,
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the cluster.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the queue was added to the workspace.`,
						},
						"create_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The operator who added the queue to the workspace.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the queue in the workspace.`,
						},
						"update_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The updater of the queue in the workspace.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the project.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the workspace associated queue.`,
						},
					},
				},
			},
		},
	}
}

func buildSecurityWorkspaceAssociatedQueuesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		res = fmt.Sprintf("%s&cluster_id=%v", res, v)
	}

	return res
}

func listSecurityWorkspaceAssociatedQueues(client *golangsdk.ServiceClient, workspaceId string,
	queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/security/permission/queue/assigned-source?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(queryParams) > 0 && queryParams[0] != "" {
		listPath += queryParams[0]
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		queueSources := utils.PathSearch("queue_sources", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, queueSources...)
		if len(queueSources) < limit {
			break
		}
		offset += len(queueSources)
	}

	return result, nil
}

func flattenSecurityWorkspaceAssociatedQueues(queueSources []interface{}) []map[string]interface{} {
	if len(queueSources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(queueSources))
	for _, qs := range queueSources {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", qs, nil),
			"source_type":     utils.PathSearch("source_type", qs, nil),
			"name":            utils.PathSearch("queue_name", qs, nil),
			"type":            utils.PathSearch("queue_type", qs, nil),
			"attribute":       utils.PathSearch("queue_attr", qs, nil),
			"connection_id":   utils.PathSearch("conn_id", qs, nil),
			"connection_name": utils.PathSearch("conn_name", qs, nil),
			"cluster_id":      utils.PathSearch("cluster_id", qs, nil),
			"cluster_name":    utils.PathSearch("cluster_name", qs, nil),
			"create_user":     utils.PathSearch("create_user", qs, nil),
			"update_user":     utils.PathSearch("update_user", qs, nil),
			"project_id":      utils.PathSearch("project_id", qs, nil),
			"description":     utils.PathSearch("description", qs, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_time", qs, float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("update_time", qs, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceSecurityWorkspaceAssociatedQueuesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	workspaceId := d.Get("workspace_id").(string)
	queueSources, err := listSecurityWorkspaceAssociatedQueues(client, workspaceId,
		buildSecurityWorkspaceAssociatedQueuesQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying workspace associated queues: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("queues", flattenSecurityWorkspaceAssociatedQueues(queueSources)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
