// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CC
// ---------------------------------------------------------------

package cc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CC POST /v3/{domain_id}/ccaas/bandwidth-packages
// @API CC DELETE /v3/{domain_id}/ccaas/bandwidth-packages/{id}
// @API CC GET /v3/{domain_id}/ccaas/bandwidth-packages/{id}
// @API CC PUT /v3/{domain_id}/ccaas/bandwidth-packages/{id}
// @API CC POST /v3/{domain_id}/ccaas/bandwidth-packages/{id}/associate
// @API CC POST /v3/{domain_id}/ccaas/bandwidth-packages/{id}/disassociate
// @API CC POST /v3/{domain_id}/ccaas/bandwidth-package/{id}/tags/action
func ResourceBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBandwidthPackageCreate,
		UpdateContext: resourceBandwidthPackageUpdate,
		ReadContext:   resourceBandwidthPackageRead,
		DeleteContext: resourceBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The bandwidth package name.`,
			},
			"local_area_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The local area ID.`,
			},
			"remote_area_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The remote area ID.`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Billing option of the bandwidth package.`,
			},
			"billing_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Billing mode of the bandwidth package.`,
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Bandwidth in the bandwidth package.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `Project ID.`,
			},
			"interflow_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Interflow mode of the bandwidth package.`,
			},
			"spec_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specification code of the bandwidth package.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description about the bandwidth package.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `ID of the enterprise project that the bandwidth package belongs to.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `ID of the resource that the bandwidth package is bound to.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Type of the resource that the bandwidth package is bound to.`,
			},
			"tags": common.TagsSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Bandwidth package status.`,
			},
		},
	}
}

func resourceBandwidthPackageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createBandwidthPackage: create a bandwidth package
	var (
		createBandwidthPackageHttpUrl = "v3/{domain_id}/ccaas/bandwidth-packages"
		createBandwidthPackageProduct = "cc"
	)
	createBandwidthPackageClient, err := cfg.NewServiceClient(createBandwidthPackageProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC Client: %s", err)
	}

	createBandwidthPackagePath := createBandwidthPackageClient.Endpoint + createBandwidthPackageHttpUrl
	createBandwidthPackagePath = strings.ReplaceAll(createBandwidthPackagePath, "{domain_id}", cfg.DomainID)

	createBandwidthPackageOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}

	createBandwidthPackageOpt.JSONBody = utils.RemoveNil(buildCreateBandwidthPackageBodyParams(d, cfg))
	createBandwidthPackageResp, err := createBandwidthPackageClient.Request("POST", createBandwidthPackagePath,
		&createBandwidthPackageOpt)
	if err != nil {
		return diag.Errorf("error creating bandwidth package: %s", err)
	}

	createBandwidthPackageRespBody, err := utils.FlattenResponse(createBandwidthPackageResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("bandwidth_package.id", createBandwidthPackageRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating bandwidth package: ID is not found in API response")
	}
	d.SetId(id)

	return resourceBandwidthPackageRead(ctx, d, meta)
}

func buildCreateBandwidthPackageBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	projectId := cfg.GetProjectID(cfg.GetRegion(d))
	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
	}
	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
			"description":           utils.ValueIgnoreEmpty(d.Get("description")),
			"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
			"local_area_id":         utils.ValueIgnoreEmpty(d.Get("local_area_id")),
			"remote_area_id":        utils.ValueIgnoreEmpty(d.Get("remote_area_id")),
			"charge_mode":           utils.ValueIgnoreEmpty(d.Get("charge_mode")),
			"billing_mode":          utils.ValueIgnoreEmpty(d.Get("billing_mode")),
			"bandwidth":             utils.ValueIgnoreEmpty(d.Get("bandwidth")),
			"project_id":            projectId,
			"resource_id":           utils.ValueIgnoreEmpty(d.Get("resource_id")),
			"resource_type":         utils.ValueIgnoreEmpty(d.Get("resource_type")),
			"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
			"interflow_mode":        utils.ValueIgnoreEmpty(d.Get("interflow_mode")),
			"spec_code":             utils.ValueIgnoreEmpty(d.Get("spec_code")),
		},
	}

	return bodyParams
}

func resourceBandwidthPackageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getBandwidthPackage: Query the bandwidth package
	var (
		getBandwidthPackageHttpUrl = "v3/{domain_id}/ccaas/bandwidth-packages/{id}"
		getBandwidthPackageProduct = "cc"
	)
	getBandwidthPackageClient, err := cfg.NewServiceClient(getBandwidthPackageProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC Client: %s", err)
	}

	getBandwidthPackagePath := getBandwidthPackageClient.Endpoint + getBandwidthPackageHttpUrl
	getBandwidthPackagePath = strings.ReplaceAll(getBandwidthPackagePath, "{domain_id}", cfg.DomainID)
	getBandwidthPackagePath = strings.ReplaceAll(getBandwidthPackagePath, "{id}", d.Id())

	getBandwidthPackageOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getBandwidthPackageResp, err := getBandwidthPackageClient.Request("GET", getBandwidthPackagePath, &getBandwidthPackageOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving bandwidth package")
	}

	getBandwidthPackageRespBody, err := utils.FlattenResponse(getBandwidthPackageResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("bandwidth_package.name", getBandwidthPackageRespBody, nil)),
		d.Set("description", utils.PathSearch("bandwidth_package.description", getBandwidthPackageRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("bandwidth_package.enterprise_project_id", getBandwidthPackageRespBody, nil)),
		d.Set("local_area_id", utils.PathSearch("bandwidth_package.local_area_id", getBandwidthPackageRespBody, nil)),
		d.Set("remote_area_id", utils.PathSearch("bandwidth_package.remote_area_id", getBandwidthPackageRespBody, nil)),
		d.Set("charge_mode", utils.PathSearch("bandwidth_package.charge_mode", getBandwidthPackageRespBody, nil)),
		d.Set("billing_mode", utils.PathSearch("bandwidth_package.billing_mode", getBandwidthPackageRespBody, nil)),
		d.Set("bandwidth", utils.PathSearch("bandwidth_package.bandwidth", getBandwidthPackageRespBody, nil)),
		d.Set("project_id", utils.PathSearch("bandwidth_package.project_id", getBandwidthPackageRespBody, nil)),
		d.Set("resource_id", utils.PathSearch("bandwidth_package.resource_id", getBandwidthPackageRespBody, nil)),
		d.Set("resource_type", utils.PathSearch("bandwidth_package.resource_type", getBandwidthPackageRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("bandwidth_package.tags", getBandwidthPackageRespBody, nil))),
		d.Set("interflow_mode", utils.PathSearch("bandwidth_package.interflow_mode", getBandwidthPackageRespBody, nil)),
		d.Set("spec_code", utils.PathSearch("bandwidth_package.spec_code", getBandwidthPackageRespBody, nil)),
		d.Set("status", utils.PathSearch("bandwidth_package.status", getBandwidthPackageRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBandwidthPackageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cc", region)
	if err != nil {
		return diag.Errorf("error creating CC Client: %s", err)
	}

	bandWidthId := d.Id()
	associateBandwidthPackageChanges := []string{
		"resource_id",
		"resource_type",
	}

	if d.HasChanges(associateBandwidthPackageChanges...) {
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

	updateBandwidthPackageChanges := []string{
		"name",
		"description",
		"bandwidth",
	}

	if d.HasChanges(updateBandwidthPackageChanges...) {
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

	updateBandwidthPackageChanges = []string{
		"billing_mode",
	}

	if d.HasChanges(updateBandwidthPackageChanges...) {
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
	if len(resourceId) == 0 {
		return nil
	}

	var associateBandwidthPackageHttpUrl = "v3/{domain_id}/ccaas/bandwidth-packages/{id}/associate"

	associateBandwidthPackagePath := client.Endpoint + associateBandwidthPackageHttpUrl
	associateBandwidthPackagePath = strings.ReplaceAll(associateBandwidthPackagePath, "{domain_id}", domainId)
	associateBandwidthPackagePath = strings.ReplaceAll(associateBandwidthPackagePath, "{id}", id)

	associateBandwidthPackageOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"resource_id":   utils.ValueIgnoreEmpty(resourceId),
			"resource_type": utils.ValueIgnoreEmpty(resourceType),
		},
	}

	associateBandwidthPackageOpt.JSONBody = utils.RemoveNil(bodyParams)
	_, err := client.Request("POST", associateBandwidthPackagePath, &associateBandwidthPackageOpt)
	if err != nil {
		return fmt.Errorf("error associating bandwidth package: %s", err)
	}
	return nil
}

func disassociateBandwidthPackage(client *golangsdk.ServiceClient, domainId, id, resourceId, resourceType string) error {
	if len(resourceId) == 0 {
		return nil
	}

	var disassociateBandwidthPackageHttpUrl = "v3/{domain_id}/ccaas/bandwidth-packages/{id}/disassociate"

	disassociateBandwidthPackagePath := client.Endpoint + disassociateBandwidthPackageHttpUrl
	disassociateBandwidthPackagePath = strings.ReplaceAll(disassociateBandwidthPackagePath, "{domain_id}", domainId)
	disassociateBandwidthPackagePath = strings.ReplaceAll(disassociateBandwidthPackagePath, "{id}", id)

	disassociateBandwidthPackageOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"resource_id":   utils.ValueIgnoreEmpty(resourceId),
			"resource_type": utils.ValueIgnoreEmpty(resourceType),
		},
	}

	disassociateBandwidthPackageOpt.JSONBody = utils.RemoveNil(bodyParams)
	_, err := client.Request("POST", disassociateBandwidthPackagePath, &disassociateBandwidthPackageOpt)
	if err != nil {
		return fmt.Errorf("error disassociating bandwidth package: %s", err)
	}
	return nil
}

func updateBandwidthPackage(client *golangsdk.ServiceClient, domainId, id string, params map[string]interface{}) error {
	var updateBandwidthPackageHttpUrl = "v3/{domain_id}/ccaas/bandwidth-packages/{id}"

	updateBandwidthPackagePath := client.Endpoint + updateBandwidthPackageHttpUrl
	updateBandwidthPackagePath = strings.ReplaceAll(updateBandwidthPackagePath, "{domain_id}", domainId)
	updateBandwidthPackagePath = strings.ReplaceAll(updateBandwidthPackagePath, "{id}", id)

	updateBandwidthPackageOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	updateBandwidthPackageOpt.JSONBody = utils.RemoveNil(params)
	_, err := client.Request("PUT", updateBandwidthPackagePath, &updateBandwidthPackageOpt)
	if err != nil {
		return fmt.Errorf("error updating bandwidth package: %s", err)
	}
	return nil
}

func updateBandwidthPackageTags(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string) error {
	var updateBandwidthPackageTagsHttpUrl = "v3/{domain_id}/ccaas/bandwidth-package/{id}/tags/action"

	updateBandwidthPackageTagsPath := client.Endpoint + updateBandwidthPackageTagsHttpUrl
	updateBandwidthPackageTagsPath = strings.ReplaceAll(updateBandwidthPackageTagsPath, "{domain_id}", domainId)
	updateBandwidthPackageTagsPath = strings.ReplaceAll(updateBandwidthPackageTagsPath, "{id}", d.Id())

	updateBandwidthPackageTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	// remove old tags
	if len(oMap) > 0 {
		updateBandwidthPackageTagsOpt.JSONBody = map[string]interface{}{
			"action": "delete",
			"tags":   utils.ExpandResourceTagsMap(oMap),
		}
		_, err := client.Request("POST", updateBandwidthPackageTagsPath, &updateBandwidthPackageTagsOpt)
		if err != nil {
			return fmt.Errorf("error updating bandwidth package: %s", err)
		}
	}

	// set new tags
	if len(nMap) > 0 {
		updateBandwidthPackageTagsOpt.JSONBody = map[string]interface{}{
			"action": "create",
			"tags":   utils.ExpandResourceTagsMap(nMap),
		}
		_, err := client.Request("POST", updateBandwidthPackageTagsPath, &updateBandwidthPackageTagsOpt)
		if err != nil {
			return fmt.Errorf("error updating bandwidth package: %s", err)
		}
	}
	return nil
}

func buildUpdateBandwidthPackageBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"name":        utils.ValueIgnoreEmpty(d.Get("name")),
			"description": utils.ValueIgnoreEmpty(d.Get("description")),
			"bandwidth":   utils.ValueIgnoreEmpty(d.Get("bandwidth")),
		},
	}
	return bodyParams
}

func buildUpdateBandwidthPackageBillingModeParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"billing_mode": utils.ValueIgnoreEmpty(d.Get("billing_mode")),
		},
	}
	return bodyParams
}

func resourceBandwidthPackageDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteBandwidthPackageHttpUrl = "v3/{domain_id}/ccaas/bandwidth-packages/{id}"
		deleteBandwidthPackageProduct = "cc"
	)
	deleteBandwidthPackageClient, err := cfg.NewServiceClient(deleteBandwidthPackageProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC Client: %s", err)
	}

	// disassociate bandwidth package first
	err = disassociateBandwidthPackage(deleteBandwidthPackageClient, cfg.DomainID, d.Id(),
		d.Get("resource_id").(string), d.Get("resource_type").(string))
	if err != nil {
		return diag.Errorf("error disassociating bandwidth package: %s", err)
	}

	deleteBandwidthPackagePath := deleteBandwidthPackageClient.Endpoint + deleteBandwidthPackageHttpUrl
	deleteBandwidthPackagePath = strings.ReplaceAll(deleteBandwidthPackagePath, "{domain_id}", cfg.DomainID)
	deleteBandwidthPackagePath = strings.ReplaceAll(deleteBandwidthPackagePath, "{id}", d.Id())

	deleteBandwidthPackageOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = deleteBandwidthPackageClient.Request("DELETE", deleteBandwidthPackagePath, &deleteBandwidthPackageOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting bandwidth package")
	}

	return nil
}
