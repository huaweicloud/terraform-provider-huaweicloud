package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v2/{project_id}/desktop-pools/statistics/by-users
func DataSourceUserDesktopPoolAssociations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserDesktopPoolAssociationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the user desktop pool associations are located.`,
			},

			// Required parameters.
			// Although this parameter is optional in the documentation, in fact, if it is not defined, an empty list will be returned.
			"user_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of user IDs to be queried.`,
			},

			// Attributes.
			"associations": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        userDesktopPoolAssociationSchema(),
				Description: `The list of user associations with desktop pools.`,
			},
		},
	}
}

func userDesktopPoolAssociationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the user.`,
			},
			"desktop_pools": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        userDesktopPoolAssociationDesktopPoolDetailsSchema(),
				Description: `The list of desktop pools associated with the user.`,
			},
		},
	}
}

func userDesktopPoolAssociationDesktopPoolDetailsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the desktop pool.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the desktop pool.`,
			},
			"is_attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether a desktop is assigned.`,
			},
		},
	}
}

func buildUserDesktopPoolAssociationsQueryParams(d *schema.ResourceData) string {
	res := ""

	userIds := d.Get("user_ids").([]interface{})
	for _, userId := range userIds {
		res = fmt.Sprintf("%s&user_ids=%v", res, userId)
	}

	return res
}

func listUserDesktopPoolAssociations(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/desktop-pools/statistics/by-users?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildUserDesktopPoolAssociationsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		users := utils.PathSearch("users", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, users...)
		if len(users) < limit {
			break
		}
		offset += len(users)
	}

	return result, nil
}

func flattenDesktopPoolAssociations(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, pool := range items {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("desktop_pool_id", pool, nil),
			"name":        utils.PathSearch("desktop_pool_name", pool, nil),
			"is_attached": utils.PathSearch("is_attached", pool, nil),
		})
	}

	return result
}

func flattenUserDesktopPoolAssociations(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"user_id": utils.PathSearch("user_id", item, nil),
			"desktop_pools": flattenDesktopPoolAssociations(utils.PathSearch("desktop_pools", item,
				make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceUserDesktopPoolAssociationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	associations, err := listUserDesktopPoolAssociations(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace user desktop pool associations: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("associations", flattenUserDesktopPoolAssociations(associations)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
