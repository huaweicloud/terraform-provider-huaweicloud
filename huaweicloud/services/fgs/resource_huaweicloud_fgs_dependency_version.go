package fgs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/fgs/v2/dependencies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph POST /v2/{project_id}/fgs/dependencies/version
// @API FunctionGraph GET /v2/{project_id}/fgs/dependencies
// @API FunctionGraph GET /v2/{project_id}/fgs/dependencies/{depend_id}/version
// @API FunctionGraph GET /v2/{project_id}/fgs/dependencies/{depend_id}/version/{version}
// @API FunctionGraph DELETE /v2/{project_id}/fgs/dependencies/{depend_id}/version/{version}
func ResourceDependencyVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDependencyVersionCreate,
		ReadContext:   resourceDependencyVersionRead,
		DeleteContext: resourceDependencyVersionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDependencyVersionImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the custom dependency version is located.",
			},
			"runtime": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The runtime of the custom dependency package version.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the custom dependency package to which the version belongs.",
			},
			"link": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The OBS bucket path where the dependency package is located.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The description of the custom dependency version.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The dependency package version.",
			},
			"version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the dependency package version.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The dependency owner, public indicates a public dependency.",
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique ID of the dependency.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The dependency size, in bytes.",
			},
			"dependency_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the dependency package corresponding to the version.",
			},
		},
	}
}

func buildDependencyVersionOpts(d *schema.ResourceData) dependencies.DependVersionOpts {
	// Since the ZIP file upload is limited in size and requires encoding, only the OBS type is supported.
	// The ZIP file uploading can also be achieved by uploading OBS objects and is more secure.
	return dependencies.DependVersionOpts{
		Name:        d.Get("name").(string),
		Runtime:     d.Get("runtime").(string),
		Description: utils.String(d.Get("description").(string)),
		Type:        "obs",
		Link:        d.Get("link").(string),
	}
}

func resourceDependencyVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	resp, err := dependencies.CreateVersion(client, buildDependencyVersionOpts(d))
	if err != nil {
		return diag.Errorf("error creating custom dependency version: %s", err)
	}
	// Using depend ID and version number as the resource ID.
	d.SetId(fmt.Sprintf("%s/%d", resp.DepId, resp.Version))

	return resourceDependencyVersionRead(ctx, d, meta)
}

func ParseDependVersionResourceId(resourceId string) (dependId, versionInfo string, err error) {
	parts := strings.Split(resourceId, "/")
	if len(parts) < 2 {
		err = fmt.Errorf("invalid ID format for dependency version resource, it must contain two parts: "+
			"dependency package information and version information, e.g. '<dependency name>/<version number>'. "+
			"but the ID that you provided does not meet this requirement '%s'", resourceId)
		return
	}
	dependId = parts[0]
	versionInfo = parts[1]
	return
}

func resourceDependencyVersionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	dependId, version, err := ParseDependVersionResourceId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := dependencies.GetVersion(client, dependId, version)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "FunctionGraph dependency version")
	}

	// FunctionGraph will store the compressed package content pointed to by the link into the new storage bucket that
	// provided by FunctionGraph and return a new link value.
	// If the ReadContext is set this value according to the query result, ForceNew behavior will be triggered the next
	// time it is applied.
	mErr := multierror.Append(
		d.Set("runtime", resp.Runtime),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("etag", resp.Etag),
		d.Set("size", resp.Size),
		d.Set("owner", resp.Owner),
		d.Set("version", resp.Version),
		d.Set("version_id", resp.ID),
		d.Set("dependency_id", resp.DepId),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting resource fields of custom dependency version (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceDependencyVersionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	dependId, version, err := ParseDependVersionResourceId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	err = dependencies.DeleteVersion(client, dependId, version)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting custom dependency version")
	}
	return nil
}

// getSpecifiedDependencyVersion is a method that queries the corresponding dependency version based on the entered ID.
// The entered ID can be in the following formats:
// + <depend_id>/<version> (Standard resource ID format)
// + <depend_id>/<version_id>
// + <depend_name>/<version> (All information that can be found through the console)
// + <depend_name>/<version_id>
func getSpecifiedDependencyVersion(client *golangsdk.ServiceClient, resourceId string) (*dependencies.DependencyVersion, error) {
	dependInfo, versionInfo, err := ParseDependVersionResourceId(resourceId)
	if err != nil {
		return nil, err
	}

	// If the input dependency package information part is not in UUID format, perform a query to obtain the
	// corresponding ID.
	if !utils.IsUUID(dependInfo) {
		opts := dependencies.ListOpts{
			Name: dependInfo,
		}
		allPages, err := dependencies.List(client, opts).AllPages()
		if err != nil {
			return nil, err
		}
		listResp, _ := dependencies.ExtractDependencies(allPages)
		if len(listResp.Dependencies) < 1 {
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Body: []byte(fmt.Sprintf("unable to find the dependency package using its name: %s", dependInfo)),
				},
			}
		}
		// Make sure the dependInfo content is the dependency ID.
		dependInfo = listResp.Dependencies[0].ID
	}

	// If the input dependency version information part is in UUID format, perform a query to obtain the specified
	// version using its ID.
	if utils.IsUUID(versionInfo) {
		opts := dependencies.ListVersionsOpts{
			DependId: dependInfo,
		}
		listResp, err := dependencies.ListVersions(client, opts)
		if err != nil {
			return nil, err
		}
		for _, dependVersion := range listResp {
			if dependVersion.ID == versionInfo {
				return &dependVersion, nil
			}
		}
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("unable to find the dependency package using its ID: %s", versionInfo)),
			},
		}
	}
	return dependencies.GetVersion(client, dependInfo, versionInfo)
}

func resourceDependencyVersionImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	// Query the corresponding dependency version based on the user's import ID.
	resp, err := getSpecifiedDependencyVersion(client, d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	d.SetId(fmt.Sprintf("%s/%d", resp.DepId, resp.Version))

	return []*schema.ResourceData{d}, nil
}
