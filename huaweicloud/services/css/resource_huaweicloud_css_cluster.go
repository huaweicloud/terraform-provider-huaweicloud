package css

import (
	"context"
	"fmt"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"
	"github.com/chnsz/golangsdk/openstack/css/v1/snapshots"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceCssCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCssClusterCreate,
		ReadContext:   resourceCssClusterRead,
		UpdateContext: resourceCssClusterUpdate,
		DeleteContext: resourceCssClusterDelete,
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
				ForceNew: true,
			},

			"engine_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "elasticsearch",
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"expect_node_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"security_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				ForceNew:  true,
			},

			"node_config": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"network_info": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"subnet_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"volume": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"backup_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"keep_days": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  7,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "snapshot",
						},
						"bucket": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"backup_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"agency": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"tags": common.TagsSchema(),

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
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
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCssClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssV1Client, err := config.CssV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating CSS V1 client: %s", err)
	}

	createClusterOpts, paramErr := buildClusterCreateParameters(d, config)
	if paramErr != nil {
		return diag.FromErr(paramErr)
	}

	r, createErr := cluster.Create(cssV1Client, *createClusterOpts)
	if createErr != nil {
		return diag.Errorf("Error creating CssClusterV1, err=%s", createErr)
	}

	clusterId := r.Cluster.Id

	createResultErr := checkClusterCreateResult(ctx, cssV1Client, clusterId, d.Timeout(schema.TimeoutCreate))
	if createResultErr != nil {
		return diag.FromErr(createResultErr)
	}
	d.SetId(clusterId)

	// enable snapshot function and set policy when "backup_strategy" was specified
	// createBackupErr := resourceCssClusterCreateBackupStrategy(d, cssV1Client)
	// if createBackupErr != nil {
	// 	return diag.FromErr(createBackupErr)
	// }

	return resourceCssClusterRead(ctx, d, meta)
}

func buildClusterCreateParameters(d *schema.ResourceData, config *config.Config) (*cluster.CreateOpts, error) {
	createClusterOpts := cluster.CreateOpts{
		Name: d.Get("name").(string),
		Datastore: &cluster.DatastoreBody{
			Type:    d.Get("engine_type").(string),
			Version: d.Get("engine_version").(string),
		},
		InstanceNum: d.Get("expect_node_num").(int),
		Instance: &cluster.InstanceBody{
			AvailabilityZone: d.Get("node_config.0.availability_zone").(string),
			FlavorRef:        d.Get("node_config.0.flavor").(string),
			Nics: cluster.InstanceNicsBody{
				NetId:           d.Get("node_config.0.network_info.0.subnet_id").(string),
				SecurityGroupId: d.Get("node_config.0.network_info.0.security_group_id").(string),
				VpcId:           d.Get("node_config.0.network_info.0.vpc_id").(string),
			},
			Volume: cluster.InstanceVolumeBody{
				Size:       d.Get("node_config.0.volume.0.size").(int),
				VolumeType: d.Get("node_config.0.volume.0.volume_type").(string),
			},
		},
		Tags:                utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		EnterpriseProjectId: config.GetEnterpriseProjectID(d),
	}
	//
	securityMode := d.Get("security_mode").(bool)
	if securityMode {
		adminPassword := d.Get("password").(string)
		if adminPassword == "" {
			return nil, fmtp.Errorf("administrator password is required in security mode")
		}
		createClusterOpts.HttpsEnable = true
		createClusterOpts.AuthorityEnable = true
		createClusterOpts.AdminPwd = adminPassword
	}

	//back_up strategy
	backupStrategy := resourceCssClusterCreateBackupStrategy(d)
	if backupStrategy != nil {
		createClusterOpts.BackupStrategy = backupStrategy
	}
	return &createClusterOpts, nil
}

func resourceCssClusterCreateBackupStrategy(d *schema.ResourceData) *cluster.BackupStrategyBody {
	backupRaw := d.Get("backup_strategy").([]interface{})
	if len(backupRaw) == 1 {
		raw := backupRaw[0].(map[string]interface{})
		opts := cluster.BackupStrategyBody{
			Prefix:   raw["prefix"].(string),
			Period:   raw["start_time"].(string),
			Keepday:  raw["keep_days"].(int),
			Bucket:   raw["bucket"].(string),
			BasePath: raw["backup_path"].(string),
			Agency:   raw["agency"].(string),
		}
		return &opts
	}
	return nil
}

func resourceCssClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssV1Client, err := config.CssV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating CSS V1 client: %s", err)
	}

	clusterDetail, createErr := cluster.Get(cssV1Client, d.Id())
	if createErr != nil {
		return diag.Errorf("Query cluster detail failed,cluster_id=%s,err=%s", d.Id(), createErr)
	}

	if err := setCssClusterProperties(d, cssV1Client, clusterDetail); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setCssClusterProperties(d *schema.ResourceData, client *golangsdk.ServiceClient,
	clusterDetail *cluster.ClusterDetailResponse) error {
	mErr := multierror.Append(
		d.Set("created", clusterDetail.Created),
		d.Set("endpoint", clusterDetail.Endpoint),
		d.Set("engine_type", clusterDetail.Datastore.Type),
		d.Set("engine_version", clusterDetail.Datastore.Version),
		d.Set("expect_node_num", len(clusterDetail.Instances)),
		d.Set("enterprise_project_id", clusterDetail.EnterpriseProjectId),
		d.Set("name", clusterDetail.Name),
		d.Set("status", clusterDetail.Status),
		setClusterNodes(d, clusterDetail.Instances),
		setClusterSecurity(d, clusterDetail),
		setClusterBackupStrategy(d, client, clusterDetail),
		utils.SetResourceTagsToState(d, client, "css-cluster"),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("Error setting vault fields: %s", err)
	}

	return nil
}

func setClusterNodes(d *schema.ResourceData, instances []cluster.ClusterDetailInstances) error {
	var result []interface{}

	for i := 0; i < len(instances); i++ {
		instance := instances[i]
		r := make(map[string]interface{})
		r["id"] = instance.Id
		r["name"] = instance.Name
		r["type"] = instance.Type
		result = append(result, r)
	}
	return d.Set("nodes", result)
}

func setClusterSecurity(d *schema.ResourceData, clusterDetail *cluster.ClusterDetailResponse) error {
	authorityEnable := clusterDetail.AuthorityEnable
	if authorityEnable {
		return d.Set("security_mode", true)
	}
	return nil
}

func setClusterBackupStrategy(d *schema.ResourceData, client *golangsdk.ServiceClient,
	clusterDetail *cluster.ClusterDetailResponse) error {
	// set backup strategy property
	policy, err := snapshots.PolicyGet(client, d.Id()).Extract()
	if err != nil {
		return fmtp.Errorf("Error extracting Cluster:backup_strategy, err: %s", err)
	}

	var strategy []map[string]interface{}

	if policy.Enable == "true" {
		strategy = []map[string]interface{}{
			{
				"prefix":      policy.Prefix,
				"start_time":  policy.Period,
				"keep_days":   policy.KeepDay,
				"bucket":      policy.Bucket,
				"backup_path": policy.BasePath,
				"agency":      policy.Agency,
			},
		}
	} else {
		strategy = nil
	}
	return d.Set("backup_strategy", strategy)
}

func resourceCssClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssV1Client, err := config.CssV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating CSS V1 client: %s", err)
	}

	clusterId := d.Id()

	//extend cluster
	opts, err := buildCssClusterV1ExtendClusterParameters(d)
	if err != nil {
		return diag.Errorf("Error building the request body of api(extend_cluster), err=%s", err)
	}
	if opts != nil {
		_, extendErr := cluster.ExtendInstanceStorage(cssV1Client, clusterId, *opts)
		if extendErr != nil {
			return diag.Errorf("Extend CSS cluster instance storage failed.cluster_id=%s,error=%s", clusterId,
				extendErr)
		}
		checkExtendErr := checkClusterExtendResult(ctx, cssV1Client, clusterId, d.Timeout(schema.TimeoutUpdate))
		if checkExtendErr != nil {
			return diag.FromErr(checkExtendErr)
		}
	}

	// update backup strategy
	if d.HasChange("backup_strategy") {
		var opts snapshots.PolicyCreateOpts

		value, ok := d.GetOk("backup_strategy")
		if !ok {
			opts = snapshots.PolicyCreateOpts{
				Prefix:  "snapshot",
				Period:  "00:00 GMT+08:00",
				KeepDay: 7,
				Enable:  "false",
			}

			errPolicy := snapshots.PolicyCreate(cssV1Client, &opts, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("Error updating backup strategy: %s", errPolicy)
			}
		} else {
			rawList := value.([]interface{})
			if len(rawList) == 1 {
				raw := rawList[0].(map[string]interface{})

				if d.HasChanges("backup_strategy.0.bucket", "backup_strategy.0.backup_path",
					"backup_strategy.0.agency") {
					// If obs is specified, update basic configurations
					obsOpts := snapshots.UpdateSnapshotSettingReq{
						Bucket:   raw["bucket"].(string),
						BasePath: raw["backup_path"].(string),
						Agency:   raw["agency"].(string),
					}
					_, err = snapshots.UpdateSnapshotSetting(cssV1Client, clusterId, obsOpts)
					if err != nil {
						return diag.Errorf("error Modifying Basic Configurations of a Cluster Snapshot: %s", err)
					}
				}

				// check backup strategy, if the policy was disabled, we should enable it
				policy, err := snapshots.PolicyGet(cssV1Client, clusterId).Extract()
				if err != nil {
					return diag.Errorf("Error extracting Cluster backup_strategy, err: %s", err)
				}

				if policy.Enable == "false" && raw["bucket"] == nil {
					// If obs is not specified,  create  basic configurations automatically
					err = snapshots.Enable(cssV1Client, d.Id()).ExtractErr()
					if err != nil {
						return diag.Errorf("Error enable snapshot function: %s", err)
					}
				}
				// update policy
				if d.HasChanges("backup_strategy.0.prefix", "backup_strategy.0.start_time",
					"backup_strategy.0.keep_days") {
					opts = snapshots.PolicyCreateOpts{
						Prefix:  raw["prefix"].(string),
						Period:  raw["start_time"].(string),
						KeepDay: raw["keep_days"].(int),
						Enable:  "true",
					}

					errPolicy := snapshots.PolicyCreate(cssV1Client, &opts, d.Id()).ExtractErr()
					if err != nil {
						return diag.Errorf("Error updating backup strategy: %s", errPolicy)
					}
				}

			}
		}

	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(cssV1Client, d, "css-cluster", clusterId)
		if tagErr != nil {
			return diag.Errorf("Error updating tags of CSS cluster:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceCssClusterRead(ctx, d, meta)
}

func buildCssClusterV1ExtendClusterParameters(rd *schema.ResourceData) (*cluster.RoleExtendReq, error) {
	var extendReq cluster.RoleExtendGrowReq

	oldv, newv := rd.GetChange("expect_node_num")
	nodesize := newv.(int) - oldv.(int)
	if nodesize < 0 {
		return nil, fmtp.Errorf("expect_node_num only supports to be extended")
	}

	extendReq.Nodesize = &nodesize

	//volume size location: reference to the Schema of css_cluster_v1
	oldDisksize, newDisksize := rd.GetChange("node_config.0.volume.0.size")
	disksize := newDisksize.(int) - oldDisksize.(int)
	if disksize < 0 {
		return nil, fmtp.Errorf("volume size only supports to be extended")
	}
	extendReq.Disksize = &disksize

	// both of nodesize and disksize can not be set to 0 simultaneously
	if nodesize == 0 && disksize == 0 {
		return nil, nil
	}

	extendReq.Type = "ess"

	return &cluster.RoleExtendReq{Grow: []cluster.RoleExtendGrowReq{extendReq}}, nil

}

func resourceCssClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssV1Client, err := config.CssV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating CSS V1 client: %s", err)
	}

	clusterId := d.Id()
	errResult := cluster.Delete(cssV1Client, clusterId)
	if errResult.Err != nil {
		return diag.Errorf("Delete CSS Cluster failed. %s", errResult.Err)
	}

	errCheckRt := checkClusterDeleteResult(ctx, cssV1Client, clusterId, d.Timeout(schema.TimeoutDelete))
	if errCheckRt != nil {
		return diag.Errorf("Failed to check the result of deletion %s", errCheckRt)
	}
	d.SetId("")
	return nil
}

func checkClusterCreateResult(ctx context.Context, cssV1Client *golangsdk.ServiceClient, clusterId string,
	timeout time.Duration) error {
	createStateConf := &resource.StateChangeConf{
		Pending: []string{cluster.ClusterStatusInProcess},
		Target:  []string{cluster.ClusterStatusAvailable},
		Refresh: func() (interface{}, string, error) {
			resp, err := cluster.Get(cssV1Client, clusterId)
			if err != nil {
				return nil, "failed", err
			}
			if resp.FailedReasons.ErrorCode != "" {
				return nil, "failed", fmtp.Errorf("error_code: %s, error_msg: %s", resp.FailedReasons.ErrorCode,
					resp.FailedReasons.ErrorMsg)
			}
			return resp, resp.Status, err
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CSS (%s) to be created: %s", clusterId, err)
	}
	return nil
}

func checkClusterDeleteResult(ctx context.Context, cssV1Client *golangsdk.ServiceClient, clusterId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			_, err := cluster.Get(cssV1Client, clusterId)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return true, "Done", nil
				}
				return nil, "", nil
			}
			return true, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CSS (%s) to be delete: %s", clusterId, err)
	}
	return nil
}

func checkClusterExtendResult(ctx context.Context, cssV1Client *golangsdk.ServiceClient, clusterId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			resp, err := cluster.Get(cssV1Client, clusterId)
			if err != nil {
				return nil, "failed", err
			}

			if resp.FailedReasons.ErrorCode != "" {
				return nil, "failed", fmtp.Errorf("error_code: %s, error_msg: %s", resp.FailedReasons.ErrorCode,
					resp.FailedReasons.ErrorMsg)
			}
			if checkCssClusterExtendResp(resp) {
				return resp, "Done", nil
			}
			return resp, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CSS (%s) to be extend: %s", clusterId, err)
	}
	return nil
}

func checkCssClusterExtendResp(detail *cluster.ClusterDetailResponse) bool {
	//actions --- the behaviors on a cluster
	if len(detail.Actions) > 0 {
		return false
	}
	if len(detail.ActionProgress) > 0 {
		return false
	}
	if len(detail.Instances) == 0 {
		return false
	}
	for _, v := range detail.Instances {
		status := v.Status
		if status != cluster.ClusterStatusAvailable {
			return false
		}
	}
	return true
}
