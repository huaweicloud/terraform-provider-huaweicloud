package identitycenter

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

var identityCenterSSOConfigurationNonUpdateParams = []string{"instance_id"}

// @API IdentityCenter POST /v1/instances/{instance_id}/sso-configuration
// @API IdentityCenter GET /v1/instances/{instance_id}/sso-configuration
func ResourceIdentityCenterSSOConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterSSOConfigurationCreateOrUpdate,
		UpdateContext: resourceIdentityCenterSSOConfigurationCreateOrUpdate,
		ReadContext:   resourceIdentityCenterSSOConfigurationRead,
		DeleteContext: resourceIdentityCenterSSOConfigurationDelete,

		CustomizeDiff: config.FlexibleForceNew(identityCenterSSOConfigurationNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"configuration_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mfa_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"no_mfa_signin_behavior": {
				Type:     schema.TypeString,
				Required: true,
			},
			"no_password_signin_behavior": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allowed_mfa_types": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 2,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"max_authentication_age": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIdentityCenterSSOConfigurationCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	instanceId := d.Get("instance_id").(string)

	var (
		httpUrl = "v1/instances/{instance_id}/sso-configuration"
		product = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateSsoConfigurationBodyParams(d)),
	}
	_, err = client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating Identity Center sso configuration: %s", err)
	}
	if d.IsNewResource() {
		d.SetId(instanceId)
	}

	return resourceIdentityCenterSSOConfigurationRead(ctx, d, meta)
}

func buildUpdateSsoConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sso_configuration": map[string]interface{}{
			"mfa_mode":                    d.Get("mfa_mode").(string),
			"no_mfa_signin_behavior":      d.Get("no_mfa_signin_behavior").(string),
			"no_password_signin_behavior": d.Get("no_password_signin_behavior").(string),
			"allowed_mfa_types":           d.Get("allowed_mfa_types").([]interface{}),
			"session_configuration": map[string]interface{}{
				"max_authentication_age": d.Get("max_authentication_age").(string),
			},
		},
		"configuration_type": utils.ValueIgnoreEmpty(d.Get("configuration_type")),
	}
	return bodyParams
}

func resourceIdentityCenterSSOConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/instances/{instance_id}/sso-configuration"
		product = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center sso configuration.")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("mfa_mode", utils.PathSearch("sso_configuration.mfa_mode", getRespBody, nil)),
		d.Set("no_mfa_signin_behavior", utils.PathSearch("sso_configuration.no_mfa_signin_behavior", getRespBody, nil)),
		d.Set("no_password_signin_behavior", utils.PathSearch("sso_configuration.no_password_signin_behavior", getRespBody, nil)),
		d.Set("allowed_mfa_types", utils.PathSearch("sso_configuration.allowed_mfa_types", getRespBody, nil)),
		d.Set("max_authentication_age", utils.PathSearch("sso_configuration.session_configuration.max_authentication_age", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterSSOConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/instances/{instance_id}/sso-configuration"
		product = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDefaultSsoConfigurationBodyParams()),
	}

	_, err = client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating Identity Center sso configuration: %s", err)
	}

	errorMsg := "Deleting sso configuration is unsupported. The resource is removed from the state," +
		" and the sso configuration is reset to default setting."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func buildDefaultSsoConfigurationBodyParams() map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sso_configuration": map[string]interface{}{
			"mfa_mode":                    "ALWAYS_ON",
			"no_mfa_signin_behavior":      "ALLOWED",
			"no_password_signin_behavior": "BLOCKED",
			"allowed_mfa_types":           []string{"TOTP"},
			"session_configuration": map[string]interface{}{
				"max_authentication_age": "PT8H",
			},
		},
		"configuration_type": "APP_AUTHENTICATION_CONFIGURATION",
	}
	return bodyParams
}
