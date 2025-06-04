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

// @API CCI GET /apis/cci/v2/namespaces/{namespace}/horizontalpodautoscalers
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/horizontalpodautoscalers/{name}
func DataSourceV2HPAs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2HPAsRead,

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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hpas": {
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
						"behavior": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scale_down": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     hpasBehaviorSchema(),
									},
									"scale_up": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     hpasBehaviorSchema(),
									},
								},
							},
						},
						"max_replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"metrics": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     hpasMetricsSchema(),
						},
						"min_replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scale_target_ref": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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
						"status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"last_transition_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"reason": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"current_metrics": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     hpasMetricsSchema(),
									},
									"current_replicas": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"desired_replicas": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"last_scale_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"observed_generation": {
										Type:     schema.TypeInt,
										Computed: true,
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

func hpasBehaviorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"select_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stabilization_window_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	return &sc
}

func hpasMetricsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"container_resource": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasContainerResourceMetricSourceSchema(),
			},
			"external": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasExternalMetricSourceSchema(),
			},
			"object": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasObjectmetricSourceSchema(),
			},
			"pods": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasPodMetricSourceSchema(),
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasResourceMetricSourceSchema(),
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func hpasContainerResourceMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"container": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasMetricTargetSchema(),
			},
		},
	}
	return &sc
}

func hpasMetricTargetSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"average_utilization": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"average_value": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func hpasExternalMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasMetricIdentifierSchema(),
			},
			"target": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasMetricTargetSchema(),
			},
		},
	}
	return &sc
}

func hpasMetricIdentifierSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
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
										Type:     schema.TypeMap,
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
		},
	}

	return &sc
}

func hpasObjectmetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"described_object": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metric": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasMetricIdentifierSchema(),
			},
			"target": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasMetricTargetSchema(),
			},
		},
	}
	return &sc
}

func hpasPodMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasMetricIdentifierSchema(),
			},
			"target": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasMetricTargetSchema(),
			},
		},
	}
	return &sc
}

func hpasResourceMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hpasMetricTargetSchema(),
			},
		},
	}
	return &sc
}

func dataSourceV2HPAsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	namespace := d.Get("namespace").(string)
	hpas := make([]interface{}, 0)
	if name, ok := d.GetOk("name"); ok {
		resp, err := GetV2HPA(client, namespace, name.(string))
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); !ok {
				return diag.Errorf("error getting the HPA from the server: %s", err)
			}
		}
		if resp != nil {
			hpas = append(hpas, resp)
		}
	} else {
		resp, err := listHPAs(client, namespace)
		if err != nil {
			return diag.Errorf("error getting the HPA list from the server: %s", err)
		}
		hpas = utils.PathSearch("items", resp, make([]interface{}, 0)).([]interface{})
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("hpas", flattenHpas(hpas)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHpas(hpas []interface{}) []interface{} {
	if len(hpas) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(hpas))
	for _, v := range hpas {
		rst = append(rst, map[string]interface{}{
			"name":               utils.PathSearch("metadata.name", v, nil),
			"namespace":          utils.PathSearch("metadata.namespace", v, nil),
			"creation_timestamp": utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":   utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                utils.PathSearch("metadata.uid", v, nil),
			"behavior":           flattenSpecBehavior(utils.PathSearch("spec.behavior", v, nil)),
			"max_replicas":       utils.PathSearch("spec.maxReplicas", v, nil),
			"metrics":            flattenMetrics(utils.PathSearch("spec.metrics", v, make([]interface{}, 0)).([]interface{})),
			"min_replicas":       utils.PathSearch("spec.minReplicas", v, nil),
			"scale_target_ref":   flattenSpecScaleTargetRef(utils.PathSearch("spec.scaleTargetRef", v, nil)),
			"status":             flattenHPAStatus(utils.PathSearch("status", v, nil)),
		})
	}
	return rst
}

func listHPAs(client *golangsdk.ServiceClient, namespace string) (interface{}, error) {
	listHPAsHttpUrl := "apis/cci/v2/namespaces/{namespace}/horizontalpodautoscalers"
	listHPAsPath := client.Endpoint + listHPAsHttpUrl
	listHPAsPath = strings.ReplaceAll(listHPAsPath, "{namespace}", namespace)
	listHPAsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listHPAsResp, err := client.Request("GET", listHPAsPath, &listHPAsOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(listHPAsResp)
}
