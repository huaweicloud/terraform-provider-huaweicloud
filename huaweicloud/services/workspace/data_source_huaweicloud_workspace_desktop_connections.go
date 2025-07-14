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

// @API Workspace GET /v2/{project_id}/connect-desktops
func DataSourceDesktopConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDesktopConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the desktop connections are located.`,
			},
			"user_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of desktop users to be queried.`,
			},
			"connect_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The connection status of the desktop.`,
			},
			"desktop_connections": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopConnectionSchema(),
				Description: `The list of desktop connections that match the query parameters.`,
			},
		},
	}
}

func desktopConnectionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the desktop.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the desktop.`,
			},
			"connect_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection status of the desktop.`,
			},
			"attach_users": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopAttachUserSchema(),
				Description: `The list of users or user groups attached to the desktop.`,
			},
		},
	}
}

func desktopAttachUserSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the user or user group.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the user or user group.`,
			},
			"user_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user group of the desktop user.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the user or user group.`,
			},
		},
	}
}

func flattenDesktopConnections(desktops []interface{}) []interface{} {
	if len(desktops) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(desktops))
	for _, desktop := range desktops {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("desktop_id", desktop, nil),
			"name":           utils.PathSearch("desktop_name", desktop, nil),
			"connect_status": utils.PathSearch("connect_status", desktop, nil),
			"attach_users":   flattenDesktopAttachUsers(utils.PathSearch("attach_users", desktop, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenDesktopAttachUsers(attachUsers []interface{}) []interface{} {
	if len(attachUsers) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(attachUsers))
	for _, user := range attachUsers {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", user, nil),
			"name":       utils.PathSearch("name", user, nil),
			"user_group": utils.PathSearch("user_group", user, nil),
			"type":       utils.PathSearch("type", user, nil),
		})
	}
	return result
}

func buildListDesktopConnectionsParams(d *schema.ResourceData) string {
	params := ""
	if v, ok := d.GetOk("user_names"); ok {
		params = fmt.Sprintf("%s&user_names=%v", params, v)
	}
	if v, ok := d.GetOk("connect_status"); ok {
		params = fmt.Sprintf("%s&connect_status=%v", params, v)
	}
	return params
}

func listDesktopConnections(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/connect-desktops?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildListDesktopConnectionsParams(d)

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

		desktopConnections := utils.PathSearch("desktops", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, desktopConnections...)
		if len(desktopConnections) < limit {
			return result, nil
		}
		offset += len(desktopConnections)
	}
}

func dataSourceDesktopConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating WorkSpace client: %s", err)
	}

	desktops, err := listDesktopConnections(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace desktop connections: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("desktop_connections", flattenDesktopConnections(desktops)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
