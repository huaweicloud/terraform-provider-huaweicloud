package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/policies"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/routerinsertion"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
)

func ResourceNetworkACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkACLCreate,
		Read:   resourceNetworkACLRead,
		Update: resourceNetworkACLUpdate,
		Delete: resourceNetworkACLDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inbound_rules": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"outbound_rules": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnets": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"inbound_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"outbound_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ports": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkACLCreate(d *schema.ResourceData, meta interface{}) error {
	var err error
	var portIds []string
	var inboundPolicyID, outboundPolicyID string

	config := meta.(*Config)
	fwClient, err := config.FwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	defer func() {
		// delete firewall policies when encounter errors
		if err != nil && inboundPolicyID != "" {
			deleteErr := policies.Delete(fwClient, inboundPolicyID).Err
			if deleteErr != nil {
				log.Printf("[WARN] Error deleting inbound firewall policy %s: %s", inboundPolicyID, deleteErr)
			}
		}

		if err != nil && outboundPolicyID != "" {
			deleteErr := policies.Delete(fwClient, outboundPolicyID).Err
			if deleteErr != nil {
				log.Printf("[WARN] Error deleting outbound firewall policy %s: %s", outboundPolicyID, deleteErr)
			}
		}
	}()

	// get port Ids from subnets
	subnetsRaw := d.Get("subnets").(*schema.Set).List()
	if len(subnetsRaw) > 0 {
		for _, v := range subnetsRaw {
			port, err := getGWPortFromSubnet(config, v.(string))
			if err != nil {
				return err
			}
			portIds = append(portIds, port)
		}
		log.Printf("[DEBUG] Will attempt to associate Firewall group with subnets: %+v", subnetsRaw)
	}

	groupName := d.Get("name").(string)
	// create inbound policy
	inboundRaw := d.Get("inbound_rules").([]interface{})
	if len(inboundRaw) > 0 {
		inboundRules := make([]string, len(inboundRaw))
		for i, r := range inboundRaw {
			inboundRules[i] = r.(string)
		}

		policyName := "inbound_policy_for_" + groupName
		policyOpts := policies.CreateOpts{
			Name:  policyName,
			Rules: inboundRules,
		}

		log.Printf("[DEBUG] Create inbound firewall policy: %#v", policyOpts)
		policy, err := policies.Create(fwClient, policyOpts).Extract()
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Firewall inbound policy created: %#v", policy)
		inboundPolicyID = policy.ID
	}

	// create outbound policy
	outboundRaw := d.Get("outbound_rules").([]interface{})
	if len(outboundRaw) > 0 {
		outboundRules := make([]string, len(outboundRaw))
		for i, r := range outboundRaw {
			outboundRules[i] = r.(string)
		}

		policyName := "outbound_policy_for_" + groupName
		policyOpts := policies.CreateOpts{
			Name:  policyName,
			Rules: outboundRules,
		}

		log.Printf("[DEBUG] Create outbound firewall policy: %#v", policyOpts)
		policy, err := policies.Create(fwClient, policyOpts).Extract()
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Firewall outbound policy created: %#v", policy)
		outboundPolicyID = policy.ID
	}

	var createOpts firewall_groups.CreateOptsBuilder
	createOpts = &firewall_groups.CreateOpts{
		Name:            groupName,
		Description:     d.Get("description").(string),
		IngressPolicyID: inboundPolicyID,
		EgressPolicyID:  outboundPolicyID,
	}

	if len(portIds) > 0 {
		createOpts = &routerinsertion.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			PortIDs:           portIds,
		}
	}

	log.Printf("[DEBUG] Create firewall group: %#v", createOpts)
	group, err := firewall_groups.Create(fwClient, createOpts).Extract()
	if err != nil {
		return err
	}

	d.SetId(group.ID)
	log.Printf("[DEBUG] waiting for Firewall group (%s) to become ACTIVE", d.Id())

	stateConf := &resource.StateChangeConf{
		// if none subnets was associated with the firewall group, the state will be "INACTIVE"
		// so we seems the "INACTIVE" as a target state.
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE", "INACTIVE"},
		Refresh:    waitForFirewallGroupActive(fwClient, group.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      2,
		MinTimeout: 2 * time.Second,
	}
	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf("Error waiting for Firewall group (%s) to become ACTIVE: %s",
			d.Id(), stateErr)
	}

	log.Printf("[DEBUG] Firewall group (%s) is active.", group.ID)
	return resourceNetworkACLRead(d, meta)
}

func resourceNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	fwClient, err := config.FwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	var fwGroup FirewallGroup
	err = firewall_groups.Get(fwClient, d.Id()).ExtractInto(&fwGroup)
	if err != nil {
		return CheckDeleted(d, err, "firewall")
	}

	log.Printf("[DEBUG] Read HuaweiCloud Firewall group %s: %#v", d.Id(), fwGroup)

	d.Set("name", fwGroup.Name)
	d.Set("status", fwGroup.Status)
	d.Set("description", fwGroup.Description)
	d.Set("inbound_policy_id", fwGroup.IngressPolicyID)
	d.Set("outbound_policy_id", fwGroup.EgressPolicyID)
	if err := d.Set("ports", fwGroup.PortIDs); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ports to state for HuaweiCloud firewall group (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceNetworkACLUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	fwClient, err := config.FwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	// first of all, inbound_policy/rules and outbound_policy/rules should be updated
	if d.HasChange("inbound_rules") {
		err := updateNetworkACLPolicyRules(d, fwClient, "inbind")
		if err != nil {
			return err
		}
	}
	if d.HasChange("outbound_rules") {
		err := updateNetworkACLPolicyRules(d, fwClient, "outbind")
		if err != nil {
			return err
		}
	}

	// update other parameters
	var changed bool
	opts := firewall_groups.UpdateOpts{}
	if d.HasChanges("name", "description") {
		changed = true
		opts.Name = d.Get("name").(string)
		opts.Description = d.Get("description").(string)
	}

	var updateOpts firewall_groups.UpdateOptsBuilder
	var portIds []string
	if d.HasChange("subnets") {
		changed = true

		// get port Ids from subnets
		subnetsRaw := d.Get("subnets").(*schema.Set).List()
		for _, v := range subnetsRaw {
			port, err := getGWPortFromSubnet(config, v.(string))
			if err != nil {
				return err
			}
			portIds = append(portIds, port)
		}
		log.Printf("[DEBUG] Will attempt to associate Firewall group with subnets: %+v", subnetsRaw)

		updateOpts = routerinsertion.UpdateOptsExt{
			UpdateOptsBuilder: opts,
			PortIDs:           portIds,
		}
	} else {
		updateOpts = opts
	}

	if changed {
		log.Printf("[DEBUG] Updating firewall with id %s: %#v", d.Id(), updateOpts)
		err = firewall_groups.Update(fwClient, d.Id(), updateOpts).Err
		if err != nil {
			return err
		}

		// if none subnets was associated with the firewall group, the state will be "INACTIVE"
		// so we seems the "INACTIVE" as a target state.
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING_CREATE", "PENDING_UPDATE"},
			Target:     []string{"ACTIVE", "INACTIVE"},
			Refresh:    waitForFirewallGroupActive(fwClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      2,
			MinTimeout: 2 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error updating firewall group (%s): %s", d.Id(), err)
		}
	}

	return resourceNetworkACLRead(d, meta)
}

func resourceNetworkACLDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy firewall group: %s", d.Id())

	config := meta.(*Config)
	fwClient, err := config.FwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	inboundPolicyID := d.Get("inbound_policy_id").(string)
	outboundPolicyID := d.Get("outbound_policy_id").(string)

	err = firewall_groups.Delete(fwClient, d.Id()).Err
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForFirewallGroupDeletion(fwClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      2,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting firewall group (%s): %s", d.Id(), err)
	}

	// delete firewall policies after the firewall group
	if inboundPolicyID != "" {
		deleteErr := policies.Delete(fwClient, inboundPolicyID).Err
		if deleteErr != nil {
			log.Printf("[WARN] Error deleting inbound firewall policy %s: %s", inboundPolicyID, deleteErr)
		}
	}

	if outboundPolicyID != "" {
		deleteErr := policies.Delete(fwClient, outboundPolicyID).Err
		if deleteErr != nil {
			log.Printf("[WARN] Error deleting outbound firewall policy %s: %s", outboundPolicyID, deleteErr)
		}
	}

	d.SetId("")
	return nil
}

func getGWPortFromSubnet(config *Config, subnetID string) (string, error) {
	var gatewayIP string
	var gatewayPort string

	subnetClient, err := config.NetworkingV1Client(config.Region)
	if err != nil {
		return "", fmt.Errorf("Error creating Huaweicloud vpc client: %s", err)
	}
	networkingClient, err := config.NetworkingV2Client(config.Region)
	if err != nil {
		return "", fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	// get Gateway IP
	n, err := subnets.Get(subnetClient, subnetID).Extract()
	if err != nil {
		return "", fmt.Errorf("Error retrieving Huaweicloud subnet %s: %s", subnetID, err)
	}
	gatewayIP = n.GatewayIP
	log.Printf("[DEBUG] the gateway IP address of subnet %s is %s", subnetID, gatewayIP)

	// list all ports in the subnet
	listOpts := ports.ListOpts{
		NetworkID: subnetID,
		//Status:    "ACTIVE",
	}
	allPages, err := ports.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return "", fmt.Errorf("Unable to list Huaweicloud ports of %s: %s", subnetID, err)
	}

	var allPorts []ports.Port
	err = ports.ExtractPortsInto(allPages, &allPorts)
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve Huaweicloud ports of %s: %s", subnetID, err)
	}

	if len(allPorts) == 0 {
		return "", fmt.Errorf("No ports was found in %s", subnetID)
	}

	// Filter IPs by the gatewayIP
	var isExist bool
	for _, p := range allPorts {
		for _, ipObject := range p.FixedIPs {
			if ipObject.IPAddress == gatewayIP {
				isExist = true
				gatewayPort = p.ID
				log.Printf("[DEBUG] the gateway port of subnet %s is %s", subnetID, gatewayPort)
				break
			}
		}
		if isExist {
			break
		}
	}
	if !isExist {
		return "", fmt.Errorf("No gateway port was found in %s", subnetID)
	}

	return gatewayPort, nil
}

func updateNetworkACLPolicyRules(d *schema.ResourceData, client *golangsdk.ServiceClient, direct string) error {
	var policyKey, rulesKey string

	if direct == "inbind" {
		policyKey = "inbound_policy_id"
		rulesKey = "inbound_rules"
	} else {
		policyKey = "outbound_policy_id"
		rulesKey = "outbound_rules"
	}

	groupName := d.Get("name").(string)
	policyID := d.Get(policyKey).(string)
	policyName := direct + "_policy_for_" + groupName

	v := d.Get(rulesKey).([]interface{})
	rulesList := make([]string, len(v))
	for i, v := range v {
		rulesList[i] = v.(string)
	}

	if policyID != "" {
		// update the firewall policy, even if the rules are empty
		policyOpts := policies.UpdateOpts{
			Name:  policyName,
			Rules: rulesList,
		}

		log.Printf("[DEBUG] updating firewall policy with id %s: %#v", policyID, policyOpts)
		err := policies.Update(client, policyID, policyOpts).Err
		if err != nil {
			return fmt.Errorf("Error updating firewall policy %s: %s", policyID, err)
		}
	} else {
		// create new firewall policy
		policyOpts := policies.CreateOpts{
			Name:  policyName,
			Rules: rulesList,
		}

		log.Printf("[DEBUG] Create firewall policy: %#v", policyOpts)
		policy, err := policies.Create(client, policyOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error creating firewall policy: %s", err)
		}

		d.Set(policyKey, policy.ID)
	}

	return nil
}
