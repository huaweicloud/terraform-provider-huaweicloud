package lts

import (
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

func buildLogMappingStreamConfigsBodyParams(streamConfigs *schema.Set) []interface{} {
	if streamConfigs.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, streamConfigs.Len())
	for _, streamConfig := range streamConfigs.List() {
		result = append(result, map[string]interface{}{
			// Required parameters
			"source_log_stream_id":   utils.PathSearch("source_log_stream_id", streamConfig, nil),
			"target_log_stream_name": utils.PathSearch("target_log_stream_name", streamConfig, nil),
			"target_log_stream_ttl":  utils.PathSearch("target_log_stream_ttl", streamConfig, nil),
			// Optional parameters
			"target_log_stream_id": utils.ValueIgnoreEmpty(utils.PathSearch("target_log_stream_id", streamConfig, nil)),
		})
	}
	return result
}

func buildLogMappingConfigsBodyParams(mappingConfigs *schema.Set) []interface{} {
	if mappingConfigs.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, mappingConfigs.Len())
	for _, groupConfig := range mappingConfigs.List() {
		result = append(result, map[string]interface{}{
			// Required parameters
			"source_log_group_id":   utils.PathSearch("source_log_group_id", groupConfig, nil),
			"target_log_group_name": utils.PathSearch("target_log_group_name", groupConfig, nil),
			// Optional parameters
			"target_log_group_id": utils.ValueIgnoreEmpty(utils.PathSearch("target_log_group_id", groupConfig, nil)),
			"log_stream_config": buildLogMappingStreamConfigsBodyParams(utils.PathSearch("log_stream_config",
				groupConfig, schema.NewSet(schema.HashString, nil)).(*schema.Set)),
		})
	}
	return result
}

func buildModifyLogConvergeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters
		"organization_id":       d.Get("organization_id"),
		"management_account_id": d.Get("management_account_id"),
		"member_account_id":     d.Get("member_account_id"),
		// Optional parameters
		"log_mapping_config":    buildLogMappingConfigsBodyParams(d.Get("log_mapping_config").(*schema.Set)),
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
		return nil, parseQueryError500(err, logConvergeNotFoundCodes)
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
		d.Set("log_mapping_config", utils.PathSearch("log_mapping_config", respBody, nil)),
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
