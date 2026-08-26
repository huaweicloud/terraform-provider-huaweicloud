package apig

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type featureConfig struct {
	UserRoutes []interface{} `json:"user_routes"`
}

var instanceRoutesNonUpdatableParams = []string{
	"instance_id",
}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/features
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/features
func ResourceInstanceRoutes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceRoutesCreate,
		ReadContext:   resourceInstanceRoutesRead,
		UpdateContext: resourceInstanceRoutesUpdate,
		DeleteContext: resourceInstanceRoutesDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceInstanceRoutesImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(instanceRoutesNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the dedicated instance and routes are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the dedicated instance to which the routes belong.",
			},
			"nexthops": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The configuration of the next hop routes.",
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true,
						Required: true,
					}),
			},
		},
	}
}

func modifyInstanceRoutes(client *golangsdk.ServiceClient, instanceId string, routes []interface{}) error {
	routeConfig := map[string]interface{}{
		"user_routes": routes,
	}
	routeBytes, err := json.Marshal(routeConfig)
	if err != nil {
		return fmt.Errorf("error parsing routes configuration: %s", err)
	}
	opts := instances.FeatureOpts{
		Name:   "route",
		Enable: utils.Bool(true),
		Config: string(routeBytes),
	}
	log.Printf("[DEBUG] The modify options of the instance routes is: %#v", opts)
	_, err = instances.UpdateFeature(client, instanceId, opts)
	if err != nil {
		return err
	}
	return nil
}

func resourceInstanceRoutesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG V2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		routes     = d.Get("nexthops").(*schema.Set)
	)
	if err := modifyInstanceRoutes(client, instanceId, routes.List()); err != nil {
		return diag.Errorf("error creating instance routes: %v", err)
	}
	d.SetId(instanceId)

	return resourceInstanceRoutesRead(ctx, d, meta)
}

func ListInstanceRoutes(client *golangsdk.ServiceClient, instanceId string) ([]interface{}, error) {
	opts := instances.ListFeaturesOpts{
		// Default value of parameter 'limit' is 20, parameter 'offset' is an invalid parameter.
		// If we omit it, we can only obtain 20 features, other features will be lost.
		Limit: 500,
	}
	resp, err := instances.ListFeatures(client, instanceId, opts)
	if err != nil {
		return nil, err
	}

	var routeConfig string
	for _, val := range resp {
		if val.Name == "route" {
			routeConfig = val.Config
			break
		}
	}

	result := utils.StringToJson(routeConfig)
	userRoutes := utils.PathSearch("user_routes", result, make([]interface{}, 0)).([]interface{})
	if len(userRoutes) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	// UserRoutes also means nexthops.
	return userRoutes, nil
}

func resourceInstanceRoutesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG V2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	nexthops, err := ListInstanceRoutes(client, instanceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting instance routes")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("nexthops", nexthops),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceInstanceRoutesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG V2 client: %s", err)
	}

	if d.HasChange("nexthops") {
		var (
			instanceId = d.Get("instance_id").(string)
			routes     = d.Get("nexthops").(*schema.Set)
		)
		if err := modifyInstanceRoutes(client, instanceId, routes.List()); err != nil {
			return diag.Errorf("error updating instance routes: %v", err)
		}
	}

	return resourceInstanceRoutesRead(ctx, d, meta)
}

func resourceInstanceRoutesDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG V2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	// The expression "{\"user_routes\":null}" has the same result as the expression"{\"user_routes\":[]}".
	if err := modifyInstanceRoutes(client, instanceId, nil); err != nil {
		return diag.Errorf("error deleting instance routes: %v", err)
	}

	return nil
}

func resourceInstanceRoutesImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	mErr := multierror.Append(nil, d.Set("instance_id", d.Id()))
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
