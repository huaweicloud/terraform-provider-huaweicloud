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

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/alerts/search
func DataSourceAlerts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlertsRead,

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
				Description: `Specifies the ID of the workspace to which the alert belongs.`,
			},
			"from_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the search start time.`,
			},
			"to_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the search end time.`,
			},
			"condition": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the search condition expression.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Specifies the condition expression list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the expression name.`,
									},
									"data": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `Specifies the expression content.`,
									},
								},
							},
						},
						"logics": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the expression logic.`,
						},
					},
				},
			},
			"alerts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The alert list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alert ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alert name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alert description.`,
						},
						"type": {
							Type:        schema.TypeList,
							Elem:        alertsTypeElem(),
							Computed:    true,
							Description: `The alert type configuration.`,
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alert level.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alert status.`,
						},
						"data_source": {
							Type:        schema.TypeList,
							Elem:        alertsDataSourceElem(),
							Computed:    true,
							Description: `The data source configuration.`,
						},
						"first_occurrence_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The first occurrence time of the alert.`,
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user name of the owner.`,
						},
						"last_occurrence_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last occurrence time of the alert.`,
						},
						"planned_closure_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The planned closure time of the alert.`,
						},
						"verification_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The verification status.`,
						},
						"stage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The stage of the alert.`,
						},
						"debugging_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Whether it's a debugging data.`,
						},
						"labels": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The labels.`,
						},
						"close_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The close reason.`,
						},
						"close_comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The close comment.`,
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name creator name.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of the data source of an alert.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the account (domain_id) to whom the data is delivered and hosted.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of project where the account to whom the data is delivered and hosted belongs to.`,
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the region where the account to whom the data is delivered and hosted belongs to.`,
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the current workspace.`,
						},
						"arrive_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data receiving time.`,
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The times of the alert occurrences.`,
						},
						"data_class_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data class ID.`,
						},
						"environment": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The coordinates of the environment where the alert was generated.`,
							Elem:        dataObjectEnvironmentElem(),
						},
						"network_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The network information.`,
							Elem:        dataObjectNetworkListElem(),
						},
						"user_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The user information.`,
							Elem:        dataObjectUserInfoElem(),
						},
						"resource_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The affected resources.`,
							Elem:        dataObjectResourceListElem(),
						},
						"process": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The process information.`,
							Elem:        dataObjectProcessElem(),
						},
						"malware": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The malware information.`,
							Elem:        dataObjectMalwareElem(),
						},
						"remediation": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The remedy measure.`,
							Elem:        dataObjectRemediationElem(),
						},
						"file_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The file information.`,
							Elem:        dataObjectFileInfoElem(),
						},
					},
				},
			},
		},
	}
}

func alertsTypeElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The category.`,
			},
			"alert_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The alert type.`,
			},
		},
	}
	return &sc
}

func alertsDataSourceElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_feature": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product feature.`,
			},
			"product_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product name.`,
			},
			"source_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The source type.`,
			},
		},
	}
	return &sc
}

func dataObjectEnvironmentElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vendor_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The environment provider. The value can be **HWCP**, **HWC**, **AWS**, **Azure**, or **GCP**.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain ID.`,
			},
			"region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The region ID. **global** is returned for global services.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID. The default value is empty for global services.`,
			},
			"cross_workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the source workspace for the data delivery.`,
			},
		},
	}
}

func dataObjectNetworkListElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The protocol.`,
			},
			"src_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source IP address.`,
			},
			"src_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source port.`,
			},
			"dest_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The destination IP address.`,
			},
			"dest_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The destination port.`,
			},
			"src_geo": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The geographical location of the source IP address.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"latitude": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The latitude of the geographical location.`,
						},
						"longitude": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The longitude of the geographical location.`,
						},
						"city_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The city code.`,
						},
						"country_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The country code.`,
						},
					},
				},
			},
			"dest_geo": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The geographical location of the destination IP address.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"longitude": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The longitude of geographical location.`,
						},
						"latitude": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The latitude of the geographical location.`,
						},
						"city_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The city code.`,
						},
						"country_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The country code.`,
						},
					},
				},
			},
			"direction": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The direction. The value can be **IN** or **OUT**.`,
			},
			"src_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source domain name.`,
			},
			"dest_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The destination domain name.`,
			},
		},
	}
}

func dataObjectUserInfoElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user ID.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user name.`,
			},
		},
	}
}

func dataObjectResourceListElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource type.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the account to which the resource belongs.`,
			},
			"region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The region ID.`,
			},
			"ep_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID.`,
			},
			"ep_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project name.`,
			},
			"tags": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource tags.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource name.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the account to which the resource belongs.`,
			},
			"provider": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cloud service name, which is the same as the provider field in the RMS service.`,
			},
		},
	}
}

func dataObjectProcessElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"process_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The process execution file path.`,
			},
			"process_pid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The process ID.`,
			},
			"process_uid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The process user ID.`,
			},
			"process_cmdline": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The process command line.`,
			},
			"process_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The process name.`,
			},
			"process_parent_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent process name.`,
			},
			"process_parent_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent process execution file path.`,
			},
			"process_parent_uid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The parent process user ID.`,
			},
			"process_child_pid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The subprocess ID.`,
			},
			"process_child_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subprocess execution file path.`,
			},
			"process_parent_cmdline": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent process command line.`,
			},
			"process_child_cmdline": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subprocess command line.`,
			},
			"process_launche_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The process start time.`,
			},
			"process_terminate_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The process end time.`,
			},
			"process_parent_pid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The parent process ID.`,
			},
			"process_child_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subprocess name.`,
			},
			"process_child_uid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The subprocess user ID.`,
			},
		},
	}
}

func dataObjectMalwareElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"malware_family": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The malicious family.`,
			},
			"malware_class": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The malware category.`,
			},
		},
	}
}

func dataObjectRemediationElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The link to the general fix information for the incident.`,
			},
			"recommendation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The recommended solution.`,
			},
		},
	}
}

func dataObjectFileInfoElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"file_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file path.`,
			},
			"file_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file content.`,
			},
			"file_new_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file new path.`,
			},
			"file_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file hash value.`,
			},
			"file_md5": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file MD5 value.`,
			},
			"file_sha256": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file SHA256 value.`,
			},
			"file_attr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file attribute.`,
			},
		},
	}
}

func dataSourceAlertsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAlerts: Query the SecMaster alerts
	var (
		listAlertsHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/alerts/search"
		listAlertsProduct = "secmaster"
	)
	client, err := cfg.NewServiceClient(listAlertsProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	listAlertsPath := client.Endpoint + listAlertsHttpUrl
	listAlertsPath = strings.ReplaceAll(listAlertsPath, "{project_id}", client.ProjectID)
	listAlertsPath = strings.ReplaceAll(listAlertsPath, "{workspace_id}", fmt.Sprintf("%v", d.Get("workspace_id")))

	listAlertsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bodyParams, err := buildAlertBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	alerts := make([]interface{}, 0)
	offset := 0
	for {
		bodyParams["offset"] = offset
		listAlertsOpt.JSONBody = bodyParams
		listAlertsResp, err := client.Request("POST", listAlertsPath, &listAlertsOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listAlertsRespBody, err := utils.FlattenResponse(listAlertsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		data := utils.PathSearch("data", listAlertsRespBody, make([]interface{}, 0)).([]interface{})
		alerts = append(alerts, data...)

		offset += len(data)
		totalCount := utils.PathSearch("total", listAlertsRespBody, float64(0))
		if int(totalCount.(float64)) == offset {
			break
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("alerts", flattenAlertsResponseBody(alerts)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAlertBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	bodyParams := map[string]interface{}{
		"limit": 1000,
	}

	if v, ok := d.GetOk("from_date"); ok {
		fromDateWithZ, err := formatInputTime(v.(string))
		if err != nil {
			return nil, err
		}

		bodyParams["from_date"] = fromDateWithZ
	}
	if v, ok := d.GetOk("to_date"); ok {
		toDateWithZ, err := formatInputTime(v.(string))
		if err != nil {
			return nil, err
		}

		bodyParams["to_date"] = toDateWithZ
	}
	if v, ok := d.GetOk("condition.0"); ok {
		bodyParams["condition"] = v
	}

	return bodyParams, nil
}

func flattenAlertsResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	alerts := make([]interface{}, len(resp))
	for i, v := range resp {
		dataObject := utils.PathSearch("data_object", v, nil)
		alerts[i] = map[string]interface{}{
			"id":                    utils.PathSearch("id", dataObject, nil),
			"name":                  utils.PathSearch("title", dataObject, nil),
			"description":           utils.PathSearch("description", dataObject, nil),
			"type":                  flattenGetAlertResponseBodyAlertType(dataObject),
			"level":                 utils.PathSearch("severity", dataObject, nil),
			"status":                utils.PathSearch("handle_status", dataObject, nil),
			"owner":                 utils.PathSearch("owner", dataObject, nil),
			"data_source":           flattenGetAlertResponseBodyDataSource(dataObject),
			"first_occurrence_time": utils.PathSearch("first_observed_time", dataObject, nil),
			"last_occurrence_time":  utils.PathSearch("last_observed_time", dataObject, nil),
			"verification_status":   utils.PathSearch("verification_state", dataObject, nil),
			"stage":                 utils.PathSearch("ipdrr_phase", dataObject, nil),
			"debugging_data":        fmt.Sprintf("%v", utils.PathSearch("simulation", dataObject, nil)),
			"labels":                utils.PathSearch("labels", dataObject, nil),
			"close_reason":          utils.PathSearch("close_reason", dataObject, nil),
			"close_comment":         utils.PathSearch("close_comment", dataObject, nil),
			"creator":               utils.PathSearch("creator", dataObject, nil),
			"created_at":            utils.PathSearch("create_time", dataObject, nil),
			"updated_at":            utils.PathSearch("update_time", dataObject, nil),
			"planned_closure_time":  utils.PathSearch("sla", dataObject, nil),
			"version":               utils.PathSearch("version", dataObject, nil),
			"domain_id":             utils.PathSearch("domain_id", dataObject, nil),
			"project_id":            utils.PathSearch("project_id", dataObject, nil),
			"region_id":             utils.PathSearch("region_id", dataObject, nil),
			"workspace_id":          utils.PathSearch("workspace_id", dataObject, nil),
			"arrive_time":           utils.PathSearch("arrive_time", dataObject, nil),
			"count":                 utils.PathSearch("count", dataObject, 0),
			"data_class_id":         utils.PathSearch("dataclass_id", dataObject, nil),
			"environment":           flattenGetAlertResponseBodyEnvironment(dataObject),
			"network_list":          flattenGetAlertResponseBodyNetworkList(dataObject),
			"user_info":             flattenGetAlertResponseBodyUserInfo(dataObject),
			"resource_list":         flattenGetAlertResponseBodyResourceList(dataObject),
			"process":               flattenGetAlertResponseBodyProcess(dataObject),
			"malware":               flattenGetAlertResponseBodyMalware(dataObject),
			"remediation":           flattenGetAlertResponseBodyRemediation(dataObject),
			"file_info":             flattenGetAlertResponseBodyFileInfo(dataObject),
		}
	}
	return alerts
}

func flattenGetAlertResponseBodyEnvironment(resp interface{}) []interface{} {
	env := utils.PathSearch("environment", resp, nil)
	if env == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"vendor_type":        utils.PathSearch("vendor_type", env, nil),
			"cross_workspace_id": utils.PathSearch("cross_workspace_id", env, nil),
			"domain_id":          utils.PathSearch("domain_id", env, nil),
			"project_id":         utils.PathSearch("project_id", env, nil),
			"region_id":          utils.PathSearch("region_id", env, nil),
		},
	}

	return rst
}

func flattenGetAlertResponseBodyNetworkList(resp interface{}) []interface{} {
	array := utils.PathSearch("network_list", resp, make([]interface{}, 0)).([]interface{})
	if len(array) == 0 {
		return array
	}

	rst := make([]interface{}, len(array))
	// The type formats of `src_port` and `dest_port` are not consistent in the API document.
	// There is no test data to support the test. They are unified into strings to avoid type conversion errors.
	for i, v := range array {
		rst[i] = map[string]interface{}{
			"protocol":    utils.PathSearch("protocol", v, nil),
			"src_ip":      utils.PathSearch("src_ip", v, nil),
			"src_port":    fmt.Sprintf("%v", utils.PathSearch("src_port", v, nil)),
			"dest_ip":     utils.PathSearch("dest_ip", v, nil),
			"dest_port":   fmt.Sprintf("%v", utils.PathSearch("dest_port", v, nil)),
			"src_geo":     flattenGetAlertResponseBodyGeo(utils.PathSearch("src_geo", v, nil)),
			"dest_geo":    flattenGetAlertResponseBodyGeo(utils.PathSearch("dest_geo", v, nil)),
			"direction":   utils.PathSearch("direction", v, nil),
			"src_domain":  utils.PathSearch("src_domain", v, nil),
			"dest_domain": utils.PathSearch("dest_domain", v, nil),
		}
	}

	return rst
}

func flattenGetAlertResponseBodyGeo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"longitude":    utils.PathSearch("longitude", resp, float64(0)),
			"latitude":     utils.PathSearch("latitude", resp, float64(0)),
			"city_code":    utils.PathSearch("city_code", resp, nil),
			"country_code": utils.PathSearch("country_code", resp, nil),
		},
	}

	return rst
}

func flattenGetAlertResponseBodyUserInfo(resp interface{}) []interface{} {
	array := utils.PathSearch("user_info", resp, make([]interface{}, 0)).([]interface{})
	if len(array) == 0 {
		return array
	}

	rst := make([]interface{}, len(array))
	for i, v := range array {
		rst[i] = map[string]interface{}{
			"user_id":   utils.PathSearch("user_id", v, nil),
			"user_name": utils.PathSearch("user_name", v, nil),
		}
	}

	return rst
}

func flattenGetAlertResponseBodyResourceList(resp interface{}) []interface{} {
	array := utils.PathSearch("resource_list", resp, make([]interface{}, 0)).([]interface{})
	if len(array) == 0 {
		return array
	}

	rst := make([]interface{}, len(array))
	for i, v := range array {
		rst[i] = map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"type":       utils.PathSearch("type", v, nil),
			"ep_id":      utils.PathSearch("ep_id", v, nil),
			"ep_name":    utils.PathSearch("ep_name", v, nil),
			"domain_id":  utils.PathSearch("domain_id", v, nil),
			"project_id": utils.PathSearch("project_id", v, nil),
			"region_id":  utils.PathSearch("region_id", v, nil),
			"provider":   utils.PathSearch("provider", v, nil),
			"tags":       utils.PathSearch("tags", v, nil),
		}
	}

	return rst
}

func flattenGetAlertResponseBodyProcess(resp interface{}) []interface{} {
	array := utils.PathSearch("process", resp, make([]interface{}, 0)).([]interface{})
	if len(array) == 0 {
		return array
	}

	rst := make([]interface{}, len(array))
	for i, v := range array {
		rst[i] = map[string]interface{}{
			"process_path":           utils.PathSearch("process_path", v, nil),
			"process_pid":            utils.PathSearch("process_pid", v, 0),
			"process_uid":            utils.PathSearch("process_uid", v, 0),
			"process_cmdline":        utils.PathSearch("process_cmdline", v, nil),
			"process_name":           utils.PathSearch("process_name", v, nil),
			"process_parent_name":    utils.PathSearch("process_parent_name", v, nil),
			"process_parent_path":    utils.PathSearch("process_parent_path", v, nil),
			"process_parent_uid":     utils.PathSearch("process_parent_uid", v, 0),
			"process_child_pid":      utils.PathSearch("process_child_pid", v, 0),
			"process_child_path":     utils.PathSearch("process_child_path", v, nil),
			"process_parent_cmdline": utils.PathSearch("process_parent_cmdline", v, nil),
			"process_child_cmdline":  utils.PathSearch("process_child_cmdline", v, nil),
			"process_launche_time":   utils.PathSearch("process_launche_time", v, nil),
			"process_terminate_time": utils.PathSearch("process_terminate_time", v, nil),
			"process_parent_pid":     utils.PathSearch("process_parent_pid", v, 0),
			"process_child_name":     utils.PathSearch("process_child_name", v, nil),
			"process_child_uid":      utils.PathSearch("process_child_uid", v, 0),
		}
	}

	return rst
}

func flattenGetAlertResponseBodyMalware(resp interface{}) []interface{} {
	malware := utils.PathSearch("malware", resp, nil)
	if malware == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"malware_family": utils.PathSearch("malware_family", malware, nil),
			"malware_class":  utils.PathSearch("malware_class", malware, nil),
		},
	}

	return rst
}

func flattenGetAlertResponseBodyRemediation(resp interface{}) []interface{} {
	remediation := utils.PathSearch("remediation", resp, nil)
	if remediation == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"url":            utils.PathSearch("url", remediation, nil),
			"recommendation": utils.PathSearch("recommendation", remediation, nil),
		},
	}

	return rst
}

func flattenGetAlertResponseBodyFileInfo(resp interface{}) []interface{} {
	array := utils.PathSearch("file_info", resp, make([]interface{}, 0)).([]interface{})
	if len(array) == 0 {
		return array
	}

	rst := make([]interface{}, len(array))
	for i, v := range array {
		rst[i] = map[string]interface{}{
			"file_path":     utils.PathSearch("file_path", v, nil),
			"file_content":  utils.PathSearch("file_content", v, nil),
			"file_new_path": utils.PathSearch("file_new_path", v, nil),
			"file_hash":     utils.PathSearch("file_hash", v, nil),
			"file_md5":      utils.PathSearch("file_md5", v, nil),
			"file_sha256":   utils.PathSearch("file_sha256", v, nil),
			"file_attr":     utils.PathSearch("file_attr", v, nil),
		}
	}

	return rst
}
