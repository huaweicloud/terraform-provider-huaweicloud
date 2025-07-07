package servicestage

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

var v3AppNotFoundCodes = []string{
	"SVCSTG.00100401",
}

// @API ServiceStage POST /v3/{project_id}/cas/applications
// @API ServiceStage GET /v3/{project_id}/cas/applications/{application_id}
// @API ServiceStage PUT /v3/{project_id}/cas/applications/{application_id}
// @API ServiceStage DELETE /v3/{project_id}/cas/applications/{application_id}
func ResourceV3Application() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3ApplicationCreate,
		ReadContext:   resourceV3ApplicationRead,
		UpdateContext: resourceV3ApplicationUpdate,
		DeleteContext: resourceV3ApplicationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the application is located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the application.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the application.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The the enterprise project ID to which the application belongs.`,
			},
			"tags": common.TagsSchema(
				`he key/value pairs to associate with the application that used to filter resource.`,
			),
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator name of the application.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the application, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the application, in RFC3339 format.`,
			},
		},
	}
}

func buildV3ApplicationCreateBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":                  d.Get("name"),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project_id": cfg.GetEnterpriseProjectID(d),
		"labels":                utils.ValueIgnoreEmpty(utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{}))),
	}
}

func resourceV3ApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/cas/applications"
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV3ApplicationCreateBodyParams(cfg, d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating application: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	appId := utils.PathSearch("id", respBody, "").(string)
	if appId == "" {
		return diag.Errorf("failed to find the application ID from the API response")
	}
	d.SetId(appId)

	return resourceV3ApplicationRead(ctx, d, meta)
}

func QueryV3Application(client *golangsdk.ServiceClient, appId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/cas/applications/{application_id}"

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{application_id}", appId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", queryPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceV3ApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		appId  = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	respBody, err := QueryV3Application(client, appId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected401ErrInto404Err(err, "error_code", v3AppNotFoundCodes...),
			fmt.Sprintf("error getting application (%s)", appId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", respBody, nil)),
		d.Set("creator", utils.PathSearch("creator", respBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("labels", respBody, make([]interface{}, 0)))),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", respBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildV3ApplicationUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"labels":      utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{}), true),
	}
}

func resourceV3ApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/cas/applications/{application_id}"
		appId   = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	if d.HasChangeExcept("enterprise_project_id") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{application_id}", appId)

		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: buildV3ApplicationUpdateBodyParams(d),
		}

		_, err = client.Request("PUT", updatePath, &opt)
		if err != nil {
			return diag.Errorf("error updating application (%s): %s", appId, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   appId,
			ResourceType: "servicestage-application",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceV3ApplicationRead(ctx, d, meta)
}

func resourceV3ApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/cas/applications/{application_id}"
		appId   = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{application_id}", appId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes: []int{204},
	}

	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected401ErrInto404Err(err, "error_code", v3AppNotFoundCodes...),
			fmt.Sprintf("error deleting application (%s)", appId))
	}
	return nil
}
