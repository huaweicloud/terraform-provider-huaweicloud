package ecs

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/servergroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API ECS POST /v1/{project_id}/cloudservers/os-server-groups
// @API ECS POST /v1/{project_id}/cloudservers/os-server-groups/{id}/action
// @API ECS GET /v1/{project_id}/cloudservers/os-server-groups/{id}
// @API ECS DELETE /v1/{project_id}/cloudservers/os-server-groups/{id}
// @API ECS GET /v1/{project_id}/cloudservers/{server_id}
func ResourceComputeServerGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeServerGroupCreate,
		ReadContext:   resourceComputeServerGroupRead,
		UpdateContext: resourceComputeServerGroupUpdate,
		DeleteContext: resourceComputeServerGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "schema: Required",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"members": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"fault_domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "schema: Internal",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildServerGroupPolicies(d *schema.ResourceData) []string {
	rawPolicies := d.Get("policies").([]interface{})
	policies := make([]string, len(rawPolicies))
	for i, raw := range rawPolicies {
		policies[i] = raw.(string)
	}
	return policies
}

func resourceComputeServerGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ecsClient, err := cfg.ComputeV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	createOpts := servergroups.CreateOpts{
		Name:     d.Get("name").(string),
		Policies: buildServerGroupPolicies(d),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	newSG, err := servergroups.Create(ecsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating ECS server group: %s", err)
	}

	d.SetId(newSG.ID)

	membersToAdd := d.Get("members").(*schema.Set)
	for _, v := range membersToAdd.List() {
		instanceId := v.(string)
		err := addServerGroupMember(ecsClient, d.Id(), instanceId)
		if err != nil {
			return diag.Errorf("error binding instance %s to ECS server group: %s", instanceId, err)
		}
	}

	return resourceComputeServerGroupRead(ctx, d, meta)
}

func resourceComputeServerGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	groupID := d.Id()
	sg, err := servergroups.Get(ecsClient, groupID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "server group")
	}

	log.Printf("[DEBUG] Retrieved server group %s: %+v", groupID, sg)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", sg.Name),
		d.Set("members", sg.Members),
		d.Set("policies", sg.Policies),
		d.Set("fault_domains", sg.FaultDomain.Names),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting server group fields: %s", err)
	}
	return nil
}

func resourceComputeServerGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ecsClient, err := cfg.ComputeV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	if d.HasChange("members") {
		oldMembers, newMembers := d.GetChange("members")
		oldMemberSet, newMemberSet := oldMembers.(*schema.Set), newMembers.(*schema.Set)
		membersToAdd := newMemberSet.Difference(oldMemberSet)
		membersToRemove := oldMemberSet.Difference(newMemberSet)

		for _, v := range membersToRemove.List() {
			instanceId := v.(string)
			err := removeServerGroupMember(ecsClient, d.Id(), instanceId)
			if err != nil {
				return diag.Errorf("error unbinding instance %s from ECS server group: %s", instanceId, err)
			}
		}

		for _, v := range membersToAdd.List() {
			instanceId := v.(string)
			err := addServerGroupMember(ecsClient, d.Id(), instanceId)
			if err != nil {
				return diag.Errorf("error binding instance %s to server group: %s", instanceId, err)
			}
		}
	}

	return resourceComputeServerGroupRead(ctx, d, meta)
}

func LockAll(ids []interface{}) {
	for _, instanceId := range ids {
		config.MutexKV.Lock(instanceId.(string))
	}
}

func UnlockAll(ids []interface{}) {
	for _, instanceId := range ids {
		config.MutexKV.Unlock(instanceId.(string))
	}
}

func resourceComputeServerGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ecsClient, err := cfg.ComputeV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	members := d.Get("members").(*schema.Set).List()
	// Make sure that no other operations on the ECS instance are performed during the unbinding process.
	LockAll(members)

	log.Printf("[DEBUG] Deleting server group %s", d.Id())
	err = servergroups.Delete(ecsClient, d.Id()).ExtractErr()
	UnlockAll(members)
	if err != nil {
		return diag.Errorf("error deleting server group: %s", err)
	}

	return nil
}

func addServerGroupMember(client *golangsdk.ServiceClient, groupID, serverID string) error {
	// the ECS instances do not support other operations when binding server groups.
	config.MutexKV.Lock(serverID)
	defer config.MutexKV.Unlock(serverID)

	addMemberOpts := servergroups.MemberOpts{
		InstanceID: serverID,
	}
	return servergroups.UpdateMember(client, addMemberOpts, "add_member", groupID).ExtractErr()
}

func removeServerGroupMember(client *golangsdk.ServiceClient, groupID, serverID string) error {
	server, err := cloudservers.Get(client, serverID).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] the compute %s is not exist, ignore to remove it from the group", serverID)
			return nil
		}
		log.Printf("[WARN] failed to retrieve compute %s: %s, try to remove it from the group", serverID, err)
	} else if server.Status == "DELETED" || server.Status == "SOFT_DELETED" {
		log.Printf("[WARN] the compute %s was removed, ignore to remove it from the group", serverID)
		return nil
	}

	// the ECS instances do not support other operations when binding server groups.
	config.MutexKV.Lock(serverID)
	defer config.MutexKV.Unlock(serverID)

	removeMemberOpts := servergroups.MemberOpts{
		InstanceID: serverID,
	}
	return servergroups.UpdateMember(client, removeMemberOpts, "remove_member", groupID).ExtractErr()
}
