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

var configurationNonUpdatableParams = []string{"config_group_id", "name"}

// @API ServiceStage POST /v3/{project_id}/cas/configs
// @API ServiceStage GET /v3/{project_id}/cas/configs/{config_id}
// @API ServiceStage PUT /v3/{project_id}/cas/configs/{config_id}
// @API ServiceStage DELETE /v3/{project_id}/cas/configs/{config_id}
func ResourceV3Configuration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3ConfigurationCreate,
		ReadContext:   resourceV3ConfigurationRead,
		UpdateContext: resourceV3ConfigurationUpdate,
		DeleteContext: resourceV3ConfigurationDelete,

		CustomizeDiff: config.FlexibleForceNew(configurationNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"config_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the configuration group to which the configuration file belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the configuration file.`,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The content of the configuration file.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the configuration file.`,
			},
			"sensitive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable data encryption.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the configuration file.`,
			},
			"components": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the environment.`,
						},
						"application_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the application.`,
						},
						"component_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the component.`,
						},
						"component_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the component.`,
						},
					},
				},
				Description: `The list of the components associated with the configuration file.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the configuration file.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the configuration file.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the configuration file, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the configuration file, in RFC3339 format.`,
			},
		},
	}
}

func buildV3ConfigurationCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"config_group_id": d.Get("config_group_id"),
		"name":            d.Get("name"),
		"content":         d.Get("content"),
		"type":            d.Get("type"),
		"sensitive":       d.Get("sensitive"),
		"description":     utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceV3ConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v3/{project_id}/cas/configs"
	)
	client, err := cfg.NewServiceClient("servicestage", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildV3ConfigurationCreateBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating configuration file: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	configurationId := utils.PathSearch("id", respBody, "").(string)
	if configurationId == "" {
		return diag.Errorf("unable to find the configuration file ID from the API response")
	}
	d.SetId(configurationId)

	return resourceV3ConfigurationRead(ctx, d, meta)
}

// GetV3ConfigurationFile is a method used to get configuration file detail by its ID.
func GetV3ConfigurationFile(client *golangsdk.ServiceClient, configurationId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/cas/configs/{config_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{config_id}", configurationId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceV3ConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		configurationId = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	respBody, err := GetV3ConfigurationFile(client, configurationId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving configuration file (%s)", configurationId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("config_group_id", utils.PathSearch("config_group_id", respBody, nil)),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("content", utils.PathSearch("content", respBody, nil)),
		d.Set("type", utils.PathSearch("type", respBody, nil)),
		d.Set("sensitive", utils.PathSearch("sensitive", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("components", flattenConfigComponents(utils.PathSearch("components", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("version", utils.PathSearch("version", respBody, nil)),
		d.Set("creator", utils.PathSearch("creator", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", respBody,
			float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", respBody,
			float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConfigComponents(components []interface{}) []interface{} {
	if len(components) == 0 {
		return nil
	}

	rest := make([]interface{}, len(components))
	for i, v := range components {
		rest[i] = map[string]interface{}{
			"environment_id": utils.PathSearch("environment_id", v, nil),
			"application_id": utils.PathSearch("application_id", v, nil),
			"component_id":   utils.PathSearch("component_id", v, nil),
			"component_name": utils.PathSearch("component_name", v, nil),
		}
	}
	return rest
}

func buildV3ConfigurationUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"config_group_id": d.Get("config_group_id"),
		"name":            d.Get("name"),
		"content":         d.Get("content"),
		"type":            d.Get("type"),
		// The `description` must be specified as an empty string to be changed to empty.
		"description": d.Get("description"),
		"sensitive":   d.Get("sensitive"),
	}
}

func resourceV3ConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		httpUrl         = "v3/{project_id}/cas/configs/{config_id}"
		configurationId = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	if d.HasChanges("content", "type", "sensitive", "description") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{config_id}", configurationId)
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf8",
			},
			JSONBody: buildV3ConfigurationUpdateBodyParams(d),
		}

		_, err = client.Request("PUT", updatePath, &opt)
		if err != nil {
			return diag.Errorf("error updating configuration file (%s): %s", configurationId, err)
		}
	}

	return resourceV3ConfigurationRead(ctx, d, meta)
}

func resourceV3ConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		httpUrl         = "v3/{project_id}/cas/configs/{config_id}"
		configurationId = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{config_id}", configurationId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	// The delete API always returns 200 status code whether config group is exist.
	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting configuration file (%s)", configurationId))
	}
	return nil
}
