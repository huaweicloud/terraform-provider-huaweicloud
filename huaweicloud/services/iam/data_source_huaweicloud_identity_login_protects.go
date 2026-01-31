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
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user id.`,
			},

			"login_protects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The login status protection information list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user id.`,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to enable login protection.`,
						},
						"verification_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The login verification method.`,
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
	var (
		httpUrl = "v3.0/OS-USER/login-protects"
	)

	listPath := iamClient.Endpoint + httpUrl
	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}
	response, err := iamClient.Request("GET", listPath, &listOpts)
	if err != nil {
		return diag.Errorf("error listing login protects: %s", err)
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.Errorf("error querying login protects: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	loginProtectsBody := utils.PathSearch("login_protects", respBody, make([]interface{}, 0)).([]interface{})
	loginProtects := make([]interface{}, 0, len(loginProtectsBody))
	for _, loginProtect := range loginProtectsBody {
		loginProtects = append(loginProtects, flattenLoginProtect(loginProtect))
	}
	if err = d.Set("login_protects", loginProtects); err != nil {
		return diag.Errorf("error setting login protects fields: %s", err)
	}
	return nil
}

func showLoginProtect(iamClient *golangsdk.ServiceClient, userId string, d *schema.ResourceData) diag.Diagnostics {
	var (
		httpUrl = "v3.0/OS-USER/users/{user_id}/login-protect"
	)
	getPath := iamClient.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{user_id}", userId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	response, err := iamClient.Request("GET", getPath, &getOpts)
	if err != nil {
		return diag.Errorf("error showing login protect: %s", err)
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.Errorf("error querying login protect: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

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

	return map[string]interface{}{
		"user_id":             utils.PathSearch("user_id", loginProtectModel, nil),
		"enabled":             utils.PathSearch("enabled", loginProtectModel, nil),
		"verification_method": utils.PathSearch("verification_method", loginProtectModel, nil),
	}
}
