package cae

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

var appResourceNotFoundCodes = []string{
	"CAE.01500003", // Application is not exist.
}

// @API CAE POST /v1/{project_id}/cae/applications
// @API CAE GET /v1/{project_id}/cae/applications/{application_id}
// @API CAE DELETE /v1/{project_id}/cae/applications/{application_id}
func ResourceApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		DeleteContext: resourceApplicationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the application is located.",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the environment to which the application belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the application.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the enterprise project to which the application belongs.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the application, in RFC3339 format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the application, in RFC3339 format.",
			},
		},
	}
}

func buildCreateApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"api_version": "v1",
		"kind":        "Application",
		"metadata": map[string]interface{}{
			"name": d.Get("name"),
		},
	}
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/cae/applications"
		envId   = d.Get("environment_id").(string)
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(envId, cfg.GetEnterpriseProjectID(d)),
		JSONBody:         utils.RemoveNil(buildCreateApplicationBodyParams(d)),
	}
	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating application: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	appId := utils.PathSearch("metadata.id", respBody, "").(string)
	if appId == "" {
		return diag.Errorf("unable to find the application ID from the API response")
	}
	d.SetId(appId)

	return resourceApplicationRead(ctx, d, meta)
}

func GetApplicationById(client *golangsdk.ServiceClient, envId, appId, epsId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cae/applications/{application_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{application_id}", appId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(envId, epsId),
	}
	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("metadata", respBody, nil), nil
}

func resourceApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		appId  = d.Id()
		envId  = d.Get("environment_id").(string)
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	app, err := GetApplicationById(client, envId, appId, cfg.GetEnterpriseProjectID(d))
	if err != nil {
		// 500 error returned if the application is not exist.
		// 400 error returned if the related environment is not exist.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(
			common.ConvertExpected500ErrInto404Err(err, "error_code", appResourceNotFoundCodes...),
			"error_code",
			envResourceNotFoundCodes...), fmt.Sprintf("error querying application by its ID (%s)", appId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", app, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at",
			app, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("updated_at",
			app, "").(string))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/cae/applications/{application_id}"
		envId   = d.Get("environment_id").(string)
		appId   = d.Id()
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{application_id}", appId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(envId, cfg.GetEnterpriseProjectID(d)),
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// 400 error returned if the related environment is not exist.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", envResourceNotFoundCodes...),
			fmt.Sprintf("error deleting application (%s)", appId))
	}
	return nil
}

func resourceApplicationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	switch len(parts) {
	case 2:
		d.SetId(parts[1])
		return []*schema.ResourceData{d}, d.Set("environment_id", parts[0])
	case 3:
		d.SetId(parts[1])
		mErr := multierror.Append(
			d.Set("environment_id", parts[0]),
			d.Set("enterprise_project_id", parts[2]),
		)
		return []*schema.ResourceData{d}, mErr.ErrorOrNil()
	}

	return nil, fmt.Errorf("invalid format specified for import ID, want '<environment_id>/<id>' or "+
		"'<environment_id>/<id>/<enterprise_project_id>', but got '%s'",
		importedId)
}
