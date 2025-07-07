package cceautopilot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var autopilotClusterNonUpdatableParams = []string{
	"name", "flavor",
	"host_network", "host_network.*.vpc", "host_network.*.subnet",
	"container_network", "container_network.*.mode",
	"annotations", "category", "type", "version", "custom_san", "enable_snat",
	"enable_swr_image_access", "enable_autopilot", "ipv6_enable",
	"service_network", "service_network.*.ipv4_cidr",
	"authentication", "authentication.*.mode",
	"kube_proxy_mode",
	"extend_param", "extend_param.*.enterprise_project_id",
	"configurations_override", "configurations_override.*.name", "configurations_override.*.configurations",
	"configurations_override.*.configurations.*.name", "configurations_override.*.configurations.*.value",
	"deletion_protection",
}

// @API CCE POST /autopilot/v3/projects/{project_id}/clusters
// @API CCE GET /autopilot/v3/projects/{project_id}/jobs/{job_id}
// @API CCE GET /autopilot/v3/projects/{project_id}/clusters/{cluster_id}
// @API CCE PUT /autopilot/v3/projects/{project_id}/clusters/{cluster_id}
// @API CCE PUT /autopilot/v3/projects/{project_id}/clusters/{cluster_id}/mastereip
// @API CCE DELETE /autopilot/v3/projects/{project_id}/clusters/{cluster_id}
// @API CCE POST /autopilot/v3/projects/{project_id}/clusters/{cluster_id}/tags/create
// @API CCE POST /autopilot/v3/projects/{project_id}/clusters/{cluster_id}/tags/delete
func ResourceAutopilotCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAutopilotClusterCreate,
		ReadContext:   resourceAutopilotClusterRead,
		UpdateContext: resourceAutopilotClusterUpdate,
		DeleteContext: resourceAutopilotClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(autopilotClusterNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_network": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subnet": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"container_network": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressVersionDiffs,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"custom_san": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"enable_snat": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_swr_image_access": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_autopilot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"eni_network": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnets": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"service_network": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_cidr": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"authentication": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"kube_proxy_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"extend_param": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"configurations_override": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"configurations": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"delete_efs": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_eni": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_net": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_obs": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_sfs30": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lts_reclaim_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"platform_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az": {
				Type:     schema.TypeString,
				Computed: true,
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
						"endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
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
	}
}

func buildClusterBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       "Cluster",
		"apiVersion": "v3",
		"metadata":   buildMetadataBodyParams(d),
		"spec":       buildSpecBodyParams(d),
	}

	return bodyParams
}

func buildMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"alias":       utils.ValueIgnoreEmpty(d.Get("alias")),
		"annotations": utils.ValueIgnoreEmpty(d.Get("annotations")),
	}

	return bodyParams
}

func buildSpecBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"category":               utils.ValueIgnoreEmpty(d.Get("category")),
		"type":                   utils.ValueIgnoreEmpty(d.Get("type")),
		"flavor":                 d.Get("flavor"),
		"version":                utils.ValueIgnoreEmpty(d.Get("version")),
		"description":            utils.ValueIgnoreEmpty(d.Get("description")),
		"customSan":              utils.ValueIgnoreEmpty(d.Get("custom_san")),
		"enableSnat":             d.Get("enable_snat"),
		"enableSWRImageAccess":   d.Get("enable_swr_image_access"),
		"enableAutopilot":        d.Get("enable_autopilot"),
		"ipv6enable":             d.Get("ipv6_enable"),
		"hostNetwork":            buildHostNetworkBodyParams(d),
		"containerNetwork":       buildContainerNetworkBodyParams(d),
		"eniNetwork":             buildEniNetworkBodyParams(d),
		"serviceNetwork":         buildServiceNetworkBodyParams(d),
		"authentication":         buildAuthenticationBodyParams(d),
		"clusterTags":            utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"kubeProxyMode":          utils.ValueIgnoreEmpty(d.Get("kube_proxy_mode")),
		"extendParam":            buildExtendParamBodyParams(d),
		"configurationsOverride": buildConfigurationsOverrideBodyParams(d),
		"deleteProtection":       d.Get("delete_protection"),
	}

	return bodyParams
}

func buildHostNetworkBodyParams(d *schema.ResourceData) map[string]interface{} {
	hostNetwork := d.Get("host_network").([]interface{})
	if len(hostNetwork) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"vpc":    utils.PathSearch("vpc", hostNetwork[0], nil),
		"subnet": utils.PathSearch("subnet", hostNetwork[0], nil),
	}

	return bodyParams
}

func buildContainerNetworkBodyParams(d *schema.ResourceData) map[string]interface{} {
	containerNetwork := d.Get("container_network").([]interface{})
	if len(containerNetwork) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"mode": utils.PathSearch("mode", containerNetwork[0], nil),
	}

	return bodyParams
}

func buildEniNetworkBodyParams(d *schema.ResourceData) map[string]interface{} {
	eniNetwork := d.Get("eni_network").([]interface{})
	if len(eniNetwork) == 0 {
		return nil
	}

	subnetsRaw := utils.PathSearch("subnets", eniNetwork[0], []interface{}{}).([]interface{})
	subnets := make([]map[string]interface{}, len(subnetsRaw))
	for i, v := range subnetsRaw {
		subnets[i] = map[string]interface{}{
			"subnetID": utils.PathSearch("subnet_id", v, nil),
		}
	}

	bodyParams := map[string]interface{}{
		"subnets": subnets,
	}

	return bodyParams
}

func buildServiceNetworkBodyParams(d *schema.ResourceData) map[string]interface{} {
	serviceNetwork := d.Get("service_network").([]interface{})
	if len(serviceNetwork) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"IPv4CIDR": utils.PathSearch("ipv4_cidr", serviceNetwork[0], nil),
	}

	return bodyParams
}

func buildAuthenticationBodyParams(d *schema.ResourceData) map[string]interface{} {
	authentication := d.Get("authentication").([]interface{})
	if len(authentication) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"mode": utils.PathSearch("mode", authentication[0], nil),
	}

	return bodyParams
}

func buildExtendParamBodyParams(d *schema.ResourceData) map[string]interface{} {
	extendParam := d.Get("extend_param").([]interface{})
	if len(extendParam) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"enterpriseProjectId": utils.PathSearch("enterprise_project_id", extendParam[0], nil),
	}

	return bodyParams
}

func buildConfigurationsOverrideBodyParams(d *schema.ResourceData) []map[string]interface{} {
	configurationsOverrideRaw := d.Get("configurations_override").([]interface{})
	if len(configurationsOverrideRaw) == 0 {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(configurationsOverrideRaw))

	for i, v := range configurationsOverrideRaw {
		bodyParams[i] = map[string]interface{}{
			"name":           utils.PathSearch("name", v, nil),
			"configurations": buildConfigurationsBodyParams(v),
		}
	}

	return bodyParams
}

func buildConfigurationsBodyParams(configurationsOverride interface{}) []map[string]interface{} {
	configurationsRaw := utils.PathSearch("configurations", configurationsOverride, []interface{}{}).([]interface{})
	if len(configurationsRaw) == 0 {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(configurationsRaw))

	for i, v := range configurationsRaw {
		bodyParams[i] = map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
	}

	return bodyParams
}

func resourceAutopilotClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createClusterHttpUrl = "autopilot/v3/projects/{project_id}/clusters"
		createClusterProduct = "cce"
	)
	createClusterClient, err := cfg.NewServiceClient(createClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	createClusterPath := createClusterClient.Endpoint + createClusterHttpUrl
	createClusterPath = strings.ReplaceAll(createClusterPath, "{project_id}", createClusterClient.ProjectID)

	createClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createOpts := buildClusterBodyParams(d)
	createClusterOpt.JSONBody = utils.RemoveNil(createOpts)
	createClusterResp, err := createClusterClient.Request("POST", createClusterPath, &createClusterOpt)
	if err != nil {
		return diag.Errorf("error creating CCE autopolit cluster: %s", err)
	}

	createClusterRespBody, err := utils.FlattenResponse(createClusterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("metadata.uid", createClusterRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CCE autopilot cluster: ID is not found in API response")
	}
	d.SetId(id)

	jobID := utils.PathSearch("status.jobID", createClusterRespBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error creating CCE autopilot cluster: jobID is not found in API response")
	}

	err = clusterJobWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate), jobID)
	if err != nil {
		return diag.Errorf("error waiting for creating CCE autopilot cluster (%s) to complete: %s", id, err)
	}

	if v, ok := d.GetOk("eip_id"); ok {
		err = clusterEipAction(createClusterClient, id, v.(string))
		if err != nil {
			return diag.Errorf("error binding EIP to CCE autopilot cluster (%s): %s", id, err)
		}
	}

	return resourceAutopilotClusterRead(ctx, d, meta)
}

func resourceAutopilotClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getClusterHttpUrl = "autopilot/v3/projects/{project_id}/clusters/{cluster_id}"
		getClusterProduct = "cce"
	)
	getClusterClient, err := cfg.NewServiceClient(getClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	getClusterPath := getClusterClient.Endpoint + getClusterHttpUrl
	getClusterPath = strings.ReplaceAll(getClusterPath, "{project_id}", getClusterClient.ProjectID)
	getClusterPath = strings.ReplaceAll(getClusterPath, "{cluster_id}", d.Id())

	getClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getClusterResp, err := getClusterClient.Request("GET", getClusterPath, &getClusterOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE autopolit cluster")
	}

	getClusterRespBody, err := utils.FlattenResponse(getClusterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// enable_snat and enable_swr_image_access not saved, because thay are not returned
	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", utils.PathSearch("metadata.name", getClusterRespBody, nil)),
		d.Set("alias", utils.PathSearch("metadata.alias", getClusterRespBody, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", getClusterRespBody, nil)),
		d.Set("category", utils.PathSearch("spec.category", getClusterRespBody, nil)),
		d.Set("type", utils.PathSearch("spec.type", getClusterRespBody, nil)),
		d.Set("flavor", utils.PathSearch("spec.flavor", getClusterRespBody, nil)),
		d.Set("version", utils.PathSearch("spec.version", getClusterRespBody, nil)),
		d.Set("description", utils.PathSearch("spec.description", getClusterRespBody, nil)),
		d.Set("custom_san", utils.PathSearch("spec.customSan", getClusterRespBody, nil)),
		d.Set("enable_autopilot", utils.PathSearch("spec.enableAutopilot", getClusterRespBody, nil)),
		d.Set("ipv6_enable", utils.PathSearch("spec.ipv6enable", getClusterRespBody, nil)),
		d.Set("host_network", flattenHostNetwork(getClusterRespBody)),
		d.Set("container_network", flattenContainerNetwork(getClusterRespBody)),
		d.Set("eni_network", flattenEniNetwork(getClusterRespBody)),
		d.Set("service_network", flattenServiceNetwork(getClusterRespBody)),
		d.Set("authentication", flattenAuthentication(getClusterRespBody)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("spec.clusterTags", getClusterRespBody, nil))),
		d.Set("kube_proxy_mode", utils.PathSearch("spec.kubeProxyMode", getClusterRespBody, nil)),
		d.Set("az", utils.PathSearch("spec.az", getClusterRespBody, nil)),
		d.Set("extend_param", flattenExtendParam(getClusterRespBody)),
		d.Set("configurations_override", flattenConfigurationsOverride(getClusterRespBody)),
		d.Set("deletion_protection", utils.PathSearch("spec.deletionProtection", getClusterRespBody, nil)),
		d.Set("platform_version", utils.PathSearch("spec.platformVersion", getClusterRespBody, nil)),
		d.Set("created_at", utils.PathSearch("metadata.creationTimestamp", getClusterRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("metadata.updateTimestamp", getClusterRespBody, nil)),
		d.Set("status", flattenStatus(getClusterRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHostNetwork(getClusterRespBody interface{}) []map[string]interface{} {
	hostNetwork := utils.PathSearch("spec.hostNetwork", getClusterRespBody, nil)
	if hostNetwork == nil {
		return nil
	}

	res := []map[string]interface{}{
		{
			"vpc":    utils.PathSearch("vpc", hostNetwork, nil),
			"subnet": utils.PathSearch("subnet", hostNetwork, nil),
		},
	}

	return res
}

func flattenContainerNetwork(getClusterRespBody interface{}) []map[string]interface{} {
	containerNetwork := utils.PathSearch("spec.containerNetwork", getClusterRespBody, nil)
	if containerNetwork == nil {
		return nil
	}

	res := []map[string]interface{}{
		{
			"mode": utils.PathSearch("mode", containerNetwork, nil),
		},
	}

	return res
}

func flattenEniNetwork(getClusterRespBody interface{}) []map[string]interface{} {
	eniNetwork := utils.PathSearch("spec.eniNetwork", getClusterRespBody, nil)
	if eniNetwork == nil {
		return nil
	}

	subnetsRaw := utils.PathSearch("subnets", eniNetwork, []interface{}{}).([]interface{})
	subnets := make([]map[string]interface{}, len(subnetsRaw))
	for i, v := range subnetsRaw {
		subnets[i] = map[string]interface{}{
			"subnet_id": utils.PathSearch("subnetID", v, nil),
		}
	}

	res := []map[string]interface{}{
		{
			"subnets": subnets,
		},
	}

	return res
}

func flattenServiceNetwork(getClusterRespBody interface{}) []map[string]interface{} {
	serviceNetwork := utils.PathSearch("spec.serviceNetwork", getClusterRespBody, nil)
	if serviceNetwork == nil {
		return nil
	}

	res := []map[string]interface{}{
		{
			"ipv4_cidr": utils.PathSearch("IPv4CIDR", serviceNetwork, nil),
		},
	}

	return res
}

func flattenAuthentication(getClusterRespBody interface{}) []map[string]interface{} {
	authentication := utils.PathSearch("spec.authentication", getClusterRespBody, nil)
	if authentication == nil {
		return nil
	}

	res := []map[string]interface{}{
		{
			"mode": utils.PathSearch("mode", authentication, nil),
		},
	}

	return res
}

func flattenExtendParam(getClusterRespBody interface{}) []map[string]interface{} {
	extendParam := utils.PathSearch("spec.extendParam", getClusterRespBody, nil)
	if extendParam == nil {
		return nil
	}

	res := []map[string]interface{}{
		{
			"enterprise_project_id": utils.PathSearch("enterpriseProjectId", extendParam, nil),
		},
	}

	return res
}

func flattenConfigurationsOverride(getClusterRespBody interface{}) []map[string]interface{} {
	configurationsOverrideRaw := utils.PathSearch("spec.configurationsOverride", getClusterRespBody, []interface{}{}).([]interface{})
	if len(configurationsOverrideRaw) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(configurationsOverrideRaw))
	for i, v := range configurationsOverrideRaw {
		res[i] = map[string]interface{}{
			"name":           utils.PathSearch("name", v, nil),
			"configurations": flattenConfigurations(v),
		}
	}

	return res
}

func flattenConfigurations(configurationsOverride interface{}) []map[string]interface{} {
	configurationsRaw := utils.PathSearch("configurations", configurationsOverride, []interface{}{}).([]interface{})
	if len(configurationsRaw) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(configurationsRaw))
	for i, v := range configurationsRaw {
		res[i] = map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
	}

	return res
}

func flattenStatus(getClusterRespBody interface{}) []map[string]interface{} {
	status := utils.PathSearch("status", getClusterRespBody, nil)
	if status == nil {
		return nil
	}

	endpointsRaw := utils.PathSearch("endpoints", status, []interface{}{}).([]interface{})
	endpoints := make([]map[string]interface{}, len(endpointsRaw))
	for i, v := range endpointsRaw {
		endpoints[i] = map[string]interface{}{
			"url":  utils.PathSearch("url", v, nil),
			"type": utils.PathSearch("type", v, nil),
		}
	}

	res := []map[string]interface{}{
		{
			"endpoints": endpoints,
			"phase":     utils.PathSearch("phase", status, nil),
		},
	}

	return res
}

func buildUpdateClusterBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{}

	if d.HasChange("alias") {
		bodyParams["metadata"] = buildUpdateMetadataBodyParams(d)
	}

	if d.HasChanges("description", "custom_san", "eni_network") {
		bodyParams["spec"] = buildUpdateSpecBodyParams(d)
	}

	return bodyParams
}

func buildUpdateMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alias": utils.ValueIgnoreEmpty(d.Get("alias")),
	}

	return bodyParams
}

func buildUpdateSpecBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{}

	if d.HasChange("description") {
		bodyParams["description"] = d.Get("description")
	}

	if d.HasChange("custom_san") {
		bodyParams["customSan"] = d.Get("custom_san")
	}

	if d.HasChange("eni_network") {
		bodyParams["eniNetwork"] = buildEniNetworkBodyParams(d)
	}

	return bodyParams
}

func buildUpdateClusterEipBodyParams(eipID string) map[string]interface{} {
	if eipID != "" {
		return map[string]interface{}{
			"spec": map[string]interface{}{
				"action": "bind",
				"spec": map[string]interface{}{
					"id": eipID,
				},
			},
		}
	}

	return map[string]interface{}{
		"spec": map[string]interface{}{
			"action": "unbind",
		},
	}
}

func clusterEipAction(updateClusterClient *golangsdk.ServiceClient, clusterID, eipID string) error {
	var updateClusterEipHttpUrl = "autopilot/v3/projects/{project_id}/clusters/{cluster_id}/mastereip"

	updateClusterEipPath := updateClusterClient.Endpoint + updateClusterEipHttpUrl
	updateClusterEipPath = strings.ReplaceAll(updateClusterEipPath, "{project_id}", updateClusterClient.ProjectID)
	updateClusterEipPath = strings.ReplaceAll(updateClusterEipPath, "{cluster_id}", clusterID)

	updateOpts := buildUpdateClusterEipBodyParams(eipID)
	updateClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(updateOpts),
	}

	_, err := updateClusterClient.Request("PUT", updateClusterEipPath, &updateClusterOpt)
	if err != nil {
		return err
	}
	return nil
}

func buildUpdateClusterTagsBodyParams(action string, tags map[string]interface{}) map[string]interface{} {
	taglist := make([]map[string]interface{}, 0, len(tags))
	if action == "create" {
		for k, v := range tags {
			tag := map[string]interface{}{
				"key":   k,
				"value": v,
			}
			taglist = append(taglist, tag)
		}
	} else {
		for k := range tags {
			tag := map[string]interface{}{
				"key": k,
			}
			taglist = append(taglist, tag)
		}
	}

	return map[string]interface{}{
		"tags": taglist,
	}
}

func clusterTagsAction(updateClusterClient *golangsdk.ServiceClient, clusterID, action string, tags map[string]interface{}) error {
	var updateClusterTagsHttpUrl = "autopilot/v3/projects/{project_id}/clusters/{cluster_id}/tags/{action}"

	updateClusterTagsPath := updateClusterClient.Endpoint + updateClusterTagsHttpUrl
	updateClusterTagsPath = strings.ReplaceAll(updateClusterTagsPath, "{project_id}", updateClusterClient.ProjectID)
	updateClusterTagsPath = strings.ReplaceAll(updateClusterTagsPath, "{cluster_id}", clusterID)
	updateClusterTagsPath = strings.ReplaceAll(updateClusterTagsPath, "{action}", action)

	updateOpts := buildUpdateClusterTagsBodyParams(action, tags)
	updateClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(updateOpts),
		OkCodes:          []int{204},
	}

	_, err := updateClusterClient.Request("POST", updateClusterTagsPath, &updateClusterOpt)
	if err != nil {
		return err
	}
	return nil
}

func resourceAutopilotClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	id := d.Id()

	var updateClusterProduct = "cce"

	updateClusterClient, err := cfg.NewServiceClient(updateClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	if d.HasChanges("alias", "description", "custom_san", "eni_network") {
		var updateClusterHttpUrl = "autopilot/v3/projects/{project_id}/clusters/{cluster_id}"

		updateClusterPath := updateClusterClient.Endpoint + updateClusterHttpUrl
		updateClusterPath = strings.ReplaceAll(updateClusterPath, "{project_id}", updateClusterClient.ProjectID)
		updateClusterPath = strings.ReplaceAll(updateClusterPath, "{cluster_id}", id)

		updateOpts := buildUpdateClusterBodyParams(d)
		updateClusterOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(updateOpts),
		}

		_, err := updateClusterClient.Request("PUT", updateClusterPath, &updateClusterOpt)
		if err != nil {
			return diag.Errorf("error updating CCE autopolit cluster: %s", err)
		}
	}

	if d.HasChange("eip_id") {
		oldEip, newEip := d.GetChange("eip_id")

		if oldEip.(string) != "" {
			err = clusterEipAction(updateClusterClient, id, "")
			if err != nil {
				return diag.Errorf("error unbinding EIP from CCE autopilot cluster (%s): %s", id, err)
			}
		}

		if newEip.(string) != "" {
			err = clusterEipAction(updateClusterClient, id, newEip.(string))
			if err != nil {
				return diag.Errorf("error binding EIP to CCE autopilot cluster (%s): %s", id, err)
			}
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		if len(oldTags.(map[string]interface{})) > 0 {
			err = clusterTagsAction(updateClusterClient, id, "delete", oldTags.(map[string]interface{}))
			if err != nil {
				return diag.Errorf("error removing tags from CCE autopilot cluster (%s): %s", id, err)
			}
		}

		if len(newTags.(map[string]interface{})) > 0 {
			err = clusterTagsAction(updateClusterClient, id, "create", newTags.(map[string]interface{}))
			if err != nil {
				return diag.Errorf("error adding tags to CCE autopilot cluster (%s): %s", id, err)
			}
		}
	}

	return resourceAutopilotClusterRead(ctx, d, meta)
}

func resourceAutopilotClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteClusterHttpUrl = "autopilot/v3/projects/{project_id}/clusters/{cluster_id}"
		deleteClusterProduct = "cce"
	)
	deleteClusterClient, err := cfg.NewServiceClient(deleteClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	deleteClusterPath := deleteClusterClient.Endpoint + deleteClusterHttpUrl
	deleteClusterPath = strings.ReplaceAll(deleteClusterPath, "{project_id}", deleteClusterClient.ProjectID)
	deleteClusterPath = strings.ReplaceAll(deleteClusterPath, "{cluster_id}", d.Id())

	deleteClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteClusteQueryParams := buildDeleteClusteQueryParams(d)
	deleteClusterPath += deleteClusteQueryParams

	deleteClusterResp, err := deleteClusterClient.Request("DELETE", deleteClusterPath, &deleteClusterOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CCE autopolit cluster")
	}

	deleteClusterRespBody, err := utils.FlattenResponse(deleteClusterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobID := utils.PathSearch("status.jobID", deleteClusterRespBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error deleting CCE autopilot cluster: jobID is not found in API response")
	}

	err = clusterJobWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete), jobID)
	if err != nil {
		return diag.Errorf("error waiting for deleting CCE autopilot cluster (%s) to complete: %s", d.Id(), err)
	}

	return nil
}

func buildDeleteClusteQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("delete_efs"); ok {
		res = fmt.Sprintf("%s&delete_efs=%v", res, v)
	}
	if v, ok := d.GetOk("delete_eni"); ok {
		res = fmt.Sprintf("%s&delete_eni=%v", res, v)
	}
	if v, ok := d.GetOk("delete_net"); ok {
		res = fmt.Sprintf("%s&delete_net=%v", res, v)
	}
	if v, ok := d.GetOk("delete_obs"); ok {
		res = fmt.Sprintf("%s&delete_obs=%v", res, v)
	}
	if v, ok := d.GetOk("delete_sfs30"); ok {
		res = fmt.Sprintf("%s&delete_sfs30=%v", res, v)
	}
	if v, ok := d.GetOk("lts_reclaim_policy"); ok {
		res = fmt.Sprintf("%s&lts_reclaim_policy=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func clusterJobWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration, jobID string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				clusterJobWaitingHttpUrl = "autopilot/v3/projects/{project_id}/jobs/{job_id}"
				clusterJobWaitingProduct = "cce"
			)
			clusterJobWaitingClient, err := cfg.NewServiceClient(clusterJobWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CCE client: %s", err)
			}

			clusterJobWaitingPath := clusterJobWaitingClient.Endpoint + clusterJobWaitingHttpUrl
			clusterJobWaitingPath = strings.ReplaceAll(clusterJobWaitingPath, "{project_id}", clusterJobWaitingClient.ProjectID)
			clusterJobWaitingPath = strings.ReplaceAll(clusterJobWaitingPath, "{job_id}", jobID)

			clusterJobWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			clusterJobWaitingResp, err := clusterJobWaitingClient.Request("GET", clusterJobWaitingPath, &clusterJobWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			clusterJobWaitingRespBody, err := utils.FlattenResponse(clusterJobWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`status.phase`, clusterJobWaitingRespBody, nil)
			if status == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `status.phase`)
			}

			targetStatus := []string{
				"Success",
			}
			if utils.StrSliceContains(targetStatus, status.(string)) {
				return clusterJobWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"Failed",
			}
			if utils.StrSliceContains(unexpectedStatus, status.(string)) {
				return clusterJobWaitingRespBody, status.(string), nil
			}

			return clusterJobWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
