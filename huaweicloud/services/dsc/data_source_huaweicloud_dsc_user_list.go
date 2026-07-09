package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v2/{project_id}/sec-ops/users
func DataSourceDscUserList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscUserListRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     userInfoSchema(),
			},
		},
	}
}

func userInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDscUserListRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "dsc"
		httpUrl  = "v2/{project_id}/sec-ops/users"
		limit    = 1000
		offset   = 0
		allUsers = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	for {
		currentPath := fmt.Sprintf("%s?offset=%d&limit=%d", requestPath, offset, limit)
		requestOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		resp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving DSC user list: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		usersRaw := utils.PathSearch("user_list", respBody, make([]interface{}, 0)).([]interface{})
		allUsers = append(allUsers, usersRaw...)
		if len(usersRaw) < limit {
			break
		}

		offset += len(usersRaw)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("user_list", flattenUserList(allUsers)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenUserList(users []interface{}) []interface{} {
	if len(users) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(users))
	for _, item := range users {
		result = append(result, map[string]interface{}{
			"user_id":   utils.PathSearch("user_id", item, nil),
			"user_name": utils.PathSearch("user_name", item, nil),
		})
	}

	return result
}
