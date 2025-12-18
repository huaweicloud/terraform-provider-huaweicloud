package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ResourceIdentityUserInfo
// @API IAM PUT /v3.0/OS-USER/users/{user_id}/info
func ResourceIdentityUserInfo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityUserInfoCreate,
		ReadContext:   resourceIdentityUserInfoRead,
		UpdateContext: resourceIdentityUserInfoUpdate,
		DeleteContext: resourceIdentityUserInfoDelete,

		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email must conform to the email format and be no longer than 255 characters.",
			},
			"mobile": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"email", "mobile"},
				Description:  "The mobile format is `<country code>-<phone number>`, such as 0086-123456789.",
			},
		},
	}
}

func resourceIdentityUserInfoCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId, err := getUserId(cfg)
	if err != nil {
		return diag.Errorf("error getUserId: %s", err)
	}
	updateUserInfoPath := iamClient.Endpoint + "v3.0/OS-USER/users/{user_id}/info"
	updateUserInfoPath = strings.ReplaceAll(updateUserInfoPath, "{user_id}", userId)

	user := map[string]string{}
	if email := d.Get("email").(string); email != "" {
		user["email"] = email
	}
	if mobile := d.Get("mobile").(string); mobile != "" {
		user["mobile"] = mobile
	}
	options := golangsdk.RequestOpts{
		OkCodes:  []int{204},
		JSONBody: map[string]interface{}{"user": user},
	}
	_, err = iamClient.Request("PUT", updateUserInfoPath, &options)
	if err != nil {
		return diag.Errorf("error updateUserInfo: %s", err)
	}
	d.SetId(userId)
	return nil
}

func resourceIdentityUserInfoRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityUserInfoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChanges("email", "mobile") {
		return resourceIdentityUserInfoCreate(ctx, d, meta)
	}
	return nil
}

func resourceIdentityUserInfoDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting email or mobile is not supported. The email and mobile are only removed from the state, " +
		"but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
