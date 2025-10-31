package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentityLoginProtects
// @API IAM GET /v3.0/OS-USER/login-protects
// @API IAM GET /v3.0/OS-USER/users/{user_id}/login-protect
func DataSourceIdentityLoginProtects() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityLoginProtectsRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"login_protects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"verification_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityLoginProtectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	userId := d.Get("user_id").(string)
	if userId == "" {
		return listLoginProtects(iamClient, d)
	}

	return showLoginProtect(iamClient, userId, d)
}

func listLoginProtects(iamClient *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	listLoginProtectsPath := iamClient.Endpoint + "v3.0/OS-USER/login-protects"
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamClient.Request("GET", listLoginProtectsPath, &options)
	if err != nil {
		return diag.Errorf("error listLoginProtects: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	loginProtectsBody := utils.PathSearch("login_protects", respBody, make([]interface{}, 0)).([]interface{})
	loginProtects := make([]interface{}, 0, len(loginProtectsBody))
	for _, loginProtect := range loginProtectsBody {
		loginProtects = append(loginProtects, flattenLoginProtect(loginProtect))
	}
	if err = d.Set("login_protects", loginProtects); err != nil {
		return diag.Errorf("error setting login_protects fields: %s", err)
	}
	return nil
}

func showLoginProtect(iamClient *golangsdk.ServiceClient, userId string, d *schema.ResourceData) diag.Diagnostics {
	showLoginProtectPath := iamClient.Endpoint + "v3.0/OS-USER/users/{user_id}/login-protect"
	showLoginProtectPath = strings.ReplaceAll(showLoginProtectPath, "{user_id}", userId)
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamClient.Request("GET", showLoginProtectPath, &options)
	if err != nil {
		return diag.Errorf("error showLoginProtect: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	loginProtectBody := utils.PathSearch("login_protect", respBody, make([]interface{}, 0))
	loginProtects := append(make([]interface{}, 0, 1), flattenLoginProtect(loginProtectBody))
	if err = d.Set("login_protects", loginProtects); err != nil {
		return diag.Errorf("error setting login_protects fields: %s", err)
	}
	return nil
}

func flattenLoginProtect(loginProtectModel interface{}) map[string]interface{} {
	if loginProtectModel == nil {
		return nil
	}
	loginProtect := make(map[string]interface{})
	loginProtect["user_id"] = utils.PathSearch("user_id", loginProtectModel, "")
	loginProtect["enabled"] = utils.PathSearch("enabled", loginProtectModel, "")
	loginProtect["verification_method"] = utils.PathSearch("verification_method", loginProtectModel, "")
	return loginProtect
}
