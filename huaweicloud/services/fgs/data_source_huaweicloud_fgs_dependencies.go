package fgs

import (
	"context"
	"fmt"
	"log"
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

// @API FunctionGraph GET /v2/{project_id}/fgs/dependencies
// DataSourceDependencies provides some parameters to filter dependent packages on the server.
func DataSourceDependencies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDependenciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the dependency packages are located.`,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"public", "private",
				}, false),
				Description: `The type of the dependency package.`,
			},
			"runtime": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The runtime of the dependency package.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the dependency package.`,
			},
			"is_versions_query_allowed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to query the versions of each dependency package.`,
			},
			"packages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dependency package.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the dependency package.`,
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The owner of the dependency package.`,
						},
						"link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OBS bucket path where the dependency package is located (FunctionGraph serivce side).`,
						},
						"etag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unique ID of the dependency package.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The size of the dependency package.`,
						},
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The file name of the stored dependency package.`,
						},
						"runtime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The runtime of the dependency package.`,
						},
						"versions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the dependency package version.",
									},
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The dependency package version.",
									},
								},
							},
							Description: `The list of the versions for the dependency package.`,
						},
					},
				},
				Description: `All dependency packages that match the filter parameters.`,
			},
		},
	}
}

func buildDependenciesQueryParams(d *schema.ResourceData) string {
	result := ""

	if pkgType, ok := d.GetOk("type"); ok {
		result = fmt.Sprintf("%s&dependency_type=%v", result, pkgType)
	}
	if pkgRuntime, ok := d.GetOk("runtime"); ok {
		result = fmt.Sprintf("%s&runtime=%v", result, pkgRuntime)
	}
	if pkgName, ok := d.GetOk("name"); ok {
		result = fmt.Sprintf("%s&name=%v", result, pkgName)
	}

	return result
}

func getDependencies(client *golangsdk.ServiceClient, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/dependencies?maxitems=100"
		marker  float64
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	if len(queryParams) > 0 {
		listPath += queryParams[0]
	}

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithMarker := fmt.Sprintf("%s&marker=%v", listPath, marker)
		requestResp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, fmt.Errorf("error querying dependency packages: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		dependencies := utils.PathSearch("dependencies", respBody, make([]interface{}, 0)).([]interface{})
		if len(dependencies) < 1 {
			break
		}
		result = append(result, dependencies...)
		// In this API, marker has the same meaning as offset.
		nextMarker := utils.PathSearch("next_marker", respBody, float64(0)).(float64)
		if nextMarker == marker || nextMarker == 0 {
			// Make sure the next marker value is correct, not the previous marker or zero (in the last page).
			break
		}
		marker = nextMarker
	}

	return result, nil
}

func dataSourceDependenciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	dependencies, err := getDependencies(client, buildDependenciesQueryParams(d))
	if err != nil {
		return diag.Errorf("error retrieving dependency packages: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		d.Set("packages", flattenDependencies(client, dependencies, d.Get("is_versions_query_allowed").(bool))),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving data source fields of the dependency packages: %s", err)
	}
	return nil
}

func flattenDependencies(client *golangsdk.ServiceClient, dependencies []interface{},
	isVersionsQueryAllowed bool) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(dependencies))

	for _, dependency := range dependencies {
		dependencyId := utils.PathSearch("id", dependency, "").(string)
		elem := map[string]interface{}{
			"id":        dependencyId,
			"name":      utils.PathSearch("name", dependency, nil),
			"owner":     utils.PathSearch("owner", dependency, nil),
			"link":      utils.PathSearch("link", dependency, nil),
			"etag":      utils.PathSearch("etag", dependency, nil),
			"size":      utils.PathSearch("size", dependency, nil),
			"file_name": utils.PathSearch("file_name", dependency, nil),
			"runtime":   utils.PathSearch("runtime", dependency, nil),
		}

		if isVersionsQueryAllowed && dependencyId != "" {
			dependencyVersions, err := getDependencyVersions(client, dependencyId)
			if err != nil {
				log.Printf("error retrieving versions under specified dependency package (%s): %s", dependencyId, err)
			} else {
				elem["versions"] = flattenDependencyVersionInfos(dependencyVersions)
			}
		}
		result = append(result, elem)
	}

	return result
}

func flattenDependencyVersionInfos(dependencyVersions []interface{}) []map[string]interface{} {
	if len(dependencyVersions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(dependencyVersions))
	for _, dependencyVersion := range dependencyVersions {
		result = append(result, map[string]interface{}{
			"id":      utils.PathSearch("id", dependencyVersion, nil),
			"version": int(utils.PathSearch("version", dependencyVersion, float64(0)).(float64)),
		})
	}

	return result
}
