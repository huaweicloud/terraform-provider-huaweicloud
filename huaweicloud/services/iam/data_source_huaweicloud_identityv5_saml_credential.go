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

// @API IAM POST /v5/agencies/assume-with-saml
func DataSourceV5SamlCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5SamlCredentialRead,

		Schema: map[string]*schema.Schema{
			"saml_response": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				Description: `Specifies the SAML authentication response (assertion) issued by the SAML identity provider, ` +
					`which is a base64 encoded XML string.`,
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
				Description: `Specifies the uniform resource name of the SAML identity provider, in the format of ` +
					`iam::{account_id}:samlProvider:{provider_name}.`,
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
			"assumed_agency_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the assumed agency session.`,
			},
			"assumed_agency_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the assumed agency session.`,
			},
			"subject": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the NameID element of the Subject element in the SAML assertion.`,
			},
			"subject_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The format of the NameID, defined by the Format attribute of the NameID element in the SAML assertion.`,
			},
			"issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the Issuer element in the SAML assertion.`,
			},
			"audience": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the Recipient attribute of the SubjectConfirmationData element in the SAML assertion.`,
			},
			"name_qualifier": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The hash value of the issuer, the account ID and the name of the SAML identity provider.`,
			},
			"idp_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the SAML identity provider.`,
			},
			"source_identity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The identity declared by the SourceIdentity attribute in the SAML assertion.`,
			},
		},
	}
}

func dataSourceV5SamlCredentialRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("sts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating STS client: %s", err)
	}

	assumeHttpUrl := "v5/agencies/assume-with-saml"
	assumePath := client.Endpoint + assumeHttpUrl
	assumeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildV5AgencyAssumeWithSamlBodyParams(d),
	}
	assumeResp, err := client.Request("POST", assumePath, &assumeOpt)
	if err != nil {
		return diag.Errorf("error assuming agency with SAML: %s", err)
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

	return setV5SamlCredentialSessionAttributes(d, assumeBody)
}

func buildV5AgencyAssumeWithSamlBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"saml_response":       d.Get("saml_response").(string),
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

func setV5SamlCredentialSessionAttributes(d *schema.ResourceData, assumeBody interface{}) diag.Diagnostics {
	mErr := multierror.Append(
		d.Set("subject", utils.PathSearch("subject", assumeBody, nil)),
		d.Set("subject_type", utils.PathSearch("subject_type", assumeBody, nil)),
		d.Set("issuer", utils.PathSearch("issuer", assumeBody, nil)),
		d.Set("audience", utils.PathSearch("audience", assumeBody, nil)),
		d.Set("name_qualifier", utils.PathSearch("name_qualifier", assumeBody, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SAML agency credential session fields: %s", err)
	}
	return nil
}
