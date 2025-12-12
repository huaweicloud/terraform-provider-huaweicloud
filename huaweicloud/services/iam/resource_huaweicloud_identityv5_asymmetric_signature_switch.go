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

// ResourceIdentityV5AsymmetricSignatureSwitch
// @API IAM PUT /v5/asymmetric-signature-switch
// @API IAM GET /v5/asymmetric-signature-switch
func ResourceIdentityV5AsymmetricSignatureSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityV5AsymmetricSignatureSwitchUpdate,
		ReadContext:   resourceIdentityV5AsymmetricSignatureSwitchRead,
		UpdateContext: resourceIdentityV5AsymmetricSignatureSwitchUpdate,
		DeleteContext: resourceIdentityV5AsymmetricSignatureSwitchDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"asymmetric_signature_switch": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceIdentityV5AsymmetricSignatureSwitchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error set IAM asymmetric signature switch: %s", err)
	}
	d.SetId(cfg.DomainID)
	return resourceIdentityV5AsymmetricSignatureSwitchRead(ctx, d, meta)
}

func resourceIdentityV5AsymmetricSignatureSwitchRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return common.CheckDeletedDiag(d, err, "error get IAM asymmetric signature switch")
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	asymmetricSignatureSwitch := utils.PathSearch("asymmetric_signature.asymmetric_signature_switch", respBody, nil)
	if asymmetricSignatureSwitch == nil {
		return diag.Errorf("error getting IAM asymmetric signature switch: asymmetric_signature_switch is not found in response")
	}
	err = d.Set("asymmetric_signature_switch", asymmetricSignatureSwitch)
	if err != nil {
		return diag.Errorf("error setting asymmetric_signature_switch fields: %s", err)
	}
	return nil
}

func resourceIdentityV5AsymmetricSignatureSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting asymmetric signature switch is not supported. " +
		"The asymmetric signature switch is only removed from the state, but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
