package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v5UserPasswordNonUpdatableParams = []string{"new_password", "old_password"}

// @API IAM POST /v5/caller-password
func ResourceV5UserPassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5UserPasswordCreate,
		ReadContext:   resourceV5UserPasswordRead,
		UpdateContext: resourceV5UserPasswordCreate,
		DeleteContext: resourceV5UserPasswordDelete,

		CustomizeDiff: config.FlexibleForceNew(v5UserPasswordNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"new_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The new password of the user.`,
			},
			"old_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The old password of the user.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceV5UserPasswordCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	changePasswordPath := iamClient.Endpoint + "v5/caller-password"
	options := golangsdk.RequestOpts{
		JSONBody: map[string]interface{}{
			"new_password": d.Get("new_password").(string),
			"old_password": d.Get("old_password").(string),
		},
	}
	_, err = iamClient.Request("POST", changePasswordPath, &options)
	if err != nil {
		return diag.Errorf("error changing user password: %s", err)
	}

	userId, err := getUserId(cfg)
	if err != nil {
		return diag.Errorf("error retrieving current user ID: %s", err)
	}

	d.SetId(userId)

	return nil
}

func resourceV5UserPasswordRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV5UserPasswordDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for changing user password. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
