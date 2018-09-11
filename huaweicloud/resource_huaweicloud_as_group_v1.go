package huaweicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/groups"
	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/instances"
)

func resourceASGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceASGroupCreate,
		Read:   resourceASGroupRead,
		Update: resourceASGroupUpdate,
		Delete: resourceASGroupDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_group_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: resourceASGroupValidateGroupName,
				ForceNew:     false,
			},
			"scaling_configuration_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"desire_instance_number": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"min_instance_number": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: false,
			},
			"max_instance_number": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: false,
			},
			"cool_down_time": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      900,
				ValidateFunc: resourceASGroupValidateCoolDownTime,
				ForceNew:     false,
				Description:  "The cooling duration, in seconds.",
			},
			"lb_listener_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: resourceASGroupValidateListenerId,
				Description:  "The system supports the binding of up to three ELB listeners, the IDs of which are separated using a comma.",
			},
			"available_zones": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: false,
			},
			"networks": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 5,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				ForceNew: false,
			},
			"security_groups": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				ForceNew: false,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"health_periodic_audit_method": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: resourceASGroupValidateHealthAuditMethod,
				ForceNew:     false,
				Default:      "NOVA_AUDIT",
			},
			"health_periodic_audit_time": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: resourceASGroupValidateHealthAuditTime,
				ForceNew:     false,
				Description:  "The health check period for instances, in minutes.",
			},
			"instance_terminate_policy": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "OLD_CONFIG_OLD_INSTANCE",
				ValidateFunc: resourceASGroupValidateTerminatePolicy,
				ForceNew:     false,
			},
			"notifications": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: false,
			},
			"delete_publicip": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: false,
			},
			"delete_instances": &schema.Schema{
				Description: "Whether to delete instances when they are removed from the AS group.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				ForceNew:    false,
			},
			"instances": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				ForceNew:    false,
				Description: "The instances id list in the as group.",
			},
		},
	}
}

type Network struct {
	ID string
}

type Group struct {
	ID string
}

func expandNetworks(Networks []Network) []groups.NetworkOpts {
	var networks []groups.NetworkOpts
	for _, v := range Networks {
		n := groups.NetworkOpts{
			ID: v.ID,
		}
		networks = append(networks, n)
	}

	return networks
}

func expandGroups(Groups []Group) []groups.SecurityGroupOpts {
	var asgroups []groups.SecurityGroupOpts
	for _, v := range Groups {
		n := groups.SecurityGroupOpts{
			ID: v.ID,
		}
		asgroups = append(asgroups, n)
	}

	return asgroups
}

func getAllAvailableZones(d *schema.ResourceData) []string {
	rawZones := d.Get("available_zones").([]interface{})
	zones := make([]string, len(rawZones))
	for i, raw := range rawZones {
		zones[i] = raw.(string)
	}
	log.Printf("[DEBUG] getAvailableZones: %#v", zones)

	return zones
}

func getAllNotifications(d *schema.ResourceData) []string {
	rawNotifications := d.Get("notifications").([]interface{})
	notifications := make([]string, len(rawNotifications))
	for i, raw := range rawNotifications {
		notifications[i] = raw.(string)
	}
	log.Printf("[DEBUG] getNotifications: %#v", notifications)

	return notifications
}

func getAllNetworks(d *schema.ResourceData, meta interface{}) []Network {
	var Networks []Network

	networks := d.Get("networks").([]interface{})
	for _, v := range networks {
		network := v.(map[string]interface{})
		networkID := network["id"].(string)
		v := Network{
			ID: networkID,
		}
		Networks = append(Networks, v)
	}

	log.Printf("[DEBUG] getNetworks: %#v", Networks)
	return Networks
}

func getAllSecurityGroups(d *schema.ResourceData, meta interface{}) []Group {
	var Groups []Group

	asgroups := d.Get("security_groups").([]interface{})
	for _, v := range asgroups {
		group := v.(map[string]interface{})
		groupID := group["id"].(string)
		v := Group{
			ID: groupID,
		}
		Groups = append(Groups, v)
	}

	log.Printf("[DEBUG] getGroups: %#v", Groups)
	return Groups
}

func getInstancesInGroup(asClient *golangsdk.ServiceClient, groupID string, opts instances.ListOptsBuilder) ([]instances.Instance, error) {
	var insList []instances.Instance
	page, err := instances.List(asClient, groupID, opts).AllPages()
	if err != nil {
		return insList, fmt.Errorf("Error getting instances in ASGroup %q: %s", groupID, err)
	}
	insList, err = page.(instances.InstancePage).Extract()
	return insList, err
}
func getInstancesIDs(allIns []instances.Instance) []string {
	var allIDs []string
	for _, ins := range allIns {
		// Maybe the instance is pending, so we can't get the id,
		// so unable to delete the instance this time, maybe next time to execute
		// terraform destroy will works
		if ins.ID != "" {
			allIDs = append(allIDs, ins.ID)
		}
	}
	log.Printf("[DEBUG] Get instances in ASGroups: %#v", allIDs)
	return allIDs
}

func getInstancesLifeStates(allIns []instances.Instance) []string {
	var allLifeStates []string
	for _, ins := range allIns {
		allLifeStates = append(allLifeStates, ins.LifeCycleStatus)
	}
	log.Printf("[DEBUG] Get instances lifecycle status in ASGroups: %#v", allLifeStates)
	return allLifeStates
}

func refreshInstancesLifeStates(asClient *golangsdk.ServiceClient, groupID string, insNum int, checkInService bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var opts instances.ListOptsBuilder
		allIns, err := getInstancesInGroup(asClient, groupID, opts)
		if err != nil {
			return nil, "ERROR", err
		}
		// maybe the instances (or some of the instances) have not put in the asg when creating
		if checkInService && len(allIns) != insNum {
			return allIns, "PENDING", err
		}
		allLifeStatus := getInstancesLifeStates(allIns)
		for _, lifeStatus := range allLifeStatus {
			log.Printf("[DEBUG] Get lifecycle status in group %s: %s", groupID, lifeStatus)
			// check for creation
			if checkInService {
				if lifeStatus == "PENDING" || lifeStatus == "REMOVING" {
					return allIns, lifeStatus, err
				}
			}
			// check for removal
			if !checkInService {
				if lifeStatus == "REMOVING" || lifeStatus != "INSERVICE" {
					return allIns, lifeStatus, err
				}
			}
		}
		if checkInService {
			return allIns, "INSERVICE", err
		}
		log.Printf("[DEBUG] Exit refreshInstancesLifeStates for %q!", groupID)
		return allIns, "", err
	}
}

func checkASGroupInstancesInService(asClient *golangsdk.ServiceClient, groupID string, insNum int, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"INSERVICE"}, //if there is no lifecyclestatus, meaning no instances in asg
		Refresh: refreshInstancesLifeStates(asClient, groupID, insNum, true),
		Timeout: timeout,
		Delay:   10 * time.Second,
	}

	_, err := stateConf.WaitForState()

	return err
}

func checkASGroupInstancesRemoved(asClient *golangsdk.ServiceClient, groupID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"REMOVING"},
		Target:  []string{""}, //if there is no lifecyclestatus, meaning no instances in asg
		Refresh: refreshInstancesLifeStates(asClient, groupID, 0, false),
		Timeout: timeout,
		Delay:   10 * time.Second,
	}

	_, err := stateConf.WaitForState()

	return err
}

func resourceASGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	asClient, err := config.autoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}
	log.Printf("[DEBUG] asClient: %#v", asClient)

	minNum := d.Get("min_instance_number").(int)
	maxNum := d.Get("max_instance_number").(int)
	desireNum := d.Get("desire_instance_number").(int)
	log.Printf("[DEBUG] Min instance number is: %#v", minNum)
	log.Printf("[DEBUG] Max instance number is: %#v", maxNum)
	log.Printf("[DEBUG] Desire instance number is: %#v", desireNum)
	if desireNum < minNum || desireNum > maxNum {
		return fmt.Errorf("Invalid parameters: it should be min_instance_number<=desire_instance_number<=max_instance_number")
	}
	var initNum int
	if desireNum > 0 {
		initNum = desireNum
	} else {
		initNum = minNum
	}
	log.Printf("[DEBUG] Init instance number is: %#v", initNum)
	networks := getAllNetworks(d, meta)
	asgNetworks := expandNetworks(networks)

	secGroups := getAllSecurityGroups(d, meta)
	asgSecGroups := expandGroups(secGroups)

	log.Printf("[DEBUG] available_zones: %#v", d.Get("available_zones"))
	createOpts := groups.CreateOpts{
		Name:                 d.Get("scaling_group_name").(string),
		ConfigurationID:      d.Get("scaling_configuration_id").(string),
		DesireInstanceNumber: desireNum,
		MinInstanceNumber:    minNum,
		MaxInstanceNumber:    maxNum,
		CoolDownTime:         d.Get("cool_down_time").(int),
		LBListenerID:         d.Get("lb_listener_id").(string),
		AvailableZones:       getAllAvailableZones(d),
		Networks:             asgNetworks,
		SecurityGroup:        asgSecGroups,
		VpcID:                d.Get("vpc_id").(string),
		HealthPeriodicAuditMethod: d.Get("health_periodic_audit_method").(string),
		HealthPeriodicAuditTime:   d.Get("health_periodic_audit_time").(int),
		InstanceTerminatePolicy:   d.Get("instance_terminate_policy").(string),
		Notifications:             getAllNotifications(d),
		IsDeletePublicip:          d.Get("delete_publicip").(bool),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	asgId, err := groups.Create(asClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating ASGroup: %s", err)
	}

	d.SetId(asgId)

	//enable asg
	enableResult := groups.Enable(asClient, asgId)
	if enableResult.Err != nil {
		return fmt.Errorf("Error enabling ASGroup %q: %s", asgId, enableResult.Err)
	}
	log.Printf("[DEBUG] Enable ASGroup %q success!", asgId)
	// check all instances are inservice
	if initNum > 0 {
		timeout := d.Timeout(schema.TimeoutCreate)
		err = checkASGroupInstancesInService(asClient, asgId, initNum, timeout)
		if err != nil {
			return fmt.Errorf("Error waiting for instances in the ASGroup %q to become inservice!!: %s", asgId, err)
		}
	}

	return resourceASGroupRead(d, meta)
}

func resourceASGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	asClient, err := config.autoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}

	asg, err := groups.Get(asClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "AS group")
	}
	log.Printf("[DEBUG] Retrieved ASGroup %q: %+v", d.Id(), asg)
	log.Printf("[DEBUG] Retrieved ASGroup %q notifications: %+v", d.Id(), asg.Notifications)
	log.Printf("[DEBUG] Retrieved ASGroup %q availablezones: %+v", d.Id(), asg.AvailableZones)
	log.Printf("[DEBUG] Retrieved ASGroup %q networks: %+v", d.Id(), asg.Networks)
	log.Printf("[DEBUG] Retrieved ASGroup %q secgroups: %+v", d.Id(), asg.SecurityGroups)

	// set properties based on the read info
	d.Set("scaling_group_name", asg.Name)
	d.Set("desire_instance_number", asg.DesireInstanceNumber)
	d.Set("min_instance_number", asg.MinInstanceNumber)
	d.Set("max_instance_number", asg.MaxInstanceNumber)
	d.Set("cool_down_time", asg.CoolDownTime)
	d.Set("lb_listener_id", asg.LBListenerID)
	d.Set("health_periodic_audit_method", asg.HealthPeriodicAuditMethod)
	d.Set("health_periodic_audit_time", asg.HealthPeriodicAuditTime)
	d.Set("instance_terminate_policy", asg.InstanceTerminatePolicy)
	d.Set("scaling_configuration_id", asg.ConfigurationID)
	d.Set("delete_publicip", asg.DeletePublicip)
	if len(asg.Notifications) >= 1 {
		d.Set("notifications", asg.Notifications)
	}

	var opts instances.ListOptsBuilder
	allIns, err := getInstancesInGroup(asClient, d.Id(), opts)
	if err != nil {
		return fmt.Errorf("Can not get the instances in ASGroup %q!!: %s", d.Id(), err)
	}
	allIDs := getInstancesIDs(allIns)
	d.Set("instances", allIDs)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceASGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	asClient, err := config.autoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}
	d.Partial(true)
	if d.HasChange("min_instance_number") || d.HasChange("max_instance_number") || d.HasChange("desire_instance_number") {
		minNum := d.Get("min_instance_number").(int)
		maxNum := d.Get("max_instance_number").(int)
		desireNum := d.Get("desire_instance_number").(int)
		log.Printf("[DEBUG] Min instance number is: %#v", minNum)
		log.Printf("[DEBUG] Max instance number is: %#v", maxNum)
		log.Printf("[DEBUG] Desire instance number is: %#v", desireNum)
		if desireNum < minNum || desireNum > maxNum {
			return fmt.Errorf("Invalid parameters: it should be min_instance_number<=desire_instance_number<=max_instance_number")
		}

	}

	networks := getAllNetworks(d, meta)
	asgNetworks := expandNetworks(networks)

	secGroups := getAllSecurityGroups(d, meta)
	asgSecGroups := expandGroups(secGroups)
	updateOpts := groups.UpdateOpts{
		Name:                      d.Get("scaling_group_name").(string),
		ConfigurationID:           d.Get("scaling_configuration_id").(string),
		DesireInstanceNumber:      d.Get("desire_instance_number").(int),
		MinInstanceNumber:         d.Get("min_instance_number").(int),
		MaxInstanceNumber:         d.Get("max_instance_number").(int),
		CoolDownTime:              d.Get("cool_down_time").(int),
		LBListenerID:              d.Get("lb_listener_id").(string),
		AvailableZones:            getAllAvailableZones(d),
		Networks:                  asgNetworks,
		SecurityGroup:             asgSecGroups,
		HealthPeriodicAuditMethod: d.Get("health_periodic_audit_method").(string),
		HealthPeriodicAuditTime:   d.Get("health_periodic_audit_time").(int),
		InstanceTerminatePolicy:   d.Get("instance_terminate_policy").(string),
		Notifications:             getAllNotifications(d),
		IsDeletePublicip:          d.Get("delete_publicip").(bool),
	}
	asgID, err := groups.Update(asClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating ASGroup %q: %s", asgID, err)
	}
	d.Partial(false)
	return resourceASGroupRead(d, meta)
}

func resourceASGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	asClient, err := config.autoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}

	log.Printf("[DEBUG] Begin to get instances of ASGroup %q", d.Id())
	var listOpts instances.ListOptsBuilder
	allIns, err := getInstancesInGroup(asClient, d.Id(), listOpts)
	if err != nil {
		return fmt.Errorf("Error listing instances of asg: %s", err)
	}
	allLifeStatus := getInstancesLifeStates(allIns)
	for _, lifeCycleState := range allLifeStatus {
		if lifeCycleState != "INSERVICE" {
			return fmt.Errorf("[DEBUG] Can't delete the ASGroup %q: There are some instances not in INSERVICE but in %s, try again latter.", d.Id(), lifeCycleState)
		}
	}
	allIDs := getInstancesIDs(allIns)
	log.Printf("[DEBUG] InstanceIDs in ASGroup %q: %+v", d.Id(), allIDs)
	log.Printf("[DEBUG] There are %d instances in ASGroup %q", len(allIDs), d.Id())
	if len(allLifeStatus) > 0 {
		min_number := d.Get("min_instance_number").(int)
		if min_number > 0 {
			return fmt.Errorf("[DEBUG] Can't delete the ASGroup %q: The instance number after the removal will less than the min number %d, modify the min number to zero first.", d.Id(), min_number)
		}
		delete_ins := d.Get("delete_instances").(string)
		log.Printf("[DEBUG] The flag delete_instances in ASGroup is %s", delete_ins)
		batchResult := instances.BatchDelete(asClient, d.Id(), allIDs, delete_ins)
		if batchResult.Err != nil {
			return fmt.Errorf("Error removing instancess of asg: %s", batchResult.Err)
		}
		log.Printf("[DEBUG] Begin to remove instances of ASGroup %q", d.Id())
		timeout := d.Timeout(schema.TimeoutDelete)
		err = checkASGroupInstancesRemoved(asClient, d.Id(), timeout)
		if err != nil {
			return fmt.Errorf(
				"[DEBUG] Error removing instances from ASGroup %q: %s", d.Id(), err)
		}
	}

	log.Printf("[DEBUG] Begin to delete ASGroup %q", d.Id())
	if delErr := groups.Delete(asClient, d.Id()).ExtractErr(); delErr != nil {
		return fmt.Errorf("Error deleting ASGroup: %s", delErr)
	}

	return nil
}

var TerminatePolices = [4]string{"OLD_CONFIG_OLD_INSTANCE", "OLD_CONFIG_NEW_INSTANCE", "OLD_INSTANCE", "NEW_INSTANCE"}

func resourceASGroupValidateTerminatePolicy(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range TerminatePolices {
		if value == TerminatePolices[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, TerminatePolices))
	return
}

func resourceASGroupValidateCoolDownTime(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if 0 <= value && value <= 86400 {
		return
	}
	errors = append(errors, fmt.Errorf("%q must be [0, 86400]", k))
	return
}

func resourceASGroupValidateListenerId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	split := strings.Split(value, ",")
	if len(split) <= 3 {
		return
	}
	errors = append(errors, fmt.Errorf("%q supports binding up to 3 ELB listeners which are separated by a comma.", k))
	return
}

var HealthAuditMethods = [2]string{"ELB_AUDIT", "NOVA_AUDIT"}

func resourceASGroupValidateHealthAuditMethod(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range HealthAuditMethods {
		if value == HealthAuditMethods[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, HealthAuditMethods))
	return
}

var HealthAuditTime = [4]int{5, 15, 60, 180}

func resourceASGroupValidateHealthAuditTime(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	for i := range HealthAuditTime {
		if value == HealthAuditTime[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, HealthAuditTime))
	return
}

func resourceASGroupValidateGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 || len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q must contain more than 1 and less than 64 characters", k))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z-_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("only alphanumeric characters, hyphens, and underscores allowed in %q", k))
	}
	return
}
