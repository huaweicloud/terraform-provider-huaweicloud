package cbh

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBH POST /v2/{project_id}/cbs/instance/{server_id}/eip/bind
// @API CBH POST /v2/{project_id}/cbs/instance/{server_id}/eip/unbind
// @API CBH POST /v2/{project_id}/cbs/instance
// @API CBH GET /v2/{project_id}/cbs/instance/list
// @API CBH PUT /v2/{project_id}/cbs/instance/password
// @API CBH PUT /v2/{project_id}/cbs/instance/{server_id}/security-groups
// @API CBH PUT /v2/{project_id}/cbs/instance
// @API CBH POST /v2/{project_id}/cbs/instance/{resource_id}/tags/action
// @API CBH GET /v2/{project_id}/cbs/instance/{resource_id}/tags
// @API CBH POST /v2/{project_id}/cbs/instance/start
// @API CBH POST /v2/{project_id}/cbs/instance/stop
// @API CBH POST /v2/{project_id}/cbs/instance/reboot
// @API CBH PUT /v2/{project_id}/cbs/instance/vpc
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
func ResourceCBHHAInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHAInstanceCreate,
		UpdateContext: resourceHAInstanceUpdate,
		ReadContext:   resourceHAInstanceRead,
		DeleteContext: resourceHAInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceHAInstanceImportState,
		},

		CustomizeDiff: customdiff.All(
			func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
				isNewMasterIPEqualOldSlaveIP, isNewSlaveIPEqualOldMasterIP := config.CheckValueInterchange(d,
					"master_private_ip", "slave_private_ip")
				isNewMasterAZEqualOldSlaveAZ, isNewSlaveAZEqualOldMasterAZ := config.CheckValueInterchange(d,
					"master_availability_zone", "slave_availability_zone")
				if isNewMasterIPEqualOldSlaveIP && isNewSlaveIPEqualOldMasterIP &&
					isNewMasterAZEqualOldSlaveAZ && isNewSlaveAZEqualOldMasterAZ {
					return fmt.Errorf("the CBH HA instance has undergone a master-slave switch, please modify the" +
						" master_private_ip, slave_private_ip, master_availability_zone, and slave_availability_zone" +
						" parameters in the script")
				}

				return nil
			},
			config.MergeDefaultTags(),
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the CBH HA instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the CBH HA instance.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the product ID of the CBH server.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a VPC.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a subnet.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the IDs of the security group.`,
			},
			"master_availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the availability zone name of the master instance.`,
			},
			"slave_availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the availability zone name of the slave instance.`,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Specifies the password for logging in to the management console.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the charging mode of the CBH HA instance.`,
			},
			"period_unit": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the charging period unit of the CBH HA instance.`,
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the charging period of the CBH HA instance.`,
			},
			"auto_renew": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether auto renew is enabled.`,
			},
			"public_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the elastic IP.`,
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether the IPv6 network is enabled.`,
			},
			"attach_disk_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the size of the additional data disk for the CBH HA instance.`,
			},
			"master_private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the private IP address of the master instance.`,
			},
			"slave_private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the private IP address of the slave instance.`,
			},
			"floating_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the floating IP address of the CBH HA instance.`,
			},
			"power_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the power action after the CBH HA instance is created.`,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the enterprise project ID to which the CBH HA instance belongs.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The elastic IP address.`,
			},
			"master_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the master instance.`,
			},
			"slave_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the slave instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the CBH HA instance.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current version of the CBH HA instance image.`,
			},
			"data_disk_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The data disk size of the CBH HA instance.`,
			},
		},
	}
}

func resourceHAInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cbh"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	orderId, err := createHAInstance(client, d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}

	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for CBH HA instance order (%s) complete: %s", orderId, err)
	}

	resourceIDs, err := common.GetResourceIDsByOrder(bssClient, orderId, 1)
	if err != nil {
		return diag.Errorf("error retrieving resource IDs of CBH HA instance order (%s): %s", orderId, err)
	}

	instances, err := getCBHInstanceList(client)
	if err != nil {
		return diag.FromErr(err)
	}

	var masterId string
	var slaveId string
	for _, resourceId := range resourceIDs {
		serverIDExpression := fmt.Sprintf("[?resource_info.resource_id == '%s']|[0].server_id", resourceId)
		serverId := utils.PathSearch(serverIDExpression, instances, "").(string)
		if serverId == "" {
			return diag.Errorf("unable to find the CBH HA instance ID from the API response")
		}

		instanceTypeExpression := fmt.Sprintf("[?server_id == '%s']|[0].ha_info.instance_type", serverId)
		instanceType := utils.PathSearch(instanceTypeExpression, instances, "").(string)
		if instanceType == "" {
			return diag.Errorf("unable to find the CBH HA instance type from the API response")
		}

		switch instanceType {
		case "master":
			masterId = serverId
		case "slave":
			slaveId = serverId
		}
	}

	if masterId == "" || slaveId == "" {
		return diag.Errorf("error creating CBH HA instance: master instance ID or slave instance ID is not found in API response")
	}

	id := generateCompositeId(masterId, slaveId)
	d.SetId(id)
	if err := waitingForHAInstanceActive(ctx, client, []string{masterId, slaveId}, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for CBH HA instance (%s) creation to active: %s", id, err)
	}

	// After successfully creating an instance using the first security group ID, check if it is necessary to update the
	// security group to the target value.
	securityGroupIDs := d.Get("security_group_id").(string)
	sgIDs := strings.Split(securityGroupIDs, ",")
	if len(sgIDs) > 1 {
		if err := updateSecurityGroup(client, masterId, sgIDs); err != nil {
			return diag.Errorf("error updating the security group after successful creation of CBH HA instance (%s): %s", id, err)
		}
	}

	// Create an instance in the shutdown state.
	if action, ok := d.GetOk("power_action"); ok {
		if action.(string) == Stop {
			if err = doHAInstancePowerAction(ctx, client, d, d.Timeout(schema.TimeoutCreate), masterId); err != nil {
				return diag.FromErr(err)
			}
		} else {
			log.Printf("[WARN] the power action (%s) is invalid after CBH HA instance created", action)
		}
	}

	return resourceHAInstanceRead(ctx, d, meta)
}

func generateCompositeId(masterId, slaveId string) string {
	ids := []string{masterId, slaveId}
	sort.Strings(ids)

	return strings.Join(ids, "/")
}

func createHAInstance(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) (string, error) {
	publicIp, err := buildCreateNetworkPublicIpBodyParam(d, cfg)
	if err != nil {
		return "", fmt.Errorf("error creating CBH HA instance: error building network public IP body: %s", err)
	}

	createInstancePath := client.Endpoint + "v2/{project_id}/cbs/instance"
	createInstancePath = strings.ReplaceAll(createInstancePath, "{project_id}", client.ProjectID)
	createInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateHAInstanceBodyParam(d, cfg.GetRegion(d), cfg.GetEnterpriseProjectID(d), publicIp)),
	}

	createInstanceResp, err := client.Request("POST", createInstancePath, &createInstanceOpt)
	if err != nil {
		return "", fmt.Errorf("error creating CBH HA instance: %s", err)
	}
	createInstanceRespBody, err := utils.FlattenResponse(createInstanceResp)
	if err != nil {
		return "", err
	}

	orderId := utils.PathSearch("order_id", createInstanceRespBody, "").(string)
	if orderId == "" {
		return "", fmt.Errorf("unable to find the order ID of the CBH HA instance from the API response")
	}

	return orderId, nil
}

func buildCreateHAInstanceBodyParam(d *schema.ResourceData, region string, epsId string, publicIp interface{}) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"specification":           d.Get("flavor_id"),
		"instance_name":           d.Get("name"),
		"password":                d.Get("password"),
		"region":                  region,
		"availability_zone":       d.Get("master_availability_zone"),
		"slave_availability_zone": d.Get("slave_availability_zone"),
		"charging_mode":           buildCreateChargingModeParam(d),
		"period_type":             buildCreatePeriodTypeParam(d),
		"period_num":              utils.ValueIgnoreEmpty(d.Get("period")),
		"is_auto_renew":           buildCreateIsAutoRenewParam(d),
		"is_auto_pay":             1,
		"network":                 buildCreateHAInstanceNetworkBodyParam(d, publicIp),
		"attach_disk_size":        utils.ValueIgnoreEmpty(d.Get("attach_disk_size")),
		"enterprise_project_id":   utils.ValueIgnoreEmpty(epsId),
	}
	// The default value of the field `ipv6_enable` is false
	if d.Get("ipv6_enable").(bool) {
		bodyParam["ipv6_enable"] = true
	}

	if tagsRaw, ok := d.GetOk("tags"); ok {
		bodyParam["tags"] = utils.ExpandResourceTags(tagsRaw.(map[string]interface{}))
	}

	return bodyParam
}

func buildCreateHAInstanceNetworkBodyParam(d *schema.ResourceData, publicIp interface{}) interface{} {
	networkBodyParam := map[string]interface{}{
		"vpc_id":          d.Get("vpc_id"),
		"subnet_id":       d.Get("subnet_id"),
		"public_ip":       publicIp,
		"security_groups": buildCreateNetworkSecurityGroupsBodyParam(d),
		"private_ip":      buildCreateHAInstanceNetworkPrivateIpBodyParam(d),
	}

	return networkBodyParam
}

func buildCreateHAInstanceNetworkPrivateIpBodyParam(d *schema.ResourceData) interface{} {
	privateIPBodyParam := make(map[string]interface{})
	if v, ok := d.GetOk("master_private_ip"); ok {
		privateIPBodyParam["ip"] = v.(string)
	}
	if v, ok := d.GetOk("slave_private_ip"); ok {
		privateIPBodyParam["slave_ip"] = v.(string)
	}
	if v, ok := d.GetOk("floating_ip"); ok {
		privateIPBodyParam["floating_ip"] = v.(string)
	}

	return privateIPBodyParam
}

// doPowerAction is a method for CBH HA instance power doing startup, shutdown and reboot actions.
func doHAInstancePowerAction(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, masterId string) error {
	doActionInstancePath := client.Endpoint + "v2/{project_id}/cbs/instance/{action}"
	doActionInstancePath = strings.ReplaceAll(doActionInstancePath, "{project_id}", client.ProjectID)
	var (
		id          = d.Id()
		ids         = strings.Split(id, "/")
		action      = d.Get("power_action").(string)
		jsonBodyMap = map[string]interface{}{
			"server_id": masterId,
		}
	)

	if action == SoftReboot || action == HardReboot {
		doActionInstancePath = strings.ReplaceAll(doActionInstancePath, "{action}", "reboot")
		jsonBodyMap["reboot_type"] = convertPowerActionValue(action)
	} else {
		doActionInstancePath = strings.ReplaceAll(doActionInstancePath, "{action}", action)
	}

	doActionInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         jsonBodyMap,
	}
	_, err := client.Request("POST", doActionInstancePath, &doActionInstanceOpt)
	if err != nil {
		return fmt.Errorf("doing power action (%s) for CBH HA instance (%s) failed: %s", action, id, err)
	}

	if err := waitingForHAInstanceTaskCompleted(ctx, client, ids, timeout); err != nil {
		return fmt.Errorf("error waiting for CBH HA instance (%s) doing action (%s) task to complete: %s", id, action, err)
	}

	if action == Start || action == SoftReboot || action == HardReboot {
		if err := waitingForHAInstanceActive(ctx, client, ids, timeout); err != nil {
			return fmt.Errorf("error waiting for CBH HA instance (%s) doing action (%s) to active: %s", id, action, err)
		}
	}

	if action == Stop {
		if err := waitingForHAInstanceShutoff(ctx, client, ids, timeout); err != nil {
			return fmt.Errorf("error waiting for CBH HA instance (%s) doing action (%s) to shutoff: %s", id, action, err)
		}
	}

	return nil
}

func waitingForHAInstanceUpdateVpcCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, masterId string) error {
	masterExpression := fmt.Sprintf("[?server_id == '%s']|[0]", masterId)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			instances, err := getCBHInstanceList(client)
			if err != nil {
				return nil, "ERROR", err
			}

			masterInstance := utils.PathSearch(masterExpression, instances, nil)
			if masterInstance == nil {
				return nil, "ERROR", golangsdk.ErrDefault404{}
			}

			vpcId := utils.PathSearch("network.vpc_id", masterInstance, "").(string)
			subnetId := utils.PathSearch("network.subnet_id", masterInstance, "").(string)
			if vpcId == d.Get("vpc_id").(string) && subnetId == d.Get("subnet_id").(string) {
				return masterInstance, "COMPLETED", nil
			}

			return masterInstance, "PENDING", nil
		},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func waitingForHAInstanceShutoff(ctx context.Context, client *golangsdk.ServiceClient, ids []string,
	timeout time.Duration) error {
	for _, serverId := range ids {
		expression := fmt.Sprintf("[?server_id == '%s']|[0]", serverId)
		unexpectedStatus := []string{"DELETING", "DELETED", "ERROR", "FROZEN"}
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING"},
			Target:  []string{"COMPLETED"},
			Refresh: func() (interface{}, string, error) {
				instances, err := getCBHInstanceList(client)
				if err != nil {
					return nil, "ERROR", err
				}

				instance := utils.PathSearch(expression, instances, nil)
				if instance == nil {
					return nil, "ERROR", golangsdk.ErrDefault404{}
				}

				status := utils.PathSearch("status_info.status", instance, "").(string)
				if status == "SHUTOFF" {
					return instance, "COMPLETED", nil
				}

				if status == "" {
					return instance, "ERROR", fmt.Errorf("status is not found in list API response")
				}

				if utils.StrSliceContains(unexpectedStatus, status) {
					return instance, status, nil
				}

				return instance, "PENDING", nil
			},
			Timeout:      timeout,
			Delay:        10 * time.Second,
			PollInterval: 10 * time.Second,
		}

		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func waitingForHAInstanceActive(ctx context.Context, client *golangsdk.ServiceClient, ids []string,
	timeout time.Duration) error {
	for _, serverId := range ids {
		expression := fmt.Sprintf("[?server_id == '%s']|[0]", serverId)
		unexpectedStatus := []string{"SHUTOFF", "DELETING", "DELETED", "ERROR", "FROZEN"}
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING"},
			Target:  []string{"COMPLETED"},
			Refresh: func() (interface{}, string, error) {
				instances, err := getCBHInstanceList(client)
				if err != nil {
					return nil, "ERROR", err
				}

				instance := utils.PathSearch(expression, instances, nil)
				if instance == nil {
					return nil, "ERROR", golangsdk.ErrDefault404{}
				}

				status := utils.PathSearch("status_info.status", instance, "").(string)
				if status == "ACTIVE" {
					return instance, "COMPLETED", nil
				}

				if status == "" {
					return instance, "ERROR", fmt.Errorf("status is not found in list API response")
				}

				if utils.StrSliceContains(unexpectedStatus, status) {
					return instance, status, nil
				}

				return instance, "PENDING", nil
			},
			Timeout:      timeout,
			Delay:        10 * time.Second,
			PollInterval: 10 * time.Second,
		}

		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func waitingForHAInstanceTaskCompleted(ctx context.Context, client *golangsdk.ServiceClient, ids []string,
	timeout time.Duration) error {
	for _, serverId := range ids {
		expression := fmt.Sprintf("[?server_id == '%s']|[0]", serverId)
		unexpectedTaskStatus := []string{"delete_wait", "frozen", "unfrozen", "updating", "configuring-ha",
			"data-migrating", "rollback", "traffic-switchover"}
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING"},
			Target:  []string{"COMPLETED"},
			Refresh: func() (interface{}, string, error) {
				instances, err := getCBHInstanceList(client)
				if err != nil {
					return nil, "ERROR", err
				}

				instance := utils.PathSearch(expression, instances, nil)
				if instance == nil {
					return nil, "ERROR", golangsdk.ErrDefault404{}
				}

				taskStatus := utils.PathSearch("status_info.task_status", instance, "").(string)
				if taskStatus == "NO_TASK" {
					return instance, "COMPLETED", nil
				}

				if taskStatus == "" {
					return instance, "ERROR", fmt.Errorf("the instnace task_status is not found in list API response")
				}

				if utils.StrSliceContains(unexpectedTaskStatus, taskStatus) {
					return instance, taskStatus, nil
				}

				return instance, "PENDING", nil
			},
			Timeout:      timeout,
			Delay:        10 * time.Second,
			PollInterval: 10 * time.Second,
		}

		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// nolint: gocyclo
func resourceHAInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cbh"
		id      = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	instances, err := getCBHInstanceList(client)
	if err != nil {
		return diag.FromErr(err)
	}

	ids := strings.Split(id, "/")
	// For master-slave instance, need to extract the master instance server ID from the ID for updating.
	var masterId string
	for _, serverId := range ids {
		instanceTypeExpression := fmt.Sprintf("[?server_id == '%s']|[0].ha_info.instance_type", serverId)
		instanceType := utils.PathSearch(instanceTypeExpression, instances, "").(string)
		if instanceType == "master" {
			masterId = serverId
		}
	}

	if masterId == "" {
		return diag.Errorf("error updating CBH HA instance: the master instance is not found in the list API response")
	}

	// Update instance can only be performed when the instance status is active.
	// Therefore, if `power_action` is start, it must be done in the first step of the update function.
	action := d.Get("power_action").(string)
	if d.HasChanges("power_action") && action == Start {
		if err = doHAInstancePowerAction(ctx, client, d, d.Timeout(schema.TimeoutUpdate), masterId); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update `vpc_id`, `subnet_id`, and private IP related parameters using the same API.
	if d.HasChanges("vpc_id", "subnet_id", "master_private_ip", "slave_private_ip", "floating_ip") {
		if err := updateHAInstanceVpc(client, d, masterId); err != nil {
			return diag.Errorf("error updating the vpc of the CBH HA instance (%s): %s", id, err)
		}

		// When updating the CBH HA instance vpc, the slave instance status always is active.
		// So, here we just need to wait for the master instance.
		if err := waitingForHAInstanceActive(ctx, client, []string{masterId}, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for CBH HA instance (%s) update vpc complete: %s", id, err)
		}

		// Due to API issues, within a short period of time after updating the CBH HA instance vpc, even if the instance
		// status is active, the value of the `vpc_id` and `subnet_id` fields is not the target value.
		// So, here we need to wait for them to reach the target value.
		if err := waitingForHAInstanceUpdateVpcCompleted(ctx, client, d, masterId); err != nil {
			return diag.Errorf("error waiting for CBH HA instance (%s) update vpc to target value: %s", id, err)
		}
	}

	if d.HasChanges("security_group_id") {
		securityGroupIDs := d.Get("security_group_id").(string)
		sgIDs := strings.Split(securityGroupIDs, ",")
		if err := updateSecurityGroup(client, masterId, sgIDs); err != nil {
			return diag.Errorf("error updating the security group of the CBH HA instance (%s): %s", id, err)
		}
	}

	if d.HasChanges("public_ip_id") {
		if err := updateHAInstancePublicIpId(client, d, masterId); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("password") {
		if err = updateHAInstancePassword(client, d, masterId); err != nil {
			return diag.FromErr(err)
		}
	}

	// Due to API limitations, update `flavor_id` and `attach_disk_size` must call the API separately.
	if d.HasChanges("flavor_id") {
		orderId, err := updateFlavorId(client, masterId, d.Get("flavor_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}

		if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}

		if err := waitingForHAInstanceActive(ctx, client, ids, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for CBH HA instance (%s) update to active: %s", id, err)
		}

		// After updating the `flavor_id` and `attach_disk_size`, the instance will automatically restart;
		// At this point, the instance is restarting, but the value of the `status` attribute is already **active**,
		// so we need to continue waiting for the instance's `task_status` attribute to become **NO_TASK**
		// before we can proceed.
		if err := waitingForHAInstanceTaskCompleted(ctx, client, ids, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for CBH instance (%s) update flavor to complete: %s", id, err)
		}
	}

	if d.HasChanges("attach_disk_size") {
		orderId, err := updateAttachDiskSize(client, masterId, int32(d.Get("attach_disk_size").(int)))
		if err != nil {
			return diag.FromErr(err)
		}

		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}

		if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}

		if err := waitingForHAInstanceActive(ctx, client, ids, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for CBH HA instance (%s) update to active: %s", id, err)
		}

		if err := waitingForHAInstanceTaskCompleted(ctx, client, ids, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for CBH HA instance (%s) update additional disk to complete: %s", id, err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}

		masterResourceId, err := getInstanceResourceIdById(client, masterId)
		if err != nil {
			return diag.FromErr(err)
		}

		if masterResourceId == "" {
			return diag.Errorf("error updating the auto-renew of the CBH HA instance (%s): "+
				"master instance resource ID is not found in list API response", id)
		}

		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), masterResourceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the CBH HA instance (%s): %s", id, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		masterResourceId, err := getInstanceResourceIdById(client, masterId)
		if err != nil {
			return diag.FromErr(err)
		}

		if masterResourceId == "" {
			return diag.Errorf("error updating the enterprise project ID of the CBH HA instance (%s): "+
				"master instance resource ID is not found in list API response", id)
		}

		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   masterResourceId,
			ResourceType: "cbh",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		masterResourceId, err := getInstanceResourceIdById(client, masterId)
		if err != nil {
			return diag.FromErr(err)
		}

		if masterResourceId == "" {
			return diag.Errorf("error updating tags of the CBH HA instance (%s): "+
				"master instance resource ID is not found in list API response", id)
		}

		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})

		if len(oMap) > 0 {
			if err = doActionInstanceTags(masterResourceId, "delete", client, oMap); err != nil {
				return diag.FromErr(err)
			}
		}

		if len(nMap) > 0 {
			if err := doActionInstanceTags(masterResourceId, "create", client, nMap); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// The stop and reboot operations should be performed after the instance updated other parameters.
	// Therefore, they must be performed at the end of the update function.
	if d.HasChange("power_action") && (action == Stop || action == SoftReboot || action == HardReboot) {
		if err = doHAInstancePowerAction(ctx, client, d, d.Timeout(schema.TimeoutUpdate), masterId); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceHAInstanceRead(ctx, d, meta)
}

func buildUpdateHAInstanceVpcNetWorkBodyParam(d *schema.ResourceData) interface{} {
	return map[string]interface{}{
		"vpc_id":    d.Get("vpc_id"),
		"subnet_id": d.Get("subnet_id"),
		// When updating VPC, security group ID is a required parameter, but it will not take effect in reality.
		"security_groups": []map[string]interface{}{
			{"id": ""},
		},
		"private_ip": buildCreateHAInstanceNetworkPrivateIpBodyParam(d),
	}
}

func updateHAInstanceVpc(client *golangsdk.ServiceClient, d *schema.ResourceData, masterId string) error {
	updateVpcPath := client.Endpoint + "v2/{project_id}/cbs/instance/vpc"
	updateVpcPath = strings.ReplaceAll(updateVpcPath, "{project_id}", client.ProjectID)
	updateVpcOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"server_id": masterId,
			"network":   buildUpdateHAInstanceVpcNetWorkBodyParam(d),
		},
	}

	_, err := client.Request("PUT", updateVpcPath, &updateVpcOpt)

	return err
}

func updateHAInstancePublicIpId(client *golangsdk.ServiceClient, d *schema.ResourceData, masterId string) error {
	oPublicIpIdRaw, nPublicIpIdRaw := d.GetChange("public_ip_id")
	oPublicIpId := strings.TrimSpace(oPublicIpIdRaw.(string))
	nPublicIpId := strings.TrimSpace(nPublicIpIdRaw.(string))

	if len(oPublicIpId) > 0 {
		if err := unbindHAInstanceEip(client, d, oPublicIpId, masterId); err != nil {
			return err
		}
	}

	if len(nPublicIpId) > 0 {
		if err := bindHAInstanceEip(client, d, nPublicIpId, masterId); err != nil {
			// if bind new eip fail, then bind the old eip to CBH HA instance
			if len(oPublicIpId) > 0 {
				if err := bindHAInstanceEip(client, d, oPublicIpId, masterId); err != nil {
					log.Printf("[WARN] error bind old EIP: %s", err)
				}
			}
			return err
		}
	}

	return nil
}

func unbindHAInstanceEip(client *golangsdk.ServiceClient, d *schema.ResourceData, publicIpId, masterId string) error {
	unbindPath := client.Endpoint + "v2/{project_id}/cbs/instance/{server_id}/eip/unbind"
	unbindPath = strings.ReplaceAll(unbindPath, "{project_id}", client.ProjectID)
	unbindPath = strings.ReplaceAll(unbindPath, "{server_id}", masterId)
	unbindEipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"publicip_id": publicIpId,
		},
	}
	_, err := client.Request("POST", unbindPath, &unbindEipOpt)
	if err != nil {
		return fmt.Errorf("error unbind EIP (%s) from CBH HA instance (%s): %s", publicIpId, d.Id(), err)
	}

	return nil
}

func bindHAInstanceEip(client *golangsdk.ServiceClient, d *schema.ResourceData, publicIpId, masterId string) error {
	bindPath := client.Endpoint + "v2/{project_id}/cbs/instance/{server_id}/eip/bind"
	bindPath = strings.ReplaceAll(bindPath, "{project_id}", client.ProjectID)
	bindPath = strings.ReplaceAll(bindPath, "{server_id}", masterId)
	unbindEipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"publicip_id": publicIpId,
		},
	}
	_, err := client.Request("POST", bindPath, &unbindEipOpt)
	if err != nil {
		return fmt.Errorf("error bind EIP (%s) to CBH HA instance (%s): %s", publicIpId, d.Id(), err)
	}

	return nil
}

func updateHAInstancePassword(client *golangsdk.ServiceClient, d *schema.ResourceData, masterId string) error {
	updatePasswordPath := client.Endpoint + "v2/{project_id}/cbs/instance/password"
	updatePasswordPath = strings.ReplaceAll(updatePasswordPath, "{project_id}", client.ProjectID)
	updatePasswordOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateHAInstanceAdminPasswordParams(d, masterId),
	}
	_, err := client.Request("PUT", updatePasswordPath, &updatePasswordOpt)
	if err != nil {
		return fmt.Errorf("error update CBH HA instance (%s) password: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateHAInstanceAdminPasswordParams(d *schema.ResourceData, masterId string) map[string]interface{} {
	return map[string]interface{}{
		"new_password": d.Get("password").(string),
		"server_id":    masterId,
	}
}

func resourceHAInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cbh"
		id      = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	instances, err := getCBHInstanceList(client)
	if err != nil {
		return diag.FromErr(err)
	}

	var masterInstance interface{}
	var slaveInstance interface{}
	ids := strings.Split(id, "/")
	for _, serverId := range ids {
		expression := fmt.Sprintf("[?server_id == '%s']|[0]", serverId)
		instance := utils.PathSearch(expression, instances, nil)
		instanceType := utils.PathSearch("ha_info.instance_type", instance, "").(string)
		switch instanceType {
		case "master":
			masterInstance = instance
		case "slave":
			slaveInstance = instance
		}
	}

	if masterInstance == nil || slaveInstance == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "CBH HA instance")
	}

	// When EIP is not configured, the query interface field will return a space string.
	publicIpId := strings.TrimSpace(utils.PathSearch("network.public_id", masterInstance, "").(string))

	masterResourceId := utils.PathSearch("resource_info.resource_id", masterInstance, "").(string)
	tags, err := getInstanceTags(masterResourceId, client)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", masterInstance, nil)),
		d.Set("flavor_id", utils.PathSearch("resource_info.specification", masterInstance, nil)),
		d.Set("vpc_id", utils.PathSearch("network.vpc_id", masterInstance, nil)),
		d.Set("subnet_id", utils.PathSearch("network.subnet_id", masterInstance, nil)),
		d.Set("security_group_id", utils.PathSearch("network.security_group_id", masterInstance, nil)),
		d.Set("master_availability_zone", utils.PathSearch("az_info.zone", masterInstance, nil)),
		d.Set("master_private_ip", utils.PathSearch("network.private_ip", masterInstance, nil)),
		d.Set("floating_ip", utils.PathSearch("network.vip", masterInstance, nil)),
		d.Set("public_ip_id", publicIpId),
		d.Set("public_ip", utils.PathSearch("network.public_ip", masterInstance, nil)),
		d.Set("master_id", utils.PathSearch("server_id", masterInstance, nil)),
		d.Set("status", utils.PathSearch("status_info.status", masterInstance, nil)),
		d.Set("data_disk_size", utils.PathSearch("resource_info.data_disk_size", masterInstance, float64(0)).(float64)),
		d.Set("version", utils.PathSearch("bastion_version", masterInstance, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", masterInstance, nil)),
		d.Set("tags", tags),

		d.Set("slave_id", utils.PathSearch("server_id", slaveInstance, nil)),
		d.Set("slave_availability_zone", utils.PathSearch("az_info.zone", slaveInstance, nil)),
		d.Set("slave_private_ip", utils.PathSearch("network.private_ip", slaveInstance, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceHAInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cbh"
		id      = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	instances, err := getCBHInstanceList(client)
	if err != nil {
		return diag.FromErr(err)
	}

	var resourceIDs []string
	var instanceList []interface{}

	ids := strings.Split(id, "/")
	for _, serverId := range ids {
		expression := fmt.Sprintf("[?server_id == '%s']|[0]", serverId)
		instance := utils.PathSearch(expression, instances, nil)
		if instance == nil {
			continue
		}
		instanceList = append(instanceList, instance)
		resourceId := utils.PathSearch("resource_info.resource_id", instance, "").(string)
		if resourceId == "" {
			continue
		}
		resourceIDs = append(resourceIDs, resourceId)
	}

	if len(instanceList) == 0 {
		// Before deleting the CBH HA instance, it is necessary to first call the query API to obtain the resource_id of
		// the instances. If the instances cannot be found, then execute the logic of checkDeleted.
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error deleting the CBH HA instance")
	}

	if len(resourceIDs) == 0 {
		return diag.Errorf("error deleting the CBH HA instance (%s): master and slave instance resource ID is not found in list API response", id)
	}

	if err = common.UnsubscribePrePaidResource(d, cfg, resourceIDs); err != nil {
		return diag.Errorf("error unsubscribe CBH HA instance (%s): %s", id, err)
	}

	if err := waitingForHAInstanceDeleted(ctx, client, ids, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for CBH HA instance (%s) deleted: %s", id, err)
	}

	return nil
}

func waitingForHAInstanceDeleted(ctx context.Context, client *golangsdk.ServiceClient, ids []string,
	timeout time.Duration) error {
	for _, serverId := range ids {
		expression := fmt.Sprintf("[?server_id == '%s']|[0]", serverId)
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING"},
			Target:  []string{"COMPLETED"},
			Refresh: func() (interface{}, string, error) {
				instances, err := getCBHInstanceList(client)
				if err != nil {
					return nil, "ERROR", err
				}

				instance := utils.PathSearch(expression, instances, nil)
				if instance == nil {
					m := map[string]string{"code": "COMPLETED"}
					return m, "COMPLETED", nil
				}

				status := utils.PathSearch("status_info.status", instance, "").(string)
				if status == "DELETED" {
					return instance, "COMPLETED", nil
				}

				return instance, "PENDING", nil
			},
			Timeout:      timeout,
			Delay:        30 * time.Second,
			PollInterval: 10 * time.Second,
		}

		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceHAInstanceImportState(_ context.Context, d *schema.ResourceData, meta interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <master_id>/<slave_id>")
	}

	var (
		masterId = parts[0]
		slaveId  = parts[1]
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "cbh"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating CBH client: %s", err)
	}

	instances, err := getCBHInstanceList(client)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	masterInstanceExpression := fmt.Sprintf("[?server_id == '%s']|[0]", masterId)
	masterInstance := utils.PathSearch(masterInstanceExpression, instances, nil)
	slaveInstanceExpression := fmt.Sprintf("[?server_id == '%s']|[0]", slaveId)
	slaveInstance := utils.PathSearch(slaveInstanceExpression, instances, nil)

	if masterInstance == nil || slaveInstance == nil {
		return []*schema.ResourceData{d}, fmt.Errorf("master instance or slave instance not found in list API response")
	}

	id := generateCompositeId(masterId, slaveId)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
