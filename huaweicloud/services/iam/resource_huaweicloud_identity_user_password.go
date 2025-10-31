package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ResourceIdentityUserPassword
// @API IAM POST /v3/users/{user_id}/password
func ResourceIdentityUserPassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityUserPasswordCreate,
		ReadContext:   resourceIdentityUserPasswordRead,
		DeleteContext: resourceIdentityUserPasswordDelete,

		Schema: map[string]*schema.Schema{
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"original_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
		},
	}
}

func resourceIdentityUserPasswordCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId, err := getUserId(cfg)
	if err != nil {
		return diag.Errorf("error getUserId: %s", err)
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

func resourceIdentityUserPasswordRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityUserPasswordDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting password is not supported. The password is only removed from the state, but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
