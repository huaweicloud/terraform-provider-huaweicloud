// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DWS
// ---------------------------------------------------------------

package dws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
				ForceNew:    true,
				Description: `The security group ID.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The availability zone in which to create the cluster instance. `,
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "schema: Required",
				RequiredWith: []string{"number_of_cn", "volume"},
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
				ForceNew:    true,
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
				ForceNew: true,
			},
			"volume": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        clusterVolumeSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "schema: Required",
			},
			"tags": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `The key/value pairs to associate with the cluster.`,
			},
			"keep_last_manual_snapshot": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The number of latest manual snapshots that need to be retained when deleting the cluster.`,
			},
			"logical_cluster_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specified whether to enable logical cluster.`,
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
				Computed:    true,
				Description: `The EIP ID.`,
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
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	createDwsClusterPath := createDwsClusterClient.Endpoint + createDwsClusterHttpUrl
	createDwsClusterPath = strings.ReplaceAll(createDwsClusterPath, "{project_id}", createDwsClusterClient.ProjectID)

	createDwsClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	createDwsClusterOpt.JSONBody = utils.RemoveNil(buildCreateDwsClusterBodyParams(d, cfg))
	createDwsClusterResp, err := createDwsClusterClient.Request("POST", createDwsClusterPath, &createDwsClusterOpt)
	if err != nil {
		return diag.Errorf("error creating DWS Cluster: %s", err)
	}

	createDwsClusterRespBody, err := utils.FlattenResponse(createDwsClusterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("cluster.id", createDwsClusterRespBody)
	if err != nil {
		return diag.Errorf("error creating DWS Cluster: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = clusterWaitingForAvailable(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the DWS cluster (%s) creation to complete: %s", d.Id(), err)
	}

	if d.Get("logical_cluster_enable").(bool) {
		if err := updateDwsLogicalClusterEnable(createDwsClusterClient, d); err != nil {
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
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	createDwsClusterPath := createDwsClusterClient.Endpoint + createDwsClusterHttpUrl
	createDwsClusterPath = strings.ReplaceAll(createDwsClusterPath, "{project_id}", createDwsClusterClient.ProjectID)

	createDwsClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	createDwsClusterOpt.JSONBody = utils.RemoveNil(buildCreateDwsClusterBodyParamsV1(d, cfg))
	createDwsClusterResp, err := createDwsClusterClient.Request("POST", createDwsClusterPath, &createDwsClusterOpt)
	if err != nil {
		return diag.Errorf("error creating DWS Cluster: %s", err)
	}

	createDwsClusterRespBody, err := utils.FlattenResponse(createDwsClusterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("cluster.id", createDwsClusterRespBody)
	if err != nil {
		return diag.Errorf("error creating DWS Cluster: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = clusterWaitingForAvailable(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the DWS cluster (%s) creation to complete: %s", d.Id(), err)
	}

	if d.Get("logical_cluster_enable").(bool) {
		if err := updateDwsLogicalClusterEnable(createDwsClusterClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDwsClusterRead(ctx, d, meta)
}

func buildCreateDwsClusterBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cluster": map[string]interface{}{
			"name":                  utils.ValueIngoreEmpty(d.Get("name")),
			"flavor":                utils.ValueIngoreEmpty(d.Get("node_type")),
			"num_node":              utils.ValueIngoreEmpty(d.Get("number_of_node")),
			"num_cn":                utils.ValueIngoreEmpty(d.Get("number_of_cn")),
			"db_name":               utils.ValueIngoreEmpty(d.Get("user_name")),
			"db_password":           utils.ValueIngoreEmpty(d.Get("user_pwd")),
			"db_port":               utils.ValueIngoreEmpty(d.Get("port")),
			"availability_zones":    []string{d.Get("availability_zone").(string)},
			"vpc_id":                utils.ValueIngoreEmpty(d.Get("vpc_id")),
			"subnet_id":             utils.ValueIngoreEmpty(d.Get("network_id")),
			"security_group_id":     utils.ValueIngoreEmpty(d.Get("security_group_id")),
			"datastore_version":     utils.ValueIngoreEmpty(d.Get("version")),
			"dss_pool_id":           utils.ValueIngoreEmpty(d.Get("dss_pool_id")),
			"enterprise_project_id": utils.ValueIngoreEmpty(common.GetEnterpriseProjectID(d, cfg)),
			"master_key_id":         utils.ValueIngoreEmpty(d.Get("kms_key_id")),
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
			"name":                  utils.ValueIngoreEmpty(d.Get("name")),
			"node_type":             utils.ValueIngoreEmpty(d.Get("node_type")),
			"number_of_node":        utils.ValueIngoreEmpty(d.Get("number_of_node")),
			"number_of_cn":          utils.ValueIngoreEmpty(d.Get("number_of_cn")),
			"user_name":             utils.ValueIngoreEmpty(d.Get("user_name")),
			"user_pwd":              utils.ValueIngoreEmpty(d.Get("user_pwd")),
			"port":                  utils.ValueIngoreEmpty(d.Get("port")),
			"availability_zone":     utils.ValueIngoreEmpty(d.Get("availability_zone")),
			"vpc_id":                utils.ValueIngoreEmpty(d.Get("vpc_id")),
			"subnet_id":             utils.ValueIngoreEmpty(d.Get("network_id")),
			"security_group_id":     utils.ValueIngoreEmpty(d.Get("security_group_id")),
			"enterprise_project_id": utils.ValueIngoreEmpty(common.GetEnterpriseProjectID(d, cfg)),
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
			"public_bind_type": utils.ValueIngoreEmpty(raw["public_bind_type"]),
			"eip_id":           utils.ValueIngoreEmpty(raw["eip_id"]),
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
			"volume":   utils.ValueIngoreEmpty(raw["type"]),
			"capacity": utils.ValueIngoreEmpty(raw["capacity"]),
		}
		return params
	}
	return nil
}

func clusterWaitingForAvailable(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// createDwsClusterWaiting: waiting cluster is available
			var (
				createDwsClusterWaitingHttpUrl = "v1.0/{project_id}/clusters/{id}"
				createDwsClusterWaitingProduct = "dws"
			)
			clusterWaitingClient, err := cfg.NewServiceClient(createDwsClusterWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating DWS Client: %s", err)
			}

			clusterWaitingPath := clusterWaitingClient.Endpoint + createDwsClusterWaitingHttpUrl
			clusterWaitingPath = strings.ReplaceAll(clusterWaitingPath, "{project_id}", clusterWaitingClient.ProjectID)
			clusterWaitingPath = strings.ReplaceAll(clusterWaitingPath, "{id}", d.Id())

			clusterWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{
					"Content-Type": "application/json;charset=UTF-8",
				},
			}

			clusterWaitingResp, err := clusterWaitingClient.Request("GET", clusterWaitingPath,
				&clusterWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			clusterWaitingRespBody, err := utils.FlattenResponse(clusterWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}

			actionProgressRaw := utils.PathSearch(`length(cluster.action_progress)`, clusterWaitingRespBody, 0.0)
			if actionProgressRaw.(float64) > 0 {
				return clusterWaitingRespBody, "PENDING", nil
			}

			statusRaw, err := jmespath.Search(`cluster.status`, clusterWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `cluster.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

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

func resourceDwsClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDwsCluster: Query the DWS cluster.
	var (
		getDwsClusterHttpUrl = "v1.0/{project_id}/clusters/{id}"
		getDwsClusterProduct = "dws"
	)
	getDwsClusterClient, err := cfg.NewServiceClient(getDwsClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	getDwsClusterPath := getDwsClusterClient.Endpoint + getDwsClusterHttpUrl
	getDwsClusterPath = strings.ReplaceAll(getDwsClusterPath, "{project_id}", getDwsClusterClient.ProjectID)
	getDwsClusterPath = strings.ReplaceAll(getDwsClusterPath, "{id}", d.Id())

	getDwsClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	getDwsClusterResp, err := getDwsClusterClient.Request("GET", getDwsClusterPath, &getDwsClusterOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, parseClusterNotFoundError(err), "error retrieving DWS Cluster")
	}

	getDwsClusterRespBody, err := utils.FlattenResponse(getDwsClusterResp)
	if err != nil {
		return diag.FromErr(err)
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
	curJson, err := jmespath.Search("cluster.public_ip", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing cluster.public_ip from response= %#v", resp)
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
	curJson, err := jmespath.Search("cluster.maintain_window", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing cluster.maintain_window from response= %#v", resp)
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

func resourceDwsClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	err := clusterWaitingForAvailable(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.Errorf("cluster (%s) state is not available to update: %s", d.Id(), err)
	}

	expandInstanceStorageChanges := []string{
		"volume.0.capacity",
	}

	if d.HasChanges(expandInstanceStorageChanges...) {
		// expandInstanceStorage: expand instance storage
		var (
			expandInstanceStorageHttpUrl = "v1.0/{project_id}/clusters/{id}/expand-instance-storage"
			expandInstanceStorageProduct = "dws"
		)
		expandInstanceStorageClient, err := cfg.NewServiceClient(expandInstanceStorageProduct, region)
		if err != nil {
			return diag.Errorf("error creating DWS Client: %s", err)
		}

		expandInstanceStoragePath := expandInstanceStorageClient.Endpoint + expandInstanceStorageHttpUrl
		expandInstanceStoragePath = strings.ReplaceAll(expandInstanceStoragePath, "{project_id}", expandInstanceStorageClient.ProjectID)
		expandInstanceStoragePath = strings.ReplaceAll(expandInstanceStoragePath, "{id}", d.Id())

		expandInstanceStorageOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=UTF-8",
			},
		}

		expandInstanceStorageOpt.JSONBody = utils.RemoveNil(buildExpandInstanceStorageBodyParams(d))
		_, err = expandInstanceStorageClient.Request("POST", expandInstanceStoragePath, &expandInstanceStorageOpt)
		if err != nil {
			return diag.Errorf("error updating DWS Cluster: %s", err)
		}
		err = clusterWaitingForAvailable(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the DWS cluster (%s) update to complete: %s", d.Id(), err)
		}
	}
	resetPasswordOfClusterChanges := []string{
		"user_pwd",
	}

	if d.HasChanges(resetPasswordOfClusterChanges...) {
		// resetPasswordOfCluster: reset password of DWS cluster
		var (
			resetPasswordOfClusterHttpUrl = "v1.0/{project_id}/clusters/{id}/reset-password"
			resetPasswordOfClusterProduct = "dws"
		)
		resetPasswordOfClusterClient, err := cfg.NewServiceClient(resetPasswordOfClusterProduct, region)
		if err != nil {
			return diag.Errorf("error creating DWS Client: %s", err)
		}

		resetPasswordOfClusterPath := resetPasswordOfClusterClient.Endpoint + resetPasswordOfClusterHttpUrl
		resetPasswordOfClusterPath = strings.ReplaceAll(resetPasswordOfClusterPath, "{project_id}", resetPasswordOfClusterClient.ProjectID)
		resetPasswordOfClusterPath = strings.ReplaceAll(resetPasswordOfClusterPath, "{id}", d.Id())

		resetPasswordOfClusterOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=UTF-8",
			},
		}

		resetPasswordOfClusterOpt.JSONBody = utils.RemoveNil(buildResetPasswordOfClusterBodyParams(d))
		_, err = resetPasswordOfClusterClient.Request("POST", resetPasswordOfClusterPath, &resetPasswordOfClusterOpt)
		if err != nil {
			return diag.Errorf("error updating DWS Cluster: %s", err)
		}

		err = clusterWaitingForAvailable(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the DWS cluster (%s) update to complete: %s", d.Id(), err)
		}
	}
	scaleOutClusterChanges := []string{
		"number_of_node",
	}

	if d.HasChanges(scaleOutClusterChanges...) {
		// scaleOutCluster: Scale out DWS cluster
		var (
			scaleOutClusterHttpUrl = "v1.0/{project_id}/clusters/{id}/resize"
			scaleOutClusterProduct = "dws"
		)
		scaleOutClusterClient, err := cfg.NewServiceClient(scaleOutClusterProduct, region)
		if err != nil {
			return diag.Errorf("error creating DWS Client: %s", err)
		}

		scaleOutClusterPath := scaleOutClusterClient.Endpoint + scaleOutClusterHttpUrl
		scaleOutClusterPath = strings.ReplaceAll(scaleOutClusterPath, "{project_id}", scaleOutClusterClient.ProjectID)
		scaleOutClusterPath = strings.ReplaceAll(scaleOutClusterPath, "{id}", d.Id())

		scaleOutClusterOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=UTF-8",
			},
		}

		scaleOutClusterOpt.JSONBody = utils.RemoveNil(buildScaleOutClusterBodyParams(d))
		_, err = scaleOutClusterClient.Request("POST", scaleOutClusterPath, &scaleOutClusterOpt)
		if err != nil {
			return diag.Errorf("error updating DWS Cluster: %s", err)
		}

		err = clusterWaitingForAvailable(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the DWS cluster (%s) update to complete: %s", d.Id(), err)
		}
	}

	// change tags
	if d.HasChange("tags") {
		clusterClient, err := cfg.NewServiceClient("dws", region)
		if err != nil {
			return diag.Errorf("error creating DWS Client: %s", err)
		}
		err = updateClusterTags(clusterClient, d, d.Id())
		if err != nil {
			return diag.Errorf("error updating tags of DWS cluster:%s, err:%s", d.Id(), err)
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

	return resourceDwsClusterRead(ctx, d, meta)
}

func buildExpandInstanceStorageBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_size": utils.ValueIngoreEmpty(d.Get("volume.0.capacity")),
	}
	return bodyParams
}

func buildResetPasswordOfClusterBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_password": utils.ValueIngoreEmpty(d.Get("user_pwd")),
	}
	return bodyParams
}

func buildScaleOutClusterBodyParams(d *schema.ResourceData) map[string]interface{} {
	oldValue, newValue := d.GetChange("number_of_node")
	num := newValue.(int) - oldValue.(int)

	bodyParams := map[string]interface{}{
		"scale_out": map[string]interface{}{
			"count": num,
		},
	}
	return bodyParams
}

func resourceDwsClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDwsCluster: delete DWS cluster
	var (
		deleteDwsClusterHttpUrl = "v1.0/{project_id}/clusters/{id}"
		deleteDwsClusterProduct = "dws"
	)
	deleteDwsClusterClient, err := cfg.NewServiceClient(deleteDwsClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	deleteDwsClusterPath := deleteDwsClusterClient.Endpoint + deleteDwsClusterHttpUrl
	deleteDwsClusterPath = strings.ReplaceAll(deleteDwsClusterPath, "{project_id}", deleteDwsClusterClient.ProjectID)
	deleteDwsClusterPath = strings.ReplaceAll(deleteDwsClusterPath, "{id}", d.Id())

	deleteDwsClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200, 202,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	deleteDwsClusterOpt.JSONBody = utils.RemoveNil(buildDeleteDwsClusterBodyParams(d))
	_, err = deleteDwsClusterClient.Request("DELETE", deleteDwsClusterPath, &deleteDwsClusterOpt)
	if err != nil {
		return diag.Errorf("error deleting DWS Cluster: %s", err)
	}

	err = deleteClusterWaitingForCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
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

func deleteClusterWaitingForCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// deleteDwsClusterWaiting: missing operation notes
			var (
				deleteDwsClusterWaitingHttpUrl = "v1.0/{project_id}/clusters/{id}"
				deleteDwsClusterWaitingProduct = "dws"
			)
			deleteDwsClusterWaitingClient, err := cfg.NewServiceClient(deleteDwsClusterWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating DWS Client: %s", err)
			}

			deleteDwsClusterWaitingPath := deleteDwsClusterWaitingClient.Endpoint + deleteDwsClusterWaitingHttpUrl
			deleteDwsClusterWaitingPath = strings.ReplaceAll(deleteDwsClusterWaitingPath, "{project_id}", deleteDwsClusterWaitingClient.ProjectID)
			deleteDwsClusterWaitingPath = strings.ReplaceAll(deleteDwsClusterWaitingPath, "{id}", d.Id())

			deleteDwsClusterWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{
					"Content-Type": "application/json;charset=UTF-8",
				},
			}

			deleteDwsClusterWaitingResp, err := deleteDwsClusterWaitingClient.Request("GET", deleteDwsClusterWaitingPath, &deleteDwsClusterWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return deleteDwsClusterWaitingResp, "COMPLETED", nil
				}

				err = parseClusterNotFoundError(err)
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return deleteDwsClusterWaitingResp, "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			deleteDwsClusterWaitingRespBody, err := utils.FlattenResponse(deleteDwsClusterWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`cluster.status`, deleteDwsClusterWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `cluster.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

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

func parseClusterNotFoundError(respErr error) error {
	var apiErr interface{}
	if errCode, ok := respErr.(golangsdk.ErrDefault401); ok {
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr != nil {
			return pErr
		}
		errCode, err := jmespath.Search(`errCode`, apiErr)
		if err != nil {
			return fmt.Errorf("error parse errorCode from response body: %s", err.Error())
		}

		if errCode == `DWS.0047` {
			return golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Body: []byte("the DWS cluster does not exist"),
				},
			}
		}
	}
	return respErr
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
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
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
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
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
