package vpcep

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/vpcep/v1/services"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	tagVPCEP        string = "endpoint"
	tagVPCEPService string = "endpoint_service"
)

func ResourceVPCEndpointService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCEndpointServiceCreate,
		ReadContext:   resourceVPCEndpointServiceRead,
		UpdateContext: resourceVPCEndpointServiceUpdate,
		DeleteContext: resourceVPCEndpointServiceDelete,

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
			"server_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VM", "LB", "VIP"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port_mapping": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "TCP",
						},
						"service_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"terminal_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-zA-Z0-9_-]{0,16}$"),
					"The name must have a maximum of 16 characters, and only contains letters, digits, underscores (_), and hyphens (-)."),
			},
			"service_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "interface",
				ValidateFunc: validation.StringInSlice([]string{"interface"}, false),
			},
			"approval": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"packet_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": common.TagsSchema(),
		},
	}
}

func expandPortMappingOpts(d *schema.ResourceData) []services.PortOpts {
	portMapping := d.Get("port_mapping").([]interface{})

	portOpts := make([]services.PortOpts, len(portMapping))
	for i, raw := range portMapping {
		port := raw.(map[string]interface{})
		portOpts[i].Protocol = port["protocol"].(string)
		portOpts[i].ServerPort = port["service_port"].(int)
		portOpts[i].ClientPort = port["terminal_port"].(int)
	}
	return portOpts
}

func doPermissionAction(client *golangsdk.ServiceClient, serviceID, action string, raw []interface{}) error {
	if len(raw) == 0 {
		return nil
	}

	permissions := make([]string, len(raw))
	for i, v := range raw {
		permissions[i] = v.(string)
	}
	permOpts := services.PermActionOpts{
		Action:      action,
		Permissions: permissions,
	}

	log.Printf("[DEBUG] %s permissions %#v to VPC endpoint service %s", action, permissions, serviceID)
	result := services.PermAction(client, serviceID, permOpts)

	return result.Err
}

func resourceVPCEndpointServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcepClient, err := config.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	approval := d.Get("approval").(bool)
	createOpts := services.CreateOpts{
		VpcID:       d.Get("vpc_id").(string),
		PortID:      d.Get("port_id").(string),
		ServerType:  d.Get("server_type").(string),
		ServiceName: d.Get("name").(string),
		ServiceType: d.Get("service_type").(string),
		Approval:    &approval,
		Ports:       expandPortMappingOpts(d),
	}
	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		createOpts.Tags = taglist
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := services.Create(vpcepClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC endpoint service: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[INFO] Waiting for VPC endpoint service(%s) to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"creating"},
		Target:       []string{"available"},
		Refresh:      waitForResourceStatus(vpcepClient, n.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf("error waiting for VPC endpoint service(%s) to become available: %s", n.ID, stateErr)
	}

	// add permissions
	raw := d.Get("permissions").(*schema.Set).List()
	err = doPermissionAction(vpcepClient, d.Id(), "add", raw)
	if err != nil {
		return diag.Errorf("error adding permissions to VPC endpoint service %s: %s", d.Id(), err)
	}

	return resourceVPCEndpointServiceRead(ctx, d, meta)
}

func resourceVPCEndpointServiceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcepClient, err := config.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	n, err := services.Get(vpcepClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC endpoint")
	}

	log.Printf("[DEBUG] retrieving VPC endpoint service: %#v", n)
	d.Set("region", region)
	d.Set("status", n.Status)
	d.Set("service_name", n.ServiceName)
	nameList := strings.Split(n.ServiceName, ".")
	if len(nameList) > 2 {
		d.Set("name", nameList[1])
	}

	d.Set("vpc_id", n.VpcID)
	d.Set("port_id", n.PortID)
	d.Set("approval", n.Approval)
	d.Set("server_type", n.ServerType)
	d.Set("service_type", n.ServiceType)

	ports := make([]map[string]interface{}, len(n.Ports))
	for i, v := range n.Ports {
		ports[i] = map[string]interface{}{
			"protocol":      v.Protocol,
			"service_port":  v.ServerPort,
			"terminal_port": v.ClientPort,
		}
	}
	d.Set("port_mapping", ports)

	// fetch tags from Services.Service
	tagmap := make(map[string]string)
	for _, val := range n.Tags {
		tagmap[val.Key] = val.Value
	}
	d.Set("tags", tagmap)

	// fetch connections
	if conns, err := flattenVPCEndpointConnections(vpcepClient, d.Id()); err == nil {
		d.Set("connections", conns)
	}

	// fetch permissions
	if perms, err := flattenVPCEndpointPermissions(vpcepClient, d.Id()); err == nil {
		d.Set("permissions", perms)
	}
	return nil
}

func resourceVPCEndpointServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcepClient, err := config.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	if d.HasChanges("name", "approval", "port_id", "port_mapping") {
		updateOpts := services.UpdateOpts{
			ServiceName: d.Get("name").(string),
		}

		if d.HasChange("approval") {
			approval := d.Get("approval").(bool)
			updateOpts.Approval = &approval
		}
		if d.HasChange("port_id") {
			updateOpts.PortID = d.Get("port_id").(string)
		}
		if d.HasChange("port_mapping") {
			updateOpts.Ports = expandPortMappingOpts(d)
		}

		_, err = services.Update(vpcepClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating VPC endpoint service: %s", err)
		}
	}

	//update tags
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(vpcepClient, d, tagVPCEPService, d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of VPC endpoint service %s: %s", d.Id(), tagErr)
		}
	}

	// update
	if d.HasChange("permissions") {
		old, new := d.GetChange("permissions")
		oldPermSet := old.(*schema.Set)
		newPermSet := new.(*schema.Set)
		added := newPermSet.Difference(oldPermSet)
		removed := oldPermSet.Difference(newPermSet)

		err = doPermissionAction(vpcepClient, d.Id(), "add", added.List())
		if err != nil {
			return diag.Errorf("error adding permissions to VPC endpoint service %s: %s", d.Id(), err)
		}

		err = doPermissionAction(vpcepClient, d.Id(), "remove", removed.List())
		if err != nil {
			return diag.Errorf("error removing permissions to VPC endpoint service %s: %s", d.Id(), err)
		}
	}

	return resourceVPCEndpointServiceRead(ctx, d, meta)
}

func resourceVPCEndpointServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcepClient, err := config.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	err = services.Delete(vpcepClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting VPC endpoint service %s: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"available", "deleting"},
		Target:       []string{"deleted"},
		Refresh:      waitForResourceStatus(vpcepClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting VPC endpoint service %s: %s", d.Id(), err)
	}

	return nil
}

func flattenVPCEndpointConnections(client *golangsdk.ServiceClient, id string) ([]map[string]interface{}, error) {
	allConns, err := services.ListConnections(client, id, nil)
	if err != nil {
		log.Printf("[WARN] Error querying connections of VPC endpoint service: %s", err)
		return nil, err
	}

	log.Printf("[DEBUG] retrieving connections of VPC endpoint service: %#v", allConns)
	connections := make([]map[string]interface{}, len(allConns))
	for i, v := range allConns {
		connections[i] = map[string]interface{}{
			"endpoint_id": v.EndpointID,
			"packet_id":   v.MarkerID,
			"domain_id":   v.DomainID,
			"status":      v.Status,
		}
	}

	return connections, nil
}

func flattenVPCEndpointPermissions(client *golangsdk.ServiceClient, id string) ([]string, error) {
	allPerms, err := services.ListPermissions(client, id)
	if err != nil {
		log.Printf("[WARN] Error querying permissions of VPC endpoint service: %s", err)
		return nil, err
	}

	log.Printf("[DEBUG] retrieving permissions of VPC endpoint service: %#v", allPerms)
	perms := make([]string, len(allPerms))
	for i, v := range allPerms {
		perms[i] = v.Permission
	}

	return perms, nil
}

func waitForResourceStatus(vpcepClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := services.Get(vpcepClient, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted VPC endpoint service %s", id)
				return n, "deleted", nil
			}
			return n, "error", err
		}

		return n, n.Status, nil
	}
}
