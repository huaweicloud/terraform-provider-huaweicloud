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
func ResourceV3LoginPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3LoginPolicyCreateOrUpdate,
		ReadContext:   resourceV3LoginPolicyRead,
		UpdateContext: resourceV3LoginPolicyCreateOrUpdate,
		DeleteContext: resourceV3LoginPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_validity_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The validity period (days) to disable users if they have not logged in within the period.`,
			},
			"custom_info_for_login": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The custom information that will be displayed upon successful login.`,
			},
			"lockout_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     15,
				Description: `The duration (minutes) to lock users out.`,
			},
			"login_failed_times": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: `The number of unsuccessful login attempts to lock users out.`,
			},
			"period_with_login_failures": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     15,
				Description: `The period (minutes) to count the number of unsuccessful login attempts.`,
			},
			"session_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
				Description: `The session timeout (minutes) that will apply if you or users created using your account
do not perform any operations within a specific period.`,
			},
			"show_recent_login_info": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to display last login information upon successful login.`,
			},
		},
	}
}

func updateV3LoginPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string) error {
	httpUrl := "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{domain_id}", domainId)
	updateOpt := golangsdk.RequestOpts{
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
	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceV3LoginPolicyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		domainId = cfg.DomainID
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	err = updateV3LoginPolicy(client, d, domainId)
	if err != nil {
		return diag.Errorf("error updating IAM login policy: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(domainId)
	}

	return resourceV3LoginPolicyRead(ctx, d, meta)
}

func GetV3LoginPolicy(client *golangsdk.ServiceClient, domainId string) (interface{}, error) {
	httpUrl := "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", domainId)
	getLoginPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getLoginPolicyOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	if utils.PathSearch("login_policy.account_validity_period", respBody, float64(0)).(float64) == 0 &&
		utils.PathSearch("login_policy.custom_info_for_login", respBody, "").(string) == "" &&
		utils.PathSearch("login_policy.lockout_duration", respBody, float64(0)).(float64) == 15 &&
		utils.PathSearch("login_policy.login_failed_times", respBody, float64(0)).(float64) == 5 &&
		utils.PathSearch("login_policy.period_with_login_failures", respBody, float64(0)).(float64) == 15 &&
		utils.PathSearch("login_policy.session_timeout", respBody, float64(0)).(float64) == 60 &&
		!utils.PathSearch("login_policy.show_recent_login_info", respBody, false).(bool) {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy",
				RequestId: "NONE",
				Body:      []byte("All configurations of login policy have been restored to the default value"),
			},
		}
	}
	return respBody, nil
}

func resourceV3LoginPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		domainId = cfg.DomainID
	)

	getLoginPolicyClient, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	respBody, err := GetV3LoginPolicy(getLoginPolicyClient, domainId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving login policy")
	}

	mErr := multierror.Append(nil,
		d.Set("account_validity_period", utils.PathSearch("login_policy.account_validity_period", respBody, nil)),
		d.Set("custom_info_for_login", utils.PathSearch("login_policy.custom_info_for_login", respBody, nil)),
		d.Set("lockout_duration", utils.PathSearch("login_policy.lockout_duration", respBody, nil)),
		d.Set("login_failed_times", utils.PathSearch("login_policy.login_failed_times", respBody, nil)),
		d.Set("period_with_login_failures", utils.PathSearch("login_policy.period_with_login_failures", respBody, nil)),
		d.Set("session_timeout", utils.PathSearch("login_policy.session_timeout", respBody, nil)),
		d.Set("show_recent_login_info", utils.PathSearch("login_policy.show_recent_login_info", respBody, nil)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving IAM login policy resource (%s) fields: %s", d.Id(), mErr)
	}
	return nil
}

func resourceV3LoginPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dproduct := "iam"
	client, err := cfg.NewServiceClient(dproduct, region)
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	httpUrl := "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy"
	restorePath := client.Endpoint + httpUrl
	restorePath = strings.ReplaceAll(restorePath, "{domain_id}", cfg.DomainID)
	restoreOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			// Default configurations of login policy to be restored
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

	_, err = client.Request("PUT", restorePath, &restoreOpt)
	if err != nil {
		return diag.Errorf("error deleting IAM login policy: %s", err)
	}
	return nil
}
