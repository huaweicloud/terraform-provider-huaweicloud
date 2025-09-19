package apig

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

// @API APIG GET /v2/{project_id}/apigw/certificates/{certificate_id}/attached-domains
func DataSourceCertificateAssociatedDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCertificateAssociatedDomainsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the domains are located.`,
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the certificate associated with the domains.`,
			},
			"url_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The associated domain name to be queried.`,
			},
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the associated domain.`,
						},
						"url_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The associated domain name.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dedicated instance to which the domain belongs.`,
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The CNAME resolution status of the domain name.`,
						},
						"min_ssl_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The minimum SSL protocol version of the domain.",
						},
						"verified_client_certificate_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether client certificate verification is enabled.`,
						},
						"api_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the API group to which the domain belongs.`,
						},
						"api_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the API group to which the domain belongs.`,
						},
					},
				},
				Description: `All domains that match the filter parameters.`,
			},
		},
	}
}

func buildCertificateAssociatedDomainsBodyParams(d *schema.ResourceData) string {
	res := ""
	if domainName, ok := d.GetOk("url_domain"); ok {
		res = fmt.Sprintf("%s&url_domain=%v", res, domainName)
	}

	return res
}

func queryCertificateAssociatedDomains(client *golangsdk.ServiceClient, d *schema.ResourceData, certificateId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/certificates/{certificate_id}/attached-domains"
		// The default limit is 20.
		limit  = 100
		offset = 0
		result = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{certificate_id}", certificateId)
	listPath += fmt.Sprintf("?limit=%d", limit)
	listPath += buildCertificateAssociatedDomainsBodyParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		attachedDomains := utils.PathSearch("bound_domains", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, attachedDomains...)
		if len(attachedDomains) < limit {
			break
		}

		offset += len(attachedDomains)
	}
	return result, nil
}

func dataSourceCertificateAssociatedDomainsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		certificateId = d.Get("certificate_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	domains, err := queryCertificateAssociatedDomains(client, d, certificateId)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("domains", flattenCertificateAssociatedDomainInfos(domains)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCertificateAssociatedDomainInfos(domains []interface{}) []interface{} {
	if len(domains) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(domains))
	for _, domain := range domains {
		result = append(result, map[string]interface{}{
			"id":                                  utils.PathSearch("id", domain, nil),
			"url_domain":                          utils.PathSearch("url_domain", domain, nil),
			"instance_id":                         utils.PathSearch("instance_id", domain, nil),
			"status":                              utils.PathSearch("status", domain, nil),
			"min_ssl_version":                     utils.PathSearch("min_ssl_version", domain, nil),
			"verified_client_certificate_enabled": utils.PathSearch("verified_client_certificate_enabled", domain, nil),
			"api_group_id":                        utils.PathSearch("api_group_id", domain, nil),
			"api_group_name":                      utils.PathSearch("api_group_name", domain, nil),
		})
	}
	return result
}
