package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM POST /v5/agencies/assume-with-oidc
func DataSourceV5OidcCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5OidcCredentialRead,

		Schema: map[string]*schema.Schema{
			"id_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Specifies the OpenID Connect ID token issued by the OIDC identity provider.`,
			},
			"agency_urn": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the uniform resource name of the trust agency to be assumed, in the format of ` +
					`iam::{account_id}:agency:{agency_name}.`,
			},
			"agency_session_name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the name of the agency session. The name contains 2 to 64 characters and can only contain letters, digits, ` +
					`underscores (_) and hyphens (-).`,
			},
			"provider_urn": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Specifies the uniform resource name of the OIDC identity provider, in the format of ` +
					`iam::{account_id}:oidcProvider:{provider_name}.`,
			},
			"duration_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(900, 43200),
				Description:  `Specifies the validity period of the temporary security credentials, in seconds. The value ranges from 900 to 43200 and defaults to 3600.`,
			},
			"policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `Specifies a custom policy in JSON format, which is used to further restrict the permissions of the temporary security credentials.`,
			},
			"policy_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of preset policy IDs, which is used to restrict the permissions of the temporary security credentials. A maximum of 20 policy IDs can be configured.`,
			},

			// Computed attributes.
			"assumed_agency": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The uniform resource name of the assumed agency session.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the assumed agency session.`,
						},
					},
				},
				Description: `The information about the assumed agency session.`,
			},
			"credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: `The access key (AK) of the temporary security credentials.`,
						},
						"secret_access_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: `The secret access key (SK) of the temporary security credentials.`,
						},
						"security_token": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: `The security token of the temporary security credentials, which is used as the request header x-security-token.`,
						},
						"expiration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The expiration time (UTC) of the temporary security credentials.`,
						},
					},
				},
				Description: `The temporary security credentials.`,
			},
			"subject_from_id_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique user identifier returned from the sub field of the ID token.`,
			},
			"audience": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The client ID returned from the aud or azp field of the ID token.`,
			},
			"idp_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the OIDC identity provider.`,
			},
			"source_identity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source identity declared by the source identity field of the ID token.`,
			},
		},
	}
}

func dataSourceV5OidcCredentialRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("sts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating STS client: %s", err)
	}

	assumeHttpUrl := "v5/agencies/assume-with-oidc"
	assumePath := client.Endpoint + assumeHttpUrl
	assumeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildV5AgencyAssumeWithOidcBodyParams(d),
	}
	assumeResp, err := client.Request("POST", assumePath, &assumeOpt)
	if err != nil {
		return diag.Errorf("error assuming agency with OIDC: %s", err)
	}

	assumeBody, err := utils.FlattenResponse(assumeResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomId)

	if diag := setV5AgencyCredentialAttributes(d, assumeBody); diag.HasError() {
		return diag
	}

	mErr := multierror.Append(
		d.Set("subject_from_id_token", utils.PathSearch("subject_from_id_token", assumeBody, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OIDC agency credential session fields: %s", err)
	}

	return nil
}

// setV5AgencyCredentialAttributes is shared between the OIDC and SAML assume
// flows and writes the API response (assumed_agency / credentials blocks and
// the top-level session fields) into the data source. The response layout
// follows the STS v5 AssumeAgencyWithSAML / AssumeAgencyWithOIDC reference:
//   - assumed_agency: { urn, id }
//   - credentials:    { access_key_id, secret_access_key, security_token, expiration }
//   - top-level:      audience, provider (idp_urn), source_identity (OIDC),
//     subject / subject_type / issuer / name_qualifier (SAML), subject_from_id_token (OIDC)

func buildV5AgencyAssumeWithOidcBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id_token":            d.Get("id_token").(string),
		"agency_urn":          d.Get("agency_urn").(string),
		"agency_session_name": d.Get("agency_session_name").(string),
	}
	if providerUrn := d.Get("provider_urn").(string); providerUrn != "" {
		bodyParams["provider_urn"] = providerUrn
	}
	if durationSeconds := d.Get("duration_seconds").(int); durationSeconds != 0 {
		bodyParams["duration_seconds"] = durationSeconds
	}
	if policy := d.Get("policy").(string); policy != "" {
		bodyParams["policy"] = policy
	}
	if policyIds := d.Get("policy_ids").([]interface{}); len(policyIds) > 0 {
		bodyParams["policy_ids"] = policyIds
	}
	return bodyParams
}

func setV5AgencyCredentialAttributes(d *schema.ResourceData, assumeBody interface{}) diag.Diagnostics {
	accessKeyID := utils.PathSearch("credentials.access_key_id", assumeBody, nil).(string)
	if accessKeyID == "" {
		return diag.Errorf("unable to find the access key from the agency credential response")
	}

	assumedAgency := []interface{}{
		map[string]interface{}{
			"urn": utils.PathSearch("assumed_agency.urn", assumeBody, nil),
			"id":  utils.PathSearch("assumed_agency.id", assumeBody, nil),
		},
	}
	credentials := []interface{}{
		map[string]interface{}{
			"access_key_id":     accessKeyID,
			"secret_access_key": utils.PathSearch("credentials.secret_access_key", assumeBody, nil),
			"security_token":    utils.PathSearch("credentials.security_token", assumeBody, nil),
			"expiration":        utils.PathSearch("credentials.expiration", assumeBody, nil),
		},
	}

	mErr := multierror.Append(
		d.Set("assumed_agency", assumedAgency),
		d.Set("credentials", credentials),
		d.Set("audience", utils.PathSearch("audience", assumeBody, nil)),
		d.Set("idp_urn", utils.PathSearch("provider", assumeBody, nil)),
		d.Set("source_identity", utils.PathSearch("source_identity", assumeBody, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting agency credential fields: %s", err)
	}

	return nil
}
