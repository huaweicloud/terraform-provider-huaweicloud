package dbss

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DBSS POST /v1/{project_id}/{instance_id}/dbss/audit/rule/sql-injections
func DataSourceDbssAuditSqlInjectionRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuditSqlInjectionRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the audit instance ID to which the SQL injection rules belong.`,
			},
			"risk_levels": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the risk level of the SQL injection rule.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the SQL injection rules.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the SQL injection rule.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the SQL injection rule.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the SQL injection rule.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the SQL injection rule.`,
						},
						"risk_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The risk level of the SQL injection rule.`,
						},
						"rank": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The rank of the SQL injection rule.`,
						},
						"feature": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SQL command characteristics of the SQL injection rule.`,
						},
						"regex": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The regular expression content of the SQL injection rule.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceAuditSqlInjectionRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		listRulesHttpUrl = "v1/{project_id}/{instance_id}/dbss/audit/rule/sql-injections"
		listRulesProduct = "dbss"
	)
	client, err := cfg.NewServiceClient(listRulesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}

	listRulesPath := client.Endpoint + listRulesHttpUrl
	listRulesPath = strings.ReplaceAll(listRulesPath, "{project_id}", client.ProjectID)
	listRulesPath = strings.ReplaceAll(listRulesPath, "{instance_id}", d.Get("instance_id").(string))

	listRulesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listRulesOpt.JSONBody = utils.RemoveNil(buildListSqlRulesBodyParams(d))

	listRulesResp, err := client.Request("POST", listRulesPath, &listRulesOpt)
	if err != nil {
		return diag.Errorf("error retrieving SQL injection rules: %s", err)
	}
	listRulesRespBody, err := utils.FlattenResponse(listRulesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	rules := utils.PathSearch("rules", listRulesRespBody, make([]interface{}, 0)).([]interface{})

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenListSqlRulesResponseBody(rules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListSqlRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"risk_levels": utils.ValueIgnoreEmpty(d.Get("risk_levels")),
	}

	return bodyParam
}

func flattenListSqlRulesResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"type":       utils.PathSearch("type", v, nil),
			"status":     utils.PathSearch("status", v, nil),
			"risk_level": utils.PathSearch("risk_level", v, nil),
			"rank":       utils.PathSearch("rank", v, nil),
			"feature":    utils.PathSearch("feature", v, nil),
			"regex":      utils.PathSearch("regex", v, nil),
		})
	}
	return rst
}
