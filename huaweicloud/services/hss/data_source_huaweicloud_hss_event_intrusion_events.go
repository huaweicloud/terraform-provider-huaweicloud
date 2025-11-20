package hss

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

// @API HSS GET /v5/{project_id}/event/events
func DataSourceEventIntrusionEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventIntrusionEventsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"handle_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_class_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"severity_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attack_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"att_ck": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_block": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_class_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"event_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"asset_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attack_phase": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attack_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"occur_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"handle_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"handle_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handle_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operate_accept_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"operate_detail_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"agent_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"process_pid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"is_parent": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"file_hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_attr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"login_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"login_user_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"keyword": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"forensic_info": {
							Type:     schema.TypeString, // JSON format
							Computed: true,
						},
						"resource_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"project_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enterprise_project_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vm_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vm_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"container_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"container_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pod_uid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pod_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_attr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"micro_service": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sys_arch": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_bit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"geo_info": {
							Type:     schema.TypeString, // JSON format
							Computed: true,
						},
						"malware_info": {
							Type:     schema.TypeString, // JSON format
							Computed: true,
						},
						"network_info": {
							Type:     schema.TypeString, // JSON format
							Computed: true,
						},
						"app_info": {
							Type:     schema.TypeString, // JSON format
							Computed: true,
						},
						"system_info": {
							Type:     schema.TypeString, // JSON format
							Computed: true,
						},
						"extend_info": {
							Type:     schema.TypeString, // JSON format
							Computed: true,
						},
						"recommendation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_abstract": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_info_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"process_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"process_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"process_pid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"process_uid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"process_username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"process_cmdline": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"process_filename": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"process_start_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"process_gid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"process_egid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"process_euid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ancestor_process_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ancestor_process_pid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ancestor_process_cmdline": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parent_process_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parent_process_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parent_process_pid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"parent_process_uid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"parent_process_cmdline": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parent_process_filename": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parent_process_start_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"parent_process_gid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"parent_process_egid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"parent_process_euid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"child_process_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"child_process_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"child_process_pid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"child_process_uid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"child_process_cmdline": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"child_process_filename": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"child_process_start_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"child_process_gid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"child_process_egid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"child_process_euid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"virt_cmd": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"virt_process_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"escape_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"escape_cmd": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"process_hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"process_file_hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parent_process_file_hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"block": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"user_info_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"user_gid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"user_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_home_dir": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"login_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"login_mode": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"login_last_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"login_fail_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"pwd_hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pwd_with_fuzzing": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pwd_used_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"pwd_min_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"pwd_max_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"pwd_warn_left_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"file_info_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_alias": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"file_mtime": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"file_atime": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"file_ctime": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"file_hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_md5": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_sha256": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_attr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_operation": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"file_action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_change_attr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_new_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_desc": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_key_word": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_dir": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"fd_info": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"fd_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"event_details": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"event_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"operate_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildEventIntrusionEventsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?category=%v&limit=1000", d.Get("category"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("last_days"); ok {
		queryParams = fmt.Sprintf("%s&last_days=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		queryParams = fmt.Sprintf("%s&private_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("public_ip"); ok {
		queryParams = fmt.Sprintf("%s&public_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("container_name"); ok {
		queryParams = fmt.Sprintf("%s&container_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("handle_status"); ok {
		queryParams = fmt.Sprintf("%s&handle_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("severity"); ok {
		queryParams = fmt.Sprintf("%s&severity=%v", queryParams, v)
	}
	if v, ok := d.GetOk("begin_time"); ok {
		queryParams = fmt.Sprintf("%s&begin_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("attack_tag"); ok {
		queryParams = fmt.Sprintf("%s&attack_tag=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_value"); ok {
		queryParams = fmt.Sprintf("%s&asset_value=%v", queryParams, v)
	}
	if v, ok := d.GetOk("att_ck"); ok {
		queryParams = fmt.Sprintf("%s&att_ck=%v", queryParams, v)
	}
	if v, ok := d.GetOk("event_name"); ok {
		queryParams = fmt.Sprintf("%s&event_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("auto_block"); ok {
		queryParams = fmt.Sprintf("%s&auto_block=%v", queryParams, v)
	}

	return queryParams
}

func buildEventIntrusionEventsListQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	eventClassIds := d.Get("event_class_ids").([]interface{})
	eventTypes := d.Get("event_types").([]interface{})
	severityList := d.Get("severity_list").([]interface{})
	tagList := d.Get("tag_list").([]interface{})

	if len(eventTypes) > 0 {
		for _, v := range eventTypes {
			queryParams = fmt.Sprintf("%s&event_types=%v", queryParams, v)
		}
	}
	if len(eventClassIds) > 0 {
		for _, v := range eventClassIds {
			queryParams = fmt.Sprintf("%s&event_class_ids=%v", queryParams, v)
		}
	}
	if len(severityList) > 0 {
		for _, v := range severityList {
			queryParams = fmt.Sprintf("%s&severity_list=%v", queryParams, v)
		}
	}
	if len(tagList) > 0 {
		for _, v := range tagList {
			queryParams = fmt.Sprintf("%s&tag_list=%v", queryParams, v)
		}
	}

	return queryParams
}

func dataSourceEventIntrusionEventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/event/events"
		result  = make([]interface{}, 0)
		offset  = 0
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildEventIntrusionEventsQueryParams(d, epsId)
	getPath += buildEventIntrusionEventsListQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%d", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving intrusion events: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenEventIntrusionEvents(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEventIntrusionEvents(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		result = append(result, map[string]interface{}{
			"event_id":            utils.PathSearch("event_id", v, nil),
			"event_class_id":      utils.PathSearch("event_class_id", v, nil),
			"event_type":          utils.PathSearch("event_type", v, nil),
			"event_name":          utils.PathSearch("event_name", v, nil),
			"severity":            utils.PathSearch("severity", v, nil),
			"container_name":      utils.PathSearch("container_name", v, nil),
			"image_name":          utils.PathSearch("image_name", v, nil),
			"host_name":           utils.PathSearch("host_name", v, nil),
			"host_id":             utils.PathSearch("host_id", v, nil),
			"private_ip":          utils.PathSearch("private_ip", v, nil),
			"public_ip":           utils.PathSearch("public_ip", v, nil),
			"os_type":             utils.PathSearch("os_type", v, nil),
			"host_status":         utils.PathSearch("host_status", v, nil),
			"agent_status":        utils.PathSearch("agent_status", v, nil),
			"protect_status":      utils.PathSearch("protect_status", v, nil),
			"asset_value":         utils.PathSearch("asset_value", v, nil),
			"attack_phase":        utils.PathSearch("attack_phase", v, nil),
			"attack_tag":          utils.PathSearch("attack_tag", v, nil),
			"occur_time":          utils.PathSearch("occur_time", v, nil),
			"handle_time":         utils.PathSearch("handle_time", v, nil),
			"handle_status":       utils.PathSearch("handle_status", v, nil),
			"handle_method":       utils.PathSearch("handle_method", v, nil),
			"handler":             utils.PathSearch("handler", v, nil),
			"operate_accept_list": utils.PathSearch("operate_accept_list", v, nil),
			"operate_detail_list": flattenEventOperateDetailList(utils.PathSearch("operate_detail_list", v,
				make([]interface{}, 0)).([]interface{})),
			"forensic_info":  utils.JsonToString(utils.PathSearch("forensic_info", v, nil)),
			"resource_info":  flattenEventResourceInfo(utils.PathSearch("resource_info", v, nil)),
			"geo_info":       utils.JsonToString(utils.PathSearch("geo_info", v, nil)),
			"malware_info":   utils.JsonToString(utils.PathSearch("malware_info", v, nil)),
			"network_info":   utils.JsonToString(utils.PathSearch("network_info", v, nil)),
			"app_info":       utils.JsonToString(utils.PathSearch("app_info", v, nil)),
			"system_info":    utils.JsonToString(utils.PathSearch("system_info", v, nil)),
			"extend_info":    utils.JsonToString(utils.PathSearch("extend_info", v, nil)),
			"recommendation": utils.PathSearch("recommendation", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"event_abstract": utils.PathSearch("event_abstract", v, nil),
			"process_info_list": flattenEventProcessInfo(utils.PathSearch("process_info_list", v,
				make([]interface{}, 0)).([]interface{})),
			"user_info_list": flattenEventUserInfo(utils.PathSearch("user_info_list", v,
				make([]interface{}, 0)).([]interface{})),
			"file_info_list": flattenEventFileInfo(utils.PathSearch("file_info_list", v,
				make([]interface{}, 0)).([]interface{})),
			"event_details": utils.PathSearch("event_details", v, nil),
			"tag_list":      utils.PathSearch("tag_list", v, nil),
			"event_count":   utils.PathSearch("event_count", v, nil),
			"operate_type":  utils.PathSearch("operate_type", v, nil),
		})
	}

	return result
}

func flattenEventOperateDetailList(detailList []interface{}) []interface{} {
	if len(detailList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(detailList))
	for _, v := range detailList {
		rst = append(rst, map[string]interface{}{
			"agent_id":        utils.PathSearch("agent_id", v, nil),
			"process_pid":     utils.PathSearch("process_pid", v, nil),
			"is_parent":       utils.PathSearch("is_parent", v, nil),
			"file_hash":       utils.PathSearch("file_hash", v, nil),
			"file_path":       utils.PathSearch("file_path", v, nil),
			"file_attr":       utils.PathSearch("file_attr", v, nil),
			"private_ip":      utils.PathSearch("private_ip", v, nil),
			"login_ip":        utils.PathSearch("login_ip", v, nil),
			"login_user_name": utils.PathSearch("login_user_name", v, nil),
			"keyword":         utils.PathSearch("keyword", v, nil),
			"hash":            utils.PathSearch("hash", v, nil),
		})
	}

	return rst
}

func flattenEventResourceInfo(resourceInfo interface{}) []map[string]interface{} {
	if resourceInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"domain_id":             utils.PathSearch("domain_id", resourceInfo, nil),
		"project_id":            utils.PathSearch("project_id", resourceInfo, nil),
		"enterprise_project_id": utils.PathSearch("enterprise_project_id", resourceInfo, nil),
		"region_name":           utils.PathSearch("region_name", resourceInfo, nil),
		"vpc_id":                utils.PathSearch("vpc_id", resourceInfo, nil),
		"cloud_id":              utils.PathSearch("cloud_id", resourceInfo, nil),
		"vm_name":               utils.PathSearch("vm_name", resourceInfo, nil),
		"vm_uuid":               utils.PathSearch("vm_uuid", resourceInfo, nil),
		"container_id":          utils.PathSearch("container_id", resourceInfo, nil),
		"container_status":      utils.PathSearch("container_status", resourceInfo, nil),
		"pod_uid":               utils.PathSearch("pod_uid", resourceInfo, nil),
		"pod_name":              utils.PathSearch("pod_name", resourceInfo, nil),
		"namespace":             utils.PathSearch("namespace", resourceInfo, nil),
		"cluster_id":            utils.PathSearch("cluster_id", resourceInfo, nil),
		"cluster_name":          utils.PathSearch("cluster_name", resourceInfo, nil),
		"image_id":              utils.PathSearch("image_id", resourceInfo, nil),
		"image_name":            utils.PathSearch("image_name", resourceInfo, nil),
		"host_attr":             utils.PathSearch("host_attr", resourceInfo, nil),
		"service":               utils.PathSearch("service", resourceInfo, nil),
		"micro_service":         utils.PathSearch("micro_service", resourceInfo, nil),
		"sys_arch":              utils.PathSearch("sys_arch", resourceInfo, nil),
		"os_bit":                utils.PathSearch("os_bit", resourceInfo, nil),
		"os_type":               utils.PathSearch("os_type", resourceInfo, nil),
		"os_name":               utils.PathSearch("os_name", resourceInfo, nil),
		"os_version":            utils.PathSearch("os_version", resourceInfo, nil),
	}

	return []map[string]interface{}{result}
}

func flattenEventProcessInfo(processInfo []interface{}) []interface{} {
	if len(processInfo) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(processInfo))
	for _, v := range processInfo {
		rst = append(rst, map[string]interface{}{
			"process_name":              utils.PathSearch("process_name", v, nil),
			"process_path":              utils.PathSearch("process_path", v, nil),
			"process_pid":               utils.PathSearch("process_pid", v, nil),
			"process_uid":               utils.PathSearch("process_uid", v, nil),
			"process_username":          utils.PathSearch("process_username", v, nil),
			"process_cmdline":           utils.PathSearch("process_cmdline", v, nil),
			"process_filename":          utils.PathSearch("process_filename", v, nil),
			"process_start_time":        utils.PathSearch("process_start_time", v, nil),
			"process_gid":               utils.PathSearch("process_gid", v, nil),
			"process_egid":              utils.PathSearch("process_egid", v, nil),
			"process_euid":              utils.PathSearch("process_euid", v, nil),
			"ancestor_process_path":     utils.PathSearch("ancestor_process_path", v, nil),
			"ancestor_process_pid":      utils.PathSearch("ancestor_process_pid", v, nil),
			"ancestor_process_cmdline":  utils.PathSearch("ancestor_process_cmdline", v, nil),
			"parent_process_name":       utils.PathSearch("parent_process_name", v, nil),
			"parent_process_path":       utils.PathSearch("parent_process_path", v, nil),
			"parent_process_pid":        utils.PathSearch("parent_process_pid", v, nil),
			"parent_process_uid":        utils.PathSearch("parent_process_uid", v, nil),
			"parent_process_cmdline":    utils.PathSearch("parent_process_cmdline", v, nil),
			"parent_process_filename":   utils.PathSearch("parent_process_filename", v, nil),
			"parent_process_start_time": utils.PathSearch("parent_process_start_time", v, nil),
			"parent_process_gid":        utils.PathSearch("parent_process_gid", v, nil),
			"parent_process_egid":       utils.PathSearch("parent_process_egid", v, nil),
			"parent_process_euid":       utils.PathSearch("parent_process_euid", v, nil),
			"child_process_name":        utils.PathSearch("child_process_name", v, nil),
			"child_process_path":        utils.PathSearch("child_process_path", v, nil),
			"child_process_pid":         utils.PathSearch("child_process_pid", v, nil),
			"child_process_uid":         utils.PathSearch("child_process_uid", v, nil),
			"child_process_cmdline":     utils.PathSearch("child_process_cmdline", v, nil),
			"child_process_filename":    utils.PathSearch("child_process_filename", v, nil),
			"child_process_start_time":  utils.PathSearch("child_process_start_time", v, nil),
			"child_process_gid":         utils.PathSearch("child_process_gid", v, nil),
			"child_process_egid":        utils.PathSearch("child_process_egid", v, nil),
			"child_process_euid":        utils.PathSearch("child_process_euid", v, nil),
			"virt_cmd":                  utils.PathSearch("virt_cmd", v, nil),
			"virt_process_name":         utils.PathSearch("virt_process_name", v, nil),
			"escape_mode":               utils.PathSearch("escape_mode", v, nil),
			"escape_cmd":                utils.PathSearch("escape_cmd", v, nil),
			"process_hash":              utils.PathSearch("process_hash", v, nil),
			"process_file_hash":         utils.PathSearch("process_file_hash", v, nil),
			"parent_process_file_hash":  utils.PathSearch("parent_process_file_hash", v, nil),
			"block":                     utils.PathSearch("block", v, nil),
		})
	}

	return rst
}

func flattenEventUserInfo(userInfo []interface{}) []interface{} {
	if len(userInfo) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(userInfo))
	for _, v := range userInfo {
		rst = append(rst, map[string]interface{}{
			"user_id":            utils.PathSearch("user_id", v, nil),
			"user_gid":           utils.PathSearch("user_gid", v, nil),
			"user_name":          utils.PathSearch("user_name", v, nil),
			"user_group_name":    utils.PathSearch("user_group_name", v, nil),
			"user_home_dir":      utils.PathSearch("user_home_dir", v, nil),
			"login_ip":           utils.PathSearch("login_ip", v, nil),
			"service_type":       utils.PathSearch("service_type", v, nil),
			"service_port":       utils.PathSearch("service_port", v, nil),
			"login_mode":         utils.PathSearch("login_mode", v, nil),
			"login_last_time":    utils.PathSearch("login_last_time", v, nil),
			"login_fail_count":   utils.PathSearch("login_fail_count", v, nil),
			"pwd_hash":           utils.PathSearch("pwd_hash", v, nil),
			"pwd_with_fuzzing":   utils.PathSearch("pwd_with_fuzzing", v, nil),
			"pwd_used_days":      utils.PathSearch("pwd_used_days", v, nil),
			"pwd_min_days":       utils.PathSearch("pwd_min_days", v, nil),
			"pwd_max_days":       utils.PathSearch("pwd_max_days", v, nil),
			"pwd_warn_left_days": utils.PathSearch("pwd_warn_left_days", v, nil),
		})
	}

	return rst
}

func flattenEventFileInfo(fileInfo []interface{}) []interface{} {
	if len(fileInfo) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(fileInfo))
	for _, v := range fileInfo {
		rst = append(rst, map[string]interface{}{
			"file_path":        utils.PathSearch("file_path", v, nil),
			"file_alias":       utils.PathSearch("file_alias", v, nil),
			"file_size":        utils.PathSearch("file_size", v, nil),
			"file_mtime":       utils.PathSearch("file_mtime", v, nil),
			"file_atime":       utils.PathSearch("file_atime", v, nil),
			"file_ctime":       utils.PathSearch("file_ctime", v, nil),
			"file_hash":        utils.PathSearch("file_hash", v, nil),
			"file_md5":         utils.PathSearch("file_md5", v, nil),
			"file_sha256":      utils.PathSearch("file_sha256", v, nil),
			"file_type":        utils.PathSearch("file_type", v, nil),
			"file_content":     utils.PathSearch("file_content", v, nil),
			"file_attr":        utils.PathSearch("file_attr", v, nil),
			"file_operation":   utils.PathSearch("file_operation", v, nil),
			"file_action":      utils.PathSearch("file_action", v, nil),
			"file_change_attr": utils.PathSearch("file_change_attr", v, nil),
			"file_new_path":    utils.PathSearch("file_new_path", v, nil),
			"file_desc":        utils.PathSearch("file_desc", v, nil),
			"file_key_word":    utils.PathSearch("file_key_word", v, nil),
			"is_dir":           utils.PathSearch("is_dir", v, nil),
			"fd_info":          utils.PathSearch("fd_info", v, nil),
			"fd_count":         utils.PathSearch("fd_count", v, nil),
		})
	}

	return rst
}
