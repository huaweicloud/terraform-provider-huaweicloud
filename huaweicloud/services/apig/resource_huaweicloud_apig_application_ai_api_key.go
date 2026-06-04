package apig

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

var applicationAiApiKeyNonUpdatableParams = []string{
	"instance_id",
	"application_id",
	"value",
	"alias",
}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/ai-api-keys
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/ai-api-keys/{ai_api_key_id}
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/ai-api-keys/{ai_api_key_id}
func ResourceApplicationAiApiKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationAiApiKeyCreate,
		ReadContext:   resourceApplicationAiApiKeyRead,
		UpdateContext: resourceApplicationAiApiKeyUpdate,
		DeleteContext: resourceApplicationAiApiKeyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationAiApiKeyImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(applicationAiApiKeyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the AI API key is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the application and AI API key belong.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application to which the AI API key belongs.`,
			},

			// Optional parameters.
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The alias of the AI API key.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: `The value of the AI API key.`,
			},

			// Attributes.
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the AI API key, in RFC3339 format.`,
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

func buildApplicationAiApiKeyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"alias":      utils.ValueIgnoreEmpty(d.Get("alias")),
		"ai_api_key": utils.ValueIgnoreEmpty(d.Get("value")),
	}
}

func createApplicationAiApiKey(client *golangsdk.ServiceClient, instanceId, appId string,
	d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/ai-api-keys"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{app_id}", appId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildApplicationAiApiKeyBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceApplicationAiApiKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	respBody, err := createApplicationAiApiKey(client, instanceId, appId, d)
	if err != nil {
		return diag.Errorf("error creating AI API key: %s", err)
	}

	resourceId := utils.PathSearch("id", respBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("unable to find the AI API key ID from the API response")
	}
	d.SetId(resourceId)

	return resourceApplicationAiApiKeyRead(ctx, d, meta)
}

func GetApplicationAiApiKey(client *golangsdk.ServiceClient, instanceId, appId, aiApiKeyId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/ai-api-keys/{ai_api_key_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{app_id}", appId)
	getPath = strings.ReplaceAll(getPath, "{ai_api_key_id}", aiApiKeyId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceApplicationAiApiKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		aiApiKeyId = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	respBody, err := GetApplicationAiApiKey(client, instanceId, appId, aiApiKeyId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying AI API key (%s) from application (%s)", aiApiKeyId, appId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("value", utils.PathSearch("ai_api_key", respBody, nil)),
		d.Set("alias", utils.PathSearch("alias", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
			respBody, "").(string))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceApplicationAiApiKeyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func deleteApplicationAiApiKey(client *golangsdk.ServiceClient, instanceId, appId, aiApiKeyId string) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/ai-api-keys/{ai_api_key_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{app_id}", appId)
	deletePath = strings.ReplaceAll(deletePath, "{ai_api_key_id}", aiApiKeyId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceApplicationAiApiKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		aiApiKeyId = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = deleteApplicationAiApiKey(client, instanceId, appId, aiApiKeyId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting AI API key (%s) from application (%s)", aiApiKeyId, appId))
	}
	return nil
}

func resourceApplicationAiApiKeyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<application_id>/<id>', "+
			"but got '%s'", importedId)
	}

	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("application_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
