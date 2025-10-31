package iam

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v3.0/OS-CREDENTIAL/credentials
// @API IAM GET /v3.0/OS-CREDENTIAL/credentials/{access_key}
func DataSourceIdentityAccessKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityAccessKeyRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"access_key"},
			},
			"access_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_id"},
			},
			"credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_use_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityAccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	getAccessKeyPath := iamClient.Endpoint + "v3.0/OS-CREDENTIAL/credentials"
	accessKey := d.Get("access_key").(string)
	if accessKey != "" {
		getAccessKeyPath = getAccessKeyPath + "/" + accessKey
		return ListAccessKey(getAccessKeyPath, iamClient, d)
	}
	userId := d.Get("user_id").(string)
	return ShowAccessKey(userId, getAccessKeyPath, iamClient, d)
}

func ListAccessKey(path string, iamClient *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	response, err := iamClient.Request("GET", path, &options)
	if err != nil {
		return diag.Errorf("iamV3Client response error : %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	credential := flattenAccessKeyList(utils.PathSearch("credential", respBody, nil))
	if err = d.Set("credentials", credential); err != nil {
		return diag.Errorf("error setting credentials fields: %s", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	return nil
}

func ShowAccessKey(userId string, path string, iamClient *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	if userId != "" {
		path = path + "?user_id=" + userId
	}
	response, err := iamClient.Request("GET", path, &options)
	if err != nil {
		return diag.Errorf("iamV3Client response error : %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	credentials := flattenAccessKey(utils.PathSearch("credentials", respBody, make([]interface{}, 0)).([]interface{}))
	if err = d.Set("credentials", credentials); err != nil {
		return diag.Errorf("error setting credentials fields: %s", err)
	}
	return nil
}

func flattenAccessKeyList(credentialBody interface{}) []interface{} {
	res := map[string]interface{}{
		"user_id":       utils.PathSearch("user_id", credentialBody, ""),
		"access":        utils.PathSearch("access", credentialBody, ""),
		"status":        utils.PathSearch("status", credentialBody, ""),
		"create_time":   utils.PathSearch("create_time", credentialBody, ""),
		"description":   utils.PathSearch("description", credentialBody, ""),
		"last_use_time": utils.PathSearch("last_use_time", credentialBody, ""),
	}
	result := append(make([]interface{}, 0, 1), res)
	return result
}

func flattenAccessKey(credentialsModel []interface{}) []map[string]interface{} {
	credentials := make([]map[string]interface{}, len(credentialsModel))
	for i, credential := range credentialsModel {
		credentials[i] = map[string]interface{}{
			"user_id":     utils.PathSearch("user_id", credential, ""),
			"access":      utils.PathSearch("access", credential, ""),
			"status":      utils.PathSearch("status", credential, ""),
			"create_time": utils.PathSearch("create_time", credential, ""),
			"description": utils.PathSearch("description", credential, ""),
		}
	}
	return credentials
}
