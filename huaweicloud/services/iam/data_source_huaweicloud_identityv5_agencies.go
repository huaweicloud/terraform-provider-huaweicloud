package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentityV5Agencies
// @API IAM GET /v5/agencies
func DataSourceIdentityV5Agencies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5AgenciesRead,

		Schema: map[string]*schema.Schema{
			"path_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agency_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"path_prefix"},
			},

			"agencies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of agencies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trust_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agency_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agency_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
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
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_session_duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5AgenciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var allAgencies []interface{}
	var marker string
	reqOpt := &golangsdk.RequestOpts{KeepResponseBody: true}
	if agencyId := d.Get("agency_id").(string); agencyId != "" {
		path := client.Endpoint + "v5/agencies/{agency_id}"
		path = strings.ReplaceAll(path, "{agency_id}", agencyId)
		r, err := client.Request("GET", path, reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving agencies: %s", err)
		}
		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}
		agency := utils.PathSearch("agency", resp, map[string]interface{}{}).(map[string]interface{})
		agencies := map[string]interface{}{
			"agencies": []interface{}{agency},
		}
		allAgencies = append(allAgencies, flattenListAgenciesV5Response(agencies)...)
	} else {
		for {
			path := client.Endpoint + "v5/agencies" + buildListAgenciesV5Params(d, marker)
			r, err := client.Request("GET", path, reqOpt)
			if err != nil {
				return diag.Errorf("error retrieving agencies: %s", err)
			}
			resp, err := utils.FlattenResponse(r)
			if err != nil {
				return diag.FromErr(err)
			}
			agencies := flattenListAgenciesV5Response(resp)
			allAgencies = append(allAgencies, agencies...)

			marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
			if marker == "" {
				break
			}
		}
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	mErr := multierror.Append(nil,
		d.Set("agencies", allAgencies),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListAgenciesV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	agencies := utils.PathSearch("agencies", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(agencies))
	for i, agency := range agencies {
		result[i] = map[string]interface{}{
			"trust_policy":         utils.PathSearch("trust_policy", agency, nil),
			"agency_id":            utils.PathSearch("agency_id", agency, nil),
			"agency_name":          utils.PathSearch("agency_name", agency, nil),
			"path":                 utils.PathSearch("path", agency, nil),
			"trust_domain_id":      utils.PathSearch("trust_domain_id", agency, nil),
			"trust_domain_name":    utils.PathSearch("trust_domain_name", agency, nil),
			"urn":                  utils.PathSearch("urn", agency, nil),
			"created_at":           utils.PathSearch("created_at", agency, nil),
			"description":          utils.PathSearch("description", agency, nil),
			"max_session_duration": utils.PathSearch("max_session_duration", agency, nil),
		}
	}
	return result
}

func buildListAgenciesV5Params(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if v, ok := d.GetOk("path_prefix"); ok {
		res = fmt.Sprintf("%s&path_prefix=%v", res, v)
	}
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}
