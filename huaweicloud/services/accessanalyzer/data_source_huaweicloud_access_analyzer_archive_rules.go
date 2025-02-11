package accessanalyzer

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

// @API AccessAnalyzer GET /v5/analyzers/{analyzer_id}/archive-rules
func DataSourceAccessAnalyzerArchiveRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccessAnalyzerArchiveRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"analyzer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"archive_rules": {
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
						"filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"criterion": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"contains": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"eq": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"neq": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"exists": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"urn": {
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

func buildListAccessAnalyzerArchiveRuleParams(marker string) string {
	res := ""

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func queryAccessAnalyzerArchiveRules(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v5/analyzers/{analyzer_id}/archive-rules"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{analyzer_id}", d.Get("analyzer_id").(string))

	queryParams := buildListAccessAnalyzerArchiveRuleParams("")

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestResp, err := client.Request("GET", listPath+queryParams, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving access analyzer archive rules: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		archiveRules := utils.PathSearch("archive_rules", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, archiveRules...)

		marker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}

		queryParams = buildListAccessAnalyzerArchiveRuleParams(marker)
	}
	return result, nil
}

func dataSourceAccessAnalyzerArchiveRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("accessanalyzer", region)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer client: %s", err)
	}

	archiveRules, err := queryAccessAnalyzerArchiveRules(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("archive_rules", flattenAccessAnalyzerArchiveRules(archiveRules)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAccessAnalyzerArchiveRules(archiveRules []interface{}) []interface{} {
	if len(archiveRules) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(archiveRules))
	for _, archiveRule := range archiveRules {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", archiveRule, nil),
			"name":       utils.PathSearch("name", archiveRule, nil),
			"urn":        utils.PathSearch("urn", archiveRule, nil),
			"created_at": utils.PathSearch("created_at", archiveRule, nil),
			"updated_at": utils.PathSearch("updated_at", archiveRule, nil),
			"filters":    flattenArchiveRuleFilters(utils.PathSearch("filters", archiveRule, nil)),
		})
	}
	return result
}

func flattenArchiveRuleFilters(filtersRaw interface{}) []map[string]interface{} {
	if filtersRaw == nil {
		return nil
	}

	filters := filtersRaw.([]interface{})
	res := make([]map[string]interface{}, len(filters))
	for i, v := range filters {
		filter := v.(map[string]interface{})
		res[i] = map[string]interface{}{
			"key":       filter["key"],
			"criterion": flattenCriterion(filter["criterion"]),
		}
	}

	return res
}
