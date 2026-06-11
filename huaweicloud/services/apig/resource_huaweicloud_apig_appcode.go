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

var apigAppcodeNonUpdatableParams = []string{
	"instance_id",
	"application_id",
	"value",
}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes/{app_code_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes/{app_code_id}
func ResourceAppcode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppcodeCreate,
		ReadContext:   resourceAppcodeRead,
		UpdateContext: resourceAppcodeUpdate,
		DeleteContext: resourceAppcodeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAppcodeImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(apigAppcodeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the APPCODE is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the application and APPCODE belong.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application to which the APPCODE belongs.`,
			},

			// Optional parameters.
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The APPCODE value (content).`,
			},

			// Attributes.
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the APPCODE, in RFC3339 format.`,
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

func buildAppcodeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"app_code": d.Get("value"),
	}
}

func createAppcode(client *golangsdk.ServiceClient, instanceId, appId string, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{app_id}", appId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAppcodeBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func autoGenerateAppcode(client *golangsdk.ServiceClient, instanceId, appId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{app_id}", appId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceAppcodeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	var respBody interface{}
	if _, ok := d.GetOk("value"); ok {
		respBody, err = createAppcode(client, instanceId, appId, d)
		if err != nil {
			return diag.Errorf("error creating APPCODE: %s", err)
		}
	} else {
		respBody, err = autoGenerateAppcode(client, instanceId, appId)
		if err != nil {
			return diag.Errorf("error auto generating APPCODE: %s", err)
		}
	}

	resourceId := utils.PathSearch("id", respBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("unable to find the APPCODE ID from the API response")
	}
	d.SetId(resourceId)

	return resourceAppcodeRead(ctx, d, meta)
}

func GetAppcode(client *golangsdk.ServiceClient, instanceId, appId, appCodeId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes/{app_code_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{app_id}", appId)
	getPath = strings.ReplaceAll(getPath, "{app_code_id}", appCodeId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceAppcodeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		appCodeId  = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	respBody, err := GetAppcode(client, instanceId, appId, appCodeId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error querying APPCODE (%s) from application (%s)", appCodeId, appId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("application_id", utils.PathSearch("app_id", respBody, nil)),
		d.Set("value", utils.PathSearch("app_code", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
			respBody, "").(string))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAppcodeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func deleteAppcode(client *golangsdk.ServiceClient, instanceId, appId, appCodeId string) error {
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

func resourceAppcodeDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		appCodeId  = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = deleteAppcode(client, instanceId, appId, appCodeId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error deleting APPCODE (%s) from application (%s)", appCodeId, appId))
	}
	return nil
}

func resourceAppcodeImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
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
