package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM POST /v3/auth/tokens
func ResourceIdentityUserToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserTokenCreate,
		ReadContext:   resourceUserTokenRead,
		DeleteContext: resourceUserTokenDelete,

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUserTokenCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	product := "identity"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	// create token
	err = createToken(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("account_name").(string) + "/" + d.Get("user_name").(string))

	return nil
}

func createToken(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	createTokenHttpUrl := "v3/auth/tokens"
	createTokenPath := client.Endpoint + createTokenHttpUrl
	createTokenOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"auth": buildCreateTokenBodyParams(d),
		},
	}

	createTokenResp, err := client.Request("POST", createTokenPath, &createTokenOpt)
	if err != nil {
		return fmt.Errorf("error retrieving IAM user token: %s", err)
	}

	d.Set("token", createTokenResp.Header.Get("X-Subject-Token"))

	createTokenRespBody, err := utils.FlattenResponse(createTokenResp)
	if err != nil {
		return fmt.Errorf("error flattening IAM user token: %s", err)
	}

	expiresAt := utils.PathSearch("token.expires_at", createTokenRespBody, nil)
	if expiresAt == nil {
		return fmt.Errorf("error retrieving IAM user token: expires_at is not found in API response")
	}
	d.Set("expires_at", expiresAt)

	return nil
}

func buildCreateTokenBodyParams(d *schema.ResourceData) map[string]interface{} {
	scope := map[string]interface{}{
		"domain": map[string]interface{}{
			"name": d.Get("account_name"),
		},
	}
	if val, ok := d.GetOk("project_name"); ok {
		scope = map[string]interface{}{
			"project": map[string]interface{}{
				"name": val,
			},
		}
	}
	bodyParams := map[string]interface{}{
		"identity": map[string]interface{}{
			"methods": []string{"password"},
			"password": map[string]interface{}{
				"user": map[string]interface{}{
					"domain": map[string]interface{}{
						"name": d.Get("account_name"),
					},
					"name":     d.Get("user_name"),
					"password": d.Get("password"),
				},
			},
		},
		"scope": scope,
	}
	return bodyParams
}

func resourceUserTokenRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	product := "identity"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	// If token expired, create a new token.
	expiresAt, err := time.ParseInLocation(`2006-01-02T15:04:05Z`, d.Get("expires_at").(string), time.UTC)
	if err != nil {
		diag.Errorf("error parsing expires at: %s", err)
	}
	if time.Now().After(expiresAt) {
		err = createToken(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceUserTokenDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting token is not supported. The token is only removed from the state, but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
