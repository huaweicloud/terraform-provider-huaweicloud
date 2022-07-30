package as

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/groups"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/instances"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/tags"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

var (
	HealthAuditMethods = []string{"ELB_AUDIT", "NOVA_AUDIT"}
	HealthAuditTime    = []int{0, 1, 5, 15, 60, 180}
	TerminatePolices   = []string{"OLD_CONFIG_OLD_INSTANCE", "OLD_CONFIG_NEW_INSTANCE", "OLD_INSTANCE", "NEW_INSTANCE"}
)

func ResourceASGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASGroupCreate,
		ReadContext:   resourceASGroupRead,
		UpdateContext: resourceASGroupUpdate,
		DeleteContext: resourceASGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa50-9a-zA-Z-_]+$"),
						"only letters, digits, underscores (_), and hyphens (-) are allowed"),
				),
			},
			"scaling_configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Required",
			},
			"desire_instance_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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
				ValidateFunc: validation.IntBetween(0, 86400),
				Description:  "The cooling duration, in seconds.",
			},
			"lbaas_listeners": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
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
				ValidateFunc: validation.StringInSlice(HealthAuditMethods, false),
				Default:      "NOVA_AUDIT",
			},
			"health_periodic_audit_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validation.IntInSlice(HealthAuditTime),
				Description:  "The health check period for instances, in minutes.",
			},
			"instance_terminate_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "OLD_CONFIG_OLD_INSTANCE",
				ValidateFunc: validation.StringInSlice(TerminatePolices, false),
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
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

			// Deprecated
			"lb_listener_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: resourceASGroupValidateListenerId,
				Description:  "The system supports the binding of up to six ELB listeners, the IDs of which are separated using a comma.",
				Deprecated:   "use lbaas_listeners instead",
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

	return zones
}

func getAllNotifications(d *schema.ResourceData) []string {
	rawNotifications := d.Get("notifications").([]interface{})
	notifications := make([]string, len(rawNotifications))
	for i, raw := range rawNotifications {
		notifications[i] = raw.(string)
	}

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

	return aslisteners
}

func getInstancesInGroup(asClient *golangsdk.ServiceClient, groupID string, opts instances.ListOptsBuilder) ([]instances.Instance, error) {
	var insList []instances.Instance
	page, err := instances.List(asClient, groupID, opts).AllPages()
	if err != nil {
		return insList, fmtp.Errorf("error getting instances in AS group %s: %s", groupID, err)
	}
	insList, err = page.(instances.InstancePage).Extract()
	return insList, err
}

func getInstancesIDs(allIns []instances.Instance) []string {
	var allIDs = make([]string, 0, len(allIns))
	for _, ins := range allIns {
		// Maybe the instance is pending, so we can't get the id,
		// so unable to delete the instance this time, maybe next time to execute
		// terraform destroy will works
		if ins.ID != "" {
			allIDs = append(allIDs, ins.ID)
		}
	}

	return allIDs
}

func getInstancesLifeStates(allIns []instances.Instance) []string {
	var allStates = make([]string, len(allIns))
	for i, ins := range allIns {
		allStates[i] = ins.LifeCycleStatus
	}

	return allStates
}

func refreshInstancesLifeStates(asClient *golangsdk.ServiceClient, groupID string, insNum int, checkInService bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		allIns, err := getInstancesInGroup(asClient, groupID, nil)
		if err != nil {
			return nil, "ERROR", err
		}
		// maybe the instances (or some of the instances) have not put in the asg when creating
		if checkInService && len(allIns) != insNum {
			return allIns, "PENDING", err
		}
		allLifeStatus := getInstancesLifeStates(allIns)
		for _, lifeStatus := range allLifeStatus {
			logp.Printf("[DEBUG] Get lifecycle status in group %s: %s", groupID, lifeStatus)
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
		return allIns, "", err
	}
}

func refreshGroupState(client *golangsdk.ServiceClient, groupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		asGroup, err := groups.Get(client, groupID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return asGroup, "DELETED", nil
			}
			return nil, "ERROR", err
		}
		return asGroup, asGroup.Status, nil
	}
}

func checkASGroupInstancesInService(ctx context.Context, client *golangsdk.ServiceClient, groupID string, insNum int, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"INSERVICE"},
		Refresh:      refreshInstancesLifeStates(client, groupID, insNum, true),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func checkASGroupInstancesRemoved(ctx context.Context, client *golangsdk.ServiceClient, groupID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"REMOVING"},
		Target:       []string{""}, // if there is no lifecyclestatus, it means that no instances in AS group
		Refresh:      refreshInstancesLifeStates(client, groupID, 0, false),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func checkASGroupRemoved(ctx context.Context, client *golangsdk.ServiceClient, groupID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"DELETED"},
		Refresh:      refreshGroupState(client, groupID),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func resourceASGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	asClient, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	minNum := d.Get("min_instance_number").(int)
	maxNum := d.Get("max_instance_number").(int)
	var desireNum int
	if v, ok := d.GetOk("desire_instance_number"); ok {
		desireNum = v.(int)
	} else {
		desireNum = minNum
	}
	logp.Printf("[DEBUG] instance number options: min(%d), max(%d), desired(%d)", minNum, maxNum, desireNum)
	if desireNum < minNum || desireNum > maxNum {
		return diag.Errorf("invalid parameters: it should be min_instance_number <= desire_instance_number <= max_instance_number")
	}

	networks := getAllNetworks(d, meta)
	asgNetworks := expandNetworks(networks)

	secGroups := getAllSecurityGroups(d, meta)
	asgSecGroups := expandGroups(secGroups)

	asgLBaaSListeners := getAllLBaaSListeners(d, meta)

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
		EnterpriseProjectID:       common.GetEnterpriseProjectID(d, config),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	asgId, err := groups.Create(asClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating AS group: %s", err)
	}

	d.SetId(asgId)

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := expandGroupsTags(tagRaw)
		if tagErr := tags.Create(asClient, asgId, taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of AS group %s: %s", asgId, tagErr)
		}
	}

	//enable asg
	if d.Get("enable").(bool) {
		enableResult := groups.Enable(asClient, asgId)
		if enableResult.Err != nil {
			return diag.Errorf("error enabling AS group %s: %s", asgId, enableResult.Err)
		}
	}

	// check all instances are inservice
	if desireNum > 0 {
		timeout := d.Timeout(schema.TimeoutCreate)
		err = checkASGroupInstancesInService(ctx, asClient, asgId, desireNum, timeout)
		if err != nil {
			return diag.Errorf("error waiting for instances in the AS group %s to become inservice: %s", asgId, err)
		}
	}

	return resourceASGroupRead(ctx, d, meta)
}

func resourceASGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	asClient, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	asg, err := groups.Get(asClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AS group")
	}
	logp.Printf("[DEBUG] Retrieved AS group %s: %#v", d.Id(), asg)

	// set properties based on the read info
	d.Set("scaling_group_name", asg.Name)
	d.Set("vpc_id", asg.VpcID)
	d.Set("status", asg.Status)

	if asg.Status == "INSERVICE" {
		d.Set("enable", true)
	} else {
		d.Set("enable", false)
	}

	var networks []map[string]interface{}
	for _, network := range asg.Networks {
		mapping := map[string]interface{}{
			"id": network.ID,
		}
		networks = append(networks, mapping)
	}
	d.Set("networks", networks)

	var securityGroups []map[string]interface{}
	for _, securityGroup := range asg.SecurityGroups {
		mapping := map[string]interface{}{
			"id": securityGroup.ID,
		}
		securityGroups = append(securityGroups, mapping)
	}

	d.Set("security_groups", securityGroups)
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
	d.Set("enterprise_project_id", asg.EnterpriseProjectID)
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

	allIns, err := getInstancesInGroup(asClient, d.Id(), nil)
	if err != nil {
		return diag.Errorf("can not get the instances in AS Group %s: %s", d.Id(), err)
	}
	allIDs := getInstancesIDs(allIns)
	d.Set("instances", allIDs)

	d.Set("region", config.GetRegion(d))

	// save group tags
	if resourceTags, err := tags.Get(asClient, d.Id()).Extract(); err == nil {
		tagmap := make(map[string]string)
		for _, val := range resourceTags.Tags {
			tagmap[val.Key] = val.Value
		}
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for AS group (%s): %s", d.Id(), err)
		}
	} else {
		logp.Printf("[WARN] Error fetching tags of AS group (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceASGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	asClient, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
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
		logp.Printf("[DEBUG] instance number options: min(%d), max(%d), desired(%d)", minNum, maxNum, desireNum)
		if desireNum < minNum || desireNum > maxNum {
			return diag.Errorf("invalid parameters: it should be min_instance_number <= desire_instance_number <= max_instance_number")
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
		EnterpriseProjectID:       common.GetEnterpriseProjectID(d, config),
	}

	logp.Printf("[DEBUG] AS Group update options: %#v", updateOpts)
	asgID, err := groups.Update(asClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating AS group %s: %s", asgID, err)
	}

	//update tags
	if d.HasChange("tags") {
		//remove old tags and set new tags
		old, new := d.GetChange("tags")
		oldRaw := old.(map[string]interface{})
		if len(oldRaw) > 0 {
			taglist := expandGroupsTags(oldRaw)
			if tagErr := tags.Delete(asClient, asgID, taglist).ExtractErr(); tagErr != nil {
				return diag.Errorf("error deleting tags of AS group %s: %s", asgID, tagErr)
			}
		}

		newRaw := new.(map[string]interface{})
		if len(newRaw) > 0 {
			taglist := expandGroupsTags(newRaw)
			if tagErr := tags.Create(asClient, asgID, taglist).ExtractErr(); tagErr != nil {
				return diag.Errorf("error setting tags of AS group %s: %s", asgID, tagErr)
			}
		}
	}

	if d.HasChange("enable") {
		if d.Get("enable").(bool) {
			enableResult := groups.Enable(asClient, asgID)
			if enableResult.Err != nil {
				return diag.Errorf("error enabling AS group %s: %s", asgID, enableResult.Err)
			}
			logp.Printf("[DEBUG] Enable AS group %s success", asgID)
		} else {
			enableResult := groups.Disable(asClient, asgID)
			if enableResult.Err != nil {
				return diag.Errorf("error disabling AS group %s: %s", asgID, enableResult.Err)
			}
			logp.Printf("[DEBUG] Disable AS group %s success", asgID)
		}
	}

	return resourceASGroupRead(ctx, d, meta)
}

func resourceASGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	asClient, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	groupID := d.Id()
	logp.Printf("[DEBUG] Begin to delete AS group %s", groupID)

	// forcibly delete an AS group
	if _, ok := d.GetOk("force_delete"); ok {
		if err := groups.ForceDelete(asClient, groupID).ExtractErr(); err != nil {
			return diag.Errorf("error deleting AS group %s: %s", groupID, err)
		}

		err = checkASGroupRemoved(ctx, asClient, groupID, timeout)
		if err != nil {
			return diag.Errorf("error deleting AS group %s: %s", groupID, err)
		}
		return nil
	}

	allIns, err := getInstancesInGroup(asClient, groupID, nil)
	if err != nil {
		return diag.Errorf("error listing instances of AS group: %s", err)
	}
	allIDs := getInstancesIDs(allIns)
	logp.Printf("[DEBUG] Instances in AS group %s: %+v", groupID, allIDs)

	allLifeStatus := getInstancesLifeStates(allIns)
	for _, lifeCycleState := range allLifeStatus {
		if lifeCycleState != "INSERVICE" {
			return diag.Errorf("can't delete the AS group %s: some instances are not in INSERVICE but in %s, "+
				"please try again latter or use force_delete option", groupID, lifeCycleState)
		}
	}

	if len(allIns) > 0 {
		minNumber := d.Get("min_instance_number").(int)
		if minNumber > 0 {
			return diag.Errorf("can't delete the AS group %s: The instance number after the removal will less than "+
				"min number %d, please modify the min number to zero or use force_delete option", groupID, minNumber)
		}

		deleteIns := d.Get("delete_instances").(string)
		logp.Printf("[DEBUG] The flag delete_instances in AS group is %s", deleteIns)
		batchResult := instances.BatchDelete(asClient, groupID, allIDs, deleteIns)
		if batchResult.Err != nil {
			return diag.Errorf("error removing instancess of AS group: %s", batchResult.Err)
		}

		err = checkASGroupInstancesRemoved(ctx, asClient, groupID, timeout)
		if err != nil {
			return diag.Errorf("error removing instances from AS group %s: %s", groupID, err)
		}
	}

	if delErr := groups.Delete(asClient, groupID).ExtractErr(); delErr != nil {
		return diag.Errorf("error deleting AS group: %s", delErr)
	}

	return nil
}

func resourceASGroupValidateListenerId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	split := strings.Split(value, ",")
	if len(split) <= 6 {
		return
	}
	errors = append(errors, fmtp.Errorf("%s supports binding up to 6 ELB listeners which are separated by a comma.", k))
	return
}
