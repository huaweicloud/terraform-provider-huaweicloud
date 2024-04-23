package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM PUT /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy
// @API IAM GET /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy
func ResourceIdentityLoginPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityLoginPolicyCreateOrUpdate,
		ReadContext:   resourceIdentityLoginPolicyRead,
		UpdateContext: resourceIdentityLoginPolicyCreateOrUpdate,
		DeleteContext: resourceIdentityLoginPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_validity_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"custom_info_for_login": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lockout_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  15,
			},
			"login_failed_times": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"period_with_login_failures": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  15,
			},
			"session_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"show_recent_login_info": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceIdentityLoginPolicyCreateOrUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	product := "iam"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	createLoginPolicyHttpUrl := "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy"
	createLoginPolicyPath := client.Endpoint + createLoginPolicyHttpUrl
	createLoginPolicyPath = strings.ReplaceAll(createLoginPolicyPath, "{domain_id}", cfg.DomainID)
	createLoginPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"login_policy": map[string]interface{}{
				"account_validity_period":    d.Get("account_validity_period"),
				"custom_info_for_login":      d.Get("custom_info_for_login"),
				"lockout_duration":           d.Get("lockout_duration"),
				"login_failed_times":         d.Get("login_failed_times"),
				"period_with_login_failures": d.Get("period_with_login_failures"),
				"session_timeout":            d.Get("session_timeout"),
				"show_recent_login_info":     d.Get("show_recent_login_info"),
			},
		},
	}
	_, err = client.Request("PUT", createLoginPolicyPath, &createLoginPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating IAM login policy: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(cfg.DomainID)
	}

	return nil
}

func resourceIdentityLoginPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	getLoginPolicyProduct := "iam"
	getLoginPolicyClient, err := cfg.NewServiceClient(getLoginPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	getLoginPolicyHttpUrl := "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy"
	getLoginPolicyPath := getLoginPolicyClient.Endpoint + getLoginPolicyHttpUrl
	getLoginPolicyPath = strings.ReplaceAll(getLoginPolicyPath, "{domain_id}", cfg.DomainID)
	getLoginPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getLoginPolicyResp, err := getLoginPolicyClient.Request("GET", getLoginPolicyPath, &getLoginPolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IAM login policy")
	}
	getLoginPolicyRespBody, err := utils.FlattenResponse(getLoginPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("account_validity_period", utils.PathSearch("login_policy.account_validity_period", getLoginPolicyRespBody, nil)),
		d.Set("custom_info_for_login", utils.PathSearch("login_policy.custom_info_for_login", getLoginPolicyRespBody, nil)),
		d.Set("lockout_duration", utils.PathSearch("login_policy.lockout_duration", getLoginPolicyRespBody, nil)),
		d.Set("login_failed_times", utils.PathSearch("login_policy.login_failed_times", getLoginPolicyRespBody, nil)),
		d.Set("period_with_login_failures", utils.PathSearch("login_policy.period_with_login_failures", getLoginPolicyRespBody, nil)),
		d.Set("session_timeout", utils.PathSearch("login_policy.session_timeout", getLoginPolicyRespBody, nil)),
		d.Set("show_recent_login_info", utils.PathSearch("login_policy.show_recent_login_info", getLoginPolicyRespBody, nil)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving IAM login policy resource (%s) fields: %s", d.Id(), mErr)
	}
	return nil
}

func resourceIdentityLoginPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dproduct := "iam"
	client, err := cfg.NewServiceClient(dproduct, region)
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	deleteLoginPolicyHttpUrl := "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy"
	deleteLoginPolicyPath := client.Endpoint + deleteLoginPolicyHttpUrl
	deleteLoginPolicyPath = strings.ReplaceAll(deleteLoginPolicyPath, "{domain_id}", cfg.DomainID)
	deleteLoginPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"login_policy": map[string]interface{}{
				"account_validity_period":    0,
				"custom_info_for_login":      "",
				"lockout_duration":           15,
				"login_failed_times":         5,
				"period_with_login_failures": 15,
				"session_timeout":            60,
				"show_recent_login_info":     false,
			},
		},
	}

	_, err = client.Request("PUT", deleteLoginPolicyPath, &deleteLoginPolicyOpt)
	if err != nil {
		return diag.Errorf("error deleting IAM login policy: %s", err)
	}

	return nil
}
