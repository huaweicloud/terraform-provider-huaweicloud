package fgs

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/fgs/v2/dependencies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API FunctionGraph GET /v2/{project_id}/fgs/dependencies
// DataSourceFunctionGraphDependencies provides some parameters to filter dependent packages on the server.
func DataSourceFunctionGraphDependencies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFunctionGraphDependenciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"public", "private",
				}, false),
			},
			"runtime": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"packages": {
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
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"etag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime": {
							Type:     schema.TypeString,
							Computed: true,
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
						},
					},
				},
			},
		},
	}
}

func dataSourceFunctionGraphDependenciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.FgsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}
	// Limit and Marker use default values.
	listOpts := dependencies.ListOpts{
		Name:           d.Get("name").(string),
		Runtime:        d.Get("runtime").(string),
		DependencyType: d.Get("type").(string),
	}

	allPages, err := dependencies.List(client, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("error retrieving dependent packages: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	resp, _ := dependencies.ExtractDependencies(allPages)
	packages := flatFunctionGraphDependencies(client, resp.Dependencies)
	mErr := multierror.Append(
		d.Set("packages", packages),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting packages of FunctionGraph dependencies: %s", err)
	}
	return nil
}

func flatFunctionGraphDependencies(client *golangsdk.ServiceClient, pkgs []dependencies.Dependency) []map[string]interface{} {
	packages := make([]map[string]interface{}, len(pkgs))

	names := schema.NewSet(schema.HashString, nil)
	for i, pkg := range pkgs {
		names.Add(pkg.Name)
		packages[i] = map[string]interface{}{
			"id":        pkg.ID,
			"name":      pkg.Name,
			"owner":     pkg.Owner,
			"link":      pkg.Link,
			"etag":      pkg.Etag,
			"size":      pkg.Size,
			"file_name": pkg.FileName,
			"runtime":   pkg.Runtime,
			"versions":  flattenPckVersions(client, pkg.ID),
		}
	}

	return packages
}

func flattenPckVersions(client *golangsdk.ServiceClient, dependencyId string) []map[string]interface{} {
	listOpts := dependencies.ListVersionsOpts{
		DependId: dependencyId,
	}
	dependencyVersions, err := dependencies.ListVersions(client, listOpts)
	if err != nil {
		log.Printf("error retrieving versions under specified dependency package (%s): %s", dependencyId, err)
		return nil
	}

	result := make([]map[string]interface{}, len(dependencyVersions))
	for i, version := range dependencyVersions {
		result[i] = map[string]interface{}{
			"id":      version.ID,
			"version": version.Version,
		}
	}

	return result
}
