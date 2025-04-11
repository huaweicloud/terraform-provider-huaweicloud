package cci

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var horizontalPodAutoscalerNonUpdatableParams = []string{"namespace", "name"}

// @API CCI POST /apis/cci/v2/namespace/{namespace}/horizontalpodautoscalers
// @API CCI GET /apis/cci/v2/namespace/{namespace}/horizontalpodautoscalers/{name}
// @API CCI PUT /apis/cci/v2/namespace/{namespace}/horizontalpodautoscalers/{name}
// @API CCI DELETE /apis/cci/v2/namespace/{namespace}/horizontalpodautoscalers/{name}
func ResourceV2HorizontalPodAutoscaler() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2HorizontalPodAutoscalerCreate,
		ReadContext:   resourceV2HorizontalPodAutoscalerRead,
		UpdateContext: resourceV2HorizontalPodAutoscalerUpdate,
		DeleteContext: resourceV2HorizontalPodAutoscalerDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2HorizontalPodAutoscalerImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(horizontalPodAutoscalerNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the CCI Image Snapshot.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the annotations of the CCI Image Snapshot.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the annotations of the CCI Image Snapshot.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the CCI Image Snapshot.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the CCI Image Snapshot.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the CCI Image Snapshot.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the CCI Image Snapshot.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the CCI Image Snapshot.`,
			},
			"behavior": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scale_down": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Elem:        horizontalPodAutoscalerBehaviorSchema(),
							Description: `Specifies the scale down of the behavior.`,
						},
						"scale_up": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Elem:        horizontalPodAutoscalerBehaviorSchema(),
							Description: `Specifies the scale up of the behavior.`,
						},
					},
				},
				Description: `Specifies the behavior of the CCI horizontal pod autoscaler.`,
			},
			"max_replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the upper limit for the number of replicas to which the autoscaler can scale up.`,
			},
			"metrics": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        horizontalPodAutoscalerMetricsSchema(),
				Description: `Specifies the metrics that can be used to calculate the desired replica count.`,
			},
			"min_replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the lower limit for the number of replicas to which the autoscaler can scale down`,
			},
			"scale_target_ref": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the API version of the referent.`,
						},
						"kind": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the kind of the referent.`,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the name of the referent.`,
						},
					},
				},
				Description: `Specifies the scale target.`,
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
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the conditions.`,
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Tthe status of the conditions.`,
									},
									"last_transition_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The last transition time of the conditions.`,
									},
									"reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The reason of the conditions.`,
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The message of the conditions.`,
									},
								},
							},
							Description: `The status.`,
						},
						"current_metrics": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        horizontalPodAutoscalerMetricsSchema(),
							Description: `The status.`,
						},
						"current_replicas": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The current replicas.`,
						},
						"desired_replicas": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The desired replicas.`,
						},
						"last_scale_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last scale time.`,
						},
						"observed_generation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The observed generation.`,
						},
					},
				},
				Description: `The status.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func horizontalPodAutoscalerBehaviorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"policies": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the window of time for which the policy should hold true.`,
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the type of the scaling policy.`,
						},
						"value": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the value, it contains the amount of change which is permitted by the policy.`,
						},
					},
				},
				Description: `Specifies the potential scaling policies which can be used during scaling.`,
			},
			"select_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies select policy that should be used.`,
			},
			"stabilization_window_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the seconds for which past recommendations should be considered while scaling up or scaling down.`,
			},
		},
	}

	return &sc
}

func horizontalPodAutoscalerMetricsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"container_resource": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        containerResourceMetricSourceSchema(),
				Description: `Specifies the container resource metric source.`,
			},
			"external": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        externalMetricSourceSchema(),
				Description: `Specifies the external metric resource.`,
			},
			"object": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        objectmetricSourceSchema(),
				Description: `Specifies the object metric resource.`,
			},
			"pods": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        podMetricSourceSchema(),
				Description: `Specifies the pod metric resource.`,
			},
			"resource": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceMetricSourceSchema(),
				Description: `Specifies the resource metric resource.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the seconds for which past recommendations should be considered while scaling up or scaling down.`,
			},
		},
	}

	return &sc
}

func containerResourceMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"container": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the container in the pods of the scaling target.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the resource in question.`,
			},
			"target": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        metricTargetSchema(),
				Description: `Specifies the container resource metric source.`,
			},
		},
	}
	return &sc
}

func metricTargetSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"average_utilization": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the target value of the resource metric across all elevant pods.`,
			},
			"average_value": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the average value of the resource.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the metric type, the value can be **Utilization**, **Value**, or **AverageValue**.`,
			},
			"value": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the value of the resource.`,
			},
		},
	}
	return &sc
}

func externalMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        metricIdentifierSchema(),
				Description: `Specifies the metric of external metric source.`,
			},
			"target": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        metricTargetSchema(),
				Description: `Specifies the target of external metric source.`,
			},
		},
	}
	return &sc
}

func metricIdentifierSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the given metric.`,
			},
			"selector": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_expressions": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `Specifies the label key that the selector applies to.`,
									},
									"operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `Specifies the operator represents a key relationship to a set of values.`,
									},
									"values": {
										Type:        schema.TypeMap,
										Optional:    true,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `Specifies the array of string values.`,
									},
								},
							},
							Description: `Specifies the match expressions of the label selector requirements.`,
						},
						"match_labels": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the match labels.`,
						},
					},
				},
				Description: `Specifies the metric of external metric source.`,
			},
		},
	}

	return &sc
}

func objectmetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"described_object": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the API version of the referent.`,
						},
						"kind": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the kind of the referent.`,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the name of the referent.`,
						},
					},
				},
				Description: `Specifies the container resource metric source.`,
			},
			"metric": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        metricIdentifierSchema(),
				Description: `Specifies the container resource metric source.`,
			},
			"target": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        metricTargetSchema(),
				Description: `Specifies the container resource metric source.`,
			},
		},
	}
	return &sc
}

func podMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        metricIdentifierSchema(),
				Description: `Specifies the container resource pod metric source.`,
			},
			"target": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        metricTargetSchema(),
				Description: `Specifies the container resource pod metric source.`,
			},
		},
	}
	return &sc
}

func resourceMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the resource in question.`,
			},
			"target": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        metricTargetSchema(),
				Description: `Specifies the container resource pod metric source.`,
			},
		},
	}
	return &sc
}

func resourceV2HorizontalPodAutoscalerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createHorizontalPodAutoscalerHttpUrl := "apis/cci/v2/namespace/{namespace}/horizontalpodautoscalers"
	createHorizontalPodAutoscalerPath := client.Endpoint + createHorizontalPodAutoscalerHttpUrl
	createHorizontalPodAutoscalerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createHorizontalPodAutoscalerOpt.JSONBody = utils.RemoveNil(buildCreateV2HorizontalPodAutoscalerParams(d))

	resp, err := client.Request("POST", createHorizontalPodAutoscalerPath, &createHorizontalPodAutoscalerOpt)
	if err != nil {
		return diag.Errorf("error creating CCI HorizontalPodAutoscaler: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if name == "" {
		return diag.Errorf("unable to find CCI HorizontalPodAutoscaler name or namespace from API response")
	}
	d.SetId(name)

	return resourceV2HorizontalPodAutoscalerRead(ctx, d, meta)
}

func buildCreateV2HorizontalPodAutoscalerParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": d.Get("name"),
		},
		"spec":   map[string]interface{}{},
		"status": map[string]interface{}{},
	}

	return bodyParams
}

func resourceV2HorizontalPodAutoscalerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	resp, err := GetHorizontalPodAutoscaler(client, ns, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 HorizontalPodAutoscaler")
	}

	mErr := multierror.Append(
		d.Set("namespace", utils.PathSearch("metadata.namespace", resp, nil)),
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("behavior", flattenSpecBehavior(utils.PathSearch("spec.behavior", resp, nil))),
		d.Set("max_replicas", utils.PathSearch("spec.maxReplicas", resp, nil)),
		d.Set("metrics", flattenMetrics(utils.PathSearch("spec.metrics", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("min_replicas", utils.PathSearch("spec.minReplicas", resp, nil)),
		d.Set("scale_target_ref", flattenSpecScaleTargetRef(utils.PathSearch("spec.scaleTargetRef", resp, nil))),
		d.Set("status", flattenHorizontalPodAutoscalerStatus(utils.PathSearch("status", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSpecScaleTargetRef(scaleTargetRef interface{}) []interface{} {
	if scaleTargetRef == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"api_version": utils.PathSearch("apiVersion", scaleTargetRef, nil),
		"kind":        utils.PathSearch("kind", scaleTargetRef, nil),
		"name":        utils.PathSearch("name", scaleTargetRef, nil),
	})

	return rst
}

func flattenSpecBehavior(behavior interface{}) []interface{} {
	if behavior == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"scale_down": flattenSpecBehaviorScale(utils.PathSearch("scaleDown", behavior, nil)),
		"scale_up":   flattenSpecBehaviorScale(utils.PathSearch("scaleUp", behavior, nil)),
	})

	return rst
}

func flattenSpecBehaviorScale(scale interface{}) []interface{} {
	if scale == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	policies := utils.PathSearch("policies", scale, make([]interface{}, 0)).([]interface{})
	rst = append(rst, map[string]interface{}{
		"policies":                     flattenSpecBehaviorScalePolicies(policies),
		"select_policy":                utils.PathSearch("selectPolicy", scale, nil),
		"stabilization_window_seconds": utils.PathSearch("stabilizationWindowSeconds", scale, nil),
	})

	return rst
}

func flattenSpecBehaviorScalePolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, len(policies))
	for i, v := range policies {
		rst[i] = map[string]interface{}{
			"period_seconds": utils.PathSearch("periodSeconds", v, nil),
			"type":           utils.PathSearch("type", v, nil),
			"value":          utils.PathSearch("value", v, nil),
		}
	}

	return rst
}

func flattenMetrics(metrics []interface{}) []interface{} {
	if len(metrics) == 0 {
		return nil
	}

	rst := make([]interface{}, len(metrics))
	for i, v := range metrics {
		rst[i] = map[string]interface{}{
			"container_resource": flattenSpecMetricsContainerResource(utils.PathSearch("containerResource", v, nil)),
			"external":           flattenSpecMetricsExternal(utils.PathSearch("external", v, nil)),
			"object":             flattenSpecMetricsObject(utils.PathSearch("object", v, nil)),
			"pods":               flattenSpecMetricsPods(utils.PathSearch("pods", v, nil)),
			"resource":           flattenSpecMetricsResource(utils.PathSearch("resource", v, nil)),
			"type":               utils.PathSearch("type", v, nil),
		}
	}

	return rst
}

func flattenSpecMetricsResource(resource interface{}) []interface{} {
	if resource == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"name":   utils.PathSearch("name", resource, nil),
		"target": flattenMetricTarget(utils.PathSearch("target", resource, nil)),
	})

	return rst
}

func flattenSpecMetricsPods(pods interface{}) []interface{} {
	if pods == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"metric": flattenMetricIdentifier(utils.PathSearch("metric", pods, nil)),
		"target": flattenMetricTarget(utils.PathSearch("target", pods, nil)),
	})

	return rst
}

func flattenSpecMetricsObject(object interface{}) []interface{} {
	if object == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"described_object": flattenSpecMetricsObjectDescribedObject(utils.PathSearch("describedObject", object, nil)),
		"metric":           flattenMetricIdentifier(utils.PathSearch("metric", object, nil)),
		"target":           flattenMetricTarget(utils.PathSearch("target", object, nil)),
	})

	return rst
}

func flattenSpecMetricsObjectDescribedObject(describedObject interface{}) []interface{} {
	if describedObject == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"api_version": flattenMetricIdentifier(utils.PathSearch("apiVersion", describedObject, nil)),
		"kind":        flattenMetricIdentifier(utils.PathSearch("kind", describedObject, nil)),
		"name":        flattenMetricTarget(utils.PathSearch("name", describedObject, nil)),
	})

	return rst
}

func flattenSpecMetricsExternal(external interface{}) []interface{} {
	if external == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"metric": flattenMetricIdentifier(utils.PathSearch("metric", external, nil)),
		"target": flattenMetricTarget(utils.PathSearch("target", external, nil)),
	})

	return rst
}

func flattenMetricIdentifier(metric interface{}) []interface{} {
	if metric == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"name":     utils.PathSearch("name", metric, nil),
		"selector": flattenMetricIdentifierSelector(utils.PathSearch("selector", metric, nil)),
	})

	return rst
}

func flattenMetricIdentifierSelector(selector interface{}) []interface{} {
	if selector == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"match_expressions": flattenMetricIdentifierSelectorMatchExpressions(utils.PathSearch("matchExpressions", selector, nil)),
		"match_labels":      utils.PathSearch("matchLabels", selector, nil),
	})

	return rst
}

func flattenMetricIdentifierSelectorMatchExpressions(matchExpressions interface{}) []interface{} {
	if matchExpressions == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"key":      flattenMetricIdentifierSelector(utils.PathSearch("key", matchExpressions, nil)),
		"operator": utils.PathSearch("operator", matchExpressions, nil),
		"values":   utils.PathSearch("values", matchExpressions, nil),
	})

	return rst
}

func flattenSpecMetricsContainerResource(containerResource interface{}) []interface{} {
	if containerResource == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"container": utils.PathSearch("container", containerResource, nil),
		"name":      (utils.PathSearch("name", containerResource, nil)),
		"target":    flattenMetricTarget(utils.PathSearch("target", containerResource, nil)),
	})

	return rst
}

func flattenMetricTarget(target interface{}) []interface{} {
	if target == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"average_utilization": utils.PathSearch("averageUtilization", target, nil),
		"average_value":       utils.PathSearch("averageValue", target, nil),
		"type":                utils.PathSearch("type", target, nil),
		"value":               utils.PathSearch("value", target, nil),
	})

	return rst
}

func flattenHorizontalPodAutoscalerStatus(status interface{}) []interface{} {
	if status == nil {
		return nil
	}

	conditions := utils.PathSearch("conditions", status, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"conditions":          flattenHorizontalPodAutoscalerStatusConditions(conditions),
		"current_metrics":     flattenMetrics(utils.PathSearch("currentMetrics", status, make([]interface{}, 0)).([]interface{})),
		"current_replicas":    utils.PathSearch("currentReplicas", status, nil),
		"desired_replicas":    utils.PathSearch("desiredReplicas", status, nil),
		"last_scale_time":     utils.PathSearch("lastScaleTime", status, nil),
		"observed_generation": utils.PathSearch("observedGeneration", status, nil),
	})

	return rst
}

func flattenHorizontalPodAutoscalerStatusConditions(conditions []interface{}) []interface{} {
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

func resourceV2HorizontalPodAutoscalerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	updateHorizontalPodAutoscalerHttpUrl := "apis/cci/v2/namespace/{namespace}/horizontalpodautoscalers/{name}"
	updateHorizontalPodAutoscalerPath := client.Endpoint + updateHorizontalPodAutoscalerHttpUrl
	updateHorizontalPodAutoscalerPath = strings.ReplaceAll(updateHorizontalPodAutoscalerPath, "{namespace}", ns)
	updateHorizontalPodAutoscalerPath = strings.ReplaceAll(updateHorizontalPodAutoscalerPath, "{name}", name)
	updateHorizontalPodAutoscalerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateHorizontalPodAutoscalerOpt.JSONBody = utils.RemoveNil(buildUpdateV2HorizontalPodAutoscalerParams(d))

	_, err = client.Request("PUT", updateHorizontalPodAutoscalerPath, &updateHorizontalPodAutoscalerOpt)
	if err != nil {
		return diag.Errorf("error updating CCI v2 HorizontalPodAutoscaler: %s", err)
	}
	return resourceV2HorizontalPodAutoscalerRead(ctx, d, meta)
}

func buildUpdateV2HorizontalPodAutoscalerParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       d.Get("kind"),
		"apiVersion": d.Get("api_version"),
		"metadata": map[string]interface{}{
			"name":              d.Get("name"),
			"namespace":         d.Get("namespace"),
			"uid":               d.Get("uid"),
			"resourceVersion":   d.Get("resource_version"),
			"creationTimestamp": d.Get("creation_timestamp"),
			"annotations":       d.Get("annotations"),
		},
		"spec":   map[string]interface{}{},
		"status": map[string]interface{}{},
	}

	return bodyParams
}

func resourceV2HorizontalPodAutoscalerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	deleteHorizontalPodAutoscalerHttpUrl := "apis/cci/v2/namespace/{namespace}/horizontalpodautoscalers/{name}"
	deleteHorizontalPodAutoscalerPath := client.Endpoint + deleteHorizontalPodAutoscalerHttpUrl
	deleteHorizontalPodAutoscalerPath = strings.ReplaceAll(deleteHorizontalPodAutoscalerPath, "{namespace}", ns)
	deleteHorizontalPodAutoscalerPath = strings.ReplaceAll(deleteHorizontalPodAutoscalerPath, "{name}", name)
	deleteHorizontalPodAutoscalerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteHorizontalPodAutoscalerPath, &deleteHorizontalPodAutoscalerOpt)
	if err != nil {
		return diag.Errorf("error deleting CCI v2 HorizontalPodAutoscaler: %s", err)
	}

	return nil
}

func resourceV2HorizontalPodAutoscalerImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	d.Set("name", d.Id())

	return []*schema.ResourceData{d}, nil
}

func GetHorizontalPodAutoscaler(client *golangsdk.ServiceClient, ns, name string) (interface{}, error) {
	getHorizontalPodAutoscalerHttpUrl := "apis/cci/v2/namespace/{namespace}/horizontalpodautoscalers/{name}"
	getHorizontalPodAutoscalerPath := client.Endpoint + getHorizontalPodAutoscalerHttpUrl
	getHorizontalPodAutoscalerPath = strings.ReplaceAll(getHorizontalPodAutoscalerPath, "{namespace}", ns)
	getHorizontalPodAutoscalerPath = strings.ReplaceAll(getHorizontalPodAutoscalerPath, "{name}", name)
	getHorizontalPodAutoscalerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getHorizontalPodAutoscalerResp, err := client.Request("GET", getHorizontalPodAutoscalerPath, &getHorizontalPodAutoscalerOpt)
	if err != nil {
		return getHorizontalPodAutoscalerResp, err
	}

	return utils.FlattenResponse(getHorizontalPodAutoscalerResp)
}
