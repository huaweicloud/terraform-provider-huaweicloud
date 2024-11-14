// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DSC
// ---------------------------------------------------------------

package dsc

import (
	"context"
	"fmt"
	"strconv"
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

const (
	cloudServiceType = "hws.service.type.sdg"

	resourceTypeBase = "hws.resource.type.dsc.base"
	resourceTypeDB   = "hws.resource.type.dsc.db"
	resourceTypeObs  = "hws.resource.type.dsc.obs"

	resourceSpecCodeProBase      = "base_professional"
	resourceSpecCodeStandardBase = "base_standard"

	resourceSpecCodeProDB      = "DB_professional"
	resourceSpecCodeStandardDB = "DB_standard"

	resourceSpecCodeProObs      = "OBS_professional"
	resourceSpecCodeStandardObs = "OBS_standard"

	resourceSizeMeasureIdObs = 47
	resourceSizeMeasureIdDB  = 30
)

// @API DSC POST /v1/{project_id}/period/order
// @API DSC GET /v1/{project_id}/period/product/specification
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS POST /v3/orders/customer-orders/pay
// @API BSS POST /v2/bills/ratings/period-resources/subscribe-rate
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
func ResourceDscInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDscInstanceCreate,
		ReadContext:   resourceDscInstanceRead,
		DeleteContext: resourceDscInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"edition": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The edition of DSC.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Billing mode.`,
			},
			"period_unit": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The charging period unit.`,
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `The charging period.`,
			},
			"auto_renew": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Whether auto renew is enabled. Valid values are "true" and "false".`,
			},
			"obs_expansion_package": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Purchase OBS expansion packages. One OBS expansion package offers 1 TB of OBS storage.`,
			},
			"database_expansion_package": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Purchase database expansion packages. One expansion package offers one database.`,
			},
		},
	}
}

func payDscInstanceOrder(cfg *config.Config, d *schema.ResourceData, orderId string) error {
	var (
		region  = cfg.GetRegion(d)
		httpUrl = "v3/orders/customer-orders/pay"
		product = "bss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating BSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: utils.RemoveNil(buildPayOrderBodyParams(orderId)),
	}
	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error paying DSC instance order (%s): %s", orderId, err)
	}
	return nil
}

func resourceDscInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/period/order"
		product = "dsc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	bodyParams, err := buildCreateDscInstanceBodyParams(d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}
	requestOpt.JSONBody = utils.RemoveNil(bodyParams)
	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating DSC instance: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderId := utils.PathSearch("order_id", respBody, "").(string)
	if orderId == "" {
		return diag.Errorf("error creating DSC instance: ID is not found in API response")
	}

	if err := payDscInstanceOrder(cfg, d, orderId); err != nil {
		return diag.FromErr(err)
	}

	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}
	err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)
	return resourceDscInstanceRead(ctx, d, meta)
}

func buildCreateDscInstanceBodyParams(d *schema.ResourceData, cfg *config.Config) (map[string]interface{}, error) {
	productInfos, err := buildCreateDscInstanceRequestBodyProductInfos(d, cfg)
	if err != nil {
		return nil, err
	}
	bodyParams := map[string]interface{}{
		"regionId":         cfg.GetRegion(d),
		"cloudServiceType": cloudServiceType,
		"periodNum":        utils.ValueIgnoreEmpty(d.Get("period")),
		"productInfos":     productInfos,
	}

	chargingMode := d.Get("charging_mode").(string)
	if chargingMode == "prePaid" {
		bodyParams["chargingMode"] = 0
	}

	periodUnit := d.Get("period_unit").(string)
	if periodUnit == "month" {
		bodyParams["periodType"] = 2
	} else {
		bodyParams["periodType"] = 3
	}

	autoRenew := d.Get("auto_renew").(string)
	if autoRenew == "true" {
		bodyParams["isAutoRenew"] = 1
	} else {
		bodyParams["isAutoRenew"] = 0
	}

	return bodyParams, nil
}

func buildCreateDscInstanceRequestBodyProductInfos(d *schema.ResourceData, cfg *config.Config) ([]map[string]interface{}, error) {
	rst := make([]map[string]interface{}, 0)
	edition := d.Get("edition").(string)

	if edition == resourceSpecCodeStandardBase {
		rst = append(rst, map[string]interface{}{
			"cloudServiceType": cloudServiceType,
			"resourceType":     resourceTypeBase,
			"resourceSpecCode": resourceSpecCodeStandardBase,
		})

		if size, ok := d.GetOk("obs_expansion_package"); ok {
			rst = append(rst, map[string]interface{}{
				"cloudServiceType":      cloudServiceType,
				"resourceType":          resourceTypeObs,
				"resourceSpecCode":      resourceSpecCodeStandardObs,
				"resourceSize":          utils.ValueIgnoreEmpty(size),
				"resourceSizeMeasureId": resourceSizeMeasureIdObs,
			})
		}

		if size, ok := d.GetOk("database_expansion_package"); ok {
			rst = append(rst, map[string]interface{}{
				"cloudServiceType":      cloudServiceType,
				"resourceType":          resourceTypeDB,
				"resourceSpecCode":      resourceSpecCodeStandardDB,
				"resourceSize":          utils.ValueIgnoreEmpty(size),
				"resourceSizeMeasureId": resourceSizeMeasureIdDB,
			})
		}
	} else {
		rst = append(rst, map[string]interface{}{
			"cloudServiceType": cloudServiceType,
			"resourceType":     resourceTypeBase,
			"resourceSpecCode": resourceSpecCodeProBase,
		})

		if size, ok := d.GetOk("obs_expansion_package"); ok {
			rst = append(rst, map[string]interface{}{
				"cloudServiceType":      cloudServiceType,
				"resourceType":          resourceTypeObs,
				"resourceSpecCode":      resourceSpecCodeProObs,
				"resourceSize":          utils.ValueIgnoreEmpty(size),
				"resourceSizeMeasureId": resourceSizeMeasureIdObs,
			})
		}

		if size, ok := d.GetOk("database_expansion_package"); ok {
			rst = append(rst, map[string]interface{}{
				"cloudServiceType":      cloudServiceType,
				"resourceType":          resourceTypeDB,
				"resourceSpecCode":      resourceSpecCodeProDB,
				"resourceSize":          utils.ValueIgnoreEmpty(size),
				"resourceSizeMeasureId": resourceSizeMeasureIdDB,
			})
		}
	}
	if err := addProductIdToProductInfo(d, cfg, rst); err != nil {
		return nil, err
	}

	return rst, nil
}

func addProductIdToProductInfo(d *schema.ResourceData, cfg *config.Config, rst []map[string]interface{}) error {
	region := cfg.GetRegion(d)
	productInfos := make([]map[string]interface{}, len(rst))
	period := d.Get("period").(int)
	periodUnit := d.Get("period_unit").(string)
	for i, product := range rst {
		// The ID is used to identify the mapping between the inquiry result and the request
		// and must be unique in an inquiry.
		id := strconv.Itoa(i + 1)
		productInfos[i] = buildProductInfo(region, periodUnit, id, period, product)
	}
	products, err := getOrderProducts(cfg, region, productInfos)
	if err != nil {
		return fmt.Errorf("error getting DSC order product ID infos: %s", err)
	}
	productInfoMap := buildProductInfoMap(products)
	for i, v := range rst {
		if productId, ok := productInfoMap[strconv.Itoa(i+1)]; ok {
			v["productId"] = productId
			continue
		}
		return fmt.Errorf("error getting DSC order product ID by product(%#v): %s", v, err)
	}
	return nil
}

func buildProductInfoMap(products []interface{}) map[string]string {
	rst := make(map[string]string)
	for _, product := range products {
		id := utils.PathSearch("id", product, "").(string)
		productId := utils.PathSearch("product_id", product, "").(string)
		rst[id] = productId
	}
	return rst
}

func getOrderProducts(cfg *config.Config, region string, productInfos []map[string]interface{}) ([]interface{}, error) {
	var (
		httpUrl = "v2/bills/ratings/period-resources/subscribe-rate"
		product = "bss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating BSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildGetProductIdBodyParams(client.ProjectID, productInfos)),
	}
	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting DSC order product infos: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	expression := "official_website_rating_result.product_rating_results"
	curJson := utils.PathSearch(expression, respBody, make([]interface{}, 0))
	return curJson.([]interface{}), nil
}

func buildGetProductIdBodyParams(projectId string, productInfos []map[string]interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_id":    projectId,
		"product_infos": productInfos,
	}
	return bodyParams
}

func buildProductInfo(region, periodUnit, id string, periodNum int, resourceInfo map[string]interface{}) map[string]interface{} {
	var periodType string
	if periodUnit == "month" {
		periodType = "2"
	} else {
		periodType = "3"
	}

	params := make(map[string]interface{})
	params["id"] = id
	params["cloud_service_type"] = resourceInfo["cloudServiceType"]
	params["resource_type"] = resourceInfo["resourceType"]
	params["resource_spec"] = resourceInfo["resourceSpecCode"]
	params["region"] = region
	params["period_type"] = periodType
	params["period_num"] = periodNum
	params["subscription_num"] = "1"
	if resourceSize, ok := resourceInfo["resourceSize"]; ok {
		params["resource_size"] = resourceSize
	}
	if resourceSizeMeasureId, ok := resourceInfo["resourceSizeMeasureId"]; ok {
		params["size_measure_id"] = resourceSizeMeasureId
	}
	return params
}

func buildPayOrderBodyParams(orderId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"order_id":     orderId,
		"use_coupon":   "NO",
		"use_discount": "NO",
	}
	return bodyParams
}

func resourceDscInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/period/product/specification"
		product = "dsc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DSC instance")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	expression := "orderInfo[?productInfo.resourceType=='hws.resource.type.dsc.base']"
	dscOrder := utils.PathSearch(expression, respBody, nil)
	if dscOrder == nil {
		return diag.Errorf("unable to find the base information about the DSC instance from the API response")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("period_unit", parsePeriodUnit(utils.PathSearch("[0].periodType", dscOrder, nil))),
		d.Set("period", utils.PathSearch("[0].periodNum", dscOrder, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func parsePeriodUnit(periodType interface{}) string {
	pUnit := fmt.Sprintf("%v", periodType)
	if pUnit == "3" {
		return "year"
	}
	return "month"
}

func resourceDscInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
		return diag.Errorf("error unsubscribing DSC order (%s): %s", d.Id(), err)
	}

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETE"},
		Refresh:      waitingForDscInstanceDeleted(client),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DSC instance deletion to complete: %s", err)
	}

	return nil
}

func waitingForDscInstanceDeleted(client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	httpUrl := "v1/{project_id}/period/product/specification"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	return func() (interface{}, string, error) {
		resp, err := client.Request("GET", requestPath, &requestOpt)
		if err != nil {
			return nil, "ERROR", fmt.Errorf("error retrieving DSC instance: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, "ERROR", err
		}

		orders := utils.PathSearch("orderInfo", respBody, make([]interface{}, 0)).([]interface{})
		if len(orders) == 0 {
			return "success_deleted", "COMPLETE", nil
		}
		return nil, "PENDING", nil
	}
}
