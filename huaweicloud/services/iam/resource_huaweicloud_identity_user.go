package iam

import (
	"context"
	"regexp"
	"strings"

	iam_users "github.com/chnsz/golangsdk/openstack/identity/v3.0/users"
	"github.com/chnsz/golangsdk/openstack/identity/v3/users"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIdentityUserV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityUserV3Create,
		ReadContext:   resourceIdentityUserV3Read,
		UpdateContext: resourceIdentityUserV3Update,
		DeleteContext: resourceIdentityUserV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"country_code"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[0-9]{0,32}$"),
					"the phone number must have a maximum of 32 digits"),
			},
			"country_code": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"phone"},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"pwd_reset": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"default", "programmatic", "console",
				}, false),
			},
			"password_strength": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_login": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityUserV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud iam client: %s", err)
	}

	if config.DomainID == "" {
		return fmtp.DiagErrorf("the domain_id must be specified in the provider configuration")
	}

	enabled := d.Get("enabled").(bool)
	reset := d.Get("pwd_reset").(bool)
	createOpts := iam_users.CreateOpts{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Email:         d.Get("email").(string),
		Phone:         d.Get("phone").(string),
		AreaCode:      d.Get("country_code").(string),
		AccessMode:    d.Get("access_type").(string),
		Enabled:       &enabled,
		PasswordReset: &reset,
		DomainID:      config.DomainID,
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	user, err := iam_users.Create(iamClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud iam user: %s", err)
	}

	d.SetId(user.ID)

	return resourceIdentityUserV3Read(ctx, d, meta)
}

func resourceIdentityUserV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud iam client: %s", err)
	}

	user, err := iam_users.Get(iamClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "user")
	}

	logp.Printf("[DEBUG] Retrieved HuaweiCloud user: %#v", user)

	mErr := multierror.Append(nil,
		d.Set("enabled", user.Enabled),
		d.Set("name", user.Name),
		d.Set("description", user.Description),
		d.Set("email", user.Email),
		d.Set("country_code", user.AreaCode),
		d.Set("access_type", user.AccessMode),
		d.Set("password_strength", user.PasswordStrength),
		d.Set("pwd_reset", user.PasswordStatus),
		d.Set("create_time", user.CreateAt),
		d.Set("last_login", user.LastLogin),
	)

	phone := strings.Split(user.Phone, "-")
	if len(phone) > 1 {
		mErr = multierror.Append(mErr, d.Set("phone", phone[1]))
	} else {
		mErr = multierror.Append(mErr, d.Set("phone", user.Phone))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting identity role fields: %s", err)
	}

	return nil
}

func resourceIdentityUserV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud iam client: %s", err)
	}

	var updateOpts iam_users.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	if d.HasChange("email") {
		updateOpts.Email = d.Get("email").(string)
	}

	if d.HasChanges("country_code", "phone") {
		updateOpts.AreaCode = d.Get("country_code").(string)
		updateOpts.Phone = d.Get("phone").(string)
	}

	if d.HasChange("access_type") {
		updateOpts.AccessMode = d.Get("access_type").(string)
	}

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	if d.HasChange("pwd_reset") {
		reset := d.Get("pwd_reset").(bool)
		updateOpts.PasswordReset = &reset
	}

	logp.Printf("[DEBUG] Update Options: %#v", updateOpts)

	// Add password here so it wouldn't go in the above log entry
	if d.HasChange("password") {
		updateOpts.Password = d.Get("password").(string)
	}

	_, err = iam_users.Update(iamClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error updating HuaweiCloud user: %s", err)
	}

	return resourceIdentityUserV3Read(ctx, d, meta)
}

func resourceIdentityUserV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	err = users.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud user: %s", err)
	}

	return nil
}
