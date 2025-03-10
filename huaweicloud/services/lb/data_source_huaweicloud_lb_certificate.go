package lb

import (
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/elb/v2/certificates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// @API ELB GET /v2/{project_id}/elb/certificates
func DataSourceLBCertificateV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLBCertificateV2Read,

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
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLBCertificateV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)

	// LoadBalancerClient catalog info: Name is "elb" and Version is "v2"
	client, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud LoadBalancer Client: %s", err)
	}

	listOpts := certificates.ListOpts{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	r, err := certificates.List(client, listOpts)
	certRst, err := r.Extract()
	if err != nil {
		return fmtp.Errorf("Unable to retrieve certificates from LoadBalancer: %s", err)
	}
	logp.Printf("[DEBUG] Get certificate list: %#v", certRst)

	if len(certRst.Certificates) > 0 {
		err = setCertificateAttributes(d, certRst.Certificates[0])
		if err != nil {
			return err
		}
	} else {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	return nil
}

func setCertificateAttributes(d *schema.ResourceData, c certificates.Certificate) error {
	d.SetId(c.Id)

	var expiration string
	tm, err := time.Parse("2006-01-02 15:04:05", c.ExpireTime)
	if err != nil {
		// If the format of ExpireTime is not expected, set the original value directly.
		expiration = c.ExpireTime
		logp.Printf("[WAIN] The format of the ExpireTime field of the LB certificate is not expected: %s",
			c.ExpireTime)
	} else {
		expiration = tm.Format("2006-01-02 15:04:05 MST")
	}

	mErr := multierror.Append(nil,
		d.Set("name", c.Name),
		d.Set("domain", c.Domain),
		d.Set("description", c.Description),
		d.Set("type", c.Type),
		d.Set("expiration", expiration),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("error setting LB Certificate fields: %s", err)
	}
	return nil
}
