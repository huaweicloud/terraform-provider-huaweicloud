package ccm

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCM POST /v1/private-certificates
// @API CCM POST /v1/private-certificates/{certificate_id}/tags/create
// @API CCM DELETE /v1/private-certificates/{certificate_id}/tags/delete
// @API CCM GET /v1/private-certificates/{certificate_id}
// @API CCM GET /v1/private-certificates/{certificate_id}/tags
// @API CCM DELETE /v1/private-certificates/{certificate_id}
func ResourcePrivateCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateCertificateCreate,
		ReadContext:   resourcePrivateCertificateRead,
		UpdateContext: resourcePrivateCertificateUpdate,
		DeleteContext: resourcePrivateCertificateDelete,
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
			"issuer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"distinguished_name": {
				Type:     schema.TypeList,
				Elem:     distinguishedName(),
				Required: true,
				ForceNew: true,
			},
			"validity": {
				Type:     schema.TypeList,
				Elem:     validity(),
				Required: true,
				ForceNew: true,
			},
			"subject_alternative_names": {
				Type:     schema.TypeList,
				Elem:     subjectAlternativeNameSchema(),
				Optional: true,
				ForceNew: true,
			},
			"key_usage": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
			},
			"server_auth": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"client_auth": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"code_signing": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"email_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"time_stamping": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"object_identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"object_identifier_value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
			"issuer_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expired_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gen_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func validity() *schema.Resource {
	sc := schema.Resource{
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

			"start_at": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
	return &sc
}
func distinguishedName() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"common_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"country": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"locality": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"organizational_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
	return &sc
}

func subjectAlternativeNameSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Type of the alternative name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Value of the corresponding alternative name type.`,
			},
		},
	}
	return &sc
}

func resourcePrivateCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                      = meta.(*config.Config)
		product                  = "ccm"
		createCertificateHttpUrl = "v1/private-certificates"
	)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	createCertificatePath := client.Endpoint + createCertificateHttpUrl
	createCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createCertificateOpt.JSONBody = utils.RemoveNil(buildCreateCertificateBodyParams(d, cfg))
	createCertificateResp, err := client.Request("POST", createCertificatePath, &createCertificateOpt)
	if err != nil {
		return diag.Errorf("error creating CCM private certificate: %s", err)
	}
	createCertificateRespBody, err := utils.FlattenResponse(createCertificateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certId := utils.PathSearch("certificate_id", createCertificateRespBody, "").(string)
	if certId == "" {
		return diag.Errorf("unable to find the CCM private certificate ID from the API response")
	}

	d.SetId(certId)

	// deal tags
	createTagsHttpUrl := "v1/private-certificates/{certificate_id}/tags/create"
	tags := d.Get("tags").(map[string]interface{})
	if err := createTags(certId, client, tags, createTagsHttpUrl, "{certificate_id}"); err != nil {
		return diag.FromErr(err)
	}

	return resourcePrivateCertificateRead(ctx, d, meta)
}

func buildCreateCertificateBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"issuer_id":                 d.Get("issuer_id"),
		"key_algorithm":             d.Get("key_algorithm"),
		"signature_algorithm":       d.Get("signature_algorithm"),
		"distinguished_name":        buildCertDistinguishedName(d.Get("distinguished_name")),
		"validity":                  buildValidity(d.Get("validity")),
		"enterprise_project_id":     utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"key_usages":                d.Get("key_usage"),
		"extended_key_usage":        buildExtendedKeyUsage(d),
		"customized_extension":      buildCustomizedExtension(d),
		"subject_alternative_names": buildSubjectAlternativeName(d),
	}
	return bodyParams
}

func buildCertDistinguishedName(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

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

func buildValidity(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	params := map[string]interface{}{
		"type":       raw["type"],
		"value":      raw["value"],
		"start_from": utils.ValueIgnoreEmpty(raw["start_at"]),
	}
	return params
}
func buildExtendedKeyUsage(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"server_auth":      d.Get("server_auth"),
		"client_auth":      d.Get("client_auth"),
		"code_signing":     d.Get("code_signing"),
		"email_protection": d.Get("email_protection"),
		"time_stamping":    d.Get("time_stamping"),
	}
	return bodyParams
}
func buildCustomizedExtension(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"object_identifier": d.Get("object_identifier"),
		"value":             d.Get("object_identifier_value"),
	}
	return bodyParams
}
func buildSubjectAlternativeName(d *schema.ResourceData) []interface{} {
	curJson := utils.PathSearch("subject_alternative_names", d, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"type":  utils.PathSearch("type", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}

func resourcePrivateCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		product = "ccm"
		region  = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})

		if len(oMap) > 0 {
			deleteTagsHttpUrl := "v1/private-certificates/{certificate_id}/tags/delete"
			if err = deleteTags(d.Id(), client, oMap, deleteTagsHttpUrl, "{certificate_id}"); err != nil {
				return diag.FromErr(err)
			}
		}

		if len(nMap) > 0 {
			createTagsHttpUrl := "v1/private-certificates/{certificate_id}/tags/create"
			if err := createTags(d.Id(), client, nMap, createTagsHttpUrl, "{certificate_id}"); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return resourcePrivateCertificateRead(ctx, d, meta)
}

func resourcePrivateCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                   = meta.(*config.Config)
		product               = "ccm"
		region                = cfg.GetRegion(d)
		getCertificateHttpUrl = "v1/private-certificates/{certificate_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	getCertificatePath := client.Endpoint + getCertificateHttpUrl
	getCertificatePath = strings.ReplaceAll(getCertificatePath, "{certificate_id}", d.Id())
	getCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getCertificateResp, err := client.Request("GET", getCertificatePath, &getCertificateOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "PCA.10010002"),
			"error retrieving CCM private certificate")
	}

	getCertificateRespBody, err := utils.FlattenResponse(getCertificateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	getTagsHttpUrl := "v1/private-certificates/{certificate_id}/tags"
	tags, err := getTags(d.Id(), client, getTagsHttpUrl, "{certificate_id}")
	if err != nil {
		return diag.FromErr(err)
	}

	created := utils.PathSearch("create_time", getCertificateRespBody, float64(0)).(float64)
	started := utils.PathSearch("not_before", getCertificateRespBody, float64(0)).(float64)
	expired := utils.PathSearch("not_after", getCertificateRespBody, float64(0)).(float64)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("distinguished_name", flattenDistinguishedName(getCertificateRespBody)),
		d.Set("issuer_id", utils.PathSearch("issuer_id", getCertificateRespBody, nil)),
		d.Set("key_algorithm", utils.PathSearch("key_algorithm", getCertificateRespBody, nil)),
		d.Set("signature_algorithm", utils.PathSearch("signature_algorithm", getCertificateRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getCertificateRespBody, nil)),
		d.Set("issuer_name", utils.PathSearch("issuer_name", getCertificateRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getCertificateRespBody, nil)),
		d.Set("gen_mode", utils.PathSearch("gen_mode", getCertificateRespBody, nil)),
		d.Set("start_at", utils.FormatTimeStampRFC3339(int64(started)/1000, false)),
		d.Set("expired_at", utils.FormatTimeStampRFC3339(int64(expired)/1000, false)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(created)/1000, false)),
		d.Set("tags", tags),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CCM private certificate fields: %s", err)
	}

	return nil
}

func resourcePrivateCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                   = meta.(*config.Config)
		product               = "ccm"
		delCertificateHttpUrl = "v1/private-certificates/{certificate_id}"
	)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	delCertificatePath := client.Endpoint + delCertificateHttpUrl
	delCertificatePath = strings.ReplaceAll(delCertificatePath, "{certificate_id}", d.Id())
	delCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	_, err = client.Request("DELETE", delCertificatePath, &delCertificateOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "PCA.10010002"),
			"error deleting CCM private certificate")
	}
	return nil
}
