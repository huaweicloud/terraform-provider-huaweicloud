package cc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			// The prepaid parameters do not return a value.
			"prepaid_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     prepaidOptionsSchema(),
			},
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func prepaidOptionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"period_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"is_auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func waitforPrepaidOrderCompleted(ctx context.Context, cfg *config.Config, d *schema.ResourceData,
	orderId string, timeout time.Duration) error {
	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}

	if err = common.WaitOrderComplete(ctx, bssClient, orderId, timeout); err != nil {
		return err
	}

	_, err = common.WaitOrderAllResourceComplete(ctx, bssClient, orderId, timeout)
	return err
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

	// We need to wait for the prepaid order to be completed; otherwise, a "404" error may occur when querying.
	if orderId := utils.PathSearch("bandwidth_package.order_id", respBody, "").(string); orderId != "" {
		if err := waitforPrepaidOrderCompleted(ctx, cfg, d, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error waiting for CC prepaid bandwidth package creation to complete (%s): %v", id, err)
		}
	}

	return resourceBandwidthPackageRead(ctx, d, meta)
}

func buildCreateBandwidthPackageProjectIdParam(d *schema.ResourceData, cfg *config.Config) string {
	if v, ok := d.GetOk("project_id"); ok {
		return v.(string)
	}

	return cfg.GetProjectID(cfg.GetRegion(d))
}

func buildCreateBandwidthPackagePrepaidOptionsParam(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"period_type":   rawMap["period_type"],
		"period_num":    rawMap["period_num"],
		"is_auto_renew": rawMap["is_auto_renew"],
		"is_auto_pay":   true,
	}
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
			"prepaid_options":       buildCreateBandwidthPackagePrepaidOptionsParam(d.Get("prepaid_options").([]interface{})),
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
		d.Set("created_at", utils.PathSearch("bandwidth_package.created_at", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("bandwidth_package.updated_at", respBody, nil)),
		d.Set("order_id", utils.PathSearch("bandwidth_package.order_id", respBody, nil)),
		d.Set("product_id", utils.PathSearch("bandwidth_package.product_id", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
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
	return err
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
	return err
}

func updateBandwidthPackageBindingObject(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData) error {
	idOld, idNew := d.GetChange("resource_id")
	typeOld, typeNew := d.GetChange("resource_type")

	if err := disassociateBandwidthPackage(client, cfg.DomainID, d.Id(), idOld.(string), typeOld.(string)); err != nil {
		return fmt.Errorf("error disassociating CC bandwidth package in update operation: %s", err)
	}

	if err := associateBandwidthPackage(client, cfg.DomainID, d.Id(), idNew.(string), typeNew.(string)); err != nil {
		return fmt.Errorf("error associating CC bandwidth package in update operation: %s", err)
	}

	return nil
}

func updateBandwidthPackage(client *golangsdk.ServiceClient, domainId, id string,
	params map[string]interface{}) (interface{}, error) {
	requestPath := client.Endpoint + "v3/{domain_id}/ccaas/bandwidth-packages/{id}"
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", domainId)
	requestPath = strings.ReplaceAll(requestPath, "{id}", id)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(params),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func buildUpdateBandwidthPackageNameAndDescription(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"name":        d.Get("name"),
			"description": utils.ValueIgnoreEmpty(d.Get("description")),
		},
	}
	return bodyParams
}

func buildUpdateBandwidthPackagePostpaidBillingMode(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"billing_mode": d.Get("billing_mode"),
			"bandwidth":    d.Get("bandwidth"),
		},
	}
	return bodyParams
}

func buildUpdateBandwidthPackagePrepaidOptions(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"period_type":   rawMap["period_type"],
		"period_num":    rawMap["period_num"],
		"is_auto_renew": rawMap["is_auto_renew"],
		"is_auto_pay":   true,
	}
}

func buildUpdateBandwidthPackagePrepaidBillingMode(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"billing_mode":    d.Get("billing_mode"),
			"prepaid_options": buildUpdateBandwidthPackagePrepaidOptions(d.Get("prepaid_options").([]interface{})),
		},
	}
	return bodyParams
}

func updateBandwidthPackageBillingMode(ctx context.Context, client *golangsdk.ServiceClient, cfg *config.Config,
	d *schema.ResourceData) error {
	var (
		billingMode             = d.Get("billing_mode").(string)
		postBillingModeArray    = []string{"5", "6"}
		prepaidBillingModeArray = []string{"1", "2"}
	)

	if utils.StrSliceContains(postBillingModeArray, billingMode) {
		// When updating to postpaid billing mode, the request body must include both `billing_mode` and `bandwidth`.
		_, err := updateBandwidthPackage(client, cfg.DomainID, d.Id(), buildUpdateBandwidthPackagePostpaidBillingMode(d))
		if err != nil {
			return fmt.Errorf("error updating CC postpaid bandwidth package billing mode: %s", err)
		}
	}

	if utils.StrSliceContains(prepaidBillingModeArray, billingMode) {
		// When updating to prepaid billing mode, the request body must include both `billing_mode` and `prepaid_options`.
		// This operation will generate a new order, and we need to track its completion.
		respBody, err := updateBandwidthPackage(client, cfg.DomainID, d.Id(), buildUpdateBandwidthPackagePrepaidBillingMode(d))
		if err != nil {
			return fmt.Errorf("error updating CC prepaid bandwidth package billing mode: %s", err)
		}

		if orderId := utils.PathSearch("bandwidth_package.order_id", respBody, "").(string); orderId != "" {
			if err := waitforPrepaidOrderCompleted(ctx, cfg, d, orderId, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return fmt.Errorf(`error waiting for CC prepaid bandwidth package billing mode update to complete (%s):
				 %v`, d.Id(), err)
			}
		}
	}

	return nil
}

func buildUpdateBandwidthPackagePostpaidBandwidth(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"bandwidth": d.Get("bandwidth"),
		},
	}
	return bodyParams
}

func buildUpdateBandwidthPackagePrepaidBandwidth(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bandwidth_package": map[string]interface{}{
			"bandwidth": d.Get("bandwidth"),
			"prepaid_options": map[string]interface{}{
				"is_auto_pay": true,
			},
		},
	}
	return bodyParams
}

func updateBandwidthPackageBandwidth(ctx context.Context, client *golangsdk.ServiceClient, cfg *config.Config,
	d *schema.ResourceData) error {
	var (
		billingMode             = d.Get("billing_mode").(string)
		postBillingModeArray    = []string{"3", "4", "5", "6"}
		prepaidBillingModeArray = []string{"1", "2"}
	)

	if utils.StrSliceContains(postBillingModeArray, billingMode) {
		// When updating the postpaid bandwidth value, the request body only need `bandwidth` field.
		_, err := updateBandwidthPackage(client, cfg.DomainID, d.Id(), buildUpdateBandwidthPackagePostpaidBandwidth(d))
		if err != nil {
			return fmt.Errorf("error updating CC postpaid bandwidth package bandwidth value: %s", err)
		}
	}

	if utils.StrSliceContains(prepaidBillingModeArray, billingMode) {
		// When updating the prepaid bandwidth value, the request body must include both `bandwidth` and `is_auto_pay` fields.
		// This operation will generate a new order, and we need to track its completion.
		respBody, err := updateBandwidthPackage(client, cfg.DomainID, d.Id(), buildUpdateBandwidthPackagePrepaidBandwidth(d))
		if err != nil {
			return fmt.Errorf("error updating CC prepaid bandwidth package bandwidth value: %s", err)
		}

		if orderId := utils.PathSearch("bandwidth_package.order_id", respBody, "").(string); orderId != "" {
			if err := waitforPrepaidOrderCompleted(ctx, cfg, d, orderId, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return fmt.Errorf("error waiting for CC prepaid bandwidth package bandwidth value update to complete (%s): %v", d.Id(), err)
			}
		}
	}

	return nil
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
		if err := updateBandwidthPackageBindingObject(client, cfg, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Editing the `name` and `description` is unaffected by postpaid and prepaid status; each is requested separately.
	if d.HasChanges("name", "description") {
		_, err := updateBandwidthPackage(client, cfg.DomainID, bandWidthId, buildUpdateBandwidthPackageNameAndDescription(d))
		if err != nil {
			return diag.Errorf("error updating CC bandwidth package name and description: %s", err)
		}
	}

	// It is necessary to distinguish the type of `billind_mode` change and handle them differently.
	if d.HasChange("billing_mode") {
		if err := updateBandwidthPackageBillingMode(ctx, client, cfg, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("bandwidth") {
		if err := updateBandwidthPackageBandwidth(ctx, client, cfg, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		err = updateBandwidthPackageTags(client, d, cfg.DomainID)
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

func deletePostpaidBandwidthPackage(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v3/{domain_id}/ccaas/bandwidth-packages/{id}"
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", cfg.DomainID)
	requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", requestPath, &requestOpt)
	return err
}

func deleteBandwidthPackage(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData) error {
	var (
		billingMode             = d.Get("billing_mode").(string)
		prepaidBillingModeArray = []string{"1", "2"}
	)

	// Instances of the prepaid class can only be deleted by unsubscribing.
	if utils.StrSliceContains(prepaidBillingModeArray, billingMode) {
		err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()})
		if err != nil {
			return fmt.Errorf("error unsubscribing CC prepaid bandwidth package: %s", err)
		}

		return nil
	}

	if err := deletePostpaidBandwidthPackage(client, cfg, d); err != nil {
		return fmt.Errorf("error deleting CC postpaid bandwidth package: %s", err)
	}

	return nil
}

func waitingForBandwidthPackageDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, cfg *config.Config) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetBandwidthPackage(client, cfg.DomainID, d.Id())
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					return "success", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceBandwidthPackageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
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

	if err := deleteBandwidthPackage(client, cfg, d); err != nil {
		return diag.FromErr(err)
	}

	// We need to confirm that the resource was successfully deleted.
	if err := waitingForBandwidthPackageDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete), cfg); err != nil {
		return diag.Errorf("error waiting for CC bandwidth package (%s) deleted: %s", d.Id(), err)
	}

	return nil
}
