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

// @API APIG GET /v2/{project_id}/apigw/certificates
func DataSourceInstanceSSLCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceSSLCertificatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the dedicated instance is located that the SSL certificates are associated with.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the SSL certificates belong.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the SSL certificate.`,
			},
			"common_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The domain name of the SSL certificate.`,
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The signature algorithm of the SSL certificate.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The visibility range of the SSL certificate.`,
			},
			"algorithm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The algorithm type of the SSL certificate(RSA, ECC, SM2).`,
			},
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All SSL certificates that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the SSL certificate.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the tenant project`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dedicated instance to which the SSL certificates belong.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the SSL certificate.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the SSL certificate.`,
						},
						"common_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name of the SSL certificate.`,
						},
						"san": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The san extended domain of the SSL certificate.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"signature_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The signature algorithm of the SSL certificate.`,
						},
						"algorithm_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The algorithm type of the SSL certificate(RSA, ECC, SM2).`,
						},
						"is_has_trusted_root_ca": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `The certificate has trusted root certificate authority or not.`,
						},
						"not_after": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The expiration date of the SSL certificate.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The create time of the SSL certificate.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the SSL certificate.`,
						},
					},
				},
			},
		},
	}
}

func buildListInstanceSSLCertificatesParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("common_name"); ok {
		res = fmt.Sprintf("%s&common_name=%v", res, v)
	}
	if v, ok := d.GetOk("signature_algorithm"); ok {
		res = fmt.Sprintf("%s&signature_algorithm=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("algorithm_type"); ok {
		res = fmt.Sprintf("%s&algorithm_type=%v", res, v)
	}
	return res
}

func listInstanceSSLCertificates(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/certificates?instance_id={instance_id}"
		instanceId = d.Get("instance_id").(string)
		offset     = 0
		limit      = 100
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)

	listPath += buildListInstanceSSLCertificatesParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&limit=%d&offset=%d", listPath, limit, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error querying associated certificates under specified dedicated instance (%s): %s", instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		certificates := utils.PathSearch("certs", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, certificates...)
		if len(certificates) < limit {
			break
		}
		offset += len(certificates)
	}
	return result, nil
}

func dataSourceInstanceSSLCertificatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	certificates, err := listInstanceSSLCertificates(client, d)
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
		d.Set("certificates", flattenAssociatedCertificates(certificates)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssociatedCertificates(certificates []interface{}) []interface{} {
	if len(certificates) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(certificates))
	for _, certificate := range certificates {
		result = append(result, map[string]interface{}{
			"id":                     utils.PathSearch("sign_id", certificate, nil),
			"name":                   utils.PathSearch("sign_name", certificate, nil),
			"type":                   utils.PathSearch("sign_type", certificate, nil),
			"instance_id":            utils.PathSearch("instance_id", certificate, nil),
			"project_id":             utils.PathSearch("project_id", certificate, nil),
			"common_name":            utils.PathSearch("common_name", certificate, nil),
			"san":                    utils.PathSearch("san", certificate, nil),
			"not_after":              utils.PathSearch("not_after", certificate, nil),
			"signature_algorithm":    utils.PathSearch("signature_algorithm", certificate, nil),
			"algorithm_type":         utils.PathSearch("algorithm_type", certificate, nil),
			"is_has_trusted_root_ca": utils.PathSearch("is_has_trusted_root_ca", certificate, nil),
			"create_time":            utils.PathSearch("create_time", certificate, nil),
			"update_time":            utils.PathSearch("update_time", certificate, nil),
		})
	}
	return result
}
