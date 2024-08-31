package aom

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

// @API AOM GET /v2/{project_id}/alert/mute-rules
func DataSourceAlarmSilenceRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmSilenceRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"silence_time": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"starts_at": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ends_at": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"scope": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeInt},
										Computed: true,
									},
								},
							},
						},
						"silence_conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operate": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeList,
													Elem:     &schema.Schema{Type: schema.TypeString},
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlarmSilenceRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	results, err := listAlarmSilenceRules(client)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID")
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenAlarmSilenceRules(results.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listAlarmSilenceRules(client *golangsdk.ServiceClient) (interface{}, error) {
	listHttpUrl := "v2/{project_id}/alert/mute-rules"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving alarm silence rules: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening alarm silence rules: %s", err)
	}

	return listRespBody, nil
}

func flattenAlarmSilenceRules(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"name":               utils.PathSearch("name", rule, nil),
			"description":        utils.PathSearch("desc", rule, nil),
			"time_zone":          utils.PathSearch("timezone", rule, nil),
			"silence_time":       flattenSilenceRuleSilenceTime(rule),
			"silence_conditions": flattenSilenceRuleSilenceConditions(rule),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_time", rule, float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("update_time", rule, float64(0)).(float64))/1000, false),
		})
	}
	return result
}
