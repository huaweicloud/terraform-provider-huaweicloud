package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentityFederationDomains
// @API IAM GET /v3/OS-FEDERATION/domains
func DataSourceIdentityFederationDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityFederationDomainsRead,

		Schema: map[string]*schema.Schema{
			"federation_token": {
				Type:     schema.TypeString,
				Required: true,
			},

			"domains": {
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
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityFederationDomainsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client := common.NewCustomClient(true, "https://iam.{region_id}.myhuaweicloud.com")
	federationDomainsPath := client.ResourceBase + "v3/OS-FEDERATION/domains"
	federationDomainsPath = strings.ReplaceAll(federationDomainsPath, "{region_id}", cfg.GetRegion(d))
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Auth-Token": d.Get("federation_token").(string),
		},
	}
	response, err := client.Request("GET", federationDomainsPath, &options)
	if err != nil {
		return diag.Errorf("error federationDomains: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generate UUID: %s", err)
	}
	d.SetId(id)
	domainsBody := utils.PathSearch("domains", respBody, make([]interface{}, 0)).([]interface{})
	domains := flattenDomains(domainsBody)
	if err = d.Set("domains", domains); err != nil {
		return diag.Errorf("error set domains filed: %s", err)
	}
	return nil
}

func flattenDomains(domainsBody []interface{}) []map[string]interface{} {
	if domainsBody == nil {
		return nil
	}
	domains := make([]map[string]interface{}, 0, len(domainsBody))
	for _, domain := range domainsBody {
		domains = append(domains, map[string]interface{}{
			"id":          utils.PathSearch("id", domain, nil),
			"name":        utils.PathSearch("name", domain, nil),
			"description": utils.PathSearch("description", domain, nil),
			"enabled":     utils.PathSearch("enabled", domain, nil),
		})
	}
	return domains
}
