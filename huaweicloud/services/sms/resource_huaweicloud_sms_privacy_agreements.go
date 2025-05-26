package sms

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SMS POST /v3/privacy-agreements
// @API SMS GET /v3/privacy-agreements
func ResourcePrivacyAgreements() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivacyAgreementsCreate,
		ReadContext:   resourcePrivacyAgreementsRead,
		UpdateContext: resourcePrivacyAgreementsUpdate,
		DeleteContext: resourcePrivacyAgreementsDelete,

		Schema: map[string]*schema.Schema{
			"flag": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourcePrivacyAgreementsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/privacy-agreements"
	)

	smsClient, err := cfg.SmsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	requestPath := smsClient.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         map[string]interface{}{},
	}

	_, err = smsClient.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating SMS privacy agreements: %s", err)
	}

	d.SetId(region)
	return resourcePrivacyAgreementsRead(ctx, d, meta)
}

func resourcePrivacyAgreementsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                         = meta.(*config.Config)
		region                      = cfg.GetRegion(d)
		getPrivacyAgreementsHttpUrl = "v3/privacy-agreements"
	)
	smsClient, err := cfg.SmsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	basePath := smsClient.Endpoint + getPrivacyAgreementsHttpUrl

	getPrivacyAgreementsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPrivacyAgreementsResp, err := smsClient.Request("GET", basePath, &getPrivacyAgreementsOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting the privacy agreements form server")
	}

	resp, err := utils.FlattenResponse(getPrivacyAgreementsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("flag", utils.PathSearch("flag", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePrivacyAgreementsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourcePrivacyAgreementsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `Deleting privacy agreements resource is not supported. The privacy agreements resource is only removed from the state.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
