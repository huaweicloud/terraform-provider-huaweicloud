package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/security"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v3.0/OS-SECURITYPOLICY/domains/{domainID}/protect-policy
// @API IAM PUT /v3.0/OS-SECURITYPOLICY/domains/{domainID}/protect-policy
func ResourceIdentityProtectionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProtectionPolicyUpdate,
		UpdateContext: resourceProtectionPolicyUpdate,
		ReadContext:   resourceProtectionPolicyRead,
		DeleteContext: resourceProtectionPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"protection_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"verification_mobile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"verification_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"self_management": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"password": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"mobile": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"email": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"self_verification": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildSelfManagement(d *schema.ResourceData) *security.AllowUserOpts {
	raw := d.Get("self_management").([]interface{})
	if len(raw) == 0 {
		// if not specified, keep the previous settings.
		return nil
	}

	item, ok := raw[0].(map[string]interface{})
	if !ok {
		return nil
	}

	allowed := security.AllowUserOpts{}
	if v, ok := item["access_key"]; ok {
		allowed.ManageAccesskey = utils.Bool(v.(bool))
	}
	if v, ok := item["password"]; ok {
		allowed.ManagePassword = utils.Bool(v.(bool))
	}
	if v, ok := item["mobile"]; ok {
		allowed.ManageMobile = utils.Bool(v.(bool))
	}
	if v, ok := item["email"]; ok {
		allowed.ManageEmail = utils.Bool(v.(bool))
	}

	return &allowed
}

func flattenSelfManagement(resp *security.AllowUserBody) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"access_key": resp.ManageAccesskey,
			"password":   resp.ManagePassword,
			"mobile":     resp.ManageMobile,
			"email":      resp.ManageEmail,
		},
	}
}

func resourceProtectionPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iamClient, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	domainID := cfg.DomainID
	enabled := d.Get("protection_enabled").(bool)
	updateOpts := &security.ProtectPolicyOpts{
		Protection: &enabled,
		Email:      utils.String(""),
		Mobile:     utils.String(""),
		AllowUser:  buildSelfManagement(d),
	}

	// verification_mobile and verification_mobile are valid when the protection is enabled
	if enabled {
		var adminCheck string
		if v, ok := d.GetOk("verification_mobile"); ok {
			adminCheck = "on"
			updateOpts.Scene = utils.String("mobile")
			updateOpts.Mobile = utils.String(v.(string))
		} else if v, ok := d.GetOk("verification_email"); ok {
			adminCheck = "on"
			updateOpts.Scene = utils.String("email")
			updateOpts.Email = utils.String(v.(string))
		} else {
			// self verification
			adminCheck = "off"
		}

		updateOpts.AdminCheck = &adminCheck
	}

	_, err = security.UpdateProtectPolicy(iamClient, updateOpts, domainID)
	if err != nil {
		return diag.Errorf("error updating the IAM protection policy: %s", err)
	}

	// set the ID only when creating
	if d.IsNewResource() {
		d.SetId(domainID)
	}

	return resourceProtectionPolicyRead(ctx, d, meta)
}

func resourceProtectionPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iamClient, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	policy, err := security.GetProtectPolicy(iamClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching the IAM protection policy")
	}

	log.Printf("[DEBUG] Retrieved the IAM protection policy: %#v", policy)
	mErr := multierror.Append(nil,
		d.Set("protection_enabled", policy.Protection),
		d.Set("verification_email", policy.Email),
		d.Set("verification_mobile", policy.Mobile),
		d.Set("self_verification", policy.AdminCheck != "on"),
		d.Set("self_management", flattenSelfManagement(&policy.AllowUser)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceProtectionPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iamClient, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	defaultOpts := &security.ProtectPolicyOpts{
		Protection: utils.Bool(false),
		AllowUser: &security.AllowUserOpts{
			ManageAccesskey: utils.Bool(true),
			ManagePassword:  utils.Bool(true),
			ManageMobile:    utils.Bool(true),
			ManageEmail:     utils.Bool(true),
		},
	}

	_, err = security.UpdateProtectPolicy(iamClient, defaultOpts, d.Id())
	if err != nil {
		return diag.Errorf("error resetting the IAM protection policy: %s", err)
	}

	return nil
}
