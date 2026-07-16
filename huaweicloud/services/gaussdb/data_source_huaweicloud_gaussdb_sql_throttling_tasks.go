package gaussdb

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

// @API GaussDB POST /v3.1/{project_id}/instances/{instance_id}/limit-task-list
func DataSourceGaussDbSqlThrottlingTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussdbOpengaussSqlThrottlingTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the GaussDB OpenGauss instance.`,
			},
			"task_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the throttling task scope.`,
			},
			"limit_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the throttling type.`,
			},
			// This parameter not take effect.
			"limit_type_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the throttling type value. Fuzzy match is supported.`,
			},
			"task_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the throttling task name. Fuzzy match is supported.`,
			},
			"sql_model": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the SQL template. Fuzzy match is supported.`,
				Deprecated:  "Deprecated",
			},
			"rule_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the rule name.`,
			},
			"sql_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter not take effect.
			"node_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start time of the throttling task in the format of **yyy-mm-ddThh:mm:ss+0000**.`,
				Deprecated:  "Deprecated",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end time of the throttling task in the format of **yyy-mm-ddThh:mm:ss+0000**.`,
				Deprecated:  "Deprecated",
			},
			"limit_task_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of throttling tasks.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance ID.`,
						},
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the throttling task ID.`,
						},
						"task_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the throttling task name.`,
						},
						"task_scope": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the throttling task scope.`,
						},
						"limit_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the throttling task type.`,
						},
						"limit_type_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the throttling task type value.`,
						},
						"sql_model": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the SQL template.`,
						},
						"key_words": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the keyword.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the throttling task status.`,
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the rule name.`,
						},
						"parallel_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the maximum concurrency.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the start time of the throttling task in the format of **yyyy-mm-ddThh:mm:ssZ**.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the end time of the throttling task in the format of **yyyy-mm-ddThh:mm:ssZ**.`,
						},
						"cpu_utilization": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the CPU usage.`,
						},
						"memory_utilization": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the memory usage.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creation time in the format of **yyyy-mm-ddThh:mm:ssZ**.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the update time in the format of **yyyy-mm-ddThh:mm:ssZ**.`,
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator.`,
						},
						"modifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the modifier.`,
						},
						"databases": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the databases of the instance.`,
						},
						"node_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the CN information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the node ID.`,
									},
									"sql_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the ID of the SQL statement executed on the node.`,
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

func buildQuerySqlThrottlingTasksBodyParams(d *schema.ResourceData, limit int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"limit":            limit,
		"task_scope":       utils.ValueIgnoreEmpty(d.Get("task_scope")),
		"limit_type":       utils.ValueIgnoreEmpty(d.Get("limit_type")),
		"limit_type_value": utils.ValueIgnoreEmpty(d.Get("limit_type_value")),
		"task_name":        utils.ValueIgnoreEmpty(d.Get("task_name")),
		"rule_name":        utils.ValueIgnoreEmpty(d.Get("rule_name")),
		"sql_id":           utils.ValueIgnoreEmpty(d.Get("sql_id")),
		"node_ids":         utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("node_ids").([]interface{}))),
	}

	return bodyParams
}

func dataSourceGaussdbOpengaussSqlThrottlingTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
		httpUrl = "v3.1/{project_id}/instances/{instance_id}/limit-task-list"
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	bodyParams := utils.RemoveNil(buildQuerySqlThrottlingTasksBodyParams(d, limit))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	}

	for {
		bodyParams["offset"] = offset
		listOpt.JSONBody = bodyParams

		resp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB SQL throttling tasks: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		tasks := utils.PathSearch("limit_task_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tasks...)

		if len(tasks) < limit {
			break
		}

		offset += len(tasks)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("limit_task_list", flattenSqlThrottlingTasks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSqlThrottlingTasks(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"instance_id":        utils.PathSearch("instance_id", v, nil),
			"task_id":            utils.PathSearch("task_id", v, nil),
			"task_name":          utils.PathSearch("task_name", v, nil),
			"task_scope":         utils.PathSearch("task_scope", v, nil),
			"limit_type":         utils.PathSearch("limit_type", v, nil),
			"limit_type_value":   utils.PathSearch("limit_type_value", v, nil),
			"sql_model":          utils.PathSearch("sql_model", v, nil),
			"key_words":          utils.PathSearch("key_words", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"rule_name":          utils.PathSearch("rule_name", v, nil),
			"parallel_size":      utils.PathSearch("parallel_size", v, nil),
			"start_time":         utils.PathSearch("start_time", v, nil),
			"end_time":           utils.PathSearch("end_time", v, nil),
			"cpu_utilization":    utils.PathSearch("cpu_utilization", v, nil),
			"memory_utilization": utils.PathSearch("memory_utilization", v, nil),
			"created_at":         utils.PathSearch("created", v, nil),
			"updated_at":         utils.PathSearch("updated", v, nil),
			"creator":            utils.PathSearch("creator", v, nil),
			"modifier":           utils.PathSearch("modifier", v, nil),
			"databases":          utils.PathSearch("databases", v, nil),
			"node_infos": flattenSqlThrottlingTasksNodeInfo(
				utils.PathSearch("node_infos", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenSqlThrottlingTasksNodeInfo(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"node_id": utils.PathSearch("node_id", v, nil),
			"sql_id":  utils.PathSearch("sql_id", v, nil),
		})
	}

	return result
}
