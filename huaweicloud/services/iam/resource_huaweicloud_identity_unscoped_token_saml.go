package iam

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceIdentityUnscopedTokenSaml
// @API IAM POST /v3.0/OS-FEDERATION/tokens
func ResourceIdentityUnscopedTokenSaml() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUnscopedTokenSamlCreate,
		ReadContext:   resourceUnscopedTokenSamlRead,
		DeleteContext: resourceUserTokenDelete,

		Schema: map[string]*schema.Schema{
			"idp_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity Provider Id",
			},
			"saml_response": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The response body returned after successful IdP authentication.",
			},
			"with_global_domain": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Whether to use a global domain name to obtain the token. Its default value is false",
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

func resourceUnscopedTokenSamlCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var client *golangsdk.ServiceClient
	var err error
	if d.Get("with_global_domain").(bool) {
		client, err = cfg.IAMV3GlobalClient()
	} else {
		client, err = cfg.IAMV3Client(cfg.GetRegion(d))
	}
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	unscopedTokenSamlPath := client.Endpoint + "v3.0/OS-FEDERATION/tokens"
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"X-Idp-Id":     d.Get("idp_id").(string),
		},
		RawBody: strings.NewReader("SAMLResponse=" + url.QueryEscape(d.Get("saml_response").(string))),
	}
	response, err := client.Request("POST", unscopedTokenSamlPath, &options)
	if err != nil {
		return diag.Errorf("error createFederationToken: %s", err)
	}
	if err = setFederationToken(d, response); err != nil {
		return diag.Errorf("error setFederationToken fields: %s", err)
	}
	return nil
}

func setFederationToken(d *schema.ResourceData, response *http.Response) error {
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return err
	}
	username := utils.PathSearch("token.user.name", respBody, "").(string)
	groupsModel := utils.PathSearch("token.user.\"OS-FEDERATION\".groups", respBody, make([]interface{}, 0)).([]interface{})
	groups := make([]string, 0, len(groupsModel))
	for _, groupModel := range groupsModel {
		groups = append(groups, utils.PathSearch("name", groupModel, "").(string))
	}

	d.SetId(d.Get("idp_id").(string) + ":" + username)
	mErr := multierror.Append(
		d.Set("token", response.Header.Get("X-Subject-Token")),
		d.Set("username", username),
		d.Set("expires_at", utils.PathSearch("token.expires_at", respBody, "").(string)),
		d.Set("groups", groups),
	)
	return mErr.ErrorOrNil()
}

func resourceUnscopedTokenSamlRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// A SAMLResponse only supports calling the interface once, so we don't try to flush the token here
	return nil
}
