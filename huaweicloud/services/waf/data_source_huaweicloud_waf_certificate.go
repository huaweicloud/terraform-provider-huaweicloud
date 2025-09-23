package waf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF GET /v1/{project_id}/waf/certificate
func DataSourceWafCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafCertificateRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// In some scenarios, the attribute value of this field will be empty in API response.
			"expiration_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// In some scenarios, the attribute value of this field will be empty in API response.
			"expired_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Deprecated; Reasons for abandonment are as follows:
			// `expire_status`: Default value of this field is empty, not zero.
			"expire_status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "Use 'expiration_status' instead. ",
				Description: `schema: Deprecated; The certificate expiration status.`,
			},
			// `expiration`: Uniformly use dates in RFC3339 format.
			"expiration": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "Use 'expired_at' instead. ",
				Description: `schema: Deprecated; The certificate expiration time.`,
			},
		},
	}
}

// Paging is not considered cause this resource is a singular resource.
func buildCertificateQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	if epsId := cfg.GetEnterpriseProjectID(d); epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%s", res, epsId)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("expiration_status"); ok {
		res = fmt.Sprintf("%s&exp_status=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func dataSourceWafCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/certificate"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildCertificateQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving WAF certificates: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificatesResp := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
	if len(certificatesResp) == 0 {
		return diag.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	id := utils.PathSearch("[0]|id", certificatesResp, "").(string)
	if id == "" {
		return diag.Errorf("error retrieving WAF certificates: Certificate ID is not found in API response")
	}

	d.SetId(id)

	expireTimestamp := utils.PathSearch("[0]|expire_time", certificatesResp, float64(0)).(float64)
	createTimestamp := utils.PathSearch("[0]|timestamp", certificatesResp, float64(0)).(float64)

	mErr := multierror.Append(
		nil,
		d.Set("name", utils.PathSearch("[0]|name", certificatesResp, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("[0]|enterprise_project_id", certificatesResp, nil)),
		d.Set("expiration_status", flattenCertificateExpirationStatus(certificatesResp)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(createTimestamp)/1000, true)),
		d.Set("expired_at", utils.FormatTimeStampRFC3339(int64(expireTimestamp)/1000, true)),
		// Keep historical code logic
		d.Set("expiration", utils.FormatTimeStampRFC3339(int64(expireTimestamp)/1000, true, "2006-01-02 15:04:05 MST")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// This field may not be returned by the API.
// When the field is not returned, the default value is empty instead of `0`.
func flattenCertificateExpirationStatus(certificatesResp interface{}) string {
	expStatus := utils.PathSearch("[0]|exp_status", certificatesResp, nil)
	if expStatus == nil {
		return ""
	}
	return fmt.Sprintf("%v", expStatus)
}
