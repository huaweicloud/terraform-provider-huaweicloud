package servicestage

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v3AppConfigurationNonUpdatableParams = []string{
	"environment_id",
	"application_id",
}

// @API ServiceStage PUT /v3/{project_id}/cas/applications/{application_id}/configuration
// @API ServiceStage GET /v3/{project_id}/cas/applications/{application_id}/configuration
// @API ServiceStage DELETE /v3/{project_id}/cas/applications/{application_id}/configuration
func ResourceV3ApplicationConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3ApplicationConfigurationCreate,
		ReadContext:   resourceV3ApplicationConfigurationRead,
		UpdateContext: resourceV3ApplicationConfigurationUpdate,
		DeleteContext: resourceV3ApplicationConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV3ApplicationConfigurationImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(v3AppConfigurationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the environment and application are located.`,
			},
			// Required parameters.
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the environment to which the configuration applies.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application to which the configuration belongs.`,
			},
			"configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"env": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The name of the environment variable.`,
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The value of the environment variable.`,
									},
								},
							},
							Description: `The list of the environment variables.`,
						},
						"assign_strategy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether the effective strategy is the continuously effective. `,
						},
					},
				},
				Description: `The configuration of the application.`,
			},
			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildV3ApplicationConfigurationEnvs(envs []interface{}) []map[string]interface{} {
	if len(envs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(envs))
	for _, env := range envs {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", env, nil),
			"value": utils.PathSearch("value", env, nil),
		})
	}
	return result
}

func updateV3ApplicationConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	var (
		httpUrl       = "v3/{project_id}/cas/applications/{application_id}/configuration"
		environmentId = d.Get("environment_id").(string)
		applicationId = d.Get("application_id").(string)
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{application_id}", applicationId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"environment_id": environmentId,
			"configuration": map[string]interface{}{
				"env":             buildV3ApplicationConfigurationEnvs(d.Get("configuration.0.env").(*schema.Set).List()),
				"assign_strategy": d.Get("configuration.0.assign_strategy"),
			},
		}),
	}

	_, err := client.Request("PUT", updatePath, &opt)
	if err != nil {
		return "", fmt.Errorf("error modifying application configuration: %s", err)
	}
	return fmt.Sprintf("%s/%s", environmentId, applicationId), nil
}

func resourceV3ApplicationConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	resourceId, err := updateV3ApplicationConfiguration(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)

	return resourceV3ApplicationConfigurationRead(ctx, d, meta)
}

func GetV3ApplicationConfiguration(client *golangsdk.ServiceClient, environmentId, applicationId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/cas/applications/{application_id}/configuration"

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{application_id}", applicationId)

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

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	configuration := utils.PathSearch(fmt.Sprintf("configuration[?environment_id=='%s']|[0].configuration",
		environmentId), respBody, nil)
	if configuration == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/{project_id}/cas/applications/{application_id}/configuration",
				RequestId: "NONE",
				Body: []byte(fmt.Sprintf("the configurations is empty in the application (%s) for the environment (%s)",
					applicationId, environmentId)),
			},
		}
	}
	return configuration, nil
}

func flattenV3ApplicationConfigurationEnvs(envs []interface{}) []map[string]interface{} {
	if len(envs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(envs))
	for _, env := range envs {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", env, nil),
			"value": utils.PathSearch("value", env, nil),
		})
	}

	return result
}

func flattenV3ApplicationConfiguration(configuration interface{}) []map[string]interface{} {
	if configuration == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"env":             flattenV3ApplicationConfigurationEnvs(utils.PathSearch("env", configuration, make([]interface{}, 0)).([]interface{})),
			"assign_strategy": utils.PathSearch("assign_strategy", configuration, nil),
		},
	}
}

func resourceV3ApplicationConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		environmentId = d.Get("environment_id").(string)
		applicationId = d.Get("application_id").(string)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	configuration, err := GetV3ApplicationConfiguration(client, environmentId, applicationId)
	if err != nil {
		// The API returns 401 error when the application not found, but the 404 error returned if the configuration not
		// found or the environment not found.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", v3AppNotFoundCodes...),
			fmt.Sprintf("error getting application (%s) configuration", applicationId),
		)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("configuration", flattenV3ApplicationConfiguration(configuration)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV3ApplicationConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	_, err = updateV3ApplicationConfiguration(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceV3ApplicationConfigurationRead(ctx, d, meta)
}

func resourceV3ApplicationConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v3/{project_id}/cas/applications/{application_id}/configuration?environment_id={environment_id}"
		environmentId = d.Get("environment_id").(string)
		applicationId = d.Get("application_id").(string)
	)

	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{application_id}", applicationId)
	deletePath = strings.ReplaceAll(deletePath, "{environment_id}", environmentId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		// The API returns 401 error when the application not found, but the 204 status code returned if the
		// configuration not found or the environment not found.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", v3AppNotFoundCodes...),
			fmt.Sprintf("error deleting application (%s) configuration", applicationId),
		)
	}
	return nil
}

func resourceV3ApplicationConfigurationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<environment_id>/<application_id>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("environment_id", parts[0]),
		d.Set("application_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
