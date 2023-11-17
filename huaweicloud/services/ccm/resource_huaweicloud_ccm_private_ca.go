package ccm

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourcePrivateCertificateAuthority() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateCACreate,
		ReadContext:   resourcePrivateCARead,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_mode": {
				Type:     schema.TypeInt,
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
		return diag.Errorf("error creating CCM Client: %s", err)
	}
	createPrivateCAPath := createPrivateCAClient.Endpoint + createPrivateCAHttpUrl
	createPrivateCAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

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
	return resourcePrivateCARead(ctx, d, meta)
}

func buildCreatePrivateCABodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":                d.Get("type"),
		"distinguished_name":  buildPrivateCARequestBodyDistinguishedName(d.Get("distinguished_name")),
		"key_algorithm":       d.Get("key_algorithm"),
		"signature_algorithm": d.Get("signature_algorithm"),
		"validity":            buildPrivateCARequestBodyValidity(d.Get("validity")),
		"issuer_id":           utils.ValueIngoreEmpty(d.Get("issuer_id")),
		"path_length":         utils.ValueIngoreEmpty(d.Get("path_length")),
		"key_usages":          buildPrivateCARequestBodyKeyUsages(d.Get("type"), d.Get("key_usages")),
		"crl_configuration": buildPrivateCARequestBodyCrlConfiguration(
			utils.ValueIngoreEmpty(d.Get("issuer_id")), d.Get("crl_configuration")),
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
		"start_from": utils.ValueIngoreEmpty(raw["started_at"]),
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

func buildPrivateCARequestBodyCrlConfiguration(issuerID interface{}, rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		name := utils.ValueIngoreEmpty(raw["crl_name"]).(string)
		if name == "" {
			name = issuerID.(string)
		}
		params := map[string]interface{}{
			"crl_name":        name,
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
		getPrivateCAHttpUrl = "v1/private-certificate-authorities/{id}"
		getPrivateCAProduct = "ccm"
	)
	getPrivateCAClient, err := cfg.NewServiceClient(getPrivateCAProduct, region)
	if err != nil {
		return diag.Errorf("error getting CCM Client: %s", err)
	}

	getPrivateCAPath := getPrivateCAClient.Endpoint + getPrivateCAHttpUrl
	getPrivateCAPath = strings.ReplaceAll(getPrivateCAPath, "{id}", d.Id())

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
		d.Set("charging_mode", utils.PathSearch("charging_mode", getPrivateCARespBody, nil)),
		d.Set("gen_mode", utils.PathSearch("gen_mode", getPrivateCARespBody, nil)),
		d.Set("serial_number", utils.PathSearch("serial_number", getPrivateCARespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", getPrivateCARespBody, float64(0)).(float64))/1000, false)),
		d.Set("expired_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("not_after", getPrivateCARespBody, float64(0)).(float64))/1000, false)),
		d.Set("free_quota", utils.PathSearch("free_quota", getPrivateCARespBody, nil)),
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

func resourcePrivateCADelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		privateCAHttpUrl = "v1/private-certificate-authorities/{id}"
		privateCAProduct = "ccm"
	)
	privateCAClient, err := cfg.NewServiceClient(privateCAProduct, region)
	if err != nil {
		return diag.Errorf("error deleting CCM Client: %s", err)
	}
	privateCAPath := privateCAClient.Endpoint + privateCAHttpUrl
	privateCAPath = strings.ReplaceAll(privateCAPath, "{id}", d.Id())
	privateCAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
			204,
		},
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
	if !(status == "EXPIRED" || status == "DISABLE") {
		disablePrivateCAPath := privateCAPath + "/disable"
		_, err = privateCAClient.Request("POST", disablePrivateCAPath, &privateCAOpt)
		if err != nil {
			return diag.Errorf("error disabling private CA: %s", err)
		}
	}

	privateCAOpt.JSONBody = utils.RemoveNil(
		map[string]interface{}{
			"pending_days": d.Get("pending_days"),
		})
	_, err = privateCAClient.Request("DELETE", privateCAPath, &privateCAOpt)
	if err != nil {
		return diag.Errorf("error deleting private CA: %s", err)
	}
	return nil
}
