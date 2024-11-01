package waf

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/waf/v1/certificates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
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
				Required: true,
			},
			"expire_status": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceWafCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	wafClient, err := conf.WafV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	expStatus := d.Get("expire_status").(int)
	listOpts := certificates.ListOpts{
		Page:                1,
		Pagesize:            5,
		Name:                d.Get("name").(string),
		ExpStatus:           &expStatus,
		EnterpriseProjectID: conf.GetEnterpriseProjectID(d),
	}

	page, err := certificates.List(wafClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("error retrieving WAF certificates: %s", err)
	}

	listCertificates, err := certificates.ExtractCertificates(page)
	if err != nil {
		return diag.Errorf("Unable to retrieve certificates: %s", err)
	}
	log.Printf("[DEBUG] Get certificate list: %#v", listCertificates)

	if len(listCertificates) == 0 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	c := listCertificates[0]
	d.SetId(c.Id)
	expires := time.Unix(int64(c.ExpireTime/1000), 0).UTC().Format("2006-01-02 15:04:05 MST")
	mErr := multierror.Append(
		nil,
		d.Set("name", c.Name),
		d.Set("expire_status", c.ExpStatus),
		d.Set("expiration", expires),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
