package dc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var routeTableNonUpdatableParams = []string{"gdgw_id", "type", "destination", "nexthop"}

// @API DC PUT /v3/{project_id}/dcaas/gdgw/{gdgw_id}/routetables
// @API DC GET /v3/{project_id}/dcaas/gdgw/{gdgw_id}/routetables
func ResourceDcGlobalGatewayRouteTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcGlobalGatewayRouteTableCreate,
		ReadContext:   resourceDcGlobalGatewayRouteTableRead,
		UpdateContext: resourceDcGlobalGatewayRouteTableUpdate,
		DeleteContext: resourceDcGlobalGatewayRouteTableDelete,

		Importer: &schema.ResourceImporter{
			StateContext: globalGatewayRouteTableImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(routeTableNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"gdgw_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination": {
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
			"obtain_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address_family": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDcGlobalGatewayRouteTableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		product = "dc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	gatewayId := d.Get("gdgw_id").(string)

	createBodyParams := utils.RemoveNil(buildCreateGlobalGatewayRouteTableBodyParams(d))
	createResp, err := updateGlobalGatewayRouteTable(client, gatewayId, createBodyParams)
	if err != nil {
		return diag.Errorf("error creating DC global gateway route table: %s", err)
	}

	searchPath := fmt.Sprintf("gdgw_routetable[?type=='%s'&&nexthop=='%s'&&destination=='%s']|[0].id",
		d.Get("type").(string), d.Get("nexthop").(string), d.Get("destination").(string))
	id := utils.PathSearch(searchPath, createResp, "").(string)
	if id == "" {
		return diag.Errorf("error creating DC global gateway route table: ID is not found in the response")
	}

	d.SetId(id)

	err = waitForGlobalGatewayRouteTableActive(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDcGlobalGatewayRouteTableRead(ctx, d, meta)
}

func buildCreateGlobalGatewayRouteTableBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":        d.Get("type"),
		"destination": d.Get("destination"),
		"nexthop":     d.Get("nexthop"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return map[string]interface{}{
		"gdgw_routetable": map[string]interface{}{
			"add_routes": []interface{}{bodyParams},
		},
	}
}

func resourceDcGlobalGatewayRouteTableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	routeTable, err := getDcGlobalGatewayRouteTable(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DC global gateway route table")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("gdgw_id", utils.PathSearch("gateway_id", routeTable, nil)),
		d.Set("type", utils.PathSearch("type", routeTable, nil)),
		d.Set("destination", utils.PathSearch("destination", routeTable, nil)),
		d.Set("nexthop", utils.PathSearch("nexthop", routeTable, nil)),
		d.Set("description", utils.PathSearch("description", routeTable, nil)),
		d.Set("obtain_mode", utils.PathSearch("obtain_mode", routeTable, nil)),
		d.Set("status", utils.PathSearch("status", routeTable, nil)),
		d.Set("address_family", utils.PathSearch("address_family", routeTable, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDcGlobalGatewayRouteTableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	updateBodyParams := buildUpdateGlobalGatewayRouteTableBodyParam(d)
	_, err = updateGlobalGatewayRouteTable(client, d.Get("gdgw_id").(string), updateBodyParams)
	if err != nil {
		return diag.Errorf("error updating DC global gateway route table: %s", err)
	}

	err = waitForGlobalGatewayRouteTableActive(ctx, client, d, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDcGlobalGatewayRouteTableRead(ctx, d, meta)
}

func buildUpdateGlobalGatewayRouteTableBodyParam(d *schema.ResourceData) interface{} {
	bodyParams := map[string]interface{}{
		"destination": d.Get("destination"),
		"nexthop":     d.Get("nexthop"),
		"description": d.Get("description"),
	}

	return map[string]interface{}{
		"gdgw_routetable": map[string]interface{}{
			"update_routes": []interface{}{bodyParams},
		},
	}
}

func resourceDcGlobalGatewayRouteTableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	deleteBodyParams := utils.RemoveNil(buildDeleteGlobalGatewayRouteTableBodyParams(d))
	_, err = updateGlobalGatewayRouteTable(client, d.Get("gdgw_id").(string), deleteBodyParams)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", []string{"DC.0001", "DC.2720"}...),
			"error deleting DC global gateway route table")
	}

	return nil
}

func updateGlobalGatewayRouteTable(client *golangsdk.ServiceClient, gatewayId string, bodyParams interface{}) (interface{}, error) {
	httpUrl := "v3/{project_id}/dcaas/gdgw/{gdgw_id}/routetables"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{gdgw_id}", gatewayId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         bodyParams,
	}

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(updateResp)
}

func buildDeleteGlobalGatewayRouteTableBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":        d.Get("type"),
		"destination": d.Get("destination"),
		"nexthop":     d.Get("nexthop"),
	}

	return map[string]interface{}{
		"gdgw_routetable": map[string]interface{}{
			"del_routes": []interface{}{bodyParams},
		},
	}
}

func waitForGlobalGatewayRouteTableActive(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"ACTIVE"},
		Refresh:      globalGatewayRouteTableRefreshFunc(client, d),
		Timeout:      timeout,
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for DC global gateway route table to active: %s ", err)
	}
	return nil
}

func globalGatewayRouteTableRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		routeTable, err := getDcGlobalGatewayRouteTable(client, d)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("status", routeTable, "").(string)
		if status == "ACTIVE" {
			return routeTable, "ACTIVE", nil
		}
		if status == "ERROR" {
			return nil, "ERROR", errors.New("error updating DC global gateway route table")
		}

		return routeTable, "PENDING", nil
	}
}

func getDcGlobalGatewayRouteTable(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpPath := "v3/{project_id}/dcaas/gdgw/{gdgw_id}/routetables"
	getPath := client.Endpoint + httpPath
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{gdgw_id}", d.Get("gdgw_id").(string))
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

	routeTable := utils.PathSearch(fmt.Sprintf("gdgw_routetables[?id=='%s']|[0]", d.Id()), getRespBody, nil)
	if routeTable == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return routeTable, nil
}

func globalGatewayRouteTableImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <gdgw_id>/<id>")
	}
	gdgwId := parts[0]
	id := parts[1]
	d.SetId(id)
	err := d.Set("gdgw_id", gdgwId)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
