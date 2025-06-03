package cbr

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var organizationPolicyNonUpdatableParams = []string{
	"operation_type",
}

// @API CBR POST /v3/{project_id}/organization-policies
// @API CBR GET /v3/{project_id}/organization-policies/{organization_policy_id}
// @API CBR PUT /v3/{project_id}/organization-policies/{organization_policy_id}
// @API CBR DELETE /v3/{project_id}/organization-policies/{organization_policy_id}
// @API CBR GET /v3/{project_id}/organization-policies
func ResourceOrganizationPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationPolicyCreate,
		ReadContext:   resourceOrganizationPolicyRead,
		UpdateContext: resourceOrganizationPolicyUpdate,
		DeleteContext: resourceOrganizationPolicyDelete,

		CustomizeDiff: config.FlexibleForceNew(organizationPolicyNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceOrganizationPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the organization policy is located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The organization policy name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The organization policy description.`,
			},
			"operation_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The organization policy type.`,
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The policy name.`,
			},
			"policy_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether the organization policy is enabled.`,
			},
			"policy_operation_definition": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day_backups": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Maximum number of daily backups that can be retained.`,
						},
						"destination_project_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The destination project ID for replication.`,
						},
						"destination_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The destination region for replication.`,
						},
						"enable_acceleration": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"true", "false",
							}, false),
							Description: `Whether to enable acceleration to shorten replication time for cross-region replication.`,
						},
						"max_backups": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Maximum number of backups that can be automatically created for a backup object.`,
						},
						"month_backups": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Maximum number of monthly backups that can be retained.`,
						},
						"retention_duration_days": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Duration of retaining a backup, in days.`,
						},
						"week_backups": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Maximum number of weekly backups that can be retained.`,
						},
						"year_backups": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Maximum number of yearly backups that can be retained.`,
						},
						"timezone": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The time zone where the user is located.`,
						},
						"full_backup_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Defines how often a full backup is performed after incremental backups.`,
						},
					},
				},
				Description: `The policy operation definition for backup and replication.`,
			},
			"policy_trigger": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"properties": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pattern": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The scheduling rules for policy execution.`,
									},
								},
							},
							Description: `The properties of policy trigger.`,
						},
					},
				},
				Description: `The policy execution time rule.`,
			},
			"effective_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The effective scope of the organization policy.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The organization policy status.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the account to which the organization policy belongs.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The account to which the organizational policy belongs.`,
			},
			// Internal parameter
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildOrganizationPolicyCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters
		"name":                        d.Get("name"),
		"operation_type":              d.Get("operation_type"),
		"policy_name":                 d.Get("policy_name"),
		"policy_enabled":              d.Get("policy_enabled"),
		"policy_operation_definition": buildPolicyOperationDefinition(d.Get("policy_operation_definition").([]interface{})),
		"policy_trigger":              buildPolicyTrigger(d.Get("policy_trigger").([]interface{})),
		// Optional parameters
		"description":     utils.ValueIgnoreEmpty(d.Get("description")),
		"effective_scope": utils.ValueIgnoreEmpty(d.Get("effective_scope")),
	}
}

func parsePolicyAccelerationEnabled(enabled string) interface{} {
	if enabled == "" {
		return nil
	}
	if enabled == "true" {
		return true
	}
	return false
}

func buildPolicyOperationDefinition(policyOpDefRaw []interface{}) map[string]interface{} {
	if len(policyOpDefRaw) < 1 {
		return nil
	}

	policyOpDef := policyOpDefRaw[0].(map[string]interface{})
	return map[string]interface{}{
		"day_backups":            utils.ValueIgnoreEmpty(utils.PathSearch("day_backups", policyOpDef, nil)),
		"destination_project_id": utils.ValueIgnoreEmpty(utils.PathSearch("destination_project_id", policyOpDef, nil)),
		"destination_region":     utils.ValueIgnoreEmpty(utils.PathSearch("destination_region", policyOpDef, nil)),
		"enable_acceleration": utils.ValueIgnoreEmpty(parsePolicyAccelerationEnabled(utils.PathSearch("enable_acceleration",
			policyOpDef, "").(string))),
		"max_backups":             utils.ValueIgnoreEmpty(utils.PathSearch("max_backups", policyOpDef, nil)),
		"month_backups":           utils.ValueIgnoreEmpty(utils.PathSearch("month_backups", policyOpDef, nil)),
		"retention_duration_days": utils.ValueIgnoreEmpty(utils.PathSearch("retention_duration_days", policyOpDef, nil)),
		"week_backups":            utils.ValueIgnoreEmpty(utils.PathSearch("week_backups", policyOpDef, nil)),
		"year_backups":            utils.ValueIgnoreEmpty(utils.PathSearch("year_backups", policyOpDef, nil)),
		"timezone":                utils.ValueIgnoreEmpty(utils.PathSearch("timezone", policyOpDef, nil)),
		"full_backup_interval":    utils.ValueIgnoreEmpty(utils.PathSearch("full_backup_interval", policyOpDef, nil)),
	}
}

func buildPolicyTrigger(policyTriggerRaw []interface{}) map[string]interface{} {
	if len(policyTriggerRaw) < 1 {
		return nil
	}

	policyTrigger := policyTriggerRaw[0].(map[string]interface{})

	return map[string]interface{}{
		"properties": buildPolicyTriggerProperties(utils.PathSearch("properties", policyTrigger,
			make([]interface{}, 0)).([]interface{})),
	}
}

func buildPolicyTriggerProperties(policyTriggerPropertiesRaw []interface{}) map[string]interface{} {
	if len(policyTriggerPropertiesRaw) < 1 {
		return nil
	}

	policyTriggerProperties := policyTriggerPropertiesRaw[0].(map[string]interface{})

	return map[string]interface{}{
		"pattern": utils.PathSearch("pattern", policyTriggerProperties, make([]interface{}, 0)).([]interface{}),
	}
}

func resourceOrganizationPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/organization-policies"
	)
	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"policy": buildOrganizationPolicyCreateBodyParams(d),
		}),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating organization policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("policy.id", respBody, "").(string)
	if policyId == "" {
		return diag.Errorf("failed to find the organization policy ID from the API response")
	}
	d.SetId(policyId)

	return resourceOrganizationPolicyRead(ctx, d, meta)
}

func GetOrganizationPolicyById(client *golangsdk.ServiceClient, policyId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/organization-policies/{organization_policy_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{organization_policy_id}", policyId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func flattenPolicyOperationDefinition(policyOpDefRaw interface{}) []map[string]interface{} {
	if policyOpDefRaw == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"day_backups":             utils.PathSearch("day_backups", policyOpDefRaw, nil),
			"destination_project_id":  utils.PathSearch("destination_project_id", policyOpDefRaw, nil),
			"destination_region":      utils.PathSearch("destination_region", policyOpDefRaw, nil),
			"enable_acceleration":     utils.PathSearch("enable_acceleration", policyOpDefRaw, nil),
			"max_backups":             utils.PathSearch("max_backups", policyOpDefRaw, nil),
			"month_backups":           utils.PathSearch("month_backups", policyOpDefRaw, nil),
			"retention_duration_days": utils.PathSearch("retention_duration_days", policyOpDefRaw, nil),
			"week_backups":            utils.PathSearch("week_backups", policyOpDefRaw, nil),
			"year_backups":            utils.PathSearch("year_backups", policyOpDefRaw, nil),
			"timezone":                utils.PathSearch("timezone", policyOpDefRaw, nil),
			"full_backup_interval":    utils.PathSearch("full_backup_interval", policyOpDefRaw, nil),
		},
	}
}

func flattenPolicyTrigger(policyTriggerRaw interface{}) []map[string]interface{} {
	if policyTriggerRaw == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"properties": flattenPolicyTriggerProperties(utils.PathSearch("properties", policyTriggerRaw, nil)),
		},
	}
}

func flattenPolicyTriggerProperties(policyTriggerPropertiesRaw interface{}) []map[string]interface{} {
	if policyTriggerPropertiesRaw == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"pattern": utils.PathSearch("pattern", policyTriggerPropertiesRaw, make([]interface{}, 0)).([]interface{}),
		},
	}
}

func resourceOrganizationPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		policyId = d.Id()
	)
	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	respBody, err := GetOrganizationPolicyById(client, policyId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting organization policy (%s)", policyId))
	}

	policy := utils.PathSearch("policy", respBody, nil)
	if policy == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error: organization policy not found")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", policy, nil)),
		d.Set("description", utils.PathSearch("description", policy, nil)),
		d.Set("operation_type", utils.PathSearch("operation_type", policy, nil)),
		d.Set("policy_name", utils.PathSearch("policy_name", policy, nil)),
		d.Set("policy_enabled", utils.PathSearch("policy_enabled", policy, nil)),
		d.Set("status", utils.PathSearch("status", policy, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", policy, nil)),
		d.Set("domain_name", utils.PathSearch("domain_name", policy, nil)),
		d.Set("policy_operation_definition", flattenPolicyOperationDefinition(utils.PathSearch("policy_operation_definition",
			policy, nil))),
		d.Set("policy_trigger", flattenPolicyTrigger(utils.PathSearch("policy_trigger", policy, nil))),
		d.Set("effective_scope", utils.PathSearch("effective_scope", policy, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildOrganizationPolicyUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":                        d.Get("name"),
		"description":                 utils.ValueIgnoreEmpty(d.Get("description")),
		"policy_name":                 d.Get("policy_name"),
		"policy_enabled":              d.Get("policy_enabled"),
		"policy_operation_definition": buildPolicyOperationDefinition(d.Get("policy_operation_definition").([]interface{})),
		"policy_trigger":              buildPolicyTrigger(d.Get("policy_trigger").([]interface{})),
		"effective_scope":             utils.ValueIgnoreEmpty(d.Get("effective_scope")),
	}
}

func resourceOrganizationPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v3/{project_id}/organization-policies/{organization_policy_id}"
		policyId = d.Id()
	)
	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{organization_policy_id}", policyId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"policy": buildOrganizationPolicyUpdateBodyParams(d),
		}),
	}

	_, err = client.Request("PUT", updatePath, &opt)
	if err != nil {
		return diag.Errorf("error updating organization policy (%s): %s", policyId, err)
	}

	return resourceOrganizationPolicyRead(ctx, d, meta)
}

func resourceOrganizationPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v3/{project_id}/organization-policies/{organization_policy_id}"
		policyId = d.Id()
	)
	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{organization_policy_id}", policyId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting organization policy (%s)", policyId))
	}
	return nil
}

func QueryOrganizationPolicies(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v3/{project_id}/organization-policies"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, fmt.Errorf("error querying organization policies: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("policies", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func getOrganizationPolicyByName(client *golangsdk.ServiceClient, name string) (interface{}, error) {
	policies, err := QueryOrganizationPolicies(client)
	if err != nil {
		return nil, err
	}

	policy := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", name), policies, nil)
	if policy == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/{project_id}/organization-policies",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the organization policy (%s) does not exist", name)),
			},
		}
	}
	return policy, nil
}

func resourceOrganizationPolicyImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		importId = d.Id()
	)

	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR client: %s", err)
	}

	// If the import ID is a UUID, use it directly
	if utils.IsUUID(importId) {
		return []*schema.ResourceData{d}, nil
	}

	// If not a UUID, treat as name and find the corresponding ID
	policy, err := getOrganizationPolicyByName(client, importId)
	if err != nil {
		return nil, fmt.Errorf("error getting organization policy by name (%s): %s", importId, err)
	}

	policyId := utils.PathSearch("id", policy, "").(string)
	if policyId == "" {
		return nil, fmt.Errorf("unable to find the organization policy ID by name (%s)", importId)
	}

	d.SetId(policyId)
	return []*schema.ResourceData{d}, nil
}
