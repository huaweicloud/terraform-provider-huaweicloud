package waf

import (
	"time"

	"github.com/chnsz/golangsdk/openstack/waf/v1/certificates"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const (
	// EXP_STATUS_NOT_EXPIRED not expired
	EXP_STATUS_NOT_EXPIRED = 0
	// EXP_STATUS_EXPIRED has expired
	EXP_STATUS_EXPIRED = 1
	// EXP_STATUS_EXPIRED_SOON will expire soon
	EXP_STATUS_EXPIRED_SOON = 2

	DEFAULT_PAGE_NUM  = 1
	DEFAULT_PAGE_SIZE = 5
)

func DataSourceWafCertificateV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceWafCertificateV1Read,

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
				Default:  EXP_STATUS_NOT_EXPIRED,
				ValidateFunc: validation.IntInSlice([]int{
					EXP_STATUS_NOT_EXPIRED, EXP_STATUS_EXPIRED, EXP_STATUS_EXPIRED_SOON,
				}),
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceWafCertificateV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	expStatus := d.Get("expire_status").(int)
	listOpts := certificates.ListOpts{
		Page:      DEFAULT_PAGE_NUM,
		Pagesize:  DEFAULT_PAGE_SIZE,
		Name:      d.Get("name").(string),
		ExpStatus: &expStatus,
	}

	page, err := certificates.List(wafClient, listOpts).AllPages()

	listCertificates, err := certificates.ExtractCertificates(page)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve certificates: %s", err)
	}
	logp.Printf("[DEBUG] Get certificate list: %#v", listCertificates)

	if len(listCertificates) > 0 {
		c := listCertificates[0]
		d.SetId(c.Id)
		d.Set("name", c.Name)
		d.Set("expire_status", c.ExpStatus)

		expires := time.Unix(int64(c.ExpireTime/1000), 0).UTC().Format("2006-01-02 15:04:05 MST")
		d.Set("expiration", expires)
	} else {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	return nil
}
