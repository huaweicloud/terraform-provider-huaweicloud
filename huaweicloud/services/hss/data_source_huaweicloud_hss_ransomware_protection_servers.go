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

// @API HSS GET /v5/{project_id}/ransomware/servers
func DataSourceRansomwareProtectionServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRansomwareProtectionServersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ransom_protection_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protect_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agent_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
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
						"host_status": {
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
						"ransom_protection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ransom_protection_fail_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failed_decoy_dir": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_error": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_code": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"error_description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"backup_protection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count_protect_event": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"count_backuped": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"agent_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vault_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vault_allocated": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vault_charging_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_policy_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"resources_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildRansomwareProtectionServersQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("os_type"); ok {
		queryParams = fmt.Sprintf("%s&os_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_ip"); ok {
		queryParams = fmt.Sprintf("%s&host_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		queryParams = fmt.Sprintf("%s&private_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_status"); ok {
		queryParams = fmt.Sprintf("%s&host_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("ransom_protection_status"); ok {
		queryParams = fmt.Sprintf("%s&ransom_protection_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("protect_policy_name"); ok {
		queryParams = fmt.Sprintf("%s&protect_policy_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("policy_name"); ok {
		queryParams = fmt.Sprintf("%s&policy_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("policy_id"); ok {
		queryParams = fmt.Sprintf("%s&policy_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("agent_status"); ok {
		queryParams = fmt.Sprintf("%s&agent_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("group_id"); ok {
		queryParams = fmt.Sprintf("%s&group_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("group_name"); ok {
		queryParams = fmt.Sprintf("%s&group_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("last_days"); ok {
		queryParams = fmt.Sprintf("%s&last_days=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceRansomwareProtectionServersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/ransomware/servers"
		epsId   = cfg.GetEnterpriseProjectID(d)
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildRansomwareProtectionServersQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving the ransomware protection servers: %s", err)
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
		d.Set("data_list", flattenRansomwareProtectionServers(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRansomwareProtectionServers(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		result = append(result, map[string]interface{}{
			"host_id":                       utils.PathSearch("host_id", v, nil),
			"agent_id":                      utils.PathSearch("agent_id", v, nil),
			"host_name":                     utils.PathSearch("host_name", v, nil),
			"host_ip":                       utils.PathSearch("host_ip", v, nil),
			"private_ip":                    utils.PathSearch("private_ip", v, nil),
			"os_type":                       utils.PathSearch("os_type", v, nil),
			"os_name":                       utils.PathSearch("os_name", v, nil),
			"host_status":                   utils.PathSearch("host_status", v, nil),
			"project_id":                    utils.PathSearch("project_id", v, nil),
			"enterprise_project_id":         utils.PathSearch("enterprise_project_id", v, nil),
			"ransom_protection_status":      utils.PathSearch("ransom_protection_status", v, nil),
			"ransom_protection_fail_reason": utils.PathSearch("ransom_protection_fail_reason", v, nil),
			"failed_decoy_dir":              utils.PathSearch("failed_decoy_dir", v, nil),
			"agent_version":                 utils.PathSearch("agent_version", v, nil),
			"protect_status":                utils.PathSearch("protect_status", v, nil),
			"group_id":                      utils.PathSearch("group_id", v, nil),
			"group_name":                    utils.PathSearch("group_name", v, nil),
			"protect_policy_id":             utils.PathSearch("protect_policy_id", v, nil),
			"protect_policy_name":           utils.PathSearch("protect_policy_name", v, nil),
			"backup_error":                  flattenBackupErrorInfo(utils.PathSearch("backup_error", v, nil)),
			"backup_protection_status":      utils.PathSearch("backup_protection_status", v, nil),
			"count_protect_event":           utils.PathSearch("count_protect_event", v, nil),
			"count_backuped":                utils.PathSearch("count_backuped", v, nil),
			"agent_status":                  utils.PathSearch("agent_status", v, nil),
			"version":                       utils.PathSearch("version", v, nil),
			"host_source":                   utils.PathSearch("host_source", v, nil),
			"vault_id":                      utils.PathSearch("vault_id", v, nil),
			"vault_name":                    utils.PathSearch("vault_name", v, nil),
			"vault_size":                    utils.PathSearch("vault_size", v, nil),
			"vault_used":                    utils.PathSearch("vault_used", v, nil),
			"vault_allocated":               utils.PathSearch("vault_allocated", v, nil),
			"vault_charging_mode":           utils.PathSearch("vault_charging_mode", v, nil),
			"vault_status":                  utils.PathSearch("vault_status", v, nil),
			"backup_policy_id":              utils.PathSearch("backup_policy_id", v, nil),
			"backup_policy_name":            utils.PathSearch("backup_policy_name", v, nil),
			"backup_policy_enabled":         utils.PathSearch("backup_policy_enabled", v, nil),
			"resources_num":                 utils.PathSearch("resources_num", v, nil),
		})
	}

	return result
}

func flattenBackupErrorInfo(errorInfo interface{}) []map[string]interface{} {
	if errorInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"error_code":        utils.PathSearch("error_code", errorInfo, nil),
		"error_description": utils.PathSearch("error_description", errorInfo, nil),
	}

	return []map[string]interface{}{result}
}
