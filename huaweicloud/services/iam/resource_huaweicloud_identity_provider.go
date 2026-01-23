package iam

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/mappings"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/metadatas"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/oidcconfig"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/protocols"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/providers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	protocolSAML = "saml"
	protocolOIDC = "oidc"

	scopeSpilt = " "
)

// @API IAM PUT /v3/OS-FEDERATION/identity_providers/{id}
// @API IAM PATCH /v3/OS-FEDERATION/identity_providers/{id}
// @API IAM GET /v3/OS-FEDERATION/identity_providers/{id}
// @API IAM DELETE /v3/OS-FEDERATION/identity_providers/{id}
// @API IAM GET /v3.0/OS-FEDERATION/identity-providers/{idp_id}/openid-connect-config
// @API IAM POST /v3.0/OS-FEDERATION/identity-providers/{idp_id}/openid-connect-config
// @API IAM PUT /v3.0/OS-FEDERATION/identity-providers/{idp_id}/openid-connect-config
// @API IAM PUT /v3/OS-FEDERATION/mappings/{id}
// @API IAM GET /v3/OS-FEDERATION/mappings/{id}
// @API IAM DELETE /v3/OS-FEDERATION/mappings/{id}
// @API IAM GET /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols
// @API IAM PUT /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}
// @API IAM GET /v3-ext/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}/metadata
// @API IAM POST /v3-ext/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}/metadata
func ResourceIdentityProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProviderCreate,
		ReadContext:   resourceIdentityProviderRead,
		UpdateContext: resourceIdentityProviderUpdate,
		DeleteContext: resourceIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"sso_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"metadata": {
				Type:      schema.TypeString,
				Optional:  true,
				StateFunc: utils.HashAndHexEncode,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Computed: true,
						},
						"scopes": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
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
				},
			},
			"conversion_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"remote": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"condition": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"login_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}

	// Create a SAML protocol provider.
	opts := providers.CreateProviderOpts{
		Enabled:     d.Get("status").(bool),
		Description: d.Get("description").(string),
		SsoType:     d.Get("sso_type").(string),
	}
	name := d.Get("name").(string)
	log.Printf("[DEBUG] Create identity options %s : %#v", name, opts)
	provider, err := providers.Create(client, name, opts)
	if err != nil {
		return diag.Errorf("error creating identity provider: %s", err)
	}

	d.SetId(provider.ID)

	if d.HasChange("protocol") {
		// Create default mapping and protocol
		err = createProtocol(client, d)
		if err != nil {
			return diag.Errorf("error creating provider protocol: %s", err)
		}

		// Import metadata, metadata only worked on saml protocol providers
		protocol := d.Get("protocol").(string)
		if protocol == protocolSAML {
			err = importMetadata(conf, d)
			if err != nil {
				return diag.Errorf("error importing matedata into identity provider: %s", err)
			}
		} else if ac, ok := d.GetOk("access_config"); ok {
			// Create access config for oidc provider.
			accessConfigArr := ac.([]interface{})
			accessConfig := accessConfigArr[0].(map[string]interface{})

			accessType := accessConfig["access_type"].(string)
			createAccessTypeOpts := oidcconfig.CreateOpts{
				AccessMode: accessType,
				IdpURL:     accessConfig["provider_url"].(string),
				ClientID:   accessConfig["client_id"].(string),
				SigningKey: accessConfig["signing_key"].(string),
			}

			if accessType == "program_console" {
				scopes := utils.ExpandToStringList(accessConfig["scopes"].([]interface{}))
				createAccessTypeOpts.Scope = strings.Join(scopes, scopeSpilt)
				createAccessTypeOpts.AuthorizationEndpoint = accessConfig["authorization_endpoint"].(string)
				createAccessTypeOpts.ResponseType = accessConfig["response_type"].(string)
				createAccessTypeOpts.ResponseMode = accessConfig["response_mode"].(string)
			}

			log.Printf("[DEBUG] Create access type of provider: %#v", opts)
			_, err = oidcconfig.Create(client, provider.ID, createAccessTypeOpts)
			if err != nil {
				return diag.Errorf("error creating the provider access config: %s", err)
			}
		}
	}

	return resourceIdentityProviderRead(ctx, d, meta)
}

// importMetadata import metadata to provider, overwrite if it exists.
func importMetadata(conf *config.Config, d *schema.ResourceData) error {
	metadata := d.Get("metadata").(string)
	if len(metadata) == 0 {
		return nil
	}

	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating IAM client without version: %s", err)
	}

	providerID := d.Get("name").(string)
	opts := metadatas.ImportOpts{
		DomainID: conf.DomainID,
		Metadata: metadata,
	}
	_, err = metadatas.Import(client, providerID, protocolSAML, opts)
	if err != nil {
		return fmt.Errorf("failed to import metadata: %s", err)
	}
	return nil
}

// createProtocol create default mapping and protocol
func createProtocol(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	providerID := d.Get("name").(string)

	// Create default mapping
	defaultConversionRules := getDefaultConversionOpts()
	mappingID := generateMappingID(providerID)
	_, err := mappings.Create(client, mappingID, *defaultConversionRules)
	if err != nil {
		return fmt.Errorf("error in creating default conversion rule: %s", err)
	}

	// Create protocol
	protocolName := d.Get("protocol").(string)
	_, err = protocols.Create(client, providerID, protocolName, mappingID)
	if err != nil {
		// If fails to create protocols, then delete the mapping.
		log.Printf("[WARN] error creating protocol and the mapping will be deleted. Error: %s", err)
		mErr := multierror.Append(err,
			mappings.Delete(client, mappingID),
		)
		return fmt.Errorf("error creating identity provider protocol: %s", mErr.Error())
	}
	return nil
}

func getDefaultConversionOpts() *mappings.MappingOption {
	localRules := []mappings.LocalRule{
		{
			User: &mappings.LocalRuleVal{
				Name: "FederationUser",
			},
		},
	}
	remoteRules := []mappings.RemoteRule{
		{
			Type: "__NAMEID__",
		},
	}

	opts := mappings.MappingOption{
		Rules: []mappings.MappingRule{
			{
				Local:  localRules,
				Remote: remoteRules,
			},
		},
	}
	return &opts
}

func resourceIdentityProviderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}

	providerID := d.Id()
	provider, err := providers.Get(client, providerID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error obtaining identity provider")
	}

	// Query the protocol name from HuaweiCloud.
	protocol := queryProtocolName(client, d)
	url := generateLoginLink(conf.Cloud, conf.DomainID, provider.ID, protocol)

	mErr := multierror.Append(nil,
		d.Set("name", provider.ID),
		d.Set("protocol", protocol),
		d.Set("sso_type", provider.SsoType),
		d.Set("status", provider.Enabled),
		d.Set("login_link", url),
		d.Set("description", provider.Description),
	)

	// Query and set conversion rules
	mappingID := generateMappingID(providerID)
	conversions, err := mappings.Get(client, mappingID)
	if err == nil {
		conversionRules := flattenConversionRulesAttr(conversions)
		err = d.Set("conversion_rules", conversionRules)
		mErr = multierror.Append(mErr, err)
	}

	// Query and set metadata of the protocol SAML provider
	if protocol == protocolSAML {
		r, err := metadatas.Get(client, providerID, protocolSAML)
		if err == nil {
			err = d.Set("metadata", utils.HashAndHexEncode(r.Data))
			mErr = multierror.Append(mErr, err)
		}
	}

	// Query and set access type of the protocol OIDC provider
	if protocol == protocolOIDC {
		accessType, err := oidcconfig.Get(client, providerID)
		if err == nil {
			scopes := strings.Split(accessType.Scope, scopeSpilt)
			accessTypeConfig := []interface{}{
				map[string]interface{}{
					"access_type":            accessType.AccessMode,
					"provider_url":           accessType.IdpURL,
					"client_id":              accessType.ClientID,
					"signing_key":            accessType.SigningKey,
					"scopes":                 scopes,
					"response_mode":          accessType.ResponseMode,
					"authorization_endpoint": accessType.AuthorizationEndpoint,
					"response_type":          accessType.ResponseType,
				},
			}

			mErr = multierror.Append(mErr,
				d.Set("access_config", accessTypeConfig),
			)
		}
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting identity provider attributes: %s", err)
	}
	return nil
}

func flattenConversionRulesAttr(conversions *mappings.IdentityMapping) []interface{} {
	conversionRules := make([]interface{}, 0, len(conversions.Rules))
	for _, v := range conversions.Rules {
		localRules := make([]map[string]interface{}, 0, len(v.Local))
		for _, localRule := range v.Local {
			r := map[string]interface{}{}
			if localRule.User != nil {
				r["username"] = localRule.User.Name
			}

			if localRule.Group != nil {
				r["group"] = localRule.Group.Name
			}
			localRules = append(localRules, r)
		}

		remoteRules := make([]map[string]interface{}, 0, len(v.Remote))
		for _, remoteRule := range v.Remote {
			r := map[string]interface{}{
				"attribute": remoteRule.Type,
			}
			if len(remoteRule.NotAnyOf) > 0 {
				r["condition"] = "not_any_of"
				r["value"] = remoteRule.NotAnyOf
			} else if len(remoteRule.AnyOneOf) > 0 {
				r["condition"] = "any_one_of"
				r["value"] = remoteRule.AnyOneOf
			}
			remoteRules = append(remoteRules, r)
		}

		rule := map[string]interface{}{
			"local":  localRules,
			"remote": remoteRules,
		}
		conversionRules = append(conversionRules, rule)
	}
	return conversionRules
}

func generateMappingID(providerID string) string {
	return "mapping_" + providerID
}

// generateLoginLink generate login link base on config.domainID.
func generateLoginLink(host, domainID, id, protocol string) string {
	// The domain name is the same as that of the console, it is converted according to the config.Cloud.
	if host == "myhuaweicloud.com" {
		host = "huaweicloud.com"
	}
	url := fmt.Sprintf("https://auth.%s/authui/federation/websso?domain_id=%s&idp=%s&protocol=%s",
		host, domainID, id, protocol)
	return url
}

func queryProtocolName(client *golangsdk.ServiceClient, d *schema.ResourceData) string {
	arr, err := protocols.List(client, d.Id())
	if err != nil {
		return ""
	}
	// The SAML protocol provider may not have protocol data,
	// so the default value is set to SAML.
	protocolName := protocolSAML
	if len(arr) > 0 {
		protocolName = arr[0].ID
	}
	return protocolName
}

func resourceIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}

	mErr := &multierror.Error{}
	if d.HasChanges("status", "description") {
		status := d.Get("status").(bool)
		description := d.Get("description").(string)
		opts := providers.UpdateOpts{
			Enabled:     &status,
			Description: &description,
		}

		log.Printf("[DEBUG] Update identity options %s : %#v", d.Id(), opts)
		_, err = providers.Update(client, d.Id(), opts)
		if err != nil {
			e := fmt.Errorf("failed to update identity provider: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if d.HasChange("metadata") {
		err = importMetadata(conf, d)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("access_config") && d.Get("protocol") == protocolOIDC {
		err = updateAccessConfig(client, d)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error updating provider: %s", err)
	}

	return resourceIdentityProviderRead(ctx, d, meta)
}

func updateAccessConfig(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	accessConfigArr := d.Get("access_config").([]interface{})
	if len(accessConfigArr) == 0 {
		return fmt.Errorf("the access_config block is required for the OIDC provider")
	}

	accessConfig := accessConfigArr[0].(map[string]interface{})
	accessType := accessConfig["access_type"].(string)
	opts := oidcconfig.UpdateOpenIDConnectConfigOpts{
		AccessMode: accessType,
		IdpURL:     accessConfig["provider_url"].(string),
		ClientID:   accessConfig["client_id"].(string),
		SigningKey: accessConfig["signing_key"].(string),
	}

	if accessType == "program_console" {
		scopes := utils.ExpandToStringList(accessConfig["scopes"].([]interface{}))
		opts.Scope = strings.Join(scopes, scopeSpilt)
		opts.AuthorizationEndpoint = accessConfig["authorization_endpoint"].(string)
		opts.ResponseType = accessConfig["response_type"].(string)
		opts.ResponseMode = accessConfig["response_mode"].(string)
	}
	log.Printf("[DEBUG] Update access type of provider: %#v", opts)
	providerID := d.Id()
	_, err := oidcconfig.Update(client, providerID, opts)
	return err
}

func resourceIdentityProviderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}

	err = providers.Delete(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IAM provider")
	}

	return nil
}
