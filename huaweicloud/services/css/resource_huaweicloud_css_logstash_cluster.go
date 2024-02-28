package css

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"

	cssv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1/model"

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
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
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
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
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
		},
	}
}

func resourceLogstashClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	createClusterOpts, paramErr := buildlogstashClusterCreateParameters(d, conf)
	if paramErr != nil {
		return diag.FromErr(paramErr)
	}

	r, err := cssV1Client.CreateCluster(createClusterOpts)
	if err != nil {
		return diag.Errorf("error creating CSS logstash cluster: %s", err)
	}

	if r.OrderId == nil {
		if r.Cluster == nil || r.Cluster.Id == nil {
			return diag.Errorf("error creating CSS logstash cluster: id is not found in API response,%#v", r)
		}
		d.SetId(*r.Cluster.Id)
	} else {
		bssClient, err := conf.BssV2Client(conf.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}

		// 1. If charging mode is PrePaid, wait for the order to be completed.
		err = common.WaitOrderComplete(ctx, bssClient, *r.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		// 2. get the resource ID, must be after order success
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, *r.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resourceId)
	}

	createResultErr := checkClusterCreateResult(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if createResultErr != nil {
		return diag.FromErr(createResultErr)
	}

	return resourceLogstashClusterRead(ctx, d, meta)
}

func buildlogstashClusterCreateParameters(d *schema.ResourceData, conf *config.Config) (*model.CreateClusterRequest, error) {
	createOpts := model.CreateClusterBody{
		Name: d.Get("name").(string),
		Datastore: &model.CreateClusterDatastoreBody{
			Type:    "logstash",
			Version: d.Get("engine_version").(string),
		},
		EnterpriseProjectId: utils.StringIgnoreEmpty(conf.GetEnterpriseProjectID(d)),
		Tags:                buildLogstashCssTags(d.Get("tags").(map[string]interface{})),
		Instance: &model.CreateClusterInstanceBody{
			FlavorRef: d.Get("node_config.0.flavor").(string),
			Nics: &model.CreateClusterInstanceNicsBody{
				VpcId:           d.Get("vpc_id").(string),
				NetId:           d.Get("subnet_id").(string),
				SecurityGroupId: d.Get("security_group_id").(string),
			},
			Volume: &model.CreateClusterInstanceVolumeBody{
				Size:       int32(d.Get("node_config.0.volume.0.size").(int)),
				VolumeType: d.Get("node_config.0.volume.0.volume_type").(string),
			},
			AvailabilityZone: utils.StringIgnoreEmpty(d.Get("availability_zone").(string)),
		},
		InstanceNum: int32(d.Get("node_config.0.instance_number").(int)),
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		createOpts.PayInfo = &model.PayInfoBody{
			Period:    int32(d.Get("period").(int)),
			IsAutoPay: utils.Int32(1),
		}

		if d.Get("period_unit").(string) == "month" {
			createOpts.PayInfo.PayModel = 2
		} else {
			createOpts.PayInfo.PayModel = 3
		}

		if d.Get("auto_renew").(string) == "true" {
			createOpts.PayInfo.IsAutoRenew = utils.Int32(1)
		}
	}

	return &model.CreateClusterRequest{Body: &model.CreateClusterReq{Cluster: &createOpts}}, nil
}

func buildLogstashCssTags(tagmap map[string]interface{}) *[]model.CreateClusterTagsBody {
	var taglist []model.CreateClusterTagsBody

	for k, v := range tagmap {
		tag := model.CreateClusterTagsBody{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	return &taglist
}

func resourceLogstashClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	clusterDetail, err := cssV1Client.ShowClusterDetail(&model.ShowClusterDetailRequest{ClusterId: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CSS logstash cluster")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", clusterDetail.Name),
		d.Set("engine_type", clusterDetail.Datastore.Type),
		d.Set("engine_version", clusterDetail.Datastore.Version),
		d.Set("enterprise_project_id", clusterDetail.EnterpriseProjectId),
		d.Set("vpc_id", clusterDetail.VpcId),
		d.Set("subnet_id", clusterDetail.SubnetId),
		d.Set("security_group_id", clusterDetail.SecurityGroupId),
		d.Set("nodes", flattenClusterNodes(clusterDetail.Instances)),
		setLogstashNodeConfigsAndAz(d, clusterDetail),
		d.Set("tags", flattenTags(clusterDetail.Tags)),
		d.Set("created_at", clusterDetail.Created),
		d.Set("endpoint", clusterDetail.Endpoint),
		d.Set("status", clusterDetail.Status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func setLogstashNodeConfigsAndAz(d *schema.ResourceData, detail *model.ShowClusterDetailResponse) error {
	if detail.Instances == nil || len(*detail.Instances) == 0 {
		return nil
	}

	var azArray []string
	for _, v := range *detail.Instances {
		azArray = append(azArray, utils.StringValue(v.AzCode))
	}
	azArray = utils.RemoveDuplicateElem(azArray)
	az := strings.Join(azArray, ",")

	mErr := multierror.Append(nil,
		d.Set("availability_zone", az),
		d.Set("node_config", []interface{}{map[string]interface{}{
			"flavor":          (*detail.Instances)[0].SpecCode,
			"instance_number": len(*detail.Instances),
			"volume": []interface{}{map[string]interface{}{
				"size":        (*detail.Instances)[0].Volume.Size,
				"volume_type": (*detail.Instances)[0].Volume.Type,
			}},
		}}),
	)
	return mErr.ErrorOrNil()
}

func resourceLogstashClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterId := d.Id()
	cssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	// update cluster name
	if d.HasChanges("name") {
		_, err = cssV1Client.UpdateClusterName(&model.UpdateClusterNameRequest{
			ClusterId: d.Id(),
			Body: &model.UpdateClusterNameReq{
				DisplayName: d.Get("name").(string),
			},
		})
		if err != nil {
			return diag.Errorf("error updating CSS logstash cluster name, cluster_id: %s, error: %s", d.Id(), err)
		}
	}

	// extend and shrink cluster
	if d.HasChanges("node_config") {
		err = modifyLogstashCluster(ctx, d, cssV1Client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		client, err := conf.CssV1Client(region)
		if err != nil {
			return diag.Errorf("error creating CSS V1 client: %s", err)
		}
		tagErr := utils.UpdateResourceTags(client, d, "css-cluster", d.Id())
		if tagErr != nil {
			return diag.Errorf("Error updating tags of CSS logstash cluster:%s, err:%s", d.Id(), tagErr)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := conf.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), clusterId); err != nil {
			return diag.Errorf("error updating the auto-renew of the CSS logstash cluster (%s): %s", clusterId, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := enterpriseprojects.MigrateResourceOpts{
			ResourceId:   clusterId,
			ResourceType: "css-cluster",
			RegionId:     region,
			ProjectId:    conf.GetProjectID(region),
		}
		if err := common.MigrateEnterpriseProject(ctx, conf, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLogstashClusterRead(ctx, d, meta)
}

func modifyLogstashCluster(ctx context.Context, d *schema.ResourceData, cssV1Client *cssv1.CssClient) error {
	if d.HasChange("node_config") {
		oldv, newv := d.GetChange("node_config.0.instance_number")
		oldNodeNum := oldv.(int)
		newNodeNum := newv.(int)
		if newNodeNum > oldNodeNum {
			return extendLogstashCluster(ctx, d, newNodeNum-oldNodeNum, cssV1Client)
		}

		// shrink
		if oldNodeNum > newNodeNum {
			azSplits := strings.Split(d.Get("availability_zone").(string), ",")
			if newNodeNum < len(azSplits) {
				return fmt.Errorf("the number of remaining nodes after scale-in" +
					" must be greater than or equal to the number of Azs")
			}
			return shrinkLogstashCluster(ctx, d, oldNodeNum-newNodeNum, cssV1Client)
		}
	}

	return nil
}

func extendLogstashCluster(ctx context.Context, d *schema.ResourceData, extendNodesize int, cssV1Client *cssv1.CssClient) error {
	opts := buildLogstashClusterV1ExtendClusterParameters(d, extendNodesize)
	_, err := cssV1Client.UpdateExtendInstanceStorage(opts)
	if err != nil {
		return fmt.Errorf("error extending CSS logstash cluster (%s) instance number failed: %s", d.Id(), err)
	}

	err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

func shrinkLogstashCluster(ctx context.Context, d *schema.ResourceData, shrinkNodesize int, cssV1Client *cssv1.CssClient) error {
	opts := buildLogstashClusterV1ShrinkClusterParameters(d, shrinkNodesize)
	_, err := cssV1Client.UpdateShrinkCluster(opts)
	if err != nil {
		return fmt.Errorf("error shrinking CSS logstash cluster (%s) instance number failed: %s", d.Id(), err)
	}

	err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

func buildLogstashClusterV1ExtendClusterParameters(d *schema.ResourceData, nodesize int) *model.UpdateExtendInstanceStorageRequest {
	return &model.UpdateExtendInstanceStorageRequest{
		ClusterId: d.Id(),
		Body: &model.RoleExtendReq{
			Grow: []model.RoleExtendGrowReq{
				{
					Type:     "lgs",
					Nodesize: int32(nodesize),
					Disksize: 0,
				},
			},
		},
	}
}

func buildLogstashClusterV1ShrinkClusterParameters(d *schema.ResourceData, nodesize int) *model.UpdateShrinkClusterRequest {
	return &model.UpdateShrinkClusterRequest{
		ClusterId: d.Id(),
		Body: &model.ShrinkClusterReq{
			Shrink: []model.ShrinkNodeReq{
				{
					Type:           "lgs",
					ReducedNodeNum: int32(nodesize),
				},
			},
		},
	}
}

func resourceLogstashClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	cssV1Client, err := conf.HcCssV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	_, err = cssV1Client.DeleteCluster(&model.DeleteClusterRequest{ClusterId: d.Id()})
	if err != nil {
		return diag.Errorf("error deleting CSS logstash cluster (%s): %s", d.Id(), err)
	}

	err = checkClusterDeleteResult(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("failed to check the result of deletion %s", err)
	}
	return nil
}
