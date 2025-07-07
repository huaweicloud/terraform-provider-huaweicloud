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

var v3EnvNotFoundCodes = []string{
	"SVCSTG.00100401",
}

// @API ServiceStage POST /v3/{project_id}/cas/environments
// @API ServiceStage GET /v3/{project_id}/cas/environments/{environment_id}
// @API ServiceStage PUT /v3/{project_id}/cas/environments/{environment_id}
// @API ServiceStage DELETE /v3/{project_id}/cas/environments/{environment_id}
func ResourceV3Environment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3EnvironmentCreate,
		ReadContext:   resourceV3EnvironmentRead,
		UpdateContext: resourceV3EnvironmentUpdate,
		DeleteContext: resourceV3EnvironmentDelete,

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
				Description: `The region where the environment is located.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The VPC ID to which the environment belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the environment.`,
			},
			"deploy_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The deploy mode of the environment.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the environment.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The the enterprise project ID to which the environment belongs.`,
			},
			"tags": common.TagsSchema(
				`The key/value pairs to associate with the environment that used to filter resource.`,
			),
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator name of the environment.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the environment, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the environment, in RFC3339 format.`,
			},
		},
	}
}

func buildV3EnvironmentCreateBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"vpc_id":                d.Get("vpc_id"),
		"name":                  d.Get("name"),
		"deploy_mode":           utils.ValueIgnoreEmpty(d.Get("deploy_mode")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project_id": cfg.GetEnterpriseProjectID(d),
		"labels":                utils.ValueIgnoreEmpty(utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{}))),
	}
}

func resourceV3EnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/cas/environments"
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
		JSONBody: utils.RemoveNil(buildV3EnvironmentCreateBodyParams(cfg, d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating environment: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	envId := utils.PathSearch("id", respBody, "").(string)
	if envId == "" {
		return diag.Errorf("unable to find the environment ID from the API response")
	}
	d.SetId(envId)

	return resourceV3EnvironmentRead(ctx, d, meta)
}

func QueryV3Environment(client *golangsdk.ServiceClient, envId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/cas/environments/{environment_id}"

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{environment_id}", envId)

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

func resourceV3EnvironmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		envId  = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	respBody, err := QueryV3Environment(client, envId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected401ErrInto404Err(err, "error_code", v3EnvNotFoundCodes...),
			fmt.Sprintf("error getting environment (%s)", envId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_id", utils.PathSearch("vpc_id", respBody, nil)),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("deploy_mode", utils.PathSearch("deploy_mode", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", respBody, nil)),
		d.Set("creator", utils.PathSearch("creator", respBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("labels", respBody, make([]interface{}, 0)))),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", respBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildV3EnvironmentUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"labels":      utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{}), true),
	}
}

func resourceV3EnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/cas/environments/{environment_id}"
		envId   = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	if d.HasChangeExcept("enterprise_project_id") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{environment_id}", envId)

		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: buildV3EnvironmentUpdateBodyParams(d),
		}

		_, err = client.Request("PUT", updatePath, &opt)
		if err != nil {
			return diag.Errorf("error updating environment (%s): %s", envId, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   envId,
			ResourceType: "servicestage-environment",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceV3EnvironmentRead(ctx, d, meta)
}

func resourceV3EnvironmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/cas/environments/{environment_id}"
		envId   = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{environment_id}", envId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes: []int{204},
	}

	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected401ErrInto404Err(err, "error_code", v3EnvNotFoundCodes...),
			fmt.Sprintf("error deleting environment (%s)", envId))
	}
	return nil
}
