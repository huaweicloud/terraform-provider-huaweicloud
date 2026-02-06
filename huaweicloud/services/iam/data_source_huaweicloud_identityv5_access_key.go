package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/users/{user_id}/access-keys
// @API IAM GET /v5/users/{user_id}/access-keys/{access_key_id}/last-used
func DataSourceV5AccessKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5AccessKeyRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the IAM user.`,
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the access key.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the access key.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the access key.`,
			},
			"last_used_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the access key was last used.`,
			},
		},
	}
}

func dataSourceV5AccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId := d.Get("user_id").(string)
	accessKeys, err := listV5AccessKeys(client, userId)
	if err != nil {
		return diag.Errorf("error retrieving access key of the user (%s): %s", userId, err)
	}

	accessKey := utils.PathSearch("[0]", accessKeys, nil)
	accessKeyId := utils.PathSearch("access_key_id", accessKey, "").(string)
	if accessKeyId == "" {
		return diag.Errorf("unable to find the access key ID from the API response: %s", err)
	}

	lastUsedAt, err := getV5AccessKeyLastUsedAt(client, userId, accessKeyId)
	if err != nil {
		// To avoid the error of this interface, causing the data source to fail to use normally, use log to record the error.
		log.Printf("[ERROR] error retrieving last used time of the access key (%s): %s", accessKeyId, err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("access_key_id", accessKeyId),
		d.Set("status", utils.PathSearch("status", accessKey, nil)),
		d.Set("created_at", utils.PathSearch("created_at", accessKey, nil)),
		d.Set("last_used_at", lastUsedAt),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
