package cceautopilot

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

func DataSourceCCEAutopilotReleases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEAutopilotReleasesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"chart_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"releases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chart_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chart_public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"chart_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type ReleasesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCCEAutopilotReleasesDSWrapper(d *schema.ResourceData, meta interface{}) *ReleasesDSWrapper {
	return &ReleasesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCCEAutopilotReleasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCCEAutopilotReleasesDSWrapper(d, meta)
	lisAutRelRst, err := wrapper.ListAutopilotReleases()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAutopilotReleasesToSchema(lisAutRelRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /autopilot/cam/v3/clusters/{cluster_id}/releases
func (w *ReleasesDSWrapper) ListAutopilotReleases() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/autopilot/cam/v3/clusters/{cluster_id}/releases"
	uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))
	params := map[string]any{
		"chart_id":  w.Get("chart_id"),
		"namespace": w.Get("namespace"),
	}

	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *ReleasesDSWrapper) listAutopilotReleasesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("releases", schemas.SliceToList(*body,
			func(release gjson.Result) any {
				return map[string]any{
					"chart_name":         release.Get("chart_name").Value(),
					"chart_public":       release.Get("chart_public").Value(),
					"chart_version":      release.Get("chart_version").Value(),
					"cluster_id":         release.Get("cluster_id").Value(),
					"cluster_name":       release.Get("cluster_name").Value(),
					"create_at":          release.Get("create_at").Value(),
					"description":        release.Get("description").Value(),
					"name":               release.Get("name").Value(),
					"namespace":          release.Get("namespace").Value(),
					"parameters":         release.Get("parameters").Value(),
					"resources":          release.Get("resources").Value(),
					"status":             release.Get("status").Value(),
					"status_description": release.Get("status_description").Value(),
					"update_at":          release.Get("update_at").Value(),
					"values":             release.Get("values").Value(),
					"version":            release.Get("version").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
