package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentityV5AgencyAttachedPolicies
// @API IAM GET /v5/agencies/{agency_id}/attached-policies
func DataSourceIdentityV5AgencyAttachedPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5AgencyAttachedPoliciesRead,

		Schema: map[string]*schema.Schema{
			"agency_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attached_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attached_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5AgencyAttachedPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	agencyId := d.Get("agency_id").(string)
	var allPolicies []interface{}
	var marker string
	var path string
	for {
		path = client.Endpoint + "v5/agencies/{agency_id}/attached-policies" + buildListAttachedPoliciesV5Params(marker)
		path = strings.ReplaceAll(path, "{agency_id}", agencyId)
		reqOpt := &golangsdk.RequestOpts{KeepResponseBody: true}
		r, err := client.Request("GET", path, reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving attached policies: %s", err)
		}
		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}
		policies := flattenListAttachedPoliciesV5Response(resp)
		allPolicies = append(allPolicies, policies...)

		marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
		if marker == "" {
			break
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	if err = d.Set("attached_policies", allPolicies); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
