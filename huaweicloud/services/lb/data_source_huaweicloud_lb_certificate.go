package lb

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/elb/v2/certificates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
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
	cfg := meta.(*config.Config)

	// LoadBalancerClient catalog info: Name is "elb" and Version is "v2"
	client, err := cfg.LoadBalancerClient(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating ELB Client: %s", err)
	}

	listOpts := certificates.ListOpts{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	r, err := certificates.List(client, listOpts)
	if err != nil {
		return fmt.Errorf("error listing certificates: %s", err)
	}
	certRst, err := r.Extract()
	if err != nil {
		return fmt.Errorf("unable to retrieve ELB certificates: %s", err)
	}
	log.Printf("[DEBUG] Get certificate list: %#v", certRst)

	if len(certRst.Certificates) < 1 {
		return errors.New("your query returned no results, please change your search criteria and try again")
	}

	return setCertificateAttributes(d, certRst.Certificates[0])
}

func setCertificateAttributes(d *schema.ResourceData, c certificates.Certificate) error {
	d.SetId(c.Id)

	var expiration string
	tm, err := time.Parse("2006-01-02 15:04:05", c.ExpireTime)
	if err != nil {
		// If the format of ExpireTime is not expected, set the original value directly.
		expiration = c.ExpireTime
		log.Printf("[WAIN] The format of the ExpireTime field of the LB certificate is not expected: %s",
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
		return fmt.Errorf("error setting LB Certificate fields: %s", err)
	}
	return nil
}
