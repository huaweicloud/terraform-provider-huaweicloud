package iam

import (
	"context"
	"errors"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/oidcconfig"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIdentityProviderOidcConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resIdeProviderOidcConfigCreate,
		ReadContext:   resIdeProviderOidcConfigRead,
		UpdateContext: resIdeProviderOidcConfigUpdate,
		DeleteContext: resIdeProviderOidcConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"program", "program_console"}, false),
			},
			"provider_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"signing_key": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"authorization_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"response_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "id_token",
			},
			"response_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "form_post",
				ValidateFunc: validation.StringInSlice([]string{"fragment", "form_post"}, false),
			},
		},
	}
}

func resIdeProviderOidcConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud IAM client without version number: %s", err)
	}

	accessType := d.Get("access_type").(string)
	opts := oidcconfig.CreateOpts{
		AccessMode: accessType,
		IdpURL:     d.Get("provider_url").(string),
		ClientID:   d.Get("client_id").(string),
		SigningKey: d.Get("signing_key").(string),
	}
	if accessType == "program_console" {
		opts.AuthorizationEndpoint = d.Get("authorization_endpoint").(string)
		opts.Scope = getScopes(d)
		opts.ResponseType = d.Get("response_type").(string)
		opts.ResponseMode = d.Get("response_mode").(string)
	}
	logp.Printf("[DEBUG] Create access type of provider: %#v", opts)

	providerID := d.Get("provider_id").(string)
	_, err = oidcconfig.Get(client, providerID)
	if err != nil && errors.As(err, &golangsdk.ErrDefault404{}) {
		_, err = oidcconfig.Create(client, providerID, opts)
	} else {
		err = updateOidcConfig(client, d)
	}

	if err != nil {
		return fmtp.DiagErrorf("Error creating provider access type: %s", err)
	}
	d.SetId(providerID)

	return resIdeProviderOidcConfigRead(ctx, d, meta)
}

func resIdeProviderOidcConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud IAM client without version number: %s", err)
	}
	// ID is the provider name
	providerID := d.Id()
	accessType, err := oidcconfig.Get(client, providerID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error in querying provider access type")
	}

	scopes := make([]string, 0)
	if len(accessType.Scope) > 0 {
		scopes = strings.Split(accessType.Scope, scopeSpilt)
	}
	mErr := multierror.Append(
		d.Set("provider_id", d.Id()),
		d.Set("access_type", accessType.AccessMode),
		d.Set("provider_url", accessType.IdpURL),
		d.Set("client_id", accessType.ClientID),
		d.Set("signing_key", accessType.SigningKey),
		d.Set("scopes", scopes),
		d.Set("authorization_endpoint", accessType.AuthorizationEndpoint),
	)

	if accessType.AccessMode == "program_console" {
		mErr = multierror.Append(
			mErr,
			d.Set("response_mode", accessType.ResponseMode),
			d.Set("response_type", accessType.ResponseType),
		)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		logp.Printf("[ERROR] Error setting identity provider config %s: %s", d.Id(), err)
		return fmtp.DiagErrorf("Error setting identity provider config : %s", err)
	}
	return nil
}

func getScopes(d *schema.ResourceData) string {
	scopes := utils.ExpandToStringList(d.Get("scopes").([]interface{}))
	str := strings.Join(scopes, scopeSpilt)
	return str
}

func updateOidcConfig(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	accessType := d.Get("access_type").(string)
	opts := oidcconfig.UpdateOpenIDConnectConfigOpts{
		AccessMode: accessType,
		IdpURL:     d.Get("provider_url").(string),
		ClientID:   d.Get("client_id").(string),
		SigningKey: d.Get("signing_key").(string),
	}

	if accessType == "program_console" {
		opts.AuthorizationEndpoint = d.Get("authorization_endpoint").(string)
		opts.Scope = getScopes(d)
		opts.ResponseType = d.Get("response_type").(string)
		opts.ResponseMode = d.Get("response_mode").(string)
	}
	logp.Printf("[DEBUG] Update access type of provider: %#v", opts)
	providerID := d.Get("provider_id").(string)
	_, err := oidcconfig.Update(client, providerID, opts)
	return err
}

func resIdeProviderOidcConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error updating HuaweiCloud IAM client without version number: %s", err)
	}
	err = updateOidcConfig(client, d)
	if err != nil {
		return fmtp.DiagErrorf("Error updating the access type of provider: %s", err)
	}

	return resIdeProviderOidcConfigRead(ctx, d, meta)
}

func resIdeProviderOidcConfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error updating HuaweiCloud IAM client without version number: %s", err)
	}
	// The oidc configuration cannot be removed from the provider, so we set it to an invalid value.
	signingKey := `{"keys": [{"use": "sig", "e": "AQAB", "kty": "RSA", "alg": "RS256", "n": "..", "kid": ".."}]}`
	opts := oidcconfig.UpdateOpenIDConnectConfigOpts{
		AccessMode:            "program_console",
		IdpURL:                "https://idp_url",
		ClientID:              "__client_id__",
		AuthorizationEndpoint: "https://endpoint",
		Scope:                 "openid",
		ResponseType:          "id_token",
		ResponseMode:          "form_post",
		SigningKey:            signingKey,
	}
	providerID := d.Get("provider_id").(string)
	_, err = oidcconfig.Update(client, providerID, opts)
	if err != nil {
		return fmtp.DiagErrorf("error deleting access type: %s", err)
	}
	d.SetId("")
	return nil
}
