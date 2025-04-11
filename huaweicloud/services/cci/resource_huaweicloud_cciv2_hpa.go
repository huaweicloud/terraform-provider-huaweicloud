package cci

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var hpaNonUpdatableParams = []string{"namespace", "name"}

// @API CCI POST /apis/cci/v2/namespaces/{namespace}/horizontalpodautoscalers
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/horizontalpodautoscalers/{name}
// @API CCI PUT /apis/cci/v2/namespaces/{namespace}/horizontalpodautoscalers/{name}
// @API CCI DELETE /apis/cci/v2/namespaces/{namespace}/horizontalpodautoscalers/{name}
func ResourceV2HPA() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2HPACreate,
		ReadContext:   resourceV2HPARead,
		UpdateContext: resourceV2HPAUpdate,
		DeleteContext: resourceV2HPADelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2HPAImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(hpaNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "cci/v2",
			},
			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "HorizontalPodAutoscaler",
			},
			"behavior": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scale_down": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     hpaBehaviorSchema(),
						},
						"scale_up": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     hpaBehaviorSchema(),
						},
					},
				},
			},
			"max_replicas": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"metrics": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     hpaMetricsSchema(),
			},
			"min_replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"scale_target_ref": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
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
							Elem:     hpaMetricsSchema(),
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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func hpaBehaviorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period_seconds": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"select_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stabilization_window_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}

	return &sc
}

func hpaMetricsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"container_resource": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     containerResourceMetricSourceSchema(),
			},
			"external": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     externalMetricSourceSchema(),
			},
			"object": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     objectmetricSourceSchema(),
			},
			"pods": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podMetricSourceSchema(),
			},
			"resources": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     resourceMetricSourceSchema(),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}

	return &sc
}

func containerResourceMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"container": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"target": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     metricTargetSchema(),
			},
		},
	}
	return &sc
}

func metricTargetSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"average_utilization": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"average_value": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func externalMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     metricIdentifierSchema(),
			},
			"target": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     metricTargetSchema(),
			},
		},
	}
	return &sc
}

func metricIdentifierSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeMap,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"match_labels": {
							Type:     schema.TypeMap,
							Optional: true,
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

func objectmetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"described_object": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"metric": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     metricIdentifierSchema(),
			},
			"target": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     metricTargetSchema(),
			},
		},
	}
	return &sc
}

func podMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     metricIdentifierSchema(),
			},
			"target": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     metricTargetSchema(),
			},
		},
	}
	return &sc
}

func resourceMetricSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"target": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     metricTargetSchema(),
			},
		},
	}
	return &sc
}

func resourceV2HPACreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createHPAHttpUrl := "apis/cci/v2/namespaces/{namespace}/horizontalpodautoscalers"
	createHPAPath := client.Endpoint + createHPAHttpUrl
	createHPAPath = strings.ReplaceAll(createHPAPath, "{namespace}", d.Get("namespace").(string))
	createHPAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createHPAOpt.JSONBody = utils.RemoveNil(buildCreateV2HPAParams(d))

	resp, err := client.Request("POST", createHPAPath, &createHPAOpt)
	if err != nil {
		return diag.Errorf("error creating CCI HPA: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ns := utils.PathSearch("metadata.namespace", respBody, "").(string)
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if ns == "" || name == "" {
		return diag.Errorf("unable to find CCI V2 HPA name or namespace from API response")
	}
	d.SetId(ns + "/" + name)

	err = waitForV2CreateOrUpdateStatus(ctx, client, ns, name, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceV2HPARead(ctx, d, meta)
}

func buildCreateV2HPAParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"apiVersion": d.Get("api_version"),
		"kind":       d.Get("kind"),
		"metadata": map[string]interface{}{
			"name":      d.Get("name"),
			"namespace": d.Get("namespace"),
		},
		"spec": map[string]interface{}{
			"maxReplicas":    d.Get("max_replicas"),
			"minReplicas":    d.Get("min_replicas"),
			"scaleTargetRef": utils.ValueIgnoreEmpty(buildV2CreateHPAScaleTargetRefParams(d.Get("scale_target_ref.0"))),
			"metrics":        utils.ValueIgnoreEmpty(buildV2CreateHPAMetricsParams(d.Get("metrics").(*schema.Set).List())),
			"behavior":       utils.ValueIgnoreEmpty(buildV2CreateHPABehaviorParams(d.Get("behavior.0"))),
		},
	}

	return bodyParams
}

func buildV2CreateHPAScaleTargetRefParams(scaleTargetRef interface{}) interface{} {
	if scaleTargetRef == nil {
		return nil
	}

	params := utils.RemoveNil(map[string]interface{}{
		"apiVersion": utils.PathSearch("api_version", scaleTargetRef, nil),
		"kind":       utils.PathSearch("kind", scaleTargetRef, nil),
		"name":       utils.PathSearch("name", scaleTargetRef, nil),
	})
	return params
}

func buildV2CreateHPAMetricsParams(metrics []interface{}) []interface{} {
	rst := make([]interface{}, len(metrics))
	for i, v := range metrics {
		rst[i] = utils.ValueIgnoreEmpty(map[string]interface{}{
			"containerResource": utils.ValueIgnoreEmpty(buildV2CreateHPAMetricsConResParams(utils.PathSearch("container_resource|[0]", v, nil))),
			"external":          utils.ValueIgnoreEmpty(buildV2CreateHPAMetricsExternalParams(utils.PathSearch("external|[0]", v, nil))),
			"object":            utils.ValueIgnoreEmpty(buildV2CreateHPAMetricsObjectParams(utils.PathSearch("object|[0]", v, nil))),
			"pods":              utils.ValueIgnoreEmpty(buildV2CreateHPAMetricsPodsParams(utils.PathSearch("pods|[0]", v, nil))),
			"resource":          utils.ValueIgnoreEmpty(buildV2CreateHPAMetricsResourceParams(utils.PathSearch("resources|[0]", v, nil))),
			"type":              utils.PathSearch("type", v, nil),
		})
	}
	return rst
}

func buildV2CreateHPAMetricsConResParams(cr interface{}) map[string]interface{} {
	if cr == nil {
		return nil
	}

	params := map[string]interface{}{
		"container": utils.PathSearch("container", cr, nil),
		"name":      utils.PathSearch("name", cr, nil),
		"target":    utils.RemoveNil(buildV2CreateHPACommonMetricTargetParams(utils.PathSearch("target|[0]", cr, nil))),
	}
	return params
}

func buildV2CreateHPACommonMetricTargetParams(pods interface{}) map[string]interface{} {
	if pods == nil {
		return nil
	}

	params := map[string]interface{}{
		"averageUtilization": utils.ValueIgnoreEmpty(utils.PathSearch("average_utilization", pods, nil)),
		"averageValue":       utils.ValueIgnoreEmpty(utils.PathSearch("average_value", pods, nil)),
		"type":               utils.ValueIgnoreEmpty(utils.PathSearch("type", pods, nil)),
		"value":              utils.ValueIgnoreEmpty(utils.PathSearch("value", pods, nil)),
	}
	return params
}

func buildV2CreateHPAMetricsExternalParams(external interface{}) map[string]interface{} {
	if external == nil {
		return nil
	}

	params := map[string]interface{}{
		"metric": utils.RemoveNil(
			buildV2CreateHPAMetricsExternalMetricIdentifierParams(utils.PathSearch("metric|[0]", external, nil))),
		"target": utils.RemoveNil(buildV2CreateHPACommonMetricTargetParams(utils.PathSearch("target|[0]", external, nil))),
	}
	return params
}

func buildV2CreateHPAMetricsExternalMetricIdentifierParams(metric interface{}) map[string]interface{} {
	if metric == nil {
		return nil
	}

	params := map[string]interface{}{
		"name": utils.PathSearch("name", metric, nil),
		"selector": utils.RemoveNil(
			buildV2CreateHPAMetricsExternalMetricIdentifierSelectorParams(utils.PathSearch("selector|[0]", metric, nil))),
	}
	return params
}

func buildV2CreateHPAMetricsExternalMetricIdentifierSelectorParams(selector interface{}) map[string]interface{} {
	if selector == nil {
		return nil
	}

	params := map[string]interface{}{
		"matchLabels": utils.ValueIgnoreEmpty(utils.PathSearch("match_labels", selector, nil)),
	}

	matchExpressions := buildV2CreateHPAMetricIdentifierSelectorMatchExpressionsParams(
		utils.PathSearch("match_expressions", selector, &schema.Set{}).(*schema.Set).List())
	if len(matchExpressions) > 0 {
		params["matchExpressions"] = matchExpressions
	}
	return params
}

func buildV2CreateHPAMetricIdentifierSelectorMatchExpressionsParams(matchExpressions []interface{}) []interface{} {
	rst := make([]interface{}, len(matchExpressions))
	for i, v := range matchExpressions {
		rst[i] = map[string]interface{}{
			"key":      utils.ValueIgnoreEmpty(utils.PathSearch("key", v, nil)),
			"operator": utils.ValueIgnoreEmpty(utils.PathSearch("operator", v, nil)),
			"values":   utils.ValueIgnoreEmpty(utils.PathSearch("values", v, &schema.Set{}).(*schema.Set).List()),
		}
	}
	return rst
}

func buildV2CreateHPAMetricsObjectParams(object interface{}) map[string]interface{} {
	if object == nil {
		return nil
	}

	params := map[string]interface{}{
		"describedObject": buildV2CreateHPAMetricsObjectDescribedParams(utils.PathSearch("described_object|[0]", object, nil)),
		"metric": utils.RemoveNil(
			buildV2CreateHPAMetricsExternalMetricIdentifierParams(utils.PathSearch("metric|[0]", object, nil))),
		"target": utils.RemoveNil(buildV2CreateHPACommonMetricTargetParams(utils.PathSearch("target|[0]", object, nil))),
	}
	return params
}

func buildV2CreateHPAMetricsObjectDescribedParams(object interface{}) map[string]interface{} {
	if object == nil {
		return nil
	}

	params := map[string]interface{}{
		"apiVersion": utils.ValueIgnoreEmpty(utils.PathSearch("api_version", object, nil)),
		"kind":       utils.ValueIgnoreEmpty(utils.PathSearch("kind", object, nil)),
		"name":       utils.ValueIgnoreEmpty(utils.PathSearch("name", object, nil)),
	}
	return params
}

func buildV2CreateHPAMetricsPodsParams(pods interface{}) interface{} {
	if pods == nil {
		return nil
	}

	params := utils.RemoveNil(map[string]interface{}{
		"metric": utils.RemoveNil(
			buildV2CreateHPAMetricsExternalMetricIdentifierParams(utils.PathSearch("metric|[0]", pods, nil))),
		"target": utils.RemoveNil(buildV2CreateHPACommonMetricTargetParams(utils.PathSearch("target|[0]", pods, nil))),
	})
	return params
}

func buildV2CreateHPAMetricsResourceParams(resources interface{}) map[string]interface{} {
	if resources == nil {
		return nil
	}

	params := map[string]interface{}{
		"name":   utils.ValueIgnoreEmpty(utils.PathSearch("name", resources, nil)),
		"target": utils.RemoveNil(buildV2CreateHPACommonMetricTargetParams(utils.PathSearch("target|[0]", resources, nil))),
	}
	return params
}

func buildV2CreateHPABehaviorParams(behavior interface{}) map[string]interface{} {
	if behavior == nil {
		return nil
	}

	params := map[string]interface{}{
		"scaleDown": buildV2CreateHPABehaviorScaleParams(utils.PathSearch("scale_down|[0]", behavior, nil)),
		"scaleUp":   buildV2CreateHPABehaviorScaleParams(utils.PathSearch("scale_up|[0]", behavior, nil)),
	}
	return params
}

func buildV2CreateHPABehaviorScaleParams(scale interface{}) map[string]interface{} {
	if scale == nil {
		return nil
	}

	policies := utils.PathSearch("policies", scale, &schema.Set{}).(*schema.Set).List()
	params := map[string]interface{}{
		"selectPolicy":               utils.ValueIgnoreEmpty(utils.PathSearch("select_policy", scale, nil)),
		"stabilizationWindowSeconds": utils.ValueIgnoreEmpty(utils.PathSearch("stabilization_window_seconds", scale, nil)),
	}
	if len(policies) > 0 {
		params["policies"] = buildV2CreateHPABehaviorScalePoliciesParams(policies)
	}
	return params
}

func buildV2CreateHPABehaviorScalePoliciesParams(policies []interface{}) []interface{} {
	rst := make([]interface{}, len(policies))
	for i, v := range policies {
		rst[i] = map[string]interface{}{
			"periodSeconds": utils.ValueIgnoreEmpty(utils.PathSearch("period_seconds", v, nil)),
			"type":          utils.ValueIgnoreEmpty(utils.PathSearch("type", v, nil)),
			"value":         utils.ValueIgnoreEmpty(utils.PathSearch("value", v, nil)),
		}
	}
	return rst
}

func resourceV2HPARead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	resp, err := GetV2HPA(client, ns, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 HPA")
	}

	mErr := multierror.Append(
		d.Set("namespace", utils.PathSearch("metadata.namespace", resp, nil)),
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("behavior", flattenSpecBehavior(utils.PathSearch("spec.behavior", resp, nil))),
		d.Set("max_replicas", utils.PathSearch("spec.maxReplicas", resp, nil)),
		d.Set("metrics", flattenMetrics(utils.PathSearch("spec.metrics", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("min_replicas", utils.PathSearch("spec.minReplicas", resp, nil)),
		d.Set("scale_target_ref", flattenSpecScaleTargetRef(utils.PathSearch("spec.scaleTargetRef", resp, nil))),
		d.Set("status", flattenHPAStatus(utils.PathSearch("status", resp, nil))),
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
			"resources":          flattenSpecMetricsResource(utils.PathSearch("resource", v, nil)),
			"type":               utils.PathSearch("type", v, nil),
		}
	}

	return rst
}

func flattenSpecMetricsResource(res interface{}) []interface{} {
	if res == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"name":   utils.PathSearch("name", res, nil),
		"target": flattenMetricTarget(utils.PathSearch("target", res, nil)),
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

func flattenHPAStatus(status interface{}) []interface{} {
	if status == nil {
		return nil
	}

	conditions := utils.PathSearch("conditions", status, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"conditions":          flattenHPAStatusConditions(conditions),
		"current_metrics":     flattenMetrics(utils.PathSearch("currentMetrics", status, make([]interface{}, 0)).([]interface{})),
		"current_replicas":    utils.PathSearch("currentReplicas", status, nil),
		"desired_replicas":    utils.PathSearch("desiredReplicas", status, nil),
		"last_scale_time":     utils.PathSearch("lastScaleTime", status, nil),
		"observed_generation": utils.PathSearch("observedGeneration", status, nil),
	})

	return rst
}

func flattenHPAStatusConditions(conditions []interface{}) []interface{} {
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

func resourceV2HPAUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	updateHPAHttpUrl := "apis/cci/v2/namespaces/{namespace}/horizontalpodautoscalers/{name}"
	updateHPAPath := client.Endpoint + updateHPAHttpUrl
	updateHPAPath = strings.ReplaceAll(updateHPAPath, "{namespace}", ns)
	updateHPAPath = strings.ReplaceAll(updateHPAPath, "{name}", name)
	updateHPAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateHPAOpt.JSONBody = utils.RemoveNil(buildCreateV2HPAParams(d))

	_, err = client.Request("PUT", updateHPAPath, &updateHPAOpt)
	if err != nil {
		return diag.Errorf("error updating CCI v2 HPA: %s", err)
	}

	err = waitForV2CreateOrUpdateStatus(ctx, client, ns, name, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceV2HPARead(ctx, d, meta)
}

func resourceV2HPADelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	deleteHPAHttpUrl := "apis/cci/v2/namespaces/{namespace}/horizontalpodautoscalers/{name}"
	deleteHPAPath := client.Endpoint + deleteHPAHttpUrl
	deleteHPAPath = strings.ReplaceAll(deleteHPAPath, "{namespace}", ns)
	deleteHPAPath = strings.ReplaceAll(deleteHPAPath, "{name}", name)
	deleteHPAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteHPAPath, &deleteHPAOpt)
	if err != nil {
		return diag.Errorf("error deleting CCI v2 HPA: %s", err)
	}

	err = waitForDeleteV2HPAStatus(ctx, client, ns, name, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForV2CreateOrUpdateStatus(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Completed"},
		Refresh:      refreshCreateOrUpdateStatus(client, ns, name),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of the CCI HPA to complete: %s", err)
	}
	return nil
}

func refreshCreateOrUpdateStatus(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetV2HPA(client, ns, name)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("status.conditions[?type=='ScalingActive'].status|[0]", resp, "").(string)
		if status == "True" {
			return resp, "Completed", nil
		}
		if status == "False" {
			message := utils.PathSearch("status.conditions[?type=='ScalingActive'].message|[0]", resp, "").(string)
			return resp, "ERROR", fmt.Errorf("error waiting for the status: %s", message)
		}

		return resp, "Pending", nil
	}
}

func resourceV2HPAImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<namespace>/<name>', but '%s'", importedId)
	}

	d.Set("namespace", parts[0])
	d.Set("name", parts[1])

	return []*schema.ResourceData{d}, nil
}

func waitForDeleteV2HPAStatus(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Deleted"},
		Refresh:      refreshDeleteV2HPAStatus(client, ns, name),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of the CCI v2 HPA to complete: %s", err)
	}
	return nil
}

func refreshDeleteV2HPAStatus(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetV2HPA(client, ns, name)
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return "", "Deleted", nil
		}
		return resp, "Pending", nil
	}
}

func GetV2HPA(client *golangsdk.ServiceClient, ns, name string) (interface{}, error) {
	getHPAHttpUrl := "apis/cci/v2/namespaces/{namespace}/horizontalpodautoscalers/{name}"
	getHPAPath := client.Endpoint + getHPAHttpUrl
	getHPAPath = strings.ReplaceAll(getHPAPath, "{namespace}", ns)
	getHPAPath = strings.ReplaceAll(getHPAPath, "{name}", name)
	getHPAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getHPAResp, err := client.Request("GET", getHPAPath, &getHPAOpt)
	if err != nil {
		return getHPAResp, err
	}

	return utils.FlattenResponse(getHPAResp)
}
