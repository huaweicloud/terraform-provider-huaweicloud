package iam

import (
	"context"
	"fmt"
	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"time"
)

// @API IAM POST /v3.0/OS-CREDENTIAL/securitytokens
func ResourceIdentitySecurityTokens() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentitySecurityTokensCreate,
		ReadContext:   resourceIdentitySecurityTokensRead,
		DeleteContext: resourceIdentitySecurityTokensDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"effect": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Allow", "Deny",
				}, false),
				ForceNew: true,
			},
			"credential": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expires_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"securitytoken": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIdentitySecurityTokensCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("Error creating IAM client: %s", err)
	}
	err = createSecurityTokens(client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("Error generating uuid: %s", err)
	}
	d.SetId(uuid)
	return nil
}

func createSecurityTokens(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	createSecurityTokensopt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"auth": buildCreateSecurityTokensBodyParams(d),
		},
	}

	resp, err := client.Request("POST", client.Endpoint+"v3.0/OS-CREDENTIAL/securitytokens", &createSecurityTokensopt)
	if err != nil {
		return fmt.Errorf("error retrieving IAM user token %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return fmt.Errorf("error flattening response %s", err)
	}
	d.Set("credential", flattenCredential(respBody))
	return nil
}

func flattenCredential(getKeyRespBody interface{}) []map[string]interface{} {
	credentialRaw := utils.PathSearch("credential", getKeyRespBody, nil)
	if credentialRaw == nil {
		return nil
	}
	res := []map[string]interface{}{
		{
			"expires_at":    utils.PathSearch("expires_at", credentialRaw, nil),
			"access":        utils.PathSearch("access", credentialRaw, nil),
			"secret":        utils.PathSearch("secret", credentialRaw, nil),
			"securitytoken": utils.PathSearch("securitytoken", credentialRaw, nil),
		},
	}
	return res
}

func buildCreateSecurityTokensBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"identity": map[string]interface{}{
			"methods": []string{"token"},
			"token":   map[string]interface{}{},
			"policy": map[string]interface{}{
				"Version": d.Get("version"),
				"Statement": []map[string]interface{}{
					{
						"Action": d.Get("action"),
						"Effect": d.Get("effect"),
					},
				},
			},
		},
	}
	return bodyParams
}

func resourceIdentitySecurityTokensRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("Error creating IAM client: %s", err)
	}
	// 安全获取credential数据
	credentialValue := d.Get("credential")
	if credentialValue == nil {
		return diag.Errorf("credential data is nil")
	}

	credentials := d.Get("credential").([]interface{})
	if len(credentials) == 0 {
		return diag.Errorf("empty credentials array")
	}

	// 获取第一个元素的expires_at
	firstCred := credentials[0].(map[string]interface{})
	expiresAtStr := firstCred["expires_at"].(string)

	// 解析时间
	expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {
		return diag.Errorf("error parsing expires_at: %v", err)
	}
	if time.Now().After(expiresAt) {
		err = createSecurityTokens(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceIdentitySecurityTokensDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting token is not supported. The token is only removed from the state, but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
