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

// ResourceIdentityV5PasswordPolicy
// @API IAM PUT /v5/password-policy
// @API IAM GET /v5/password-policy
func ResourceIdentityV5PasswordPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityV5PasswordPolicyUpdate,
		ReadContext:   resourceIdentityV5PasswordPolicyRead,
		UpdateContext: resourceIdentityV5PasswordPolicyUpdate,
		DeleteContext: resourceIdentityV5PasswordPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"maximum_consecutive_identical_chars": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"minimum_password_age": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"minimum_password_length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"password_reuse_prevention": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"password_not_username_or_invert": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"password_validity_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"password_char_combination": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"allow_user_to_change_password": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"maximum_password_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"password_requirements": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityV5PasswordPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	createPasswordPolicyHttpUrl := "v5/password-policy"
	createPasswordPolicyPath := iamClient.Endpoint + createPasswordPolicyHttpUrl
	createPasswordPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreatePasswordPolicyOptBodyParams(d),
	}

	_, err = iamClient.Request("PUT", createPasswordPolicyPath, &createPasswordPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating IAM password policy: %s", err)
	}
	if d.IsNewResource() {
		d.SetId(cfg.DomainID)
	}
	return resourceIdentityV5PasswordPolicyRead(ctx, d, meta)
}

func buildCreatePasswordPolicyOptBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"maximum_consecutive_identical_chars": d.Get("maximum_consecutive_identical_chars"),
		"minimum_password_age":                d.Get("minimum_password_age"),
		"minimum_password_length":             d.Get("minimum_password_length"),
		"password_reuse_prevention":           d.Get("password_reuse_prevention"),
		"password_not_username_or_invert":     d.Get("password_not_username_or_invert"),
		"password_validity_period":            d.Get("password_validity_period"),
		"password_char_combination":           d.Get("password_char_combination"),
		"allow_user_to_change_password":       d.Get("allow_user_to_change_password"),
	}
	return bodyParams
}

func resourceIdentityV5PasswordPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	getPasswordPolicyHttpUrl := "v5/password-policy"
	getPasswordPolicyPath := iamClient.Endpoint + getPasswordPolicyHttpUrl
	getPasswordPolicyOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getPasswordPolicyResp, err := iamClient.Request("GET", getPasswordPolicyPath, &getPasswordPolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching the IAM account password policy")
	}
	getPasswordPolicyRespBody, err := utils.FlattenResponse(getPasswordPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("maximum_consecutive_identical_chars", utils.PathSearch("password_policy.maximum_consecutive_identical_chars",
			getPasswordPolicyRespBody, nil)),
		d.Set("maximum_password_length", utils.PathSearch("password_policy.maximum_password_length", getPasswordPolicyRespBody, nil)),
		d.Set("minimum_password_age", utils.PathSearch("password_policy.minimum_password_age", getPasswordPolicyRespBody, nil)),
		d.Set("minimum_password_length", utils.PathSearch("password_policy.minimum_password_length", getPasswordPolicyRespBody, nil)),
		d.Set("password_reuse_prevention", utils.PathSearch("password_policy.password_reuse_prevention", getPasswordPolicyRespBody, nil)),
		d.Set("password_not_username_or_invert", utils.PathSearch("password_policy.password_not_username_or_invert", getPasswordPolicyRespBody, nil)),
		d.Set("password_requirements", utils.PathSearch("password_policy.password_requirements", getPasswordPolicyRespBody, nil)),
		d.Set("password_validity_period", utils.PathSearch("password_policy.password_validity_period", getPasswordPolicyRespBody, nil)),
		d.Set("password_char_combination", utils.PathSearch("password_policy.password_char_combination", getPasswordPolicyRespBody, nil)),
		d.Set("allow_user_to_change_password", utils.PathSearch("password_policy.allow_user_to_change_password", getPasswordPolicyRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityV5PasswordPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	deletePasswordPolicyHttpUrl := "v5/password-policy"
	deletePasswordPolicyPath := iamClient.Endpoint + deletePasswordPolicyHttpUrl
	deletePasswordPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"maximum_consecutive_identical_chars": 0,
			"minimum_password_age":                0,
			"minimum_password_length":             8,
			"password_reuse_prevention":           1,
			"password_not_username_or_invert":     true,
			"password_validity_period":            0,
			"password_char_combination":           2,
			"allow_user_to_change_password":       true,
		},
	}

	_, err = iamClient.Request("PUT", deletePasswordPolicyPath, &deletePasswordPolicyOpt)
	if err != nil {
		return diag.Errorf("error resetting the IAM account password policy: %s", err)
	}
	return nil
}
