package swrenterprise

import (
	"context"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// nolint:revive
// @API SWR GET /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}/artifacts/{reference}/vulnerabilities
func DataSourceSwrEnterpriseInstanceArtifactVulnerabilities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrEnterpriseInstanceArtifactVulnerabilitiesRead,

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
			"reports": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the reports.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the report type.`,
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the report content.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"generated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the generation time.`,
									},
									"severity": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the report severity.`,
									},
									"scanner":         resourceSchemeArtifactVulnerabilitiesContentScanner(),
									"vulnerabilities": resourceSchemeArtifactVulnerabilitiesContentVulnerabilities(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSchemeArtifactVulnerabilitiesContentScanner() *schema.Schema {
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

func resourceSchemeArtifactVulnerabilitiesContentVulnerabilities() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: `Indicates the vulnerabilities.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `Indicates the vulnerability ID.`,
				},
				"package": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `Indicates the package name which has vulnerability.`,
				},
				"version": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `Indicates the package version which has vulnerability.`,
				},
				"fix_version": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `Indicates the fixed package version of the vulnerability.`,
				},
				"severity": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `Indicates the severity of the vulnerability.`,
				},
				"description": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `Indicates the description of the vulnerability.`,
				},
				"links": {
					Type:        schema.TypeList,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: `Indicates the links of the vulnerability.`,
				},
				"artifact_digests": {
					Type:        schema.TypeList,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: `Indicates the digests containing this vulnerability.`,
				},
				"cwe_ids": {
					Type:        schema.TypeList,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: `Indicates the CWE IDs of the vulnerability.`,
				},
				"preferred_cvss": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: `Indicates the vulnerability scoring and attack analysis based on CVSS3 and CVSS2.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"score_v2": {
								Type:        schema.TypeFloat,
								Computed:    true,
								Description: `Indicates the vCVSS-2 score of the vulnerability.`,
							},
							"score_v3": {
								Type:        schema.TypeFloat,
								Computed:    true,
								Description: `Indicates the CVSS-3 score of the vulnerability.`,
							},
							"vector_v2": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: `Indicates the CVSS-2 attack vector of the vulnerability.`,
							},
							"vector_v3": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: `Indicates the CVSS-3 attack vector of the vulnerability.`,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceSwrEnterpriseInstanceArtifactVulnerabilitiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	// nolint:revive
	listHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}/artifacts/{reference}/vulnerabilities"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{namespace_name}", d.Get("namespace_name").(string))
	listPath = strings.ReplaceAll(listPath, "{repository_name}", url.PathEscape(strings.ReplaceAll(d.Get("repository_name").(string), "/", "%2F")))
	listPath = strings.ReplaceAll(listPath, "{reference}", d.Get("reference").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error querying SWR artifact vulnerabilities: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.Errorf("error flattening SWR artifact vulnerabilities response: %s", err)
	}

	reports := listRespBody.(map[string]interface{})
	results := make([]map[string]interface{}, 0, len(reports))

	for reportType, report := range reports {
		results = append(results, map[string]interface{}{
			"type":    reportType,
			"content": flattenArtifactVulnerabilitiesContent(report),
		})
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("reports", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenArtifactVulnerabilitiesContent(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}
	rst := map[string]interface{}{
		"generated_at": utils.PathSearch("generated_at", params, nil),
		"severity":     utils.PathSearch("severity", params, nil),
		"scanner":      flattenArtifactVulnerabilitiesContentScanner(utils.PathSearch("scanner", params, nil)),
		"vulnerabilities": flattenArtifactVulnerabilitiesContentVulnerabilities(
			utils.PathSearch("vulnerabilities", params, make([]interface{}, 0)).([]interface{})),
	}

	return []map[string]interface{}{rst}
}

func flattenArtifactVulnerabilitiesContentScanner(params interface{}) []map[string]interface{} {
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

func flattenArtifactVulnerabilitiesContentVulnerabilities(paramsList []interface{}) []map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"id":               utils.PathSearch("id", params, nil),
			"package":          utils.PathSearch("package", params, nil),
			"version":          utils.PathSearch("version", params, nil),
			"fix_version":      utils.PathSearch("fix_version", params, nil),
			"severity":         utils.PathSearch("severity", params, nil),
			"description":      utils.PathSearch("description", params, nil),
			"links":            utils.PathSearch("links", params, nil),
			"artifact_digests": utils.PathSearch("artifact_digests", params, nil),
			"cwe_ids":          utils.PathSearch("cwe_ids", params, nil),
			"preferred_cvss": flattenArtifactVulnerabilitiesContentVulnerabilitiesPreferredCvss(
				utils.PathSearch("preferred_cvss", params, nil)),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenArtifactVulnerabilitiesContentVulnerabilitiesPreferredCvss(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}
	rst := map[string]interface{}{
		"score_v2":  utils.PathSearch("score_v2", params, nil),
		"score_v3":  utils.PathSearch("score_v3", params, nil),
		"vector_v2": utils.PathSearch("vector_v2", params, nil),
		"vector_v3": utils.PathSearch("vector_v3", params, nil),
	}

	return []map[string]interface{}{rst}
}
