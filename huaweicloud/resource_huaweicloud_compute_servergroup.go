package huaweicloud

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/servergroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceComputeServerGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeServerGroupCreate,
		Read:   resourceComputeServerGroupRead,
		Update: resourceComputeServerGroupUpdate,
		Delete: resourceComputeServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"policies": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"members": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"fault_domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceComputeServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	ecsClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	createOpts := servergroups.CreateOpts{
		Name:     d.Get("name").(string),
		Policies: resourceServerGroupPolicies(d),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	newSG, err := servergroups.Create(ecsClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating ServerGroup: %s", err)
	}

	d.SetId(newSG.ID)

	membersToAdd := d.Get("members").(*schema.Set)
	for _, v := range membersToAdd.List() {
		var addMemberOpts servergroups.MemberOpts
		addMemberOpts.InstanceID = v.(string)
		if err := servergroups.UpdateMember(ecsClient, addMemberOpts, "add_member", d.Id()).ExtractErr(); err != nil {
			return fmtp.Errorf("Error to add an instance to ECS server group, err=%s", err)
		}
	}

	return resourceComputeServerGroupRead(d, meta)
}

func resourceComputeServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	ecsClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	sg, err := servergroups.Get(ecsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "server group")
	}

	logp.Printf("[DEBUG] Retrieved ServerGroup %s: %+v", d.Id(), sg)

	policies := make([]string, len(sg.Policies))
	for i, p := range sg.Policies {
		policies[i] = p
	}
	d.Set("policies", policies)
	d.Set("name", sg.Name)
	d.Set("members", sg.Members)
	d.Set("fault_domains", sg.FaultDomain.Names)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	ecsClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute V1 client: %s", err)
	}

	if d.HasChange("members") {
		oldMembers, newMembers := d.GetChange("members")
		oldMemberSet, newMemberSet := oldMembers.(*schema.Set), newMembers.(*schema.Set)
		membersToAdd := newMemberSet.Difference(oldMemberSet)
		membersToRemove := oldMemberSet.Difference(newMemberSet)

		for _, v := range membersToAdd.List() {
			var addMemberOpts servergroups.MemberOpts
			addMemberOpts.InstanceID = v.(string)
			if err := servergroups.UpdateMember(ecsClient, addMemberOpts, "add_member", d.Id()).ExtractErr(); err != nil {
				return fmtp.Errorf("Error to add a instance to ECS server group, err=%s", err)
			}
		}

		for _, v := range membersToRemove.List() {
			instanceID := v.(string)
			server, err := cloudservers.Get(ecsClient, instanceID).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					logp.Printf("[WARN] the compute %s is not exist, ignore to remove it from the group", instanceID)
					continue
				}
				logp.Printf("[WARN] failed to retrieve compute %s: %s, try to remove it from the group", instanceID, err)
			} else {
				if server.Status == "DELETED" {
					logp.Printf("[WARN] the compute %s was removed, ignore to remove it from the group", instanceID)
					continue
				}
			}

			var removeMemberOpts servergroups.MemberOpts
			removeMemberOpts.InstanceID = instanceID
			if err := servergroups.UpdateMember(ecsClient, removeMemberOpts, "remove_member", d.Id()).ExtractErr(); err != nil {
				return fmtp.Errorf("Error to remove a instance from ECS server group, err=%s", err)
			}
		}
	}

	return resourceComputeServerGroupRead(d, meta)
}

func resourceComputeServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	ecsClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	logp.Printf("[DEBUG] Deleting ServerGroup %s", d.Id())
	if err := servergroups.Delete(ecsClient, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Error deleting ServerGroup: %s", err)
	}

	return nil
}

func resourceServerGroupPolicies(d *schema.ResourceData) []string {
	rawPolicies := d.Get("policies").([]interface{})
	policies := make([]string, len(rawPolicies))
	for i, raw := range rawPolicies {
		policies[i] = raw.(string)
	}
	return policies
}
