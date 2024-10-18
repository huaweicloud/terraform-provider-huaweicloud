package lts

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var logConvergeNotFoundCodes = []string{
	"LTS.2504", // The member does not found.
}

// @API LTS PUT /v1/{project_id}/lts/log-converge-config
// @API LTS GET /v1/{project_id}/lts/log-converge-config/{member_account_id}
func ResourceLogConverge() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogConvergeCreate,
		ReadContext:   resourceLogConvergeRead,
		UpdateContext: resourceLogConvergeUpdate,
		DeleteContext: resourceLogConvergeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceLogConvergeImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The organization ID to which the converged logs belong.`,
			},
			"management_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The administrator account ID used to manage log converge.`,
			},
			"member_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The member account ID to which the converged logs belong.`,
			},
			"log_mapping_config": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_log_group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the log group for source side.",
						},
						"target_log_group_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the log group for target side.",
						},
						"target_log_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The ID of the log group for target side.",
						},
						"log_stream_config": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_log_stream_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The ID of the log stream for source side.",
									},
									"target_log_stream_ttl": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The TTL of the log stream for target side.",
									},
									"target_log_stream_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the log group for target side.",
									},
									"target_log_stream_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The ID of the log stream for target side.",
									},
									"target_log_stream_eps_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The enterprise project ID of the log stream for target side.",
									},
								},
							},
							Description: `The log streams converged under the current log group.`,
						},
					},
				},
				Description: `The log converge configurations.`,
				// Since this structure uses full coverage logic, only the overall changes need to be identified.
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					oldConfig, newConfig := d.GetChange("log_mapping_config")
					return buildLogMappingConfigCompareObj(oldConfig.(*schema.Set)) == buildLogMappingConfigCompareObj(newConfig.(*schema.Set))
				},
			},
			"management_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The administrator project ID that required for first-time use.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the log converge configuration.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the log converge configuration.",
			},
		},
	}
}

// Use fields other than target_log_stream_id to calculate hash value.
// Since the target_log_stream_name field is a required field, and the ID is strongly related to the name (the service
// will verify it). If the ID is changed, the name will definitely change accordingly.
// So, only the ID needs to be calculated.
func buildLogStreamConfigCompareObj(obj *schema.Set) int {
	var buf bytes.Buffer
	for _, val := range obj.List() {
		m := val.(map[string]interface{})
		if m["source_log_stream_id"] != nil {
			buf.WriteString(fmt.Sprintf("%v-", m["source_log_stream_id"]))
		}
		if m["target_log_stream_ttl"] != nil {
			buf.WriteString(fmt.Sprintf("%v-", m["target_log_stream_ttl"]))
		}
		if m["target_log_stream_name"] != nil {
			buf.WriteString(fmt.Sprintf("%v-", m["target_log_stream_name"]))
		}
	}
	return hashcode.String(buf.String())
}

// Use fields other than target_log_group_id to calculate hash value.
// Since the target_log_group_name field is a required field, and the ID is strongly related to the name (the service
// will verify it). If the ID is changed, the name will definitely change accordingly.
// So, only the ID needs to be calculated.
func buildLogMappingConfigCompareObj(obj *schema.Set) int {
	var buf bytes.Buffer
	for _, val := range obj.List() {
		m := val.(map[string]interface{})
		if m["source_log_group_id"] != nil {
			buf.WriteString(fmt.Sprintf("%v-", m["source_log_group_id"]))
		}
		if m["target_log_group_name"] != nil {
			buf.WriteString(fmt.Sprintf("%v-", m["target_log_group_name"]))
		}
		if cfg := m["log_stream_config"]; cfg != nil {
			configs := cfg.(*schema.Set)
			buf.WriteString(fmt.Sprintf("%v-", buildLogStreamConfigCompareObj(configs)))
		}
	}
	return hashcode.String(buf.String())
}

func buildLogStreamConfigsBodyParams(oldStreamConfigs, newStreamConfigs *schema.Set) []interface{} {
	result := make([]interface{}, 0, newStreamConfigs.Len())

	for _, streamConfig := range newStreamConfigs.List() {
		newSourceLogStreamId := utils.PathSearch("source_log_stream_id", streamConfig, "").(string)
		newTargetLogStreamName := utils.PathSearch("target_log_stream_name", streamConfig, "").(string)
		oldStreamConfig := findTargetObjFromSetBySpecStr(oldStreamConfigs, "source_log_stream_id", newSourceLogStreamId)

		configElem := map[string]interface{}{
			// Required parameters
			"source_log_stream_id": newSourceLogStreamId,
			// If the group does not exist, the log group name that you want to create just ok.
			"target_log_stream_name": newTargetLogStreamName,
			"target_log_stream_ttl":  utils.PathSearch("target_log_stream_ttl", streamConfig, nil),
		}
		if newTargetStreamId := utils.PathSearch("target_log_stream_id", streamConfig, "").(string); newTargetStreamId != "" {
			// If the value of the parameter 'target_log_stream_id' in the new stream configuration is not empty, using this value.
			configElem["target_log_stream_id"] = newTargetStreamId
		} else if oldStreamConfig != nil && utils.PathSearch("target_log_stream_name", oldStreamConfig, "").(string) == newTargetLogStreamName {
			// Both new values of the parameter 'target_log_stream_name' and 'target_log_stream_name' are exist in the
			// old config, means the config is going to update, not create a new one.
			// Find the value of the parameter 'target_log_stream_id' in the old stream configuration.
			configElem["target_log_stream_id"] = utils.PathSearch("target_log_stream_id", oldStreamConfig, "")
		}
		result = append(result, configElem)
	}
	return result
}

func findTargetObjFromSetBySpecStr(targets *schema.Set, specKey, specStr string) interface{} {
	if targets.Len() > 0 {
		for _, val := range targets.List() {
			if utils.PathSearch(specKey, val, "").(string) == specStr {
				return val
			}
		}
	}
	return nil
}

func buildLogMappingConfigsBodyParams(oldMappingConfigs, newMappingConfigs *schema.Set) []interface{} {
	result := make([]interface{}, 0, newMappingConfigs.Len())

	for _, mappingConfig := range newMappingConfigs.List() {
		newSourceLogGroupId := utils.PathSearch("source_log_group_id", mappingConfig, "").(string)
		newTargetLogGroupName := utils.PathSearch("target_log_group_name", mappingConfig, "").(string)
		oldMappingConfig := findTargetObjFromSetBySpecStr(oldMappingConfigs, "source_log_group_id", newSourceLogGroupId)

		configElem := map[string]interface{}{
			// Required parameters ;
			"source_log_group_id":   newSourceLogGroupId,
			"target_log_group_name": newTargetLogGroupName,
			// Optional parameters
			"log_stream_config": buildLogStreamConfigsBodyParams(
				utils.PathSearch("log_stream_config", oldMappingConfig, schema.NewSet(schema.HashString, nil)).(*schema.Set),
				utils.PathSearch("log_stream_config", mappingConfig, schema.NewSet(schema.HashString, nil)).(*schema.Set),
			),
		}
		if newTargetLogGroupId := utils.PathSearch("target_log_group_id", mappingConfig, "").(string); newTargetLogGroupId != "" {
			// If the value of the parameter 'target_log_group_id' in the new mapping configuration is not empty, using this value.
			configElem["target_log_group_id"] = newTargetLogGroupId
		} else if oldMappingConfig != nil && utils.PathSearch("target_log_group_name", oldMappingConfig, "").(string) == newTargetLogGroupName {
			// Both new values of the parameter 'source_log_group_id' and 'target_log_group_name' are exist in the
			// old config, means the config is going to update, not create a new one.
			// If the group exist, the request body can only input the log group ID.
			configElem["target_log_group_id"] = utils.PathSearch("target_log_group_id", oldMappingConfig, "")
		}
		result = append(result, configElem)
	}
	return result
}

func buildModifyLogConvergeBodyParams(d *schema.ResourceData) map[string]interface{} {
	oldM, newM := d.GetChange("log_mapping_config")

	return map[string]interface{}{
		// Required parameters
		"organization_id":       d.Get("organization_id"),
		"management_account_id": d.Get("management_account_id"),
		"member_account_id":     d.Get("member_account_id"),
		"log_mapping_config":    buildLogMappingConfigsBodyParams(oldM.(*schema.Set), newM.(*schema.Set)),
		// Optional parameters
		"management_project_id": utils.ValueIgnoreEmpty(d.Get("management_project_id")),
	}
}

func modifyLogConvergeConfigs(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/lts/log-converge-config"

	modifyPath := client.Endpoint + httpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildModifyLogConvergeBodyParams(d)),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("PUT", modifyPath, &opts)
	return err
}

func resourceLogConvergeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		organizationId  = d.Get("organization_id").(string)
		memberAccountId = d.Get("member_account_id").(string)
	)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	err = modifyLogConvergeConfigs(client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", organizationId, memberAccountId))

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      logConvergeStateRefreshFunc(client, memberAccountId, []string{"done"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLogConvergeRead(ctx, d, meta)
}

func GetLogConvergeConfigsById(client *golangsdk.ServiceClient, memberAccountId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/lts/log-converge-config/{member_account_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{member_account_id}", memberAccountId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, common.ConvertExpected500ErrInto404Err(err, "error_code", logConvergeNotFoundCodes...)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return respBody, err
	}
	if len(utils.PathSearch("log_mapping_config", respBody, make([]interface{}, 0)).([]interface{})) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return respBody, nil
}

func flattenLogMappingConfigs(configs []interface{}) []interface{} {
	if len(configs) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(configs))
	for _, val := range configs {
		result = append(result, map[string]interface{}{
			"source_log_group_id":   utils.PathSearch("source_log_group_id", val, nil),
			"target_log_group_name": utils.PathSearch("target_log_group_name", val, nil),
			"target_log_group_id":   utils.PathSearch("target_log_group_id", val, nil),
			"log_stream_config":     flattenLogStreamConfigs(utils.PathSearch("log_stream_config", val, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenLogStreamConfigs(configs []interface{}) []interface{} {
	if len(configs) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(configs))
	for _, val := range configs {
		result = append(result, map[string]interface{}{
			"source_log_stream_id":     utils.PathSearch("source_log_stream_id", val, nil),
			"target_log_stream_ttl":    utils.PathSearch("target_log_stream_ttl", val, nil),
			"target_log_stream_name":   utils.PathSearch("target_log_stream_name", val, nil),
			"target_log_stream_id":     utils.PathSearch("target_log_stream_id", val, nil),
			"target_log_stream_eps_id": utils.PathSearch("target_log_stream_eps_id", val, nil),
		})
	}
	return result
}

func resourceLogConvergeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		memberAccountId = d.Get("member_account_id").(string)
	)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	respBody, err := GetLogConvergeConfigsById(client, memberAccountId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the configuration of log converge")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("organization_id", utils.PathSearch("organization_id", respBody, nil)),
		d.Set("member_account_id", utils.PathSearch("member_account_id", respBody, nil)),
		d.Set("management_account_id", utils.PathSearch("management_account_id", respBody, nil)),
		d.Set("log_mapping_config", flattenLogMappingConfigs(utils.PathSearch("log_mapping_config",
			respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("management_project_id", utils.PathSearch("management_project_id", respBody, nil)),
		// Attributes
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", respBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("update_time", respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLogConvergeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		memberAccountId = d.Get("member_account_id").(string)
	)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	err = modifyLogConvergeConfigs(client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      logConvergeStateRefreshFunc(client, memberAccountId, []string{"done"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLogConvergeRead(ctx, d, meta)
}

func buildDeleteLogConvergeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters
		"organization_id":       d.Get("organization_id"),
		"management_account_id": d.Get("management_account_id"),
		"member_account_id":     d.Get("member_account_id"),
		// Optional parameters
		"log_mapping_config": make([]interface{}, 0),
	}
}

func deleteLogConvergeConfigs(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/lts/log-converge-config"

	modifyPath := client.Endpoint + httpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteLogConvergeBodyParams(d),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("PUT", modifyPath, &opts)
	return err
}

func resourceLogConvergeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		memberAccountId = d.Get("member_account_id").(string)
	)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	err = deleteLogConvergeConfigs(client, d)
	if err != nil {
		return diag.Errorf("error deleting the configuration of LTS log converge: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      logConvergeStateRefreshFunc(client, memberAccountId, nil),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return diag.FromErr(err)
}

func logConvergeStateRefreshFunc(client *golangsdk.ServiceClient, memberAccountId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetLogConvergeConfigsById(client, memberAccountId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "Resource Not Found", "COMPLETED", nil
			}
			return resp, "ERROR", err
		}

		if utils.StrSliceContains(targets, utils.PathSearch("status", resp, "null").(string)) {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}

func resourceLogConvergeImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be '<organization_id>/<member_account_id>', "+
			"but got '%s'", importedId)
	}
	return []*schema.ResourceData{d}, d.Set("member_account_id", parts[1])
}
