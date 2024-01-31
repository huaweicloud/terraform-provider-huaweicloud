package sdrs

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/sdrs/v1/domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API SDRS GET /v1/{project_id}/active-domains
func DataSourceSDRSDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSDRSDomainRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSDRSDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	v, err := domains.Get(client).Extract()
	if err != nil {
		return diag.Errorf("error retrieving SDRS domains, %s", err)
	}

	domain, err := flattenSDRSDomain(d, v.Domains)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(domain.Id)
	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", domain.Name),
		d.Set("description", domain.Description),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

// flattenSDRSDomain Currently, there is only one element in the response domains.
func flattenSDRSDomain(d *schema.ResourceData, domainList []domains.Domain) (*domains.Domain, error) {
	name := d.Get("name").(string)
	filterDomains := make([]domains.Domain, 0, len(domainList))
	for _, domain := range domainList {
		if name != "" && name != domain.Name {
			continue
		}
		filterDomains = append(filterDomains, domain)
	}
	if len(filterDomains) < 1 {
		return nil, fmt.Errorf("your query returned no results. " +
			"Please change your search criteria and try again")
	}
	return &filterDomains[0], nil
}
