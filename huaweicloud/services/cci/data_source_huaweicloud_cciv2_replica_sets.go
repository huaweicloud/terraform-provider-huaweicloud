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

// @API CCI GET /apis/cci/v2/namespaces/{namespace}/replicasets
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/replicasets/{name}
func DataSourceV2ReplicaSets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2ReplicaSetsRead,

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
				Computed: true,
			},
			"replica_sets": {
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
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finalizers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
						"min_ready_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"selector": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     replicaSetsLabelSelectorSchema(),
						},
						"template": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     replicaSetsTemplateSchema(),
						},
						"status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     replicaSetsStatusSchema(),
						},
					},
				},
			},
		},
	}
}

func replicaSetsStatusSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"available_replicas": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fully_labeled_replicas": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"observed_generation": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ready_replicas": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replicas": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replicaSetsStatusConditionsSchema(),
			},
		},
	}
	return &sc
}

func replicaSetsStatusConditionsSchema() *schema.Resource {
	sc := schema.Resource{
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
	}
	return &sc
}

func replicaSetsLabelSelectorSchema() *schema.Resource {
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
				Elem:     replicaSetsMatchExpressionsSchema(),
			},
		},
	}
	return &sc
}

func replicaSetsMatchExpressionsSchema() *schema.Resource {
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

func replicaSetsTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replicaSetsTemplateMetadataSchema(),
			},
			"spec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replicaSetsTemplateSpecSchema(),
			},
		},
	}
	return &sc
}

func replicaSetsTemplateMetadataSchema() *schema.Resource {
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

func replicaSetsTemplateSpecSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"containers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"env": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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
					},
				},
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
							Elem:     replicaSetsNodeAffinitySchema(),
						},
						"pod_anti_affinity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     replicaSetsPodAntiAffinitySchema(),
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

func replicaSetsNodeAffinitySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"required_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replicaSetsNodeSelectorSchema(),
			},
		},
	}
	return &sc
}

func replicaSetsNodeSelectorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_selector_terms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replicaSetsNodeSelectorTermSchema(),
			},
		},
	}

	return &sc
}

func replicaSetsNodeSelectorTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_expressions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replicaSetsNodeSelectorRequirementSchema(),
			},
		},
	}

	return &sc
}

func replicaSetsNodeSelectorRequirementSchema() *schema.Resource {
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

func replicaSetsPodAntiAffinitySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"preferred_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replicaSetsWeightedPodAffinityTermSchema(),
			},
			"required_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replicaSetsPodAffinityTermSchema(),
			},
		},
	}
	return &sc
}

func replicaSetsPodAffinityTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"label_selector": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replicaSetsLabelSelectorSchema(),
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

func replicaSetsWeightedPodAffinityTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"pod_affinity_term": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replicaSetsPodAffinityTermSchema(),
			},
			"weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceV2ReplicaSetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	namespace := d.Get("namespace").(string)
	replicaSets := make([]interface{}, 0)
	if name, ok := d.GetOk("name"); ok {
		resp, err := GetReplicaSet(client, namespace, name.(string))
		if err != nil {
			return diag.Errorf("error getting the replica set from the server: %s", err)
		}
		replicaSets = append(replicaSets, resp)
	} else {
		resp, err := listReplicaSets(client, namespace)
		if err != nil {
			return diag.Errorf("error getting the replica sets from the server: %s", err)
		}
		replicaSets = utils.PathSearch("items", resp, make([]interface{}, 0)).([]interface{})
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("replica_sets", flattenReplicaSets(replicaSets)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenReplicaSets(replicaSets []interface{}) []interface{} {
	if len(replicaSets) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(replicaSets))
	for _, v := range replicaSets {
		rst = append(rst, map[string]interface{}{
			"name":               utils.PathSearch("metadata.name", v, nil),
			"api_version":        utils.PathSearch("apiVersion", v, nil),
			"kind":               utils.PathSearch("kind", v, nil),
			"annotations":        utils.PathSearch("metadata.annotations", v, nil),
			"labels":             utils.PathSearch("metadata.labels", v, nil),
			"creation_timestamp": utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":   utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                utils.PathSearch("metadata.uid", v, nil),
			"min_ready_seconds":  utils.PathSearch("spec.min_ready_seconds", v, nil),
			"replicas":           utils.PathSearch("spec.replicas", v, nil),
			"selector":           flattenLabelSelector(utils.PathSearch("spec.selector", v, nil)),
			"template":           flattenSpecTemplate(utils.PathSearch("spec.template", v, nil)),
			"status":             flattenReplicaSetsStatus(utils.PathSearch("status", v, nil)),
		})
	}
	return rst
}

func flattenReplicaSetsStatus(status interface{}) []map[string]interface{} {
	if status == nil {
		return nil
	}

	conditions := utils.PathSearch("conditions", status, make([]interface{}, 0)).([]interface{})
	rst := []map[string]interface{}{
		{
			"available_replicas":     utils.PathSearch("availableReplicas", status, nil),
			"fully_labeled_replicas": utils.PathSearch("fullyLabeledReplicas", status, nil),
			"observed_generation":    utils.PathSearch("observedGeneration", status, nil),
			"ready_replicas":         utils.PathSearch("readyReplicas", status, nil),
			"replicas":               utils.PathSearch("replicas", status, nil),
			"conditions":             flattenReplicaSetsStatusConditions(conditions),
		},
	}

	return rst
}

func flattenReplicaSetsStatusConditions(conditions []interface{}) []interface{} {
	if len(conditions) == 0 {
		return nil
	}

	rst := make([]interface{}, len(conditions))
	for i, v := range conditions {
		rst[i] = map[string]interface{}{
			"type":                 utils.PathSearch("type", v, nil),
			"status":               utils.PathSearch("status", v, nil),
			"last_transition_time": utils.PathSearch("lastTransitionTime", v, nil),
			"reason":               utils.PathSearch("reason", v, nil),
			"message":              utils.PathSearch("message", v, nil),
		}
	}

	return rst
}

func GetReplicaSet(client *golangsdk.ServiceClient, namespace, name string) (interface{}, error) {
	getReplicaSetHttpUrl := "apis/cci/v2/namespaces/{namespace}/replicasets/{name}"
	getReplicaSetPath := client.Endpoint + getReplicaSetHttpUrl
	getReplicaSetPath = strings.ReplaceAll(getReplicaSetPath, "{namespace}", namespace)
	getReplicaSetPath = strings.ReplaceAll(getReplicaSetPath, "{name}", name)
	getReplicaSetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getReplicaSetResp, err := client.Request("GET", getReplicaSetPath, &getReplicaSetOpt)
	if err != nil {
		return getReplicaSetResp, err
	}

	return utils.FlattenResponse(getReplicaSetResp)
}

func listReplicaSets(client *golangsdk.ServiceClient, namespace string) (interface{}, error) {
	listReplicaSetsHttpUrl := "apis/cci/v2/namespaces/{namespace}/replicasets"
	listReplicaSetsPath := client.Endpoint + listReplicaSetsHttpUrl
	listReplicaSetsPath = strings.ReplaceAll(listReplicaSetsPath, "{namespace}", namespace)
	listReplicaSetsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listReplicaSetsResp, err := client.Request("GET", listReplicaSetsPath, &listReplicaSetsOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(listReplicaSetsResp)
}
