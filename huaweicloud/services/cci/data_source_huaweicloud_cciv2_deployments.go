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

// @API CCI GET /apis/cci/v2/namespaces/{namespace}/deployments
func DataSourceV2Deployments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2DeploymentsRead,

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
			"deployments": {
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
						"replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_ready_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"progress_deadline_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"selector": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     deploymentsLabelSelectorSchema(),
						},
						"template": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     deploymentsTemplateSchema(),
						},
						"strategy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rolling_update": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
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
						"resource_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"generation": {
							Type:     schema.TypeInt,
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
									"observed_generation": {
										Type:     schema.TypeInt,
										Computed: true,
									},
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
												"last_update_time": {
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func deploymentsLabelSelectorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"match_expressions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     deploymentsMatchExpressionsSchema(),
			},
		},
	}
	return &sc
}

func deploymentsMatchExpressionsSchema() *schema.Resource {
	sc := schema.Resource{
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
	}

	return &sc
}

func deploymentsTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     deploymentsTemplateMetadataSchema(),
			},
			"spec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     deploymentsTemplateSpecSchema(),
			},
		},
	}
	return &sc
}

func deploymentsTemplateMetadataSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
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
		},
	}
	return &sc
}

func deploymentsTemplateSpecSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"containers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersSchema(),
			},
			"dns_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_deadline_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"overhead": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"restart_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheduler_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"set_hostname_as_pqdn": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"share_process_namespace": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"termination_grace_period_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"affinity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_affinity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     deploymentsNodeAffinitySchema(),
						},
						"pod_anti_affinity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     deploymentsPodAntiAffinitySchema(),
						},
					},
				},
			},
			"image_pull_secrets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
	return &sc
}

func deploymentsNodeAffinitySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"required_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     deploymentsNodeSelectorSchema(),
			},
		},
	}
	return &sc
}

func deploymentsNodeSelectorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_selector_terms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     deploymentsNodeSelectorTermSchema(),
			},
		},
	}

	return &sc
}

func deploymentsNodeSelectorTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_expressions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     deploymentsNodeSelectorRequirementSchema(),
			},
		},
	}

	return &sc
}

func deploymentsNodeSelectorRequirementSchema() *schema.Resource {
	sc := schema.Resource{
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
	}

	return &sc
}

func deploymentsPodAntiAffinitySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"preferred_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     deploymentsWeightedPodAffinityTermSchema(),
			},
			"required_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     deploymentsPodAffinityTermSchema(),
			},
		},
	}
	return &sc
}

func deploymentsPodAffinityTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"label_selector": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     labelSelectorSchema(),
			},
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"topology_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func deploymentsWeightedPodAffinityTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"pod_affinity_term": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podAffinityTermSchema(),
			},
			"weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceV2DeploymentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}
	listDeploymentsHttpUrl := "apis/cci/v2/namespaces/{namespace}/deployments"
	listDeploymentsPath := client.Endpoint + listDeploymentsHttpUrl
	listDeploymentsPath = strings.ReplaceAll(listDeploymentsPath, "{namespace}", d.Get("namespace").(string))
	listDeploymentsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listDeploymentsResp, err := client.Request("GET", listDeploymentsPath, &listDeploymentsOpt)
	if err != nil {
		return diag.Errorf("error getting CCI deployments list: %s", err)
	}

	listDeploymentsRespBody, err := utils.FlattenResponse(listDeploymentsResp)
	if err != nil {
		return diag.Errorf("error retrieving CCI deployments: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	deployments := utils.PathSearch("items", listDeploymentsRespBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("deployments", flattenDeployments(deployments)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDeployments(deployments []interface{}) []interface{} {
	if len(deployments) == 0 {
		return nil
	}

	rst := make([]interface{}, len(deployments))
	for i, v := range deployments {
		rst[i] = map[string]interface{}{
			"name":                      utils.PathSearch("metadata.name", v, nil),
			"namespace":                 utils.PathSearch("metadata.namespace", v, nil),
			"annotations":               utils.PathSearch("metadata.annotations", v, nil),
			"creation_timestamp":        utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":          utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                       utils.PathSearch("metadata.uid", v, nil),
			"generation":                utils.PathSearch("metadata.generation", v, nil),
			"replicas":                  int(utils.PathSearch("spec.replicas", v, float64(0)).(float64)),
			"min_ready_seconds":         int(utils.PathSearch("spec.minReadySeconds", v, float64(0)).(float64)),
			"progress_deadline_seconds": int(utils.PathSearch("spec.progressDeadlineSeconds", v, float64(0)).(float64)),
			"selector":                  flattenLabelSelector(utils.PathSearch("spec.selector", v, nil)),
			"template":                  flattenSpecTemplate(utils.PathSearch("spec.template", v, nil)),
			"strategy":                  flattenSpecStrategy(utils.PathSearch("spec.strategy", v, nil)),
			"status":                    flattenDeploymentStatus(utils.PathSearch("status", v, nil)),
		}
	}
	return rst
}
