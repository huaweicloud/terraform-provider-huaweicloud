package ccm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

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
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crl_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"obs_bucket_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"valid_days": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"crl_dis_point": {
							Type:     schema.TypeString,
							Computed: true,
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

func resourcePrivateCACreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Get("type").(string) == "SUBORDINATE" {
		if _, ok := d.GetOk("issuer_id"); !ok {
			return diag.Errorf("error: required parameter [%s] for creating subordinate CA is not set", "issuer_id")
		}
	}

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		createPrivateCAHttpUrl = "v1/private-certificate-authorities"
		createPrivateCAProduct = "ccm"
	)
	createPrivateCAClient, err := cfg.NewServiceClient(createPrivateCAProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}
	createPrivateCAPath := createPrivateCAClient.Endpoint + createPrivateCAHttpUrl
	createPrivateCAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// if charging mode is pre-paid, need to order private CA first and then activate it.
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		// order private CA
		orderPrivateCAPath := createPrivateCAPath + "/order"
		orderPrivateCAOpt := createPrivateCAOpt
		parms, err := buildOrderPrivateCABodyParams(d)
		if err != nil {
			return diag.FromErr(err)
		}
		orderPrivateCAOpt.JSONBody = utils.RemoveNil(parms)
		orderPrivateCAResp, err := createPrivateCAClient.Request("POST", orderPrivateCAPath, &orderPrivateCAOpt)
		if err != nil {
			return diag.Errorf("error orderring CCM private CA: %s", err)
		}
		orderPrivateCARespBody, err := utils.FlattenResponse(orderPrivateCAResp)
		if err != nil {
			return diag.FromErr(err)
		}
		ids, err := jmespath.Search("ca_ids", orderPrivateCARespBody)
		if err != nil {
			return diag.Errorf("error orderring CCM private CA: ID is not found in API response")
		}
		id := ids.([]interface{})[0]
		d.SetId(id.(string))
		orderID, err := jmespath.Search("order_id", orderPrivateCARespBody)
		if err != nil {
			return diag.Errorf("error orderring CCM private CA: order ID is not found in API response")
		}

		// wait for order success
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderID.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderID.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for order resource %s complete: %s", orderID.(string), err)
		}

		// activate private CA
		activePrivateCAPath := createPrivateCAPath + "/{ca_id}/activate"
		activePrivateCAPath = strings.ReplaceAll(activePrivateCAPath, "{ca_id}", d.Id())
		activePrivateCAOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		activePrivateCAOpt.JSONBody = utils.RemoveNil(buildCreatePrivateCABodyParams(d, cfg))
		_, err = createPrivateCAClient.Request("POST", activePrivateCAPath, &activePrivateCAOpt)
		if err != nil {
			return diag.Errorf("error activating CCM private CA: %s", err)
		}

		createTagsHttpUrl := "v1/private-certificate-authorities/{ca_id}/tags/create"
		tags := d.Get("tags").(map[string]interface{})
		if err := createTags(id.(string), createPrivateCAClient, tags, createTagsHttpUrl, "{ca_id}"); err != nil {
			return diag.FromErr(err)
		}

		if d.Get("action").(string) == "disable" {
			if err := disablePrivateCA(createPrivateCAClient, d); err != nil {
				return diag.FromErr(err)
			}
		}

		return resourcePrivateCARead(ctx, d, meta)
	}

	// if charging mode is post-paid, then directly create
	createPrivateCAOpt.JSONBody = utils.RemoveNil(buildCreatePrivateCABodyParams(d, cfg))
	createPrivateCAResp, err := createPrivateCAClient.Request("POST", createPrivateCAPath, &createPrivateCAOpt)
	if err != nil {
		return diag.Errorf("error creating private CA: %s", err)
	}
	createPrivateCARespBody, err := utils.FlattenResponse(createPrivateCAResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("ca_id", createPrivateCARespBody)
	if err != nil {
		return diag.Errorf("error creating CCM private CA: ID is not found in API response")
	}
	d.SetId(id.(string))

	createTagsHttpUrl := "v1/private-certificate-authorities/{ca_id}/tags/create"
	tags := d.Get("tags").(map[string]interface{})
	if err := createTags(id.(string), createPrivateCAClient, tags, createTagsHttpUrl, "{ca_id}"); err != nil {
		return diag.FromErr(err)
	}

	if d.Get("action").(string) == "disable" {
		if err := disablePrivateCA(createPrivateCAClient, d); err != nil {
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

func buildOrderPrivateCABodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	rawArray := d.Get("validity").([]interface{})
	raw := rawArray[0].(map[string]interface{})

	// precheck period type
	rawTyepe := raw["type"].(string)
	if rawTyepe == "DAY" || rawTyepe == "HOUR" {
		return nil, fmt.Errorf("error: required parameter [%s] for creating pre-paid CA is invalid",
			"validity.period_type")
	}

	periodType := 2
	if rawTyepe == "YEAR" {
		periodType = 3
	}
	autoRenew := 0
	if val, ok := d.GetOk("auto_renew"); ok && val == "true" {
		autoRenew = 1
	}

	var prodecutInfos []map[string]interface{}
	prodecutInfos = append(prodecutInfos, map[string]interface{}{
		"cloud_service_type": "hws.service.type.ccm",
		"resource_type":      "hws.resource.type.pca.duration",
		"resource_spec_code": "ca.duration",
	})
	bodyParams := map[string]interface{}{
		"cloud_service_type": "hws.service.type.ccm",
		"charging_mode":      0,
		"period_type":        periodType,
		"period_num":         raw["value"].(int),
		"is_auto_renew":      autoRenew,
		"is_auto_pay":        1,
		"subscription_num":   1,
		"product_infos":      prodecutInfos,
	}
	return bodyParams, nil
}

func buildCreatePrivateCABodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
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

func buildPrivateCARequestBodyDistinguishedName(rawParams interface{}) map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	raw := rawArray[0].(map[string]interface{})
	params := map[string]interface{}{
		"common_name":         raw["common_name"],
		"country":             raw["country"],
		"state":               raw["state"],
		"locality":            raw["locality"],
		"organization":        raw["organization"],
		"organizational_unit": raw["organizational_unit"],
	}
	return params
}

func buildPrivateCARequestBodyValidity(rawParams interface{}) map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	raw := rawArray[0].(map[string]interface{})
	params := map[string]interface{}{
		"type":       raw["type"],
		"value":      raw["value"],
		"start_from": utils.ValueIgnoreEmpty(raw["started_at"]),
	}
	return params
}

func buildPrivateCARequestBodyKeyUsages(caType interface{}, rawParams interface{}) []interface{} {
	rawArray, _ := rawParams.([]interface{})
	if caType.(string) == "ROOT" || len(rawArray) == 0 {
		return []interface{}{"digitalSignature", "keyCertSign", "cRLSign"}
	}
	return rawArray
}

func buildPrivateCARequestBodyCrlConfiguration(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"enabled":         true,
			"crl_name":        raw["crl_name"],
			"obs_bucket_name": raw["obs_bucket_name"],
			"valid_days":      raw["valid_days"],
		}
		return params
	}
	return nil
}

func resourcePrivateCARead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getPrivateCAHttpUrl = "v1/private-certificate-authorities/{ca_id}"
		getPrivateCAProduct = "ccm"
	)
	getPrivateCAClient, err := cfg.NewServiceClient(getPrivateCAProduct, region)
	if err != nil {
		return diag.Errorf("error getting CCM client: %s", err)
	}

	getPrivateCAPath := getPrivateCAClient.Endpoint + getPrivateCAHttpUrl
	getPrivateCAPath = strings.ReplaceAll(getPrivateCAPath, "{ca_id}", d.Id())

	getPrivateCAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getPrivateCAResp, err := getPrivateCAClient.Request("GET", getPrivateCAPath, &getPrivateCAOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving private CA")
	}
	getPrivateCARespBody, err := utils.FlattenResponse(getPrivateCAResp)
	if err != nil {
		return diag.FromErr(err)
	}

	status := utils.PathSearch("status", getPrivateCARespBody, nil).(string)
	if status == "DELETED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving private CA")
	}

	chargingMode := "prePaid"
	if utils.PathSearch("charging_mode", getPrivateCARespBody, float64(0)).(float64) == 1 {
		chargingMode = "postPaid"
	}

	getTagsHttpUrl := "v1/private-certificate-authorities/{ca_id}/tags"
	tags, err := getTags(d.Id(), getPrivateCAClient, getTagsHttpUrl, "{ca_id}")
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("type", utils.PathSearch("type", getPrivateCARespBody, nil)),
		d.Set("distinguished_name", flattenDistinguishedName(getPrivateCARespBody)),
		d.Set("key_algorithm", utils.PathSearch("key_algorithm", getPrivateCARespBody, nil)),
		d.Set("signature_algorithm", utils.PathSearch("signature_algorithm", getPrivateCARespBody, nil)),
		d.Set("crl_configuration", flattenCrlConfiguration(getPrivateCARespBody)),
		d.Set("issuer_id", utils.PathSearch("issuer_id", getPrivateCARespBody, nil)),
		d.Set("issuer_name", utils.PathSearch("issuer_name", getPrivateCARespBody, nil)),
		d.Set("path_length", utils.PathSearch("path_length", getPrivateCARespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getPrivateCARespBody, nil)),
		d.Set("status", utils.PathSearch("status", getPrivateCARespBody, nil)),
		d.Set("charging_mode", chargingMode),
		d.Set("gen_mode", utils.PathSearch("gen_mode", getPrivateCARespBody, nil)),
		d.Set("serial_number", utils.PathSearch("serial_number", getPrivateCARespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", getPrivateCARespBody, float64(0)).(float64))/1000, false)),
		d.Set("expired_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("not_after", getPrivateCARespBody, float64(0)).(float64))/1000, false)),
		d.Set("free_quota", utils.PathSearch("free_quota", getPrivateCARespBody, nil)),
		d.Set("tags", tags),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CCM private CA fields: %s", err)
	}

	return nil
}

func flattenDistinguishedName(resp interface{}) []interface{} {
	curJson := utils.PathSearch("distinguished_name", resp, make([]interface{}, 0))
	curArray := curJson.(map[string]interface{})
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"common_name":         curArray["common_name"],
		"country":             curArray["country"],
		"state":               curArray["state"],
		"locality":            curArray["locality"],
		"organization":        curArray["organization"],
		"organizational_unit": curArray["organizational_unit"],
	})
	return rst
}

func flattenCrlConfiguration(resp interface{}) []interface{} {
	curJson := utils.PathSearch("crl_configuration", resp, make([]interface{}, 0))
	curArray := curJson.(map[string]interface{})
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"crl_name":        curArray["crl_name"],
		"obs_bucket_name": curArray["obs_bucket_name"],
		"valid_days":      curArray["valid_days"],
		"crl_dis_point":   curArray["crl_dis_point"],
	})
	return rst
}

func resourcePrivateCADelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		privateCAHttpUrl = "v1/private-certificate-authorities/{ca_id}"
		privateCAProduct = "ccm"
	)
	privateCAClient, err := cfg.NewServiceClient(privateCAProduct, region)
	if err != nil {
		return diag.Errorf("error deleting CCM client: %s", err)
	}
	privateCAPath := privateCAClient.Endpoint + privateCAHttpUrl
	privateCAPath = strings.ReplaceAll(privateCAPath, "{ca_id}", d.Id())
	privateCAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
			204,
		},
	}

	// if charging mode is pre-paid, unsubscribe the order.
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribing CCM private CA: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"ACTIVED", "DISABLED", "EXPIRED"},
			Target:       []string{"DELETED"},
			Refresh:      privateCAStatusRefreshFunc(d.Id(), region, cfg),
			Timeout:      d.Timeout(schema.TimeoutDelete),
			Delay:        15 * time.Second,
			PollInterval: 10 * time.Second,
		}

		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error deleting private CA: %s", err)
		}

		return nil
	}

	// get and check CA status, if not expired or disable then disable it.
	getPrivateCAResp, err := privateCAClient.Request("GET", privateCAPath, &privateCAOpt)
	if err != nil {
		return diag.Errorf("error getting private CA: %s", err)
	}
	getPrivateCARespBody, err := utils.FlattenResponse(getPrivateCAResp)
	if err != nil {
		return diag.FromErr(err)
	}
	status := utils.PathSearch("status", getPrivateCARespBody, nil).(string)
	if !(status == "EXPIRED" || status == "DISABLED") {
		disablePrivateCAPath := privateCAPath + "/disable"
		_, err = privateCAClient.Request("POST", disablePrivateCAPath, &privateCAOpt)
		if err != nil {
			return diag.Errorf("error disabling private CA: %s", err)
		}
	}

	pendingDays := d.Get("pending_days")
	privateCAPath += fmt.Sprintf("?pending_days=%v", pendingDays)

	_, err = privateCAClient.Request("DELETE", privateCAPath, &privateCAOpt)
	if err != nil {
		return diag.Errorf("error deleting private CA: %s", err)
	}
	return nil
}

func privateCAStatusRefreshFunc(id, region string, cfg *config.Config) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPrivateCAHttpUrl := "v1/private-certificate-authorities/{ca_id}"
		getPrivateCAProduct := "ccm"
		getPrivateCAClient, err := cfg.NewServiceClient(getPrivateCAProduct, region)
		if err != nil {
			return nil, "", fmt.Errorf("error creating CCM client: %s", err)
		}

		getPrivateCAPath := getPrivateCAClient.Endpoint + getPrivateCAHttpUrl
		getPrivateCAPath = strings.ReplaceAll(getPrivateCAPath, "{ca_id}", id)
		getPrivateCAOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getPrivateCAResp, err := getPrivateCAClient.Request("GET", getPrivateCAPath, &getPrivateCAOpt)
		if err != nil && hasErrorCode(err, "PCA.10010002") {
			return getPrivateCAResp, "DELETED", nil
		}
		getPrivateCARespBody, err := utils.FlattenResponse(getPrivateCAResp)
		if err != nil {
			return nil, "", err
		}
		status := utils.PathSearch("status", getPrivateCARespBody, "")
		return getPrivateCARespBody, status.(string), nil
	}
}

func resourcePrivateCAUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	privateCAClient, err := cfg.NewServiceClient("ccm", region)
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
			if err = deleteTags(d.Id(), privateCAClient, oMap, deleteTagsHttpUrl, "{ca_id}"); err != nil {
				return diag.FromErr(err)
			}
		}

		// set new tags
		if len(nMap) > 0 {
			createTagsHttpUrl := "v1/private-certificate-authorities/{ca_id}/tags/create"
			if err := createTags(d.Id(), privateCAClient, nMap, createTagsHttpUrl, "{ca_id}"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("action") {
		var actionErr error
		switch d.Get("action").(string) {
		case "enable":
			actionErr = enablePrivateCA(privateCAClient, d)
		case "disable":
			actionErr = disablePrivateCA(privateCAClient, d)
		}

		if actionErr != nil {
			return diag.FromErr(actionErr)
		}
	}
	return resourcePrivateCARead(ctx, d, meta)
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
