package iam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

var mappingNonUpdatableParams = []string{"provider_id"}

// @API IAM PUT /v3/OS-FEDERATION/mappings/{id}
// @API IAM PATCH /v3/OS-FEDERATION/mappings/{id}
// @API IAM GET /v3/OS-FEDERATION/mappings
// @API IAM GET /v3/OS-FEDERATION/mappings/{id}
// @API IAM GET /v3/OS-FEDERATION/identity_providers/{id}
func ResourceIAMProviderMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIAMProviderMappingCreate,
		ReadContext:   resourceIAMProviderMappingRead,
		UpdateContext: resourceIAMProviderMappingUpdate,
		DeleteContext: resourceIAMProviderMappingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(mappingNonUpdatableParams),
		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mapping_rules": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(o, n)
					return equal
				},
			},
		},
	}
}

func resourceIAMProviderMappingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}
	providerID := d.Get("provider_id").(string)
	mappingID := generateMappingID(providerID)

	// Check if the mappingID exists, update if it exists, otherwise create it.
	r, err := mappings.List(client).AllPages()
	err404 := golangsdk.ErrDefault404{}
	if err != nil && !errors.As(err, &err404) {
		return diag.Errorf("error in querying or extract mappings: %s", err)
	}

	providerMappings, err := mappings.ExtractMappings(r)
	if err != nil {
		return diag.Errorf("error in extracting provider mappings: %s", err)
	}

	filterData, err := utils.FilterSliceWithField(providerMappings, map[string]interface{}{
		"ID": mappingID,
	})
	if err != nil {
		return diag.Errorf("error in filtering mappings: %s", err)
	}

	mappingRules := d.Get("mapping_rules").(string)
	mappingOpts, err := buildMappingOpts(mappingRules)
	if err != nil {
		return diag.FromErr(err)
	}
	// Create the mapping if it does not exist, otherwise update it.
	if len(filterData) == 0 {
		err = createMapping(client, mappingID, mappingOpts)
	} else {
		err = updateMapping(client, mappingID, mappingOpts)
	}
	if err != nil {
		return diag.Errorf("error in creating/updating mapping: %s", err)
	}

	d.SetId(mappingID)
	return resourceIAMProviderMappingRead(ctx, d, meta)
}

func createMapping(client *golangsdk.ServiceClient, mappingID string, mappingOpts interface{}) error {
	createMappingHttpUrl := "v3/OS-FEDERATION/mappings/{id}"
	createMappingPath := client.Endpoint + createMappingHttpUrl
	createMappingPath = strings.ReplaceAll(createMappingPath, "{id}", mappingID)

	createMappingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         mappingOpts,
	}

	_, err := client.Request("PUT", createMappingPath, &createMappingOpt)
	if err != nil {
		return err
	}

	return nil
}

func updateMapping(client *golangsdk.ServiceClient, mappingID string, mappingOpts interface{}) error {
	updateMappingHttpUrl := "v3/OS-FEDERATION/mappings/{id}"
	updateMappingPath := client.Endpoint + updateMappingHttpUrl
	updateMappingPath = strings.ReplaceAll(updateMappingPath, "{id}", mappingID)

	updateMappingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         mappingOpts,
	}

	_, err := client.Request("PATCH", updateMappingPath, &updateMappingOpt)
	if err != nil {
		return err
	}

	return nil
}

func buildMappingOpts(mappingRules string) (interface{}, error) {
	var rules interface{}
	err := json.Unmarshal([]byte(mappingRules), &rules)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling rules, please check the format of the mapping rules: %s", err)
	}

	res := map[string]interface{}{
		"mapping": map[string]interface{}{
			"rules": rules,
		},
	}
	return res, nil
}

func resourceIAMProviderMappingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	mappingID := d.Id()

	client, err := cfg.IAMNoVersionClient(region)
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}

	getMappingHttpUrl := "v3/OS-FEDERATION/mappings/{id}"
	getMappingPath := client.Endpoint + getMappingHttpUrl
	getMappingPath = strings.ReplaceAll(getMappingPath, "{id}", mappingID)
	getMappingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getMappingResp, err := client.Request("GET", getMappingPath, &getMappingOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving identity provider mapping")
	}

	getMappingRespBody, err := utils.FlattenResponse(getMappingResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mapping := utils.PathSearch("mapping", getMappingRespBody, nil)
	if mapping == nil {
		return diag.Errorf("error getting identity provider mapping: mapping is not found in API response")
	}

	rules := utils.PathSearch("rules", mapping, nil)

	mappingRules, err := json.Marshal(rules)
	if err != nil {
		return diag.Errorf("error marshaling rules: %s", err)
	}

	providerID := strings.ReplaceAll(mappingID, "mapping_", "")
	mErr := multierror.Append(
		d.Set("provider_id", providerID),
		d.Set("mapping_rules", string(mappingRules)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting identity provider mapping rules: %s", mErr)
	}
	return nil
}

func resourceIAMProviderMappingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}

	mappingRules := d.Get("mapping_rules").(string)
	mappingRuleOpts, err := buildMappingOpts(mappingRules)
	if err != nil {
		return diag.FromErr(err)
	}

	mappingID := d.Id()
	err = updateMapping(client, mappingID, mappingRuleOpts)
	if err != nil {
		return diag.Errorf("failed to update the provider mapping rules: %s", err)
	}

	return resourceIAMProviderMappingRead(ctx, d, meta)
}

func resourceIAMProviderMappingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client without version: %s", err)
	}

	providerID := d.Get("provider_id").(string)
	_, err = providers.Get(client, providerID)
	if err != nil && errors.As(err, &golangsdk.ErrDefault404{}) {
		d.SetId("")
		return nil
	}

	mappingID := d.Id()
	opts, err := buildMappingOpts(defaultMapping)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateMapping(client, mappingID, opts)
	if err != nil {
		return diag.Errorf("error resetting provider mapping rules to default value" +
			"(the mapping rules can not be deleted, it can be reset to default value).")
	}

	return nil
}

const defaultMapping = "[\r\n{\r\n\"local\":[\r\n{\r\n\"user\":{\r\n\"name\":\"FederationUser\"\r\n}\r\n}\r\n]," +
	"\r\n\"remote\":[\r\n{\r\n\"type\":\"_NAMEID__\"\r\n}\r\n]\r\n}\r\n]"
