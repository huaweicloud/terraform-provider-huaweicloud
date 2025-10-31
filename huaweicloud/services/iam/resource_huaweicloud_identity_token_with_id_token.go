package iam

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ResourceIdentityTokenWithIdToken
// @API IAM POST /v3.0/OS-AUTH/id-token/tokens
func ResourceIdentityTokenWithIdToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTokenWithIdTokenCreate,
		ReadContext:   resourceTokenWithIdTokenRead,
		DeleteContext: resourceUserTokenDelete,

		Schema: map[string]*schema.Schema{
			"idp_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity Provider Id",
			},
			"id_token": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID Token of the OpenID Connect Identity Provider",
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"domain_name"},
			},

			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTokenWithIdTokenCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	tokenWithIdTokenPath := client.Endpoint + "v3.0/OS-AUTH/id-token/tokens"
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Idp-Id": d.Get("idp_id").(string),
		},
		JSONBody: map[string]interface{}{
			"auth": buildTokenWithIdTokenBody(d),
		},
	}
	response, err := client.Request("POST", tokenWithIdTokenPath, &options)
	if err != nil {
		return diag.Errorf("error createFederationToken: %s", err)
	}
	if err = setFederationToken(d, response); err != nil {
		return diag.Errorf("error setFederationToken fields: %s", err)
	}
	return nil
}

func buildTokenWithIdTokenBody(d *schema.ResourceData) map[string]interface{} {
	var scope map[string]interface{}
	if domainName, ok := d.GetOk("domain_name"); ok {
		scope = map[string]interface{}{
			"domain": map[string]interface{}{
				"name": domainName,
			},
		}
	}
	if projectName, ok := d.GetOk("project_name"); ok {
		scope = map[string]interface{}{
			"project": map[string]interface{}{
				"name": projectName,
			},
		}
	}
	body := map[string]interface{}{
		"id_token": map[string]interface{}{
			"id": d.Get("id_token").(string),
		},
		"scope": scope,
	}
	return body
}

func resourceTokenWithIdTokenRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// If token expired, try to create a new token.
	expiresAt, err := time.ParseInLocation(`2006-01-02T15:04:05Z`, d.Get("expires_at").(string), time.UTC)
	if err != nil {
		return diag.Errorf("error parsing expires at: %s", err)
	}
	if time.Now().After(expiresAt) {
		return resourceTokenWithIdTokenCreate(c, d, meta)
	}
	return nil
}
