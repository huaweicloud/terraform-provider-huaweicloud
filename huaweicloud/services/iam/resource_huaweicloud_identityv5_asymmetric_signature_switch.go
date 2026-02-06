package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM PUT /v5/asymmetric-signature-switch
// @API IAM GET /v5/asymmetric-signature-switch
func ResourceV5AsymmetricSignatureSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5AsymmetricSignatureSwitchUpdate,
		ReadContext:   resourceV5AsymmetricSignatureSwitchRead,
		UpdateContext: resourceV5AsymmetricSignatureSwitchUpdate,
		DeleteContext: resourceV5AsymmetricSignatureSwitchDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"asymmetric_signature_switch": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether to enable the asymmetric signature function.`,
			},
		},
	}
}

func resourceV5AsymmetricSignatureSwitchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	path := iamClient.Endpoint + "v5/asymmetric-signature-switch"
	reqOpts := golangsdk.RequestOpts{
		OkCodes: []int{204},
		JSONBody: map[string]interface{}{
			"asymmetric_signature": map[string]interface{}{
				"asymmetric_signature_switch": d.Get("asymmetric_signature_switch").(bool),
			},
		},
	}
	_, err = iamClient.Request("PUT", path, &reqOpts)
	if err != nil {
		return diag.Errorf("unable to set the asymmetric signature switch: %s", err)
	}

	d.SetId(cfg.DomainID)
	return resourceV5AsymmetricSignatureSwitchRead(ctx, d, meta)
}

func resourceV5AsymmetricSignatureSwitchRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	path := iamClient.Endpoint + "v5/asymmetric-signature-switch"
	reqOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := iamClient.Request("GET", path, &reqOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting the asymmetric signature switch")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	asymmetricSignatureSwitch := utils.PathSearch("asymmetric_signature.asymmetric_signature_switch", respBody, nil)
	if asymmetricSignatureSwitch == nil {
		return diag.Errorf("unable to find the asymmetric signature switch in the API response")
	}

	err = d.Set("asymmetric_signature_switch", asymmetricSignatureSwitch)
	if err != nil {
		return diag.Errorf("error setting asymmetric_signature_switch fields: %s", err)
	}
	return nil
}

func resourceV5AsymmetricSignatureSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to enable or disable the asymmetric signature function.
Deleting this resource will not clear the corresponding request record, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
