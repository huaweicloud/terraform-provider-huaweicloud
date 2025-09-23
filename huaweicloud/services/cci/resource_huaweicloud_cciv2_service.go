package cci

import (
	"context"
	"fmt"
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

var serviceNonUpdatableParams = []string{"namespace", "name"}

// @API CCI POST /apis/cci/v2/namespaces/{namespace}/services
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/services/{name}
// @API CCI PUT /apis/cci/v2/namespaces/{namespace}/services/{name}
// @API CCI DELETE /apis/cci/v2/namespaces/{namespace}/services/{name}
func ResourceV2Service() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2ServiceCreate,
		ReadContext:   resourceV2ServiceRead,
		UpdateContext: resourceV2ServiceUpdate,
		DeleteContext: resourceV2ServiceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2ServiceImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(serviceNonUpdatableParams),

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
				Description: `Specifies the namespace of the CCI Service.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the CCI Service.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the annotations of the CCI Service.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the labels of the CCI Service.`,
			},
			"ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The app protocol.`,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The name.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The port.`,
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The protocol.`,
						},
						"target_port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The target port.`,
						},
					},
				},
				Description: `Specifies the ports of the CCI Service.`,
			},
			"selector": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the selector of the CCI Service.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of the CCI Service.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the CCI Service.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the CCI Service.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the namespace.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the namespace.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the namespace.`,
			},
			"finalizers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The finalizers of the namespace.`,
			},
			"cluster_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the cluster IP of the CCI Service.`,
			},
			"cluster_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the cluster IPs of the CCI Service.`,
			},
			"external_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The external name of the CCI Service.`,
			},
			"ip_families": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The IP families of the CCI Service.`,
			},
			"ip_family_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP family policy of the CCI Service.`,
			},
			"load_balancer_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The load balancer IP of the CCI Service.`,
			},
			"publish_not_ready_addresses": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the publish is not ready addresses of the CCI Service.`,
			},
			"session_affinity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The session affinity of the CCI Service.`,
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
										Description: `The type.`,
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Tthe status.`,
									},
									"observe_generation": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The observe generation.`,
									},
									"last_transition_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The last transition time.`,
									},
									"reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The reason.`,
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The message.`,
									},
								},
							},
							Description: `Tthe conditions of the CCI Service.`,
						},
						"loadbalancer": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ingress": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The IP of the loadbalancer.`,
												},
												"ports": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"error": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: `The error.`,
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: `The port.`,
															},
															"protocol": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: `The protocol.`,
															},
														},
													},
													Description: `The ports of the loadbalancer.`,
												},
											},
										},
										Description: `The ingress of the loadbalancer.`,
									},
								},
							},
							Description: `The loadbalancer of the CCI Service.`,
						},
					},
				},
				Description: `The status of the namespace.`,
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

func resourceV2ServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createServiceHttpUrl := "apis/cci/v2/namespaces/{namespace}/services"
	createServicePath := client.Endpoint + createServiceHttpUrl
	createServicePath = strings.ReplaceAll(createServicePath, "{namespace}", d.Get("namespace").(string))
	createServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createServiceOpt.JSONBody = utils.RemoveNil(buildCreateV2ServiceParams(d))

	resp, err := client.Request("POST", createServicePath, &createServiceOpt)
	if err != nil {
		return diag.Errorf("error creating CCI Service: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	ns := utils.PathSearch("metadata.namespace", respBody, "").(string)
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if ns == "" || name == "" {
		return diag.Errorf("unable to find CCI Service name or namespace from API response")
	}
	d.SetId(ns + "/" + name)

	return resourceV2ServiceRead(ctx, d, meta)
}

func buildCreateV2ServiceParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":        d.Get("name"),
			"namespace":   d.Get("namespace"),
			"annotations": d.Get("annotations"),
		},
		"spec": map[string]interface{}{
			"ports":    buildCreateV2ServiceSpecPorts(d.Get("ports").(*schema.Set).List()),
			"selector": d.Get("selector"),
			"type":     d.Get("type"),
		},
	}

	return bodyParams
}

func buildCreateV2ServiceSpecPorts(ports []interface{}) []interface{} {
	if len(ports) == 0 {
		return nil
	}

	rst := make([]interface{}, len(ports))
	for i, v := range ports {
		rst[i] = map[string]interface{}{
			"appProtocol": utils.ValueIgnoreEmpty(utils.PathSearch("app_protocol", v, nil)),
			"name":        utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"port":        utils.PathSearch("port", v, nil),
			"protocol":    utils.ValueIgnoreEmpty(utils.PathSearch("protocol", v, nil)),
			"targetPort":  utils.ValueIgnoreEmpty(utils.PathSearch("target_port", v, nil)),
		}
	}

	return rst
}

func resourceV2ServiceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	resp, err := GetService(client, ns, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 Service")
	}

	mErr := multierror.Append(
		d.Set("namespace", utils.PathSearch("metadata.namespace", resp, nil)),
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("labels", utils.PathSearch("metadata.labels", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("finalizers", utils.PathSearch("metadata.finalizers", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("cluster_ip", utils.PathSearch("spec.clusterIP", resp, nil)),
		d.Set("cluster_ips", utils.PathSearch("spec.clusterIPs", resp, nil)),
		d.Set("external_name", utils.PathSearch("spec.externalName", resp, nil)),
		d.Set("ip_families", utils.PathSearch("spec.ipFamilies", resp, nil)),
		d.Set("ip_family_policy", utils.PathSearch("spec.ipFamilyPolicy", resp, nil)),
		d.Set("load_balancer_ip", utils.PathSearch("spec.loadBalancerIP", resp, nil)),
		d.Set("publish_not_ready_addresses", utils.PathSearch("spec.publishNotReadyAddresses", resp, nil)),
		d.Set("selector", utils.PathSearch("spec.selector", resp, nil)),
		d.Set("session_affinity", utils.PathSearch("spec.sessionAffinity", resp, nil)),
		d.Set("type", utils.PathSearch("spec.type", resp, nil)),
		d.Set("ports", flattenServiceSpecPorts(utils.PathSearch("spec.ports", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("status", flattenServiceStatus(utils.PathSearch("status", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenServiceSpecPorts(ports []interface{}) []interface{} {
	if len(ports) == 0 {
		return nil
	}

	rst := make([]interface{}, len(ports))
	for i, v := range ports {
		rst[i] = map[string]interface{}{
			"app_protocol": utils.PathSearch("appProtocol", v, nil),
			"name":         utils.PathSearch("name", v, nil),
			"port":         utils.PathSearch("port", v, nil),
			"protocol":     utils.PathSearch("protocol", v, nil),
			"target_port":  utils.PathSearch("targetPort", v, nil),
		}
	}

	return rst
}

func flattenServiceStatus(status interface{}) []map[string]interface{} {
	if status == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"loadbalancer": flattenServiceStatusLoadBalancer(utils.PathSearch("loadbalancer", status, nil)),
			"conditions":   flattenServiceStatusConditions(utils.PathSearch("conditions", status, make([]interface{}, 0)).([]interface{})),
		},
	}

	return rst
}

func flattenServiceStatusConditions(conditions []interface{}) []interface{} {
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

func flattenServiceStatusLoadBalancer(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"ingress": flattenServiceStatusLoadBalancerIngress(utils.PathSearch("ingress", resp, make([]interface{}, 0)).([]interface{})),
	})

	return rst
}

func flattenServiceStatusLoadBalancerIngress(ingress []interface{}) []interface{} {
	if len(ingress) == 0 {
		return nil
	}

	rst := make([]interface{}, len(ingress))
	for i, v := range ingress {
		rst[i] = map[string]interface{}{
			"ip":    utils.PathSearch("ip", v, nil),
			"ports": flattenServiceStatusLoadBalancerIngressPorts(utils.PathSearch("ports", v, make([]interface{}, 0)).([]interface{})),
		}
	}

	return rst
}

func flattenServiceStatusLoadBalancerIngressPorts(ports []interface{}) []interface{} {
	if len(ports) == 0 {
		return nil
	}

	rst := make([]interface{}, len(ports))
	for i, v := range ports {
		rst[i] = map[string]interface{}{
			"error":    utils.PathSearch("error", v, nil),
			"port":     utils.PathSearch("port", v, nil),
			"protocol": utils.PathSearch("protocol", v, nil),
		}
	}

	return rst
}

func resourceV2ServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	updateServiceHttpUrl := "apis/cci/v2/namespaces/{namespace}/services/{name}"
	updateServicePath := client.Endpoint + updateServiceHttpUrl
	updateServicePath = strings.ReplaceAll(updateServicePath, "{namespace}", d.Get("namespace").(string))
	updateServicePath = strings.ReplaceAll(updateServicePath, "{name}", d.Get("name").(string))
	updateServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateServiceOpt.JSONBody = utils.RemoveNil(buildUpdateV2ServiceParams(d))

	_, err = client.Request("PUT", updateServicePath, &updateServiceOpt)
	if err != nil {
		return diag.Errorf("error updating CCI v2 Service: %s", err)
	}
	return resourceV2ServiceRead(ctx, d, meta)
}

func buildUpdateV2ServiceParams(d *schema.ResourceData) map[string]interface{} {
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
			"labels":            d.Get("labels"),
			"finalizers":        d.Get("finalizers"),
		},
		"spec": map[string]interface{}{
			"ports":           buildCreateV2ServiceSpecPorts(d.Get("ports").(*schema.Set).List()),
			"selector":        utils.ValueIgnoreEmpty(d.Get("selector")),
			"type":            utils.ValueIgnoreEmpty(d.Get("type")),
			"sessionAffinity": utils.ValueIgnoreEmpty(d.Get("session_affinity")),
		},
		"status": map[string]interface{}{
			"loadBalancer": buildServiceStatusLoadBalancer(d.Get("status.0.loadbalancer.0")),
			"conditions":   buildServiceStatusConditions(d.Get("status.0.conditions").([]interface{})),
		},
	}

	return bodyParams
}

func buildServiceStatusLoadBalancer(loadBalancer interface{}) interface{} {
	if loadBalancer == nil {
		return nil
	}

	rst := map[string]interface{}{
		"ingress": buildServiceStatusLoadBalancerIngress(utils.PathSearch("ingress", loadBalancer, make([]interface{}, 0)).([]interface{})),
	}

	return rst
}

func buildServiceStatusLoadBalancerIngress(ingress []interface{}) []interface{} {
	if len(ingress) == 0 {
		return nil
	}

	rst := make([]interface{}, len(ingress))
	for i, v := range ingress {
		rst[i] = map[string]interface{}{
			"ip":    utils.PathSearch("ip", v, nil),
			"ports": buildServiceStatusLoadBalancerIngressPorts(utils.PathSearch("ports", v, make([]interface{}, 0)).([]interface{})),
		}
	}

	return rst
}

func buildServiceStatusLoadBalancerIngressPorts(ports []interface{}) []interface{} {
	if len(ports) == 0 {
		return nil
	}

	rst := make([]interface{}, len(ports))
	for i, v := range ports {
		rst[i] = map[string]interface{}{
			"error":    utils.PathSearch("error", v, nil),
			"port":     utils.PathSearch("port", v, nil),
			"protocol": utils.PathSearch("protocol", v, nil),
		}
	}

	return rst
}

func buildServiceStatusConditions(conditions []interface{}) []interface{} {
	if len(conditions) == 0 {
		return nil
	}

	rst := make([]interface{}, len(conditions))
	for i, v := range conditions {
		rst[i] = map[string]interface{}{
			"type":               utils.PathSearch("type", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"observeGeneration":  utils.PathSearch("observe_generation", v, nil),
			"lastTransitionTime": utils.PathSearch("last_transition_time", v, nil),
			"reason":             utils.PathSearch("reason", v, nil),
			"message":            utils.PathSearch("message", v, nil),
		}
	}

	return rst
}

func resourceV2ServiceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	deleteServiceHttpUrl := "apis/cci/v2/namespaces/{namespace}/services/{name}"
	deleteServicePath := client.Endpoint + deleteServiceHttpUrl
	deleteServicePath = strings.ReplaceAll(deleteServicePath, "{namespace}", ns)
	deleteServicePath = strings.ReplaceAll(deleteServicePath, "{name}", name)
	deleteServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteServicePath, &deleteServiceOpt)
	if err != nil {
		return diag.Errorf("error deleting CCI v2 Service: %s", err)
	}

	return nil
}

func resourceV2ServiceImportState(_ context.Context, d *schema.ResourceData,
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

func GetService(client *golangsdk.ServiceClient, namespace, name string) (interface{}, error) {
	getServiceHttpUrl := "apis/cci/v2/namespaces/{namespace}/services/{name}"
	getServicePath := client.Endpoint + getServiceHttpUrl
	getServicePath = strings.ReplaceAll(getServicePath, "{namespace}", namespace)
	getServicePath = strings.ReplaceAll(getServicePath, "{name}", name)
	getServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getServiceResp, err := client.Request("GET", getServicePath, &getServiceOpt)
	if err != nil {
		return getServiceResp, err
	}

	return utils.FlattenResponse(getServiceResp)
}
