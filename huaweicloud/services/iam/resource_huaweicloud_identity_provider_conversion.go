package iam

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/mappings"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/providers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM PUT /v3/OS-FEDERATION/mappings/{id}
// @API IAM PATCH /v3/OS-FEDERATION/mappings/{id}
// @API IAM GET /v3/OS-FEDERATION/mappings
// @API IAM GET /v3/OS-FEDERATION/mappings/{id}
// @API IAM GET /v3/OS-FEDERATION/identity_providers/{id}
func ResourceV3ProviderConversion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3ProviderConversionCreate,
		ReadContext:   resourceV3ProviderConversionRead,
		UpdateContext: resourceV3ProviderConversionUpdate,
		DeleteContext: resourceV3ProviderConversionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the identity provider used to manage the conversion rules.`,
			},
			"conversion_rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `The identity conversion rules of the identity provider.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `The federated user information on the cloud platform.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The name of a federated user on the cloud platform.`,
									},
									"group": {
										Type:     schema.TypeString,
										Optional: true,
										Description: `The user group to which the federated user belongs on the
                                               cloud platform.`,
									},
								},
							},
						},
						"remote": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `Federated user information in the IDP system.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The attribute in the IDP assertion.`,
									},
									"condition": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"any_one_of", "not_any_of"}, false),
										Description:  `The condition of conversion rule.`,
									},
									"value": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Description: `The rule is matched only if the specified strings appear in the
                                            attribute type.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildV3ProviderConversionRules(conversionRules []interface{}) mappings.MappingOption {
	rules := make([]mappings.MappingRule, 0, len(conversionRules))

	for _, cr := range conversionRules {
		convRule := cr.(map[string]interface{})

		// build local rules
		local := convRule["local"].([]interface{})
		localRules := make([]mappings.LocalRule, 0, len(local))
		for _, v := range local {
			lRule := v.(map[string]interface{})
			var r mappings.LocalRule
			username, ok := lRule["username"]
			if ok && len(username.(string)) > 0 {
				r.User = &mappings.LocalRuleVal{
					Name: lRule["username"].(string),
				}
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

func resourceV3ProviderConversionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}
	providerID := d.Get("provider_id").(string)
	mappingID := generateV3ProviderMappingID(providerID)

	// Check if the mappingID exists, update if it exists, otherwise create it.
	r, err := mappings.List(client).AllPages()
	err404 := golangsdk.ErrDefault404{}
	if err != nil && !errors.As(err, &err404) {
		return diag.Errorf("error in querying or extract conversions: %s", err)
	}

	conversions, err := mappings.ExtractMappings(r)
	if err != nil {
		return diag.Errorf("error in extracting provider conversions: %s", err)
	}

	filterData, err := utils.FilterSliceWithField(conversions, map[string]interface{}{
		"ID": mappingID,
	})
	if err != nil {
		return diag.Errorf("error in filtering conversions: %s", err)
	}

	conversionRules := d.Get("conversion_rules").([]interface{})
	mappingOpts := buildV3ProviderConversionRules(conversionRules)
	// Create the mapping if it does not exist, otherwise update it.
	if len(filterData) == 0 {
		_, err = mappings.Create(client, mappingID, mappingOpts)
	} else {
		_, err = mappings.Update(client, mappingID, mappingOpts)
	}
	if err != nil {
		return diag.Errorf("error in creating/updating mapping: %s", err)
	}

	d.SetId(mappingID)
	return resourceV3ProviderConversionRead(ctx, d, meta)
}

func resourceV3ProviderConversionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}

	conversionID := d.Id()
	conversions, err := mappings.Get(client, conversionID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error in querying conversion rules")
	}

	providerID := strings.ReplaceAll(conversionID, "mapping_", "")
	mErr := multierror.Append(
		d.Set("provider_id", providerID),
		d.Set("conversion_rules", flattenV3ProviderConversionRules(conversions)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting identity provider conversion rules: %s", mErr)
	}
	return nil
}

func resourceV3ProviderConversionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}

	conversionRules := d.Get("conversion_rules").([]interface{})
	conversionRuleOpts := buildV3ProviderConversionRules(conversionRules)
	conversionID := d.Id()
	_, err = mappings.Update(client, conversionID, conversionRuleOpts)
	if err != nil {
		return diag.Errorf("failed to update the provider conversion rules: %s", err)
	}

	return resourceV3ProviderConversionRead(ctx, d, meta)
}

func resourceV3ProviderConversionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}

	providerID := d.Get("provider_id").(string)
	_, err = providers.Get(client, providerID)
	if err != nil && errors.As(err, &golangsdk.ErrDefault404{}) {
		d.SetId("")
		return nil
	}

	conversionID := d.Id()
	opts := getDefaultV3ProviderConversionOpts()
	_, err = mappings.Update(client, conversionID, *opts)

	if err != nil {
		return diag.Errorf("error resetting provider conversion rules to default value" +
			"(the conversion rules can not be deleted, it can be reset to default value).")
	}

	return nil
}
