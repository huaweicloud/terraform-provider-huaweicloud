package iam

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v3UnscopedTokenSamlNonUpdatableParams = []string{
	"idp_id",
	"saml_response",
	"with_global_domain",
}

// @API IAM POST /v3.0/OS-FEDERATION/tokens
func ResourceV3UnscopedTokenSaml() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3UnscopedTokenSamlCreate,
		ReadContext:   resourceV3UnscopedTokenSamlRead,
		UpdateContext: resourceV3UnscopedTokenSamlUpdate,
		DeleteContext: resourceV3UnscopedTokenSamlDelete,

		CustomizeDiff: config.FlexibleForceNew(v3UnscopedTokenSamlNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"idp_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The identity provider id.`,
			},
			"saml_response": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The response body returned after successful IDP authentication.`,
			},
			"with_global_domain": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether to use a global domain name to obtain the token.`,
			},

			// Attributes
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unscoped token.`,
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user of token.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The group list of the user.`,
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the token will expire.`,
			},

			// Internal
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceV3UnscopedTokenSamlCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error creating unscoped token: %s", err)
	}
	if err = setFederationToken(d, response); err != nil {
		return diag.Errorf("error setting federation token fields: %s", err)
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

func resourceV3UnscopedTokenSamlRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// A SAMLResponse only supports calling the interface once, so we don't try to flush the token here
	return nil
}

func resourceV3UnscopedTokenSamlUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV3UnscopedTokenSamlDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for creating unscoped token. Deleting this resource will
    not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
