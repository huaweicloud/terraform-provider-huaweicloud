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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

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
				ValidateFunc: validation.StringInSlice([]string{
					resourceSpecCodeStandardBase, resourceSpecCodeProBase,
				}, false),
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Billing mode.`,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid",
				}, false),
			},
			"period_unit": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The charging period unit.`,
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
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

func resourceDscInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDscInstance: create a DSC.
	var (
		createDscInstanceHttpUrl = "v1/{project_id}/period/order"
		createDscInstanceProduct = "dsc"
	)
	createDscInstanceClient, err := cfg.NewServiceClient(createDscInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DscInstance Client: %s", err)
	}

	createDscInstancePath := createDscInstanceClient.Endpoint + createDscInstanceHttpUrl
	createDscInstancePath = strings.ReplaceAll(createDscInstancePath, "{project_id}", createDscInstanceClient.ProjectID)

	createDscInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	bodyParams, err := buildCreateDscInstanceBodyParams(d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}
	createDscInstanceOpt.JSONBody = utils.RemoveNil(bodyParams)
	createDscInstanceResp, err := createDscInstanceClient.Request("POST", createDscInstancePath, &createDscInstanceOpt)
	if err != nil {
		return diag.Errorf("error creating DscInstance: %s", err)
	}

	createDscInstanceRespBody, err := utils.FlattenResponse(createDscInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderId, err := jmespath.Search("order_id", createDscInstanceRespBody)
	if err != nil {
		return diag.Errorf("error creating DscInstance: ID is not found in API response")
	}

	// auto pay
	var (
		payOrderHttpUrl = "v3/orders/customer-orders/pay"
		payOrderProduct = "bss"
	)
	payOrderClient, err := cfg.NewServiceClient(payOrderProduct, region)
	if err != nil {
		return diag.Errorf("error creating BSS Client: %s", err)
	}

	payOrderPath := payOrderClient.Endpoint + payOrderHttpUrl

	payOrderOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	payOrderOpt.JSONBody = utils.RemoveNil(buildPayOrderBodyParams(orderId.(string)))
	_, err = payOrderClient.Request("POST", payOrderPath, &payOrderOpt)
	if err != nil {
		return diag.Errorf("error pay order=%s: %s", d.Id(), err)
	}

	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}
	err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
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
		getOrderProductsHttpUrl = "v2/bills/ratings/period-resources/subscribe-rate"
		getOrderProductsProduct = "bss"
	)
	getOrderProductsClient, err := cfg.NewServiceClient(getOrderProductsProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating BSS Client: %s", err)
	}

	getOrderProductsPath := getOrderProductsClient.Endpoint + getOrderProductsHttpUrl

	getOrderProductsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOrderProductsOpt.JSONBody = utils.RemoveNil(buildGetProductIdBodyParams(getOrderProductsClient.ProjectID,
		productInfos))
	getOrderProductResp, err := getOrderProductsClient.Request("POST", getOrderProductsPath, &getOrderProductsOpt)

	if err != nil {
		return nil, fmt.Errorf("error getting DSC order product infos: %s", err)
	}

	getOrderProductRespBody, err := utils.FlattenResponse(getOrderProductResp)
	if err != nil {
		return nil, err
	}
	curJson := utils.PathSearch("official_website_rating_result.product_rating_results",
		getOrderProductRespBody, make([]interface{}, 0))
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDscInstance: Query the DSC instance
	var (
		getDscInstanceHttpUrl = "v1/{project_id}/period/product/specification"
		getDscInstanceProduct = "dsc"
	)
	getDscInstanceClient, err := cfg.NewServiceClient(getDscInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DscInstance Client: %s", err)
	}

	getDscInstancePath := getDscInstanceClient.Endpoint + getDscInstanceHttpUrl
	getDscInstancePath = strings.ReplaceAll(getDscInstancePath, "{project_id}", getDscInstanceClient.ProjectID)

	getDscInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDscInstanceResp, err := getDscInstanceClient.Request("GET", getDscInstancePath, &getDscInstanceOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DscInstance")
	}

	getDscInstanceRespBody, err := utils.FlattenResponse(getDscInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dscOrder, err := jmespath.Search("orderInfo[?productInfo.resourceType=='hws.resource.type.dsc.base']",
		getDscInstanceRespBody)
	if err != nil {
		return diag.Errorf("error getting the instance info: base info not found in API response")
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
	cfg := meta.(*config.Config)

	if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
		return diag.Errorf("Error unsubscribing DSC order = %s: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETE"},
		Refresh:      waitForBmsInstanceDelete(ctx, d, meta),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting DSC instance: %s", err)
	}

	return nil
}

func waitForBmsInstanceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	return func() (interface{}, string, error) {
		// getDscInstance: Query the DSC instance
		var (
			getDscInstanceHttpUrl = "v1/{project_id}/period/product/specification"
			getDscInstanceProduct = "dsc"
		)
		getDscInstanceClient, err := cfg.NewServiceClient(getDscInstanceProduct, region)
		if err != nil {
			return nil, "error", fmt.Errorf("error creating DscInstance Client: %s", err)
		}

		getDscInstancePath := getDscInstanceClient.Endpoint + getDscInstanceHttpUrl
		getDscInstancePath = strings.ReplaceAll(getDscInstancePath, "{project_id}", getDscInstanceClient.ProjectID)

		getDscInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		getDscInstanceResp, err := getDscInstanceClient.Request("GET", getDscInstancePath, &getDscInstanceOpt)

		if err != nil {
			return nil, "error", fmt.Errorf("error retrieving DscInstance: %s", err)
		}

		getDscInstanceRespBody, err := utils.FlattenResponse(getDscInstanceResp)
		if err != nil {
			return nil, "error", fmt.Errorf("error retrieving DscInstance: %s", err)
		}

		orderInfo := utils.PathSearch("orderInfo", getDscInstanceRespBody, []interface{}{})
		orders := orderInfo.([]interface{})
		if len(orders) == 0 {
			return orders, "COMPLETE", nil
		}
		return nil, "PENDING", nil
	}
}
