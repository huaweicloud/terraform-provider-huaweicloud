package dds

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dds/v3/instances"
	"github.com/chnsz/golangsdk/openstack/dds/v3/jobs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ctxType string

// @API DDS POST /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/instances
// @API DDS POST /v3/{project_id}/instances/{id}/tags/action
// @API DDS GET /v3/{project_id}/instances/{id}/tags
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/modify-name
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/reset-password
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/modify-security-group
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/switch-ssl
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/modify-port
// @API DDS POST /v3/{project_id}/instances/{instance_id}/enlarge-volume
// @API DDS POST /v3/{project_id}/instances/{instance_id}/enlarge
// @API DDS POST /v3/{project_id}/instances/{instance_id}/resize
// @API DDS GET /v3/{project_id}/jobs
// @API DDS DELETE /v3/{project_id}/instances/{serverID}
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/remark
// @API DDS POST /v3/{project_id}/instances/{instance_id}/migrate
// @API DDS GET /v3/{project_id}/instances/{instance_id}/backups/policy
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/backups/policy
// @API DDS GET /v3/{project_id}/instances/{instance_id}/monitoring-by-seconds/switch
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/monitoring-by-seconds/switch
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/replica-set/name
// @API DDS GET /v3/{project_id}/instances/{instance_id}/replica-set/name
// @API DDS PUT /v3/{project_id}/configurations/{config_id}/apply
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/slowlog-desensitization/{status}
// @API DDS GET /v3/{project_id}/instances/{instance_id}/slowlog-desensitization/status
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/balancer/active-window
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/balancer/{action}
// @API DDS GET /v3/{project_id}/instances/{instance_id}/balancer
// @API DDS POST /v3/{project_id}/instances/{instance_id}/replicaset-node
// @API DDS POST /v3/{project_id}/instances/{instance_id}/client-network
// @API DDS GET /v3/{project_id}/instances/{instance_id}/client-network
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/maintenance-window
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
func ResourceDdsInstanceV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsInstanceV3Create,
		ReadContext:   resourceDdsInstanceV3Read,
		UpdateContext: resourceDdsInstanceV3Update,
		DeleteContext: resourceDdsInstanceV3Delete,
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
			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"storage_engine": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Required:         true,
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
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"disk_encryption_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"flavor": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"num": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"storage": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"spec_code": {
							Type:     schema.TypeString,
							Required: true,
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
							Required: true,
						},
						"period": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"balancer_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"balancer_active_begin": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"balancer_active_end"},
			},
			"balancer_active_end": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"balancer_active_begin"},
			},
			"maintain_begin": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"maintain_end"},
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"second_level_monitoring_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replica_set_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"slow_log_desensitization": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"client_network_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenew(nil),
			"auto_pay":      common.SchemaAutoPay(nil),
			"tags":          common.TagsSchema(),
			"db_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodes": {
							Type:     schema.TypeList,
							Elem:     ddsInstanceInstanceNodeSchema(),
							Computed: true,
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// deprecated
			"nodes": {
				Type:        schema.TypeList,
				Elem:        ddsInstanceInstanceNodeSchema(),
				Computed:    true,
				Description: `This field is deprecated.`,
			},
		},
	}
}

func resourceDdsDataStore(d *schema.ResourceData) instances.DataStore {
	var dataStore instances.DataStore
	datastoreRaw := d.Get("datastore").([]interface{})
	log.Printf("[DEBUG] datastoreRaw: %+v", datastoreRaw)
	if len(datastoreRaw) == 1 {
		dataStore.Type = datastoreRaw[0].(map[string]interface{})["type"].(string)
		dataStore.Version = datastoreRaw[0].(map[string]interface{})["version"].(string)
		dataStore.StorageEngine = datastoreRaw[0].(map[string]interface{})["storage_engine"].(string)
	}
	log.Printf("[DEBUG] datastore: %+v", dataStore)
	return dataStore
}

func resourceDdsConfiguration(d *schema.ResourceData) []instances.Configuration {
	var configurations []instances.Configuration
	configurationRaw := d.Get("configuration").([]interface{})
	log.Printf("[DEBUG] configurationRaw: %+v", configurationRaw)
	for i := range configurationRaw {
		configuration := configurationRaw[i].(map[string]interface{})
		flavorReq := instances.Configuration{
			Type: configuration["type"].(string),
			Id:   configuration["id"].(string),
		}
		configurations = append(configurations, flavorReq)
	}
	log.Printf("[DEBUG] configurations: %+v", configurations)
	return configurations
}

func resourceDdsFlavors(d *schema.ResourceData) []instances.Flavor {
	var flavors []instances.Flavor
	flavorRaw := d.Get("flavor").([]interface{})
	log.Printf("[DEBUG] flavorRaw: %+v", flavorRaw)
	for i := range flavorRaw {
		flavor := flavorRaw[i].(map[string]interface{})
		flavorReq := instances.Flavor{
			Type:     flavor["type"].(string),
			Num:      flavor["num"].(int),
			Storage:  flavor["storage"].(string),
			Size:     flavor["size"].(int),
			SpecCode: flavor["spec_code"].(string),
		}
		flavors = append(flavors, flavorReq)
	}
	log.Printf("[DEBUG] flavors: %+v", flavors)
	return flavors
}

func resourceDdsBackupStrategy(d *schema.ResourceData) instances.BackupStrategy {
	var backupStrategy instances.BackupStrategy
	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	startTime := "00:00-01:00"
	keepDays := 7
	period := "1,2,3,4,5,6,7"
	if len(backupStrategyRaw) == 1 {
		startTime = backupStrategyRaw[0].(map[string]interface{})["start_time"].(string)
		keepDays = backupStrategyRaw[0].(map[string]interface{})["keep_days"].(int)
		if periodRaw := backupStrategyRaw[0].(map[string]interface{})["period"].(string); periodRaw != "" {
			period = periodRaw
		}
	}
	backupStrategy.KeepDays = &keepDays
	backupStrategy.StartTime = startTime
	backupStrategy.Period = period
	return backupStrategy
}

func ddsInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := instances.ListInstanceOpts{
			Id: instanceID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return nil, "", err
		}
		instancesList, err := instances.ExtractInstances(allPages)
		if err != nil {
			return nil, "", err
		}

		if instancesList.TotalCount == 0 {
			var instance instances.InstanceResponse
			return instance, "deleted", nil
		}
		insts := instancesList.Instances

		status := insts[0].Status
		// wait for updating
		if status == "normal" && len(insts[0].Actions) > 0 {
			status = "updating"
		}
		return insts[0], status, nil
	}
}

func buildChargeInfoParams(d *schema.ResourceData) instances.ChargeInfo {
	chargeInfo := instances.ChargeInfo{
		ChargeMode: d.Get("charging_mode").(string),
		PeriodType: d.Get("period_unit").(string),
		PeriodNum:  d.Get("period").(int),
	}
	if d.Get("auto_pay").(string) != "false" {
		chargeInfo.IsAutoPay = true
	}
	if d.Get("auto_renew").(string) == "true" {
		chargeInfo.IsAutoRenew = true
	}
	return chargeInfo
}

func resourceDdsInstanceV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating DDS client: %s ", err)
	}

	createOpts := instances.CreateOpts{
		Name:                d.Get("name").(string),
		DataStore:           resourceDdsDataStore(d),
		Region:              conf.GetRegion(d),
		AvailabilityZone:    d.Get("availability_zone").(string),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		DiskEncryptionId:    d.Get("disk_encryption_id").(string),
		Mode:                d.Get("mode").(string),
		Configuration:       resourceDdsConfiguration(d),
		Flavor:              resourceDdsFlavors(d),
		BackupStrategy:      resourceDdsBackupStrategy(d),
		EnterpriseProjectID: conf.GetEnterpriseProjectID(d),
	}
	if d.Get("ssl").(bool) {
		createOpts.Ssl = "1"
	} else {
		createOpts.Ssl = "0"
	}
	if d.Get("charging_mode").(string) == "prePaid" {
		chargeInfo := buildChargeInfoParams(d)
		createOpts.ChargeInfo = &chargeInfo
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	if val, ok := d.GetOk("port"); ok {
		createOpts.Port = strconv.Itoa(val.(int))
	}

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error getting instance from result: %s ", err)
	}
	log.Printf("[DEBUG] Create : instance %s: %#v", instance.Id, instance)

	if instance.OrderId != "" {
		bssClient, err := conf.BssV2Client(conf.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, instance.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, instance.OrderId,
			d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resourceId)
	} else {
		d.SetId(instance.Id)
	}

	if description, ok := d.GetOk("description"); ok {
		opt := instances.RemarkOpts{
			Remark: description.(string),
		}
		err = instances.UpdateRemark(client, d.Id(), opt)
		if err != nil {
			return diag.Errorf("error adding description of the DDS instance (%s) : %s ", d.Id(), err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "updating"},
		Target:     []string{"normal"},
		Refresh:    ddsInstanceStateRefreshFunc(client, instance.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      120 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for instance (%s) to become ready: %s ",
			instance.Id, err)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "instances", instance.Id, taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("Error setting tags of DDS instance %s: %s", instance.Id, tagErr)
		}
	}

	// since the POST method has no `period`, update backup strategy for it
	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	if len(backupStrategyRaw) == 1 {
		period := backupStrategyRaw[0].(map[string]interface{})["period"].(string)
		if period != "" && !isEqualPeriod(period, "1,2,3,4,5,6,7") {
			if err := createBackupStrategy(ctx, client, d); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if ranges, ok := d.GetOk("client_network_ranges"); ok {
		err = updateClientNetworkRanges(ctx, client, d.Timeout(schema.TimeoutCreate), instance.Id, ranges)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if secondLevelMonitoringEnabled := d.Get("second_level_monitoring_enabled").(bool); secondLevelMonitoringEnabled {
		err = UpdateSecondsLevelMonitoring(ctx, client, d.Timeout(schema.TimeoutCreate), instance.Id,
			secondLevelMonitoringEnabled)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if slowLogDesensitization := d.Get("slow_log_desensitization").(string); slowLogDesensitization == "off" {
		err = UpdateSlowLogStatus(ctx, client, d.Timeout(schema.TimeoutCreate), instance.Id, slowLogDesensitization)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if status, ok := d.GetOk("balancer_status"); ok && status == "stop" {
		err = updateBalancerStatus(ctx, client, d.Timeout(schema.TimeoutCreate), instance.Id, status.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if begin, ok := d.GetOk("balancer_active_begin"); ok {
		err = updateBalancerActiveWindow(ctx, client, d.Timeout(schema.TimeoutCreate), instance.Id, begin.(string),
			d.Get("balancer_active_end").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if replicaSetName, ok := d.GetOk("replica_set_name"); ok && replicaSetName.(string) != "replica" {
		err = updateReplicaSetName(ctx, client, d.Timeout(schema.TimeoutCreate), instance.Id, replicaSetName.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if begin, ok := d.GetOk("maintain_begin"); ok {
		windowOpts := instances.ChangeMaintenanceWindowOpts{
			StartTime: begin.(string),
			EndTime:   d.Get("maintain_end").(string),
		}
		err = instances.UpdateMaintenanceWindow(client, instance.Id, windowOpts)
		if err != nil {
			return diag.Errorf("error setting maintenance window of the DDS instance %s: %s", instance.Id, err)
		}
	}

	return resourceDdsInstanceV3Read(ctx, d, meta)
}

func isEqualPeriod(old, new string) bool {
	if len(old) != len(new) {
		return false
	}
	oldArray := strings.Split(old, ",")
	newArray := strings.Split(new, ",")
	sort.Strings(oldArray)
	sort.Strings(newArray)

	return reflect.DeepEqual(oldArray, newArray)
}

func createBackupStrategy(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	retryFunc := func() (interface{}, bool, error) {
		_, err := instances.CreateBackupPolicy(client, d.Id(), resourceDdsBackupStrategy(d))
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("error creating backup strategy of the DDS instance: %s ", err)
	}

	return nil
}

func updateBalancerStatus(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceId, action string) error {
	retryFunc := func() (interface{}, bool, error) {
		resp, err := instances.UpdateBalancerSwicth(client, instanceId, action)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating balancer switch: %s", err)
	}
	resp := r.(*instances.CommonResp)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      JobStateRefreshFunc(client, resp.JobId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s ", resp.JobId, err)
	}

	return nil
}

func updateBalancerActiveWindow(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceId, begin, end string) error {
	opt := instances.BalancerActiveWindowOpts{
		StartTime: begin,
		StopTime:  end,
	}
	retryFunc := func() (interface{}, bool, error) {
		resp, err := instances.UpdateBalancerActiveWindow(client, instanceId, opt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating balancer active window: %s", err)
	}
	resp := r.(*instances.CommonResp)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      JobStateRefreshFunc(client, resp.JobId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s", resp.JobId, err)
	}

	return nil
}

func updateClientNetworkRanges(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceId string, ranges interface{}) error {
	clientNetworkRanges := utils.ExpandToStringList(ranges.(*schema.Set).List())
	opt := instances.UpdateClientNetworkOpts{
		ClientNetworkRanges: &clientNetworkRanges,
	}
	retryFunc := func() (interface{}, bool, error) {
		err := instances.UpdateClientNetWorkRanges(client, instanceId, opt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating client network ranges: %s", err)
	}

	return nil
}

func UpdateSecondsLevelMonitoring(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceId string, enabled bool) error {
	retryFunc := func() (interface{}, bool, error) {
		_, err := instances.UpdateSecondsLevelMonitoring(client, instanceId, enabled)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating second level monitoring of the DDS instance %s: %s ", instanceId, err)
	}

	return nil
}

func UpdateSlowLogStatus(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceId, slowLogStatus string) error {
	retryFunc := func() (interface{}, bool, error) {
		err := instances.UpdateSlowLogStatus(client, instanceId, slowLogStatus)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating slow log desensitization of the DDS instance %s: %s", instanceId, err)
	}

	return nil
}

func updateReplicaSetName(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceId, replicaSetName string) error {
	opt := instances.ReplicaSetNameOpts{
		Name: replicaSetName,
	}
	retryFunc := func() (interface{}, bool, error) {
		resp, err := instances.UpdateReplicaSetName(client, instanceId, opt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating replica set name: %s", err)
	}
	resp := r.(*instances.CommonResp)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      JobStateRefreshFunc(client, resp.JobId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s", resp.JobId, err)
	}

	return nil
}

func resourceDdsInstanceV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating DDS client: %s", err)
	}

	instanceID := d.Id()
	opts := instances.ListInstanceOpts{
		Id: instanceID,
	}
	allPages, err := instances.List(client, &opts).AllPages()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DdsInstance")
	}
	instanceList, err := instances.ExtractInstances(allPages)
	if err != nil {
		return diag.Errorf("Error extracting DDS instance: %s", err)
	}
	if instanceList.TotalCount == 0 {
		log.Printf("[WARN] DDS instance (%s) was not found", instanceID)
		d.SetId("")
		return nil
	}
	insts := instanceList.Instances
	instanceObj := insts[0]

	log.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instanceObj)

	mErr := multierror.Append(
		d.Set("region", instanceObj.Region),
		d.Set("name", instanceObj.Name),
		d.Set("vpc_id", instanceObj.VpcId),
		d.Set("subnet_id", instanceObj.SubnetId),
		d.Set("security_group_id", instanceObj.SecurityGroupId),
		d.Set("disk_encryption_id", instanceObj.DiskEncryptionId),
		d.Set("mode", instanceObj.Mode),
		d.Set("db_username", instanceObj.DbUserName),
		d.Set("status", instanceObj.Status),
		d.Set("enterprise_project_id", instanceObj.EnterpriseProjectID),
		d.Set("nodes", flattenDdsInstanceV3Nodes(instanceObj)),
		d.Set("groups", flattenDdsInstanceV3Groups(instanceObj)),
		d.Set("description", instanceObj.Remark),
		d.Set("created_at", instanceObj.Created),
		d.Set("updated_at", instanceObj.Updated),
		d.Set("time_zone", instanceObj.TimeZone),
	)

	chargingMode := "postPaid"
	if instanceObj.PayMode == "1" {
		chargingMode = "prePaid"
	}
	mErr = multierror.Append(mErr, d.Set("charging_mode", chargingMode))

	port, err := strconv.Atoi(instanceObj.Port)
	if err != nil {
		log.Printf("[WARNING] Port %s invalid, Type conversion error: %s", instanceObj.Port, err)
	}
	mErr = multierror.Append(mErr, d.Set("port", port))

	sslEnable := true
	if instanceObj.Ssl == 0 {
		sslEnable = false
	}
	mErr = multierror.Append(mErr, d.Set("ssl", sslEnable))

	datastoreList := make([]map[string]interface{}, 0, 1)
	datastore := map[string]interface{}{
		"type":           instanceObj.DataStore.Type,
		"version":        instanceObj.DataStore.Version,
		"storage_engine": instanceObj.Engine,
	}
	datastoreList = append(datastoreList, datastore)
	mErr = multierror.Append(mErr, d.Set("datastore", datastoreList))

	// set backup strategy
	backupStrategyResp, err := instances.GetBackupPolicy(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	backupStrategyList := make([]map[string]interface{}, 0, 1)
	backupStrategy := map[string]interface{}{
		"start_time": backupStrategyResp.BackupPolicy.StartTime,
		"keep_days":  backupStrategyResp.BackupPolicy.KeepDays,
		"period":     backupStrategyResp.BackupPolicy.Period,
	}
	backupStrategyList = append(backupStrategyList, backupStrategy)
	mErr = multierror.Append(mErr, d.Set("backup_strategy", backupStrategyList))

	// set maintenance window
	windows := strings.Split(instanceObj.MaintenanceWindow, "-")
	if len(windows) != 2 {
		return diag.Errorf("invalid format of maintenance window, must be <start_time>-<end_time>")
	}
	mErr = multierror.Append(mErr,
		d.Set("maintain_begin", windows[0]),
		d.Set("maintain_end", windows[1]),
	)

	// save tags
	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagmap))
	} else {
		log.Printf("[WARN] Error fetching tags of DDS instance (%s): %s", d.Id(), err)
	}

	// get second level monitoring
	secondsLevelMonitoring, err := instances.GetSecondsLevelMonitoring(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(mErr, d.Set("second_level_monitoring_enabled", secondsLevelMonitoring.Enabled))

	// save balancer
	if d.Get("mode").(string) == "Sharding" {
		balancer, err := instances.GetBalancer(client, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		balancerStatus := "stop"
		if balancer.IsOpen {
			balancerStatus = "start"
		}
		mErr = multierror.Append(mErr,
			d.Set("balancer_status", balancerStatus),
			d.Set("balancer_active_begin", balancer.ActiveWindow.StartTime),
			d.Set("balancer_active_end", balancer.ActiveWindow.StopTime),
		)
	}

	// set replica set name and client network
	if d.Get("mode").(string) == "ReplicaSet" {
		replicaSetName, err := instances.GetReplicaSetName(client, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		clientNetworkRanges, err := instances.GetClientNetWorkRanges(client, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		mErr = multierror.Append(mErr,
			d.Set("replica_set_name", replicaSetName.Name),
			d.Set("client_network_ranges", clientNetworkRanges.ClientNetworkRanges),
		)
	}

	// set slow log desensitization
	slowLog, err := instances.GetSlowLogStatus(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	mErr = multierror.Append(mErr, d.Set("slow_log_desensitization", slowLog.Status))

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("Error setting dds instance fields: %s", err)
	}

	if ctx.Value(ctxType("configurationIdChanged")) == "true" {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Configuration ID Changed",
				Detail:   "Configuration ID changed, please check whether the instance needs to be restarted.",
			},
		}
	}

	return nil
}

func JobStateRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := jobs.Get(client, jobId)
		if err != nil {
			return nil, "", err
		}

		return resp, resp.Status, nil
	}
}

func waitForInstanceReady(ctx context.Context, client *golangsdk.ServiceClient, instanceId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating"},
		Target:     []string{"normal"},
		Refresh:    ddsInstanceStateRefreshFunc(client, instanceId),
		Timeout:    timeout,
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to become ready: %s ", instanceId, err)
	}

	return nil
}

func resourceDdsInstanceV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	instanceId := d.Id()
	region := cfg.GetRegion(d)
	client, err := cfg.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("Error creating DDS client: %s ", err)
	}

	if d.HasChange("description") {
		opt := instances.RemarkOpts{
			Remark: d.Get("description").(string),
		}
		err = instances.UpdateRemark(client, instanceId, opt)
		if err != nil {
			return diag.Errorf("error updating description of the DDS instance (%s) : %s ", instanceId, err)
		}
	}

	var opts []instances.UpdateOpt
	if d.HasChange("name") {
		opt := instances.UpdateOpt{
			Param:  "new_instance_name",
			Value:  d.Get("name").(string),
			Action: "modify-name",
			Method: "put",
		}
		opts = append(opts, opt)
	}

	if d.HasChange("password") {
		opt := instances.UpdateOpt{
			Param:  "user_pwd",
			Value:  d.Get("password").(string),
			Action: "reset-password",
			Method: "put",
		}
		opts = append(opts, opt)
	}

	if d.HasChange("security_group_id") {
		opt := instances.UpdateOpt{
			Param:  "security_group_id",
			Value:  d.Get("security_group_id").(string),
			Action: "modify-security-group",
			Method: "post",
		}
		opts = append(opts, opt)
	}

	if d.HasChange("backup_strategy") {
		backupStrategy := resourceDdsBackupStrategy(d)
		opt := instances.UpdateOpt{
			Param:  "backup_policy",
			Value:  backupStrategy,
			Action: "backups/policy",
			Method: "put",
		}
		opts = append(opts, opt)
	}

	if d.HasChange("ssl") {
		opt := instances.UpdateOpt{
			Param:  "ssl_option",
			Action: "switch-ssl",
			Method: "post",
		}
		if d.Get("ssl").(bool) {
			opt.Value = "1"
		} else {
			opt.Value = "0"
		}
		opts = append(opts, opt)
	}

	if len(opts) > 0 {
		retryFunc := func() (interface{}, bool, error) {
			resp, err := instances.Update(client, instanceId, opts).Extract()
			retry, err := handleMultiOperationsError(err)
			return resp, retry, err
		}
		r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})

		if err != nil {
			return diag.Errorf("Error updating instance from result: %s ", err)
		}
		resp := r.(*instances.UpdateResp)
		if resp.OrderId != "" {
			bssClient, err := cfg.BssV2Client(region)
			if err != nil {
				return diag.Errorf("error creating BSS v2 client: %s", err)
			}
			err = common.WaitOrderComplete(ctx, bssClient, resp.OrderId, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("balancer_status") {
		err = updateBalancerStatus(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceId, d.Get("balancer_status").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("balancer_active_begin") {
		err = updateBalancerActiveWindow(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceId,
			d.Get("balancer_active_begin").(string), d.Get("balancer_active_end").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("replica_set_name") {
		err = updateReplicaSetName(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceId, d.Get("replica_set_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("port") {
		retryFunc := func() (interface{}, bool, error) {
			resp, err := instances.UpdatePort(client, instanceId, d.Get("port").(int))
			retry, err := handleMultiOperationsError(err)
			return resp, retry, err
		}
		r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error updating database access port: %s", err)
		}
		resp := r.(*instances.PortUpdateResp)
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"Running"},
			Target:       []string{"Completed"},
			Refresh:      JobStateRefreshFunc(client, resp.JobId),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			PollInterval: 10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for the job (%s) completed: %s ", resp.JobId, err)
		}
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "instances", instanceId)
		if tagErr != nil {
			return diag.Errorf("Error updating tags of DDS instance:%s, err:%s", instanceId, tagErr)
		}
	}

	if d.HasChange("maintain_begin") {
		windowOpts := instances.ChangeMaintenanceWindowOpts{
			StartTime: d.Get("maintain_begin").(string),
			EndTime:   d.Get("maintain_end").(string),
		}
		retryFunc := func() (interface{}, bool, error) {
			err = instances.UpdateMaintenanceWindow(client, instanceId, windowOpts)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error setting maintenance window of the DDS instance %s: %s", instanceId, err)
		}
	}

	if d.HasChange("client_network_ranges") {
		err = updateClientNetworkRanges(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceId, d.Get("client_network_ranges"))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("second_level_monitoring_enabled") {
		err = UpdateSecondsLevelMonitoring(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceId,
			d.Get("second_level_monitoring_enabled").(bool))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("slow_log_desensitization") {
		err = UpdateSlowLogStatus(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceId,
			d.Get("slow_log_desensitization").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update flavor
	if d.HasChange("flavor") {
		for i := range d.Get("flavor").([]interface{}) {
			numIndex := fmt.Sprintf("flavor.%d.num", i)
			volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
			specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)

			// The update operation of the number must at the last, lest the new node already has new size or spec-code.
			if d.HasChange(volumeSizeIndex) {
				err := flavorSizeUpdate(ctx, cfg, client, d, i)
				if err != nil {
					return diag.FromErr(err)
				}
			}
			if d.HasChange(specCodeIndex) {
				err := flavorSpecCodeUpdate(ctx, cfg, client, d, i)
				if err != nil {
					return diag.FromErr(err)
				}
			}
			if d.HasChange(numIndex) {
				err := flavorNumUpdate(ctx, cfg, client, d, i)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	// update configuration
	if d.HasChange("configuration") {
		for i := range d.Get("configuration").([]interface{}) {
			configTypePath := fmt.Sprintf("configuration.%d.type", i)
			configIdPath := fmt.Sprintf("configuration.%d.id", i)
			if d.HasChange(configIdPath) {
				// If the DB instance type is cluster and the shard or config parameter template is to be changed, the
				// param is the group ID. If the parameter template of the mongos node is changed, the param is the
				// node ID. If the DB instance to be changed is a replica set instance, the param is the instance ID.
				var ids []string
				configType := d.Get(configTypePath).(string)
				if configType == "replica" {
					ids = []string{instanceId}
				} else {
					ids, err = getDdsInstanceV3GroupIDOrNodeID(client, instanceId, configType)
					if err != nil {
						return diag.FromErr(err)
					}
				}

				// update config ID
				if ctx, err = applyConfigurationToEntity(ctx, client, d.Timeout(schema.TimeoutUpdate), d.Id(),
					d.Get(configIdPath).(string), ids); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "dds",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("availability_zone") {
		azOpt := instances.AvailabilityZoneOpts{
			TargetAzs: d.Get("availability_zone").(string),
		}
		retryFunc := func() (interface{}, bool, error) {
			resp, err := instances.UpdateAvailabilityZone(client, instanceId, azOpt)
			retry, err := handleMultiOperationsError(err)
			return resp, retry, err
		}
		r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error updating availability zone: %s", err)
		}
		resp := r.(*instances.AvailabilityZoneResp)
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"Running"},
			Target:       []string{"Completed"},
			Refresh:      JobStateRefreshFunc(client, resp.JobId),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        10 * time.Second,
			PollInterval: 10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for the job (%s) completed: %s ", resp.JobId, err)
		}
	}

	return resourceDdsInstanceV3Read(ctx, d, meta)
}

func resourceDdsInstanceV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating DDS client: %s ", err)
	}

	instanceId := d.Id()
	// for prePaid mode, we should unsubscribe the resource
	if d.Get("charging_mode").(string) == "prePaid" {
		retryFunc := func() (interface{}, bool, error) {
			err = common.UnsubscribePrePaidResource(d, conf, []string{instanceId})
			retry, err := handleDeletionError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     ddsInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error unsubscribing DDS instance : %s", err)
		}
	} else {
		retryFunc := func() (interface{}, bool, error) {
			result := instances.Delete(client, instanceId)
			retry, err := handleDeletionError(result.Err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     ddsInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"normal", "abnormal", "frozen", "createfail", "enlargefail", "data_disk_full"},
		Target:     []string{"deleted"},
		Refresh:    ddsInstanceStateRefreshFunc(client, instanceId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for instance (%s) to be deleted: %s ",
			instanceId, err)
	}
	log.Printf("[DEBUG] Successfully deleted instance %s", instanceId)
	return nil
}

func flattenDdsInstanceV3Groups(dds instances.InstanceResponse) interface{} {
	nodesList := make([]map[string]interface{}, len(dds.Groups))
	for i, group := range dds.Groups {
		node := map[string]interface{}{
			"id":     group.Id,
			"name":   group.Name,
			"type":   group.Type,
			"status": group.Status,
			"size":   group.Volume.Size,
			"used":   group.Volume.Used,
			"nodes":  flattenDdsInstanceGroupNodes(group.Nodes),
		}
		nodesList[i] = node
	}
	return nodesList
}

func flattenDdsInstanceGroupNodes(nodes []instances.Nodes) interface{} {
	nodesList := make([]map[string]interface{}, len(nodes))
	for i, node := range nodes {
		node := map[string]interface{}{
			"id":         node.Id,
			"name":       node.Name,
			"role":       node.Role,
			"status":     node.Status,
			"private_ip": node.PrivateIP,
			"public_ip":  node.PublicIP,
		}
		nodesList[i] = node
	}
	return nodesList
}

func flattenDdsInstanceV3Nodes(dds instances.InstanceResponse) interface{} {
	nodesList := make([]map[string]interface{}, 0)
	for _, group := range dds.Groups {
		groupType := group.Type
		for _, Node := range group.Nodes {
			node := map[string]interface{}{
				"type":       groupType,
				"id":         Node.Id,
				"name":       Node.Name,
				"role":       Node.Role,
				"status":     Node.Status,
				"private_ip": Node.PrivateIP,
				"public_ip":  Node.PublicIP,
			}
			nodesList = append(nodesList, node)
		}
	}
	return nodesList
}

func getDdsInstanceV3GroupIDOrNodeID(client *golangsdk.ServiceClient, instanceID, getTpye string) ([]string, error) {
	ids := make([]string, 0)

	opts := instances.ListInstanceOpts{
		Id: instanceID,
	}
	allPages, err := instances.List(client, &opts).AllPages()
	if err != nil {
		return ids, fmt.Errorf("error fetching DDS instance: %s", err)
	}
	instanceList, err := instances.ExtractInstances(allPages)
	if err != nil {
		return ids, fmt.Errorf("error extracting DDS instance: %s", err)
	}
	if instanceList.TotalCount == 0 {
		log.Printf("[WARN] DDS instance (%s) was not found", instanceID)
		return ids, nil
	}
	insts := instanceList.Instances
	instanceObj := insts[0]

	log.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instanceObj)

	switch getTpye {
	case "shard", "config":
		for _, group := range instanceObj.Groups {
			if group.Type == getTpye {
				ids = append(ids, group.Id)
			}
		}
	case "mongos":
		for _, group := range instanceObj.Groups {
			if group.Type == "mongos" {
				for _, node := range group.Nodes {
					ids = append(ids, node.Id)
				}
			}
		}
	}

	return ids, nil
}

func flavorUpdate(ctx context.Context, conf *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData,
	opts []instances.UpdateOpt) error {
	retryFunc := func() (interface{}, bool, error) {
		resp, err := instances.Update(client, d.Id(), opts).Extract()
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating instance from result: %s ", err)
	}
	resp := r.(*instances.UpdateResp)
	if resp.OrderId != "" {
		bssClient, err := conf.BssV2Client(conf.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, resp.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}

	err = waitForInstanceReady(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func flavorNumUpdate(ctx context.Context, conf *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)
	if groupType != "mongos" && groupType != "shard" && groupType != "replica" {
		return fmt.Errorf("error updating instance: %s does not support adding nodes", groupType)
	}
	specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)
	volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
	volumeSize := d.Get(volumeSizeIndex).(int)
	numIndex := fmt.Sprintf("flavor.%d.num", i)
	oldNumRaw, newNumRaw := d.GetChange(numIndex)
	oldNum := oldNumRaw.(int)
	newNum := newNumRaw.(int)
	if newNum < oldNum {
		return fmt.Errorf("error updating instance: the new num(%d) must be greater than the old num(%d)", newNum, oldNum)
	}

	var numUpdateOpts []instances.UpdateOpt
	if groupType == "replica" {
		updateNodeNumOpts := instances.UpdateReplicaSetNodeNumOpts{
			Num: newNum - oldNum,
		}
		if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
			updateNodeNumOpts.IsAutoPay = true
		}
		opt := instances.UpdateOpt{
			Param:  "",
			Value:  updateNodeNumOpts,
			Action: "replicaset-node",
			Method: "post",
		}
		numUpdateOpts = append(numUpdateOpts, opt)
	} else {
		updateNodeNumOpts := instances.UpdateNodeNumOpts{
			Type:     groupType,
			SpecCode: d.Get(specCodeIndex).(string),
			Num:      newNum - oldNum,
		}
		if groupType == "shard" {
			volume := instances.VolumeOpts{
				Size: &volumeSize,
			}
			updateNodeNumOpts.Volume = &volume
		}
		if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
			updateNodeNumOpts.IsAutoPay = true
		}
		opt := instances.UpdateOpt{
			Param:  "",
			Value:  updateNodeNumOpts,
			Action: "enlarge",
			Method: "post",
		}
		numUpdateOpts = append(numUpdateOpts, opt)
	}
	err := flavorUpdate(ctx, conf, client, d, numUpdateOpts)
	if err != nil {
		return err
	}
	return nil
}

func flavorSizeUpdate(ctx context.Context, conf *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
	oldSizeRaw, newSizeRaw := d.GetChange(volumeSizeIndex)
	oldSize := oldSizeRaw.(int)
	newSize := newSizeRaw.(int)
	if newSize < oldSize {
		return fmt.Errorf("error updating instance: the new size(%d) must be greater than the old size(%d)", newSize, oldSize)
	}
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)
	if groupType != "replica" && groupType != "single" && groupType != "shard" {
		return fmt.Errorf("error updating instance: %s does not support scaling up storage space", groupType)
	}

	if groupType == "shard" {
		groupIDs, err := getDdsInstanceV3GroupIDOrNodeID(client, d.Id(), "shard")
		if err != nil {
			return err
		}

		for _, groupID := range groupIDs {
			var sizeUpdateOpts []instances.UpdateOpt
			updateVolumeOpts := instances.UpdateVolumeOpts{
				Volume: instances.VolumeOpts{
					GroupID: groupID,
					Size:    &newSize,
				},
			}
			if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
				updateVolumeOpts.IsAutoPay = true
			}
			opt := instances.UpdateOpt{
				Param:  "",
				Value:  updateVolumeOpts,
				Action: "enlarge-volume",
				Method: "post",
			}
			sizeUpdateOpts = append(sizeUpdateOpts, opt)
			err := flavorUpdate(ctx, conf, client, d, sizeUpdateOpts)
			if err != nil {
				return err
			}
		}
	} else {
		var sizeUpdateOpts []instances.UpdateOpt
		updateVolumeOpts := instances.UpdateVolumeOpts{
			Volume: instances.VolumeOpts{
				Size: &newSize,
			},
		}
		if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
			updateVolumeOpts.IsAutoPay = true
		}
		opt := instances.UpdateOpt{
			Param:  "",
			Value:  updateVolumeOpts,
			Action: "enlarge-volume",
			Method: "post",
		}
		sizeUpdateOpts = append(sizeUpdateOpts, opt)
		err := flavorUpdate(ctx, conf, client, d, sizeUpdateOpts)
		if err != nil {
			return err
		}
	}
	return nil
}

func flavorSpecCodeUpdate(ctx context.Context, conf *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)

	if utils.StrSliceContains([]string{"mongos", "shard", "config"}, groupType) {
		nodeIDs, err := getDdsInstanceV3GroupIDOrNodeID(client, d.Id(), groupType)
		if err != nil {
			return err
		}
		for _, ID := range nodeIDs {
			var specUpdateOpts []instances.UpdateOpt
			updateSpecOpts := instances.UpdateSpecOpts{
				Resize: instances.SpecOpts{
					TargetType:     groupType,
					TargetID:       ID,
					TargetSpecCode: d.Get(specCodeIndex).(string),
				},
			}
			if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
				updateSpecOpts.IsAutoPay = true
			}
			opt := instances.UpdateOpt{
				Param:  "",
				Value:  updateSpecOpts,
				Action: "resize",
				Method: "post",
			}
			specUpdateOpts = append(specUpdateOpts, opt)
			err := flavorUpdate(ctx, conf, client, d, specUpdateOpts)
			if err != nil {
				return err
			}
		}
	} else {
		var specUpdateOpts []instances.UpdateOpt
		updateSpecOpts := instances.UpdateSpecOpts{
			Resize: instances.SpecOpts{
				TargetID:       d.Id(),
				TargetSpecCode: d.Get(specCodeIndex).(string),
			},
		}
		if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
			updateSpecOpts.IsAutoPay = true
		}
		opt := instances.UpdateOpt{
			Param:  "",
			Value:  updateSpecOpts,
			Action: "resize",
			Method: "post",
		}
		specUpdateOpts = append(specUpdateOpts, opt)
		err := flavorUpdate(ctx, conf, client, d, specUpdateOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyConfigurationToEntity(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instID, configID string, entityIDs []string) (context.Context, error) {
	applyConfigurationToEntityHttpUrl := "v3/{project_id}/configurations/{config_id}/apply"
	applyConfigurationToEntityPath := client.Endpoint + applyConfigurationToEntityHttpUrl
	applyConfigurationToEntityPath = strings.ReplaceAll(applyConfigurationToEntityPath, "{project_id}", client.ProjectID)
	applyConfigurationToEntityPath = strings.ReplaceAll(applyConfigurationToEntityPath, "{config_id}", configID)

	applyConfigurationToEntityOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"entity_ids": entityIDs,
		},
	}

	// retry, the job_id in return is useless
	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("PUT", applyConfigurationToEntityPath, &applyConfigurationToEntityOpt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}

	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instID),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return ctx, fmt.Errorf("error apply configuration(%s): %s", configID, err)
	}

	// wait for job complete
	err = waitForInstanceReady(ctx, client, instID, timeout)
	if err != nil {
		return ctx, err
	}

	// Sending configurationIdChanged to Read to warn users the instance needs a reboot.
	ctx = context.WithValue(ctx, ctxType("configurationIdChanged"), "true")

	return ctx, nil
}
