package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// DataSourceIdentityCheckGroupMembership
// @API IAM HEAD /v3/groups/{group_id}/users/{user_id}
func DataSourceIdentityCheckGroupMembership() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCheckGroupMembershipRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"result": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityCheckGroupMembershipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	userId := d.Get("user_id").(string)
	checkGroupMembershipPath := iamClient.Endpoint + "v3/groups/{group_id}/users/{user_id}"
	checkGroupMembershipPath = strings.ReplaceAll(checkGroupMembershipPath, "{group_id}", groupId)
	checkGroupMembershipPath = strings.ReplaceAll(checkGroupMembershipPath, "{user_id}", userId)
	options := golangsdk.RequestOpts{
		OkCodes: []int{204, 404},
	}
	response, err := iamClient.Request("HEAD", checkGroupMembershipPath, &options)
	if err != nil {
		return diag.Errorf("error checkGroupMembership: %s", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generate UUID: %s", err)
	}
	d.SetId(id)
	if response.StatusCode == 204 {
		err = d.Set("result", true)
	} else if response.StatusCode == 404 {
		err = d.Set("result", false)
	}
	if err != nil {
		return diag.Errorf("error set result filed: %s", err)
	}
	return nil
}
