package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/routetables"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC DELETE /v1/{project_id}/routetables/{id}
// @API VPC GET /v1/{project_id}/routetables/{id}
// @API VPC PUT /v1/{project_id}/routetables/{id}
// @API VPC POST /v1/{project_id}/routetables
// @API VPC POST /v1/{project_id}/routetables/{id}/action
func ResourceVPCRouteTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcRouteTableCreate,
		ReadContext:   resourceVpcRouteTableRead,
		UpdateContext: resourceVpcRouteTableUpdate,
		DeleteContext: resourceVpcRouteTableDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
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
		return diag.Errorf("error creating VPC client: %s", err)
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

	log.Printf("[DEBUG] VPC route table create options: %#v", createOpts)
	routeTable, err := routetables.Create(vpcClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC route table: %s", err)
	}

	d.SetId(routeTable.ID)

	if v, ok := d.GetOk("subnets"); ok {
		subnets := utils.ExpandToStringList(v.(*schema.Set).List())
		err = associateRouteTableSubnets(vpcClient, d.Id(), subnets)
		if err != nil {
			return diag.Errorf("error associating subnets with VPC route table %s: %s", d.Id(), err)
		}
	}

	if len(allRouteOpts) > MaxCreateRoutes {
		updateOpts := routetables.UpdateOpts{
			Routes: map[string][]routetables.RouteOpts{
				"add": allRouteOpts,
			},
		}

		log.Printf("[DEBUG] add routes to VPC route table %s: %#v", d.Id(), updateOpts)
		_, err = routetables.Update(vpcClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error creating VPC route: %s", err)
		}
	}

	return resourceVpcRouteTableRead(ctx, d, meta)

}

func resourceVpcRouteTableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
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
		return diag.Errorf("error saving VPC route table: %s", err)
	}

	return nil
}

func resourceVpcRouteTableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
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

		var filteredMod []interface{}

		processed := map[string]bool{}

		for _, addItem := range addRaws.List() {
			addOpts := addItem.(map[string]interface{})
			addDest := addOpts["destination"].(string)

			for _, delItem := range delRaws.List() {
				delOpts := delItem.(map[string]interface{})
				delDest := delOpts["destination"].(string)

				if addDest == delDest && !processed[addDest] {
					filteredMod = append(filteredMod, addItem)
					processed[addDest] = true
					break
				}
			}
		}

		filteredAdd := []interface{}{}
		for _, item := range addRaws.List() {
			opts := item.(map[string]interface{})
			if !processed[opts["destination"].(string)] {
				filteredAdd = append(filteredAdd, item)
			}
		}

		filteredDel := []interface{}{}
		for _, item := range delRaws.List() {
			opts := item.(map[string]interface{})
			if !processed[opts["destination"].(string)] {
				filteredDel = append(filteredDel, item)
			}
		}

		if len(filteredDel) > 0 {
			delRouteOpts := make([]routetables.RouteOpts, len(filteredDel))
			for i, item := range filteredDel {
				opts := item.(map[string]interface{})
				delRouteOpts[i] = routetables.RouteOpts{
					Type:        opts["type"].(string),
					NextHop:     opts["nexthop"].(string),
					Destination: opts["destination"].(string),
				}
			}
			routesOpts["del"] = delRouteOpts
		}

		if len(filteredAdd) > 0 {
			addRouteOpts := make([]routetables.RouteOpts, len(filteredAdd))
			for i, item := range filteredAdd {
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

		if len(filteredMod) > 0 {
			modRouteOpts := make([]routetables.RouteOpts, len(filteredMod))
			for i, item := range filteredMod {
				opts := item.(map[string]interface{})
				desc := opts["description"].(string)
				modRouteOpts[i] = routetables.RouteOpts{
					Type:        opts["type"].(string),
					NextHop:     opts["nexthop"].(string),
					Destination: opts["destination"].(string),
					Description: &desc,
				}
			}
			routesOpts["mod"] = modRouteOpts
		}

		updateOpts.Routes = routesOpts
	}

	if changed {
		log.Printf("[DEBUG] VPC route table update options: %#v", updateOpts)
		if _, err := routetables.Update(vpcClient, d.Id(), updateOpts).Extract(); err != nil {
			return diag.Errorf("error updating VPC route table: %s", err)
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
				return diag.Errorf("error disassociating subnets with VPC route table %s: %s", d.Id(), err)
			}
		}

		associateSubnets := utils.ExpandToStringList(associateRaws.List())
		if len(associateSubnets) > 0 {
			err = associateRouteTableSubnets(vpcClient, d.Id(), associateSubnets)
			if err != nil {
				return diag.Errorf("error associating subnets with VPC route table %s: %s", d.Id(), err)
			}
		}
	}

	return resourceVpcRouteTableRead(ctx, d, meta)
}

func resourceVpcRouteTableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	if v, ok := d.GetOk("subnets"); ok {
		subnets := utils.ExpandToStringList(v.(*schema.Set).List())
		err = disassociateRouteTableSubnets(vpcClient, d.Id(), subnets)
		if err != nil {
			return diag.Errorf("error disassociating subnets with VPC route table %s: %s", d.Id(), err)
		}
	}

	err = routetables.Delete(vpcClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting VPC route table: %s", err)
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
		return fmt.Errorf("action should be associate or disassociate, but got %s", action)
	}

	actionOpts := routetables.ActionOpts{
		Subnets: opts,
	}

	log.Printf("[DEBUG] %s subnets %v with VPC route table %s", action, subnets, id)
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
