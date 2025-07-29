// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DWS
// ---------------------------------------------------------------

package dws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	PublicBindTypeAuto         = "auto_assign"
	PublicBindTypeNotUse       = "not_use"
	PublicBindTypeBindExisting = "bind_existing"
)

const ClusterIdIllegalErrCode = "DWS.0001"

// @API DWS POST /v1.0/{project_id}/clusters
// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}
// @API DWS POST /v1.0/{project_id}/clusters/{cluster_id}/expand-instance-storage
// @API DWS POST /v1.0/{project_id}/clusters/{cluster_id}/reset-password
// @API DWS POST /v1.0/{project_id}/clusters/{cluster_id}/resize
// @API DWS DELETE /v1.0/{project_id}/clusters/{cluster_id}
// @API DWS POST /v1.0/{project_id}/clusters/{cluster_id}/tags/batch-create
// @API DWS POST /v1.0/{project_id}/clusters/{cluster_id}/tags/batch-delete
// @API DWS GET /v1.0/{project_id}/job/{job_id}
// @API DWS POST /v2/{project_id}/clusters
// @API DWS PUT /v2/{project_id}/clusters/{cluster_id}/logical-clusters/enable
// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/elbs/{elb_id}
// @API DWS DELETE /v2/{project_id}/clusters/{cluster_id}/elbs/{elb_id}
// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/lts-logs/enable
// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/lts-logs/disable
// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/eips/{eip_id}
// @API DWS DELETE /v2/{project_id}/clusters/{cluster_id}/eips/{eip_id}
// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/description
// @API DWS PUT /v1/{project_id}/clusters/{cluster_id}/security-group
// @API DWS POST /v1.0/{project_id}/clusters/{cluster_id}/cluster-shrink
func ResourceDwsCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDwsClusterCreate,
		UpdateContext: resourceDwsClusterUpdate,
		ReadContext:   resourceDwsClusterRead,
		DeleteContext: resourceDwsClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The cluster name.`,
			},
			"node_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The flavor of the cluster.`,
			},
			"number_of_node": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Number of nodes in a cluster.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Administrator username for logging in to a data warehouse cluster.`,
			},
			"user_pwd": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Administrator password for logging in to a data warehouse cluster.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The VPC ID.`,
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The subnet ID.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The security group ID.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The availability zone in which to create the cluster instance. `,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "schema: Required",
			},
			"number_of_cn": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "schema: Required",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     8000,
				Description: "Service port of a cluster (8000 to 10000). The default value is 8000.",
			},
			"dss_pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Dedicated storage pool ID.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The enterprise project ID.`,
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The KMS key ID.`,
			},
			"public_ip": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     clusterPublicIpSchema(),
				Optional: true,
				Computed: true,
			},
			"volume": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     clusterVolumeSchema(),
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
			"keep_last_manual_snapshot": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The number of latest manual snapshots that need to be retained when deleting the cluster.`,
			},
			"logical_cluster_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable logical cluster.`,
			},
			"elb_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the ELB load balancer.`,
			},
			"lts_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable LTS.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the cluster.`,
			},
			"force_backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to automatically execute snapshot when shrinking the number of nodes.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cluster status.`,
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the cluster.`,
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updated time of the cluster.`,
			},
			"endpoints": {
				Type:        schema.TypeList,
				Elem:        clusterEndpointSchema(),
				Computed:    true,
				Description: `Private network connection information about the cluster.`,
			},
			"public_endpoints": {
				Type:        schema.TypeList,
				Elem:        clusterPublicEndpointSchema(),
				Computed:    true,
				Description: `Public network connection information about the cluster.`,
			},
			"recent_event": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sub_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Sub-status of clusters in the AVAILABLE state.`,
			},
			"task_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster management task.`,
			},
			"private_ip": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `List of private network IP addresses.`,
			},
			"maintain_window": {
				Type:        schema.TypeList,
				Elem:        clusterMaintainWindowSchema(),
				Computed:    true,
				Description: `Cluster maintenance window.`,
			},
			"elb": {
				Type:        schema.TypeList,
				Elem:        clusterElbSchema(),
				Computed:    true,
				Description: `The ELB information bound to the cluster.`,
			},
		},
	}
}

func clusterPublicIpSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"public_bind_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The bind type of public IP.`,
			},
			"eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The EIP ID.`,
				DiffSuppressFunc: func(_, _, newVal string, d *schema.ResourceData) bool {
					// If "public_bind_type" is set to "auto_assign", the EIP will be automatically bound, the EIP Will be triggered to change.
					if v, ok := d.GetOk("public_ip.0.public_bind_type"); ok {
						return v.(string) == PublicBindTypeAuto && newVal == ""
					}
					return false
				},
			},
		},
	}
	return &sc
}

func clusterVolumeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The volume type.`,
			},
			"capacity": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The capacity size, in GB.`,
			},
		},
	}
	return &sc
}

func clusterEndpointSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"connect_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Private network connection information.`,
			},
			"jdbc_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `JDBC URL. Format: jdbc:postgresql://<connect_info>/<YOUR_DATABASE_NAME>`,
			},
		},
	}
	return &sc
}

func clusterPublicEndpointSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"public_connect_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Public network connection information.`,
			},
			"jdbc_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `JDBC URL. Format: jdbc:postgresql://<public_connect_info>/<YOUR_DATABASE_NAME>`,
			},
		},
	}
	return &sc
}

func clusterMaintainWindowSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"day": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Maintenance time in each week in the unit of day.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Maintenance start time in HH:mm format. The time zone is GMT+0.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Maintenance end time in HH:mm format. The time zone is GMT+0.`,
			},
		},
	}
	return &sc
}

func clusterElbSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the ELB load balancer.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the ELB load balancer.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public IP address of the ELB load balancer.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The private IP address of the ELB load balancer.`,
			},
			"private_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The private endpoint of the ELB load balancer.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of VPC to which the ELB load balancer belongs.`,
			},
			"private_ip_v6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IPv6 address of the ELB load balancer.`,
			},
		},
	}
	return &sc
}

func resourceDwsClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if _, ok := d.GetOk("version"); ok {
		return resourceDwsClusterCreateV2(ctx, d, meta)
	}

	return resourceDwsClusterCreateV1(ctx, d, meta)
}

func updateDwsLogicalClusterEnable(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	switchDwsClusterPath := client.Endpoint + "v2/{project_id}/clusters/{cluster_id}/logical-clusters/enable"
	switchDwsClusterPath = strings.ReplaceAll(switchDwsClusterPath, "{project_id}", client.ProjectID)
	switchDwsClusterPath = strings.ReplaceAll(switchDwsClusterPath, "{cluster_id}", d.Id())

	switchDwsClusterOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"enable": d.Get("logical_cluster_enable"),
		},
	}
	_, err := client.Request("PUT", switchDwsClusterPath, &switchDwsClusterOpt)
	if err != nil {
		return fmt.Errorf("error updating DWS logical cluster switch: %s", err)
	}
	return nil
}

func resourceDwsClusterCreateV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDwsCluster: create a DWS cluster.
	var (
		createDwsClusterHttpUrl = "v2/{project_id}/clusters"
		createDwsClusterProduct = "dws"
	)
	createDwsClusterClient, err := cfg.NewServiceClient(createDwsClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createDwsClusterPath := createDwsClusterClient.Endpoint + createDwsClusterHttpUrl
	createDwsClusterPath = strings.ReplaceAll(createDwsClusterPath, "{project_id}", createDwsClusterClient.ProjectID)

	createDwsClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	createDwsClusterOpt.JSONBody = utils.RemoveNil(buildCreateDwsClusterBodyParams(d, cfg))
	createDwsClusterResp, err := createDwsClusterClient.Request("POST", createDwsClusterPath, &createDwsClusterOpt)
	if err != nil {
		return diag.Errorf("error creating DWS cluster: %s", err)
	}

	createDwsClusterRespBody, err := utils.FlattenResponse(createDwsClusterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := utils.PathSearch("cluster.id", createDwsClusterRespBody, "").(string)
	if clusterId == "" {
		return diag.Errorf("unable to find the DWS Cluster ID from the API response")
	}
	d.SetId(clusterId)

	err = clusterWaitingForAvailable(ctx, d, createDwsClusterClient, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the DWS cluster (%s) creation to complete: %s", clusterId, err)
	}

	if d.Get("logical_cluster_enable").(bool) {
		if err := updateDwsLogicalClusterEnable(createDwsClusterClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// The cluster binding ELB load balaner.
	if v, ok := d.GetOk("elb_id"); ok {
		elbId := v.(string)
		err := bindElb(ctx, d, createDwsClusterClient, elbId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// If lts_enable is true, enable LTS.
	if d.Get("lts_enable").(bool) {
		err = enableOrDisableLts(d, createDwsClusterClient)
		if err != nil {
			return diag.Errorf("error enable LTS for DWS cluster: %s", err)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		if err := updateDescription(createDwsClusterClient, clusterId, v.(string)); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceDwsClusterRead(ctx, d, meta)
}

func resourceDwsClusterCreateV1(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDwsCluster: create a DWS cluster.
	var (
		createDwsClusterHttpUrl = "v1.0/{project_id}/clusters"
		createDwsClusterProduct = "dws"
	)
	createDwsClusterClient, err := cfg.NewServiceClient(createDwsClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createDwsClusterPath := createDwsClusterClient.Endpoint + createDwsClusterHttpUrl
	createDwsClusterPath = strings.ReplaceAll(createDwsClusterPath, "{project_id}", createDwsClusterClient.ProjectID)

	createDwsClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	createDwsClusterOpt.JSONBody = utils.RemoveNil(buildCreateDwsClusterBodyParamsV1(d, cfg))
	createDwsClusterResp, err := createDwsClusterClient.Request("POST", createDwsClusterPath, &createDwsClusterOpt)
	if err != nil {
		return diag.Errorf("error creating DWS cluster: %s", err)
	}

	createDwsClusterRespBody, err := utils.FlattenResponse(createDwsClusterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := utils.PathSearch("cluster.id", createDwsClusterRespBody, "").(string)
	if clusterId == "" {
		return diag.Errorf("unable to find the DWS Cluster ID from the API response")
	}
	d.SetId(clusterId)

	err = clusterWaitingForAvailable(ctx, d, createDwsClusterClient, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the DWS cluster (%s) creation to complete: %s", clusterId, err)
	}

	if d.Get("logical_cluster_enable").(bool) {
		if err := updateDwsLogicalClusterEnable(createDwsClusterClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// The cluster binding ELB load balaner.
	if v, ok := d.GetOk("elb_id"); ok {
		elbId := v.(string)
		err := bindElb(ctx, d, createDwsClusterClient, elbId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		if err := updateDescription(createDwsClusterClient, clusterId, v.(string)); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceDwsClusterRead(ctx, d, meta)
}

func buildCreateDwsClusterBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	availabilityZones := strings.Split(d.Get("availability_zone").(string), ",")
	bodyParams := map[string]interface{}{
		"cluster": map[string]interface{}{
			"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
			"flavor":                utils.ValueIgnoreEmpty(d.Get("node_type")),
			"num_node":              utils.ValueIgnoreEmpty(d.Get("number_of_node")),
			"num_cn":                utils.ValueIgnoreEmpty(d.Get("number_of_cn")),
			"db_name":               utils.ValueIgnoreEmpty(d.Get("user_name")),
			"db_password":           utils.ValueIgnoreEmpty(d.Get("user_pwd")),
			"db_port":               utils.ValueIgnoreEmpty(d.Get("port")),
			"availability_zones":    availabilityZones,
			"vpc_id":                utils.ValueIgnoreEmpty(d.Get("vpc_id")),
			"subnet_id":             utils.ValueIgnoreEmpty(d.Get("network_id")),
			"security_group_id":     utils.ValueIgnoreEmpty(d.Get("security_group_id")),
			"datastore_version":     utils.ValueIgnoreEmpty(d.Get("version")),
			"dss_pool_id":           utils.ValueIgnoreEmpty(d.Get("dss_pool_id")),
			"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
			"master_key_id":         utils.ValueIgnoreEmpty(d.Get("kms_key_id")),
			"public_ip":             buildCreateDwsClusterReqBodyPublicIp(d.Get("public_ip")),
			"volume":                buildCreateDwsClusterReqBodyVolume(d.Get("volume")),
			"tags":                  utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		},
	}
	return bodyParams
}

func buildCreateDwsClusterBodyParamsV1(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cluster": map[string]interface{}{
			"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
			"node_type":             utils.ValueIgnoreEmpty(d.Get("node_type")),
			"number_of_node":        utils.ValueIgnoreEmpty(d.Get("number_of_node")),
			"number_of_cn":          utils.ValueIgnoreEmpty(d.Get("number_of_cn")),
			"user_name":             utils.ValueIgnoreEmpty(d.Get("user_name")),
			"user_pwd":              utils.ValueIgnoreEmpty(d.Get("user_pwd")),
			"port":                  utils.ValueIgnoreEmpty(d.Get("port")),
			"availability_zone":     utils.ValueIgnoreEmpty(d.Get("availability_zone")),
			"vpc_id":                utils.ValueIgnoreEmpty(d.Get("vpc_id")),
			"subnet_id":             utils.ValueIgnoreEmpty(d.Get("network_id")),
			"security_group_id":     utils.ValueIgnoreEmpty(d.Get("security_group_id")),
			"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
			"public_ip":             buildCreateDwsClusterReqBodyPublicIp(d.Get("public_ip")),
			"tags":                  utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		},
	}
	return bodyParams
}

func buildCreateDwsClusterReqBodyPublicIp(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"public_bind_type": utils.ValueIgnoreEmpty(raw["public_bind_type"]),
			"eip_id":           utils.ValueIgnoreEmpty(raw["eip_id"]),
		}
		return params
	}
	return nil
}

func buildCreateDwsClusterReqBodyVolume(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"volume":   utils.ValueIgnoreEmpty(raw["type"]),
			"capacity": utils.ValueIgnoreEmpty(raw["capacity"]),
		}
		return params
	}
	return nil
}

func clusterWaitingForAvailable(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			clusterWaitingRespBody, err := GetClusterInfoByClusterId(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			actionProgressRaw := utils.PathSearch(`length(cluster.action_progress)`, clusterWaitingRespBody, 0.0)
			if actionProgressRaw.(float64) > 0 {
				return clusterWaitingRespBody, "PENDING", nil
			}

			status := utils.PathSearch(`cluster.status`, clusterWaitingRespBody, "").(string)

			targetStatus := []string{
				"AVAILABLE",
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return clusterWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"FAILED",
				"CREATE_FAILED",
				"CREATION FAILED",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return clusterWaitingRespBody, status, nil
			}

			return clusterWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// GetClusterInfoByClusterId is a method that used to query DWS cluster detail.
func GetClusterInfoByClusterId(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	getDwsClusterHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}"
	getDwsClusterPath := client.Endpoint + getDwsClusterHttpUrl
	getDwsClusterPath = strings.ReplaceAll(getDwsClusterPath, "{project_id}", client.ProjectID)
	getDwsClusterPath = strings.ReplaceAll(getDwsClusterPath, "{cluster_id}", clusterId)

	getDwsClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	getDwsClusterResp, err := client.Request("GET", getDwsClusterPath, &getDwsClusterOpt)
	if err != nil {
		return nil, parseClusterNotFoundError(err)
	}
	return utils.FlattenResponse(getDwsClusterResp)
}

func resourceDwsClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDwsCluster: Query the DWS cluster.
	getDwsClusterClient, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	getDwsClusterRespBody, err := GetClusterInfoByClusterId(getDwsClusterClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DWS cluster")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("cluster.name", getDwsClusterRespBody, nil)),
		d.Set("status", utils.PathSearch("cluster.status", getDwsClusterRespBody, nil)),
		d.Set("version", utils.PathSearch("cluster.version", getDwsClusterRespBody, nil)),
		d.Set("created", utils.PathSearch("cluster.created", getDwsClusterRespBody, nil)),
		d.Set("updated", utils.PathSearch("cluster.updated", getDwsClusterRespBody, nil)),
		d.Set("port", utils.PathSearch("cluster.port", getDwsClusterRespBody, nil)),
		d.Set("endpoints", flattenGetDwsClusterRespBodyEndpoint(getDwsClusterRespBody)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("cluster.tags", getDwsClusterRespBody, nil))),
		d.Set("user_name", utils.PathSearch("cluster.user_name", getDwsClusterRespBody, nil)),
		d.Set("number_of_node", utils.PathSearch("cluster.number_of_node", getDwsClusterRespBody, nil)),
		d.Set("availability_zone", utils.PathSearch("cluster.availability_zone", getDwsClusterRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("cluster.enterprise_project_id", getDwsClusterRespBody, nil)),
		d.Set("node_type", utils.PathSearch("cluster.node_type", getDwsClusterRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("cluster.vpc_id", getDwsClusterRespBody, nil)),
		d.Set("network_id", utils.PathSearch("cluster.subnet_id", getDwsClusterRespBody, nil)),
		d.Set("security_group_id", utils.PathSearch("cluster.security_group_id", getDwsClusterRespBody, nil)),
		d.Set("public_ip", flattenGetDwsClusterRespBodyPublicIp(getDwsClusterRespBody)),
		d.Set("public_endpoints", flattenGetDwsClusterRespBodyPublicEndpoint(getDwsClusterRespBody)),
		d.Set("sub_status", utils.PathSearch("cluster.sub_status", getDwsClusterRespBody, nil)),
		d.Set("task_status", utils.PathSearch("cluster.task_status", getDwsClusterRespBody, nil)),
		d.Set("recent_event", utils.PathSearch("cluster.recent_event", getDwsClusterRespBody, nil)),
		d.Set("private_ip", utils.PathSearch("cluster.private_ip", getDwsClusterRespBody, nil)),
		d.Set("maintain_window", flattenGetDwsClusterRespBodyMaintainWindow(getDwsClusterRespBody)),
		d.Set("elb", flattenGetDwsClusterRespBodyElb(getDwsClusterRespBody)),
		d.Set("description", utils.PathSearch("cluster.cluster_description_info", getDwsClusterRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDwsClusterRespBodyEndpoint(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("cluster.endpoints", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"connect_info": utils.PathSearch("connect_info", v, nil),
			"jdbc_url":     utils.PathSearch("jdbc_url", v, nil),
		})
	}
	return rst
}

func flattenGetDwsClusterRespBodyPublicIp(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("cluster.public_ip", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"public_bind_type": utils.PathSearch("public_bind_type", curJson, nil),
			"eip_id":           utils.PathSearch("eip_id", curJson, nil),
		},
	}
	return rst
}

func flattenGetDwsClusterRespBodyPublicEndpoint(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("cluster.public_endpoints", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"public_connect_info": utils.PathSearch("public_connect_info", v, nil),
			"jdbc_url":            utils.PathSearch("jdbc_url", v, nil),
		})
	}
	return rst
}

func flattenGetDwsClusterRespBodyMaintainWindow(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("cluster.maintain_window", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"day":        utils.PathSearch("day", curJson, nil),
			"start_time": utils.PathSearch("start_time", curJson, nil),
			"end_time":   utils.PathSearch("end_time", curJson, nil),
		},
	}
	return rst
}

func flattenGetDwsClusterRespBodyElb(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("cluster.elb", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"name":             utils.PathSearch("name", curJson, nil),
			"id":               utils.PathSearch("id", curJson, nil),
			"public_ip":        utils.PathSearch("public_ip", curJson, nil),
			"private_ip":       utils.PathSearch("private_ip", curJson, nil),
			"private_endpoint": utils.PathSearch("private_endpoint", curJson, nil),
			"vpc_id":           utils.PathSearch("vpc_id", curJson, nil),
			"private_ip_v6":    utils.PathSearch("private_ip_v6", curJson, nil),
		},
	}

	return rst
}

func resourceDwsClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	clusterId := d.Id()
	clusterClient, clientErr := cfg.NewServiceClient("dws", region)
	if clientErr != nil {
		return diag.Errorf("error creating DWS client: %s", clientErr)
	}

	err := clusterWaitingForAvailable(ctx, d, clusterClient, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.Errorf("cluster (%s) state is not available to update: %s", clusterId, err)
	}

	expandInstanceStorageChanges := []string{
		"volume.0.capacity",
	}

	if d.HasChanges(expandInstanceStorageChanges...) {
		// expandInstanceStorage: expand instance storage
		var (
			expandInstanceStorageHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/expand-instance-storage"
		)

		expandInstanceStoragePath := clusterClient.Endpoint + expandInstanceStorageHttpUrl
		expandInstanceStoragePath = strings.ReplaceAll(expandInstanceStoragePath, "{project_id}", clusterClient.ProjectID)
		expandInstanceStoragePath = strings.ReplaceAll(expandInstanceStoragePath, "{cluster_id}", d.Id())

		expandInstanceStorageOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      requestOpts.MoreHeaders,
		}

		expandInstanceStorageOpt.JSONBody = utils.RemoveNil(buildExpandInstanceStorageBodyParams(d))
		_, err = clusterClient.Request("POST", expandInstanceStoragePath, &expandInstanceStorageOpt)
		if err != nil {
			return diag.Errorf("error updating DWS cluster: %s", err)
		}
		err = clusterWaitingForAvailable(ctx, d, clusterClient, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the DWS cluster (%s) update to complete: %s", clusterId, err)
		}
	}
	resetPasswordOfClusterChanges := []string{
		"user_pwd",
	}

	if d.HasChanges(resetPasswordOfClusterChanges...) {
		// resetPasswordOfCluster: reset password of DWS cluster
		var (
			resetPasswordOfClusterHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/reset-password"
		)

		resetPasswordOfClusterPath := clusterClient.Endpoint + resetPasswordOfClusterHttpUrl
		resetPasswordOfClusterPath = strings.ReplaceAll(resetPasswordOfClusterPath, "{project_id}", clusterClient.ProjectID)
		resetPasswordOfClusterPath = strings.ReplaceAll(resetPasswordOfClusterPath, "{cluster_id}", d.Id())

		resetPasswordOfClusterOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      requestOpts.MoreHeaders,
		}

		resetPasswordOfClusterOpt.JSONBody = utils.RemoveNil(buildResetPasswordOfClusterBodyParams(d))
		_, err = clusterClient.Request("POST", resetPasswordOfClusterPath, &resetPasswordOfClusterOpt)
		if err != nil {
			return diag.Errorf("error updating DWS cluster: %s", err)
		}

		err = clusterWaitingForAvailable(ctx, d, clusterClient, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the DWS cluster (%s) update to complete: %s", clusterId, err)
		}
	}

	if d.HasChange("number_of_node") {
		if err := updateClusterNodes(ctx, clusterClient, d, clusterId); err != nil {
			return diag.FromErr(err)
		}
	}

	// change tags
	if d.HasChange("tags") {
		err = updateClusterTags(clusterClient, d, clusterId)
		if err != nil {
			return diag.Errorf("error updating tags of DWS cluster:%s, err:%s", clusterId, err)
		}
	}

	if d.HasChange("logical_cluster_enable") {
		client, err := cfg.NewServiceClient("dws", region)
		if err != nil {
			return diag.Errorf("error creating DWS client: %s", err)
		}

		if err := updateDwsLogicalClusterEnable(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("elb_id") {
		oldElbIdRaw, newElbIdRaw := d.GetChange("elb_id")
		oldElbId := oldElbIdRaw.(string)
		newElbId := newElbIdRaw.(string)

		err = unbindElb(ctx, d, clusterClient, oldElbId)
		if err != nil {
			return diag.FromErr(err)
		}

		err = bindElb(ctx, d, clusterClient, newElbId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   clusterId,
			ResourceType: "dws_clusters",
			RegionId:     region,
			ProjectId:    clusterClient.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("lts_enable") {
		err = enableOrDisableLts(d, clusterClient)
		if err != nil {
			err = parseLtsError(err)
			if err != nil {
				return common.CheckDeletedDiag(d, err,
					fmt.Sprintf("error modifying LTS for DWS cluster, the expected LTS enable status is: %v", d.Get("lts_enable").(bool)))
			}
		}
	}

	if d.HasChange("public_ip.0.eip_id") {
		if err := updateEip(clusterClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("description") {
		if err := updateDescription(clusterClient, clusterId, d.Get("description").(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("security_group_id") {
		if err := updateSecurityGroup(clusterClient, clusterId, d.Get("security_group_id").(string)); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceDwsClusterRead(ctx, d, meta)
}

func buildExpandInstanceStorageBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_size": utils.ValueIgnoreEmpty(d.Get("volume.0.capacity")),
	}
	return bodyParams
}

func buildResetPasswordOfClusterBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_password": utils.ValueIgnoreEmpty(d.Get("user_pwd")),
	}
	return bodyParams
}

func updateClusterNodes(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, clusterId string) error {
	var (
		oldValue, newValue = d.GetChange("number_of_node")
		nodeNum            = newValue.(int) - oldValue.(int)
		err                error
	)

	if nodeNum > 0 {
		err = scaleOutCluster(client, clusterId, nodeNum)
	} else {
		err = scaleInCluster(client, d, clusterId, -nodeNum)
	}

	if err != nil {
		return err
	}

	return waitClusterTaskStateCompleted(ctx, client, d.Timeout(schema.TimeoutUpdate), clusterId)
}

func scaleOutCluster(client *golangsdk.ServiceClient, clusterId string, scaleOutNodes int) error {
	scaleOutClusterHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/resize"
	scaleOutClusterPath := client.Endpoint + scaleOutClusterHttpUrl
	scaleOutClusterPath = strings.ReplaceAll(scaleOutClusterPath, "{project_id}", client.ProjectID)
	scaleOutClusterPath = strings.ReplaceAll(scaleOutClusterPath, "{cluster_id}", clusterId)

	scaleOutClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
		JSONBody: map[string]interface{}{
			"scale_out": map[string]interface{}{
				"count": scaleOutNodes,
			},
		},
	}
	_, err := client.Request("POST", scaleOutClusterPath, &scaleOutClusterOpt)
	if err != nil {
		return fmt.Errorf("error extending nodes of the DWS cluster (%s) : %s", clusterId, err)
	}
	return nil
}

func buildScaleInBodyParams(d *schema.ResourceData, shrinkNum int, datastoreType string) map[string]interface{} {
	return map[string]interface{}{
		"shrink_number": shrinkNum,
		"force_backup":  d.Get("force_backup"),
		"type":          datastoreType,
	}
}

func scaleInCluster(client *golangsdk.ServiceClient, d *schema.ResourceData, clusterId string, shrinkNum int) error {
	clusterResp, err := GetClusterInfoByClusterId(client, clusterId)
	if err != nil {
		return err
	}

	datastoreType := utils.PathSearch("cluster.datastore_type", clusterResp, "").(string)
	if datastoreType == "" {
		return fmt.Errorf("unable to get datastore type of the cluster (%s)", clusterId)
	}

	httpUrl := "v1.0/{project_id}/clusters/{cluster_id}/cluster-shrink"
	scaleInPath := client.Endpoint + httpUrl
	scaleInPath = strings.ReplaceAll(scaleInPath, "{project_id}", client.ProjectID)
	scaleInPath = strings.ReplaceAll(scaleInPath, "{cluster_id}", d.Id())

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
		JSONBody:         utils.RemoveNil(buildScaleInBodyParams(d, shrinkNum, datastoreType)),
	}
	resp, err := client.Request("POST", scaleInPath, &opt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	// a. In some cases, the status code is 200 when the scaling-in fails, such as: The scale-in number is incorrect.
	// b. The two situations of successful scale-in are as follows:
	//    1. When the status id `200`, the error_code` is `DWS.0000`, it means the scale-in is successful.
	//    2. When the status id `202`, the job_id` is not empty, it means the scale-in is successful.
	errCode := utils.PathSearch("error_code", respBody, "").(string)
	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if errCode == "DWS.0000" || jobId != "" {
		return nil
	}

	errMsg := utils.PathSearch("error_msg", respBody, "").(string)
	return fmt.Errorf("error shrinking nodes of the DWS cluster (%s) : %s", clusterId, errMsg)
}

func resourceDwsClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDwsCluster: delete DWS cluster
	var (
		deleteDwsClusterHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}"
		deleteDwsClusterProduct = "dws"
	)
	deleteDwsClusterClient, err := cfg.NewServiceClient(deleteDwsClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	deleteDwsClusterPath := deleteDwsClusterClient.Endpoint + deleteDwsClusterHttpUrl
	deleteDwsClusterPath = strings.ReplaceAll(deleteDwsClusterPath, "{project_id}", deleteDwsClusterClient.ProjectID)
	deleteDwsClusterPath = strings.ReplaceAll(deleteDwsClusterPath, "{cluster_id}", d.Id())

	deleteDwsClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	deleteDwsClusterOpt.JSONBody = utils.RemoveNil(buildDeleteDwsClusterBodyParams(d))
	_, err = deleteDwsClusterClient.Request("DELETE", deleteDwsClusterPath, &deleteDwsClusterOpt)
	if err != nil {
		return diag.Errorf("error deleting DWS cluster: %s", err)
	}

	err = deleteClusterWaitingForCompleted(ctx, d, deleteDwsClusterClient, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the DWS cluster (%s) deletion to complete: %s", d.Id(), err)
	}
	return nil
}

func buildDeleteDwsClusterBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"keep_last_manual_snapshot": d.Get("keep_last_manual_snapshot"),
	}
	return bodyParams
}

func deleteClusterWaitingForCompleted(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			// deleteDwsClusterWaiting: missing operation notes
			deleteDwsClusterWaitingRespBody, err := GetClusterInfoByClusterId(client, d.Id())
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			status := utils.PathSearch(`cluster.status`, deleteDwsClusterWaitingRespBody, "").(string)

			targetStatus := []string{
				"DELETED",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return deleteDwsClusterWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"FAILED",
				"DELETE_FAILED",
				"FROZEN",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return deleteDwsClusterWaitingRespBody, status, nil
			}

			return deleteDwsClusterWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func addClusterTags(client *golangsdk.ServiceClient, clusterId string, rawTags []tags.ResourceTag) error {
	var (
		addTagsHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/tags/batch-create"
	)

	addTagsPath := client.Endpoint + addTagsHttpUrl
	addTagsPath = strings.ReplaceAll(addTagsPath, "{project_id}", client.ProjectID)
	addTagsPath = strings.ReplaceAll(addTagsPath, "{cluster_id}", clusterId)

	addTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	addTagsOpt.JSONBody = map[string]interface{}{
		"tags": rawTags,
	}
	_, err := client.Request("POST", addTagsPath, &addTagsOpt)
	if err != nil {
		return fmt.Errorf("error adding tags of DWS cluster: %s", err)
	}

	return nil
}

func deleteClusterTags(client *golangsdk.ServiceClient, clusterId string, rawTags []tags.ResourceTag) error {
	var (
		deleteTagsHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/tags/batch-delete"
	)

	deleteTagsPath := client.Endpoint + deleteTagsHttpUrl
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{project_id}", client.ProjectID)
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{cluster_id}", clusterId)

	deleteTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	deleteTagsOpt.JSONBody = map[string]interface{}{
		"tags": rawTags,
	}
	_, err := client.Request("POST", deleteTagsPath, &deleteTagsOpt)
	if err != nil {
		return fmt.Errorf("error deleting tags of DWS cluster: %s", err)
	}

	return nil
}

func updateClusterTags(client *golangsdk.ServiceClient, d *schema.ResourceData, id string) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	// remove old tags
	if len(oMap) > 0 {
		taglist := utils.ExpandResourceTags(oMap)
		err := deleteClusterTags(client, id, taglist)
		if err != nil {
			return err
		}
	}

	// set new tags
	if len(nMap) > 0 {
		taglist := utils.ExpandResourceTags(nMap)
		err := addClusterTags(client, id, taglist)
		if err != nil {
			return err
		}
	}

	return nil
}

func bindElb(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, elbId string) error {
	if elbId == "" {
		return nil
	}

	bindElbHttpUrl := "v2/{project_id}/clusters/{cluster_id}/elbs/{elb_id}"

	bindElbPath := client.Endpoint + bindElbHttpUrl
	bindElbPath = strings.ReplaceAll(bindElbPath, "{project_id}", client.ProjectID)
	bindElbPath = strings.ReplaceAll(bindElbPath, "{cluster_id}", d.Id())
	bindElbPath = strings.ReplaceAll(bindElbPath, "{elb_id}", elbId)

	bindElbOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	bindElbResp, err := client.Request("POST", bindElbPath, &bindElbOpt)
	if err != nil {
		return fmt.Errorf("error binding ELB to DWS cluster: %s", err)
	}

	bindElbRespBody, err := utils.FlattenResponse(bindElbResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", bindElbRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error binding ELB to DWS cluster: job ID is not found in API response")
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"INIT"},
		Target:       []string{"SUCCESS"},
		Refresh:      jobStatusRefreshFunc(client, jobId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        60 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for binding ELB to DWS cluster: %s", err)
	}

	return nil
}

func unbindElb(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, elbId string) error {
	if elbId == "" {
		return nil
	}

	unbindElbHttpUrl := "v2/{project_id}/clusters/{cluster_id}/elbs/{elb_id}"

	unbindElbPath := client.Endpoint + unbindElbHttpUrl
	unbindElbPath = strings.ReplaceAll(unbindElbPath, "{project_id}", client.ProjectID)
	unbindElbPath = strings.ReplaceAll(unbindElbPath, "{cluster_id}", d.Id())
	unbindElbPath = strings.ReplaceAll(unbindElbPath, "{elb_id}", elbId)

	bindElbOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	unbindElbResp, err := client.Request("DELETE", unbindElbPath, &bindElbOpt)
	if err != nil {
		return fmt.Errorf("error unbinding ELB from DWS cluster: %s", err)
	}

	unbindElbRespBody, err := utils.FlattenResponse(unbindElbResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", unbindElbRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error unbinding ELB from DWS cluster: job ID is not found in API response: %s", jobId)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"INIT"},
		Target:       []string{"SUCCESS"},
		Refresh:      jobStatusRefreshFunc(client, jobId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        60 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for unbinding ELB from DWS cluster: %s", err)
	}

	return nil
}

func updateEip(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		oldEipRaw, newEipRaw = d.GetChange("public_ip.0.eip_id")
		oldEipId             = oldEipRaw.(string)
		newEipId             = newEipRaw.(string)
		clusterId            = d.Id()
	)

	path := client.Endpoint + "v2/{project_id}/clusters/{cluster_id}/eips/{eip_id}"
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{cluster_id}", clusterId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	if oldEipId != "" {
		unBindEipPath := strings.ReplaceAll(path, "{eip_id}", oldEipId)
		_, err := client.Request("DELETE", unBindEipPath, &opt)
		if err != nil {
			return fmt.Errorf("error unbinding EIP (%s) from DWS instance (%s): %s", oldEipId, clusterId, err)
		}
	}

	if newEipId != "" {
		bindEipPath := strings.ReplaceAll(path, "{eip_id}", newEipId)
		_, err := client.Request("POST", bindEipPath, &opt)
		if err != nil {
			return fmt.Errorf("error binding EIP (%s) to DWS instance (%s): %s", newEipId, clusterId, err)
		}
	}
	return nil
}

func jobStatusRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getJobStatusHttpUrl := "v1.0/{project_id}/job/{job_id}"

		getJobStatusPath := client.Endpoint + getJobStatusHttpUrl
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{project_id}", client.ProjectID)
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{job_id}", jobId)

		getJobStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      requestOpts.MoreHeaders,
		}
		getJobStatusResp, err := client.Request("GET", getJobStatusPath, &getJobStatusOpt)
		if err != nil {
			return getJobStatusResp, "FAIL", err
		}

		getJobStatusRespBody, err := utils.FlattenResponse(getJobStatusResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("status", getJobStatusRespBody, "")
		if status.(string) == "FAIL" {
			failedCode := utils.PathSearch("failed_code", getJobStatusRespBody, "")
			failedDetail := utils.PathSearch("failed_detail", getJobStatusRespBody, "")
			return nil, "", fmt.Errorf("DWS cluster binding ELB job failed,"+
				" job ID: %s, failed_code: %s, failed_detail: %s", jobId, failedCode, failedDetail)
		}
		return getJobStatusRespBody, status.(string), nil
	}
}

func enableOrDisableLts(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	ltsHttpUrl := "v1/{project_id}/clusters/{cluster_id}/lts-logs/{action}"

	ltsPath := client.Endpoint + ltsHttpUrl
	ltsPath = strings.ReplaceAll(ltsPath, "{project_id}", client.ProjectID)
	ltsPath = strings.ReplaceAll(ltsPath, "{cluster_id}", d.Id())
	operate := d.Get("lts_enable").(bool)
	if operate {
		ltsPath = strings.ReplaceAll(ltsPath, "{action}", "enable")
	} else {
		ltsPath = strings.ReplaceAll(ltsPath, "{action}", "disable")
	}
	ltsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
		JSONBody:         map[string]interface{}{},
	}

	_, err := client.Request("POST", ltsPath, &ltsOpt)
	if err != nil {
		return err
	}

	return nil
}

func parseLtsError(err error) error {
	var errCode400 golangsdk.ErrDefault400
	if errors.As(err, &errCode400) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode400.Body, &apiError); jsonErr != nil {
			if decodeRes, decodeErr := base64.URLEncoding.DecodeString(string(errCode400.Body)); decodeErr == nil {
				if jsonErr = json.Unmarshal(decodeRes, &apiError); jsonErr != nil {
					return err
				}
			}
		}
		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return err
		}
		// error code DWS.7107 means the cluster LTS is disable; DWS.0015 means the cluster not exists.
		if errorCode == "DWS.7107" {
			return nil
		}
		if errorCode == "DWS.0015" {
			return golangsdk.ErrDefault404(errCode400)
		}
	}
	return err
}

func parseClusterNotFoundError(err error) error {
	parsedErr := common.ConvertExpected401ErrInto404Err(err, "error_code", "DWS.0047")
	if _, ok := parsedErr.(golangsdk.ErrDefault404); ok {
		return parsedErr
	}

	// "DWS.0015": The cluster ID does not exist (standard UUID format). Status code is 403.
	parsedErr = common.ConvertExpected403ErrInto404Err(err, "error_code", "DWS.0015")
	// "DWS.3027": The cluster was deleted after it was created. Status code is 404.
	if _, ok := parsedErr.(golangsdk.ErrDefault404); ok {
		return parsedErr
	}
	return err
}

func updateDescription(client *golangsdk.ServiceClient, clusterId, description string) error {
	httpUrl := "v1/{project_id}/clusters/{cluster_id}/description"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{cluster_id}", clusterId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		// The description can be changed to empty.
		JSONBody: map[string]interface{}{
			"description_info": description,
		},
	}
	_, err := client.Request("POST", path, &opt)
	if err != nil {
		return fmt.Errorf("unable to set description for the cluster (%s): %s", clusterId, err)
	}

	return nil
}

func updateSecurityGroup(client *golangsdk.ServiceClient, clusterId, securityGroupId string) error {
	httpUrl := "v1/{project_id}/clusters/{cluster_id}/security-group"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{cluster_id}", clusterId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"security_groups": []string{securityGroupId},
		},
	}
	_, err := client.Request("PUT", path, &opt)
	if err != nil {
		return fmt.Errorf("error updating security group for the cluster (%s): %s", clusterId, err)
	}

	return nil
}
