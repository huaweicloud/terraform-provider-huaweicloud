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

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCM GET /v3/scm/certificates
func DataSourceCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceCertificatesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deploy_support": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Elem:     certificatesCertificateSchema(),
				Computed: true,
			},
		},
	}
}

func certificatesCertificateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sans": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deploy_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"brand": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"validity_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wildcard_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func datasourceCertificatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v3/scm/certificates"
		product = "scm"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath += buildListCertificatesQueryParams(d)
	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving CCM SSL certificates: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("certificates", filterListCertificatesBodyCertificate(
			flattenListCertificatesBodyCertificate(listRespBody), d)),
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
