package swrenterprise

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}/artifacts/{reference}
func DataSourceSwrEnterpriseInstanceArtifactDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrEnterpriseInstanceArtifactDetailsRead,

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
				Description: `Specifies the enterprise instance ID.`,
			},
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace name.`,
			},
			"repository_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the repository name.`,
			},
			"reference": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the artifact digest.`,
			},
			"with_scan_overview": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  `Specifies whether to return the scan overview infos.`,
			},
			"artifact_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the artifact ID.`,
			},
			"namespace_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the namespace ID.`,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the repository ID.`,
			},
			"media_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the media type.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the artifact size, unit is byte.`,
			},
			"digest": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the digest.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the artifact type.`,
			},
			"manifest_media_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the manifest media type.`,
			},
			"pull_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last pull time.`,
			},
			"push_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last push time.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the artifact version tags.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the tag ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the tag name.`,
						},
						"repository_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the repository ID.`,
						},
						"artifact_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the artifact ID.`,
						},
						"push_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the push time of the tag.`,
						},
						"pull_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the pull time of the tag.`,
						},
					},
				},
			},
			"scan_overview": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the report scan overview.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the report type.`,
						},
						"overview": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the report content.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"report_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the report ID.`,
									},
									"scan_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the scan status.`,
									},
									"severity": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the severity.`,
									},
									"duration": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the duration.`,
									},
									"start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the start time of the report.`,
									},
									"end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the end time of the report.`,
									},
									"complete_percent": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the completed percent.`,
									},
									"scanner": resourceSchemeArtifactDetailsContentScanner(),
									"summary": resourceSchemeArtifactDetailsContentVulnerabilities(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSchemeArtifactDetailsContentScanner() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: `Indicates the scanner infos.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `Indicates the scanner name.`,
				},
				"vendor": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `Indicates the vendor of the scanner.`,
				},
				"version": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `Indicates the version of the scanner.`,
				},
			},
		},
	}
}

func resourceSchemeArtifactDetailsContentVulnerabilities() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: `Indicates the vulnerabilities summary.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"total": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: `Indicates the vulnerability counts.`,
				},
				"fixable": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: `Indicates the fixable count of the vulnerability.`,
				},
				"summary": {
					Type:        schema.TypeMap,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeInt},
					Description: `Indicates the summary of the different level vulnerability.`,
				},
			},
		},
	}
}

func dataSourceSwrEnterpriseInstanceArtifactDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}/artifacts/{reference}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{namespace_name}", d.Get("namespace_name").(string))
	getPath = strings.ReplaceAll(getPath, "{repository_name}", url.PathEscape(strings.ReplaceAll(d.Get("repository_name").(string), "/", "%2F")))
	getPath = strings.ReplaceAll(getPath, "{reference}", d.Get("reference").(string))

	if v, ok := d.GetOk("with_scan_overview"); ok {
		getPath += fmt.Sprintf("?with_scan_overview=%v", v)
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error querying SWR artifact details: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flattening SWR artifact details response: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("artifact_id", utils.PathSearch("id", getRespBody, nil)),
		d.Set("namespace_id", utils.PathSearch("namespace_id", getRespBody, nil)),
		d.Set("repository_id", utils.PathSearch("repository_id", getRespBody, nil)),
		d.Set("media_type", utils.PathSearch("media_type", getRespBody, nil)),
		d.Set("size", utils.PathSearch("size", getRespBody, nil)),
		d.Set("digest", utils.PathSearch("digest", getRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRespBody, nil)),
		d.Set("manifest_media_type", utils.PathSearch("manifest_media_type", getRespBody, nil)),
		d.Set("pull_time", utils.PathSearch("pull_time", getRespBody, nil)),
		d.Set("push_time", utils.PathSearch("push_time", getRespBody, nil)),
		d.Set("tags", flattenArtifactDetailsTags(
			utils.PathSearch("tags", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("scan_overview", flattenArtifactDetailsScanOverview(utils.PathSearch("scan_overview", getRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenArtifactDetailsTags(paramsList []interface{}) []map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"id":            utils.PathSearch("id", params, nil),
			"name":          utils.PathSearch("name", params, nil),
			"repository_id": utils.PathSearch("repository_id", params, nil),
			"artifact_id":   utils.PathSearch("artifact_id", params, nil),
			"push_time":     utils.PathSearch("push_time", params, nil),
			"pull_time":     utils.PathSearch("pull_time", params, nil),
		}

		rst = append(rst, m)
	}

	return rst
}

func flattenArtifactDetailsScanOverview(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}

	reports := params.(map[string]interface{})
	results := make([]map[string]interface{}, 0)

	for reportType, report := range reports {
		if strings.Contains(reportType, "vulnerability.report") {
			results = append(results, map[string]interface{}{
				"type":     reportType,
				"overview": flattenArtifactVulnerabilitiesOverview(report),
			})
		}
	}

	return results
}

func flattenArtifactVulnerabilitiesOverview(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}

	rst := map[string]interface{}{
		"report_id":        utils.PathSearch("report_id", params, nil),
		"scan_status":      utils.PathSearch("scan_status", params, nil),
		"severity":         utils.PathSearch("severity", params, nil),
		"duration":         utils.PathSearch("duration", params, nil),
		"start_time":       utils.PathSearch("start_time", params, nil),
		"end_time":         utils.PathSearch("end_time", params, nil),
		"complete_percent": utils.PathSearch("complete_percent", params, nil),
		"scanner":          flattenArtifactDetailsContentScanner(utils.PathSearch("scanner", params, nil)),
		"summary":          flattenArtifactDetailsOverviewSummary(utils.PathSearch("summary", params, nil)),
	}

	return []map[string]interface{}{rst}
}

func flattenArtifactDetailsContentScanner(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}

	rst := map[string]interface{}{
		"name":    utils.PathSearch("name", params, nil),
		"vendor":  utils.PathSearch("vendor", params, nil),
		"version": utils.PathSearch("version", params, nil),
	}

	return []map[string]interface{}{rst}
}

func flattenArtifactDetailsOverviewSummary(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}

	rst := map[string]interface{}{
		"total":   utils.PathSearch("total", params, nil),
		"fixable": utils.PathSearch("fixable", params, nil),
		"summary": utils.PathSearch("summary", params, nil),
	}

	return []map[string]interface{}{rst}
}
