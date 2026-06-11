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

type SecretAction string

const (
	SecretActionReset SecretAction = "RESET"
)

var applicationNonUpdatableParams = []string{
	"instance_id",
}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apps
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes/{app_code_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/apps/secret/{app_id}
func ResourceApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(applicationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the application is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the application belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The application name.`,
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The application description.`,
			},
			"app_codes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				MaxItems:    5,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The array of one or more application codes that the application has.`,
			},
			"app_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: `The APP key.`,
			},
			"app_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: `The APP secret.`,
			},
			"secret_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The secret action to be done for the application.`,
			},

			// Attributes.
			"registration_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The registration time of the application, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the application, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true,
						Required: true,
					}),
			},
		},
	}
}

func buildApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":       d.Get("name"),
		"remark":     d.Get("description"),
		"app_key":    utils.ValueIgnoreEmpty(d.Get("app_key")),
		"app_secret": utils.ValueIgnoreEmpty(d.Get("app_secret")),
	}
}

func createApplication(client *golangsdk.ServiceClient, instanceId string, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildApplicationBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func createApplicationAppcode(client *golangsdk.ServiceClient, instanceId, appId, code string) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{app_id}", appId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"app_code": code,
		},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func createApplicationAppcodes(client *golangsdk.ServiceClient, instanceId, appId string, codes []interface{}) error {
	for _, code := range codes {
		if err := createApplicationAppcode(client, instanceId, appId, code.(string)); err != nil {
			return fmt.Errorf("error creating application code: %s", err)
		}
	}
	return nil
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	respBody, err := createApplication(client, instanceId, d)
	if err != nil {
		return diag.Errorf("error creating dedicated application: %s", err)
	}

	appId := utils.PathSearch("id", respBody, "").(string)
	if appId == "" {
		return diag.Errorf("unable to find the application ID from the API response")
	}
	d.SetId(appId)

	if v, ok := d.GetOk("app_codes"); ok {
		if err := createApplicationAppcodes(client, instanceId, appId, v.(*schema.Set).List()); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceApplicationRead(ctx, d, meta)
}

func GetApplication(client *golangsdk.ServiceClient, instanceId, appId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{app_id}", appId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func flattenApplicationAppcodes(appcodes []interface{}) []interface{} {
	if len(appcodes) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(appcodes))
	for _, appcode := range appcodes {
		result = append(result, utils.PathSearch("app_code", appcode, nil))
	}
	return result
}

func resourceApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		appId      = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	respBody, err := GetApplication(client, instanceId, appId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "dedicated application")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("remark", respBody, nil)),
		d.Set("registration_time", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("register_time",
			respBody, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("update_time",
			respBody, "").(string))/1000, false)),
		d.Set("app_key", utils.PathSearch("app_key", respBody, nil)),
		d.Set("app_secret", utils.PathSearch("app_secret", respBody, nil)),
	)
	if appcodes, err := listAppcodes(client, instanceId, appId); err != nil {
		mErr = multierror.Append(mErr, err)
	} else {
		// The application code is sort by create time on server, not code.
		mErr = multierror.Append(mErr, d.Set("app_codes",
			schema.NewSet(schema.HashString, flattenApplicationAppcodes(appcodes))))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateApplication(client *golangsdk.ServiceClient, instanceId, appId string, d *schema.ResourceData) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)
	updatePath = strings.ReplaceAll(updatePath, "{app_id}", appId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildApplicationBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func getAppcodeIdByValue(appcodes []interface{}, code string) (string, bool) {
	for _, appcode := range appcodes {
		if utils.PathSearch("app_code", appcode, "").(string) == code {
			return utils.PathSearch("id", appcode, "").(string), true
		}
	}
	return "", false
}

func deleteApplicationAppcode(client *golangsdk.ServiceClient, instanceId, appId, appCodeId string) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes/{app_code_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{app_id}", appId)
	deletePath = strings.ReplaceAll(deletePath, "{app_code_id}", appCodeId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func removeApplicationAppcodes(client *golangsdk.ServiceClient, instanceId, appId string, codes []interface{}) error {
	appcodes, err := listAppcodes(client, instanceId, appId)
	if err != nil {
		return fmt.Errorf("error retrieving application codes: %s", err)
	}
	for _, code := range codes {
		appCodeId, ok := getAppcodeIdByValue(appcodes, code.(string))
		if !ok {
			continue
		}
		if err := deleteApplicationAppcode(client, instanceId, appId, appCodeId); err != nil {
			return fmt.Errorf("error removing code (%v) from the application (%s): %s", code, appId, err)
		}
	}
	return nil
}

func updateApplicationAppcodes(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		instanceId       = d.Get("instance_id").(string)
		appId            = d.Id()
		oldRaws, newRaws = d.GetChange("app_codes")

		addRaws    = newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
		removeRaws = oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	)
	if removeRaws.Len() > 0 {
		if err := removeApplicationAppcodes(client, instanceId, appId, removeRaws.List()); err != nil {
			return err
		}
	}

	if addRaws.Len() > 0 {
		if err := createApplicationAppcodes(client, instanceId, appId, addRaws.List()); err != nil {
			return err
		}
	}
	return nil
}

func resetApplicationSecret(client *golangsdk.ServiceClient, instanceId, appId string) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/secret/{app_id}"
	resetPath := client.Endpoint + httpUrl
	resetPath = strings.ReplaceAll(resetPath, "{project_id}", client.ProjectID)
	resetPath = strings.ReplaceAll(resetPath, "{instance_id}", instanceId)
	resetPath = strings.ReplaceAll(resetPath, "{app_id}", appId)

	resetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         map[string]interface{}{},
	}

	_, err := client.Request("PUT", resetPath, &resetOpt)
	return err
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		appId      = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	if d.HasChanges("name", "description", "app_key", "app_secret") {
		if err := updateApplication(client, instanceId, appId, d); err != nil {
			return diag.Errorf("error updating dedicated application (%s): %s", appId, err)
		}
	}
	if d.HasChange("app_codes") {
		if err := updateApplicationAppcodes(client, d); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("secret_action") {
		if v, ok := d.GetOk("secret_action"); ok && v.(string) == string(SecretActionReset) {
			if err := resetApplicationSecret(client, instanceId, appId); err != nil {
				return diag.Errorf("error reseting application secret: %s", err)
			}
		}
	}
	return resourceApplicationRead(ctx, d, meta)
}

func deleteApplication(client *golangsdk.ServiceClient, instanceId, appId string) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{app_id}", appId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		appId      = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = deleteApplication(client, instanceId, appId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error deleting application (%s) from the instance (%s)", appId, instanceId))
	}

	return nil
}

func resourceApplicationImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
