package fgs

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph GET /v2/{project_id}/fgs/dependencies/{depend_id}/version
func DataSourceDependencieVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDependencieVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the dependency package and the versions are located.`,
			},
			"dependency_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dependency package to which the versions belong.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the dependency package version.`,
			},
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The version of the dependency package.`,
			},
			"runtime": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The runtime of the dependency package version.`,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dependency package version.`,
						},
						"version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The dependency package version.`,
						},
						"dependency_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dependency package corresponding to the version.`,
						},
						"dependency_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the dependency package corresponding to the version.`,
						},
						"runtime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The runtime of the dependency package version.`,
						},
						"link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OBS bucket path where the dependency package version is located.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The size of the ZIP file used by the dependency package version, in bytes.`,
						},
						"etag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unique ID of the dependency.`,
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The dependency owner, public indicates a public dependency.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the dependency package version.`,
						},
					},
				},
				Description: `All dependency package versions that match the filter parameters.`,
			},
		},
	}
}

func dataSourceDependencieVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		dependId = d.Get("dependency_id").(string)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	dependVersions, err := getDependencyVersions(client, dependId)
	if err != nil {
		return diag.Errorf("error retrieving versions under specified dependency package (%s): %s", dependId, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("versions", filterVersions(flattenDependencyVersions(dependVersions), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterVersions(versions []map[string]interface{}, d *schema.ResourceData) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(versions))
	for _, version := range versions {
		if param, ok := d.GetOk("version_id"); ok && fmt.Sprint(param) != utils.PathSearch("id", version, "").(string) {
			continue
		}

		if param, ok := d.GetOk("version"); ok && param.(int) != utils.PathSearch("version", version, int(0)).(int) {
			continue
		}

		if param, ok := d.GetOk("runtime"); ok && fmt.Sprint(param) != utils.PathSearch("runtime", version, "").(string) {
			continue
		}

		rst = append(rst, version)
	}
	return rst
}

func flattenDependencyVersions(versions []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(versions))

	for _, version := range versions {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", version, ""),
			"version":         int(utils.PathSearch("version", version, float64(0)).(float64)),
			"dependency_id":   utils.PathSearch("dep_id", version, nil),
			"dependency_name": utils.PathSearch("name", version, nil),
			"runtime":         utils.PathSearch("runtime", version, ""),
			"link":            utils.PathSearch("link", version, nil),
			"size":            utils.PathSearch("size", version, nil),
			"etag":            utils.PathSearch("etag", version, nil),
			"owner":           utils.PathSearch("owner", version, nil),
			"description":     utils.PathSearch("description", version, nil),
		})
	}

	return result
}
