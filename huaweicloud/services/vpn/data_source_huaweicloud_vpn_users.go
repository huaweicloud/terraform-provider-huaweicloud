// Generated by PMS #377
package vpn

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceVpnUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnUsersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"vpn_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a VPN server.`,
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The user list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The username.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user description.`,
						},
						"user_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the user group to which a user belongs.`,
						},
						"user_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the user group to which a user belongs.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time.`,
						},
					},
				},
			},
		},
	}
}

type UsersDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newUsersDSWrapper(d *schema.ResourceData, meta interface{}) *UsersDSWrapper {
	return &UsersDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceVpnUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newUsersDSWrapper(d, meta)
	listVpnUsersRst, err := wrapper.ListVpnUsers()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listVpnUsersToSchema(listVpnUsersRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API VPN GET /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users
func (w *UsersDSWrapper) ListVpnUsers() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "vpn")
	if err != nil {
		return nil, err
	}

	uri := "/v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users"
	uri = strings.ReplaceAll(uri, "{vpn_server_id}", w.Get("vpn_server_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		MarkerPager("users", "page_info.next_marker", "marker").
		Request().
		Result()
}

func (w *UsersDSWrapper) listVpnUsersToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("users", schemas.SliceToList(body.Get("users"),
			func(users gjson.Result) any {
				return map[string]any{
					"id":              users.Get("id").Value(),
					"name":            users.Get("name").Value(),
					"description":     users.Get("description").Value(),
					"user_group_id":   users.Get("user_group_id").Value(),
					"user_group_name": users.Get("user_group_name").Value(),
					"created_at":      users.Get("created_at").Value(),
					"updated_at":      users.Get("updated_at").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}