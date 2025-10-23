package cdn

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
		ReadContext: dataSourceDomainCertificatesRead,

		Schema: map[string]*schema.Schema{
			// Optional parameters
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the domain.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project that the resource belongs.`,
			},

			// Attributes
			"domain_certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        domainCertificateSchema(),
				Description: `The list of certificates that are associated with the queried domain.`,
			},
		},
	}
}

func domainCertificateSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the domain.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the domain.`,
			},
			"certificate_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the certificate.`,
			},
			"certificate_body": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The content of the certificate.`,
			},
			"certificate_source": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The source type of the certificate.`,
			},
			"expire_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expiration time, in RFC3339 format.`,
			},
			"https_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Whether the HTTPS certificate is enabled.`,
			},
			"force_redirect_https": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Whether client requests are forced to be redirected.`,
			},
			"http2_enabled": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Whether the HTTP/2 protocol is used.`,
			},
		},
	}
}

func buildCertificateQueryParams(conf *config.Config, d *schema.ResourceData) string {
	res := ""

	if epsId := conf.GetEnterpriseProjectID(d); epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&domain_name=%v", res, v)
	}

	return res
}

func dataSourceDomainCertificatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		httpUrl    = "v1.0/cdn/domains/https-certificate-info?page_size={page_size}"
		pageNumber = 1
		pageSize   = 10
		result     = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	listPathWithSize := client.Endpoint + httpUrl
	listPathWithSize = strings.ReplaceAll(listPathWithSize, "{page_size}", strconv.Itoa(pageSize))
	listPathWithSize += buildCertificateQueryParams(cfg, d)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithPage := fmt.Sprintf("%s&page_number=%v", listPathWithSize, pageNumber)
		resp, err := client.Request("GET", listPathWithPage, &opt)
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
		result = append(result, domainCertificates...)
		pageNumber++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("domain_certificates", flattenListCertificateDomainsBody(result)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListCertificateDomainsBody(domainCertificates []interface{}) []interface{} {
	result := make([]interface{}, 0, len(domainCertificates))

	for _, v := range domainCertificates {
		expirationTime := utils.PathSearch("expiration_time", v, 0)
		result = append(result, map[string]interface{}{
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

	return result
}
