package live

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LIVE PUT /v1/{project_id}/domain/geo-blocking
// @API LIVE GET /v1/{project_id}/domain/geo-blocking
func ResourceGeoBlocking() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeoBlockingCreate,
		ReadContext:   resourceGeoBlockingRead,
		UpdateContext: resourceGeoBlockingUpdate,
		DeleteContext: resourceGeoBlockingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGeoBlockingImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the streaming domain name.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the application name.`,
			},
			// This field is specially designed to be mandatory compared to the parameters of openapi.
			// This design is intended to prevent the terraform lifecycle from being disrupted.
			"area_whitelist": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of supported areas. An empty list indicates no restriction.`,
			},
		},
	}
}

func buildUpdateGeoBlockingBodyParams(d *schema.ResourceData, areaWhitelist []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"app":            d.Get("app_name"),
		"area_whitelist": areaWhitelist,
	}
}

func updateGeoBlocking(client *golangsdk.ServiceClient, d *schema.ResourceData, areaWhitelist []interface{}) error {
	requestPath := client.Endpoint + "v1/{project_id}/domain/geo-blocking"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?play_domain=%s", d.Get("domain_name").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildUpdateGeoBlockingBodyParams(d, areaWhitelist),
	}
	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceGeoBlockingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		product       = "live"
		areaWhitelist = d.Get("area_whitelist").(*schema.Set).List()
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	resourceID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating Live geo blocking resource ID: %s", err)
	}

	if err := updateGeoBlocking(client, d, areaWhitelist); err != nil {
		return diag.Errorf("error creating Live geo blocking: %s", err)
	}

	d.SetId(resourceID)

	return resourceGeoBlockingRead(ctx, d, meta)
}

func ReadGeoBlocking(client *golangsdk.ServiceClient, domainDomain string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/domain/geo-blocking"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?play_domain=%s", domainDomain)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceGeoBlockingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "live"
		appName = d.Get("app_name").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	respBody, err := ReadGeoBlocking(client, d.Get("domain_name").(string))
	if err != nil {
		// When the domain does not exist, it will respond with a `400` status code.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "LIVE.100011001"),
			"error retrieving Live geo blocking")
	}

	expression := fmt.Sprintf("apps[?app == '%s']|[0].area_whitelist", appName)
	areaWhitelist := utils.PathSearch(expression, respBody, nil)
	if areaWhitelist == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("domain_name", utils.PathSearch("play_domain", respBody, nil)),
		d.Set("app_name", appName),
		d.Set("area_whitelist", areaWhitelist),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGeoBlockingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		product       = "live"
		areaWhitelist = d.Get("area_whitelist").(*schema.Set).List()
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	if err := updateGeoBlocking(client, d, areaWhitelist); err != nil {
		return diag.Errorf("error updating Live geo blocking: %s", err)
	}
	return resourceGeoBlockingRead(ctx, d, meta)
}

func resourceGeoBlockingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "live"
		// Passing an empty list to this field means no restriction.
		areaWhitelist = make([]interface{}, 0)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	if err := updateGeoBlocking(client, d, areaWhitelist); err != nil {
		// When the domain does not exist, it will respond with a `400` status code.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "LIVE.100011001"),
			"error deleting Live geo blocking")
	}
	return nil
}

func resourceGeoBlockingImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate ID: %s", err)
	}

	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<domain_name>/<app_name>', but got '%s'",
			importedId)
	}

	d.SetId(resourceId)
	mErr := multierror.Append(nil,
		d.Set("domain_name", parts[0]),
		d.Set("app_name", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
