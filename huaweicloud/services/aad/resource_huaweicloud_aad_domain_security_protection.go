package aad

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableDomainSecurityProtectionParams = []string{
	"domain_id",
	"waf_switch",
	"cc_switch",
}

// @API AAD POST /v1/{project_id}/aad/external/domains/switch
func ResourceDomainSecurityProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainSecurityProtectionCreate,
		ReadContext:   resourceDomainSecurityProtectionRead,
		UpdateContext: resourceDomainSecurityProtectionUpdate,
		DeleteContext: resourceDomainSecurityProtectionDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableDomainSecurityProtectionParams),

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the domain ID.`,
			},
			"waf_switch": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies whether to enable basic web protection.`,
			},
			"cc_switch": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies whether to enable CC protection.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildDomainSecurityProtectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"domain_id":  d.Get("domain_id"),
		"waf_switch": d.Get("waf_switch"),
		"cc_switch":  d.Get("cc_switch"),
	}
}

func resourceDomainSecurityProtectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/aad/external/domains/switch"
		product  = "aad"
		domainID = d.Get("domain_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDomainSecurityProtectionBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating AAD domain security protection: %s", err)
	}

	d.SetId(domainID)
	return resourceDomainSecurityProtectionRead(ctx, d, meta)
}

func resourceDomainSecurityProtectionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDomainSecurityProtectionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDomainSecurityProtectionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to modify AAD security protection. Deleting this 
resource will not change the current AAD security protection, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
