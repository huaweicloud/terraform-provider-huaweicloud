package deprecated

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/routes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceVPCRouteV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcRouteV2Create,
		ReadContext:   resourceVpcRouteV2Read,
		DeleteContext: resourceVpcRouteV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		DeprecationMessage: "use huaweicloud_vpc_route resource instead",

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
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nexthop": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: utils.ValidateCIDR,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVpcRouteV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcRouteClient, err := config.NetworkingV2Client(config.GetRegion(d))

	if err != nil {
		return diag.Errorf("error creating vpc route client: %s", err)
	}

	createOpts := routes.CreateOpts{
		Type:        d.Get("type").(string),
		NextHop:     d.Get("nexthop").(string),
		Destination: d.Get("destination").(string),
		VPC_ID:      d.Get("vpc_id").(string),
	}

	n, err := routes.Create(vpcRouteClient, createOpts).Extract()

	if err != nil {
		return diag.Errorf("error creating VPC route: %s", err)
	}
	d.SetId(n.RouteID)

	log.Printf("[INFO] Vpc Route ID: %s", n.RouteID)

	d.SetId(n.RouteID)

	return resourceVpcRouteV2Read(ctx, d, meta)

}

func resourceVpcRouteV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcRouteClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Vpc route client: %s", err)
	}

	n, err := routes.Get(vpcRouteClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error retrieving Vpc route: %s", err)
	}

	d.Set("type", n.Type)
	d.Set("nexthop", n.NextHop)
	d.Set("destination", n.Destination)
	d.Set("vpc_id", n.VPC_ID)
	d.Set("region", config.GetRegion(d))

	return nil
}

func resourceVpcRouteV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := meta.(*config.Config)
	vpcRouteClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating vpc route: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcRouteDelete(vpcRouteClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting Vpc route: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcRouteDelete(vpcRouteClient *golangsdk.ServiceClient, routeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := routes.Get(vpcRouteClient, routeId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted vpc route %s", routeId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = routes.Delete(vpcRouteClient, routeId).ExtractErr()
		log.Printf("[DEBUG] Value if error: %v", err)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted vpc route %s", routeId)
				return r, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault409); ok {
				return r, "ACTIVE", nil
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}
