package fgs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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
				Description: `The region where the custom dependency version is located.`,
			},
			"runtime": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The runtime of the custom dependency package version.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the custom dependency package to which the version belongs.`,
			},
			"link": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The OBS bucket path where the dependency package is located.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The description of the custom dependency version.`,
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The dependency package version.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the dependency package version.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The dependency owner, public indicates a public dependency.`,
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique ID of the dependency.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The dependency size, in bytes.`,
			},
			"dependency_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the dependency package corresponding to the version.`,
			},
		},
	}
}

func buildCreateDependencyVersionBodyParams(d *schema.ResourceData) map[string]interface{} {
	// Since the ZIP file upload is limited in size and requires encoding, only the OBS type is supported.
	// The ZIP file uploading can also be achieved by uploading OBS objects and is more secure.
	return map[string]interface{}{
		// Required parameters.
		"name":        d.Get("name").(string),
		"runtime":     d.Get("runtime").(string),
		"depend_type": "obs",
		"depend_link": d.Get("link").(string),
		// Optional parameters.
		"description": d.Get("description").(string),
	}
}

func resourceDependencyVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/fgs/dependencies/version"
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateDependencyVersionBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating dependency package version: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dependId := utils.PathSearch("dep_id", respBody, "").(string)
	dependVersion := utils.PathSearch("version", respBody, float64(0)).(float64)
	if dependId == "" || dependVersion == 0 {
		return diag.Errorf("unable to find the dependency package version ID or version from the API response")
	}

	// Using depend ID and version number as the resource ID.
	d.SetId(fmt.Sprintf("%s/%v", dependId, dependVersion))

	return resourceDependencyVersionRead(ctx, d, meta)
}

func parseDependVersionResourceId(resourceId string) (dependId, versionInfo string) {
	parts := strings.Split(resourceId, "/")
	if len(parts) < 2 {
		log.Printf("[ERROR] invalid ID format for dependency package version resource, it must contain two parts: "+
			"dependency package information and version information, e.g. '<dependency name>/<version number>'. "+
			"but the ID that you provided does not meet this requirement '%s'", resourceId)
		return
	}
	dependId = parts[0]
	versionInfo = parts[1]
	return
}

func GetDependencyVersionById(client *golangsdk.ServiceClient, resourceId string) (interface{}, error) {
	var (
		httpUrl                 = "v2/{project_id}/fgs/dependencies/{depend_id}/version/{version}"
		dependId, dependVersion = parseDependVersionResourceId(resourceId)
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{depend_id}", dependId)
	getPath = strings.ReplaceAll(getPath, "{version}", dependVersion)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying dependency package version (%s): %s", dependVersion, err)
	}
	return utils.FlattenResponse(requestResp)
}

func resourceDependencyVersionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		resourceId = d.Id()
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	respBody, err := GetDependencyVersionById(client, resourceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "FunctionGraph dependency package version")
	}

	// FunctionGraph will store the compressed package content pointed to by the link into the new storage bucket that
	// provided by FunctionGraph and return a new link value.
	// If the ReadContext is set this value according to the query result, ForceNew behavior will be triggered the next
	// time it is applied.
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("runtime", utils.PathSearch("runtime", respBody, nil)),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("etag", utils.PathSearch("etag", respBody, nil)),
		d.Set("size", utils.PathSearch("size", respBody, nil)),
		d.Set("owner", utils.PathSearch("owner", respBody, nil)),
		d.Set("version", utils.PathSearch("version", respBody, nil)),
		d.Set("version_id", utils.PathSearch("id", respBody, nil)),
		d.Set("dependency_id", utils.PathSearch("dep_id", respBody, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting resource fields of custom dependency package version (%s): %s",
			resourceId, err)
	}

	return nil
}

func resourceDependencyVersionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                     = meta.(*config.Config)
		region                  = cfg.GetRegion(d)
		httpUrl                 = "v2/{project_id}/fgs/dependencies/{depend_id}/version/{version}"
		resourceId              = d.Id()
		dependId, dependVersion = parseDependVersionResourceId(resourceId)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{depend_id}", dependId)
	deletePath = strings.ReplaceAll(deletePath, "{version}", dependVersion)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting custom dependency package version")
	}
	return nil
}

func getDependencyVersions(client *golangsdk.ServiceClient, dependId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/dependencies/{depend_id}/version?maxitems=100"
		marker  float64
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{depend_id}", dependId)
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
			return nil, fmt.Errorf("error querying dependency package versions: %s", err)
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

// getSpecifiedDependencyVersion is a method that queries the corresponding dependency version based on the entered ID.
// The entered ID can be in the following formats:
// + <depend_id>/<version> (Standard resource ID format)
// + <depend_id>/<version_id>
// + <depend_name>/<version> (All information that can be found through the console)
// + <depend_name>/<version_id>
func refreshSpecifiedDependencyVersion(client *golangsdk.ServiceClient, resourceId string) (dependId, dependVersion string, err error) {
	dependId, dependVersion = parseDependVersionResourceId(resourceId)

	var result []interface{}
	// If the input dependency package information part is not in UUID format, perform a query to obtain the
	// corresponding ID.
	if !utils.IsUUID(dependId) {
		result, err = getDependencies(client)
		if err != nil {
			return dependId, dependVersion, err
		}
		dependId = utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].id", dependId), result, "").(string)
		if dependId == "" {
			return dependId, dependVersion, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Body: []byte(fmt.Sprintf("unable to find the dependency package using its name: %s", dependId)),
				},
			}
		}
	}

	// If the input dependency version information part is in UUID format, perform a query to obtain the specified
	// version using its ID.
	if utils.IsUUID(dependVersion) {
		result, err := getDependencyVersions(client, dependId)
		if err != nil {
			return dependId, dependVersion, err
		}
		dependVersion = fmt.Sprint(utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0].version", dependVersion),
			result, float64(0)).(float64))
		if dependVersion == "" {
			return dependId, dependVersion, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Body: []byte(fmt.Sprintf("unable to find the dependency package version using its ID: %s", dependVersion)),
				},
			}
		}
	}

	return dependId, dependVersion, nil
}

func resourceDependencyVersionImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}

	// Query the corresponding dependency version based on the user's import ID.
	dependId, dependVersion, err := refreshSpecifiedDependencyVersion(client, d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	d.SetId(fmt.Sprintf("%s/%s", dependId, dependVersion))

	return []*schema.ResourceData{d}, nil
}
