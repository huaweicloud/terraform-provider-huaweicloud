// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CCM
// ---------------------------------------------------------------

package ccm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCM GET /v3/scm/certificates
func DataSourceCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceCertificatesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Certificate status.`,
				ValidateFunc: validation.StringInSlice([]string{
					"ALL", "PAID", "ISSUED", "CHECKING", "CANCELCHECKING", "UNPASSED", "EXPIRED", "REVOKING", "REVOKED",
					"UPLOAD", "CHECKING_ORG", "ISSUING", "SUPPLEMENTCHECKING",
				}, false),
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The enterprise project id of the project.`,
			},
			"deploy_support": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to support deployment.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Certificate name.`,
			},
			"certificates": {
				Type:        schema.TypeList,
				Elem:        certificatesCertificateSchema(),
				Computed:    true,
				Description: `Certificate list. For details, see Data structure of the Certificate field.`,
			},
		},
	}
}

func certificatesCertificateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Certificate ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Certificate name.`,
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Domain name associated with the certificate.`,
			},
			"sans": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Additional domain name associated with the certificate.`,
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Signature algorithm.`,
			},
			"deploy_support": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to support deployment.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Certificate type.`,
			},
			"brand": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Certificate authority.`,
			},
			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Certificate expiration time.`,
			},
			"domain_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Domain name type.`,
			},
			"validity_period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Certificate validity period, in months.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Certificate status.`,
			},
			"domain_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of domain names that can be associated with the certificate.`,
			},
			"wildcard_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of wildcard domain names that can be associated with the certificate.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Certificate description.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project id of the project.`,
			},
		},
	}
	return &sc
}

func resourceCertificatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// listCertificates: Query the List of SCM certificates.
	var (
		listCertificatesHttpUrl = "v3/scm/certificates"
		listCertificatesProduct = "scm"
	)
	listCertificatesClient, err := conf.NewServiceClient(listCertificatesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Certificates Client: %s", err)
	}

	listCertificatesPath := listCertificatesClient.Endpoint + listCertificatesHttpUrl

	listCertificatesqueryParams := buildListCertificatesQueryParams(d)
	listCertificatesPath += listCertificatesqueryParams

	listCertificatesResp, err := pagination.ListAllItems(
		listCertificatesClient,
		"offset",
		listCertificatesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Certificates")
	}

	listCertificatesRespJson, err := json.Marshal(listCertificatesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listCertificatesRespBody interface{}
	err = json.Unmarshal(listCertificatesRespJson, &listCertificatesRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("certificates", filterListCertificatesBodyCertificate(
			flattenListCertificatesBodyCertificate(listCertificatesRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListCertificatesBodyCertificate(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("certificates", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"domain":                utils.PathSearch("domain", v, nil),
			"sans":                  utils.PathSearch("sans", v, nil),
			"signature_algorithm":   utils.PathSearch("signature_algorithm", v, nil),
			"deploy_support":        utils.PathSearch("deploy_support", v, nil),
			"type":                  utils.PathSearch("type", v, nil),
			"brand":                 utils.PathSearch("brand", v, nil),
			"expire_time":           utils.PathSearch("expire_time", v, nil),
			"domain_type":           utils.PathSearch("domain_type", v, nil),
			"validity_period":       utils.PathSearch("validity_period", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"domain_count":          utils.PathSearch("domain_count", v, nil),
			"wildcard_count":        utils.PathSearch("wildcard_count", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
		})
	}
	return rst
}

func filterListCertificatesBodyCertificate(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("name"); ok && param != utils.PathSearch("name", v, nil) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListCertificatesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}

	if v, ok := d.GetOk("deploy_support"); ok {
		res = fmt.Sprintf("%s&deploy_support=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
