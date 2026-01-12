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

// @API HSS GET /v5/{project_id}/ransomware/backup/policies
func DataSourceRansomwareBackupPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRansomwareBackupPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_definition": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day_backups": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_backups": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"month_backups": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"retention_duration_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"timezone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"week_backups": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"year_backups": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"trigger": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"properties": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"pattern": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"start_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildRansomwareBackupPoliciesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("policy_id"); ok {
		queryParams = fmt.Sprintf("%s&policy_id=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

// Due to the ineffective API paging parameters, this data source is temporarily not using paging queries.
func dataSourceRansomwareBackupPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/ransomware/backup/policies"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildRansomwareBackupPoliciesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS ransomware backup policies: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", utils.PathSearch("total_num", respBody, nil)),
		d.Set("data_list", flattenRansomwareBackupPoliciesDataList(
			utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRansomwareBackupPoliciesDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"enabled":              utils.PathSearch("enabled", v, nil),
			"id":                   utils.PathSearch("id", v, nil),
			"name":                 utils.PathSearch("name", v, nil),
			"operation_type":       utils.PathSearch("operation_type", v, nil),
			"operation_definition": flattenOperationDefinition(utils.PathSearch("operation_definition", v, nil)),
			"trigger":              flattenTrigger(utils.PathSearch("trigger", v, nil)),
		})
	}

	return result
}

func flattenOperationDefinition(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"day_backups":             utils.PathSearch("day_backups", resp, nil),
			"max_backups":             utils.PathSearch("max_backups", resp, nil),
			"month_backups":           utils.PathSearch("month_backups", resp, nil),
			"retention_duration_days": utils.PathSearch("retention_duration_days", resp, nil),
			"timezone":                utils.PathSearch("timezone", resp, nil),
			"week_backups":            utils.PathSearch("week_backups", resp, nil),
			"year_backups":            utils.PathSearch("year_backups", resp, nil),
		},
	}
}

func flattenTrigger(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":         utils.PathSearch("id", resp, nil),
			"name":       utils.PathSearch("name", resp, nil),
			"type":       utils.PathSearch("type", resp, nil),
			"properties": flattenTriggerProperties(utils.PathSearch("properties", resp, nil)),
		},
	}
}

func flattenTriggerProperties(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"pattern":    utils.ExpandToStringList(utils.PathSearch("pattern", resp, make([]interface{}, 0)).([]interface{})),
			"start_time": utils.PathSearch("start_time", resp, nil),
		},
	}
}
