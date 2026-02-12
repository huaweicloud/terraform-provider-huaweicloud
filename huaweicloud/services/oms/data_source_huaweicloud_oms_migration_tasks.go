package oms

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

// @API OMS GET /v2/{project_id}/tasks
func DataSourceMigrationTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMigrationTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth_policy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_bandwidth": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"complete_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dst_node": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"save_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"enable_failed_object_recording": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_kms": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_metadata_migration": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_restore": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"error_reason": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"error_msg": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"fail_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"failed_object_record": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"result": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"list_file_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"error_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_query_over": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"left_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"migrate_since": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"migrate_speed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"real_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"skipped_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"src_node": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"app_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"object_key": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"list_file": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"list_file_key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"obs_bucket": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"list_file_num": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						}, "status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"successful_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"task_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"smn_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notify_result": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"notify_error_message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"topic_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"source_cdn": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"authentication_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"success_record_error_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"skip_record_error_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"object_overwrite_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dst_storage_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consistency_check": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_requester_pays": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"task_priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildMigrationTasksQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("group_id"); ok {
		queryParams = fmt.Sprintf("%s&group_id=%v", queryParams, v)
	}

	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}

	return queryParams
}

func listMigrationTasks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/tasks?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{limit}", strconv.Itoa(limit))
	getPath += buildMigrationTasksQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return nil, err
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}

		tasks := utils.PathSearch("tasks", getRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tasks...)
		if len(tasks) < limit {
			break
		}

		offset += len(tasks)
	}

	return result, nil
}

func dataSourceMigrationTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("oms", region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	resp, err := listMigrationTasks(client, d)
	if err != nil {
		return diag.Errorf("error retrieving migration tasks: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", flattenMigrationTasks(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMigrationTasks(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"bandwidth_policy": flattenBandwidthPolicyResp(
				utils.PathSearch("bandwidth_policy", v, make([]interface{}, 0)).([]interface{})),
			"complete_size":                  utils.PathSearch("complete_size", v, nil),
			"description":                    utils.PathSearch("description", v, nil),
			"dst_node":                       flattenDstNodeResp(utils.PathSearch("dst_node", v, nil)),
			"enable_failed_object_recording": utils.PathSearch("enable_failed_object_recording", v, nil),
			"enable_kms":                     utils.PathSearch("enable_kms", v, nil),
			"enable_metadata_migration":      utils.PathSearch("enable_metadata_migration", v, nil),
			"enable_restore":                 utils.PathSearch("enable_restore", v, nil),
			"error_reason":                   flattenErrorReasonResp(utils.PathSearch("error_reason", v, nil)),
			"fail_num":                       utils.PathSearch("fail_num", v, nil),
			"failed_object_record":           flattenFailObjectRecordResp(utils.PathSearch("failed_object_record", v, nil)),
			"group_id":                       utils.PathSearch("group_id", v, nil),
			"id":                             utils.PathSearch("id", v, nil),
			"is_query_over":                  utils.PathSearch("is_query_over", v, nil),
			"left_time":                      utils.PathSearch("left_time", v, nil),
			"migrate_since":                  utils.PathSearch("migrate_since", v, nil),
			"migrate_speed":                  utils.PathSearch("migrate_speed", v, nil),
			"name":                           utils.PathSearch("name", v, nil),
			"progress":                       utils.PathSearch("progress", v, nil),
			"real_size":                      utils.PathSearch("real_size", v, nil),
			"skipped_num":                    utils.PathSearch("skipped_num", v, nil),
			"src_node":                       flattenSrcNodeResp(utils.PathSearch("src_node", v, nil)),
			"start_time":                     utils.PathSearch("start_time", v, nil),
			"status":                         utils.PathSearch("status", v, nil),
			"successful_num":                 utils.PathSearch("successful_num", v, nil),
			"task_type":                      utils.PathSearch("task_type", v, nil),
			"group_type":                     utils.PathSearch("group_type", v, nil),
			"total_num":                      utils.PathSearch("total_num", v, nil),
			"total_size":                     utils.PathSearch("total_size", v, nil),
			"total_time":                     utils.PathSearch("total_time", v, nil),
			"smn_info":                       flattenSmnConfigInfo(utils.PathSearch("smn_info", v, nil)),
			"source_cdn":                     flattenSourceCdnResp(utils.PathSearch("source_cdn", v, nil)),
			"success_record_error_reason":    utils.PathSearch("success_record_error_reason", v, nil),
			"skip_record_error_reason":       utils.PathSearch("skip_record_error_reason", v, nil),
			"object_overwrite_mode":          utils.PathSearch("object_overwrite_mode", v, nil),
			"dst_storage_policy":             utils.PathSearch("dst_storage_policy", v, nil),
			"consistency_check":              utils.PathSearch("consistency_check", v, nil),
			"enable_requester_pays":          utils.PathSearch("enable_requester_pays", v, nil),
			"task_priority":                  utils.PathSearch("task_priority", v, nil),
		})
	}

	return result
}
