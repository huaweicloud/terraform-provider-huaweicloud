package identitycenter

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IDENTITYCENTER GET /v1/identity-stores/{identity_store_id}/users/{user_id}/sessions
func DataSourceIdentityCenterUserSessions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCenterUserSessionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"session_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"session_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"session_not_valid_after": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_agent": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type UserSessionsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newUserSessionsDSWrapper(d *schema.ResourceData, meta interface{}) *UserSessionsDSWrapper {
	return &UserSessionsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceIdentityCenterUserSessionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newUserSessionsDSWrapper(d, meta)
	sessions, err := wrapper.getUserSessions()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.getUserSessionsToSchema(sessions)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func (w *UserSessionsDSWrapper) getUserSessions() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "identitystore")
	if err != nil {
		return nil, err
	}

	uri := "/v1/identity-stores/{identity_store_id}/users/{user_id}/sessions"
	uri = strings.ReplaceAll(uri, "{identity_store_id}", w.Get("identity_store_id").(string))
	uri = strings.ReplaceAll(uri, "{user_id}", w.Get("user_id").(string))

	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *UserSessionsDSWrapper) getUserSessionsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("session_list", schemas.SliceToList(body.Get("session_list"),
			func(sessions gjson.Result) any {
				return map[string]any{
					"ip_address":              sessions.Get("ip_address").Value(),
					"session_id":              sessions.Get("session_id").Value(),
					"user_agent":              sessions.Get("user_agent").Value(),
					"creation_time":           w.setCreDate(sessions),
					"session_not_valid_after": w.setNotValidDate(sessions),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*UserSessionsDSWrapper) setCreDate(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339((data.Get("creation_time").Int())/1000, true)
}

func (*UserSessionsDSWrapper) setNotValidDate(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339((data.Get("session_not_valid_after").Int())/1000, true)
}
