package cdn

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var domainOwnerVerifyNonUpdatableParams = []string{"domain_name", "verify_type"}

// @API CDN POST /v1.0/cdn/configuration/domains/{domain_name}/verify-owner
func ResourceDomainOwnerVerify() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainOwnerVerifyCreate,
		ReadContext:   resourceDomainOwnerVerifyRead,
		UpdateContext: resourceDomainOwnerVerifyUpdate,
		DeleteContext: resourceDomainOwnerVerifyDelete,

		CustomizeDiff: config.FlexibleForceNew(domainOwnerVerifyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The domain name to be verified.`,
			},

			// Optional parameters.
			"verify_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The verification method.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildDomainOwnerVerifyBodyParams(d *schema.ResourceData) map[string]interface{} {
	verifyType := "all"

	if v, ok := d.GetOk("verify_type"); ok {
		verifyType = v.(string)
	}

	return map[string]interface{}{
		"verify_type": verifyType,
	}
}

func resourceDomainOwnerVerifyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		httpUrl    = "v1.0/cdn/configuration/domains/{domain_name}/verify-owner"
		domainName = d.Get("domain_name").(string)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	verifyPath := client.Endpoint + httpUrl
	verifyPath = strings.ReplaceAll(verifyPath, "{domain_name}", domainName)
	verifyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildDomainOwnerVerifyBodyParams(d),
	}

	resp, err := client.Request("POST", verifyPath, &verifyOpt)
	if err != nil {
		return diag.Errorf("error verifying domain owner: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	result := utils.PathSearch("result", respBody, false).(bool)
	if !result {
		return diag.Errorf("domain ownership verification failed for domain: %s", domainName)
	}

	return resourceDomainOwnerVerifyRead(ctx, d, meta)
}

func resourceDomainOwnerVerifyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDomainOwnerVerifyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDomainOwnerVerifyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This is resource only a one-time action resource for verify domain owner. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
