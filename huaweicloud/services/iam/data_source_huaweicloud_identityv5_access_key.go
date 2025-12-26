package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/users/{user_id}/access-keys
// @API IAM GET /v5/users/{user_id}/access-keys/{access_key_id}/last-used
func DataSourceIdentityV5AccessKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5AccessKeyRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_used_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityV5AccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId := d.Get("user_id").(string)

	accessKey, err := getAccessKeyV5(client, userId)
	if err != nil {
		return diag.Errorf("error retrieving access key: %s", err)
	}
	if accessKey == nil {
		return diag.Errorf("not found access key : %s", err)
	}

	accessKeyId := utils.PathSearch("access_key_id", accessKey, "").(string)

	lastUsedAt, err := getAccessKeyLastUsedV5(client, userId, accessKeyId)
	if err != nil {
		return diag.Errorf("error retrieving access key last used time: %s", err)
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("access_key_id", accessKeyId),
		d.Set("status", utils.PathSearch("status", accessKey, nil)),
		d.Set("created_at", utils.PathSearch("created_at", accessKey, nil)),
		d.Set("last_used_at", lastUsedAt),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getAccessKeyV5(client *golangsdk.ServiceClient, userId string) (interface{}, error) {
	path := client.Endpoint + "v5/users/" + userId + "/access-keys"
	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	r, err := client.Request("GET", path, &reqOpt)
	if err != nil {
		return nil, err
	}

	resp, err := utils.FlattenResponse(r)
	if err != nil {
		return nil, err
	}

	accessKey := utils.PathSearch("access_keys", resp, make([]interface{}, 0)).([]interface{})
	if len(accessKey) < 1 {
		return nil, nil
	}

	return accessKey[0], nil
}

func getAccessKeyLastUsedV5(client *golangsdk.ServiceClient, userId, accessKeyId string) (string, error) {
	path := client.Endpoint + "v5/users/" + userId + "/access-keys/" + accessKeyId + "/last-used"
	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	r, err := client.Request("GET", path, &reqOpt)
	if err != nil {
		if errCode, ok := err.(golangsdk.ErrDefault404); ok {
			fmt.Printf("Access key last used info not found: %s\n", errCode.Error())
			return "", nil
		}
		return "", err
	}

	resp, err := utils.FlattenResponse(r)
	if err != nil {
		return "", err
	}

	lastUsedAt := utils.PathSearch("access_key_last_used.last_used_at", resp, "").(string)
	return lastUsedAt, nil
}
