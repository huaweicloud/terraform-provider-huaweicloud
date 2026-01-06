package ecs

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	bandwidthsv1 "github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/chnsz/golangsdk/openstack/networking/v2/bandwidths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const publicIPv6Type string = "5_dualStack"

// @API ECS GET /v1/{project_id}/cloudservers/{server_id}
// @API EIP PUT /v1/{project_id}/publicips/{publicip_id}
// @API EIP GET /v1/{project_id}/publicips
// @API EIP POST /v2.0/{project_id}/bandwidths/{bandwidth_id}/insert
// @API EIP POST /v2.0/{project_id}/bandwidths/{bandwidth_id}/remove
// @API EIP GET /v2.0/{project_id}/bandwidths/{bandwidth_id}
func ResourceComputeEIPAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeEIPAssociateCreate,
		ReadContext:   resourceComputeEIPAssociateRead,
		DeleteContext: resourceComputeEIPAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceComputeEIPAssociateImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
				ExactlyOneOf: []string{"bandwidth_id"},
			},
			"bandwidth_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"fixed_ip"},
			},
			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func eipAssociateRefreshFunc(client *golangsdk.ServiceClient, serverId, publicIP string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := cloudservers.Get(client, serverId).Extract()
		if err != nil {
			return nil, "ERROR", err
		}
		for _, addresses := range resp.Addresses {
			for _, address := range addresses {
				if address.Type == "floating" && address.Addr == publicIP {
					return resp, "COMPLETED", nil
				}
			}
		}
		return resp, "PENDING", nil
	}
}

func bandwidthAssociateRefreshFunc(client *golangsdk.ServiceClient, bwID, ipv6PortID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := bandwidthsv1.Get(client, bwID).Extract()
		if err != nil {
			return nil, "ERROR", err
		}

		for _, item := range resp.PublicipInfo {
			if item.PublicipId == ipv6PortID {
				return resp, "COMPLETED", nil
			}
		}
		return resp, "PENDING", nil
	}
}

func resourceComputeEIPAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	var associateID string
	var refreshFunc resource.StateRefreshFunc

	instanceID := d.Get("instance_id").(string)
	fixedIP := d.Get("fixed_ip").(string)

	if _, ok := d.GetOk("bandwidth_id"); ok {
		// fixed_ip must be a valid IPv6 address when combining with bandwidth_id
		if utils.IsIPv4Address(fixedIP) {
			return diag.Errorf("the fixed_ip must be a valid IPv6 address, got: %s", fixedIP)
		}
	}

	// get port id
	portID, privateIP, err := getComputeInstancePortIDbyFixedIP(ecsClient, instanceID, fixedIP)
	if err != nil {
		return diag.Errorf("error getting port ID of compute instance: %s", err)
	}

	if v, ok := d.GetOk("public_ip"); ok {
		vpcClient, err := cfg.NetworkingV1Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC client: %s", err)
		}

		// get EIP id
		eipAddr := v.(string)
		epsID := "all_granted_eps"
		eipID, err := common.GetEipIDbyAddress(vpcClient, eipAddr, epsID)
		if err != nil {
			return diag.Errorf("error getting EIP: %s", err)
		}

		err = bindPortToEIP(vpcClient, eipID, portID)
		if err != nil {
			return diag.Errorf("error associating port %s to EIP: %s", portID, err)
		}

		associateID = fmt.Sprintf("%s/%s/%s", eipAddr, instanceID, privateIP)
		refreshFunc = eipAssociateRefreshFunc(ecsClient, instanceID, eipAddr)
	} else {
		bwClient, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating bandwidth v2.0 client: %s", err)
		}

		bwID := d.Get("bandwidth_id").(string)
		err = insertPortToBandwidth(bwClient, bwID, portID)
		if err != nil {
			return diag.Errorf("error associating IPv6 port %s to bandwidth: %s", portID, err)
		}

		associateID = fmt.Sprintf("%s/%s/%s", bwID, instanceID, privateIP)
		refreshFunc = bandwidthAssociateRefreshFunc(bwClient, bwID, portID)
	}

	d.SetId(associateID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshFunc,
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for EIP association to complete: %s", err)
	}

	return resourceComputeEIPAssociateRead(ctx, d, meta)
}

func resourceComputeEIPAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	var associated bool
	var publicID string
	instanceID := d.Get("instance_id").(string)
	fixedIP := d.Get("fixed_ip").(string)

	// get port id of compute instance
	portID, privateIP, err := getComputeInstancePortIDbyFixedIP(ecsClient, instanceID, fixedIP)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "EIP associate")
	}

	if v, ok := d.GetOk("public_ip"); ok {
		eipAddr := v.(string)
		publicID = eipAddr
		epsID := "all_granted_eps"
		eipInfo, err := getFloatingIPbyAddress(vpcClient, eipAddr, epsID)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "EIP associate")
		}

		if eipInfo.PortID == portID {
			associated = true
		}
	} else {
		bwID := d.Get("bandwidth_id").(string)
		publicID = bwID
		band, err := bandwidthsv1.Get(vpcClient, bwID).Extract()
		if err != nil {
			return common.CheckDeletedDiag(d, err, "bandwidth")
		}

		for _, ipInfo := range band.PublicipInfo {
			if ipInfo.PublicipId == portID {
				associated = true
				break
			}
		}
	}

	if !associated {
		log.Printf("[WARN] the resource is not associated with the specified EIP or bandwidth")
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not associated",
				Detail:   "the resource is not associated with the specified EIP or bandwidth and will be removed in Terraform state.",
			},
		}
	}

	id := fmt.Sprintf("%s/%s/%s", publicID, instanceID, privateIP)
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("fixed_ip", privateIP),
		d.Set("port_id", portID),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceComputeEIPAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	fixedIP := d.Get("fixed_ip").(string)

	// get port id of compute instance
	portID, _, err := getComputeInstancePortIDbyFixedIP(ecsClient, instanceID, fixedIP)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "eip associate")
	}

	if v, ok := d.GetOk("public_ip"); ok {
		vpcClient, err := cfg.NetworkingV1Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC client: %s", err)
		}

		eipAddr := v.(string)
		epsID := "all_granted_eps"
		eipID, err := common.GetEipIDbyAddress(vpcClient, eipAddr, epsID)
		if err != nil {
			return diag.Errorf("error getting EIP: %s", err)
		}

		err = unbindPortFromEIP(vpcClient, eipID, portID)
		if err != nil {
			return diag.Errorf("error disassociating Floating IP: %s", err)
		}
	} else {
		bwClient, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating bandwidth v2.0 client: %s", err)
		}

		bwID := d.Get("bandwidth_id").(string)
		err = removePortFromBandwidth(bwClient, bwID, portID)
		if err != nil {
			return diag.Errorf("error disassociating IPv6 port %s from bandwidth: %s", portID, err)
		}
	}

	return nil
}

func getComputeInstancePortIDbyFixedIP(client *golangsdk.ServiceClient, instanceId, fixedIP string) (portId, privateIp string, err error) {
	var instance *cloudservers.CloudServer
	instance, err = cloudservers.Get(client, instanceId).Extract()
	if err != nil {
		return
	}

	if instance.Status == "DELETED" || instance.Status == "SOFT_DELETED" {
		err = golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte("the ECS instance has been deleted"),
			},
		}
		return
	}

	for _, networkAddresses := range instance.Addresses {
		for _, address := range networkAddresses {
			if address.Type == "fixed" {
				if fixedIP == "" || address.Addr == fixedIP {
					portId = address.PortID
					privateIp = address.Addr
					return
				}
			}
		}
	}

	err = golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte("the port ID does not exist"),
		},
	}
	return
}

func getFloatingIPbyAddress(client *golangsdk.ServiceClient, floatingIP, epsID string) (*eips.PublicIp, error) {
	listOpts := &eips.ListOpts{
		PublicIp:            []string{floatingIP},
		EnterpriseProjectId: epsID,
	}

	pages, err := eips.List(client, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return nil, err
	}

	if len(allEips) != 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("can not find the EIP by %s", floatingIP)),
			},
		}
	}

	return &allEips[0], nil
}

func insertPortToBandwidth(client *golangsdk.ServiceClient, bwID, portID string) error {
	insertOpts := bandwidths.BandWidthInsertOpts{
		PublicipInfo: []bandwidths.PublicIpInfoID{
			{
				PublicIPID:   portID,
				PublicIPType: publicIPv6Type,
			},
		},
	}

	log.Printf("[DEBUG] Insert port %s to bandwidth %s", portID, bwID)
	_, err := bandwidths.Insert(client, bwID, insertOpts).Extract()
	if err != nil {
		return fmt.Errorf("error inserting %s into bandwidth %s: %s", portID, bwID, err)
	}
	return nil
}

func removePortFromBandwidth(client *golangsdk.ServiceClient, bwID, portID string) error {
	removalChargeMode := "bandwidth"
	removalSize := 5
	removeOpts := bandwidths.BandWidthRemoveOpts{
		ChargeMode: removalChargeMode,
		Size:       &removalSize,
		PublicipInfo: []bandwidths.PublicIpInfoID{
			{
				PublicIPID:   portID,
				PublicIPType: publicIPv6Type,
			},
		},
	}

	log.Printf("[DEBUG] Remove port %s from bandwidth %s", portID, bwID)
	err := bandwidths.Remove(client, bwID, removeOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("error removing %s from bandwidth: %s", portID, err)
	}
	return nil
}

func bindPortToEIP(client *golangsdk.ServiceClient, eipID, portID string) error {
	log.Printf("[DEBUG] Bind port %s to EIP %s", portID, eipID)
	return actionOnPort(client, eipID, portID)
}

func unbindPortFromEIP(client *golangsdk.ServiceClient, eipID, portID string) error {
	log.Printf("[DEBUG] Unbind port %s from EIP: %s", portID, eipID)
	return actionOnPort(client, eipID, "")
}

func actionOnPort(client *golangsdk.ServiceClient, eipID, portID string) error {
	updateOpts := eips.UpdateOpts{
		PortID: portID,
	}
	_, err := eips.Update(client, eipID, updateOpts).Extract()
	if err != nil {
		return err
	}

	return nil
}

func parseComputeFloatingIPAssociateID(id string) (publicID, instanceID, fixedIP string, err error) {
	idParts := strings.Split(id, "/")
	// fixed_ip is optional, but we want users to give this value in import scenario
	if len(idParts) != 3 && len(idParts) != 2 {
		err = fmt.Errorf("unable to parse the resource ID, must be <eip address or bandwidth_id>/<instance_id>/<fixed_ip> format")
		return
	}

	publicID = idParts[0]
	instanceID = idParts[1]
	if len(idParts) == 3 {
		fixedIP = idParts[2]
	}

	return
}

func resourceComputeEIPAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	publicID, instanceID, fixedIP, err := parseComputeFloatingIPAssociateID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("instance_id", instanceID)
	d.Set("fixed_ip", fixedIP)
	parsedIP := net.ParseIP(publicID)
	if parsedIP != nil {
		d.Set("public_ip", publicID)
	} else {
		d.Set("bandwidth_id", publicID)
	}

	return []*schema.ResourceData{d}, nil
}
