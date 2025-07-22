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

// @API IdentityStore POST /v1/clients
func ResourceIdentityCenterClient() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterClientCreate,
		ReadContext:   resourceIdentityCenterClientRead,
		DeleteContext: resourceIdentityCenterClientDelete,
		Description:   "schema: Internal",
		Schema: map[string]*schema.Schema{
			"client_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token_endpoint_auth_method": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"grant_types": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"response_types": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_secret_expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterClientCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return registerClient(d, meta)
}

func registerClient(d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		createIdentityCenterClientHttpUrl = "v1/clients"
		createIdentityCenterClientProduct = "identityoidc"
	)
	createIdentityCenterClientClient, err := cfg.NewServiceClient(createIdentityCenterClientProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createIdentityCenterClientPath := createIdentityCenterClientClient.Endpoint + createIdentityCenterClientHttpUrl
	createIdentityCenterClientOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createIdentityCenterClientOpt.JSONBody = utils.RemoveNil(buildCreateIdentityCenterClientBodyParams(d))
	createIdentityCenterClientResp, err := createIdentityCenterClientClient.Request("POST",
		createIdentityCenterClientPath, &createIdentityCenterClientOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createIdentityCenterClientRespBody, err := utils.FlattenResponse(createIdentityCenterClientResp)
	if err != nil {
		return diag.FromErr(err)
	}

	clientId := utils.PathSearch("client_info.client_id", createIdentityCenterClientRespBody, "").(string)
	if clientId == "" {
		return diag.Errorf("unable to find the Identity Center Client ID from the API response")
	}
	d.SetId(clientId)
	d.Set("client_id", clientId)
	clientSecret := utils.PathSearch("client_info.client_secret", createIdentityCenterClientRespBody, "").(string)
	if clientSecret == "" {
		return diag.Errorf("unable to find the Identity Center Client ID from the API response")
	}
	d.Set("client_secret", clientSecret)

	expiredAtFloat := utils.PathSearch("client_info.client_secret_expires_at", createIdentityCenterClientRespBody, 0).(float64)
	expiredAt := int64(expiredAtFloat)
	if expiredAt < time.Now().Unix() {
		return diag.Errorf("unable to find the Identity Center client_secret_expires_at from the API response")
	}
	expiredAtString := strconv.FormatInt(expiredAt, 10)
	d.Set("client_secret_expires_at", expiredAtString)
	return nil
}

func buildCreateIdentityCenterClientBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"client_name":                utils.ValueIgnoreEmpty(d.Get("client_name")),
		"client_type":                utils.ValueIgnoreEmpty(d.Get("client_type")),
		"token_endpoint_auth_method": utils.ValueIgnoreEmpty(d.Get("token_endpoint_auth_method")),
		"scopes":                     d.Get("scopes"),
		"grant_types":                d.Get("grant_types"),
		"response_types":             d.Get("response_types"),
	}
	return bodyParams
}

func resourceIdentityCenterClientRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// If client expired, create a new token.
	expiresAtString := d.Get("client_secret_expires_at").(string)
	if expiresAtString == "" {
		expiresAtString = "1"
	}
	expiresAt, err := strconv.ParseInt(expiresAtString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	if time.Now().Unix() > expiresAt {
		log.Println("client has expired, create a new client")
		registerClient(d, meta)
	}
	return nil
}

func resourceIdentityCenterClientDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting client is not supported."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
