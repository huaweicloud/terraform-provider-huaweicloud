package cci

import (
	"context"
	"fmt"
	"log"
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

var deploymentNonUpdatableParams = []string{"namespace", "name"}

// @API CCI POST /apis/cci/v2/namespaces/{namespace}/deployments
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/deployments/{name}
// @API CCI PUT /apis/cci/v2/namespaces/{namespace}/deployments/{name}
// @API CCI DELETE /apis/cci/v2/namespaces/{namespace}/deployments/{name}
func ResourceV2Deployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2DeploymentCreate,
		ReadContext:   resourceV2DeploymentRead,
		UpdateContext: resourceV2DeploymentUpdate,
		DeleteContext: resourceV2DeploymentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2DeploymentImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(deploymentNonUpdatableParams),

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
			"replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"min_ready_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"progress_deadline_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"selector": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     labelSelectorSchema(),
			},
			"template": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     deploymentTemplateSchema(),
			},
			"strategy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"rolling_update": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"delete_propagation_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func labelSelectorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"match_expressions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     matchExpressionsSchema(),
			},
		},
	}
	return &sc
}

func matchExpressionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

	return &sc
}

func deploymentTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     deploymentTemplateMetadataSchema(),
			},
			"spec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     deploymentTemplateSpecSchema(),
			},
		},
	}
	return &sc
}

func deploymentTemplateMetadataSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func deploymentTemplateSpecSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"containers": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     podContainersSchema(),
			},
			"dns_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"active_deadline_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"overhead": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"restart_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scheduler_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"set_hostname_as_pqdn": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"share_process_namespace": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"termination_grace_period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"affinity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_affinity": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     nodeAffinitySchema(),
						},
						"pod_anti_affinity": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     podAntiAffinitySchema(),
						},
					},
				},
			},
			"image_pull_secrets": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
	return &sc
}

func nodeAffinitySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"required_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     nodeSelectorSchema(),
			},
		},
	}
	return &sc
}

func nodeSelectorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_selector_terms": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     nodeSelectorTermSchema(),
			},
		},
	}

	return &sc
}

func nodeSelectorTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_expressions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     nodeSelectorRequirementSchema(),
			},
		},
	}

	return &sc
}

func nodeSelectorRequirementSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

	return &sc
}

func podAntiAffinitySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"preferred_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     weightedPodAffinityTermSchema(),
			},
			"required_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     podAffinityTermSchema(),
			},
		},
	}
	return &sc
}

func weightedPodAffinityTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"pod_affinity_term": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     podAffinityTermSchema(),
			},
			"weight": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
	return &sc
}

func podAffinityTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"label_selector": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     labelSelectorSchema(),
			},
			"namespaces": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"topology_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func resourceV2DeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createV2DeploymentHttpUrl := "apis/cci/v2/namespaces/{namespace}/deployments"
	createV2DeploymentPath := client.Endpoint + createV2DeploymentHttpUrl
	createV2DeploymentPath = strings.ReplaceAll(createV2DeploymentPath, "{namespace}", d.Get("namespace").(string))
	createV2DeploymentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createV2DeploymentOpt.JSONBody = utils.RemoveNil(buildCreateV2DeploymentParams(d))

	resp, err := client.Request("POST", createV2DeploymentPath, &createV2DeploymentOpt)
	if err != nil {
		return diag.Errorf("error creating CCI V2 Deployment: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	ns := utils.PathSearch("metadata.namespace", respBody, "").(string)
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if ns == "" || name == "" {
		return diag.Errorf("unable to find CCI V2 Deployment name or namespace from API response")
	}
	d.SetId(ns + "/" + name)

	err = waitForV2DeploymentStatusToCompleted(ctx, client, ns, name, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceV2DeploymentRead(ctx, d, meta)
}

func buildCreateV2DeploymentParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       "Deployment",
		"apiVersion": "cci/v2",
		"metadata": map[string]interface{}{
			"name":      d.Get("name"),
			"namespace": d.Get("namespace"),
		},
		"spec": map[string]interface{}{
			"replicas":                d.Get("replicas"),
			"minReadySeconds":         utils.ValueIgnoreEmpty(d.Get("min_ready_seconds")),
			"progressDeadlineSeconds": utils.ValueIgnoreEmpty(d.Get("progress_deadline_seconds")),
			"selector":                buildLabelSelectorParams(d.Get("selector.0")),
			"template":                buildV2DeploymentTemplateParams(d.Get("template.0")),
			"strategy":                buildV2DeploymentStrategyParams(d.Get("strategy.0")),
		},
	}

	return bodyParams
}

func buildV2DeploymentStrategyParams(strategy interface{}) map[string]interface{} {
	if strategy == nil {
		return nil
	}

	return map[string]interface{}{
		"type":           utils.PathSearch("type", strategy, nil),
		"rolling_update": utils.PathSearch("rolling_update", strategy, nil),
	}
}

func buildV2DeploymentTemplateParams(template interface{}) map[string]interface{} {
	if template == nil {
		return nil
	}

	metadata := utils.PathSearch("metadata|[0]", template, nil)
	spec := utils.PathSearch("spec|[0]", template, nil)
	return map[string]interface{}{
		"metadata": buildV2DeploymentTemplateMetadataParams(metadata),
		"spec":     buildV2DeploymentTemplateSpecParams(spec),
	}
}

func buildV2DeploymentTemplateSpecParams(spec interface{}) map[string]interface{} {
	if spec == nil {
		return nil
	}

	containers := utils.PathSearch("containers", spec, make([]interface{}, 0)).([]interface{})
	affinity := utils.PathSearch("affinity|[0]", spec, nil)
	imagePullSecrets := utils.PathSearch("image_pull_secrets", spec, &schema.Set{}).(*schema.Set).List()
	rst := map[string]interface{}{
		"containers":                    buildV2PodContainersParams(containers),
		"dnsPolicy":                     utils.ValueIgnoreEmpty(utils.PathSearch("dns_policy", spec, nil)),
		"activeDeadlineSeconds":         utils.ValueIgnoreEmpty(utils.PathSearch("active_deadline_seconds", spec, nil)),
		"hostname":                      utils.ValueIgnoreEmpty(utils.PathSearch("hostname", spec, nil)),
		"nodeName":                      utils.ValueIgnoreEmpty(utils.PathSearch("node_name", spec, nil)),
		"overhead":                      utils.ValueIgnoreEmpty(utils.PathSearch("overhead", spec, nil)),
		"restartPolicy":                 utils.ValueIgnoreEmpty(utils.PathSearch("restart_policy", spec, nil)),
		"schedulerName":                 utils.ValueIgnoreEmpty(utils.PathSearch("scheduler_name", spec, nil)),
		"setHostnameAsPQDN":             utils.ValueIgnoreEmpty(utils.PathSearch("set_hostname_as_pqdn", spec, nil)),
		"shareProcessNamespace":         utils.ValueIgnoreEmpty(utils.PathSearch("share_process_namespace", spec, nil)),
		"terminationGracePeriodSeconds": utils.ValueIgnoreEmpty(utils.PathSearch("termination_grace_period_seconds", spec, nil)),
		"affinity":                      buildV2DeploymentTemplateSpecAffinityParams(affinity),
		"imagePullSecrets":              buildImagePullSecretsParams(imagePullSecrets),
	}

	return rst
}

func buildImagePullSecretsParams(imagePullSecrets []interface{}) []interface{} {
	if len(imagePullSecrets) == 0 {
		return nil
	}

	rst := make([]interface{}, len(imagePullSecrets))
	for i, v := range imagePullSecrets {
		rst[i] = map[string]interface{}{
			"name": utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
		}
	}
	return rst
}

func buildV2DeploymentTemplateSpecAffinityParams(affinity interface{}) map[string]interface{} {
	if affinity == nil {
		return nil
	}
	return map[string]interface{}{
		"nodeAffinity":    buildNodeAffinityParams(utils.PathSearch("node_affinity|[0]", affinity, nil)),
		"podAntiAffinity": buildPodAntiAffinityParams(utils.PathSearch("pod_anti_affinity|[0]", affinity, nil)),
	}
}

func buildNodeAffinityParams(nodeAffinity interface{}) map[string]interface{} {
	if nodeAffinity == nil {
		return nil
	}
	nodeSelector := utils.PathSearch("required_during_scheduling_ignored_during_execution|[0]", nodeAffinity, nil)
	return map[string]interface{}{
		"requiredDuringSchedulingIgnoredDuringExecution": buildNodeSelectorParams(nodeSelector),
	}
}

func buildNodeSelectorParams(nodeSelector interface{}) map[string]interface{} {
	if nodeSelector == nil {
		return nil
	}
	nodeSelectorTerms := utils.PathSearch("node_selector_terms", nodeSelector, &schema.Set{}).(*schema.Set).List()
	return map[string]interface{}{
		"nodeSelectorTerms": buildNodeSelectorTermsParams(nodeSelectorTerms),
	}
}

func buildNodeSelectorTermsParams(nodeSelectorTerms []interface{}) []interface{} {
	if len(nodeSelectorTerms) == 0 {
		return nil
	}

	rst := make([]interface{}, len(nodeSelectorTerms))
	for i, v := range nodeSelectorTerms {
		matchExpressions := utils.PathSearch("match_expressions", v, &schema.Set{}).(*schema.Set).List()
		rst[i] = map[string]interface{}{
			"matchExpressions": buildV2DeploymentSelectorMatchExpressionsParams(matchExpressions),
		}
	}
	return rst
}

func buildPodAntiAffinityParams(podAntiAffinity interface{}) map[string]interface{} {
	if podAntiAffinity == nil {
		return nil
	}

	preferred := utils.PathSearch("preferred_during_scheduling_ignored_during_execution", podAntiAffinity, nil)
	required := utils.PathSearch("required_during_scheduling_ignored_during_execution", podAntiAffinity, nil)
	return map[string]interface{}{
		"preferredDuringSchedulingIgnoredDuringExecution": buildWeightedPodAffinityTermParams(preferred.(*schema.Set).List()),
		"requiredDuringSchedulingIgnoredDuringExecution":  buildPodAffinityTermsParams(required.(*schema.Set).List()),
	}
}

func buildPodAffinityTermsParams(podAffinityTerms []interface{}) []interface{} {
	if len(podAffinityTerms) == 0 {
		return nil
	}

	params := make([]interface{}, len(podAffinityTerms))
	for i, v := range podAffinityTerms {
		params[i] = buildPodAffinityTermParams(v)
	}

	return params
}

func buildWeightedPodAffinityTermParams(preferred []interface{}) []interface{} {
	if len(preferred) == 0 {
		return nil
	}

	rst := make([]interface{}, len(preferred))
	for i, v := range preferred {
		podAffinityTerm := utils.PathSearch("pod_affinity_term|[0]", v, &schema.Set{}).(*schema.Set).List()
		rst[i] = map[string]interface{}{
			"pod_affinity_term": buildPodAffinityTermParams(podAffinityTerm),
			"weight":            utils.PathSearch("weight", v, nil),
		}
	}
	return rst
}

func buildPodAffinityTermParams(podAffinityTerm interface{}) map[string]interface{} {
	if podAffinityTerm == nil {
		return nil
	}
	labelSelector := utils.PathSearch("label_selector|[0]", podAffinityTerm, nil)
	return map[string]interface{}{
		"labelSelector": buildLabelSelectorParams(labelSelector),
		"namespaces":    utils.PathSearch("namespaces", podAffinityTerm, &schema.Set{}).(*schema.Set).List(),
		"topologyKey":   utils.PathSearch("topology_key", podAffinityTerm, nil),
	}
}

func buildContainerEnvParams(env []interface{}) []interface{} {
	rst := make([]interface{}, len(env))
	for i, v := range env {
		rst[i] = map[string]interface{}{
			"name":  utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"value": utils.ValueIgnoreEmpty(utils.PathSearch("value", v, nil)),
		}
	}
	return rst
}

func buildContainerResourceParams(resources interface{}) map[string]interface{} {
	if resources == nil {
		return nil
	}
	return map[string]interface{}{
		"limits":   utils.ValueIgnoreEmpty(utils.PathSearch("limits", resources, nil)),
		"requests": utils.ValueIgnoreEmpty(utils.PathSearch("requests", resources, nil)),
	}
}

func buildV2DeploymentTemplateMetadataParams(metadata interface{}) map[string]interface{} {
	if metadata == nil {
		return nil
	}
	return map[string]interface{}{
		"labels":      utils.ValueIgnoreEmpty(utils.PathSearch("labels", metadata, nil)),
		"annotations": utils.ValueIgnoreEmpty(utils.PathSearch("annotations", metadata, nil)),
	}
}

func buildLabelSelectorParams(selector interface{}) interface{} {
	if selector == nil {
		return nil
	}

	params := utils.RemoveNil(map[string]interface{}{
		"matchLabels": utils.PathSearch("match_labels", selector, nil),
	})
	matchExpressions := utils.PathSearch("match_expressions", selector, &schema.Set{}).(*schema.Set).List()
	if len(matchExpressions) > 0 {
		params["matchExpressions"] = buildV2DeploymentSelectorMatchExpressionsParams(matchExpressions)
	}
	return params
}

func buildV2DeploymentSelectorMatchExpressionsParams(matchExpressions []interface{}) []interface{} {
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

func resourceV2DeploymentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	resp, err := GetV2Deployment(client, ns, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 Deployment")
	}

	mErr := multierror.Append(
		d.Set("namespace", utils.PathSearch("metadata.namespace", resp, nil)),
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("generation", utils.PathSearch("metadata.generation", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("replicas", int(utils.PathSearch("spec.replicas", resp, float64(0)).(float64))),
		d.Set("min_ready_seconds", int(utils.PathSearch("spec.minReadySeconds", resp, float64(0)).(float64))),
		d.Set("progress_deadline_seconds", int(utils.PathSearch("spec.progressDeadlineSeconds", resp, float64(0)).(float64))),
		d.Set("selector", flattenLabelSelector(utils.PathSearch("spec.selector", resp, nil))),
		d.Set("template", flattenSpecTemplate(utils.PathSearch("spec.template", resp, nil))),
		d.Set("delete_propagation_policy", d.Get("delete_propagation_policy")),
		d.Set("strategy", flattenSpecStrategy(utils.PathSearch("spec.strategy", resp, nil))),
		d.Set("progress_deadline_seconds", int(utils.PathSearch("spec.progressDeadlineSeconds", resp, float64(0)).(float64))),
		d.Set("status", flattenDeploymentStatus(utils.PathSearch("status", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLabelSelector(selector interface{}) []interface{} {
	if selector == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"match_labels": utils.PathSearch("matchLabels", selector, nil),
		"match_expressions": flattenMatchExpressions(
			utils.PathSearch("matchExpressions", selector, make([]interface{}, 0)).([]interface{})),
	})

	return rst
}

func flattenMatchExpressions(matchExpressions []interface{}) []interface{} {
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

func flattenSpecTemplate(template interface{}) []map[string]interface{} {
	if template == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"metadata": flattenTemplateMetadata(utils.PathSearch("metadata", template, nil)),
			"spec":     flattenTemplateSpec(utils.PathSearch("spec", template, nil)),
		},
	}

	return rst
}

func flattenTemplateSpec(spec interface{}) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	containers := utils.PathSearch("containers", spec, make([]interface{}, 0)).([]interface{})
	affinity := utils.PathSearch("affinity", spec, nil)
	imagePullSecrets := utils.PathSearch("imagePullSecrets", spec, make([]interface{}, 0)).([]interface{})

	rst := []map[string]interface{}{
		{
			"containers":                       flattenPodContainers(containers),
			"dns_policy":                       utils.PathSearch("dnsPolicy", spec, nil),
			"active_deadline_seconds":          utils.PathSearch("activeDeadlineSeconds", spec, nil),
			"hostname":                         utils.PathSearch("hostname", spec, nil),
			"node_name":                        utils.PathSearch("nodeName", spec, nil),
			"overhead":                         utils.PathSearch("overhead", spec, nil),
			"restart_policy":                   utils.PathSearch("restartPolicy", spec, nil),
			"scheduler_name":                   utils.PathSearch("schedulerName", spec, nil),
			"set_hostname_as_pqdn":             utils.PathSearch("setHostnameAsPQDN", spec, nil),
			"share_process_namespace":          utils.PathSearch("shareProcessNamespace", spec, nil),
			"termination_grace_period_seconds": utils.PathSearch("terminationGracePeriodSeconds", spec, nil),
			"affinity":                         flattenTemplateSpecAffinity(affinity),
			"image_pull_secrets":               flattenTemplateImagePullSecrets(imagePullSecrets),
		},
	}

	return rst
}

func flattenTemplateImagePullSecrets(imagePullSecrets []interface{}) []interface{} {
	if len(imagePullSecrets) == 0 {
		return nil
	}

	rst := make([]interface{}, len(imagePullSecrets))
	for i, v := range imagePullSecrets {
		rst[i] = map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
		}
	}
	return rst
}

func flattenTemplateSpecAffinity(affinity interface{}) []map[string]interface{} {
	if affinity == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"node_affinity":     flattenNodeAffinity(utils.PathSearch("nodeAffinity", affinity, nil)),
			"pod_anti_affinity": flattenPodAntiAffinity(utils.PathSearch("podAntiAffinity", affinity, nil)),
		},
	}
}

func flattenPodAntiAffinity(podAntiAffinity interface{}) map[string]interface{} {
	if podAntiAffinity == nil {
		return nil
	}

	preferred := utils.PathSearch("preferredDuringSchedulingIgnoredDuringExecution", podAntiAffinity, make([]interface{}, 0)).([]interface{})
	required := utils.PathSearch("requiredDuringSchedulingIgnoredDuringExecution", podAntiAffinity, make([]interface{}, 0)).([]interface{})
	return map[string]interface{}{
		"preferred_during_scheduling_ignored_during_execution": flattenWeightedPodAffinityTerms(preferred),
		"required_during_scheduling_ignored_during_execution":  flattenPodAffinityTerms(required),
	}
}

func flattenWeightedPodAffinityTerms(preferred []interface{}) []interface{} {
	if len(preferred) == 0 {
		return nil
	}

	rst := make([]interface{}, len(preferred))
	for i, v := range preferred {
		podAffinityTerm := utils.PathSearch("pod_affinity_term", v, nil)
		rst[i] = map[string]interface{}{
			"pod_affinity_term": flattenPodAffinityTerm(podAffinityTerm),
			"weight":            utils.PathSearch("weight", v, nil),
		}
	}
	return rst
}

func flattenPodAffinityTerms(required []interface{}) []interface{} {
	if len(required) == 0 {
		return nil
	}

	params := make([]interface{}, len(required))
	for i, v := range required {
		params[i] = flattenPodAffinityTerm(v)
	}

	return params
}

func flattenPodAffinityTerm(podAffinityTerm interface{}) map[string]interface{} {
	if podAffinityTerm == nil {
		return nil
	}
	labelSelector := utils.PathSearch("labelSelector", podAffinityTerm, nil)
	return map[string]interface{}{
		"label_selector": flattenLabelSelector(labelSelector),
		"namespaces":     utils.PathSearch("namespaces", podAffinityTerm, nil),
		"topology_key":   utils.PathSearch("topologyKey", podAffinityTerm, nil),
	}
}

func flattenNodeAffinity(nodeAffinity interface{}) []map[string]interface{} {
	if nodeAffinity == nil {
		return nil
	}
	nodeSelector := utils.PathSearch("requiredDuringSchedulingIgnoredDuringExecution", nodeAffinity, nil)
	return []map[string]interface{}{
		{
			"required_during_scheduling_ignored_during_execution": flattenNodeSelector(nodeSelector),
		},
	}
}

func flattenNodeSelector(nodeSelector interface{}) []map[string]interface{} {
	if nodeSelector == nil {
		return nil
	}
	nodeSelectorTerms := utils.PathSearch("nodeSelectorTerms", nodeSelector, make([]interface{}, 0)).([]interface{})
	return []map[string]interface{}{
		{
			"node_selector_terms": flattenNodeSelectorTerms(nodeSelectorTerms),
		},
	}
}

func flattenNodeSelectorTerms(nodeSelectorTerms []interface{}) []interface{} {
	if len(nodeSelectorTerms) == 0 {
		return nil
	}

	rst := make([]interface{}, len(nodeSelectorTerms))
	for i, v := range nodeSelectorTerms {
		matchExpressions := utils.PathSearch("matchExpressions", v, make([]interface{}, 0)).([]interface{})
		rst[i] = map[string]interface{}{
			"match_expressions": flattenMatchExpressions(matchExpressions),
		}
	}
	return rst
}

func flattenContainerEnv(env []interface{}) []interface{} {
	if len(env) == 0 {
		return nil
	}

	rst := make([]interface{}, len(env))
	for i, v := range env {
		rst[i] = map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
	}
	return rst
}

func flattenContainerResource(resources interface{}) []map[string]interface{} {
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

func flattenTemplateMetadata(metadata interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"labels":      utils.PathSearch("labels", metadata, nil),
			"annotations": utils.PathSearch("annotations", metadata, nil),
		},
	}

	return rst
}

func flattenSpecStrategy(strategy interface{}) []map[string]interface{} {
	if strategy == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"type":           utils.PathSearch("type", strategy, nil),
			"rolling_update": utils.PathSearch("rollingUpdate", strategy, nil),
		},
	}

	return rst
}

func flattenDeploymentStatus(status interface{}) []map[string]interface{} {
	if status == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"observed_generation": utils.PathSearch("observedGeneration", status, nil),
			"conditions":          flattenDeploymentStatusConditions(utils.PathSearch("conditions", status, make([]interface{}, 0)).([]interface{})),
		},
	}

	return rst
}

func flattenDeploymentStatusConditions(conditions []interface{}) []interface{} {
	if len(conditions) == 0 {
		return nil
	}

	rst := make([]interface{}, len(conditions))
	for i, v := range conditions {
		rst[i] = map[string]interface{}{
			"type":                 utils.PathSearch("type", v, nil),
			"status":               utils.PathSearch("status", v, nil),
			"last_update_time":     utils.PathSearch("lastUpdateTime", v, nil),
			"last_transition_time": utils.PathSearch("lastTransitionTime", v, nil),
			"reason":               utils.PathSearch("reason", v, nil),
			"message":              utils.PathSearch("message", v, nil),
		}
	}

	return rst
}

func resourceV2DeploymentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	updateV2DeploymentHttpUrl := "apis/cci/v2/namespaces/{namespace}/deployments/{name}"
	updateV2DeploymentPath := client.Endpoint + updateV2DeploymentHttpUrl
	updateV2DeploymentPath = strings.ReplaceAll(updateV2DeploymentPath, "{namespace}", ns)
	updateV2DeploymentPath = strings.ReplaceAll(updateV2DeploymentPath, "{name}", name)
	updateV2DeploymentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateV2DeploymentOpt.JSONBody = utils.RemoveNil(buildCreateV2DeploymentParams(d))

	_, err = client.Request("PUT", updateV2DeploymentPath, &updateV2DeploymentOpt)
	if err != nil {
		return diag.Errorf("error updating CCI v2 Deployment: %s", err)
	}

	err = waitForV2DeploymentStatusToCompleted(ctx, client, ns, name, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceV2DeploymentRead(ctx, d, meta)
}

func resourceV2DeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	deleteV2DeploymentHttpUrl := "apis/cci/v2/namespaces/{namespace}/deployments/{name}"
	deleteV2DeploymentPath := client.Endpoint + deleteV2DeploymentHttpUrl
	deleteV2DeploymentPath = strings.ReplaceAll(deleteV2DeploymentPath, "{namespace}", ns)
	deleteV2DeploymentPath = strings.ReplaceAll(deleteV2DeploymentPath, "{name}", name)
	deleteV2DeploymentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	if v, ok := d.GetOk("delete_propagation_policy"); ok {
		deleteV2DeploymentOpt.JSONBody = map[string]interface{}{
			"kind":              "DeleteOptions",
			"apiVersion":        "v1",
			"propagationPolicy": v.(string),
		}
	}
	_, err = client.Request("DELETE", deleteV2DeploymentPath, &deleteV2DeploymentOpt)
	if err != nil {
		return diag.Errorf("error deleting CCI v2 Deployment: %s", err)
	}

	err = waitForDeleteV2DeploymentStatus(ctx, client, ns, name, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForV2DeploymentStatusToCompleted(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Completed"},
		Refresh:      refreshCreateV2DeploymentStatus(client, ns, name),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of the CCI Deployment to complete: %s", err)
	}
	return nil
}

func refreshCreateV2DeploymentStatus(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetV2Deployment(client, ns, name)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("status.conditions[?type=='Available'].status|[0]", resp, "").(string)
		if status == "True" {
			return resp, "Completed", nil
		}

		return resp, "Pending", nil
	}
}

func waitForDeleteV2DeploymentStatus(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Deleted"},
		Refresh:      refreshDeleteV2DeploymentStatus(client, ns, name),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of the CCI v2 Deployment to complete: %s", err)
	}
	return nil
}

func refreshDeleteV2DeploymentStatus(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetV2Deployment(client, ns, name)
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[DEBUG] successfully deleted CCI Deployment: %s", name)
			return "", "Deleted", nil
		}
		return resp, "Pending", nil
	}
}

func resourceV2DeploymentImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<namespace>/<name>', but '%s'", importedId)
	}

	d.Set("namespace", parts[0])
	d.Set("name", parts[1])

	return []*schema.ResourceData{d}, nil
}

func GetV2Deployment(client *golangsdk.ServiceClient, namespace, name string) (interface{}, error) {
	getV2DeploymentHttpUrl := "apis/cci/v2/namespaces/{namespace}/deployments/{name}"
	getV2DeploymentPath := client.Endpoint + getV2DeploymentHttpUrl
	getV2DeploymentPath = strings.ReplaceAll(getV2DeploymentPath, "{namespace}", namespace)
	getV2DeploymentPath = strings.ReplaceAll(getV2DeploymentPath, "{name}", name)
	getV2DeploymentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getV2DeploymentResp, err := client.Request("GET", getV2DeploymentPath, &getV2DeploymentOpt)
	if err != nil {
		return getV2DeploymentResp, err
	}

	return utils.FlattenResponse(getV2DeploymentResp)
}
