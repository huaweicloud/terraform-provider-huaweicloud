package iam

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

// ResourceIdentityV5LoginPolicy
// @API IAM PUT /v5/login-policy
// @API IAM GET /v5/login-policy
func ResourceIdentityV5LoginPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityV5LoginPolicyUpdate,
		ReadContext:   resourceIdentityV5LoginPolicyRead,
		UpdateContext: resourceIdentityV5LoginPolicyUpdate,
		DeleteContext: resourceIdentityV5LoginPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_validity_period": {
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
			},
			"login_failed_times": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"period_with_login_failures": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"session_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"show_recent_login_info": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"allow_address_netmasks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_netmask": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"allow_ip_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_range": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"allow_ip_ranges_ipv6": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_range": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceIdentityV5LoginPolicyUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	createLoginPolicyHttpUrl := "v5/login-policy"
	createLoginPolicyPath := iamClient.Endpoint + createLoginPolicyHttpUrl
	createLoginPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateLoginPolicyBodyParams(d),
	}

	_, err = iamClient.Request("PUT", createLoginPolicyPath, &createLoginPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating IAM login policy: %s", err)
	}
	if d.IsNewResource() {
		d.SetId(cfg.DomainID)
	}
	return nil
}

func buildCreateLoginPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_validity_period":       d.Get("user_validity_period"),
		"custom_info_for_login":      d.Get("custom_info_for_login"),
		"lockout_duration":           d.Get("lockout_duration"),
		"login_failed_times":         d.Get("login_failed_times"),
		"period_with_login_failures": d.Get("period_with_login_failures"),
		"session_timeout":            d.Get("session_timeout"),
		"show_recent_login_info":     d.Get("show_recent_login_info"),
		"allow_address_netmasks":     d.Get("allow_address_netmasks"),
		"allow_ip_ranges":            d.Get("allow_ip_ranges"),
		"allow_ip_ranges_ipv6":       d.Get("allow_ip_ranges_ipv6"),
	}
	return bodyParams
}

func resourceIdentityV5LoginPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	getLoginPolicyHttpUrl := "v5/login-policy"
	getLoginPolicyPath := iamClient.Endpoint + getLoginPolicyHttpUrl
	getLoginPolicyOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getLoginPolicyResp, err := iamClient.Request("GET", getLoginPolicyPath, &getLoginPolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IAM login policy")
	}
	getLoginPolicyRespBody, err := utils.FlattenResponse(getLoginPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("user_validity_period", utils.PathSearch("login_policy.user_validity_period", getLoginPolicyRespBody, nil)),
		d.Set("custom_info_for_login", utils.PathSearch("login_policy.custom_info_for_login", getLoginPolicyRespBody, nil)),
		d.Set("lockout_duration", utils.PathSearch("login_policy.lockout_duration", getLoginPolicyRespBody, nil)),
		d.Set("login_failed_times", utils.PathSearch("login_policy.login_failed_times", getLoginPolicyRespBody, nil)),
		d.Set("period_with_login_failures", utils.PathSearch("login_policy.period_with_login_failures", getLoginPolicyRespBody, nil)),
		d.Set("session_timeout", utils.PathSearch("login_policy.session_timeout", getLoginPolicyRespBody, nil)),
		d.Set("show_recent_login_info", utils.PathSearch("login_policy.show_recent_login_info", getLoginPolicyRespBody, nil)),
		d.Set("allow_address_netmasks", utils.PathSearch("login_policy.allow_address_netmasks", getLoginPolicyRespBody, nil)),
		d.Set("allow_ip_ranges", utils.PathSearch("login_policy.allow_ip_ranges", getLoginPolicyRespBody, nil)),
		d.Set("allow_ip_ranges_ipv6", utils.PathSearch("login_policy.allow_ip_ranges_ipv6", getLoginPolicyRespBody, nil)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving IAM login policy resource (%s) fields: %s", d.Id(), mErr)
	}
	return nil
}

func resourceIdentityV5LoginPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	deleteLoginPolicyHttpUrl := "v5/login-policy"
	deleteLoginPolicyPath := iamClient.Endpoint + deleteLoginPolicyHttpUrl
	allowIpRanges := make([]map[string]interface{}, 1)
	allowIpRanges[0] = map[string]interface{}{
		"ip_range":    "0.0.0.0-255.255.255.255",
		"description": "",
	}
	allowIpRangesIpv6 := make([]map[string]interface{}, 1)
	allowIpRangesIpv6[0] = map[string]interface{}{
		"ip_range":    "0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF",
		"description": "",
	}
	deleteLoginPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"user_validity_period":       0,
			"custom_info_for_login":      "",
			"lockout_duration":           15,
			"login_failed_times":         5,
			"period_with_login_failures": 15,
			"session_timeout":            60,
			"show_recent_login_info":     false,
			"allow_address_netmasks":     make([]map[string]interface{}, 0),
			"allow_ip_ranges":            allowIpRanges,
			"allow_ip_ranges_ipv6":       allowIpRangesIpv6,
		},
	}

	_, err = iamClient.Request("PUT", deleteLoginPolicyPath, &deleteLoginPolicyOpt)
	if err != nil {
		return diag.Errorf("error deleting IAM login policy: %s", err)
	}
	return nil
}
