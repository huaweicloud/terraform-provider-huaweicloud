package waf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF GET /v1/{project_id}/waf/certificate
func DataSourceWafCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafCertificatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the data source. If omitted, the provider-level region will be used.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID of WAF certificate.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of certificate.`,
			},
			"host": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to obtain the domain name for which the certificate is used.`,
			},
			"expiration_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the certificate expiration status.`,
			},
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The certificate list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The certificate ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The certificate name.`,
						},
						"expiration_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The certificate expiration status.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the certificate was uploaded.`,
						},
						"expired_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the certificate expires.`,
						},
						"bind_host": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The domain information bound to the certificate.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The domain ID.`,
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The domain name.`,
									},
									"mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The special domain pattern.`,
									},
									"waf_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The deployment mode of WAF instance that is used for the domain name.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildWafCertificatesQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := "?pagesize=100"
	if epsId := cfg.GetEnterpriseProjectID(d); epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%s", res, epsId)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	// The type of this field is `boolean`, the default value is false, the default value is ignored.
	if v, ok := d.GetOk("host"); ok {
		res = fmt.Sprintf("%s&host=%v", res, v)
	}
	if v, ok := d.GetOk("expiration_status"); ok {
		res = fmt.Sprintf("%s&exp_status=%v", res, v)
	}

	return res
}

func dataSourceWafCertificatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		mErr              *multierror.Error
		httpUrl           = "v1/{project_id}/waf/certificate"
		product           = "waf"
		totalCertificates []interface{}
		currentPage       = 1
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildWafCertificatesQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithPage := fmt.Sprintf("%s&page=%d", requestPath, currentPage)
		resp, err := client.Request("GET", requestPathWithPage, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving WAF certificates: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		certificatesResp := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(certificatesResp) == 0 {
			break
		}

		totalCertificates = append(totalCertificates, certificatesResp...)
		currentPage++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("certificates", flattenWafCertificatesResp(totalCertificates)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataSourceCertificateExpirationStatus(certificateResp interface{}) string {
	expStatus := utils.PathSearch("exp_status", certificateResp, nil)
	if expStatus == nil {
		return ""
	}
	return fmt.Sprintf("%v", expStatus)
}

func flattenDataSourceCertificateBindHost(bindHosts []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(bindHosts))
	for _, v := range bindHosts {
		rst = append(rst, map[string]interface{}{
			"id":       utils.PathSearch("id", v, nil),
			"domain":   utils.PathSearch("hostname", v, nil),
			"mode":     utils.PathSearch("mode", v, nil),
			"waf_type": utils.PathSearch("waf_type", v, nil),
		})
	}
	return rst
}

func flattenWafCertificatesResp(totalCertificates []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(totalCertificates))
	for _, v := range totalCertificates {
		expireTimestamp := utils.PathSearch("expire_time", v, float64(0)).(float64)
		createTimestamp := utils.PathSearch("timestamp", v, float64(0)).(float64)
		bindHosts := utils.PathSearch("bind_host", v, make([]interface{}, 0)).([]interface{})
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"expiration_status": flattenDataSourceCertificateExpirationStatus(v),
			"created_at":        utils.FormatTimeStampRFC3339(int64(createTimestamp)/1000, true),
			"expired_at":        utils.FormatTimeStampRFC3339(int64(expireTimestamp)/1000, true),
			"bind_host":         flattenDataSourceCertificateBindHost(bindHosts),
		})
	}
	return rst
}
