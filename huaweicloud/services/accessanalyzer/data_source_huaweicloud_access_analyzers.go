package accessanalyzer

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AccessAnalyzer GET /v5/analyzers
func DataSourceAccessAnalyzers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccessAnalyzerRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the analyzer.`,
			},
			"analyzers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the analyzers.`,
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
						"configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unused_access": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unused_access_age": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_reason": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"details": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_analyzed_resource": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_resource_analyzed_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildListAccessAnalyzerParams(d *schema.ResourceData, marker string) string {
	res := ""
	if analyzerType, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, analyzerType)
	}

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func queryAccessAnalyzer(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v5/analyzers"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl

	queryParams := buildListAccessAnalyzerParams(d, "")

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestResp, err := client.Request("GET", listPath+queryParams, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving access analyzers: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		analyzers := utils.PathSearch("analyzers", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, analyzers...)

		marker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}

		queryParams = buildListAccessAnalyzerParams(d, marker)
	}
	return result, nil
}

func dataSourceAccessAnalyzerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("accessanalyzer", region)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer client: %s", err)
	}

	analyzers, err := queryAccessAnalyzer(client, d)
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
		d.Set("analyzers", flattenAccessAnalyzer(analyzers)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAccessAnalyzer(analyzers []interface{}) []interface{} {
	if len(analyzers) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(analyzers))
	for _, analyzer := range analyzers {
		result = append(result, map[string]interface{}{
			"id":                        utils.PathSearch("id", analyzer, nil),
			"name":                      utils.PathSearch("name", analyzer, nil),
			"type":                      utils.PathSearch("type", analyzer, nil),
			"configuration":             flattenConfiguration(utils.PathSearch("configuration", analyzer, nil)),
			"tags":                      utils.FlattenTagsToMap(utils.PathSearch("tags", analyzer, make([]interface{}, 0))),
			"status":                    utils.PathSearch("status", analyzer, nil),
			"status_reason":             flattenStatusReason(utils.PathSearch("status_reason", analyzer, nil)),
			"urn":                       utils.PathSearch("urn", analyzer, nil),
			"organization_id":           utils.PathSearch("organization_id", analyzer, nil),
			"last_analyzed_resource":    utils.PathSearch("last_analyzed_resource", analyzer, nil),
			"last_resource_analyzed_at": utils.PathSearch("last_resource_analyzed_at", analyzer, nil),
			"created_at":                utils.PathSearch("created_at", analyzer, nil),
		})
	}
	return result
}

func flattenStatusReason(statusReason interface{}) []map[string]interface{} {
	if statusReason == nil {
		return nil
	}

	res := []map[string]interface{}{
		{
			"code":    utils.PathSearch("code", statusReason, nil),
			"details": utils.PathSearch("details", statusReason, nil),
		},
	}

	return res
}
