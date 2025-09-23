package dns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var privateZoneAssociateNonUpdatableParams = []string{
	"zone_id", "router_id", "router_region",
}

// @API DNS POST /v2/zones/{zone_id}/associaterouter
// @API DNS POST /v2/zones/{zone_id}/disassociaterouter
// @API DNS GET /v2/zones/{zone_id}
func ResourceDNSPrivateZoneAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSPrivateZoneAssociateCreate,
		UpdateContext: resourceDNSPrivateZoneAssociateUpdate,
		ReadContext:   resourceDNSPrivateZoneAssociateRead,
		DeleteContext: resourceDNSPrivateZoneAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDNSPrivateZoneAssociateImportStateFunc,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(privateZoneAssociateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the zone to which the record set belongs.`,
			},
			"router_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the associated VPC.`,
			},
			"router_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region of the VPC.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the associated VPC.`,
			},
		},
	}
}

func resourceDNSPrivateZoneAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	zoneId := d.Get("zone_id").(string)
	routerId := d.Get("router_id").(string)

	httpUrl := "v2/zones/{zone_id}/associaterouter"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{zone_id}", zoneId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrDeleteDNSPrivateZoneAssociateBodyParams(d, region)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error associating DNS private zone with VPC: %s", err)
	}

	id := fmt.Sprintf("%s/%s", zoneId, routerId)
	d.SetId(id)

	stateRouterConf := &resource.StateChangeConf{
		Target:     []string{"COMPLETED"},
		Pending:    []string{"PENDING"},
		Refresh:    waitForDNSZoneRouterStatus(client, zoneId, routerId),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateRouterConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for associating zone (%s) to router (%s) become ACTIVE: %s",
			d.Id(), routerId, err)
	}

	return resourceDNSPrivateZoneAssociateRead(ctx, d, meta)
}

func buildCreateOrDeleteDNSPrivateZoneAssociateBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"router_id": d.Get("router_id"),
	}

	if v, ok := d.GetOk("router_region"); ok {
		bodyParams["router_region"] = v
	} else {
		bodyParams["router_region"] = region
	}

	return map[string]interface{}{
		"router": bodyParams,
	}
}

func resourceDNSPrivateZoneAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDNSPrivateZoneAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	getRespBody, err := getDNSZone(client, d.Get("zone_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS zone")
	}

	searchPath := fmt.Sprintf("routers[?router_id=='%s']|[0]", d.Get("router_id").(string))
	router := utils.PathSearch(searchPath, getRespBody, nil)
	if router == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving DNS zone router")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("router_id", utils.PathSearch("router_id", router, nil)),
		d.Set("router_region", utils.PathSearch("router_region", router, nil)),
		d.Set("status", utils.PathSearch("status", router, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getDNSZone(client *golangsdk.ServiceClient, zoneId string) (interface{}, error) {
	httpUrl := "v2/zones/{zone_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{zone_id}", zoneId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func resourceDNSPrivateZoneAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	zoneId := d.Get("zone_id").(string)
	routerId := d.Get("router_id").(string)

	httpUrl := "v2/zones/{zone_id}/disassociaterouter"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{zone_id}", zoneId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrDeleteDNSPrivateZoneAssociateBodyParams(d, region)),
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "code", "DNS.0707"),
			"error disassociating DNS private zone with VPC")
	}

	stateRouterConf := &resource.StateChangeConf{
		Target:     []string{"DELETED"},
		Pending:    []string{"COMPLETED", "PENDING"},
		Refresh:    waitForDNSZoneRouterStatus(client, zoneId, routerId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateRouterConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for disassociating zone (%s) to router (%s): %s", d.Id(), routerId, err)
	}

	return nil
}

func waitForDNSZoneRouterStatus(client *golangsdk.ServiceClient, zoneId string, routerId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		zone, err := getDNSZone(client, zoneId)
		if err != nil {
			return nil, "ERROR", err
		}

		searchPath := fmt.Sprintf("routers[?router_id=='%s']|[0]", routerId)
		router := utils.PathSearch(searchPath, zone, nil)
		if router != nil {
			status := utils.PathSearch("status", router, "").(string)
			if status == "ACTIVE" {
				return zone, "COMPLETED", nil
			}

			return zone, "PENDING", nil
		}

		return zone, "DELETED", nil
	}
}

func resourceDNSPrivateZoneAssociateImportStateFunc(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<zone_id>/<router_id>', but got '%s'", d.Id())
	}
	mErr := multierror.Append(nil,
		d.Set("zone_id", parts[0]),
		d.Set("router_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
