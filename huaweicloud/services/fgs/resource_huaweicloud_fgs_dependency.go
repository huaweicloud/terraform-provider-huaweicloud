package fgs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph POST /v2/{project_id}/fgs/dependencies
// @API FunctionGraph GET /v2/{project_id}/fgs/dependencies/{depend_id}
// @API FunctionGraph PUT /v2/{project_id}/fgs/dependencies/{depend_id}
// @API FunctionGraph DELETE /v2/{project_id}/fgs/dependencies/{depend_id}
func ResourceDependency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDependencyCreate,
		ReadContext:   resourceDependencyRead,
		UpdateContext: resourceDependencyUpdate,
		DeleteContext: resourceDependencyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the dependency package is located.`,
			},
			"runtime": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The runtime of the dependency package.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the dependency package.`,
			},
			"link": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The OBS storage URL of the dependency package.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the dependency package.`,
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The etag of the dependency package.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The owner name of the dependency package.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The capacity of the dependency package.`,
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The version of the dependency package.`,
			},
		},
	}
}

func buildCreateDependencyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name").(string),
		"runtime":     d.Get("runtime").(string),
		"description": d.Get("description").(string),
		"depend_type": "obs",
		"depend_link": d.Get("link").(string),
	}
}

func resourceDependencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/fgs/dependencies"
	)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDependencyBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph custom dependency: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dependencyId := utils.PathSearch("id", respBody, "").(string)
	if dependencyId == "" {
		return diag.Errorf("unable to find the dependency ID from the API response")
	}
	d.SetId(dependencyId)

	return resourceDependencyRead(ctx, d, meta)
}

func GetDependencyById(client *golangsdk.ServiceClient, dependencyId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/fgs/dependencies/{depend_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{depend_id}", dependencyId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceDependencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		dependencyId = d.Id()
	)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	respBody, err := GetDependencyById(client, dependencyId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying custom dependency (%s)", dependencyId))
	}

	mErr := multierror.Append(
		d.Set("runtime", utils.PathSearch("runtime", respBody, nil)),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("link", utils.PathSearch("link", respBody, nil)),
		d.Set("etag", utils.PathSearch("etag", respBody, nil)),
		d.Set("size", utils.PathSearch("size", respBody, nil)),
		d.Set("owner", utils.PathSearch("owner", respBody, nil)),
		d.Set("version", utils.PathSearch("version", respBody, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting resource fields of custom dependency (%s): %s", d.Id(), err)
	}
	return nil
}

func buildUpdateDependencyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name").(string),
		"runtime":     d.Get("runtime").(string),
		"description": d.Get("description").(string),
		"depend_type": "obs",
		"depend_link": d.Get("link").(string),
	}
}

func resourceDependencyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v2/{project_id}/fgs/dependencies/{depend_id}"
		dependencyId = d.Id()
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{depend_id}", dependencyId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildUpdateDependencyBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating custom dependency: %s", err)
	}
	return resourceDependencyRead(ctx, d, meta)
}

func resourceDependencyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v2/{project_id}/fgs/dependencies/{depend_id}"
		dependencyId = d.Id()
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{depend_id}", dependencyId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting custom dependency (%s)", dependencyId))
	}
	return nil
}
