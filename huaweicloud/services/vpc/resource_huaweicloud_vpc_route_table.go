package vpc

import (
	"context"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/routetables"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceVPCRouteTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcRouteTableCreate,
		ReadContext:   resourceVpcRouteTableRead,
		UpdateContext: resourceVpcRouteTableUpdate,
		DeleteContext: resourceVpcRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"subnets": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"route": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 200,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: utils.ValidateCIDR,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"nexthop": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// MaxCreateRoutes is the limitation of creating API
const MaxCreateRoutes int = 5

func resourceVpcRouteTableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating VPC client: %s", err)
	}

	createOpts := routetables.CreateOpts{
		VpcID:       d.Get("vpc_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	allRouteOpts := buildVpcRTRoutes(d)
	if len(allRouteOpts) <= MaxCreateRoutes {
		createOpts.Routes = allRouteOpts
	}

	logp.Printf("[DEBUG] VPC route table create options: %#v", createOpts)
	routeTable, err := routetables.Create(vpcClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating VPC route table: %s", err)
	}

	d.SetId(routeTable.ID)

	if v, ok := d.GetOk("subnets"); ok {
		subnets := utils.ExpandToStringList(v.(*schema.Set).List())
		err = associateRouteTableSubnets(vpcClient, d.Id(), subnets)
		if err != nil {
			return fmtp.DiagErrorf("Error associating subnets with VPC route table %s: %s", d.Id(), err)
		}
	}

	if len(allRouteOpts) > MaxCreateRoutes {
		updateOpts := routetables.UpdateOpts{
			Routes: map[string][]routetables.RouteOpts{
				"add": allRouteOpts,
			},
		}

		logp.Printf("[DEBUG] add routes to VPC route table %s: %#v", d.Id(), updateOpts)
		_, err = routetables.Update(vpcClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.DiagErrorf("Error creating VPC route: %s", err)
		}
	}

	return resourceVpcRouteTableRead(ctx, d, meta)

}

func resourceVpcRouteTableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating VPC client: %s", err)
	}

	routeTable, err := routetables.Get(vpcClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC route table")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("vpc_id", routeTable.VpcID),
		d.Set("name", routeTable.Name),
		d.Set("description", routeTable.Description),
		d.Set("route", expandVpcRTRoutes(routeTable.Routes)),
		d.Set("subnets", expandVpcRTSubnets(routeTable.Subnets)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error saving VPC route table: %s", err)
	}

	return nil
}

func resourceVpcRouteTableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating VPC client: %s", err)
	}

	var changed bool
	var updateOpts routetables.UpdateOpts
	if d.HasChanges("name", "description") {
		changed = true
		desc := d.Get("description").(string)
		updateOpts.Description = &desc
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("route") {
		changed = true
		routesOpts := map[string][]routetables.RouteOpts{}

		old, new := d.GetChange("route")
		addRaws := new.(*schema.Set).Difference(old.(*schema.Set))
		delRaws := old.(*schema.Set).Difference(new.(*schema.Set))

		if delLen := delRaws.Len(); delLen > 0 {
			delRouteOpts := make([]routetables.RouteOpts, delLen)
			for i, item := range delRaws.List() {
				opts := item.(map[string]interface{})
				delRouteOpts[i] = routetables.RouteOpts{
					Type:        opts["type"].(string),
					NextHop:     opts["nexthop"].(string),
					Destination: opts["destination"].(string),
				}
			}
			routesOpts["del"] = delRouteOpts
		}

		if addLen := addRaws.Len(); addLen > 0 {
			addRouteOpts := make([]routetables.RouteOpts, addLen)
			for i, item := range addRaws.List() {
				opts := item.(map[string]interface{})
				desc := opts["description"].(string)
				addRouteOpts[i] = routetables.RouteOpts{
					Type:        opts["type"].(string),
					NextHop:     opts["nexthop"].(string),
					Destination: opts["destination"].(string),
					Description: &desc,
				}
			}
			routesOpts["add"] = addRouteOpts
		}
		updateOpts.Routes = routesOpts
	}

	if changed {
		logp.Printf("[DEBUG] VPC route table update options: %#v", updateOpts)
		if _, err := routetables.Update(vpcClient, d.Id(), updateOpts).Extract(); err != nil {
			return fmtp.DiagErrorf("Error updating VPC route table: %s", err)
		}
	}

	if d.HasChange("subnets") {
		old, new := d.GetChange("subnets")
		associateRaws := new.(*schema.Set).Difference(old.(*schema.Set))
		disassociateRaws := old.(*schema.Set).Difference(new.(*schema.Set))

		disassociateSubnets := utils.ExpandToStringList(disassociateRaws.List())
		if len(disassociateSubnets) > 0 {
			err = disassociateRouteTableSubnets(vpcClient, d.Id(), disassociateSubnets)
			if err != nil {
				return fmtp.DiagErrorf("Error disassociating subnets with VPC route table %s: %s", d.Id(), err)
			}
		}

		associateSubnets := utils.ExpandToStringList(associateRaws.List())
		if len(associateSubnets) > 0 {
			err = associateRouteTableSubnets(vpcClient, d.Id(), associateSubnets)
			if err != nil {
				return fmtp.DiagErrorf("Error associating subnets with VPC route table %s: %s", d.Id(), err)
			}
		}
	}

	return resourceVpcRouteTableRead(ctx, d, meta)
}

func resourceVpcRouteTableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating VPC client: %s", err)
	}

	if v, ok := d.GetOk("subnets"); ok {
		subnets := utils.ExpandToStringList(v.(*schema.Set).List())
		err = disassociateRouteTableSubnets(vpcClient, d.Id(), subnets)
		if err != nil {
			return fmtp.DiagErrorf("Error disassociating subnets with VPC route table %s: %s", d.Id(), err)
		}
	}

	err = routetables.Delete(vpcClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting VPC route table: %s", err)
	}

	d.SetId("")
	return nil
}

func associateRouteTableSubnets(client *golangsdk.ServiceClient, id string, subnets []string) error {
	return ipmlVpcRTSubnetsAction(client, id, "associate", subnets)
}

func disassociateRouteTableSubnets(client *golangsdk.ServiceClient, id string, subnets []string) error {
	return ipmlVpcRTSubnetsAction(client, id, "disassociate", subnets)
}

func ipmlVpcRTSubnetsAction(client *golangsdk.ServiceClient, id, action string, subnets []string) error {
	var opts routetables.ActionSubnetsOpts
	switch action {
	case "associate":
		opts.Associate = subnets
	case "disassociate":
		opts.Disassociate = subnets
	default:
		return fmtp.Errorf("action should be associate or disassociate, but got %s", action)
	}

	actionOpts := routetables.ActionOpts{
		Subnets: opts,
	}

	logp.Printf("[DEBUG] %s subnets %v with VPC route table %s", action, subnets, id)
	_, err := routetables.Action(client, id, actionOpts).Extract()
	return err
}

func buildVpcRTRoutes(d *schema.ResourceData) []routetables.RouteOpts {
	rawRoutes := d.Get("route").(*schema.Set).List()
	routeOpts := make([]routetables.RouteOpts, len(rawRoutes))

	for i, raw := range rawRoutes {
		opts := raw.(map[string]interface{})
		routeDesc := opts["description"].(string)
		routeOpts[i] = routetables.RouteOpts{
			Type:        opts["type"].(string),
			NextHop:     opts["nexthop"].(string),
			Destination: opts["destination"].(string),
			Description: &routeDesc,
		}
	}

	return routeOpts
}

func expandVpcRTRoutes(routes []routetables.Route) []map[string]interface{} {
	rtRules := make([]map[string]interface{}, 0, len(routes))

	for _, item := range routes {
		// ignore local rule as it can not be modified
		if item.Type == "local" {
			continue
		}

		acessRule := map[string]interface{}{
			"destination": item.DestinationCIDR,
			"type":        item.Type,
			"nexthop":     item.NextHop,
			"description": item.Description,
		}
		rtRules = append(rtRules, acessRule)
	}

	return rtRules
}

func expandVpcRTSubnets(subnets []routetables.Subnet) []string {
	rtSubnets := make([]string, len(subnets))

	for i, item := range subnets {
		rtSubnets[i] = item.ID
	}

	return rtSubnets
}
