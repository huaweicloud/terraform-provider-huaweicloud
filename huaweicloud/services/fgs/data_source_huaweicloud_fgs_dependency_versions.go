package fgs

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/dependencies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API FunctionGraph GET /v2/{project_id}/fgs/dependencies/{depend_id}/version
func DataSourceDependencieVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDependencieVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dependency_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the dependency package to which the versions belong.",
			},
			"version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the dependency package version.",
			},
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the version of the dependency package.",
			},
			"runtime": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the runtime of the dependency package version.",
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
						"dependency_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the dependency package corresponding to the version.",
						},
						"dependency_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the dependency package corresponding to the version.",
						},
						"runtime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The runtime of the dependency package version.",
						},
						"link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The OBS bucket path where the dependency package version is located.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the ZIP file used by the dependency package version, in bytes.",
						},
						"etag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the dependency.",
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dependency owner, public indicates a public dependency.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the dependency package version.",
						},
					},
				},
			},
		},
	}
}

func dataSourceDependencieVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.FgsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	dependencyId := d.Get("dependency_id").(string)
	listOpts := dependencies.ListVersionsOpts{
		DependId: dependencyId,
	}
	dependencyVersions, err := dependencies.ListVersions(client, listOpts)
	if err != nil {
		return diag.Errorf("error retrieving versions under specified dependency package (%s): %s", dependencyId, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("versions", filterVersions(flattenVersions(dependencyVersions), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterVersions(versions []map[string]interface{}, d *schema.ResourceData) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(versions))
	for _, v := range versions {
		if param, ok := d.GetOk("version_id"); ok && fmt.Sprint(param) != v["id"] {
			continue
		}

		if param, ok := d.GetOk("version"); ok && param.(int) != v["version"] {
			continue
		}

		if param, ok := d.GetOk("runtime"); ok && fmt.Sprint(param) != v["runtime"] {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func flattenVersions(versions []dependencies.DependencyVersion) []map[string]interface{} {
	result := make([]map[string]interface{}, len(versions))
	for i, version := range versions {
		result[i] = map[string]interface{}{
			"id":              version.ID,
			"version":         version.Version,
			"dependency_id":   version.DepId,
			"dependency_name": version.Name,
			"runtime":         version.Runtime,
			"link":            version.Link,
			"size":            version.Size,
			"etag":            version.Etag,
			"owner":           version.Owner,
			"description":     version.Description,
		}
	}

	return result
}
