package cci

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI GET /apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims
func DataSourceV2PersistentVolumeClaims() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2PersistentVolumeClaimsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pvcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"access_modes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limits": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"requests": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"selector": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_expressions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"values": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"match_labels": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"storage_class_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"valume_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finalizers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceV2PersistentVolumeClaimsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}
	listPvcsHttpUrl := "apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims"
	listPvcsPath := client.Endpoint + listPvcsHttpUrl
	listPvcsPath = strings.ReplaceAll(listPvcsPath, "{namespace}", d.Get("namespace").(string))
	listPvcsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listPvcsResp, err := client.Request("GET", listPvcsPath, &listPvcsOpt)
	if err != nil {
		return diag.Errorf("error getting CCI pvc list: %s", err)
	}

	listPvcsRespBody, err := utils.FlattenResponse(listPvcsResp)
	if err != nil {
		return diag.Errorf("error retrieving CCI pvcs: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	pvcs := utils.PathSearch("items", listPvcsRespBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("pvcs", flattenPvcs(pvcs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPvcs(pvcs []interface{}) []interface{} {
	if len(pvcs) == 0 {
		return nil
	}

	rst := make([]interface{}, len(pvcs))
	for i, v := range pvcs {
		rst[i] = map[string]interface{}{
			"name":               utils.PathSearch("metadata.name", v, nil),
			"namespace":          utils.PathSearch("metadata.namespace", v, nil),
			"annotations":        utils.PathSearch("metadata.annotations", v, nil),
			"labels":             utils.PathSearch("metadata.labels", v, nil),
			"creation_timestamp": utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":   utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                utils.PathSearch("metadata.uid", v, nil),
			"finalizers":         utils.PathSearch("metadata.finalizers", v, nil),
			"status":             utils.PathSearch("status.phase", v, nil),
			"access_modes":       utils.PathSearch("spec.accessModes", v, nil),
			"storage_class_name": utils.PathSearch("spec.storageClassName", v, nil),
			"volume_mode":        utils.PathSearch("spec.volumeMode", v, nil),
			"valume_name":        utils.PathSearch("spec.valumeName", v, nil),
			"selector":           flattenPvcsLabelSelector(utils.PathSearch("spec.selector", v, nil)),
			"resources":          flattenPvcResources(utils.PathSearch("spec.resources", v, nil)),
		}
	}
	return rst
}

func flattenPvcsLabelSelector(selector interface{}) []interface{} {
	if selector == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"match_labels": utils.PathSearch("matchLabels", selector, nil),
		"match_expressions": flattenPvcsMatchExpressions(
			utils.PathSearch("matchExpressions", selector, make([]interface{}, 0)).([]interface{})),
	})

	return rst
}

func flattenPvcsMatchExpressions(matchExpressions []interface{}) []interface{} {
	if len(matchExpressions) == 0 {
		return nil
	}

	rst := make([]interface{}, len(matchExpressions))
	for i, v := range matchExpressions {
		rst[i] = map[string]interface{}{
			"key":      utils.PathSearch("key", v, nil),
			"operator": utils.PathSearch("operator", v, nil),
			"values":   utils.PathSearch("values", v, nil),
		}
	}
	return rst
}

func flattenPvcResources(resources interface{}) []map[string]interface{} {
	if resources == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"limits":   utils.PathSearch("limits", resources, nil),
			"requests": utils.PathSearch("requests", resources, nil),
		},
	}
}
