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

// ResourceIdentityV5UserPassword
// @API IAM POST /v5/caller-password
func ResourceIdentityV5UserPassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityV5UserPasswordCreate,
		ReadContext:   resourceIdentityV5UserPasswordRead,
		UpdateContext: resourceIdentityV5UserPasswordCreate,
		DeleteContext: resourceIdentityV5UserPasswordDelete,

		CustomizeDiff: config.FlexibleForceNew(v5UserPasswordNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"new_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"old_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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

func resourceIdentityV5UserPasswordCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	changePasswordPath := iamClient.Endpoint + "v5/caller-password"
	options := golangsdk.RequestOpts{
		OkCodes: []int{200},
		JSONBody: map[string]interface{}{
			"new_password": d.Get("new_password").(string),
			"old_password": d.Get("old_password").(string),
		},
	}
	_, err = iamClient.Request("POST", changePasswordPath, &options)
	if err != nil {
		return diag.Errorf("error change password: %s", err)
	}
	userId, err := getUserId(cfg)
	if err != nil {
		return diag.Errorf("error retrieving user id: %s", err)
	}
	d.SetId(userId)
	return nil
}

func resourceIdentityV5UserPasswordRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityV5UserPasswordDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting password is not supported. The password is only removed from the state, but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
