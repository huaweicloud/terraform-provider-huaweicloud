package iam

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ResourceIdentityUnscopedTokenWithIdToken
// @API IAM POST /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}/auth
func ResourceIdentityUnscopedTokenWithIdToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUnscopedTokenWithIdTokenCreate,
		ReadContext:   resourceUnscopedTokenWithIdTokenRead,
		DeleteContext: resourceUserTokenDelete,

		Schema: map[string]*schema.Schema{
			"idp_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The identity provider id.",
			},
			"protocol_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The protocol id.",
			},
			"id_token": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The security token of the OpenID Connect Identity Provider, format is Bearer {ID Token}.",
			},

			// Attributes
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user name.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The details of the obtained token.",
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The user group information.",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration time.",
			},
		},
	}
}

func resourceUnscopedTokenWithIdTokenCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client := common.NewCustomClient(true, "https://iam.{region_id}.myhuaweicloud.com")
	unscopedTokenWithIdTokenPath := client.ResourceBase + "v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}/auth"
	unscopedTokenWithIdTokenPath = strings.ReplaceAll(unscopedTokenWithIdTokenPath, "{region_id}", cfg.GetRegion(d))
	unscopedTokenWithIdTokenPath = strings.ReplaceAll(unscopedTokenWithIdTokenPath, "{idp_id}", d.Get("idp_id").(string))
	unscopedTokenWithIdTokenPath = strings.ReplaceAll(unscopedTokenWithIdTokenPath, "{protocol_id}", d.Get("protocol_id").(string))
	authorization := "Bearer " + d.Get("id_token").(string)
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Authorization": authorization,
		},
	}
	response, err := client.Request("POST", unscopedTokenWithIdTokenPath, &options)
	if err != nil {
		return diag.Errorf("error createFederationToken: %s", err)
	}
	if err = setFederationToken(d, response); err != nil {
		return diag.Errorf("error setFederationToken fields: %s", err)
	}
	return nil
}

func resourceUnscopedTokenWithIdTokenRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// If token expired, try to create a new token.
	expiresAt, err := time.ParseInLocation(`2006-01-02T15:04:05Z`, d.Get("expires_at").(string), time.UTC)
	if err != nil {
		return diag.Errorf("error parsing expires at: %s", err)
	}
	if time.Now().After(expiresAt) {
		return resourceUnscopedTokenWithIdTokenCreate(c, d, meta)
	}
	return nil
}
