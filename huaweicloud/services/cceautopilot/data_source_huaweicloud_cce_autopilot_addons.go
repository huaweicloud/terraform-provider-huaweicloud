package cceautopilot

import (
	"context"

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

func DataSourceCceAutopilotAddons() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCceAutopilotAddonsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster.`,
			},
			"addon_template_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The template name of addon.`,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alias": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"annotations": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"creation_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"labels": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"update_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"spec": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addon_template_labels": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"addon_template_logo": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"addon_template_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"addon_template_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"current_version": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     itemsStatusCurrentVersionElem(),
									},
									"is_rollbackable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"previous_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_versions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func itemsStatusCurrentVersionElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"input": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"stable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_version": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"translate": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"update_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type AddonsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newAddonsDSWrapper(d *schema.ResourceData, meta interface{}) *AddonsDSWrapper {
	return &AddonsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCceAutopilotAddonsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newAddonsDSWrapper(d, meta)
	lisAutAddInsRst, err := wrapper.ListAutopilotAddons()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAutopilotAddonsToSchema(lisAutAddInsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /autopilot/v3/addons
func (w *AddonsDSWrapper) ListAutopilotAddons() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/autopilot/v3/addons"
	params := map[string]any{
		"cluster_id":          w.Get("cluster_id"),
		"addon_template_name": w.Get("addon_template_name"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("instances", "offset", "limit", 100).
		Request().
		Result()
}

func (w *AddonsDSWrapper) listAutopilotAddonsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("items", schemas.SliceToList(body.Get("items"),
			func(items gjson.Result) any {
				return map[string]any{
					"metadata": schemas.SliceToList(items.Get("metadata"),
						func(metadata gjson.Result) any {
							return map[string]any{
								"alias":              metadata.Get("alias").Value(),
								"annotations":        schemas.MapToStrMap(metadata.Get("annotations")),
								"creation_timestamp": metadata.Get("creationTimestamp").Value(),
								"labels":             schemas.MapToStrMap(metadata.Get("labels")),
								"name":               metadata.Get("name").Value(),
								"uid":                metadata.Get("uid").Value(),
								"update_timestamp":   metadata.Get("updateTimestamp").Value(),
							}
						},
					),
					"spec": schemas.SliceToList(items.Get("spec"),
						func(spec gjson.Result) any {
							return map[string]any{
								"addon_template_labels": schemas.SliceToStrList(spec.Get("addonTemplateLabels")),
								"addon_template_logo":   spec.Get("addonTemplateLogo").Value(),
								"addon_template_name":   spec.Get("addonTemplateName").Value(),
								"addon_template_type":   spec.Get("addonTemplateType").Value(),
								"cluster_id":            spec.Get("clusterID").Value(),
								"description":           spec.Get("description").Value(),
								"values":                schemas.SliceToStrList(spec.Get("values")),
								"version":               spec.Get("version").Value(),
							}
						},
					),
					"status": schemas.SliceToList(items.Get("status"),
						func(status gjson.Result) any {
							return map[string]any{
								"reason":           status.Get("Reason").Value(),
								"current_version":  w.setIteStaCurVersion(status),
								"is_rollbackable":  status.Get("isRollbackable").Value(),
								"message":          status.Get("message").Value(),
								"previous_version": status.Get("previousVersion").Value(),
								"status":           status.Get("status").Value(),
								"target_versions":  schemas.SliceToStrList(status.Get("targetVersions")),
							}
						},
					),
				}
			},
		)),
	)

	return mErr.ErrorOrNil()
}

func (*AddonsDSWrapper) setIteStaCurVersion(status gjson.Result) any {
	return schemas.SliceToList(status.Get("currentVersion"), func(currentVersion gjson.Result) any {
		return map[string]any{
			"creation_timestamp": currentVersion.Get("creationTimestamp").Value(),
			"input":              currentVersion.Get("input").Value(),
			"stable":             currentVersion.Get("stable").Value(),
			"support_versions": schemas.SliceToList(currentVersion.Get("supportVersions"),
				func(supportVersions gjson.Result) any {
					return map[string]any{
						"category":        schemas.SliceToStrList(supportVersions.Get("category")),
						"cluster_type":    supportVersions.Get("clusterType").Value(),
						"cluster_version": schemas.SliceToStrList(supportVersions.Get("clusterVersion")),
					}
				},
			),
			"translate":        currentVersion.Get("translate").Value(),
			"update_timestamp": currentVersion.Get("updateTimestamp").Value(),
			"version":          currentVersion.Get("version").Value(),
		}
	})
}
