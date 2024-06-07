package as

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/groups"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/instances"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	HealthAuditMethods = []string{"ELB_AUDIT", "NOVA_AUDIT"}
	HealthAuditTime    = []int{0, 1, 5, 15, 60, 180}
	TerminatePolices   = []string{"OLD_CONFIG_OLD_INSTANCE", "OLD_CONFIG_NEW_INSTANCE", "OLD_INSTANCE", "NEW_INSTANCE"}
)

// @API AS GET /autoscaling-api/v1/{project_id}/scaling_group/{id}
// @API AS PUT /autoscaling-api/v1/{project_id}/scaling_group/{id}
// @API AS DELETE /autoscaling-api/v1/{project_id}/scaling_group/{id}
// @API AS POST /autoscaling-api/v1/{project_id}/scaling_group_instance/{groupID}/action
// @API AS GET /autoscaling-api/v1/{project_id}/scaling_group_instance/{groupID}/list
// @API AS POST /autoscaling-api/v1/{project_id}/scaling_group_tag/{id}/tags/action
// @API AS GET /autoscaling-api/v1/{project_id}/scaling_group_tag/{id}/tags
// @API AS POST /autoscaling-api/v1/{project_id}/scaling_group
// @API AS POST /autoscaling-api/v1/{project_id}/scaling_group/{id}/action
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
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
						"ipv6_enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ipv6_bandwidth_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source_dest_check": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"security_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"multi_az_scaling_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"health_periodic_audit_grace_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 86400),
				Description:  "The health check grace period for instances, in seconds.",
			},
			"instance_terminate_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "OLD_CONFIG_OLD_INSTANCE",
				ValidateFunc: validation.StringInSlice(TerminatePolices, false),
			},
			"agency_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"available_zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "schema: Deprecated; use availability_zones instead",
			},
			"notifications": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "schema: Deprecated; The notification mode has been canceled",
			},
		},
	}
}

func buildNetworksOpts(networks []interface{}) []groups.NetworkOpts {
	res := make([]groups.NetworkOpts, len(networks))
	for i, v := range networks {
		item := v.(map[string]interface{})
		res[i] = groups.NetworkOpts{
			ID:         item["id"].(string),
			IPv6Enable: item["ipv6_enable"].(bool),
		}

		if id, ok := item["ipv6_bandwidth_id"]; ok && id.(string) != "" {
			res[i].IPv6BandWidth = &groups.BandWidthOpts{
				ID: id.(string),
			}
		}

		if item["source_dest_check"].(bool) {
			// Cancle all allowed-address-pairs to enable the source/destination check
			res[i].AllowedAddressPairs = make([]groups.AddressPairOpts, 0)
		} else {
			// Update the allowed-address-pairs to 1.1.1.1/0
			// to disable the source/destination check
			addressPairs := groups.AddressPairOpts{
				IpAddress: "1.1.1.1/0",
			}
			res[i].AllowedAddressPairs = []groups.AddressPairOpts{addressPairs}
		}
	}

	return res
}

func buildSecurityGroupsOpts(secGroups []interface{}) []groups.SecurityGroupOpts {
	if len(secGroups) == 0 {
		return nil
	}

	res := make([]groups.SecurityGroupOpts, len(secGroups))
	for i, v := range secGroups {
		item := v.(map[string]interface{})
		res[i] = groups.SecurityGroupOpts{
			ID: item["id"].(string),
		}
	}

	return res
}

func buildLBaaSListenersOpts(listeners []interface{}) []groups.LBaaSListenerOpts {
	if len(listeners) == 0 {
		return nil
	}

	res := make([]groups.LBaaSListenerOpts, len(listeners))
	for i, v := range listeners {
		item := v.(map[string]interface{})
		res[i] = groups.LBaaSListenerOpts{
			PoolID:       item["pool_id"].(string),
			ProtocolPort: item["protocol_port"].(int),
			Weight:       item["weight"].(int),
		}
	}

	return res
}

func buildAvailabilityZonesOpts(d *schema.ResourceData) []string {
	var rawZones []interface{}
	v1, ok1 := d.GetOk("availability_zones")
	v2, ok2 := d.GetOk("available_zones")

	if ok1 {
		rawZones = v1.([]interface{})
	} else if ok2 {
		rawZones = v2.([]interface{})
	}

	zones := make([]string, len(rawZones))
	for i, raw := range rawZones {
		zones[i] = raw.(string)
	}

	return zones
}

func expandGroupsTags(tagmap map[string]interface{}) []tags.ResourceTag {
	taglist := make([]tags.ResourceTag, 0, len(tagmap))
	for k, v := range tagmap {
		tag := tags.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	return taglist
}

func getInstancesInGroup(asClient *golangsdk.ServiceClient, groupID string, opts instances.ListOptsBuilder) ([]instances.Instance, error) {
	var insList []instances.Instance
	page, err := instances.List(asClient, groupID, opts).AllPages()
	if err != nil {
		return insList, fmt.Errorf("error getting instances in AS group %s: %s", groupID, err)
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
		// maybe the instances (or some of the instances) have not put in the autoscaling group when creating
		if checkInService && len(allIns) != insNum {
			return allIns, "PENDING", err
		}
		allLifeStatus := getInstancesLifeStates(allIns)
		for _, lifeStatus := range allLifeStatus {
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
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(parseGroupResponseError(err), &errDefault404) {
				return asGroup, "DELETED", nil
			}
			return nil, "ERROR", err
		}
		return asGroup, asGroup.Status, nil
	}
}

// When the AS group does not exist, the response body example of the details interface is as follows:
// {"error":{"code":"AS.2007","message":"The AS group does not exist."}}
func parseGroupResponseError(err error) error {
	var errCode golangsdk.ErrDefault400
	if errors.As(err, &errCode) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return err
		}
		errorCode, errorCodeErr := jmespath.Search("error.code", apiError)
		if errorCodeErr != nil || errorCode == nil {
			return err
		}

		if errorCode.(string) == "AS.2007" {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return err
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
	conf := meta.(*config.Config)
	asClient, err := conf.AutoscalingV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	minNum := d.Get("min_instance_number").(int)
	maxNum := d.Get("max_instance_number").(int)
	desireNum := minNum
	if v, ok := d.GetOk("desire_instance_number"); ok {
		desireNum = v.(int)
	}
	log.Printf("[DEBUG] instance number options: min(%d), max(%d), desired(%d)", minNum, maxNum, desireNum)
	if desireNum < minNum || desireNum > maxNum {
		return diag.Errorf("invalid parameters: it should be min_instance_number <= desire_instance_number <= max_instance_number")
	}

	createOpts := groups.CreateOpts{
		Name:                      d.Get("scaling_group_name").(string),
		ConfigurationID:           d.Get("scaling_configuration_id").(string),
		DesireInstanceNumber:      desireNum,
		MinInstanceNumber:         minNum,
		MaxInstanceNumber:         maxNum,
		CoolDownTime:              d.Get("cool_down_time").(int),
		LBListenerID:              d.Get("lb_listener_id").(string),
		AvailableZones:            buildAvailabilityZonesOpts(d),
		LBaaSListeners:            buildLBaaSListenersOpts(d.Get("lbaas_listeners").([]interface{})),
		Networks:                  buildNetworksOpts(d.Get("networks").([]interface{})),
		SecurityGroup:             buildSecurityGroupsOpts(d.Get("security_groups").([]interface{})),
		Notifications:             utils.ExpandToStringList(d.Get("notifications").([]interface{})),
		VpcID:                     d.Get("vpc_id").(string),
		HealthPeriodicAuditMethod: d.Get("health_periodic_audit_method").(string),
		HealthPeriodicAuditTime:   d.Get("health_periodic_audit_time").(int),
		HealthPeriodicAuditGrace:  d.Get("health_periodic_audit_grace_period").(int),
		InstanceTerminatePolicy:   d.Get("instance_terminate_policy").(string),
		MultiAZPriorityPolicy:     d.Get("multi_az_scaling_policy").(string),
		Description:               d.Get("description").(string),
		IamAgencyName:             d.Get("agency_name").(string),
		IsDeletePublicip:          d.Get("delete_publicip").(bool),
		EnterpriseProjectID:       common.GetEnterpriseProjectID(d, conf),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	asgId, err := groups.Create(asClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating AS group: %s", err)
	}

	d.SetId(asgId)

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := expandGroupsTags(tagRaw)
		if tagErr := tags.Create(asClient, asgId, taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of AS group %s: %s", asgId, tagErr)
		}
	}

	// the autoscaling group is disabled after creating
	if d.Get("enable").(bool) {
		enableResult := groups.Enable(asClient, asgId)
		if enableResult.Err != nil {
			return diag.Errorf("error enabling AS group %s: %s", asgId, enableResult.Err)
		}
	}

	// check all instances are inservice
	if desireNum > 0 {
		err = checkASGroupInstancesInService(ctx, asClient, asgId, desireNum, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for instances in the AS group %s to become inservice: %s", asgId, err)
		}
	}

	return resourceASGroupRead(ctx, d, meta)
}

func resourceASGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	groupID := d.Id()
	asg, err := groups.Get(asClient, groupID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, parseGroupResponseError(err), "AS group")
	}

	log.Printf("[DEBUG] Retrieved AS group %s: %#v", groupID, asg)
	allIns, err := getInstancesInGroup(asClient, groupID, nil)
	if err != nil {
		return diag.Errorf("can not get the instances in AS Group %s: %s", groupID, err)
	}
	allIDs := getInstancesIDs(allIns)

	// set properties based on the read info
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("scaling_group_name", asg.Name),
		d.Set("scaling_configuration_id", asg.ConfigurationID),
		d.Set("vpc_id", asg.VpcID),
		d.Set("status", asg.Status),
		d.Set("enable", asg.Status == "INSERVICE"),
		d.Set("current_instance_number", asg.ActualInstanceNumber),
		d.Set("desire_instance_number", asg.DesireInstanceNumber),
		d.Set("min_instance_number", asg.MinInstanceNumber),
		d.Set("max_instance_number", asg.MaxInstanceNumber),
		d.Set("cool_down_time", asg.CoolDownTime),
		d.Set("lb_listener_id", asg.LBListenerID),
		d.Set("health_periodic_audit_method", asg.HealthPeriodicAuditMethod),
		d.Set("health_periodic_audit_time", asg.HealthPeriodicAuditTime),
		d.Set("health_periodic_audit_grace_period", asg.HealthPeriodicAuditGrace),
		d.Set("instance_terminate_policy", asg.InstanceTerminatePolicy),
		d.Set("delete_publicip", asg.DeletePublicip),
		d.Set("enterprise_project_id", asg.EnterpriseProjectID),
		d.Set("availability_zones", asg.AvailableZones),
		d.Set("multi_az_scaling_policy", asg.MultiAZPriorityPolicy),
		d.Set("description", asg.Description),
		d.Set("agency_name", asg.IamAgencyName),
		d.Set("notifications", asg.Notifications),
		d.Set("instances", allIDs),
		d.Set("networks", flattenNetworks(asg.Networks)),
		d.Set("security_groups", flattenSecurityGroups(asg.SecurityGroups)),
		d.Set("lbaas_listeners", flattenLBaaSListeners(asg.LBaaSListeners)),
	)

	// save group tags
	if resourceTags, err := tags.Get(asClient, groupID).Extract(); err == nil {
		tagmap := make(map[string]string)
		for _, val := range resourceTags.Tags {
			tagmap[val.Key] = val.Value
		}
		mErr = multierror.Append(mErr, d.Set("tags", tagmap))
	} else {
		log.Printf("[WARN] Error fetching tags of AS group (%s): %s", groupID, err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNetworks(networks []groups.Network) []map[string]interface{} {
	res := make([]map[string]interface{}, len(networks))
	for i, item := range networks {
		res[i] = map[string]interface{}{
			"id":                item.ID,
			"ipv6_enable":       item.IPv6Enable,
			"ipv6_bandwidth_id": item.IPv6BandWidth.ID,
			"source_dest_check": len(item.AllowedAddressPairs) == 0,
		}
	}
	return res
}

func flattenSecurityGroups(sgs []groups.SecurityGroup) []map[string]interface{} {
	res := make([]map[string]interface{}, len(sgs))
	for i, item := range sgs {
		res[i] = map[string]interface{}{
			"id": item.ID,
		}
	}
	return res
}

func flattenLBaaSListeners(listeners []groups.LBaaSListener) []map[string]interface{} {
	if len(listeners) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(listeners))
	for i, item := range listeners {
		res[i] = map[string]interface{}{
			"pool_id":       item.PoolID,
			"protocol_port": item.ProtocolPort,
			"weight":        item.Weight,
		}
	}
	return res
}

func resourceASGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	minNum := d.Get("min_instance_number").(int)
	maxNum := d.Get("max_instance_number").(int)
	desireNum := minNum
	if v, ok := d.GetOk("desire_instance_number"); ok {
		desireNum = v.(int)
	}
	if d.HasChanges("min_instance_number", "max_instance_number", "desire_instance_number") {
		log.Printf("[DEBUG] instance number options: min(%d), max(%d), desired(%d)", minNum, maxNum, desireNum)
		if desireNum < minNum || desireNum > maxNum {
			return diag.Errorf("invalid parameters: it should be min_instance_number <= desire_instance_number <= max_instance_number")
		}
	}

	updateOpts := groups.UpdateOpts{
		Name:                      d.Get("scaling_group_name").(string),
		ConfigurationID:           d.Get("scaling_configuration_id").(string),
		DesireInstanceNumber:      desireNum,
		MinInstanceNumber:         minNum,
		MaxInstanceNumber:         maxNum,
		CoolDownTime:              d.Get("cool_down_time").(int),
		LBListenerID:              d.Get("lb_listener_id").(string),
		AvailableZones:            buildAvailabilityZonesOpts(d),
		LBaaSListeners:            buildLBaaSListenersOpts(d.Get("lbaas_listeners").([]interface{})),
		Networks:                  buildNetworksOpts(d.Get("networks").([]interface{})),
		SecurityGroup:             buildSecurityGroupsOpts(d.Get("security_groups").([]interface{})),
		Notifications:             utils.ExpandToStringList(d.Get("notifications").([]interface{})),
		HealthPeriodicAuditMethod: d.Get("health_periodic_audit_method").(string),
		HealthPeriodicAuditTime:   d.Get("health_periodic_audit_time").(int),
		HealthPeriodicAuditGrace:  d.Get("health_periodic_audit_grace_period").(int),
		InstanceTerminatePolicy:   d.Get("instance_terminate_policy").(string),
		MultiAZPriorityPolicy:     d.Get("multi_az_scaling_policy").(string),
		IsDeletePublicip:          d.Get("delete_publicip").(bool),
		Description:               utils.String(d.Get("description").(string)),
		EnterpriseProjectID:       common.GetEnterpriseProjectID(d, conf),
	}

	if d.HasChange("agency_name") {
		updateOpts.IamAgencyName = d.Get("agency_name").(string)
	}

	log.Printf("[DEBUG] AS Group update options: %#v", updateOpts)
	asgID, err := groups.Update(asClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating AS group %s: %s", asgID, err)
	}

	// update tags
	if d.HasChange("tags") {
		// remove oldTag tags and set newTag tags
		oldTag, newTag := d.GetChange("tags")
		oldRaw := oldTag.(map[string]interface{})
		if len(oldRaw) > 0 {
			taglist := expandGroupsTags(oldRaw)
			if tagErr := tags.Delete(asClient, asgID, taglist).ExtractErr(); tagErr != nil {
				return diag.Errorf("error deleting tags of AS group %s: %s", asgID, tagErr)
			}
		}

		newRaw := newTag.(map[string]interface{})
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
			log.Printf("[DEBUG] Enable AS group %s success", asgID)
		} else {
			enableResult := groups.Disable(asClient, asgID)
			if enableResult.Err != nil {
				return diag.Errorf("error disabling AS group %s: %s", asgID, enableResult.Err)
			}
			log.Printf("[DEBUG] Disable AS group %s success", asgID)
		}
	}

	return resourceASGroupRead(ctx, d, meta)
}

func resourceASGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	asClient, err := conf.AutoscalingV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	groupID := d.Id()
	timeout := d.Timeout(schema.TimeoutDelete)

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
	log.Printf("[DEBUG] Instances in AS group %s: %+v", groupID, allIDs)

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
	errors = append(errors, fmt.Errorf("%s supports binding up to 6 ELB listeners which are separated by a comma", k))
	return
}
