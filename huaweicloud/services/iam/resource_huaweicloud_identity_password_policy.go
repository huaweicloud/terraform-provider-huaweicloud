package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/security"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM PUT /v3.0/OS-SECURITYPOLICY/domains/{domainID}/password-policy
// @API IAM GET /v3.0/OS-SECURITYPOLICY/domains/{domainID}/password-policy
func ResourceIdentityPasswordPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePasswordPolicyUpdate,
		UpdateContext: resourcePasswordPolicyUpdate,
		ReadContext:   resourcePasswordPolicyRead,
		DeleteContext: resourcePasswordPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"password_char_combination": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
			"minimum_password_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  8,
			},
			"maximum_consecutive_identical_chars": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"number_of_recent_passwords_disallowed": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"password_validity_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"minimum_password_age": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"password_not_username_or_invert": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"maximum_password_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourcePasswordPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iamClient, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	domainID := cfg.DomainID
	updateOpts := &security.PasswordPolicyOpts{
		MinCharCombination:             utils.Int(d.Get("password_char_combination").(int)),
		MinPasswordLength:              utils.Int(d.Get("minimum_password_length").(int)),
		MaxConsecutiveIdenticalChars:   utils.Int(d.Get("maximum_consecutive_identical_chars").(int)),
		RecentPasswordsDisallowedCount: utils.Int(d.Get("number_of_recent_passwords_disallowed").(int)),
		PasswordValidityPeriod:         utils.Int(d.Get("password_validity_period").(int)),
		MinPasswordAge:                 utils.Int(d.Get("minimum_password_age").(int)),
		PasswordNotUsernameOrInvert:    utils.Bool(d.Get("password_not_username_or_invert").(bool)),
	}

	_, err = security.UpdatePasswordPolicy(iamClient, updateOpts, domainID)
	if err != nil {
		return diag.Errorf("error updating the IAM account password policy: %s", err)
	}

	// set the ID only when creating
	if d.IsNewResource() {
		d.SetId(domainID)
	}

	return resourcePasswordPolicyRead(ctx, d, meta)
}

func resourcePasswordPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iamClient, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	policy, err := security.GetPasswordPolicy(iamClient, d.Id())
	if err != nil {
		return diag.Errorf("error fetching the IAM account password policy")
	}

	log.Printf("[DEBUG] Retrieved the IAM account password policy: %#v", policy)
	mErr := multierror.Append(nil,
		d.Set("password_char_combination", policy.MinCharCombination),
		d.Set("minimum_password_length", policy.MinPasswordLength),
		d.Set("maximum_consecutive_identical_chars", policy.MaxConsecutiveIdenticalChars),
		d.Set("number_of_recent_passwords_disallowed", policy.RecentPasswordsDisallowedCount),
		d.Set("password_validity_period", policy.PasswordValidityPeriod),
		d.Set("minimum_password_age", policy.MinPasswordAge),
		d.Set("password_not_username_or_invert", policy.PasswordNotUsernameOrInvert),
		d.Set("maximum_password_length", policy.MaxPasswordLength),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePasswordPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iamClient, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	defaultOpts := &security.PasswordPolicyOpts{
		MinCharCombination:             utils.Int(2),
		MinPasswordLength:              utils.Int(8),
		MaxConsecutiveIdenticalChars:   utils.Int(0),
		RecentPasswordsDisallowedCount: utils.Int(1),
		PasswordValidityPeriod:         utils.Int(0),
		MinPasswordAge:                 utils.Int(0),
		PasswordNotUsernameOrInvert:    utils.Bool(true),
	}

	_, err = security.UpdatePasswordPolicy(iamClient, defaultOpts, d.Id())
	if err != nil {
		return diag.Errorf("error resetting the IAM account password policy: %s", err)
	}

	return nil
}
