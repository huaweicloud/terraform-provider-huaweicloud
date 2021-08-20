package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/servergroups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceComputeServerGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeServerGroupV2Create,
		Read:   resourceComputeServerGroupV2Read,
		Update: resourceComputeServerGroupV2Update,
		Delete: resourceComputeServerGroupV2Delete,
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
				ForceNew: true,
				Required: true,
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

func resourceComputeServerGroupV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	createOpts := servergroups.CreateOpts{
		Name:     d.Get("name").(string),
		Policies: resourceServerGroupPoliciesV2(d),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	newSG, err := servergroups.Create(computeClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating ServerGroup: %s", err)
	}

	d.SetId(newSG.ID)

	clv1, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute V1 client: %s", err)
	}
	membersToAdd := d.Get("members").(*schema.Set)
	for _, v := range membersToAdd.List() {
		var addMemberOpts servergroups.MemberOpts
		addMemberOpts.InstanceUUid = v.(string)
		if err := servergroups.UpdateMember(clv1, addMemberOpts, "add_member", d.Id()).ExtractErr(); err != nil {
			return fmtp.Errorf("Error to add a instance to ECS server group, err=%s", err)
		}
	}

	return resourceComputeServerGroupV2Read(d, meta)
}

func resourceComputeServerGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	sg, err := servergroups.Get(computeClient, d.Id()).Extract()
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

func resourceComputeServerGroupV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	clv1, err := config.ComputeV1Client(GetRegion(d, config))
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
			addMemberOpts.InstanceUUid = v.(string)
			if err := servergroups.UpdateMember(clv1, addMemberOpts, "add_member", d.Id()).ExtractErr(); err != nil {
				return fmtp.Errorf("Error to add a instance to ECS server group, err=%s", err)
			}
		}

		for _, v := range membersToRemove.List() {
			var removeMemberOpts servergroups.MemberOpts
			removeMemberOpts.InstanceUUid = v.(string)
			if err := servergroups.UpdateMember(clv1, removeMemberOpts, "remove_member", d.Id()).ExtractErr(); err != nil {
				return fmtp.Errorf("Error to remove a instance from ECS server group, err=%s", err)
			}
		}
	}

	return resourceComputeServerGroupV2Read(d, meta)
}

func resourceComputeServerGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	logp.Printf("[DEBUG] Deleting ServerGroup %s", d.Id())
	if err := servergroups.Delete(computeClient, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Error deleting ServerGroup: %s", err)
	}

	return nil
}

func resourceServerGroupPoliciesV2(d *schema.ResourceData) []string {
	rawPolicies := d.Get("policies").([]interface{})
	policies := make([]string, len(rawPolicies))
	for i, raw := range rawPolicies {
		policies[i] = raw.(string)
	}
	return policies
}
