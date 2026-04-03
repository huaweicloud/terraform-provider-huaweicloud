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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/tasks
func DataSourceSecmasterTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"business_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"from_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aopengine_task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifier_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifier_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"approveuser_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"approveuser_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"approver": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"definition_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"note": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"due_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"review_comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"view_parameters": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handle_parameters": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"business_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"related_object": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attachment_id_list": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"comments": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataTasksCommentsSchema(),
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"due_handle": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataTasksCommentsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildSecmasterTasksQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=100"

	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}
	if v, ok := d.GetOk("note"); ok {
		queryParams = fmt.Sprintf("%s&note=%v", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("business_type"); ok {
		queryParams = fmt.Sprintf("%s&business_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("creator_name"); ok {
		queryParams = fmt.Sprintf("%s&creator_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("query_type"); ok {
		queryParams = fmt.Sprintf("%s&query_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("from_date"); ok {
		queryParams = fmt.Sprintf("%s&from_date=%v", queryParams, v)
	}
	if v, ok := d.GetOk("to_date"); ok {
		queryParams = fmt.Sprintf("%s&to_date=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceSecmasterTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/tasks"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath += buildSecmasterTasksQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster tasks: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("tasks", flattenSecmasterTasks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecmasterTasks(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"id":                 utils.PathSearch("id", v, nil),
			"aopengine_task_id":  utils.PathSearch("aopengine_task_id", v, nil),
			"name":               utils.PathSearch("name", v, nil),
			"project_id":         utils.PathSearch("project_id", v, nil),
			"description":        utils.PathSearch("description", v, nil),
			"create_time":        utils.PathSearch("create_time", v, nil),
			"creator_id":         utils.PathSearch("creator_id", v, nil),
			"creator_name":       utils.PathSearch("creator_name", v, nil),
			"update_time":        utils.PathSearch("update_time", v, nil),
			"modifier_id":        utils.PathSearch("modifier_id", v, nil),
			"modifier_name":      utils.PathSearch("modifier_name", v, nil),
			"approveuser_id":     utils.PathSearch("approveuser_id", v, nil),
			"approveuser_name":   utils.PathSearch("approveuser_name", v, nil),
			"approver":           utils.PathSearch("approver", v, nil),
			"notes":              utils.PathSearch("notes", v, nil),
			"definition_key":     utils.PathSearch("definition_key", v, nil),
			"note":               utils.PathSearch("note", v, nil),
			"due_date":           utils.PathSearch("due_date", v, nil),
			"action_id":          utils.PathSearch("action_id", v, nil),
			"action_version_id":  utils.PathSearch("action_version_id", v, nil),
			"action_instance_id": utils.PathSearch("action_instance_id", v, nil),
			"workspace_id":       utils.PathSearch("workspace_id", v, nil),
			"review_comments":    utils.PathSearch("review_comments", v, nil),
			"view_parameters":    utils.PathSearch("view_parameters", v, nil),
			"handle_parameters":  utils.PathSearch("handle_parameters", v, nil),
			"business_type":      utils.PathSearch("business_type", v, nil),
			"related_object":     utils.PathSearch("related_object", v, nil),
			"attachment_id_list": utils.PathSearch("attachment_id_list", v, nil),
			"comments":           flattenSecmasterTasksComments(v),
			"status":             utils.PathSearch("status", v, nil),
			"due_handle":         utils.PathSearch("due_handle", v, nil),
		})
	}

	return rst
}

func flattenSecmasterTasksComments(respBody interface{}) []interface{} {
	commentsRespArray := utils.PathSearch("comments", respBody, make([]interface{}, 0)).([]interface{})
	if len(commentsRespArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(commentsRespArray))
	for _, v := range commentsRespArray {
		rst = append(rst, map[string]interface{}{
			"id":        utils.PathSearch("id", v, nil),
			"message":   utils.PathSearch("message", v, nil),
			"user_id":   utils.PathSearch("user_id", v, nil),
			"user_name": utils.PathSearch("user_name", v, nil),
			"time":      utils.PathSearch("time", v, nil),
		})
	}

	return rst
}
