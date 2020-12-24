package huaweicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/groups"
	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/instances"
	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/tags"
)

func ResourceASGroup() *schema.Resource {
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
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: resourceASGroupValidateGroupName,
			},
			"scaling_configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"desire_instance_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_instance_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"max_instance_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"cool_down_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ValidateFunc: resourceASGroupValidateCoolDownTime,
				Description:  "The cooling duration, in seconds.",
			},
			"lb_listener_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: resourceASGroupValidateListenerId,
				Description:  "The system supports the binding of up to six ELB listeners, the IDs of which are separated using a comma.",
				Deprecated:   "use lbaas_listeners instead",
			},
			"lbaas_listeners": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      6,
				ConflictsWith: []string{"lb_listener_id"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pool_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
					},
				},
			},
			"available_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"networks": {
				Type:     schema.TypeList,
				MaxItems: 5,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"security_groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"health_periodic_audit_method": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: resourceASGroupValidateHealthAuditMethod,
				Default:      "NOVA_AUDIT",
			},
			"health_periodic_audit_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: resourceASGroupValidateHealthAuditTime,
				Description:  "The health check period for instances, in minutes.",
			},
			"instance_terminate_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "OLD_CONFIG_OLD_INSTANCE",
				ValidateFunc: resourceASGroupValidateTerminatePolicy,
			},
			"notifications": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"delete_publicip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"delete_instances": {
				Description: "Whether to delete instances when they are removed from the AS group.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The instances id list in the as group.",
			},
			"current_instance_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func expandGroupsTags(tagmap map[string]interface{}) []tags.ResourceTag {
	var taglist []tags.ResourceTag

	for k, v := range tagmap {
		tag := tags.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	return taglist
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

func getAllLBaaSListeners(d *schema.ResourceData, meta interface{}) []groups.LBaaSListenerOpts {
	var aslisteners []groups.LBaaSListenerOpts

	listeners := d.Get("lbaas_listeners").([]interface{})
	for _, v := range listeners {
		listener := v.(map[string]interface{})
		s := groups.LBaaSListenerOpts{
			PoolID:       listener["pool_id"].(string),
			ProtocolPort: listener["protocol_port"].(int),
			Weight:       listener["weight"].(int),
		}
		aslisteners = append(aslisteners, s)
	}

	log.Printf("[DEBUG] getAllLBaaSListeners: %#v", aslisteners)
	return aslisteners
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
	asClient, err := config.AutoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}
	log.Printf("[DEBUG] asClient: %#v", asClient)

	minNum := d.Get("min_instance_number").(int)
	maxNum := d.Get("max_instance_number").(int)
	var desireNum int
	if v, ok := d.GetOk("desire_instance_number"); ok {
		desireNum = v.(int)
	} else {
		desireNum = minNum
	}
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

	asgLBaaSListeners := getAllLBaaSListeners(d, meta)

	log.Printf("[DEBUG] available_zones: %#v", d.Get("available_zones"))
	createOpts := groups.CreateOpts{
		Name:                      d.Get("scaling_group_name").(string),
		ConfigurationID:           d.Get("scaling_configuration_id").(string),
		DesireInstanceNumber:      desireNum,
		MinInstanceNumber:         minNum,
		MaxInstanceNumber:         maxNum,
		CoolDownTime:              d.Get("cool_down_time").(int),
		LBListenerID:              d.Get("lb_listener_id").(string),
		LBaaSListeners:            asgLBaaSListeners,
		AvailableZones:            getAllAvailableZones(d),
		Networks:                  asgNetworks,
		SecurityGroup:             asgSecGroups,
		VpcID:                     d.Get("vpc_id").(string),
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

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := expandGroupsTags(tagRaw)
		if tagErr := tags.Create(asClient, asgId, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of ASGroup %q: %s", asgId, tagErr)
		}
	}

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
	asClient, err := config.AutoscalingV1Client(GetRegion(d, config))
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
	log.Printf("[DEBUG] Retrieved ASGroup %q lbaaslisteners: %+v", d.Id(), asg.LBaaSListeners)

	// set properties based on the read info
	d.Set("scaling_group_name", asg.Name)
	d.Set("status", asg.Status)
	d.Set("current_instance_number", asg.ActualInstanceNumber)
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
	if len(asg.LBaaSListeners) >= 1 {
		listeners := make([]map[string]interface{}, len(asg.LBaaSListeners))
		for i, listener := range asg.LBaaSListeners {
			listeners[i] = make(map[string]interface{})
			listeners[i]["pool_id"] = listener.PoolID
			listeners[i]["protocol_port"] = listener.ProtocolPort
			listeners[i]["weight"] = listener.Weight
		}
		d.Set("lbaas_listeners", listeners)
	}

	var opts instances.ListOptsBuilder
	allIns, err := getInstancesInGroup(asClient, d.Id(), opts)
	if err != nil {
		return fmt.Errorf("Can not get the instances in ASGroup %q!!: %s", d.Id(), err)
	}
	allIDs := getInstancesIDs(allIns)
	d.Set("instances", allIDs)

	d.Set("region", GetRegion(d, config))

	// save group tags
	resourceTags, err := tags.Get(asClient, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching HuaweiCloud ASGroup tags: %s", err)
	}

	tagmap := make(map[string]string)
	for _, val := range resourceTags.Tags {
		tagmap[val.Key] = val.Value
	}
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("Error saving tags for HuaweiCloud ASGroup (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceASGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	asClient, err := config.AutoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}
	var desireNum int
	minNum := d.Get("min_instance_number").(int)
	maxNum := d.Get("max_instance_number").(int)
	if v, ok := d.GetOk("desire_instance_number"); ok {
		desireNum = v.(int)
	} else {
		desireNum = minNum
	}
	if d.HasChanges("min_instance_number", "max_instance_number", "desire_instance_number") {
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

	asgLBaaSListeners := getAllLBaaSListeners(d, meta)
	updateOpts := groups.UpdateOpts{
		Name:                      d.Get("scaling_group_name").(string),
		ConfigurationID:           d.Get("scaling_configuration_id").(string),
		DesireInstanceNumber:      desireNum,
		MinInstanceNumber:         minNum,
		MaxInstanceNumber:         maxNum,
		CoolDownTime:              d.Get("cool_down_time").(int),
		LBListenerID:              d.Get("lb_listener_id").(string),
		LBaaSListeners:            asgLBaaSListeners,
		AvailableZones:            getAllAvailableZones(d),
		Networks:                  asgNetworks,
		SecurityGroup:             asgSecGroups,
		HealthPeriodicAuditMethod: d.Get("health_periodic_audit_method").(string),
		HealthPeriodicAuditTime:   d.Get("health_periodic_audit_time").(int),
		InstanceTerminatePolicy:   d.Get("instance_terminate_policy").(string),
		Notifications:             getAllNotifications(d),
		IsDeletePublicip:          d.Get("delete_publicip").(bool),
	}

	log.Printf("[DEBUG] AS Group update options: %#v", updateOpts)
	asgID, err := groups.Update(asClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating ASGroup %q: %s", asgID, err)
	}

	//update tags
	if d.HasChange("tags") {
		//remove old tags and set new tags
		old, new := d.GetChange("tags")
		oldRaw := old.(map[string]interface{})
		if len(oldRaw) > 0 {
			taglist := expandGroupsTags(oldRaw)
			if tagErr := tags.Delete(asClient, asgID, taglist).ExtractErr(); tagErr != nil {
				return fmt.Errorf("Error deleting tags of ASGroup %q: %s", asgID, tagErr)
			}
		}

		newRaw := new.(map[string]interface{})
		if len(newRaw) > 0 {
			taglist := expandGroupsTags(newRaw)
			if tagErr := tags.Create(asClient, asgID, taglist).ExtractErr(); tagErr != nil {
				return fmt.Errorf("Error setting tags of ASGroup %q: %s", asgID, tagErr)
			}
		}
	}

	return resourceASGroupRead(d, meta)
}

func resourceASGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	asClient, err := config.AutoscalingV1Client(GetRegion(d, config))
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
	if len(split) <= 6 {
		return
	}
	errors = append(errors, fmt.Errorf("%q supports binding up to 6 ELB listeners which are separated by a comma.", k))
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

//lintignore:V001
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
