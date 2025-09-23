package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/domains/https-certificate-info
func DataSourceDomainCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDomainCertificatesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_certificates": {
				Type:     schema.TypeList,
				Elem:     domainCertificateSchema(),
				Computed: true,
			},
		},
	}
}

func domainCertificateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_source": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"expire_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"https_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"force_redirect_https": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http2_enabled": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildCertificateQueryParams(conf *config.Config, d *schema.ResourceData) string {
	rst := "?page_size=10"
	if epsId := conf.GetEnterpriseProjectID(d); epsId != "" {
		rst += fmt.Sprintf("&enterprise_project_id=%v", epsId)
	}

	if v, ok := d.GetOk("name"); ok {
		rst += fmt.Sprintf("&domain_name=%v", v)
	}
	return rst
}

func resourceDomainCertificatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf         = meta.(*config.Config)
		region       = conf.GetRegion(d)
		mErr         *multierror.Error
		httpUrl      = "v1.0/cdn/domains/https-certificate-info"
		product      = "cdn"
		currentTotal = 1
		rst          = make([]interface{}, 0)
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildCertificateQueryParams(conf, d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&page_number=%v", requestPath, currentTotal)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.FromErr(err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		domainCertificates := utils.PathSearch("https", respBody, make([]interface{}, 0)).([]interface{})
		if len(domainCertificates) == 0 {
			break
		}
		rst = append(rst, domainCertificates...)
		currentTotal++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("domain_certificates", flattenListCertificateDomainsBody(rst)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListCertificateDomainsBody(domainCertificates []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(domainCertificates))
	for _, v := range domainCertificates {
		expirationTime := utils.PathSearch("expiration_time", v, 0)
		rst = append(rst, map[string]interface{}{
			"domain_id":            utils.PathSearch("domain_id", v, nil),
			"domain_name":          utils.PathSearch("domain_name", v, nil),
			"certificate_name":     utils.PathSearch("cert_name", v, nil),
			"certificate_body":     utils.PathSearch("certificate", v, nil),
			"certificate_source":   utils.PathSearch("certificate_type", v, nil),
			"expire_at":            utils.FormatTimeStampRFC3339(int64(expirationTime.(float64))/1000, false),
			"https_status":         utils.PathSearch("https_status", v, nil),
			"force_redirect_https": utils.PathSearch("force_redirect_https", v, nil),
			"http2_enabled":        utils.PathSearch("http2", v, nil),
		})
	}
	return rst
}
