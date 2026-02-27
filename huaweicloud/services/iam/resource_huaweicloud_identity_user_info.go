package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v3UserInfoNonUpdatableParams = []string{
	"email",
	"mobile",
}

// @API IAM PUT /v3.0/OS-USER/users/{user_id}/info
func ResourceV3UserInfo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3UserInfoCreate,
		ReadContext:   resourceV3UserInfoRead,
		UpdateContext: resourceV3UserInfoUpdate,
		DeleteContext: resourceV3UserInfoDelete,

		CustomizeDiff: config.FlexibleForceNew(v3UserInfoNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email of the user.",
			},
			"mobile": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"email", "mobile"},
				Description:  "The mobile phone number of the user.",
			},

			// Internal
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceV3UserInfoCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId, err := getUserId(cfg)
	if err != nil {
		return diag.Errorf("error getting user ID: %s", err)
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

func resourceV3UserInfoRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV3UserInfoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChanges("email", "mobile") {
		return resourceV3UserInfoCreate(ctx, d, meta)
	}
	return nil
}

func resourceV3UserInfoDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for creating user information. Deleting this resource will
    not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
