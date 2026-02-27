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

// @API IAM POST /v3/users/{user_id}/password
func ResourceV3UserPassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3UserPasswordCreate,
		ReadContext:   resourceV3UserPasswordRead,
		UpdateContext: resourceV3UserPasswordUpdate,
		DeleteContext: resourceV3UserPasswordDelete,

		Schema: map[string]*schema.Schema{
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The new password of the IAM user.",
			},
			"original_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The original password of the IAM user.",
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

func resourceV3UserPasswordCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId, err := getUserId(cfg)
	if err != nil {
		return diag.Errorf("error getting user ID: %s", err)
	}
	updateUserPasswordPath := iamClient.Endpoint + "v3/users/{user_id}/password"
	updateUserPasswordPath = strings.ReplaceAll(updateUserPasswordPath, "{user_id}", userId)
	options := golangsdk.RequestOpts{
		OkCodes: []int{204},
		JSONBody: map[string]interface{}{
			"user": map[string]interface{}{
				"password":          d.Get("password").(string),
				"original_password": d.Get("original_password").(string),
			},
		},
	}
	_, err = iamClient.Request("POST", updateUserPasswordPath, &options)
	if err != nil {
		return diag.Errorf("error updateUserPassword: %s", err)
	}
	d.SetId(userId)
	return nil
}

func resourceV3UserPasswordRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV3UserPasswordUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV3UserPasswordDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for modifying password of user. Deleting this resource will
    not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
