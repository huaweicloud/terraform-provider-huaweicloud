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

// @API IAM PUT /v5/password-policy
// @API IAM GET /v5/password-policy
func ResourceV5PasswordPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5PasswordPolicyUpdate,
		ReadContext:   resourceV5PasswordPolicyRead,
		UpdateContext: resourceV5PasswordPolicyUpdate,
		DeleteContext: resourceV5PasswordPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"maximum_consecutive_identical_chars": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The maximum number of consecutive identical characters.`,
			},
			"minimum_password_age": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The minimum password usage time, in minutes.`,
			},
			"minimum_password_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     8,
				Description: `The minimum number of characters that a password must contain.`,
			},
			"password_reuse_prevention": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: `The password cannot be repeated with historical for a certain number of times.`,
			},
			"password_not_username_or_invert": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether the password can be the username or the username spelled backwards.`,
			},
			"password_validity_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The password validity period, in days.`,
			},
			"password_char_combination": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: `The minimum number of character types that a password must contain.`,
			},
			"allow_user_to_change_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether IAM users are allowed to change their own passwords.`,
			},
			"maximum_password_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of characters that a password can contain.`,
			},
			"password_requirements": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The requirements of character that passwords must include.`,
			},
		},
	}
}

func resourceV5PasswordPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	httpUrl := "v5/password-policy"
	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateV5PasswordPolicyBodyParams(d),
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IAM password policy: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(cfg.DomainID)
	}

	return resourceV5PasswordPolicyRead(ctx, d, meta)
}

func buildCreateV5PasswordPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"maximum_consecutive_identical_chars": d.Get("maximum_consecutive_identical_chars"),
		"minimum_password_age":                d.Get("minimum_password_age"),
		"minimum_password_length":             d.Get("minimum_password_length"),
		"password_reuse_prevention":           d.Get("password_reuse_prevention"),
		"password_not_username_or_invert":     d.Get("password_not_username_or_invert"),
		"password_validity_period":            d.Get("password_validity_period"),
		"password_char_combination":           d.Get("password_char_combination"),
		"allow_user_to_change_password":       d.Get("allow_user_to_change_password"),
	}
}

func GetV5PasswordPolicy(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v5/password-policy"
	getPath := client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	passwordPolicy := utils.PathSearch("password_policy", respBody, nil)
	if utils.PathSearch("maximum_consecutive_identical_chars", passwordPolicy, float64(0)).(float64) != 0 ||
		utils.PathSearch("minimum_password_age", passwordPolicy, float64(0)).(float64) != 0 ||
		utils.PathSearch("minimum_password_length", passwordPolicy, float64(0)).(float64) != 8 ||
		utils.PathSearch("password_reuse_prevention", passwordPolicy, float64(0)).(float64) != 1 ||
		!utils.PathSearch("password_not_username_or_invert", passwordPolicy, false).(bool) ||
		utils.PathSearch("password_validity_period", passwordPolicy, float64(0)).(float64) != 0 ||
		utils.PathSearch("password_char_combination", passwordPolicy, float64(0)).(float64) != 2 ||
		!utils.PathSearch("allow_user_to_change_password", passwordPolicy, false).(bool) {
		return passwordPolicy, nil
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v5/password-policy",
			RequestId: "NONE",
			Body:      []byte("All configurations of password policy have been restored to the default value"),
		},
	}
}

func resourceV5PasswordPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	passwordPolicy, err := GetV5PasswordPolicy(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching the IAM account password policy")
	}

	mErr := multierror.Append(
		d.Set("maximum_consecutive_identical_chars", utils.PathSearch("maximum_consecutive_identical_chars", passwordPolicy, nil)),
		d.Set("maximum_password_length", utils.PathSearch("maximum_password_length", passwordPolicy, nil)),
		d.Set("minimum_password_age", utils.PathSearch("minimum_password_age", passwordPolicy, nil)),
		d.Set("minimum_password_length", utils.PathSearch("minimum_password_length", passwordPolicy, nil)),
		d.Set("password_reuse_prevention", utils.PathSearch("password_reuse_prevention", passwordPolicy, nil)),
		d.Set("password_not_username_or_invert", utils.PathSearch("password_not_username_or_invert", passwordPolicy, nil)),
		d.Set("password_requirements", utils.PathSearch("password_requirements", passwordPolicy, nil)),
		d.Set("password_validity_period", utils.PathSearch("password_validity_period", passwordPolicy, nil)),
		d.Set("password_char_combination", utils.PathSearch("password_char_combination", passwordPolicy, nil)),
		d.Set("allow_user_to_change_password", utils.PathSearch("allow_user_to_change_password", passwordPolicy, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV5PasswordPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deletePasswordPolicyHttpUrl := "v5/password-policy"
	deletePasswordPolicyPath := client.Endpoint + deletePasswordPolicyHttpUrl
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

	_, err = client.Request("PUT", deletePasswordPolicyPath, &deletePasswordPolicyOpt)
	if err != nil {
		return diag.Errorf("error resetting the IAM account password policy: %s", err)
	}
	return nil
}
