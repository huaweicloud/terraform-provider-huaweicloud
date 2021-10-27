package iam

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/mappings"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/metadatas"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/oidcconfig"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/protocols"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/providers"

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

const (
	protocolSAML = "saml"
	protocolOIDC = "oidc"
	scopeSpilt   = " "
)

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
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w-]{1,64}$`),
					"The maximum length is 64 characters. "+
						"Only letters, digits, underscores (_), and hyphens (-) are allowed"),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{protocolSAML, protocolOIDC}, false),
			},
			"conversion_rules": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
									"group": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"remote": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 10,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute": {
										Type:     schema.TypeString,
										Required: true,
									},
									"condition": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"any_one_of", "not_any_of"}, false),
									},
									"value": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "enabled",
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"sso_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "virtual_user_sso",
				ValidateFunc: validation.StringInSlice([]string{"virtual_user_sso", "iam_user_sso"}, false),
			},
			"access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"program", "program_console"}, false),
			},
			"provider_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"signing_key": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"metadata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
		return fmtp.DiagErrorf("Error creating HuaweiCloud IAM without version client: %s", err)
	}
	// Create a SAML protocol provider.
	opts := providers.CreateProviderOpts{
		Enabled:     getEnabledVal(d),
		SsoType:     d.Get("sso_type").(string),
		Description: d.Get("description").(string),
	}
	name := d.Get("name").(string)
	logp.Printf("[DEBUG] Create identity options %s : %#v", name, opts)
	provider, err := providers.Create(client, name, opts)
	if err != nil {
		logp.Printf("[ERROR] Failed to create identity provider: %s", err)
		return diag.FromErr(err)
	}
	d.SetId(provider.ID)

	protocol := d.Get("protocol").(string)
	// Create protocol and mapping
	err = createMappingAndProtocol(client, d)
	// Create the provider access type, it only worked on the OIDC protocol.
	if err == nil && protocol == protocolOIDC {
		err = createAccessType(conf, d)
	}
	// If create mapping or protocol or access type fails, delete the provider.
	if err != nil {
		logp.Printf("[ERROR] Failed to create mapping or protocol or access type, delete the provider: %s,", err)
		resourceIdentityProviderDelete(ctx, d, meta)
		return diag.FromErr(err)
	}

	// Import metadata, metadata only worked on saml protocol provider
	if protocol == protocolSAML {
		err = importMetadata(conf, d)
		if err != nil {
			resourceIdentityProviderDelete(ctx, d, meta)
			return diag.FromErr(err)
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
		return fmtp.Errorf("Error creating HuaweiCloud IAM client without version number: %s", err)
	}

	providerID := d.Get("name").(string)
	opts := metadatas.ImportOpts{
		DomainID: conf.DomainID,
		Metadata: metadata,
	}
	_, err = metadatas.Import(client, providerID, protocolSAML, opts)
	if err != nil {
		return fmtp.Errorf("failed to import metadata: %s", err)
	}
	return nil
}

func getEnabledVal(d *schema.ResourceData) bool {
	return d.Get("status").(string) == "enabled"
}

func getScopes(d *schema.ResourceData) string {
	scopes := utils.ExpandToStringList(d.Get("scopes").([]interface{}))
	str := strings.Join(scopes, scopeSpilt)
	return str
}

func createAccessType(conf *config.Config, d *schema.ResourceData) error {
	// access_type is not required, if access_type is empty do not set it.
	if _, ok := d.GetOk("access_type"); !ok {
		return nil
	}

	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IAM client without version number: %s", err)
	}

	opts := oidcconfig.CreateOpts{
		AccessMode:            d.Get("access_type").(string),
		IdpURL:                d.Get("provider_url").(string),
		ClientID:              d.Get("client_id").(string),
		AuthorizationEndpoint: d.Get("authorization_endpoint").(string),
		Scope:                 getScopes(d),
		ResponseType:          d.Get("response_type").(string),
		ResponseMode:          d.Get("response_mode").(string),
		SigningKey:            d.Get("signing_key").(string),
	}
	logp.Printf("[DEBUG] Create access type of provider: %#v", opts)
	providerID := d.Get("name").(string)
	_, err = oidcconfig.Create(client, providerID, opts)
	if err != nil {
		return fmtp.Errorf("failed to create access type: %s", err)
	}
	return nil
}

func buildMappingRules(conversionRules []interface{}) mappings.MappingOption {
	rules := make([]mappings.MappingRule, 0, len(conversionRules))

	for _, cr := range conversionRules {
		convRule := cr.(map[string]interface{})

		// build local rules
		local := convRule["local"].([]interface{})
		localRules := make([]mappings.LocalRule, 0, len(local))
		for _, v := range local {
			lRule := v.(map[string]interface{})
			r := mappings.LocalRule{
				User: mappings.LocalRuleVal{
					Name: lRule["username"].(string),
				},
			}
			group, ok := lRule["group"]
			if ok && len(group.(string)) > 0 {
				r.Group = &mappings.LocalRuleVal{
					Name: group.(string),
				}
			}
			localRules = append(localRules, r)
		}
		// build remote rule
		remote := convRule["remote"].([]interface{})
		remoteRules := make([]mappings.RemoteRule, 0, len(remote))
		for _, v := range remote {
			rRule := v.(map[string]interface{})
			r := mappings.RemoteRule{
				Type: rRule["attribute"].(string),
			}
			if condition, ok := rRule["condition"]; ok {
				values := utils.ExpandToStringList(rRule["value"].([]interface{}))
				if condition.(string) == "any_one_of" {
					r.AnyOneOf = values
				} else if condition.(string) == "not_any_of" {
					r.NotAnyOf = values
				}
			}
			remoteRules = append(remoteRules, r)
		}

		rule := mappings.MappingRule{
			Local:  localRules,
			Remote: remoteRules,
		}
		rules = append(rules, rule)
	}
	opts := mappings.MappingOption{
		Rules: rules,
	}
	return opts
}

// createMappingAndProtocol create mapping and protocol
func createMappingAndProtocol(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	providerID := d.Get("name").(string)
	mappingID := "mapping_" + providerID

	conversionRules := d.Get("conversion_rules").([]interface{})
	if len(conversionRules) == 0 {
		return fmtp.Errorf("conversion_rules is required for identity provider")
	}

	// Check if the mappingID exists, update if it exists, otherwise create it.
	r, err := mappings.List(client).AllPages()
	if err != nil {
		return fmtp.Errorf("error in querying mapping: %s", err)
	}
	mapArr, err := mappings.ExtractMappings(r)
	mErr := multierror.Append(nil, err)
	filterData, err := utils.FilterSliceWithField(mapArr, map[string]interface{}{
		"ID": mappingID,
	})
	mErr = multierror.Append(mErr, err)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("error in querying mapping: %s", mErr)
	}

	mappingOpts := buildMappingRules(conversionRules)
	// Create the mapping if it does not exist, otherwise update it.
	if len(filterData) == 0 {
		_, err = mappings.Create(client, mappingID, mappingOpts)
	} else {
		_, err = mappings.Update(client, mappingID, mappingOpts)
	}
	if err != nil {
		return fmtp.Errorf("error in creating/updating mapping: %s", err)
	}

	protocolName := d.Get("protocol").(string)
	// Create protocol for the provider
	_, err = protocols.Create(client, providerID, protocolName, mappingID)
	if err != nil {
		// If fails to create protocols, then delete the mapping.
		mErr = multierror.Append(
			mErr,
			err,
			mappings.Delete(client, mappingID),
		)
		logp.Printf("[ERROR] Failed to create protocol, so delete the mapping that has been created.\n%s", mErr)
		return fmtp.Errorf("error in creating provider protocol: %s", mErr.Error())
	}
	return nil
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

func resourceIdentityProviderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud IAM client without version number: %s", err)
	}

	provider, err := providers.Get(client, d.Id())
	if err != nil {
		logp.Printf("[ERROR] Error in obtaining identity provider: %s", err)
		return diag.FromErr(err)
	}

	// Query the protocol name from HuaweiCloud.
	protocol := queryProtocolName(client, d)
	url := generateLoginLink(conf.Cloud, conf.DomainID, provider.ID, protocol)

	status := "disabled"
	if provider.Enabled {
		status = "enabled"
	}
	mErr := multierror.Append(err,
		d.Set("name", provider.ID),
		d.Set("protocol", protocol),
		d.Set("sso_type", provider.SsoType),
		d.Set("status", status),
		d.Set("login_link", url),
		getAndSetConversionRulesAttrs(client, d),
	)
	d.Set("description", provider.Description)
	// Get and set access type to attribute.
	getAndSetAccessTypeAttrs(client, d)
	// Get and set metadata to attribute
	getAndSetMetadataAttr(client, d)

	if err = mErr.ErrorOrNil(); err != nil {
		logp.Printf("[ERROR] Error in setting identity provider attributes %s: %s", d.Id(), err)
		return diag.FromErr(err)
	}
	return nil
}

func getAndSetMetadataAttr(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	providerID := d.Get("name").(string)
	r, err := metadatas.Get(client, providerID, protocolSAML)
	if err != nil {
		if errors.As(err, &golangsdk.ErrDefault404{}) {
			return nil
		}
		return err
	}
	d.Set("metadata", r.Data)
	return nil
}

func getAndSetAccessTypeAttrs(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	// The access type only work on OIDC protocol.
	if protocol, ok := d.GetOk("protocol"); ok && protocol == protocolSAML {
		// Because response_mode and response_type has default values, so set default values for them.
		d.Set("response_mode", "form_post")
		d.Set("response_type", "id_token")
		return nil
	}

	providerID := d.Get("name").(string)
	accType, err := oidcconfig.Get(client, providerID)
	if err != nil {
		if errors.As(err, &golangsdk.ErrDefault404{}) {
			return nil
		}
		return err
	}
	mErr := multierror.Append(
		d.Set("access_type", accType.AccessMode),
		d.Set("provider_url", accType.IdpURL),
		d.Set("client_id", accType.ClientID),
		d.Set("signing_key", accType.SigningKey),
	)
	scopes := strings.Split(accType.Scope, scopeSpilt)
	d.Set("scopes", scopes)
	d.Set("response_mode", accType.ResponseMode)
	d.Set("authorization_endpoint", accType.AuthorizationEndpoint)
	d.Set("response_type", accType.ResponseType)

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("failed to set access type attributes: %s", err)
	}
	return nil
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

func getAndSetConversionRulesAttrs(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	mappingID := "mapping_" + d.Id()
	mapping, err := mappings.Get(client, mappingID)
	if err != nil {
		if errors.As(err, &golangsdk.ErrDefault404{}) {
			return fmtp.Errorf("The conversion rule is required, but not found, " +
				"please go to the console to verify the data.")
		}
		return fmtp.Errorf("error in querying conversion rules: %s", err)
	}

	conversionRules := make([]interface{}, 0, len(mapping.Rules))
	for _, v := range mapping.Rules {
		localRules := make([]map[string]interface{}, 0, len(v.Local))
		for _, localRule := range v.Local {
			r := map[string]interface{}{
				"username": localRule.User.Name,
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

	err = d.Set("conversion_rules", conversionRules)
	return err
}

func resourceIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud IAM client without version number: %s", err)
	}

	mErr := &multierror.Error{}
	if d.HasChanges("status", "description") {
		status := getEnabledVal(d)
		description := d.Get("description").(string)
		opts := providers.UpdateOpts{
			Enabled:     &status,
			Description: &description,
		}
		logp.Printf("[DEBUG] Update identity options %s : %#v", d.Id(), opts)

		_, err = providers.Update(client, d.Id(), opts)
		if err != nil {
			e := fmtp.Errorf("Failed to update identity provider: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if d.HasChange("metadata") {
		err = importMetadata(conf, d)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("conversion_rules") {
		conversionRules := d.Get("conversion_rules").([]interface{})
		mappingOpts := buildMappingRules(conversionRules)
		mappingID := "mapping_" + d.Id()
		_, err = mappings.Update(client, mappingID, mappingOpts)
		if err != nil {
			e := fmtp.Errorf("Failed to update the mapping of provider: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if d.Get("protocol").(string) == protocolOIDC && d.HasChanges("access_type", "provider_url", "client_id",
		"authorization_endpoint", "scopes", "response_type", "response_mode", "signing_key") {
		err = updateAccessType(client, d)
		if err != nil {
			e := fmtp.Errorf("Error in updating the access type of provider: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error in updating provider: %s", err)
	}

	return resourceIdentityProviderRead(ctx, d, meta)
}

func updateAccessType(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	opts := oidcconfig.UpdateOpenIDConnectConfigOpts{
		AccessMode:            d.Get("access_type").(string),
		IdpURL:                d.Get("provider_url").(string),
		ClientID:              d.Get("client_id").(string),
		AuthorizationEndpoint: d.Get("authorization_endpoint").(string),
		Scope:                 getScopes(d),
		ResponseType:          d.Get("response_type").(string),
		ResponseMode:          d.Get("response_mode").(string),
		SigningKey:            d.Get("signing_key").(string),
	}
	logp.Printf("[DEBUG] Update access type of provider: %#v", opts)
	providerID := d.Get("name").(string)
	_, err := oidcconfig.Update(client, providerID, opts)
	if err != nil {
		return fmtp.Errorf("failed to update access type: %s", err)
	}
	return nil
}

func resourceIdentityProviderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud IAM client without version number: %s", err)
	}

	err = providers.Delete(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error deleting HuaweiCloud identity provider")
	}
	d.SetId("")
	return nil
}
