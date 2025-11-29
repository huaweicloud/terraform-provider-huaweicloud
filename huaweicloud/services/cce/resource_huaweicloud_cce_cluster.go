package cce

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/aom/v1/icagents"
	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"
	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/chnsz/golangsdk/openstack/common/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE POST /api/v3/projects/{project_id}/clusters
// @API CCE GET /api/v3/projects/{project_id}/clusters/{id}
// @API CCE PUT /api/v3/projects/{project_id}/clusters/{id}
// @API CCE DELETE /api/v3/projects/{project_id}/clusters/{id}
// @API CCE GET /api/v3/projects/{project_id}/jobs/{job_id}
// @API CCE POST /api/v3/projects/{project_id}/clusters/{id}/operation/{action}
// @API CCE POST /api/v3/projects/{project_id}/clusters/{id}/clustercert
// @API CCE PUT /api/v3/projects/{project_id}/clusters/{id}/mastereip
// @API CCE POST /api/v3/projects/{project_id}/clusters/{id}/tags/create
// @API CCE POST /api/v3/projects/{project_id}/clusters/{id}/tags/delete
// @API CCE POST /api/v3/projects/{project_id}/clusters/{id}/operation/resize
// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/operation/hibernate
// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/operation/awake
// @API BSS GET /V2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources/filter
// @API AOM POST /svcstg/icmgr/v1/{project_id}/agents

// ResourceCCEClusterV3 defines the CCE cluster resource schema and functions.
// Deprecated: It's a deprecated function, please refer to the function 'ResourceCluster'.
func ResourceCCEClusterV3() *schema.Resource {
	return ResourceCluster()
}

var associateDeleteSchema *schema.Schema = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
	ValidateFunc: validation.StringInSlice([]string{
		"true", "try", "false",
	}, true),
	ConflictsWith: []string{"delete_all"},
}

var associateDeleteSchemaInternal *schema.Schema = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
	ValidateFunc: validation.StringInSlice([]string{
		"true", "try", "false",
	}, true),
	ConflictsWith: []string{"delete_all"},
	Description:   "schema: Internal",
}

// ResourceCluster defines the CCE cluster resource schema and functions.
func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		//request and response parameters
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
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_version": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressVersionDiffs,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "VirtualMachine",
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "schema: Internal",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "schema: Internal",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"highway_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "schema: Internal",
			},
			"container_network_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_network_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"eni_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "the IPv4 subnet ID of the subnet where the ENI resides",
			},
			"eni_subnet_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Computed",
			},
			"enable_distribute_management": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"authentication_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "rbac",
			},
			"authenticating_proxy_ca": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"authenticating_proxy_cert", "authenticating_proxy_private_key"},
			},
			"authenticating_proxy_cert": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"authenticating_proxy_ca", "authenticating_proxy_private_key"},
			},
			"authenticating_proxy_private_key": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"authenticating_proxy_ca", "authenticating_proxy_cert"},
			},
			"multi_az": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"masters"},
			},
			"masters": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				MaxItems:      3,
				ConflictsWith: []string{"multi_az"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"eip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateIP,
			},
			"service_network_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"kube_proxy_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"extend_param": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "schema: Internal",
			},
			"extend_params": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"multi_az",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_az": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"dss_master_volumes": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"fix_pool_mask": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"dec_master_flavor": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"docker_umask_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cpu_manager_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"hibernate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"component_configurations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"configurations": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"encryption_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"custom_san": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),

			// charge info: charging_mode, period_unit, period, auto_renew, auto_pay
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			"auto_pay":      common.SchemaAutoPay(nil),

			"delete_efs": associateDeleteSchema,
			"delete_eni": associateDeleteSchemaInternal,
			"delete_evs": associateDeleteSchema,
			"delete_net": associateDeleteSchemaInternal,
			"delete_obs": associateDeleteSchema,
			"delete_sfs": associateDeleteSchema,
			"delete_all": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "try", "false",
				}, true),
			},
			"lts_reclaim_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kube_config_raw": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_istio": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{
					Computed: true,
				}),
			},
			"certificate_clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_authority_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"certificate_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_certificate_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_key_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			// Deprecated
			"billing_mode": {
				Type:       schema.TypeInt,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "use charging_mode instead",
			},
		},
	}
}

func resourceClusterLabels(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceClusterTags(d *schema.ResourceData) []tags.ResourceTag {
	tagRaw := d.Get("tags").(map[string]interface{})
	return utils.ExpandResourceTags(tagRaw)
}

func resourceClusterAnnotations(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("annotations").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceClusterExtendParam(d *schema.ResourceData) map[string]interface{} {
	extendParam := make(map[string]interface{})
	if v, ok := d.GetOk("extend_param"); ok {
		for key, val := range v.(map[string]interface{}) {
			extendParam[key] = val.(string)
		}
	}

	if multiAZ, ok := d.GetOk("multi_az"); ok && multiAZ.(bool) {
		extendParam["clusterAZ"] = "multi_az"
	}

	return extendParam
}

func resourceClusterExtendParams(extendParamsRaw []interface{}) map[string]interface{} {
	if len(extendParamsRaw) != 1 {
		return nil
	}

	if extendParams, ok := extendParamsRaw[0].(map[string]interface{}); ok {
		res := map[string]interface{}{
			"clusterAZ":                      utils.ValueIgnoreEmpty(extendParams["cluster_az"]),
			"dssMasterVolumes":               utils.ValueIgnoreEmpty(extendParams["dss_master_volumes"]),
			"alpha.cce/fixPoolMask":          utils.ValueIgnoreEmpty(extendParams["fix_pool_mask"]),
			"decMasterFlavor":                utils.ValueIgnoreEmpty(extendParams["dec_master_flavor"]),
			"dockerUmaskMode":                utils.ValueIgnoreEmpty(extendParams["docker_umask_mode"]),
			"kubernetes.io/cpuManagerPolicy": utils.ValueIgnoreEmpty(extendParams["cpu_manager_policy"]),
		}

		return res
	}

	return nil
}

func buildResourceClusterExtendParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	res := make(map[string]interface{})
	extendParam := resourceClusterExtendParam(d)
	extendParams := resourceClusterExtendParams(d.Get("extend_params").([]interface{}))

	// defaults to use extend_params
	if len(extendParam) != 0 {
		for k, v := range extendParam {
			res[k] = v
		}
	} else {
		for k, v := range extendParams {
			res[k] = v
		}
	}

	if eip, ok := d.GetOk("eip"); ok {
		res["clusterExternalIP"] = eip.(string)
	}

	epsID := cfg.GetEnterpriseProjectID(d)
	if epsID != "" {
		res["enterpriseProjectId"] = epsID
	}

	// assemble the charge info
	var isPrePaid bool
	var billingMode int
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		isPrePaid = true
	}
	if v, ok := d.GetOk("billing_mode"); ok {
		billingMode = v.(int)
	}
	if isPrePaid || billingMode == 1 {
		res["isAutoRenew"] = "false"
		res["isAutoPay"] = common.GetAutoPay(d)
	}

	if v, ok := d.GetOk("period_unit"); ok {
		res["periodType"] = v.(string)
	}
	if v, ok := d.GetOk("period"); ok {
		res["periodNum"] = v.(int)
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		res["isAutoRenew"] = v.(string)
	}

	return utils.RemoveNil(res)
}

func resourceClusterMasters(d *schema.ResourceData) ([]clusters.MasterSpec, error) {
	if v, ok := d.GetOk("masters"); ok {
		flavorId := d.Get("flavor_id").(string)
		mastersRaw := v.([]interface{})
		if strings.Contains(flavorId, "s1") && len(mastersRaw) != 1 {
			return nil, fmt.Errorf("error creating CCE cluster: "+
				"single-master cluster need 1 az for master node, but got %d", len(mastersRaw))
		}
		if strings.Contains(flavorId, "s2") && len(mastersRaw) != 3 {
			return nil, fmt.Errorf("error creating CCE cluster: "+
				"high-availability cluster need 3 az for master nodes, but got %d", len(mastersRaw))
		}
		masters := make([]clusters.MasterSpec, len(mastersRaw))
		for i, raw := range mastersRaw {
			rawMap := raw.(map[string]interface{})
			masters[i] = clusters.MasterSpec{
				MasterAZ: rawMap["availability_zone"].(string),
			}
		}
		return masters, nil
	}

	return nil, nil
}

func buildContainerNetworkCidrsOpts(cidrs string) []clusters.CidrSpec {
	if cidrs == "" {
		return nil
	}

	cidrList := strings.Split(cidrs, ",")

	res := make([]clusters.CidrSpec, len(cidrList))
	for i, cidr := range cidrList {
		res[i] = clusters.CidrSpec{
			Cidr: cidr,
		}
	}

	return res
}

func buildEniNetworkOpts(eniSubnetID string) *clusters.EniNetworkSpec {
	if eniSubnetID == "" {
		return nil
	}

	subnetIDs := strings.Split(eniSubnetID, ",")
	subnets := make([]clusters.EniSubnetSpec, len(subnetIDs))
	for i, subnetID := range subnetIDs {
		subnets[i] = clusters.EniSubnetSpec{
			SubnetID: subnetID,
		}
	}

	eniNetwork := clusters.EniNetworkSpec{
		Subnets: subnets,
	}

	return &eniNetwork
}

func buildResourceClusterConfigurationsOverride(componentConfigurationsRaw []interface{}) ([]clusters.PackageConfiguration, error) {
	if len(componentConfigurationsRaw) == 0 {
		return nil, nil
	}

	res := make([]clusters.PackageConfiguration, len(componentConfigurationsRaw))
	for i, v := range componentConfigurationsRaw {
		if componentConfiguration, ok := v.(map[string]interface{}); ok {
			res[i] = clusters.PackageConfiguration{
				Name: componentConfiguration["name"].(string),
			}

			if configurations := componentConfiguration["configurations"].(string); configurations != "" {
				err := json.Unmarshal([]byte(configurations), &res[i].Configurations)
				if err != nil {
					err = fmt.Errorf("error unmarshalling configurations of %s: %s", componentConfiguration["name"].(string), err)
					return nil, err
				}
			}
		}
	}

	return res, nil
}

func buildResourceClusterEncryptionConfig(d *schema.ResourceData) *clusters.EncryptionConfig {
	encryptionConfigRaw, ok := d.GetOk("encryption_config")
	if !ok {
		return nil
	}

	encryptionConfig := encryptionConfigRaw.([]interface{})[0]

	res := clusters.EncryptionConfig{
		Mode:     utils.PathSearch("mode", encryptionConfig, "").(string),
		KmsKeyID: utils.PathSearch("kms_key_id", encryptionConfig, "").(string),
	}

	return &res
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	cceClient, err := config.CceV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}
	icAgentClient, err := config.AomV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM v1 client: %s", err)
	}

	authenticating_proxy := make(map[string]string)
	if common.HasFilledOpt(d, "authenticating_proxy_ca") {
		authenticating_proxy["ca"] = utils.TryBase64EncodeString(d.Get("authenticating_proxy_ca").(string))
		authenticating_proxy["cert"] = utils.TryBase64EncodeString(d.Get("authenticating_proxy_cert").(string))
		authenticating_proxy["privateKey"] = utils.TryBase64EncodeString(d.Get("authenticating_proxy_private_key").(string))
	}

	billingMode := 0
	if d.Get("charging_mode").(string) == "prePaid" || d.Get("billing_mode").(int) == 1 {
		billingMode = 1
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return diag.FromErr(err)
		}
	}

	clusterName := d.Get("name").(string)
	createOpts := clusters.CreateOpts{
		Kind:       "Cluster",
		ApiVersion: "v3",
		Metadata: clusters.CreateMetaData{
			Name:        clusterName,
			Alias:       d.Get("alias").(string),
			Labels:      resourceClusterLabels(d),
			Annotations: resourceClusterAnnotations(d),
			Timezone:    d.Get("timezone").(string),
		},

		Spec: clusters.Spec{
			Type:        d.Get("cluster_type").(string),
			Flavor:      d.Get("flavor_id").(string),
			Version:     d.Get("cluster_version").(string),
			Description: d.Get("description").(string),
			HostNetwork: clusters.HostNetworkSpec{
				VpcId:         d.Get("vpc_id").(string),
				SubnetId:      d.Get("subnet_id").(string),
				HighwaySubnet: d.Get("highway_subnet_id").(string),
				SecurityGroup: d.Get("security_group_id").(string),
			},
			ContainerNetwork: clusters.ContainerNetworkSpec{
				Mode:  d.Get("container_network_type").(string),
				Cidrs: buildContainerNetworkCidrsOpts(d.Get("container_network_cidr").(string)),
			},
			EniNetwork: buildEniNetworkOpts(d.Get("eni_subnet_id").(string)),
			Authentication: clusters.AuthenticationSpec{
				Mode:                d.Get("authentication_mode").(string),
				AuthenticatingProxy: authenticating_proxy,
			},
			BillingMode:      billingMode,
			ExtendParam:      buildResourceClusterExtendParams(d, config),
			ClusterTags:      resourceClusterTags(d),
			CustomSan:        utils.ExpandToStringList(d.Get("custom_san").([]interface{})),
			IPv6Enable:       d.Get("ipv6_enable").(bool),
			KubeProxyMode:    d.Get("kube_proxy_mode").(string),
			EncryptionConfig: buildResourceClusterEncryptionConfig(d),
		},
	}

	if _, ok := d.GetOk("enable_distribute_management"); ok {
		createOpts.Spec.EnableDistMgt = d.Get("enable_distribute_management").(bool)
	}

	if v, ok := d.GetOk("service_network_cidr"); ok {
		serviceNetwork := clusters.ServiceNetwork{
			IPv4Cidr: v.(string),
		}
		createOpts.Spec.ServiceNetwork = &serviceNetwork
	}

	masters, err := resourceClusterMasters(d)
	if err != nil {
		return diag.FromErr(err)
	}
	createOpts.Spec.Masters = masters

	componentConfigurations, err := buildResourceClusterConfigurationsOverride(
		d.Get("component_configurations").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createOpts.Spec.ConfigurationsOverride = componentConfigurations

	s, err := clusters.Create(cceClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating CCE cluster: %s", err)
	}

	if orderId, ok := s.Spec.ExtendParam["orderID"]; ok && orderId != "" {
		bssClient, err := config.BssV2Client(config.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resourceId)
	} else {
		jobID := s.Status.JobID
		if jobID == "" {
			return diag.Errorf("error fetching job ID after creating CCE cluster: %s", clusterName)
		}

		clusterID, err := getClusterIDFromJob(ctx, cceClient, jobID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(clusterID)
	}

	log.Printf("[DEBUG] Waiting for CCE cluster (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		// The statuses of pending phase include "Creating".
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      clusterStateRefreshFunc(cceClient, d.Id(), []string{"Available"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating CCE cluster: %s", err)
	}

	log.Printf("[DEBUG] Installing ICAgent for CCE cluster (%s)", d.Id())
	installParam := icagents.InstallParam{
		ClusterId: d.Id(),
		NameSpace: "default",
	}
	result := icagents.Create(icAgentClient, installParam)
	var diags diag.Diagnostics
	if result.Err != nil {
		diagIcagent := diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error installing ICAgent",
			Detail:   fmt.Sprintf("error installing ICAgent in CCE cluster: %s", result.Err),
		}
		diags = append(diags, diagIcagent)
	}

	// create a hibernating cluster
	if d.Get("hibernate").(bool) {
		err = resourceClusterHibernate(ctx, d, cceClient)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	diags = append(diags, resourceClusterRead(ctx, d, meta)...)

	return diags

}

func resourceClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	cceClient, err := config.CceV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	n, err := clusters.Get(cceClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CCE cluster")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", n.Metadata.Name),
		d.Set("alias", n.Metadata.Alias),
		d.Set("timezone", n.Metadata.Timezone),
		d.Set("status", n.Status.Phase),
		d.Set("flavor_id", n.Spec.Flavor),
		d.Set("cluster_version", n.Spec.Version),
		d.Set("cluster_type", n.Spec.Type),
		d.Set("description", n.Spec.Description),
		d.Set("vpc_id", n.Spec.HostNetwork.VpcId),
		d.Set("subnet_id", n.Spec.HostNetwork.SubnetId),
		d.Set("highway_subnet_id", n.Spec.HostNetwork.HighwaySubnet),
		d.Set("container_network_type", n.Spec.ContainerNetwork.Mode),
		d.Set("container_network_cidr", flattenContainerNetworkCidrs(n.Spec.ContainerNetwork)),
		d.Set("eni_subnet_id", flattenEniSubnetID(n.Spec.EniNetwork)),
		d.Set("eni_subnet_cidr", n.Spec.EniNetwork.Cidr),
		d.Set("authentication_mode", n.Spec.Authentication.Mode),
		d.Set("security_group_id", n.Spec.HostNetwork.SecurityGroup),
		d.Set("enterprise_project_id", n.Spec.ExtendParam["enterpriseProjectId"]),
		d.Set("service_network_cidr", n.Spec.ServiceNetwork.IPv4Cidr),
		d.Set("billing_mode", n.Spec.BillingMode),
		d.Set("tags", utils.TagsToMap(n.Spec.ClusterTags)),
		d.Set("ipv6_enable", n.Spec.IPv6Enable),
		d.Set("enable_distribute_management", n.Spec.EnableDistMgt),
		d.Set("kube_proxy_mode", n.Spec.KubeProxyMode),
		d.Set("support_istio", n.Spec.SupportIstio),
		d.Set("custom_san", n.Spec.CustomSan),
		d.Set("category", n.Spec.Category),
		d.Set("encryption_config", flattenEncrytionConfig(n.Spec.EncryptionConfig)),
	)

	if n.Spec.BillingMode != 0 {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "prePaid"))
	}

	// duration -1 is equal to the maximum value 1827 days
	opts := clusters.GetCertOpts{Duration: -1}
	r := clusters.GetCert(cceClient, d.Id(), opts)

	kubeConfigRaw, err := utils.JsonMarshal(r.Body)

	if err != nil {
		log.Printf("error marshaling r.Body: %s", err)
	}

	mErr = multierror.Append(mErr, d.Set("kube_config_raw", string(kubeConfigRaw)))

	cert, err := r.Extract()

	if err != nil {
		log.Printf("error retrieving CCE cluster certificate: %s", err)
	}

	//Set Certificate Clusters
	var clusterList []map[string]interface{}
	for _, clusterObj := range cert.Clusters {
		clusterCert := make(map[string]interface{})
		clusterCert["name"] = clusterObj.Name
		clusterCert["server"] = clusterObj.Cluster.Server
		clusterCert["certificate_authority_data"] = clusterObj.Cluster.CertAuthorityData
		clusterList = append(clusterList, clusterCert)
	}
	mErr = multierror.Append(mErr, d.Set("certificate_clusters", clusterList))

	//Set Certificate Users
	var userList []map[string]interface{}
	for _, userObj := range cert.Users {
		userCert := make(map[string]interface{})
		userCert["name"] = userObj.Name
		userCert["client_certificate_data"] = userObj.User.ClientCertData
		userCert["client_key_data"] = userObj.User.ClientKeyData
		userList = append(userList, userCert)
	}
	mErr = multierror.Append(mErr, d.Set("certificate_users", userList))

	// Set masters
	var masterList []map[string]interface{}
	for _, masterObj := range n.Spec.Masters {
		master := make(map[string]interface{})
		master["availability_zone"] = masterObj.MasterAZ
		masterList = append(masterList, master)
	}
	mErr = multierror.Append(mErr, d.Set("masters", masterList))

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CCE cluster fields: %s", err)
	}

	return nil
}

func flattenContainerNetworkCidrs(containerNetwork clusters.ContainerNetworkSpec) string {
	cidrs := containerNetwork.Cidrs
	if len(cidrs) != 0 {
		cidrList := make([]string, len(cidrs))
		for i, v := range cidrs {
			cidrList[i] = v.Cidr
		}

		return strings.Join(cidrList, ",")
	}

	return containerNetwork.Cidr
}

func flattenEniSubnetID(eniNetwork *clusters.EniNetworkSpec) string {
	if eniNetwork == nil {
		return ""
	}

	subnets := eniNetwork.Subnets
	subnetIDs := make([]string, len(subnets))
	for i, v := range subnets {
		subnetIDs[i] = v.SubnetID
	}

	return strings.Join(subnetIDs, ",")
}

func flattenEncrytionConfig(encrytionConfig *clusters.EncryptionConfig) []map[string]interface{} {
	if encrytionConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"mode":       encrytionConfig.Mode,
			"kms_key_id": encrytionConfig.KmsKeyID,
		},
	}
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	cceClient, err := cfg.CceV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	clusterId := d.Id()
	updateOpts := clusters.UpdateOpts{}

	if d.HasChange("alias") {
		updateOpts.Metadata = &clusters.UpdateMetadata{
			Alias: d.Get("alias").(string),
		}
	}

	if d.HasChanges("description") {
		updateOpts.Spec.Description = d.Get("description").(string)
	}

	if d.HasChange("container_network_cidr") {
		o, n := d.GetChange("container_network_cidr")
		oldCidr := o.(string)
		newCidr := n.(string)

		if len(newCidr) < len(oldCidr) || newCidr[0:len(oldCidr)] != oldCidr {
			return diag.Errorf("error updating CCE cluster: " +
				"the container_network_cidr can only be updated incrementally," +
				" and the new value must contains the old value as a prefix")
		}

		// only incremental part can be contained in the request
		updateOpts.Spec.ContainerNetwork = &clusters.UpdateContainerNetworkSpec{
			Cidrs: buildContainerNetworkCidrsOpts(newCidr[len(oldCidr)+1:]),
		}
	}

	if d.HasChanges("eni_subnet_id") {
		updateOpts.Spec.EniNetwork = buildEniNetworkOpts(d.Get("eni_subnet_id").(string))
	}

	if d.HasChange("security_group_id") {
		updateOpts.Spec.HostNetwork = &clusters.UpdateHostNetworkSpec{
			SecurityGroup: d.Get("security_group_id").(string),
		}
	}

	if d.HasChange("custom_san") {
		updateOpts.Spec.CustomSan = utils.ExpandToStringList(d.Get("custom_san").([]interface{}))
	}

	if !reflect.DeepEqual(updateOpts, clusters.UpdateOpts{}) {
		_, err = clusters.Update(cceClient, clusterId, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating CCE cluster: %s", err)
		}
	}

	if d.HasChange("flavor_id") {
		err := resourceClusterResize(ctx, cfg, d, cceClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("hibernate") {
		if d.Get("hibernate").(bool) {
			err = resourceClusterHibernate(ctx, d, cceClient)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err = resourceClusterAwake(ctx, d, cceClient)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("eip") {
		eipClient, err := cfg.NetworkingV1Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v1 client: %s", err)
		}

		oldEip, newEip := d.GetChange("eip")
		if oldEip.(string) != "" {
			err = resourceClusterEipAction(cceClient, eipClient, clusterId, oldEip.(string), "unbind")
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if newEip.(string) != "" {
			err = resourceClusterEipAction(cceClient, eipClient, clusterId, newEip.(string), "bind")
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("tags") {
		// remove old tags and set new tags
		oldTags, newTags := d.GetChange("tags")
		oldTagsRaw := oldTags.(map[string]interface{})
		if len(oldTagsRaw) > 0 {
			taglist := utils.ExpandResourceTags(oldTagsRaw)
			if tagErr := clusters.RemoveTags(cceClient, clusterId, taglist).ExtractErr(); tagErr != nil {
				return diag.Errorf("error deleting tags of CCE cluster %s: %s", clusterId, tagErr)
			}
		}

		newTagsRaw := newTags.(map[string]interface{})
		if len(newTagsRaw) > 0 {
			taglist := utils.ExpandResourceTags(newTagsRaw)
			if tagErr := clusters.AddTags(cceClient, clusterId, taglist).ExtractErr(); tagErr != nil {
				return diag.Errorf("error setting tags of CCE cluster %s: %s", clusterId, tagErr)
			}
		}
	}

	if d.HasChange("component_configurations") {
		var (
			updateClusterConfigurationsHttpUrl = "api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/master/configuration"
		)

		updateClusterConfigurationsPath := cceClient.Endpoint + updateClusterConfigurationsHttpUrl
		updateClusterConfigurationsPath = strings.ReplaceAll(updateClusterConfigurationsPath, "{project_id}", cfg.GetProjectID(region))
		updateClusterConfigurationsPath = strings.ReplaceAll(updateClusterConfigurationsPath, "{cluster_id}", clusterId)

		updateClusterConfigurationsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		bodyParams, err := buildUpdateClusterConfigurationsBodyParams(d)
		if err != nil {
			return diag.FromErr(err)
		}
		updateClusterConfigurationsOpt.JSONBody = bodyParams
		_, err = cceClient.Request("PUT", updateClusterConfigurationsPath, &updateClusterConfigurationsOpt)
		if err != nil {
			return diag.Errorf("error updating CCE cluster configurations: %s", err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), clusterId); err != nil {
			return diag.Errorf("error updating the auto-renew of the CCE cluster (%s): %s", clusterId, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   clusterId,
			ResourceType: "cce-cluster",
			RegionId:     region,
			ProjectId:    cceClient.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceClusterRead(ctx, d, meta)
}

func buildUpdateClusterConfigurationsBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	configurationsPackages, err := buildConfigurationsPackagesBodyParams(d)
	if err != nil {
		return nil, err
	}
	bodyParams := map[string]interface{}{
		"apiVersion": "v3",
		"kind":       "Configuration",
		"metadata": map[string]interface{}{
			"name": "configuration",
		},
		"spec": map[string]interface{}{
			"packages": configurationsPackages,
		},
	}
	return bodyParams, nil
}

func buildConfigurationsPackagesBodyParams(d *schema.ResourceData) ([]map[string]interface{}, error) {
	packagesRaw := d.Get("component_configurations").([]interface{})
	bodyParams := make([]map[string]interface{}, len(packagesRaw))
	for i, v := range packagesRaw {
		packageRaw := v.(map[string]interface{})
		bodyParams[i] = map[string]interface{}{
			"name": packageRaw["name"],
		}

		if configurationsRaw := packageRaw["configurations"].(string); configurationsRaw != "" {
			var configurations interface{}
			err := json.Unmarshal([]byte(configurationsRaw), &configurations)
			if err != nil {
				err = fmt.Errorf("error unmarshalling configurations of %s: %s", packageRaw["name"].(string), err)
				return nil, err
			}
			bodyParams[i]["configurations"] = configurations
		}
	}

	return bodyParams, nil
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	cceClient, err := config.CceV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	// for prePaid mode, we should unsubscribe the resource
	if d.Get("charging_mode").(string) == "prePaid" || d.Get("billing_mode").(int) == 1 {
		if err := common.UnsubscribePrePaidResource(d, config, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribing CCE cluster: %s", err)
		}
	} else {
		deleteOpts := clusters.DeleteOpts{}
		if v, ok := d.GetOk("delete_all"); ok && v.(string) != "false" {
			deleteOpt := d.Get("delete_all").(string)
			deleteOpts.DeleteEfs = deleteOpt
			deleteOpts.DeleteEvs = deleteOpt
			deleteOpts.DeleteObs = deleteOpt
			deleteOpts.DeleteSfs = deleteOpt
			deleteOpts.DeleteSfs30 = deleteOpt
		} else {
			deleteOpts.DeleteEfs = d.Get("delete_efs").(string)
			deleteOpts.DeleteENI = d.Get("delete_eni").(string)
			deleteOpts.DeleteEvs = d.Get("delete_evs").(string)
			deleteOpts.DeleteNet = d.Get("delete_net").(string)
			deleteOpts.DeleteObs = d.Get("delete_obs").(string)
			// delete_sfs indecates delete SFS and SFS3.0 together
			deleteOpts.DeleteSfs = d.Get("delete_sfs").(string)
			deleteOpts.DeleteSfs30 = d.Get("delete_sfs").(string)
		}

		deleteOpts.LtsReclaimPolicy = d.Get("lts_reclaim_policy").(string)
		err = clusters.DeleteWithOpts(cceClient, d.Id(), deleteOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error deleting CCE cluster: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		// The statuses of pending phase includes "Deleting", "Available" and "Unavailable".
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      clusterStateRefreshFunc(cceClient, d.Id(), nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)

	if err != nil {
		return diag.Errorf("error deleting CCE cluster: %s", err)
	}

	d.SetId("")
	return nil
}

func clusterStateRefreshFunc(cceClient *golangsdk.ServiceClient, clusterId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Expect the status of CCE cluster to be any one of the status list: %v", targets)
		resp, err := clusters.Get(cceClient, clusterId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] The cluster (%s) has been deleted", clusterId)
				return resp, "COMPLETED", nil
			}
			return nil, "ERROR", err
		}

		invalidStatuses := []string{"Error", "Shelved", "Unknow"}
		if utils.IsStrContainsSliceElement(resp.Status.Phase, invalidStatuses, true, true) {
			return resp, "ERROR", fmt.Errorf("unexpected status: %s", resp.Status.Phase)
		}

		if utils.StrSliceContains(targets, resp.Status.Phase) {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}

func getClusterIDFromJob(ctx context.Context, client *golangsdk.ServiceClient, jobID string, timeout time.Duration) (string, error) {
	stateJob := &resource.StateChangeConf{
		Pending:      []string{"Initializing", "Running"},
		Target:       []string{"Success"},
		Refresh:      waitForJobStatus(client, jobID),
		Timeout:      timeout,
		Delay:        150 * time.Second,
		PollInterval: 20 * time.Second,
	}

	v, err := stateJob.WaitForStateContext(ctx)
	if err != nil {
		if job, ok := v.(*nodes.Job); ok {
			return "", fmt.Errorf("error waiting for job (%s) to become success: %s, reason: %s",
				jobID, err, job.Status.Reason)
		} else {
			return "", fmt.Errorf("error waiting for job (%s) to become success: %s", jobID, err)
		}

	}

	job := v.(*nodes.Job)
	clusterID := job.Spec.ClusterID
	if clusterID == "" {
		return "", fmt.Errorf("error fetching CCE cluster ID")
	}
	return clusterID, nil
}

func resourceClusterResize(ctx context.Context, cfg *config.Config, d *schema.ResourceData, cceClient *golangsdk.ServiceClient) error {
	clusterID := d.Id()

	var decMasterFlavor string
	extendParams := resourceClusterExtendParams(d.Get("extend_params").([]interface{}))
	if v, ok := extendParams["decMasterFlavor"]; ok {
		decMasterFlavor = v.(string)
	}

	opts := clusters.ResizeOpts{
		FavorResize: d.Get("flavor_id").(string),
		ExtendParam: &clusters.ResizeExtendParam{
			DecMasterFlavor: decMasterFlavor,
		},
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		opts.ExtendParam.IsAutoPay = common.GetAutoPay(d)
	}

	resp, err := clusters.Resize(cceClient, clusterID, opts)
	if err != nil {
		return fmt.Errorf("error resizing CCE cluster: %s", err)
	}

	if resp.OrderID != "" {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, resp.OrderID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Waiting for CCE cluster (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      clusterStateRefreshFunc(cceClient, d.Id(), []string{"Available"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error resizing CCE cluster: %s", err)
	}
	return nil
}

func resourceClusterHibernate(ctx context.Context, d *schema.ResourceData, cceClient *golangsdk.ServiceClient) error {
	clusterID := d.Id()
	err := clusters.Operation(cceClient, clusterID, "hibernate").ExtractErr()
	if err != nil {
		return fmt.Errorf("error hibernating CCE cluster: %s", err)
	}

	log.Printf("[DEBUG] Waiting for CCE cluster (%s) to become hibernate", clusterID)
	stateConf := &resource.StateChangeConf{
		// The statuses of pending phase includes "Available" and "Hibernating".
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      clusterStateRefreshFunc(cceClient, clusterID, []string{"Hibernation"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error hibernating CCE cluster: %s", err)
	}
	return nil
}

func resourceClusterAwake(ctx context.Context, d *schema.ResourceData, cceClient *golangsdk.ServiceClient) error {
	clusterID := d.Id()
	err := clusters.Operation(cceClient, clusterID, "awake").ExtractErr()
	if err != nil {
		return fmt.Errorf("error awaking CCE cluster: %s", err)
	}

	log.Printf("[DEBUG] Waiting for CCE cluster (%s) to become available", clusterID)
	stateConf := &resource.StateChangeConf{
		// The statuses of pending phase include "Awaking".
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      clusterStateRefreshFunc(cceClient, clusterID, []string{"Available"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        100 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error awaking CCE cluster: %s", err)
	}
	return nil
}

func resourceClusterEipAction(cceClient, eipClient *golangsdk.ServiceClient,
	clusterID, eip, action string) error {
	eipID, err := common.GetEipIDbyAddress(eipClient, eip, "all_granted_eps")
	if err != nil {
		return fmt.Errorf("error fetching EIP ID: %s", err)
	}

	opts := clusters.UpdateIpOpts{
		Action: action,
		Spec: clusters.IpSpec{
			ID: eipID,
		},
	}

	err = clusters.UpdateMasterIp(cceClient, clusterID, opts).ExtractErr()
	if err != nil {
		return fmt.Errorf("error %sing the public IP of CCE cluster: %s", action, err)
	}
	return nil
}
