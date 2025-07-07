package css

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSS POST /v1.0/{project_id}/clusters
// @API CSS DELETE /v1.0/{project_id}/clusters/{cluster_id}
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/role_extend
// @API CSS POST /v1.0/extend/{project_id}/clusters/{cluster_id}/role/shrink
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/changename
// @API CSS POST /v1.0/{project_id}/{resource_type}/{cluster_id}/tags/action
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/route
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/route
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/sg/change
// @API CSS POST /v1.0/{project_id}/cluster/{cluster_id}/period
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
func ResourceLogstashCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogstashClusterCreate,
		ReadContext:   resourceLogstashClusterRead,
		UpdateContext: resourceLogstashClusterUpdate,
		DeleteContext: resourceLogstashClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"instance_number": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"volume": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
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
				Required: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// charging_mode, period_unit and period only support changing post-paid to pre-paid billing mode.
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"period"},
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"period_unit"},
			},
			"auto_renew": common.SchemaAutoRenewUpdatable(nil),
			"routes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_net_mask": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"engine_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_period": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceLogstashClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	createLogstashClusterHttpUrl := "v1.0/{project_id}/clusters"
	createLogstashClusterPath := client.Endpoint + createLogstashClusterHttpUrl
	createLogstashClusterPath = strings.ReplaceAll(createLogstashClusterPath, "{project_id}", client.ProjectID)

	createLogstashClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	bodyParams, paramErr := buildlogstashClusterCreateParameters(d, conf)
	if paramErr != nil {
		return diag.FromErr(paramErr)
	}
	createLogstashClusterOpt.JSONBody = utils.RemoveNil(bodyParams)
	createLogstashClusterResp, err := client.Request("POST", createLogstashClusterPath, &createLogstashClusterOpt)
	if err != nil {
		return diag.Errorf("error creating CSS logstash cluster: %s", err)
	}

	createLogstashClusterRespBody, err := utils.FlattenResponse(createLogstashClusterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := utils.PathSearch("cluster.id", createLogstashClusterRespBody, "").(string)
	orderId := utils.PathSearch("orderId", createLogstashClusterRespBody, "").(string)

	if orderId == "" {
		if clusterId == "" {
			return diag.Errorf("error creating CSS cluster: id is not found in API response,%#v", createLogstashClusterRespBody)
		}
		d.SetId(clusterId)
	} else {
		bssClient, err := conf.BssV2Client(conf.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}

		// 1. If charging mode is PrePaid, wait for the order to be completed.
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		// 2. get the resource ID, must be after order success
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resourceId)
	}

	createResultErr := checkClusterCreateResult(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if createResultErr != nil {
		return diag.FromErr(createResultErr)
	}

	if v, ok := d.GetOk("routes"); ok {
		err := updateClusterRoute(conf, d, v.(*schema.Set).List(), "add_ip")
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceLogstashClusterRead(ctx, d, meta)
}

func buildlogstashClusterCreateParameters(d *schema.ResourceData, conf *config.Config) (map[string]interface{}, error) {
	cluster := map[string]interface{}{
		"name": d.Get("name").(string),
		"datastore": map[string]interface{}{
			"type":    "logstash",
			"version": d.Get("engine_version"),
		},
		"enterprise_project_id": utils.ValueIgnoreEmpty(conf.GetEnterpriseProjectID(d)),
		"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"instance": map[string]interface{}{
			"flavorRef": d.Get("node_config.0.flavor"),
			"nics": map[string]interface{}{
				"vpcId":           d.Get("vpc_id"),
				"netId":           d.Get("subnet_id"),
				"securityGroupId": d.Get("security_group_id"),
			},
			"volume": map[string]interface{}{
				"size":        d.Get("node_config.0.volume.0.size"),
				"volume_type": d.Get("node_config.0.volume.0.volume_type"),
			},
			"availability_zone": utils.ValueIgnoreEmpty(d.Get("availability_zone")),
		},
		"instanceNum": d.Get("node_config.0.instance_number"),
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		payInfo := map[string]interface{}{
			"period":    d.Get("period"),
			"isAutoPay": 1,
		}

		payModel := d.Get("period_unit").(string)
		if payModel == "month" {
			payInfo["payModel"] = 2
		} else {
			payInfo["payModel"] = 3
		}

		if d.Get("auto_renew").(string) == "true" {
			payInfo["isAutoRenew"] = 1
		}
		cluster["payInfo"] = payInfo
	}

	bodyParams := map[string]interface{}{
		"cluster": cluster,
	}
	return bodyParams, nil
}

func resourceLogstashClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	v1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	clusterDetail, err := getClusterDetails(v1Client, d.Id())
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"), "error getting CSS cluster")
	}

	getRoutesRespBody, err := getClusterRoute(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", clusterDetail, nil)),
		d.Set("engine_type", utils.PathSearch("datastore.type", clusterDetail, nil)),
		d.Set("engine_version", utils.PathSearch("datastore.version", clusterDetail, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterpriseProjectId", clusterDetail, nil)),
		d.Set("vpc_id", utils.PathSearch("vpcId", clusterDetail, nil)),
		d.Set("subnet_id", utils.PathSearch("subnetId", clusterDetail, nil)),
		d.Set("security_group_id", utils.PathSearch("securityGroupId", clusterDetail, nil)),
		d.Set("nodes", flattenClusterNodes(utils.PathSearch("instances", clusterDetail, make([]interface{}, 0)).([]interface{}))),
		setLogstashNodeConfigsAndAz(clusterDetail, d),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", clusterDetail, nil))),
		d.Set("created_at", utils.PathSearch("created", clusterDetail, nil)),
		d.Set("endpoint", utils.PathSearch("endpoint", clusterDetail, nil)),
		d.Set("status", utils.PathSearch("status", clusterDetail, nil)),
		d.Set("routes", flattenGetRoute(getRoutesRespBody)),
		d.Set("updated_at", utils.PathSearch("updated", clusterDetail, nil)),
		d.Set("is_period", utils.PathSearch("period", clusterDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func setLogstashNodeConfigsAndAz(clusterDetail interface{}, d *schema.ResourceData) error {
	instances := utils.PathSearch("instances", clusterDetail, make([]interface{}, 0)).([]interface{})
	if len(instances) == 0 {
		return nil
	}

	var azArray []string
	for _, v := range instances {
		azArray = append(azArray, utils.PathSearch("azCode", v, "").(string))
	}
	azArray = utils.RemoveDuplicateElem(azArray)
	az := strings.Join(azArray, ",")

	mErr := multierror.Append(nil,
		d.Set("availability_zone", az),
		d.Set("node_config", []interface{}{map[string]interface{}{
			"flavor":          utils.PathSearch("[0].specCode", instances, nil),
			"instance_number": len(instances),
			"volume": []interface{}{map[string]interface{}{
				"size":        int(utils.PathSearch("[0].volume.size", instances, float64(0)).(float64)),
				"volume_type": utils.PathSearch("[0].volume.type", instances, nil),
			}},
		}}),
	)
	return mErr.ErrorOrNil()
}

func resourceLogstashClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	clusterId := d.Id()
	client, err := cfg.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	// update cluster name
	if d.HasChanges("name") {
		err := updateClusterName(client, clusterId, d.Get("name").(string))
		if err != nil {
			return diag.Errorf("error updating CSS logstash cluster name, cluster_id: %s, error: %s", clusterId, err)
		}
	}

	// extend and shrink cluster
	if d.HasChanges("node_config") {
		err = modifyLogstashCluster(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update security group ID
	if d.HasChange("security_group_id") {
		err = updateSecurityGroup(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		client, err := cfg.CssV1Client(region)
		if err != nil {
			return diag.Errorf("error creating CSS V1 client: %s", err)
		}
		tagErr := utils.UpdateResourceTags(client, d, "css-cluster", clusterId)
		if tagErr != nil {
			return diag.Errorf("Error updating tags of CSS logstash cluster:%s, err:%s", clusterId, tagErr)
		}
	}

	if d.HasChanges("charging_mode", "auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}

		err = updateChangingModeOrAutoRenew(ctx, d, client, bssClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   clusterId,
			ResourceType: "css-cluster",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("routes") {
		oldRaws, newRaws := d.GetChange("routes")
		addRaws := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
		delRaws := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
		err := updateClusterRoute(cfg, d, delRaws.List(), "del_ip")
		if err != nil {
			return diag.FromErr(err)
		}
		err = updateClusterRoute(cfg, d, addRaws.List(), "add_ip")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLogstashClusterRead(ctx, d, meta)
}

func updateClusterName(client *golangsdk.ServiceClient, clusterId, displayName string) error {
	updateClusterNameHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/changename"
	updateClusterNamePath := client.Endpoint + updateClusterNameHttpUrl
	updateClusterNamePath = strings.ReplaceAll(updateClusterNamePath, "{project_id}", client.ProjectID)
	updateClusterNamePath = strings.ReplaceAll(updateClusterNamePath, "{cluster_id}", clusterId)

	updateClusterNameOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"displayName": displayName,
		},
	}

	_, err := client.Request("POST", updateClusterNamePath, &updateClusterNameOpt)
	if err != nil {
		return err
	}
	return nil
}

func modifyLogstashCluster(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if d.HasChange("node_config") {
		oldv, newv := d.GetChange("node_config.0.instance_number")
		oldNodeNum := oldv.(int)
		newNodeNum := newv.(int)
		if newNodeNum > oldNodeNum {
			return extendLogstashCluster(ctx, d, newNodeNum-oldNodeNum, client)
		}

		// shrink
		if oldNodeNum > newNodeNum {
			azSplits := strings.Split(d.Get("availability_zone").(string), ",")
			if newNodeNum < len(azSplits) {
				return fmt.Errorf("the number of remaining nodes after scale-in" +
					" must be greater than or equal to the number of Azs")
			}
			return shrinkLogstashCluster(ctx, d, oldNodeNum-newNodeNum, client)
		}
	}

	return nil
}

func extendLogstashCluster(ctx context.Context, d *schema.ResourceData, extendNodesize int, client *golangsdk.ServiceClient) error {
	extendLogstashClusterHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/role_extend"
	extendLogstashClusterPath := client.Endpoint + extendLogstashClusterHttpUrl
	extendLogstashClusterPath = strings.ReplaceAll(extendLogstashClusterPath, "{project_id}", client.ProjectID)
	extendLogstashClusterPath = strings.ReplaceAll(extendLogstashClusterPath, "{cluster_id}", d.Id())

	extendLogstashClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"grow": []map[string]interface{}{
				{
					"type":     "lgs",
					"nodesize": extendNodesize,
					"disksize": 0,
				},
			},
			"isAutoPay": 1,
		},
	}

	_, err := client.Request("POST", extendLogstashClusterPath, &extendLogstashClusterOpt)
	if err != nil {
		return fmt.Errorf("error extending CSS logstash cluster (%s) instance number failed: %s", d.Id(), err)
	}

	err = checkClusterOperationResult(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

func shrinkLogstashCluster(ctx context.Context, d *schema.ResourceData, shrinkNodesize int, client *golangsdk.ServiceClient) error {
	shrinkLogstashClusterHttpUrl := "v1.0/extend/{project_id}/clusters/{cluster_id}/role/shrink"
	shrinkLogstashClusterPath := client.Endpoint + shrinkLogstashClusterHttpUrl
	shrinkLogstashClusterPath = strings.ReplaceAll(shrinkLogstashClusterPath, "{project_id}", client.ProjectID)
	shrinkLogstashClusterPath = strings.ReplaceAll(shrinkLogstashClusterPath, "{cluster_id}", d.Id())

	shrinkLogstashClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"shrink": []map[string]interface{}{
				{
					"type":           "lgs",
					"reducedNodeNum": shrinkNodesize,
				},
			},
		},
	}

	_, err := client.Request("POST", shrinkLogstashClusterPath, &shrinkLogstashClusterOpt)
	if err != nil {
		return fmt.Errorf("error shrinking CSS logstash cluster (%s) instance number failed: %s", d.Id(), err)
	}

	err = checkClusterOperationResult(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

func resourceLogstashClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterId := d.Id()
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	deleteClusterHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}"
	deleteClusterPath := client.Endpoint + deleteClusterHttpUrl
	deleteClusterPath = strings.ReplaceAll(deleteClusterPath, "{project_id}", client.ProjectID)
	deleteClusterPath = strings.ReplaceAll(deleteClusterPath, "{cluster_id}", clusterId)

	deleteClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteClusterPath, &deleteClusterOpt)
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting CSS logstash cluster (%s)", d.Id()))
	}

	err = checkClusterDeleteResult(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("failed to check the result of deletion %s", err)
	}
	return nil
}

func updateClusterRoute(conf *config.Config, d *schema.ResourceData, routes []interface{}, actionType string) error {
	if len(routes) == 0 {
		return nil
	}

	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return fmt.Errorf("error creating CSS V1 client: %s", err)
	}

	updateRouteHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/route"
	updateRoutePath := cssV1Client.Endpoint + updateRouteHttpUrl
	updateRoutePath = strings.ReplaceAll(updateRoutePath, "{project_id}", cssV1Client.ProjectID)
	updateRoutePath = strings.ReplaceAll(updateRoutePath, "{cluster_id}", d.Id())

	updateRouteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for _, route := range routes {
		ipAddress := utils.PathSearch("ip_address", route, "").(string)
		ipNetMask := utils.PathSearch("ip_net_mask", route, "").(string)
		updateRouteOpt.JSONBody = map[string]interface{}{
			"configtype":  actionType,
			"configkey":   ipAddress,
			"configvalue": ipNetMask,
		}
		_, err = cssV1Client.Request("POST", updateRoutePath, &updateRouteOpt)
		if err != nil {
			return fmt.Errorf("error updating CSS logstash route, cluster_id: %s, error: %s", d.Id(), err)
		}
	}

	return nil
}

func getClusterRoute(conf *config.Config, d *schema.ResourceData) (interface{}, error) {
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS V1 client: %s", err)
	}

	getRouteHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/route"
	getRoutePath := cssV1Client.Endpoint + getRouteHttpUrl
	getRoutePath = strings.ReplaceAll(getRoutePath, "{project_id}", cssV1Client.ProjectID)
	getRoutePath = strings.ReplaceAll(getRoutePath, "{cluster_id}", d.Id())

	updateRouteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getRouteResp, err := cssV1Client.Request("GET", getRoutePath, &updateRouteOpt)
	if err != nil {
		return nil, fmt.Errorf("error get CSS logstash routes, cluster_id: %s, error: %s", d.Id(), err)
	}

	getRouteRespBody, err := utils.FlattenResponse(getRouteResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CSS logstash routes response: %s", err)
	}

	return getRouteRespBody, nil
}

func flattenGetRoute(resp interface{}) []interface{} {
	routes := utils.PathSearch("routeResps", resp, make([]interface{}, 0)).([]interface{})

	rst := make([]interface{}, len(routes))
	for i, v := range routes {
		rst[i] = map[string]interface{}{
			"ip_address":  utils.PathSearch("ipAddress", v, nil),
			"ip_net_mask": utils.PathSearch("ipNetMask", v, nil),
		}
	}
	return rst
}

func updateChangingModeOrAutoRenew(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient) error {
	clusterID := d.Id()
	autoRenew := d.Get("auto_renew").(string)
	if d.HasChange("charging_mode") {
		if d.Get("charging_mode").(string) == "postPaid" {
			return fmt.Errorf("error updating the charging mode of the CSS cluster (%s): %s", clusterID,
				"only support changing the CSS cluster form post-paid to pre-paid")
		}

		changeOpts := map[string]interface{}{
			"periodNum": d.Get("period"),
			"isAutoPay": 1,
		}

		if d.Get("period_unit").(string) == "month" {
			changeOpts["periodType"] = 2
		} else {
			changeOpts["periodType"] = 3
		}

		if autoRenew == "true" {
			changeOpts["isAutoRenew"] = 1
		}

		updateClusterToPeriodHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/publickibana/bandwidth"
		updateClusterToPeriodPath := client.Endpoint + updateClusterToPeriodHttpUrl
		updateClusterToPeriodPath = strings.ReplaceAll(updateClusterToPeriodPath, "{project_id}", client.ProjectID)
		updateClusterToPeriodPath = strings.ReplaceAll(updateClusterToPeriodPath, "{cluster_id}", clusterID)

		updateClusterToPeriodOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			JSONBody:         changeOpts,
		}

		resp, err := client.Request("POST", updateClusterToPeriodPath, &updateClusterToPeriodOpt)
		if err != nil {
			return fmt.Errorf("error updating the CSS cluster (%s) form post-paid to pre-paid: %s", clusterID, err)
		}

		orderId := utils.PathSearch("OrderId", resp, "").(string)
		_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	} else if d.HasChange("auto_renew") {
		if err := common.UpdateAutoRenew(bssClient, autoRenew, clusterID); err != nil {
			return fmt.Errorf("error updating the auto-renew of the CSS cluster (%s): %s", clusterID, err)
		}
	}

	return nil
}
