package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/identity/v3/agency"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IAM GET /v3.0/OS-AGENCY/agencies
func DataSourceIdentityAgencies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityAgenciesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trust_domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agencies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trust_domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trust_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityAgenciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	iamClient, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	domainID := cfg.DomainID
	if domainID == "" {
		return diag.Errorf("the domain_id must be specified in the provider configuration")
	}
	opts := agency.ListOpts{
		DomainID:      domainID,
		Name:          d.Get("name").(string),
		TrustDomainId: d.Get("trust_domain_id").(string),
	}
	pages, err := agency.List(iamClient, opts).AllPages()
	if err != nil {
		return diag.Errorf("error querying IAM agencies: %s", err)
	}
	resp, err := agency.ExtractList(pages)
	if err != nil {
		return diag.Errorf("error querying IAM agencies: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)
	mErr := multierror.Append(nil,
		d.Set("agencies", flattenAgencies(resp)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM agencies fields: %s", err)
	}

	return nil
}

func flattenAgencies(agencies []agency.Agency) []map[string]interface{} {
	if len(agencies) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(agencies))
	for i, val := range agencies {
		result[i] = map[string]interface{}{
			"id":                val.ID,
			"name":              val.Name,
			"trust_domain_id":   val.DelegatedDomainID,
			"trust_domain_name": val.DelegatedDomainName,
			"created_at":        val.CreateTime,
			"expired_at":        val.ExpireTime,
			"description":       val.Description,
			"duration":          val.Duration,
		}
	}
	return result
}
