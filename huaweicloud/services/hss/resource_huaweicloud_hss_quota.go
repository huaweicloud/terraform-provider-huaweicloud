// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product HSS
// ---------------------------------------------------------------
// Due to bugs in HuaweiCloud SKD, automatic generation writing is used.

package hss

import (
	"context"
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

// @API HSS POST /v5/{project_id}/quotas/orders
// @API HSS GET /v5/{project_id}/billing/quotas-detail
// @API HSS POST /v5/{project_id}/{resource_type}/{resource_id}/tags/create
// @API HSS DELETE /v5/{project_id}/{resource_type}/{resource_id}/tags/{key}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{resource_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{resource_id}
func ResourceQuota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceQuotaCreate,
		ReadContext:   resourceQuotaRead,
		UpdateContext: resourceQuotaUpdate,
		DeleteContext: resourceQuotaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"auto_renew": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Attributes
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shared_quota": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_trial_quota": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enterprise_project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateQuotaQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func buildCreatePeriodUnitParam(d *schema.ResourceData) interface{} {
	if d.Get("period_unit").(string) == "month" {
		return 2
	}

	if d.Get("period_unit").(string) == "year" {
		return 3
	}

	return nil
}

func buildCreateIsAutoRenewParam(d *schema.ResourceData) interface{} {
	if d.Get("auto_renew").(string) == "true" {
		return true
	}

	if d.Get("auto_renew").(string) == "false" {
		return false
	}

	return nil
}

func buildCreateQuotaBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"resource_spec_code": d.Get("version"),
		"period_type":        buildCreatePeriodUnitParam(d),
		"period_num":         d.Get("period").(int),
		"is_auto_renew":      buildCreateIsAutoRenewParam(d),
		"is_auto_pay":        1,
		// Currently, only one creation is supported.
		"subscription_num": 1,
	}

	return bodyParams
}

func resourceQuotaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	createPath := client.Endpoint + "v5/{project_id}/quotas/orders"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildCreateQuotaQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateQuotaBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating HSS quota: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderId := utils.PathSearch("order_id", createRespBody, "").(string)
	if orderId == "" {
		return diag.Errorf("unable to find the order ID of the HSS quota from the API response")
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for HSS quota order (%s) complete: %s", orderId, err)
	}

	d.SetId(resourceId)

	if tagsRaw, ok := d.GetOk("tags"); ok {
		err = createQuotaTags(client, resourceId, tagsRaw.(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceQuotaRead(ctx, d, meta)
}

func buildGetQuotaQueryParams(epsId, quotaId string) string {
	return fmt.Sprintf("?enterprise_project_id=%v&resource_id=%v", epsId, quotaId)
}

func GetQuotaById(client *golangsdk.ServiceClient, epsId, quotaId string) (interface{}, error) {
	getPath := client.Endpoint + "v5/{project_id}/billing/quotas-detail"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGetQuotaQueryParams(epsId, quotaId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HSS quota, %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	quotaResp := utils.PathSearch("data_list[0]", getRespBody, nil)

	return quotaResp, nil
}

func resourceQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		quotaId = d.Id()
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	quota, err := GetQuotaById(client, epsId, quotaId)
	if err != nil {
		return diag.FromErr(err)
	}

	if quota == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "HSS quota")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("version", utils.PathSearch("version", quota, nil)),
		d.Set("status", utils.PathSearch("quota_status", quota, nil)),
		d.Set("used_status", utils.PathSearch("used_status", quota, nil)),
		d.Set("host_id", utils.PathSearch("host_id", quota, nil)),
		d.Set("host_name", utils.PathSearch("host_name", quota, nil)),
		d.Set("charging_mode", flattenChargingMode(utils.PathSearch("charging_mode", quota, "").(string))),
		d.Set("expire_time", utils.PathSearch("expire_time", quota, nil)),
		d.Set("shared_quota", utils.PathSearch("shared_quota", quota, nil)),
		d.Set("is_trial_quota", utils.PathSearch("is_trial_quota", quota, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", quota, nil)),
		d.Set("enterprise_project_name", utils.PathSearch("enterprise_project_name", quota, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", quota, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceQuotaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		quotaId = d.Id()
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}

		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), quotaId); err != nil {
			return diag.Errorf("error updating the auto_renew of the HSS quota (%s): %s", quotaId, err)
		}
	}

	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})

		if len(oMap) > 0 {
			oldKeys := getOldTagKeys(oMap)
			if err := utils.DeleteResourceTagsWithKeys(client, oldKeys, "hss", quotaId); err != nil {
				return diag.FromErr(err)
			}
		}
		if len(nMap) > 0 {
			if err := createQuotaTags(client, quotaId, nMap); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   quotaId,
			ResourceType: "hss",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceQuotaRead(ctx, d, meta)
}

func createQuotaTags(client *golangsdk.ServiceClient, quotaId string, tagsMap map[string]interface{}) error {
	createTagsPath := client.Endpoint + "v5/{project_id}/{resource_type}/{resource_id}/tags/create"
	createTagsPath = strings.ReplaceAll(createTagsPath, "{project_id}", client.ProjectID)
	createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_type}", "hss")
	createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_id}", quotaId)
	createTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createTagsOpt.JSONBody = map[string]interface{}{
		"tags": utils.ExpandResourceTags(tagsMap),
	}

	_, err := client.Request("POST", createTagsPath, &createTagsOpt)
	if err != nil {
		return fmt.Errorf("error setting tags of the HSS quota (%s): %s", quotaId, err)
	}

	return nil
}

func getOldTagKeys(oRaw map[string]interface{}) []string {
	var tagKeys []string
	if len(oRaw) > 0 {
		for k := range oRaw {
			tagKeys = append(tagKeys, k)
		}
	}

	return tagKeys
}

func resourceQuotaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		quotaId = d.Id()
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if err = common.UnsubscribePrePaidResource(d, cfg, []string{quotaId}); err != nil {
		// When the resource does not exist, the API for unsubscribing prePaid resource will return a `400` status code,
		// and the response body is as follows:
		// {"error_code": "CBC.30000067",
		// "error_msg": "Unsubscription not supported. This resource has been deleted or the subscription to this resource has
		// not been synchronized to ..."}
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "CBC.30000067"), "error unsubscribe HSS quota")
	}

	if err := waitingForQuotaDeleted(ctx, client, epsId, quotaId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for HSS quota (%s) deleted: %s", quotaId, err)
	}

	return nil
}

func waitingForQuotaDeleted(ctx context.Context, client *golangsdk.ServiceClient, epsId, quotaId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			quota, err := GetQuotaById(client, epsId, quotaId)
			if err != nil {
				return nil, "ERROR", err
			}

			if quota == nil {
				m := map[string]string{"code": "COMPLETED"}
				return m, "COMPLETED", nil
			}

			return quota, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
