package ccm

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

// @API CCM POST /v1/private-certificate-authorities
// @API CCM POST /v1/private-certificate-authorities/order
// @API CCM POST /v1/private-certificate-authorities/{ca_id}/activate
// @API CCM POST /v1/private-certificate-authorities/{ca_id}/tags/create
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API CCM GET /v1/private-certificate-authorities/{ca_id}
// @API CCM GET /v1/private-certificate-authorities/{ca_id}/tags
// @API CCM DELETE /v1/private-certificate-authorities/{ca_id}
// @API CCM POST /v1/private-certificate-authorities/{ca_id}/enable
// @API CCM POST /v1/private-certificate-authorities/{ca_id}/disable
// @API CCM DELETE /v1/private-certificate-authorities/{ca_id}/tags/delete
// @API CCM POST /v1/private-certificate-authorities/{ca_id}/crl/enable
// @API CCM POST /v1/private-certificate-authorities/{ca_id}/crl/disable
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourcePrivateCertificateAuthority() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateCACreate,
		ReadContext:   resourcePrivateCARead,
		UpdateContext: resourcePrivateCAUpdate,
		DeleteContext: resourcePrivateCADelete,
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
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"distinguished_name": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"common_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"country": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"state": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"locality": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"organizational_unit": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"key_algorithm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"validity": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"started_at": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"pending_days": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"issuer_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"path_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_usages": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"crl_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"obs_bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"valid_days": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"crl_dis_point": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"crl_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "schema: Computed",
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"auto_renew":    common.SchemaAutoRenew(nil),
			"tags":          common.TagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gen_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expired_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"free_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildPrepaidPrivateCABodyParams(d *schema.ResourceData, raw map[string]interface{},
	cfg *config.Config) map[string]interface{} {
	periodType := 2
	if raw["type"].(string) == "YEAR" {
		periodType = 3
	}

	autoRenew := 0
	if val, ok := d.GetOk("auto_renew"); ok && val == "true" {
		autoRenew = 1
	}

	var productInfos []map[string]interface{}
	productInfos = append(productInfos, map[string]interface{}{
		"cloud_service_type": "hws.service.type.ccm",
		"resource_type":      "hws.resource.type.pca.duration",
		"resource_spec_code": "ca.duration",
	})

	bodyParams := map[string]interface{}{
		"cloud_service_type":    "hws.service.type.ccm",
		"charging_mode":         0,
		"period_type":           periodType,
		"period_num":            raw["value"].(int),
		"is_auto_renew":         autoRenew,
		"is_auto_pay":           1,
		"subscription_num":      1,
		"enterprise_project_id": cfg.GetEnterpriseProjectID(d),
		"product_infos":         productInfos,
	}
	return bodyParams
}

func createPrepaidPrivateCA(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) (interface{}, error) {
	var (
		httpUrl  = "v1/private-certificate-authorities/order"
		rawArray = d.Get("validity").([]interface{})
		raw      = rawArray[0].(map[string]interface{})
		rawType  = raw["type"].(string)
	)

	if rawType != "YEAR" && rawType != "MONTH" {
		return nil, fmt.Errorf("the validity type value (%s) is invalid, only `YEAR` or `MONTH` is supported when"+
			" creating a prepaid private CA", rawType)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildPrepaidPrivateCABodyParams(d, raw, cfg),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, fmt.Errorf("error creating CCM prepaid private CA: %s", err)
	}
	return utils.FlattenResponse(createResp)
}

func waitingForOrderComplete(ctx context.Context, bssClient *golangsdk.ServiceClient, orderID string, timeout time.Duration) error {
	if err := common.WaitOrderComplete(ctx, bssClient, orderID, timeout); err != nil {
		return err
	}

	if _, err := common.WaitOrderResourceComplete(ctx, bssClient, orderID, timeout); err != nil {
		return fmt.Errorf("error waiting for order resource %s complete: %s", orderID, err)
	}
	return nil
}

func buildPrivateCARequestBodyDistinguishedName(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok && len(rawArray) > 0 {
		raw := rawArray[0].(map[string]interface{})
		return map[string]interface{}{
			"common_name":         raw["common_name"],
			"country":             raw["country"],
			"state":               raw["state"],
			"locality":            raw["locality"],
			"organization":        raw["organization"],
			"organizational_unit": raw["organizational_unit"],
		}
	}
	return nil
}

func buildPrivateCARequestBodyValidity(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok && len(rawArray) > 0 {
		raw := rawArray[0].(map[string]interface{})
		return map[string]interface{}{
			"type":       raw["type"],
			"value":      raw["value"],
			"start_from": utils.ValueIgnoreEmpty(raw["started_at"]),
		}
	}
	return nil
}

func buildPrivateCARequestBodyKeyUsages(caType interface{}, rawParams interface{}) []interface{} {
	rawArray, _ := rawParams.([]interface{})
	if caType.(string) == "ROOT" || len(rawArray) == 0 {
		return []interface{}{"digitalSignature", "keyCertSign", "cRLSign"}
	}
	return rawArray
}

func buildPrivateCARequestBodyCrlConfiguration(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}
	rawMap := rawArray[0].(map[string]interface{})
	if !rawMap["enabled"].(bool) {
		return nil
	}

	return map[string]interface{}{
		"enabled":         rawMap["enabled"],
		"crl_name":        rawMap["crl_name"],
		"obs_bucket_name": rawMap["obs_bucket_name"],
		"valid_days":      rawMap["valid_days"],
	}
}

func buildCreateOrActivatePrivateCABodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":                  d.Get("type"),
		"distinguished_name":    buildPrivateCARequestBodyDistinguishedName(d.Get("distinguished_name")),
		"key_algorithm":         d.Get("key_algorithm"),
		"signature_algorithm":   d.Get("signature_algorithm"),
		"validity":              buildPrivateCARequestBodyValidity(d.Get("validity")),
		"issuer_id":             utils.ValueIgnoreEmpty(d.Get("issuer_id")),
		"path_length":           utils.ValueIgnoreEmpty(d.Get("path_length")),
		"key_usages":            buildPrivateCARequestBodyKeyUsages(d.Get("type"), d.Get("key_usages")),
		"crl_configuration":     buildPrivateCARequestBodyCrlConfiguration(d.Get("crl_configuration")),
		"enterprise_project_id": cfg.GetEnterpriseProjectID(d),
	}
	return bodyParams
}

func activatePrivateCA(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	activatePath := client.Endpoint + "v1/private-certificate-authorities/{ca_id}/activate"
	activatePath = strings.ReplaceAll(activatePath, "{ca_id}", d.Id())
	activateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: utils.RemoveNil(buildCreateOrActivatePrivateCABodyParams(d, cfg)),
	}

	_, err := client.Request("POST", activatePath, &activateOpt)
	if err != nil {
		return fmt.Errorf("error activating CCM private CA: %s", err)
	}
	return nil
}

func createPostpaidPrivateCA(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) (interface{}, error) {
	createPath := client.Endpoint + "v1/private-certificate-authorities"
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrActivatePrivateCABodyParams(d, cfg)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, fmt.Errorf("error creating CCM postpaid private CA: %s", err)
	}
	return utils.FlattenResponse(createResp)
}

func resourcePrivateCACreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if _, ok := d.GetOk("issuer_id"); !ok && d.Get("type").(string) == "SUBORDINATE" {
		return diag.Errorf("the parameter `issuer_id` is required when creating a subordinate private CA")
	}

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("ccm", region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		createRespBody, err := createPrepaidPrivateCA(client, d, cfg)
		if err != nil {
			return diag.FromErr(err)
		}

		caId := utils.PathSearch("ca_ids|[0]", createRespBody, "").(string)
		if caId == "" {
			return diag.Errorf("unable to find the CCM prepaid private CA ID from the API response")
		}
		d.SetId(caId)

		orderId := utils.PathSearch("order_id", createRespBody, "").(string)
		if orderId == "" {
			return diag.Errorf("unable to find the order ID of the CCM prepaid private CA from the API response")
		}

		// wait for order success
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		if err := waitingForOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}

		// activate private CA
		if err := activatePrivateCA(client, d, cfg); err != nil {
			return diag.FromErr(err)
		}
	} else {
		createPrivateCARespBody, err := createPostpaidPrivateCA(client, d, cfg)
		if err != nil {
			return diag.FromErr(err)
		}

		caId := utils.PathSearch("ca_id", createPrivateCARespBody, "").(string)
		if caId == "" {
			return diag.Errorf("unable to find the CCM postpaid private CA ID from the API response")
		}
		d.SetId(caId)
	}

	// create tags
	createTagsHttpUrl := "v1/private-certificate-authorities/{ca_id}/tags/create"
	tags := d.Get("tags").(map[string]interface{})
	if err := createTags(d.Id(), client, tags, createTagsHttpUrl, "{ca_id}"); err != nil {
		return diag.FromErr(err)
	}

	// disable private CA
	if d.Get("action").(string) == "disable" {
		if err := disablePrivateCA(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePrivateCARead(ctx, d, meta)
}

func enablePrivateCA(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	enableHttpUrl := "v1/private-certificate-authorities/{ca_id}/enable"
	enablePath := client.Endpoint + enableHttpUrl
	enablePath = strings.ReplaceAll(enablePath, "{ca_id}", d.Id())
	enableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}

	_, err := client.Request("POST", enablePath, &enableOpt)
	if err != nil {
		return fmt.Errorf("error enabling CCM private CA: %s", err)
	}
	return nil
}

func disablePrivateCA(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	disableHttpUrl := "v1/private-certificate-authorities/{ca_id}/disable"
	disablePath := client.Endpoint + disableHttpUrl
	disablePath = strings.ReplaceAll(disablePath, "{ca_id}", d.Id())
	disableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}

	_, err := client.Request("POST", disablePath, &disableOpt)
	if err != nil {
		return fmt.Errorf("error disabling CCM private CA: %s", err)
	}
	return nil
}

func readPrivateCA(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getPrivateCAPath := client.Endpoint + "v1/private-certificate-authorities/{ca_id}"
	getPrivateCAPath = strings.ReplaceAll(getPrivateCAPath, "{ca_id}", d.Id())
	getPrivateCAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPrivateCAPath, &getPrivateCAOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(getResp)
}

func resourcePrivateCARead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ccm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	getRespBody, err := readPrivateCA(client, d)
	if err != nil {
		// When the resource does not exist, the response status code of the query API is 400. The response body example
		// is: {"error_code": "PCA.10010002","error_msg": "XXX"}
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "PCA.10010002"),
			"error retrieving CCM private CA")
	}

	status := utils.PathSearch("status", getRespBody, nil).(string)
	if status == "DELETED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving private CA")
	}

	getTagsHttpUrl := "v1/private-certificate-authorities/{ca_id}/tags"
	tags, err := getTags(d.Id(), client, getTagsHttpUrl, "{ca_id}")
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("type", utils.PathSearch("type", getRespBody, nil)),
		d.Set("distinguished_name", flattenDistinguishedName(getRespBody)),
		d.Set("key_algorithm", utils.PathSearch("key_algorithm", getRespBody, nil)),
		d.Set("signature_algorithm", utils.PathSearch("signature_algorithm", getRespBody, nil)),
		d.Set("crl_configuration", flattenCrlConfiguration(getRespBody)),
		d.Set("issuer_id", utils.PathSearch("issuer_id", getRespBody, nil)),
		d.Set("issuer_name", utils.PathSearch("issuer_name", getRespBody, nil)),
		d.Set("path_length", utils.PathSearch("path_length", getRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("charging_mode", flattenChargingMode(getRespBody)),
		d.Set("gen_mode", utils.PathSearch("gen_mode", getRespBody, nil)),
		d.Set("serial_number", utils.PathSearch("serial_number", getRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", getRespBody, float64(0)).(float64))/1000, false)),
		d.Set("expired_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("not_after", getRespBody, float64(0)).(float64))/1000, false)),
		d.Set("free_quota", utils.PathSearch("free_quota", getRespBody, nil)),
		d.Set("tags", tags),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CCM private CA fields: %s", err)
	}

	return nil
}

func flattenChargingMode(getRespBody interface{}) string {
	if utils.PathSearch("charging_mode", getRespBody, float64(0)).(float64) == 1 {
		return "postPaid"
	}
	return "prePaid"
}

func flattenDistinguishedName(resp interface{}) []interface{} {
	curJson := utils.PathSearch("distinguished_name", resp, make(map[string]interface{}))
	rawMap := curJson.(map[string]interface{})
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"common_name":         rawMap["common_name"],
		"country":             rawMap["country"],
		"state":               rawMap["state"],
		"locality":            rawMap["locality"],
		"organization":        rawMap["organization"],
		"organizational_unit": rawMap["organizational_unit"],
	})
	return rst
}

func flattenCrlConfiguration(resp interface{}) []interface{} {
	curJson := utils.PathSearch("crl_configuration", resp, make(map[string]interface{}))
	rawMap := curJson.(map[string]interface{})
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"enabled":         rawMap["enabled"],
		"crl_name":        rawMap["crl_name"],
		"obs_bucket_name": rawMap["obs_bucket_name"],
		"valid_days":      rawMap["valid_days"],
		"crl_dis_point":   rawMap["crl_dis_point"],
	})
	return rst
}

func deletePrepaidPrivateCA(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) diag.Diagnostics {
	if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
		// When the resource does not exist, the response status code of the query API is 400. The response body example
		// is: {"error_code": "CBC.30000067","error_msg": "XXX"}
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CBC.30000067"),
			"error unsubscribing CCM private CA")
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			getRespBody, err := readPrivateCA(client, d)
			if err != nil {
				convertErr := common.ConvertExpected400ErrInto404Err(err, "error_code", "PCA.10010002")
				var err404 golangsdk.ErrDefault404
				if errors.As(convertErr, &err404) {
					return "deleted", "COMPLETED", nil
				}
				return getRespBody, "ERROR", err
			}

			status := utils.PathSearch("status", getRespBody, "").(string)
			if status == "" {
				return getRespBody, "ERROR", fmt.Errorf("attribute `status` is not found in API response")
			}

			if status == "DELETED" {
				return "deleted", "COMPLETED", nil
			}
			return "continue", "PENDING", nil
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        15 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCM prepaid private CA (%s) to be deleted: %s", d.Id(), err)
	}
	return nil
}

func deletePostpaidPrivateCA(client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	getRespBody, err := readPrivateCA(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "PCA.10010002"),
			"error retrieving CCM private CA")
	}
	status := utils.PathSearch("status", getRespBody, "").(string)
	// Only CA in `PENDING` or `DISABLED` status can be deleted.
	// When the CA is in `ACTIVED` or `EXPIRED` status, it needs to be disabled first and then deleted.
	if status == "ACTIVED" || status == "EXPIRED" {
		if err := disablePrivateCA(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if status == "DELETED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving CCM private CA")
	}

	deletePath := client.Endpoint + "v1/private-certificate-authorities/{ca_id}"
	deletePath = strings.ReplaceAll(deletePath, "{ca_id}", d.Id())
	deletePath += fmt.Sprintf("?pending_days=%v", d.Get("pending_days"))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
			204,
		},
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting postpaid private CA: %s", err)
	}
	return nil
}

func resourcePrivateCADelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ccm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	// if charging mode is pre-paid, unsubscribe the order.
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		return deletePrepaidPrivateCA(ctx, client, d, cfg)
	}
	return deletePostpaidPrivateCA(client, d)
}

func resourcePrivateCAUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("ccm", region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	// update tags
	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})

		// remove old tags
		if len(oMap) > 0 {
			deleteTagsHttpUrl := "v1/private-certificate-authorities/{ca_id}/tags/delete"
			if err = deleteTags(d.Id(), client, oMap, deleteTagsHttpUrl, "{ca_id}"); err != nil {
				return diag.FromErr(err)
			}
		}

		// set new tags
		if len(nMap) > 0 {
			createTagsHttpUrl := "v1/private-certificate-authorities/{ca_id}/tags/create"
			if err := createTags(d.Id(), client, nMap, createTagsHttpUrl, "{ca_id}"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("crl_configuration.0.enabled") {
		if err := updateCRLConfiguration(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("action") {
		var actionErr error
		switch d.Get("action").(string) {
		case "enable":
			actionErr = enablePrivateCA(client, d)
		case "disable":
			actionErr = disablePrivateCA(client, d)
		}

		if actionErr != nil {
			return diag.FromErr(actionErr)
		}
	}
	return resourcePrivateCARead(ctx, d, meta)
}

func buildEnableConfigurationBodyParams(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"crl_name":        utils.ValueIgnoreEmpty(rawMap["crl_name"]),
		"obs_bucket_name": rawMap["obs_bucket_name"],
		"valid_days":      rawMap["valid_days"],
	}
}

func enableCRLConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	enableHttpUrl := "v1/private-certificate-authorities/{ca_id}/crl/enable"
	enablePath := client.Endpoint + enableHttpUrl
	enablePath = strings.ReplaceAll(enablePath, "{ca_id}", d.Id())
	enableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
		JSONBody:         utils.RemoveNil(buildEnableConfigurationBodyParams(d.Get("crl_configuration"))),
	}

	_, err := client.Request("POST", enablePath, &enableOpt)
	if err != nil {
		return fmt.Errorf("error enabling CRL configuration of CCM private CA (%s): %s", d.Id(), err)
	}
	return nil
}

func disableCRLConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	disableHttpUrl := "v1/private-certificate-authorities/{ca_id}/crl/disable"
	disablePath := client.Endpoint + disableHttpUrl
	disablePath = strings.ReplaceAll(disablePath, "{ca_id}", d.Id())
	disableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}

	_, err := client.Request("POST", disablePath, &disableOpt)
	if err != nil {
		return fmt.Errorf("error disabling CRL configuration of CCM private CA (%s): %s", d.Id(), err)
	}
	return nil
}

func updateCRLConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	enabled := d.Get("crl_configuration.0.enabled").(bool)
	if enabled {
		return enableCRLConfiguration(client, d)
	}

	return disableCRLConfiguration(client, d)
}

func createTags(id string, createTagsClient *golangsdk.ServiceClient, tags map[string]interface{},
	createTagsHttpUrl, idParamName string) error {
	if len(tags) > 0 {
		createTagsPath := createTagsClient.Endpoint + createTagsHttpUrl
		createTagsPath = strings.ReplaceAll(createTagsPath, idParamName, id)
		createTagsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		createTagsOpt.JSONBody = map[string]interface{}{
			"tags": utils.ExpandResourceTags(tags),
		}

		_, err := createTagsClient.Request("POST", createTagsPath, &createTagsOpt)
		if err != nil {
			return fmt.Errorf("error creating tags: %s", err)
		}
	}
	return nil
}

func getTags(id string, getTagsClient *golangsdk.ServiceClient, getTagsHttpUrl, idParamName string) (
	map[string]interface{}, error) {
	getTagsPath := getTagsClient.Endpoint + getTagsHttpUrl
	getTagsPath = strings.ReplaceAll(getTagsPath, idParamName, id)
	getTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTagsResp, err := getTagsClient.Request("GET", getTagsPath, &getTagsOpt)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags: %s", err)
	}
	getTagsRespBody, err := utils.FlattenResponse(getTagsResp)
	if err != nil {
		return nil, err
	}
	tags := utils.PathSearch("tags", getTagsRespBody, make([]interface{}, 0)).([]interface{})
	result := make(map[string]interface{})
	for _, val := range tags {
		valMap := val.(map[string]interface{})
		result[valMap["key"].(string)] = valMap["value"]
	}

	return result, nil
}

func deleteTags(id string, deleteTagsClient *golangsdk.ServiceClient, tags map[string]interface{},
	deleteTagsHttpUrl, idParamName string) error {
	if len(tags) == 0 {
		return nil
	}
	deleteTagsPath := deleteTagsClient.Endpoint + deleteTagsHttpUrl
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, idParamName, id)
	deleteTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	deleteTagsOpt.JSONBody = map[string]interface{}{
		"tags": utils.ExpandResourceTags(tags),
	}

	_, err := deleteTagsClient.Request("DELETE", deleteTagsPath, &deleteTagsOpt)
	if err != nil {
		return fmt.Errorf("error deleting tags: %s", err)
	}
	return nil
}
