package vpcep

import (
	"context"
	"fmt"
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

// @API VPCEP DELETE /v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}
// @API VPCEP GET /v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}
// @API VPCEP PUT /v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}
// @API VPCEP POST /v1/{project_id}/vpc-endpoint-services
// @API VPCEP GET /v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/connections
// @API VPCEP POST /v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/permissions/action
// @API VPCEP GET /v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/permissions
// @API VPCEP POST /v1/{project_id}/{resource_type}/{resource_id}/tags/action
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

		CustomizeDiff: config.MergeDefaultTags(),

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
			"organization_permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
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
			"tcp_proxy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"snat_network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			// This field is not tested due to insufficient testing conditions.
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The IPv4 address or domain name of the server in the interface type VLAN scenario.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
			// This field is not tested due to insufficient testing conditions.
			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The dedicated cluster ID associated with the VPC endpoint service.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
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

func updatePermissions(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	if d.HasChange("permissions") {
		oldVal, newVal := d.GetChange("permissions")
		oldPermSet := oldVal.(*schema.Set)
		newPermSet := newVal.(*schema.Set)
		added := newPermSet.Difference(oldPermSet)
		removed := oldPermSet.Difference(newPermSet)

		if err := doPermissionAdd(client, added.List(), d.Id(), "domainId"); err != nil {
			return fmt.Errorf("error adding domain permissions to VPC endpoint service %s: %s", d.Id(), err)
		}

		if err := doPermissionRemove(client, removed.List(), d.Id(), "domainId"); err != nil {
			return fmt.Errorf("error removing domain permissions to VPC endpoint service %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("organization_permissions") {
		oldVal, newVal := d.GetChange("organization_permissions")
		oldPermSet := oldVal.(*schema.Set)
		newPermSet := newVal.(*schema.Set)
		added := newPermSet.Difference(oldPermSet)
		removed := oldPermSet.Difference(newPermSet)

		if err := doPermissionAdd(client, added.List(), d.Id(), "orgPath"); err != nil {
			return fmt.Errorf("error adding organization permissions to VPC endpoint service %s: %s", d.Id(), err)
		}

		if err := doPermissionRemove(client, removed.List(), d.Id(), "orgPath"); err != nil {
			return fmt.Errorf("error removing organization permissions to VPC endpoint service %s: %s", d.Id(), err)
		}
	}

	return nil
}

func doPermissionAdd(client *golangsdk.ServiceClient, raw []interface{}, serviceID, permissionType string) error {
	if len(raw) == 0 {
		return nil
	}
	permissions := utils.ExpandToStringList(raw)
	permOpts := services.PermActionOpts{
		Action:         "add",
		PermissionType: permissionType,
		Permissions:    permissions,
	}
	result := services.PermAction(client, serviceID, permOpts)
	return result.Err
}

func doPermissionRemove(client *golangsdk.ServiceClient, raw []interface{}, serviceID, permissionType string) error {
	if len(raw) == 0 {
		return nil
	}
	permissions := utils.ExpandToStringList(raw)
	permOpts := services.PermActionOpts{
		Action:         "remove",
		PermissionType: permissionType,
		Permissions:    permissions,
	}
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
		VpcID:         d.Get("vpc_id").(string),
		PortID:        d.Get("port_id").(string),
		ServerType:    d.Get("server_type").(string),
		ServiceName:   d.Get("name").(string),
		ServiceType:   d.Get("service_type").(string),
		Description:   d.Get("description").(string),
		TCPProxy:      d.Get("tcp_proxy").(string),
		IpVersion:     d.Get("ip_version").(string),
		SnatNetworkId: d.Get("snat_network_id").(string),
		IpAddress:     d.Get("ip_address").(string),
		PoolId:        d.Get("pool_id").(string),
		Approval:      utils.Bool(d.Get("approval").(bool)),
		Ports:         buildPortMappingOpts(d),
		Tags:          utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}

	// The European station does not support this parameter, so set it separately.
	if d.Get("enable_policy").(bool) {
		createOpts.EnablePolicy = utils.Bool(true)
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

	if err := updatePermissions(vpcepClient, d); err != nil {
		return diag.FromErr(err)
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
		d.Set("tcp_proxy", n.TCPProxy),
		d.Set("port_mapping", flattenVPCEndpointServicePorts(n)),
		d.Set("tags", utils.TagsToMap(n.Tags)),
		d.Set("enable_policy", n.EnablePolicy),
		d.Set("ip_version", n.IpVersion),
		d.Set("snat_network_id", n.SnatNetworkId),
		d.Set("ip_address", n.IpAddress),
		d.Set("pool_id", n.PoolId),
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
	perms, orgPerms, err := flattenVPCEndpointPermissions(vpcepClient, d.Id())
	if err == nil {
		mErr = multierror.Append(mErr, d.Set("permissions", perms))
		mErr = multierror.Append(mErr, d.Set("organization_permissions", orgPerms))
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

	if d.HasChanges("name", "approval", "port_id", "port_mapping", "description", "tcp_proxy", "ip_address") {
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
		if d.HasChange("tcp_proxy") {
			updateOpts.TCPProxy = d.Get("tcp_proxy").(string)
		}
		if d.HasChange("ip_address") {
			updateOpts.IpAddress = d.Get("ip_address").(string)
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
	if err := updatePermissions(vpcepClient, d); err != nil {
		return diag.FromErr(err)
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
	allConnections, err := services.ListAllConnections(client, id, nil)
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

func flattenVPCEndpointPermissions(client *golangsdk.ServiceClient, id string) (perms []string, orgPerms []string, err error) {
	allPerms, err := services.ListAllPermissions(client, id, nil)
	if err != nil {
		log.Printf("[WARN] Error querying permissions of VPC endpoint service: %s", err)
		return
	}

	log.Printf("[DEBUG] retrieving permissions of VPC endpoint service: %#v", allPerms)
	for _, v := range allPerms {
		if v.PermissionType == "domainId" {
			perms = append(perms, v.Permission)
		}
		if v.PermissionType == "orgPath" {
			orgPerms = append(orgPerms, v.Permission)
		}
	}

	return
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
