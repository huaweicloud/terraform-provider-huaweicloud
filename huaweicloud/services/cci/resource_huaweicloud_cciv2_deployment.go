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

var deploymentNonUpdatableParams = []string{"namespace", "name", "ip_families", "subnets", "subnets.*.subnet_id"}

// @API CCI POST /apis/cci/v2/namespaces/{namespace}/deployment
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/deployment/{name}
// @API CCI PUT /apis/cci/v2/namespaces/{namespace}/deployment/{name}
// @API CCI DELETE /apis/cci/v2/namespaces/{namespace}/deployment/{name}
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
			Create: schema.DefaultTimeout(10 * time.Minute),
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
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace of the CCI.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the CCI deployment.`,
			},
			"replicas": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the replicas of the CCI deployment.`,
			},
			"selector": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_labels": {
							Type:        schema.TypeMap,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the match labels of the CCI deployment selector.`,
						},
					},
				},
				Description: `Specifies the selector of the CCI deployment.`,
			},
			"template": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        deploymentTemplateSchema(),
				Description: `Specifies the template of the CCI deployment.`,
			},
			"strategy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the type of the CCI deployment strategy.`,
						},
						"rolling_update": {
							Type:        schema.TypeMap,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the rolling update config of the CCI deployment strategy.`,
						},
					},
				},
				Description: `Specifies the strategy of the CCI deployment.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the CCI deployment.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the CCI deployment.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the metadata annotations of the CCI deployment.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the CCI deployment.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the CCI deployment.`,
			},
			"generation": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The generation of the CCI deployment.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the CCI deployment.`,
			},
			"progress_deadline_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The progress deadline seconds of the CCI deployment.`,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"observed_generation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The observed generation of the CCI deployment.`,
						},
						"conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the CCI deployment conditions.`,
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Tthe status of the CCI deployment conditions.`,
									},
									"last_update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The last update time of the CCI deployment conditions.`,
									},
									"last_transition_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The last transition time of the CCI deployment conditions.`,
									},
									"reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The reason of the CCI deployment conditions.`,
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The message of the CCI deployment conditions.`,
									},
								},
							},
							Description: `Tthe conditions of the CCI deployment.`,
						},
					},
				},
				Description: `The status of the CCI deployment.`,
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

func deploymentTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the metadata labels of the CCI deployment template.`,
			},
			"metadata_annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the metadata annotations of the CCI deployment template.`,
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					oldVal, newVal := d.GetChange("annotations")
					for key, value := range newVal.(map[string]interface{}) {
						if mapValue, exists := oldVal.(map[string]interface{})[key]; exists && mapValue == value {
							continue
						}
						return false
					}
					return true
				},
			},
			"spec_containers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the subnet ID of the CCI network.`,
						},
						"image": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the subnet ID of the CCI network.`,
						},
						"env": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the subnet ID of the CCI network.`,
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the subnet ID of the CCI network.`,
									},
								},
							},
							Description: `Specifies the container environment of the CCI deployment.`,
						},
						"resources": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limits": {
										Type:        schema.TypeMap,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `Specifies the subnet ID of the CCI network.`,
									},
									"requests": {
										Type:        schema.TypeMap,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `Specifies the subnet ID of the CCI network.`,
									},
								},
							},
							Description: `Specifies the container of the CCI deployment.`,
						},
					},
				},
				Description: `Specifies the IP families of the CCI deployment.`,
			},
			"spec_dns_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the security group IDs of the CCI deployment.`,
			},
			"spec_image_pull_secrets": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateJsonString,
				Description:  `Specifies the subnets of the CCI deployment.`,
			},
			"spec_affinity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateJsonString,
				Description:  `Specifies the subnets of the CCI deployment.`,
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

	createNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/networks"
	createNetworkPath := client.Endpoint + createNetworkHttpUrl
	createNetworkPath = strings.ReplaceAll(createNetworkPath, "{namespace}", d.Get("namespace").(string))
	createNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createNetworkOpt.JSONBody = utils.RemoveNil(buildCreateV2DeploymentParams(d))

	resp, err := client.Request("POST", createNetworkPath, &createNetworkOpt)
	if err != nil {
		return diag.Errorf("error creating CCI Network: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	ns := utils.PathSearch("metadata.namespace", respBody, "").(string)
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if ns == "" || name == "" {
		return diag.Errorf("unable to find CCI Network name or namespace from API response")
	}
	d.SetId(ns + "/" + name)

	err = waitForCreateV2DeploymentStatus(ctx, client, ns, name, d.Timeout(schema.TimeoutCreate))
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
			"replicas": d.Get("replicas"),
			"selector": buildV2DeploymentSelectorParams(d.Get("selector")),
			"template": buildV2DeploymentTemplateParams(d),
			// "strategy":     buildV2DeploymentStrategyParams(d.Get("strategy")),
		},
	}

	return bodyParams
}

func buildV2DeploymentTemplateParams(d *schema.ResourceData) map[string]interface{} {
	template := map[string]interface{}{
		"metadata": buildV2DeploymentTemplateMetadataParams(d),
		"spec":     buildV2DeploymentTemplateSpecParams(d),
	}

	return template
}

func buildV2DeploymentTemplateSpecParams(d *schema.ResourceData) map[string]interface{} {
	metadata := map[string]interface{}{
		"containers":       buildV2DeploymentTemplateSpecContainersParams(d),
		"dnsPolicy":        d.Get("template.0.spec_dns_policy"),
		"imagePullSecrets": d.Get("template.0.spec_image_pull_secrets"),
		"affinity":         d.Get("template.0.affinity"),
	}

	return metadata
}

func buildV2DeploymentTemplateSpecContainersParams(d *schema.ResourceData) []interface{} {
	containers := d.Get("template.0.spec_containers").([]interface{})
	if len(containers) == 0 {
		return nil
	}
	containersParams := make([]interface{}, 0, len(containers))
	for i, v := range containers {
		containersParams[i] = map[string]interface{}{
			"name":      utils.PathSearch("name", v, nil),
			"image":     utils.PathSearch("image", v, nil),
			"env":       utils.PathSearch("env", v, nil),
			"resources": utils.PathSearch("resources", v, nil),
		}
	}

	return containersParams
}

func buildV2DeploymentTemplateMetadataParams(d *schema.ResourceData) map[string]interface{} {
	metadata := map[string]interface{}{
		"labels":      d.Get("template.0.metadata_labels"),
		"annotations": d.Get("template.0.metadata_annotations"),
	}

	return metadata
}

func buildV2DeploymentSelectorParams(selector interface{}) interface{} {
	if selector == nil {
		return nil
	}

	params := map[string]interface{}{
		"matchLabels": utils.PathSearch("[0].match_labels", selector, nil),
	}
	return params
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
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 network")
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
		d.Set("selector", flattenSpecSelector(utils.PathSearch("spec.selector", resp, nil))),
		d.Set("template", flattenSpecTemplate(utils.PathSearch("spec.template", resp, nil))),
		d.Set("strategy", flattenSpecStrategy(utils.PathSearch("spec.strategy", resp, nil))),
		d.Set("progress_deadline_seconds", int(utils.PathSearch("spec.progressDeadlineSeconds", resp, float64(0)).(float64))),
		d.Set("status", flattenDeploymentStatus(utils.PathSearch("status", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSpecSelector(selector interface{}) []map[string]interface{} {
	if selector == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"match_labels": utils.PathSearch("matchLabels", selector, nil),
		},
	}

	return rst
}

func flattenSpecTemplate(template interface{}) []map[string]interface{} {
	if template == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"metadata": utils.PathSearch("metadata", template, nil),
			"spec":     utils.PathSearch("spec", template, nil),
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

	updateNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/networks/{name}"
	updateNetworkPath := client.Endpoint + updateNetworkHttpUrl
	updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{namespace}", d.Get("namespace").(string))
	updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{name}", d.Get("name").(string))
	updateNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateNetworkOpt.JSONBody = utils.RemoveNil(buildCreateV2DeploymentParams(d))

	_, err = client.Request("PUT", updateNetworkPath, &updateNetworkOpt)
	if err != nil {
		return diag.Errorf("error updating CCI v2 Network: %s", err)
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
	deleteNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/networks/{name}"
	deleteNetworkPath := client.Endpoint + deleteNetworkHttpUrl
	deleteNetworkPath = strings.ReplaceAll(deleteNetworkPath, "{namespace}", ns)
	deleteNetworkPath = strings.ReplaceAll(deleteNetworkPath, "{name}", name)
	deleteNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteNetworkPath, &deleteNetworkOpt)
	if err != nil {
		return diag.Errorf("error deleting CCI v2 network: %s", err)
	}

	err = waitForDeleteV2DeploymentStatus(ctx, client, ns, name, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForCreateV2DeploymentStatus(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Ready"},
		Refresh:      refreshCreateV2DeploymentStatus(client, ns, name),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of the CCI network to complete: %s", err)
	}
	return nil
}

func refreshCreateV2DeploymentStatus(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetV2Deployment(client, ns, name)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("status.status", resp, "").(string)
		if status != "Ready" {
			return resp, "Pending", nil
		}

		return resp, status, nil
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
		return fmt.Errorf("error waiting for the status of the CCI v2 network to complete: %s", err)
	}
	return nil
}

func refreshDeleteV2DeploymentStatus(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetV2Deployment(client, ns, name)
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[DEBUG] successfully deleted CCI network: %s", name)
			return "", "Deleted", nil
		}
		return resp, "Pending", nil
	}
}

func resourceV2DeploymentImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<namespace>/<name>', but '%s'", importedId)
	}

	d.Set("namespace", parts[0])
	d.Set("name", parts[1])

	return []*schema.ResourceData{d}, nil
}

func GetV2Deployment(client *golangsdk.ServiceClient, namespace, name string) (interface{}, error) {
	getV2DeploymentHttpUrl := "apis/cci/v2/namespaces/{namespace}/deployment/{name}"
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
