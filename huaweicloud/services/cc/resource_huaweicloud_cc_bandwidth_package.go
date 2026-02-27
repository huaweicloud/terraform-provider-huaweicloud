package cc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableBandwidthPackageParams = []string{
	"local_area_id",
	"remote_area_id",
	"charge_mode",
	"project_id",
	"interflow_mode",
	"spec_code",
}

// @API CC POST /v3/{domain_id}/ccaas/bandwidth-packages
// @API CC DELETE /v3/{domain_id}/ccaas/bandwidth-packages/{id}
// @API CC GET /v3/{domain_id}/ccaas/bandwidth-packages/{id}
// @API CC PUT /v3/{domain_id}/ccaas/bandwidth-packages/{id}
// @API CC POST /v3/{domain_id}/ccaas/bandwidth-packages/{id}/associate
// @API CC POST /v3/{domain_id}/ccaas/bandwidth-packages/{id}/disassociate
// @API CC POST /v3/{domain_id}/ccaas/bandwidth-packages/{id}/tag
// @API CC POST /v3/{domain_id}/ccaas/bandwidth-packages/{id}/untag
func ResourceBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBandwidthPackageCreate,
		UpdateContext: resourceBandwidthPackageUpdate,
		ReadContext:   resourceBandwidthPackageRead,
		DeleteContext: resourceBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(nonUpdatableBandwidthPackageParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"local_area_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_area_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"billing_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"interflow_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBandwidthPackageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{domain_id}/ccaas/bandwidth-packages"
		product = "cc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", cfg.DomainID)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateBandwidthPackageBodyParams(d, cfg)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating CC bandwidth package: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("bandwidth_package.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating bandwidth package: ID is not found in API response")
	}
	d.SetId(id)

	return resourceBandwidthPackageRead(ctx, d, meta)
}

func buildCreateBandwidthPackageProjectIdParam(d *schema.ResourceData, cfg *config.Config) string {
	if v, ok := d.GetOk("project_id"); ok {
		return v.(string)
	}

	return cfg.GetProjectID(cfg.GetRegion(d))
}

func buildCreateBandwidthPackageBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"name":                  d.Get("name"),
			"local_area_id":         d.Get("local_area_id"),
			"remote_area_id":        d.Get("remote_area_id"),
			"charge_mode":           d.Get("charge_mode"),
			"billing_mode":          d.Get("billing_mode"),
			"bandwidth":             d.Get("bandwidth"),
			"project_id":            buildCreateBandwidthPackageProjectIdParam(d, cfg),
			"interflow_mode":        utils.ValueIgnoreEmpty(d.Get("interflow_mode")),
			"spec_code":             utils.ValueIgnoreEmpty(d.Get("spec_code")),
			"description":           utils.ValueIgnoreEmpty(d.Get("description")),
			"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
			"resource_id":           utils.ValueIgnoreEmpty(d.Get("resource_id")),
			"resource_type":         utils.ValueIgnoreEmpty(d.Get("resource_type")),
			"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		},
	}
}

func GetBandwidthPackage(client *golangsdk.ServiceClient, domainId, packageId string) (interface{}, error) {
	requestPath := client.Endpoint + "v3/{domain_id}/ccaas/bandwidth-packages/{id}"
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", domainId)
	requestPath = strings.ReplaceAll(requestPath, "{id}", packageId)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceBandwidthPackageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	respBody, err := GetBandwidthPackage(client, cfg.DomainID, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CC bandwidth package")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("bandwidth_package.name", respBody, nil)),
		d.Set("local_area_id", utils.PathSearch("bandwidth_package.local_area_id", respBody, nil)),
		d.Set("remote_area_id", utils.PathSearch("bandwidth_package.remote_area_id", respBody, nil)),
		d.Set("charge_mode", utils.PathSearch("bandwidth_package.charge_mode", respBody, nil)),
		d.Set("billing_mode", utils.PathSearch("bandwidth_package.billing_mode", respBody, nil)),
		d.Set("bandwidth", utils.PathSearch("bandwidth_package.bandwidth", respBody, nil)),
		d.Set("project_id", utils.PathSearch("bandwidth_package.project_id", respBody, nil)),
		d.Set("interflow_mode", utils.PathSearch("bandwidth_package.interflow_mode", respBody, nil)),
		d.Set("spec_code", utils.PathSearch("bandwidth_package.spec_code", respBody, nil)),
		d.Set("description", utils.PathSearch("bandwidth_package.description", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("bandwidth_package.enterprise_project_id", respBody, nil)),
		d.Set("resource_id", utils.PathSearch("bandwidth_package.resource_id", respBody, nil)),
		d.Set("resource_type", utils.PathSearch("bandwidth_package.resource_type", respBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("bandwidth_package.tags", respBody, nil))),
		d.Set("status", utils.PathSearch("bandwidth_package.status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBandwidthPackageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		bandWidthId = d.Id()
	)

	client, err := cfg.NewServiceClient("cc", region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	if d.HasChanges("resource_id", "resource_type") {
		idOld, idNew := d.GetChange("resource_id")
		typeOld, typeNew := d.GetChange("resource_type")

		err = disassociateBandwidthPackage(client, cfg.DomainID, bandWidthId, idOld.(string), typeOld.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		err = associateBandwidthPackage(client, cfg.DomainID, bandWidthId, idNew.(string), typeNew.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("name", "description", "bandwidth") {
		err = updateBandwidthPackage(client, cfg.DomainID, bandWidthId, buildUpdateBandwidthPackageBodyParams(d))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		err = updateBandwidthPackageTags(client, d, cfg.DomainID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Editing the `billing_mode` field has special restrictions:
	// 1. When an edit API request is initiated, the value of `billing_mode` can only be one of [`1`, `2`, `5`, `6`, `7`, `8`].
	// 2. When an edit API request includes `billing_mode`, the API requires `bandwidth` to be passed in.
	// For the reasons mentioned above, the `billing_mode` need to be edited by calling the API separately.
	if d.HasChange("billing_mode") {
		err = updateBandwidthPackage(client, cfg.DomainID, bandWidthId, buildUpdateBandwidthPackageBillingModeParams(d))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   bandWidthId,
			ResourceType: "bwp",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceBandwidthPackageRead(ctx, d, meta)
}

func associateBandwidthPackage(client *golangsdk.ServiceClient, domainId, id, resourceId, resourceType string) error {
	if resourceId == "" || resourceType == "" {
		return nil
	}

	requestPath := client.Endpoint + "v3/{domain_id}/ccaas/bandwidth-packages/{id}/associate"
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", domainId)
	requestPath = strings.ReplaceAll(requestPath, "{id}", id)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"bandwidth_package": map[string]interface{}{
				"resource_id":   resourceId,
				"resource_type": resourceType,
			},
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error associating CC bandwidth package: %s", err)
	}
	return nil
}

func disassociateBandwidthPackage(client *golangsdk.ServiceClient, domainId, id, resourceId, resourceType string) error {
	if resourceId == "" || resourceType == "" {
		return nil
	}

	requestPath := client.Endpoint + "v3/{domain_id}/ccaas/bandwidth-packages/{id}/disassociate"
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", domainId)
	requestPath = strings.ReplaceAll(requestPath, "{id}", id)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"bandwidth_package": map[string]interface{}{
				"resource_id":   resourceId,
				"resource_type": resourceType,
			},
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error disassociating CC bandwidth package: %s", err)
	}
	return nil
}

func updateBandwidthPackage(client *golangsdk.ServiceClient, domainId, id string, params map[string]interface{}) error {
	requestPath := client.Endpoint + "v3/{domain_id}/ccaas/bandwidth-packages/{id}"
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", domainId)
	requestPath = strings.ReplaceAll(requestPath, "{id}", id)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(params),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating CC bandwidth package: %s", err)
	}
	return nil
}

func updateBandwidthPackageTags(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	// remove old tags
	if len(oMap) > 0 {
		err := deleteBandwidthPackageTags(client, d, domainId, oMap)
		if err != nil {
			return err
		}
	}

	// set new tags
	if len(nMap) > 0 {
		err := addBandwidthPackageTags(client, d, domainId, nMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func addBandwidthPackageTags(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string,
	tags map[string]interface{}) error {
	addPath := client.Endpoint + "v3/{domain_id}/ccaas/bandwidth-packages/{id}/tag"
	addPath = strings.ReplaceAll(addPath, "{domain_id}", domainId)
	addPath = strings.ReplaceAll(addPath, "{id}", d.Id())
	addOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
			201,
			204,
		},
		JSONBody: map[string]interface{}{
			"tags": utils.ExpandResourceTagsMap(tags),
		},
	}

	_, err := client.Request("POST", addPath, &addOpt)
	if err != nil {
		return fmt.Errorf("error adding CC bandwidth package tags: %s", err)
	}

	return nil
}

func deleteBandwidthPackageTags(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string,
	tags map[string]interface{}) error {
	deletePath := client.Endpoint + "v3/{domain_id}/ccaas/bandwidth-packages/{id}/untag"
	deletePath = strings.ReplaceAll(deletePath, "{domain_id}", domainId)
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
			201,
			204,
		},
		JSONBody: map[string]interface{}{
			"tags": utils.ExpandResourceTagsMap(tags),
		},
	}

	_, err := client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return fmt.Errorf("error deleting CC bandwidth package tags: %s", err)
	}

	return nil
}

func buildUpdateBandwidthPackageBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"name":        d.Get("name"),
			"bandwidth":   d.Get("bandwidth"),
			"description": utils.ValueIgnoreEmpty(d.Get("description")),
		},
	}
	return bodyParams
}

func buildUpdateBandwidthPackageBillingModeParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"billing_mode": d.Get("billing_mode"),
			"bandwidth":    d.Get("bandwidth"),
		},
	}
	return bodyParams
}

func resourceBandwidthPackageDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v3/{domain_id}/ccaas/bandwidth-packages/{id}"
		product      = "cc"
		resourceId   = d.Get("resource_id").(string)
		resourceType = d.Get("resource_type").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	// disassociate bandwidth package first
	err = disassociateBandwidthPackage(client, cfg.DomainID, d.Id(), resourceId, resourceType)
	if err != nil {
		return diag.Errorf("error disassociating CC bandwidth package: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", cfg.DomainID)
	requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CC bandwidth package")
	}

	return nil
}
