package gaussdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bss/v2/orders"
	"github.com/chnsz/golangsdk/openstack/opengauss/v3/backups"
	"github.com/chnsz/golangsdk/openstack/opengauss/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type HaMode string

type ConsistencyType string

const (
	HaModeDistributed HaMode = "enterprise"
	HAModeCentralized HaMode = "centralization_standard"

	ConsistencyTypeStrong   ConsistencyType = "strong"
	ConsistencyTypeEventual ConsistencyType = "eventual"
)

// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB POST /v3/{project_id}/instances
// @API GaussDB PUT /v3/{project_id}/instances/{instance_id}/name
// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/password
// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/action
// @API GaussDB PUT /v3/{project_id}/instances/{instance_id}/backups/policy
// @API GaussDB DELETE /v3/{project_id}/instances/{instance_id}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceOpenGaussInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussInstanceCreate,
		ReadContext:   resourceOpenGaussInstanceRead,
		UpdateContext: resourceOpenGaussInstanceUpdate,
		DeleteContext: resourceOpenGaussInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, v interface{}) error {
			if d.HasChange("coordinator_num") {
				return d.SetNewComputed("private_ips")
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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
			"ha": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: utils.SuppressCaseDiffs,
						},
						"replication_mode": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"consistency": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(ConsistencyTypeStrong), string(ConsistencyTypeEventual),
							}, true),
							DiffSuppressFunc: utils.SuppressCaseDiffs,
						},
					},
				},
			},
			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"sharding_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"coordinator_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"replica_num": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					2, 3,
				}),
				Default: 3,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "UTC+08:00",
			},
			"datastore": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"backup_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
						},
					},
				},
			},
			"force_import": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),

			// Attributes
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"db_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_strategy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_window": {
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
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceOpenGaussDataStore(d *schema.ResourceData) instances.DataStoreOpt {
	var db instances.DataStoreOpt

	datastoreRaw := d.Get("datastore").([]interface{})
	if len(datastoreRaw) == 1 {
		datastore := datastoreRaw[0].(map[string]interface{})
		db.Type = datastore["engine"].(string)
		db.Version = datastore["version"].(string)
	} else {
		db.Type = "GaussDB(for openGauss)"
	}
	return db
}

func resourceOpenGaussBackupStrategy(d *schema.ResourceData) *instances.BackupStrategyOpt {
	var backupOpt instances.BackupStrategyOpt

	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	if len(backupStrategyRaw) == 1 {
		strategy := backupStrategyRaw[0].(map[string]interface{})
		backupOpt.StartTime = strategy["start_time"].(string)
		backupOpt.KeepDays = strategy["keep_days"].(int)
		return &backupOpt
	}

	return nil
}

func OpenGaussInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.GetInstanceByID(client, instanceID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}

		return v, v.Status, nil
	}
}

func buildOpenGaussInstanceCreateOpts(d *schema.ResourceData,
	cfg *config.Config) (instances.CreateGaussDBOpts, error) {
	createOpts := instances.CreateGaussDBOpts{
		Name:                d.Get("name").(string),
		Flavor:              d.Get("flavor").(string),
		Region:              cfg.GetRegion(d),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		Port:                d.Get("port").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		TimeZone:            d.Get("time_zone").(string),
		AvailabilityZone:    d.Get("availability_zone").(string),
		ConfigurationId:     d.Get("configuration_id").(string),
		ShardingNum:         d.Get("sharding_num").(int),
		CoordinatorNum:      d.Get("coordinator_num").(int),
		ReplicaNum:          d.Get("replica_num").(int),
		DataStore:           resourceOpenGaussDataStore(d),
		BackupStrategy:      resourceOpenGaussBackupStrategy(d),
	}

	var dnNum = 1
	haRaw := d.Get("ha").([]interface{})
	log.Printf("[DEBUG] The HA structure is: %#v", haRaw)
	ha := haRaw[0].(map[string]interface{})
	mode := ha["mode"].(string)
	createOpts.Ha = &instances.HaOpt{
		Mode:            mode,
		ReplicationMode: ha["replication_mode"].(string),
		Consistency:     ha["consistency"].(string),
	}
	if mode == string(HaModeDistributed) {
		dnNum = d.Get("sharding_num").(int)
	}
	if mode == string(HAModeCentralized) {
		dnNum = d.Get("replica_num").(int) + 1
	}

	volumeRaw := d.Get("volume").([]interface{})
	if len(volumeRaw) > 0 {
		log.Printf("[DEBUG] The volume structure is: %#v", volumeRaw)
		volume := volumeRaw[0].(map[string]interface{})
		dnSize := volume["size"].(int)
		volumeSize := dnSize * dnNum
		createOpts.Volume = instances.VolumeOpt{
			Type: volume["type"].(string),
			Size: volumeSize,
		}
	}
	log.Printf("[DEBUG] The createOpts object is: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	if d.Get("charging_mode").(string) == "prePaid" {
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return createOpts, err
		}
		createOpts.ChargeInfo = &instances.ChargeInfo{
			ChargeMode:  "prePaid",
			PeriodType:  d.Get("period_unit").(string),
			PeriodNum:   d.Get("period").(int),
			IsAutoRenew: d.Get("auto_renew").(string),
			IsAutoPay:   "true",
		}
	}
	return createOpts, nil
}

func resourceOpenGaussInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.OpenGaussV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB v3 client: %s ", err)
	}

	// If force_import set, try to import it instead of creating
	if common.HasFilledOpt(d, "force_import") {
		log.Printf("[DEBUG] the Gaussdb opengauss instance force_import is set, try to import it instead of creating")
		listOpts := instances.ListGaussDBInstanceOpts{
			Name: d.Get("name").(string),
		}
		pages, err := instances.List(client, listOpts).AllPages()
		if err != nil {
			return diag.FromErr(err)
		}

		allInstances, err := instances.ExtractGaussDBInstances(pages)
		if err != nil {
			return diag.Errorf("unable to retrieve instances: %s", err)
		}
		if allInstances.TotalCount > 0 {
			instance := allInstances.Instances[0]
			log.Printf("[DEBUG] found existing opengauss instance %s with name %s", instance.Id, instance.Name)
			d.SetId(instance.Id)
			return resourceOpenGaussInstanceRead(ctx, d, meta)
		}
	}

	createOpts, err := buildOpenGaussInstanceCreateOpts(d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating OpenGauss instance: %s", err)
	}

	if resp.OrderId != "" {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, resp.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, resp.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resourceId)
	} else {
		d.SetId(resp.Instance.Id)
	}

	// waiting for the instance to become ready
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"BUILD", "BACKING UP"},
		Target:                    []string{"ACTIVE"},
		Refresh:                   OpenGaussInstanceStateRefreshFunc(client, d.Id()),
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     20 * time.Second,
		PollInterval:              20 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to become ready: %s", d.Id(), err)
	}

	// This is a workaround to avoid db connection issue
	time.Sleep(360 * time.Second) // lintignore:R018

	return resourceOpenGaussInstanceRead(ctx, d, meta)
}

func flattenOpenGaussDataStore(dataStore instances.DataStoreOpt) []map[string]interface{} {
	if dataStore == (instances.DataStoreOpt{}) {
		return nil
	}
	return []map[string]interface{}{
		{
			"version": dataStore.Version,
			"engine":  dataStore.Type,
		},
	}
}

func flattenOpenGaussBackupStrategy(backupStrategy instances.BackupStrategyOpt) []map[string]interface{} {
	if backupStrategy == (instances.BackupStrategyOpt{}) {
		return nil
	}
	return []map[string]interface{}{
		{
			"start_time": backupStrategy.StartTime,
			"keep_days":  backupStrategy.KeepDays,
		},
	}
}

func flattenOpenGaussVolume(volume instances.VolumeOpt, dnNum int) []map[string]interface{} {
	if volume == (instances.VolumeOpt{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"type": volume.Type,
			"size": volume.Size / dnNum,
		},
	}
}

func setOpenGaussNodesAndRelatedNumbers(d *schema.ResourceData, instance instances.GaussDBInstance,
	dnNum *int) error {
	var (
		shardingNum    = 0
		coordinatorNum = 0
	)

	nodesList := make([]map[string]interface{}, 0, 1)
	for _, raw := range instance.Nodes {
		node := map[string]interface{}{
			"id":                raw.Id,
			"name":              raw.Name,
			"status":            raw.Status,
			"role":              raw.Role,
			"availability_zone": raw.AvailabilityZone,
		}
		nodesList = append(nodesList, node)

		if strings.Contains(raw.Name, "_gaussdbv5cn") {
			coordinatorNum++
		} else if strings.Contains(raw.Name, "_gaussdbv5dn") {
			shardingNum++
		}
	}

	if shardingNum > 0 && coordinatorNum > 0 {
		*dnNum = shardingNum / d.Get("replica_num").(int)
		return multierror.Append(nil,
			d.Set("nodes", nodesList),
			d.Set("sharding_num", dnNum),
			d.Set("coordinator_num", coordinatorNum),
		).ErrorOrNil()
	}
	// If the HA mode is centralized, the HA structure of API response is nil.
	*dnNum = instance.ReplicaNum + 1
	return multierror.Append(nil,
		d.Set("nodes", nodesList),
		d.Set("replica_num", instance.ReplicaNum),
	).ErrorOrNil()
}

func setOpenGaussPrivateIpsAndEndpoints(d *schema.ResourceData, privateIps []string, port int) error {
	if len(privateIps) < 1 {
		return nil
	}

	privateIp := privateIps[0]
	ipList := strings.Split(privateIp, "/")
	endpoints := []string{}
	for i := 0; i < len(ipList); i++ {
		ipList[i] = strings.Trim(ipList[i], " ")
		endpoint := fmt.Sprintf("%s:%d", ipList[i], port)
		endpoints = append(endpoints, endpoint)
	}

	return multierror.Append(nil,
		d.Set("private_ips", ipList),
		d.Set("endpoints", endpoints),
	).ErrorOrNil()
}

func resourceOpenGaussInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.OpenGaussV3Client(region)
	if err != nil {
		return diag.Errorf("error creating GaussDB v3 client: %s ", err)
	}

	instanceID := d.Id()
	instance, err := instances.GetInstanceByID(client, instanceID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenGauss instance")
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}

	var dnNum = 1
	log.Printf("[DEBUG] retrieved instance (%s): %#v", instanceID, instance)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", instance.Name),
		d.Set("status", instance.Status),
		d.Set("type", instance.Type),
		d.Set("vpc_id", instance.VpcId),
		d.Set("subnet_id", instance.SubnetId),
		d.Set("security_group_id", instance.SecurityGroupId),
		d.Set("db_user_name", instance.DbUserName),
		d.Set("time_zone", instance.TimeZone),
		d.Set("flavor", instance.FlavorRef),
		d.Set("port", strconv.Itoa(instance.Port)),
		d.Set("switch_strategy", instance.SwitchStrategy),
		d.Set("maintenance_window", instance.MaintenanceWindow),
		d.Set("public_ips", instance.PublicIps),
		d.Set("charging_mode", instance.ChargeInfo.ChargeMode),
		d.Set("datastore", flattenOpenGaussDataStore(instance.DataStore)),
		d.Set("backup_strategy", flattenOpenGaussBackupStrategy(instance.BackupStrategy)),
		setOpenGaussNodesAndRelatedNumbers(d, instance, &dnNum),
		d.Set("volume", flattenOpenGaussVolume(instance.Volume, dnNum)),
		setOpenGaussPrivateIpsAndEndpoints(d, instance.PrivateIps, instance.Port),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting OpenGauss instance fields: %s", mErr.ErrorOrNil())
	}
	return nil
}

func expandOpenGaussShardingNumber(ctx context.Context, cfg *config.Config, client *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	old, newnum := d.GetChange("sharding_num")
	if newnum.(int) < old.(int) {
		return fmt.Errorf("error expanding shard for instance: new num must be larger than the old one")
	}
	expandSize := newnum.(int) - old.(int)
	opts := instances.UpdateOpts{
		ExpandCluster: &instances.UpdateClusterOpts{
			Shard: &instances.Shard{
				Count: expandSize,
			},
		},
		IsAutoPay: "true",
	}
	log.Printf("[DEBUG] the updateOpts object of sharding number is: %#v", opts)
	return updateVolumeAndRelatedHaNumbers(ctx, cfg, client, d, opts)
}

func expandOpenGaussCoordinatorNumber(ctx context.Context, cfg *config.Config, client *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	old, newnum := d.GetChange("coordinator_num")
	if newnum.(int) < old.(int) {
		return fmt.Errorf("error expanding coordinator for instance: new number must be larger than the old one")
	}
	expandSize := newnum.(int) - old.(int)

	var coordinators []instances.Coordinator
	azlist := strings.Split(d.Get("availability_zone").(string), ",")
	for i := 0; i < expandSize; i++ {
		coordinator := instances.Coordinator{
			AzCode: azlist[0],
		}
		coordinators = append(coordinators, coordinator)
	}
	opts := instances.UpdateOpts{
		ExpandCluster: &instances.UpdateClusterOpts{
			Coordinators: coordinators,
		},
		IsAutoPay: "true",
	}
	log.Printf("[DEBUG] the updateOpts object of coordinator number is: %#v", opts)
	return updateVolumeAndRelatedHaNumbers(ctx, cfg, client, d, opts)
}

func updateOpenGaussVolumeSize(ctx context.Context, cfg *config.Config, client *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	volumeRaw := d.Get("volume").([]interface{})
	dnSize := volumeRaw[0].(map[string]interface{})["size"].(int)
	dnNum := 1
	if d.Get("ha.0.mode").(string) == string(HaModeDistributed) {
		dnNum = d.Get("sharding_num").(int)
	}
	if d.Get("ha.0.mode").(string) == string(HAModeCentralized) {
		dnNum = d.Get("replica_num").(int) + 1
	}
	opts := instances.UpdateOpts{
		EnlargeVolume: &instances.UpdateVolumeOpts{
			Size: dnSize * dnNum,
		},
		IsAutoPay: "true",
	}
	log.Printf("[DEBUG] the updateOpts object of volume size is: %#v", opts)
	return updateVolumeAndRelatedHaNumbers(ctx, cfg, client, d, opts)
}

func updateVolumeAndRelatedHaNumbers(ctx context.Context, cfg *config.Config, client *golangsdk.ServiceClient,
	d *schema.ResourceData, opts instances.UpdateOpts) error {
	instanceId := d.Id()
	resp, err := instances.Update(client, instanceId, opts)
	if err != nil {
		return fmt.Errorf("error updating instance (%s): %s", instanceId, err)
	}
	if resp.OrderId != "" {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}
		if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), resp.OrderId); err != nil {
			return err
		}
	}
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"MODIFYING", "EXPANDING", "BACKING UP"},
		Target:                    []string{"ACTIVE"},
		Refresh:                   OpenGaussInstanceStateRefreshFunc(client, instanceId),
		Timeout:                   d.Timeout(schema.TimeoutUpdate),
		Delay:                     20 * time.Second,
		PollInterval:              20 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) status to active: %s ", instanceId, err)
	}

	return nil
}

func resourceOpenGaussInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.OpenGaussV3Client(region)
	if err != nil {
		return diag.Errorf("error creating GaussDB v3 client: %s ", err)
	}

	instanceId := d.Id()
	log.Printf("[DEBUG] updating OpenGaussDB instances %s", instanceId)

	if d.HasChange("name") {
		renameOpts := instances.RenameOpts{
			Name: d.Get("name").(string),
		}
		_, err = instances.Rename(client, renameOpts, instanceId).Extract()
		if err != nil {
			return diag.Errorf("error updating name for instance (%s): %s", instanceId, err)
		}
	}

	if d.HasChange("password") {
		restorePasswordOpts := instances.RestorePasswordOpts{
			Password: d.Get("password").(string),
		}
		r := golangsdk.ErrResult{}
		r.Result = instances.RestorePassword(client, restorePasswordOpts, instanceId)
		if r.ExtractErr() != nil {
			return diag.Errorf("error updating password for instance (%s): %s ", instanceId, r.Err)
		}
	}

	if d.HasChange("sharding_num") {
		if err := expandOpenGaussShardingNumber(ctx, cfg, client, d); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("coordinator_num") {
		if err := expandOpenGaussCoordinatorNumber(ctx, cfg, client, d); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("volume") {
		if err := updateOpenGaussVolumeSize(ctx, cfg, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("backup_strategy") {
		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keepDays := rawMap["keep_days"].(int)

		updateOpts := backups.UpdateOpts{
			KeepDays:           &keepDays,
			StartTime:          rawMap["start_time"].(string),
			Period:             "1,2,3,4,5,6,7", // Fixed to "1,2,3,4,5,6,7"
			DifferentialPeriod: "30",            // Fixed to "30"
		}

		log.Printf("[DEBUG] the updateOpts object of backup_strategy parameter is: %#v", updateOpts)
		err = backups.Update(client, instanceId, updateOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating backup_strategy: %s", err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), instanceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", instanceId, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "gaussdb",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceOpenGaussInstanceRead(ctx, d, meta)
}

func resourceOpenGaussInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.OpenGaussV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB v3 client: %s ", err)
	}

	instanceId := d.Id()
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{instanceId}); err != nil {
			return diag.Errorf("error unsubscribe OpenGauss instance: %s", err)
		}
	} else {
		result := instances.Delete(client, instanceId)
		if result.Err != nil {
			return common.CheckDeletedDiag(d, result.Err, "OpenGauss instance")
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:        []string{"ACTIVE", "BACKING UP", "FAILED"},
		Target:         []string{"DELETED"},
		Refresh:        OpenGaussInstanceStateRefreshFunc(client, instanceId),
		Timeout:        d.Timeout(schema.TimeoutDelete),
		Delay:          60 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 2,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to be deleted: %s", instanceId, err)
	}
	log.Printf("[DEBUG] instance deleted successfully %s", instanceId)
	return nil
}
