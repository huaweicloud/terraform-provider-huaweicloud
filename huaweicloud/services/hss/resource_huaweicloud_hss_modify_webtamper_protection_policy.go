package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var modifyWebtamperPolicyNonUpdatableParams = []string{
	"host_id",
	"protect_dir_info",
	"protect_dir_info.*.protect_dir_list",
	"protect_dir_info.*.protect_dir_list.*.protect_dir",
	"protect_dir_info.*.protect_dir_list.*.local_backup_dir",
	"protect_dir_info.*.protect_dir_list.*.exclude_child_dir",
	"protect_dir_info.*.protect_dir_list.*.exclude_file_path",
	"protect_dir_info.*.exclude_file_type",
	"protect_dir_info.*.protect_mode",
	"enable_timing_off",
	"timing_off_config_info",
	"timing_off_config_info.*.week_off_list",
	"timing_off_config_info.*.timing_range_list",
	"timing_off_config_info.*.timing_range_list.*.time_range",
	"timing_off_config_info.*.timing_range_list.*.description",
	"enable_rasp_protect",
	"rasp_path",
	"enable_privileged_process",
	"privileged_process_info",
	"privileged_process_info.*.privileged_process_path_list",
	"privileged_process_info.*.privileged_child_status",
	"enterprise_project_id",
}

// @API HSS PUT /v5/{project_id}/webtamper/{host_id}/policy
func ResourceModifyWebtamperProtectionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModifyWebtamperProtectionPolicyCreate,
		ReadContext:   resourceModifyWebtamperProtectionPolicyRead,
		UpdateContext: resourceModifyWebtamperProtectionPolicyUpdate,
		DeleteContext: resourceModifyWebtamperProtectionPolicyDelete,

		CustomizeDiff: config.FlexibleForceNew(modifyWebtamperPolicyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protect_dir_info": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protect_dir_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protect_dir": {
										Type:     schema.TypeString,
										Required: true,
									},
									"local_backup_dir": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"exclude_child_dir": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"exclude_file_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"exclude_file_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protect_mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"enable_timing_off": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"timing_off_config_info": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"week_off_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"timing_range_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_range": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"enable_rasp_protect": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rasp_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_privileged_process": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"privileged_process_info": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"privileged_process_path_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"privileged_child_status": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildModifyWebtamperPolicyQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	return ""
}

func buildModifyWebtamperPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"protect_dir_info":          buildProtectDirInfoBodyParams(d.Get("protect_dir_info").([]interface{})),
		"enable_timing_off":         d.Get("enable_timing_off"),
		"timing_off_config_info":    buildTimeOffConfigInfoBodyParams(d.Get("timing_off_config_info").([]interface{})),
		"enable_rasp_protect":       d.Get("enable_rasp_protect"),
		"rasp_path":                 utils.ValueIgnoreEmpty(d.Get("rasp_path")),
		"enable_privileged_process": d.Get("enable_privileged_process"),
		"privileged_process_info":   buildPrivilegedProcessInfoBodyParams(d.Get("privileged_process_info").([]interface{})),
	}

	return bodyParams
}

func buildProtectDirInfoBodyParams(protectDirInfo []interface{}) map[string]interface{} {
	if len(protectDirInfo) == 0 {
		return nil
	}

	rawInfo, ok := protectDirInfo[0].(map[string]interface{})
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"protect_dir_list":  buildProtectDirListBodyParams(rawInfo["protect_dir_list"].([]interface{})),
		"exclude_file_type": utils.ValueIgnoreEmpty(rawInfo["exclude_file_type"]),
		"protect_mode":      utils.ValueIgnoreEmpty(rawInfo["protect_mode"]),
	}

	return bodyParams
}

func buildProtectDirListBodyParams(dirList []interface{}) []map[string]interface{} {
	if len(dirList) == 0 {
		return nil
	}

	dirInfo := make([]map[string]interface{}, 0, len(dirList))
	for _, v := range dirList {
		rawInfo, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"protect_dir":       rawInfo["protect_dir"],
			"local_backup_dir":  utils.ValueIgnoreEmpty(rawInfo["local_backup_dir"]),
			"exclude_child_dir": utils.ValueIgnoreEmpty(rawInfo["exclude_child_dir"]),
			"exclude_file_path": utils.ValueIgnoreEmpty(rawInfo["exclude_file_path"]),
		}
		dirInfo = append(dirInfo, params)
	}

	return dirInfo
}

func buildTimeOffConfigInfoBodyParams(configInfo []interface{}) map[string]interface{} {
	if len(configInfo) == 0 {
		return nil
	}

	rawInfo, ok := configInfo[0].(map[string]interface{})
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"week_off_list":     utils.ExpandToIntList(rawInfo["week_off_list"].([]interface{})),
		"timing_range_list": buildTimeRangeListBodyParams(rawInfo["timing_range_list"].([]interface{})),
	}

	return bodyParams
}

func buildTimeRangeListBodyParams(timeList []interface{}) []map[string]interface{} {
	if len(timeList) == 0 {
		return nil
	}

	rangeInfo := make([]map[string]interface{}, 0, len(timeList))
	for _, v := range timeList {
		rawInfo, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"time_range":  rawInfo["time_range"],
			"description": utils.ValueIgnoreEmpty(rawInfo["description"]),
		}
		rangeInfo = append(rangeInfo, params)
	}

	return rangeInfo
}

func buildPrivilegedProcessInfoBodyParams(processInfo []interface{}) map[string]interface{} {
	if len(processInfo) == 0 {
		return nil
	}

	rawInfo, ok := processInfo[0].(map[string]interface{})
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"privileged_process_path_list": utils.ExpandToStringList(rawInfo["privileged_process_path_list"].([]interface{})),
		"privileged_child_status":      rawInfo["privileged_child_status"],
	}

	return bodyParams
}

func resourceModifyWebtamperProtectionPolicyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/webtamper/{host_id}/policy"
		hostId  = d.Get("host_id").(string)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{host_id}", hostId)
	requestPath += buildModifyWebtamperPolicyQueryParams(epsId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildModifyWebtamperPolicyBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating the web tamper protection policy: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return nil
}

func resourceModifyWebtamperProtectionPolicyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceModifyWebtamperProtectionPolicyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceModifyWebtamperProtectionPolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to update the web tamper protection policy. Deleting
	  this resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
