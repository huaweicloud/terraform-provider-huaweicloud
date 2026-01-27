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
)

func DataSourceCceAutopilotClusterUpgradeInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCceAutopilotClusterUpgradeInfoRead,

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
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"update_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_timestamp": {
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
						"last_upgrade_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"phase": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"progress": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"completion_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"version_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"release": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"patch": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"suggest_patch": {
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
						"upgrade_feature_gates": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"support_upgrade_page_v4": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"phase": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"completion_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type ClusterUpgradeInfoDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCceAutopilotClusterUpgradeInfoDSWrapper(d *schema.ResourceData, meta interface{}) *ClusterUpgradeInfoDSWrapper {
	return &ClusterUpgradeInfoDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCceAutopilotClusterUpgradeInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCceAutopilotClusterUpgradeInfoDSWrapper(d, meta)
	showAutCluUpgInfRst, err := wrapper.ShowAutopilotClusterUpgradeInfo()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showAutopilotClusterUpgradeInfoToSchema(showAutCluUpgInfRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /autopilot/v3/projects/{project_id}/clusters/{cluster_id}/upgradeinfo
func (w *ClusterUpgradeInfoDSWrapper) ShowAutopilotClusterUpgradeInfo() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/upgradeinfo"
	uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *ClusterUpgradeInfoDSWrapper) showAutopilotClusterUpgradeInfoToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("kind", body.Get("kind").Value()),
		d.Set("api_version", body.Get("apiVersion").Value()),
		d.Set("metadata", schemas.SliceToList(body.Get("metadata"),
			func(metadata gjson.Result) any {
				return map[string]any{
					"uid":                metadata.Get("uid").Value(),
					"name":               metadata.Get("name").Value(),
					"labels":             schemas.MapToStrMap(metadata.Get("labels")),
					"annotations":        schemas.MapToStrMap(metadata.Get("annotations")),
					"update_timestamp":   metadata.Get("updateTimestamp").Value(),
					"creation_timestamp": metadata.Get("creationTimestamp").Value(),
				}
			},
		)),
		d.Set("spec", schemas.SliceToList(body.Get("spec"),
			func(spec gjson.Result) any {
				return map[string]any{
					"last_upgrade_info": schemas.SliceToList(spec.Get("lastUpgradeInfo"),
						func(lastinfo gjson.Result) any {
							return map[string]any{
								"phase":           lastinfo.Get("phase").Value(),
								"progress":        lastinfo.Get("progress").Value(),
								"completion_time": lastinfo.Get("completionTime").Value(),
							}
						},
					),
					"version_info": schemas.SliceToList(spec.Get("versionInfo"),
						func(versioninfo gjson.Result) any {
							return map[string]any{
								"release":         versioninfo.Get("release").Value(),
								"patch":           versioninfo.Get("patch").Value(),
								"suggest_patch":   versioninfo.Get("suggestPatch").Value(),
								"target_versions": schemas.SliceToStrList(versioninfo.Get("targetVersions")),
							}
						},
					),
					"upgrade_feature_gates": schemas.SliceToList(spec.Get("upgradeFeatureGates"),
						func(featuregates gjson.Result) any {
							return map[string]any{
								"support_upgrade_page_v4": featuregates.Get("supportUpgradePageV4").Value(),
							}
						},
					),
				}
			},
		)),
		d.Set("status", schemas.SliceToList(body.Get("status"),
			func(status gjson.Result) any {
				return map[string]any{
					"phase":           status.Get("phase").Value(),
					"progress":        status.Get("progress").Value(),
					"completion_time": status.Get("completionTime").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
