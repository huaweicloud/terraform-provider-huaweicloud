package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v2/{project_id}/ous
func DataSourceOus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOusRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the OUs are located.`,
			},
			"ou_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the OU.`,
			},
			"ous": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the OU.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the OU.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the AD domain.`,
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The AD domain name to which the OU belongs.`,
						},
						"ou_dn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The distinguished name (DN) of the OU.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the OU.`,
						},
					},
				},
				Description: `The list of OUs that match the filter parameters.`,
			},
		},
	}
}

func buildOusQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("ou_name"); ok {
		res = fmt.Sprintf("%s&ou_name=%v", res, v)
	}
	return res
}

func dataSourceOusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	ous, err := listOus(client, buildOusQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying OUs: %s", err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("ous", flattenOus(ous)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOus(ous []interface{}) []interface{} {
	if len(ous) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(ous))
	for _, ou := range ous {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", ou, nil),
			"name":        utils.PathSearch("ou_name", ou, nil),
			"domain_id":   utils.PathSearch("domain_id", ou, nil),
			"domain":      utils.PathSearch("domain", ou, nil),
			"ou_dn":       utils.PathSearch("ou_dn", ou, nil),
			"description": utils.PathSearch("description", ou, nil),
		})
	}

	return result
}
