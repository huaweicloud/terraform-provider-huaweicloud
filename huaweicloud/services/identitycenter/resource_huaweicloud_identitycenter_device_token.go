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
	"strconv"
	"time"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityStore POST /v1/tokens
func ResourceIdentityCenterDeviceToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterDeviceTokenCreate,
		ReadContext:   resourceIdentityCenterDeviceTokenRead,
		DeleteContext: resourceIdentityCenterDeviceTokenDelete,
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
			"code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"device_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"grant_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"redirect_uri": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"refresh_token": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"access_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_in": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterDeviceTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createIdentityCenterClient: create IdentityCenter device token
	var (
		createIdentityCenterDeviceTokenHttpUrl = "v1/tokens"
		createIdentityCenterDeviceTokenProduct = "identityoidc"
	)
	createIdentityCenterDeviceTokenClient, err := cfg.NewServiceClient(createIdentityCenterDeviceTokenProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center device token: %s", err)
	}

	createIdentityCenterDeviceTokenPath := createIdentityCenterDeviceTokenClient.Endpoint + createIdentityCenterDeviceTokenHttpUrl
	createIdentityCenterDeviceTokenOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	log.Println(createIdentityCenterDeviceTokenOpt)
	createIdentityCenterDeviceTokenOpt.JSONBody = utils.RemoveNil(buildCreateIdentityCenterDeviceTokenBodyParams(d))
	createIdentityCenterClientResp, err := createIdentityCenterDeviceTokenClient.Request("POST",
		createIdentityCenterDeviceTokenPath, &createIdentityCenterDeviceTokenOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center device token: %s", err)
	}

	createIdentityCenterDeviceTokenRespBody, err := utils.FlattenResponse(createIdentityCenterClientResp)
	if err != nil {
		return diag.FromErr(err)
	}

	token := utils.PathSearch("token_info.access_token", createIdentityCenterDeviceTokenRespBody, "").(string)
	if token == "" {
		return diag.Errorf("unable to find the Identity Center access_token from the API response")
	}
	d.SetId(token)
	d.Set("access_token", token)

	expiredAtFloat := utils.PathSearch("token_info.expires_in", createIdentityCenterDeviceTokenRespBody, 0).(float64)
	expiredAt := int64(expiredAtFloat)
	if expiredAt < time.Now().Unix() {
		return diag.Errorf("unable to find the Identity Center expires_in from the API response")
	}
	expiredAtString := strconv.FormatInt(expiredAt, 10)
	d.Set("expires_in", expiredAtString)

	idToken := utils.PathSearch("token_info.id_token", createIdentityCenterDeviceTokenRespBody, "").(string)
	d.Set("id_token", idToken)

	refreshToken := utils.PathSearch("token_info.refresh_token", createIdentityCenterDeviceTokenRespBody, "").(string)
	d.Set("refresh_token", refreshToken)
	return nil
}

func buildCreateIdentityCenterDeviceTokenBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"client_id":     utils.ValueIgnoreEmpty(d.Get("client_id")),
		"client_secret": utils.ValueIgnoreEmpty(d.Get("client_secret")),
		"code":          utils.ValueIgnoreEmpty(d.Get("code")),
		"device_code":   utils.ValueIgnoreEmpty(d.Get("device_code")),
		"grant_type":    utils.ValueIgnoreEmpty(d.Get("grant_type")),
		"redirect_uri":  utils.ValueIgnoreEmpty(d.Get("redirect_uri")),
		"refresh_token": utils.ValueIgnoreEmpty(d.Get("refresh_token")),
		"scopes":        d.Get("scopes"),
	}
	return bodyParams
}

func resourceIdentityCenterDeviceTokenRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// If client expired, create a new token.
	expiresAtString := d.Get("expires_in").(string)
	if expiresAtString == "" {
		expiresAtString = "1"
	}
	expiresAt, err := strconv.ParseInt(expiresAtString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	if time.Now().Unix() > expiresAt {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "token has expired, please provide a new device code to generate a new token",
			},
		}
	}
	return nil
}

func resourceIdentityCenterDeviceTokenDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting token is not supported."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
