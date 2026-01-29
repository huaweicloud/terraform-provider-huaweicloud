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

// @API IAM PUT /v5/login-policy
// @API IAM GET /v5/login-policy
func ResourceV5LoginPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5LoginPolicyUpdate,
		ReadContext:   resourceV5LoginPolicyRead,
		UpdateContext: resourceV5LoginPolicyUpdate,
		DeleteContext: resourceV5LoginPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_validity_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The validity period to disable users, in days.`,
			},
			"custom_info_for_login": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The custom information that will be displayed upon successful login.`,
			},
			"lockout_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The lockout duration after multiple failed login attempts, in minutes.`,
			},
			"login_failed_times": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The number of consecutive failed login attempts before the account is locked.`,
			},
			"period_with_login_failures": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The period to reset the account lockout counter, in minutes.`,
			},
			"session_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The session timeout duration after user login, in minutes.`,
			},
			"show_recent_login_info": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to display the last login information.`,
			},
			"allow_address_netmasks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_netmask": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The IP address or network segment.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description information`,
						},
					},
				},
				Description: `IP address list or network segment list that are allowed to access.`,
			},
			"allow_ip_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_range": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The IP address range.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description information.`,
						},
					},
				},
				Description: `The IP address range list that are allowed to access.`,
			},
			"allow_ip_ranges_ipv6": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_range": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The IPv6 address range.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description information.`,
						},
					},
				},
				Description: `The IPv6 address range list that are allowed to access.`,
			},
		},
	}
}

func resourceV5LoginPolicyUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createLoginPolicyHttpUrl := "v5/login-policy"
	createLoginPolicyPath := iamClient.Endpoint + createLoginPolicyHttpUrl
	createLoginPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateV5LoginPolicyBodyParams(d),
	}

	_, err = iamClient.Request("PUT", createLoginPolicyPath, &createLoginPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating login policy: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(cfg.DomainID)
	}
	return nil
}

func buildCreateV5LoginPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
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

func resourceV5LoginPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	getLoginPolicyHttpUrl := "v5/login-policy"
	getLoginPolicyPath := iamClient.Endpoint + getLoginPolicyHttpUrl
	getLoginPolicyOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getLoginPolicyResp, err := iamClient.Request("GET", getLoginPolicyPath, &getLoginPolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving login policy")
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
		return diag.Errorf("error saving login policy resource (%s) fields: %s", d.Id(), mErr)
	}

	return nil
}

func resourceV5LoginPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
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
		return diag.Errorf("error deleting login policy: %s", err)
	}
	return nil
}
