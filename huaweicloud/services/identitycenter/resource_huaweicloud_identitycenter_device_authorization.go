// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IdentityCenter
// ---------------------------------------------------------------

package identitycenter

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityStore POST /v1/device/authorize
func ResourceIdentityCenterDeviceAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterDeviceAuthorizationCreate,
		ReadContext:   resourceIdentityCenterDeviceAuthorizationRead,
		DeleteContext: resourceIdentityCenterDeviceAuthorizationDelete,
		Description:   "schema: Internal",
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_secret": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"start_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"device_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"verification_uri_complete": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterDeviceAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return startDeviceAuthorization(d, meta)
}

func startDeviceAuthorization(d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createIdentityCenterClient: create IdentityCenter client
	var (
		createIdentityCenterClientHttpUrl              = "v1/device/authorize"
		createIdentityCenterDeviceAuthorizationProduct = "identityoidc"
	)
	createIdentityCenterDeviceAuthorizationClient, err := cfg.NewServiceClient(createIdentityCenterDeviceAuthorizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center device authorization: %s", err)
	}

	createIdentityCenterDeviceAuthorizationPath := createIdentityCenterDeviceAuthorizationClient.Endpoint + createIdentityCenterClientHttpUrl
	createIdentityCenterDeviceAuthorizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	log.Println(createIdentityCenterDeviceAuthorizationOpt)
	createIdentityCenterDeviceAuthorizationOpt.JSONBody = utils.RemoveNil(buildCreateIdentityCenterDeviceAuthorizationBodyParams(d))
	createIdentityCenterClientResp, err := createIdentityCenterDeviceAuthorizationClient.Request("POST",
		createIdentityCenterDeviceAuthorizationPath, &createIdentityCenterDeviceAuthorizationOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center device authorization: %s", err)
	}

	createIdentityCenterDeviceAuthorizationRespBody, err := utils.FlattenResponse(createIdentityCenterClientResp)
	if err != nil {
		return diag.FromErr(err)
	}

	deviceCode := utils.PathSearch("device_code", createIdentityCenterDeviceAuthorizationRespBody, "").(string)
	if deviceCode == "" {
		return diag.Errorf("unable to find the Identity Center device_code from the API response")
	}
	d.SetId(deviceCode)
	d.Set("device_code", deviceCode)
	verificationUriComplete := utils.PathSearch("verification_uri_complete", createIdentityCenterDeviceAuthorizationRespBody, "").(string)
	if verificationUriComplete == "" {
		return diag.Errorf("unable to find the Identity Center verification_uri_complete from the API response")
	}
	d.Set("verification_uri_complete", verificationUriComplete)
	return nil
}

func buildCreateIdentityCenterDeviceAuthorizationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"client_id":     utils.ValueIgnoreEmpty(d.Get("client_id")),
		"client_secret": utils.ValueIgnoreEmpty(d.Get("client_secret")),
		"start_url":     utils.ValueIgnoreEmpty(d.Get("start_url")),
	}
	return bodyParams
}

func resourceIdentityCenterDeviceAuthorizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// everytime generate a new device code
	startDeviceAuthorization(d, meta)
	return nil
}

func resourceIdentityCenterDeviceAuthorizationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting device authorization is not supported."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
