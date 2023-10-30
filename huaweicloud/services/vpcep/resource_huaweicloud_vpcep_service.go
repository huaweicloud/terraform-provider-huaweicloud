package vpcep

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/vpcep/v1/services"

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "schema: Required",
						},
						"terminal_port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "schema: Required",
						},
					},
				},
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "interface",
				Description: "schema: Computed",
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildPortMappingOpts(d *schema.ResourceData) []services.PortOpts {
	portMapping := d.Get("port_mapping").([]interface{})
	portOpts := make([]services.PortOpts, len(portMapping))
	for i, raw := range portMapping {
		port := raw.(map[string]interface{})
		portOpts[i] = services.PortOpts{
			Protocol:   port["protocol"].(string),
			ServerPort: port["service_port"].(int),
			ClientPort: port["terminal_port"].(int),
		}
	}
	return portOpts
}

func doPermissionAction(client *golangsdk.ServiceClient, serviceID, action string, raw []interface{}) error {
	if len(raw) == 0 {
		return nil
	}
	permissions := utils.ExpandToStringList(raw)
	permOpts := services.PermActionOpts{
		Action:      action,
		Permissions: permissions,
	}

	log.Printf("[DEBUG] %s permissions %#v to VPC endpoint service %s", action, permissions, serviceID)
	result := services.PermAction(client, serviceID, permOpts)
	return result.Err
}

func resourceVPCEndpointServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	createOpts := services.CreateOpts{
		VpcID:        d.Get("vpc_id").(string),
		PortID:       d.Get("port_id").(string),
		ServerType:   d.Get("server_type").(string),
		ServiceName:  d.Get("name").(string),
		ServiceType:  d.Get("service_type").(string),
		Description:  d.Get("description").(string),
		Approval:     utils.Bool(d.Get("approval").(bool)),
		Ports:        buildPortMappingOpts(d),
		Tags:         utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		EnablePolicy: utils.Bool(d.Get("enable_policy").(bool)),
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	n, err := services.Get(vpcepClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC endpoint service")
	}

	log.Printf("[DEBUG] retrieving VPC endpoint service: %#v", n)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("status", n.Status),
		d.Set("service_name", n.ServiceName),
		d.Set("vpc_id", n.VpcID),
		d.Set("port_id", n.PortID),
		d.Set("approval", n.Approval),
		d.Set("server_type", n.ServerType),
		d.Set("service_type", n.ServiceType),
		d.Set("description", n.Description),
		d.Set("port_mapping", flattenVPCEndpointServicePorts(n)),
		d.Set("tags", utils.TagsToMap(n.Tags)),
		d.Set("enable_policy", n.EnablePolicy),
	)

	nameList := strings.Split(n.ServiceName, ".")
	if len(nameList) > 2 {
		mErr = multierror.Append(mErr, d.Set("name", nameList[1]))
	}

	// fetch connections
	if connections, err := flattenVPCEndpointConnections(vpcepClient, d.Id()); err == nil {
		mErr = multierror.Append(mErr, d.Set("connections", connections))
	}

	// fetch permissions
	if perms, err := flattenVPCEndpointPermissions(vpcepClient, d.Id()); err == nil {
		mErr = multierror.Append(mErr, d.Set("permissions", perms))
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVPCEndpointServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	if d.HasChanges("name", "approval", "port_id", "port_mapping", "description") {
		updateOpts := services.UpdateOpts{
			ServiceName: d.Get("name").(string),
			Description: utils.String(d.Get("description").(string)),
		}

		if d.HasChange("approval") {
			updateOpts.Approval = utils.Bool(d.Get("approval").(bool))
		}
		if d.HasChange("port_id") {
			updateOpts.PortID = d.Get("port_id").(string)
		}
		if d.HasChange("port_mapping") {
			updateOpts.Ports = buildPortMappingOpts(d)
		}

		_, err = services.Update(vpcepClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating VPC endpoint service: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(vpcepClient, d, tagVPCEPService, d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of VPC endpoint service %s: %s", d.Id(), tagErr)
		}
	}

	// update permissions
	if d.HasChange("permissions") {
		oldVal, newVal := d.GetChange("permissions")
		oldPermSet := oldVal.(*schema.Set)
		newPermSet := newVal.(*schema.Set)
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
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

func flattenVPCEndpointServicePorts(n *services.Service) []map[string]interface{} {
	ports := make([]map[string]interface{}, len(n.Ports))
	for i, v := range n.Ports {
		ports[i] = map[string]interface{}{
			"protocol":      v.Protocol,
			"service_port":  v.ServerPort,
			"terminal_port": v.ClientPort,
		}
	}
	return ports
}

func flattenVPCEndpointConnections(client *golangsdk.ServiceClient, id string) ([]map[string]interface{}, error) {
	allConnections, err := services.ListConnections(client, id, nil)
	if err != nil {
		log.Printf("[WARN] Error querying connections of VPC endpoint service: %s", err)
		return nil, err
	}

	log.Printf("[DEBUG] retrieving connections of VPC endpoint service: %#v", allConnections)
	connections := make([]map[string]interface{}, len(allConnections))
	for i, v := range allConnections {
		connections[i] = map[string]interface{}{
			"endpoint_id": v.EndpointID,
			"packet_id":   v.MarkerID,
			"domain_id":   v.DomainID,
			"status":      v.Status,
			"description": v.Description,
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
