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

// @API HSS GET /v5/{project_id}/backup/{policy_id}
func DataSourceBackupPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupPolicyRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource. If omitted, the provider-level region will be used.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the backup policy ID.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the enterprise project ID.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the backup policy is enabled.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the backup policy.",
			},
			"operation_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The backup type.",
			},
			"operation_definition": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataOperationDefinitionSchema(),
				Description: "The policy attribute.",
			},
			"trigger": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataTriggerSchema(),
				Description: "The backup policy scheduling rule.",
			},
		},
	}
}

func dataOperationDefinitionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"day_backups": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of daily backups that can be retained.",
			},
			"max_backups": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of backups that can be automatically created for a backup object.",
			},
			"month_backups": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of monthly backups that can be retained.",
			},
			"retention_duration_days": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The duration of retaining a backup, in days.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time zone where the user is located.",
			},
			"week_backups": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of weekly backups that can be retained.",
			},
			"year_backups": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of yearly backups that can be retained.",
			},
		},
	}
}

func dataTriggerSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scheduler ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scheduler name.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scheduler type.",
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataTriggerPropertiesSchema(),
				Description: "The scheduler attribute.",
			},
		},
	}
}

func dataTriggerPropertiesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pattern": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The scheduling policy.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scheduler start time.",
			},
		},
	}
}

func buildBackupPolicyQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	return ""
}

func dataSourceBackupPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v5/{project_id}/backup/{policy_id}"
		policyId = d.Get("policy_id").(string)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", policyId)
	getPath += buildBackupPolicyQueryParams(d, cfg)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the HSS backup policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("enabled", utils.PathSearch("enabled", respBody, nil)),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("operation_type", utils.PathSearch("operation_type", respBody, nil)),
		d.Set("operation_definition", flattenOperationDefinitionAttribute(utils.PathSearch("operation_definition", respBody, nil))),
		d.Set("trigger", flattenTriggerAttribute(utils.PathSearch("trigger", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOperationDefinitionAttribute(rawMap interface{}) []interface{} {
	if rawMap == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"day_backups":             utils.PathSearch("day_backups", rawMap, nil),
			"max_backups":             utils.PathSearch("max_backups", rawMap, nil),
			"month_backups":           utils.PathSearch("month_backups", rawMap, nil),
			"retention_duration_days": utils.PathSearch("retention_duration_days", rawMap, nil),
			"timezone":                utils.PathSearch("timezone", rawMap, nil),
			"week_backups":            utils.PathSearch("week_backups", rawMap, nil),
			"year_backups":            utils.PathSearch("year_backups", rawMap, nil),
		},
	}
}

func flattenTriggerAttribute(rawMap interface{}) []interface{} {
	if rawMap == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":         utils.PathSearch("id", rawMap, nil),
			"name":       utils.PathSearch("name", rawMap, nil),
			"type":       utils.PathSearch("type", rawMap, nil),
			"properties": flattenTriggerPropertiesAttribute(utils.PathSearch("properties", rawMap, nil)),
		},
	}
}

func flattenTriggerPropertiesAttribute(rawMap interface{}) []interface{} {
	if rawMap == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"pattern":    utils.PathSearch("pattern", rawMap, nil),
			"start_time": utils.PathSearch("start_time", rawMap, nil),
		},
	}
}
