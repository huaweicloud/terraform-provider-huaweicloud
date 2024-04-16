package rds

import (
	"context"
	"encoding/json"
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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/auditlog-policy
func DataSourceRdsSqlAuditTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsSqlAuditTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS instance.`,
			},
			"operation_types": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the list of the operation type.`,
			},
			"operations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of the audit operations.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the type of the operation.`,
						},
						"actions": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `Indicates the list of the operation actions.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsSqlAuditTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)

	if err != nil {
		return diag.Errorf("error retrieving RDS SQL audit operations: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	operations, err := flattenAuditOperationsBody(d, getRespBody)
	if err != nil {
		return diag.FromErr(err)
	}
	mErr := multierror.Append(
		d.Set("region", cfg.GetRegion(d)),
		d.Set("operations", operations),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAuditOperationsBody(d *schema.ResourceData, resp interface{}) ([]interface{}, error) {
	if resp == nil {
		return nil, nil
	}
	allAuditLogAction := utils.PathSearch("all_audit_log_action", resp, nil)
	if allAuditLogAction == nil {
		return nil, fmt.Errorf("unable to find the all_audit_log_action from the response")
	}
	var auditTypes interface{}
	err := json.Unmarshal([]byte(allAuditLogAction.(string)), &auditTypes)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling RDS SQL audit types: %s", err)
	}
	operationsTypes := d.Get("operation_types").([]interface{})
	operationsTypeMap := make(map[string]bool)
	for _, operationsType := range operationsTypes {
		operationsTypeMap[operationsType.(string)] = true
	}
	rst := make([]interface{}, 0)
	for k, v := range auditTypes.(map[string]interface{}) {
		if len(operationsTypes) > 0 && !operationsTypeMap[k] {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"type":    k,
			"actions": strings.Split(v.(string), ","),
		})
	}
	return rst, nil
}
