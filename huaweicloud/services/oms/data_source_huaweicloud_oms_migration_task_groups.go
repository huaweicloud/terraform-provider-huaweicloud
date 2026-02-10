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

// @API OMS GET /v2/{project_id}/taskgroups
func DataSourceMigrationTaskGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMigrationTaskGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"taskgroups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
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
						"enable_metadata_migration": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_failed_object_recording": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_restore": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_kms": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"task_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
						"smn_config": {
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
						"migrate_since": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"migrate_speed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_task_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_task_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"failed_task_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"complete_task_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"paused_task_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"executing_task_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"waiting_task_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_complete_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"success_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"fail_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"skip_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_complete_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"complete_size": {
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
					},
				},
			},
		},
	}
}

func buildMigrationTaskGroupsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}

	return queryParams
}

func listMigrationTaskGroups(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/taskgroups?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{limit}", strconv.Itoa(limit))
	getPath += buildMigrationTaskGroupsQueryParams(d)
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

		taskGroups := utils.PathSearch("taskgroups", getRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, taskGroups...)
		if len(taskGroups) < limit {
			break
		}

		offset += len(taskGroups)
	}

	return result, nil
}

func dataSourceMigrationTaskGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("oms", region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	resp, err := listMigrationTaskGroups(client, d)
	if err != nil {
		return diag.Errorf("error retrieving migration task groups: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("taskgroups", flattenMigrationTaskGroups(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMigrationTaskGroups(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"group_id":                       utils.PathSearch("group_id", v, nil),
			"status":                         utils.PathSearch("status", v, nil),
			"error_reason":                   flattenErrorReasonResp(utils.PathSearch("error_reason", v, nil)),
			"src_node":                       flattenSrcNodeResp(utils.PathSearch("src_node", v, nil)),
			"description":                    utils.PathSearch("description", v, nil),
			"dst_node":                       flattenDstNodeResp(utils.PathSearch("dst_node", v, nil)),
			"enable_metadata_migration":      utils.PathSearch("enable_metadata_migration", v, nil),
			"enable_failed_object_recording": utils.PathSearch("enable_failed_object_recording", v, nil),
			"enable_restore":                 utils.PathSearch("enable_restore", v, nil),
			"enable_kms":                     utils.PathSearch("enable_kms", v, nil),
			"task_type":                      utils.PathSearch("task_type", v, nil),
			"bandwidth_policy": flattenBandwidthPolicyResp(
				utils.PathSearch("bandwidth_policy", v, make([]interface{}, 0)).([]interface{})),
			"smn_config":            flattenSmnConfigInfo(utils.PathSearch("smn_config", v, nil)),
			"source_cdn":            flattenSourceCdnResp(utils.PathSearch("source_cdn", v, nil)),
			"migrate_since":         utils.PathSearch("migrate_since", v, nil),
			"migrate_speed":         utils.PathSearch("migrate_speed", v, nil),
			"total_time":            utils.PathSearch("total_time", v, nil),
			"start_time":            utils.PathSearch("start_time", v, nil),
			"total_task_num":        utils.PathSearch("total_task_num", v, nil),
			"create_task_num":       utils.PathSearch("create_task_num", v, nil),
			"failed_task_num":       utils.PathSearch("failed_task_num", v, nil),
			"complete_task_num":     utils.PathSearch("complete_task_num", v, nil),
			"paused_task_num":       utils.PathSearch("paused_task_num", v, nil),
			"executing_task_num":    utils.PathSearch("executing_task_num", v, nil),
			"waiting_task_num":      utils.PathSearch("waiting_task_num", v, nil),
			"total_num":             utils.PathSearch("total_num", v, nil),
			"create_complete_num":   utils.PathSearch("create_complete_num", v, nil),
			"success_num":           utils.PathSearch("success_num", v, nil),
			"fail_num":              utils.PathSearch("fail_num", v, nil),
			"skip_num":              utils.PathSearch("skip_num", v, nil),
			"total_size":            utils.PathSearch("total_size", v, nil),
			"create_complete_size":  utils.PathSearch("create_complete_size", v, nil),
			"complete_size":         utils.PathSearch("complete_size", v, nil),
			"failed_object_record":  flattenFailObjectRecordResp(utils.PathSearch("failed_object_record", v, nil)),
			"object_overwrite_mode": utils.PathSearch("object_overwrite_mode", v, nil),
			"dst_storage_policy":    utils.PathSearch("dst_storage_policy", v, nil),
			"consistency_check":     utils.PathSearch("consistency_check", v, nil),
			"enable_requester_pays": utils.PathSearch("enable_requester_pays", v, nil),
		})
	}

	return result
}

func flattenErrorReasonResp(errorReasonResp interface{}) []map[string]interface{} {
	if errorReasonResp == nil {
		return nil
	}

	result := map[string]interface{}{
		"error_code": utils.PathSearch("error_code", errorReasonResp, nil),
		"error_msg":  utils.PathSearch("error_msg", errorReasonResp, nil),
	}

	return []map[string]interface{}{result}
}

func flattenSrcNodeResp(srcNodeResp interface{}) []map[string]interface{} {
	if srcNodeResp == nil {
		return nil
	}

	result := map[string]interface{}{
		"bucket":     utils.PathSearch("bucket", srcNodeResp, nil),
		"cloud_type": utils.PathSearch("cloud_type", srcNodeResp, nil),
		"region":     utils.PathSearch("region", srcNodeResp, nil),
		"app_id":     utils.PathSearch("app_id", srcNodeResp, nil),
		"object_key": utils.PathSearch("object_key", srcNodeResp, nil),
		"list_file":  flattenListFile(utils.PathSearch("list_file", srcNodeResp, nil)),
	}

	return []map[string]interface{}{result}
}

func flattenListFile(listFile interface{}) []map[string]interface{} {
	if listFile == nil {
		return nil
	}

	result := map[string]interface{}{
		"list_file_key": utils.PathSearch("list_file_key", listFile, nil),
		"obs_bucket":    utils.PathSearch("obs_bucket", listFile, nil),
		"list_file_num": utils.PathSearch("list_file_num", listFile, nil),
	}

	return []map[string]interface{}{result}
}

func flattenDstNodeResp(dstNodeResp interface{}) []map[string]interface{} {
	if dstNodeResp == nil {
		return nil
	}

	result := map[string]interface{}{
		"bucket":      utils.PathSearch("bucket", dstNodeResp, nil),
		"region":      utils.PathSearch("region", dstNodeResp, nil),
		"save_prefix": utils.PathSearch("save_prefix", dstNodeResp, nil),
	}

	return []map[string]interface{}{result}
}

func flattenBandwidthPolicyResp(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"max_bandwidth": utils.PathSearch("max_bandwidth", v, nil),
			"start":         utils.PathSearch("start", v, nil),
			"end":           utils.PathSearch("end", v, nil),
		})
	}

	return result
}

func flattenSmnConfigInfo(smnInfo interface{}) []map[string]interface{} {
	if smnInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"notify_result":        utils.PathSearch("notify_result", smnInfo, nil),
		"notify_error_message": utils.PathSearch("notify_error_message", smnInfo, nil),
		"topic_name":           utils.PathSearch("topic_name", smnInfo, nil),
	}

	return []map[string]interface{}{result}
}

func flattenFailObjectRecordResp(sourceCdnResp interface{}) []map[string]interface{} {
	if sourceCdnResp == nil {
		return nil
	}

	result := map[string]interface{}{
		"result":        utils.PathSearch("result", sourceCdnResp, nil),
		"list_file_key": utils.PathSearch("list_file_key", sourceCdnResp, nil),
		"error_code":    utils.PathSearch("error_code", sourceCdnResp, nil),
	}

	return []map[string]interface{}{result}
}
