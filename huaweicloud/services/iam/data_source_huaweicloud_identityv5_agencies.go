package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/agencies
// @API IAM GET /v5/agencies/{agency_id}
func DataSourceV5Agencies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5AgenciesRead,

		Schema: map[string]*schema.Schema{
			"path_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The path prefix of the agency.`,
			},
			"agency_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"path_prefix"},
				Description:   `The ID of the agency.`,
			},
			"agencies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the agencies.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agency_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the agency.`,
						},
						"agency_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the agency.`,
						},
						"trust_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The trust policy of the agency.`,
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The path of the agency.`,
						},
						"trust_domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The account ID of the trusted domain.`,
						},
						"trust_domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The account name of the trusted domain.`,
						},
						"urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URN of the agency.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the agency or trust agency.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the agency or trust agency.`,
						},
						"max_session_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum session duration of the agency.`,
						},
					},
				},
			},
		},
	}
}

func buildV5AgenciesQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("path_prefix"); ok {
		return fmt.Sprintf("&path_prefix=%v", v)
	}

	return ""
}

func listV5Agencies(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v5/agencies"
		result  = make([]interface{}, 0)
		// The default limit is 100, the maximum limit is 200.
		limit  = 200
		marker = ""
	)

	listPath := client.Endpoint + httpUrl
	listPath = fmt.Sprintf("%s?limit=%d%s", listPath, limit, buildV5AgenciesQueryParams(d))
	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &reqOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		agencies := utils.PathSearch("agencies", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, agencies...)
		if len(agencies) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func getV5AgencyById(client *golangsdk.ServiceClient, agencyId string) (interface{}, error) {
	getPath := client.Endpoint + "v5/agencies/{agency_id}"
	getPath = strings.ReplaceAll(getPath, "{agency_id}", agencyId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	r, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	resp, err := utils.FlattenResponse(r)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("agency", resp, map[string]interface{}{}), nil
}

func dataSourceV5AgenciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var agencies []interface{}
	if agencyId := d.Get("agency_id").(string); agencyId != "" {
		agency, err := getV5AgencyById(client, agencyId)
		if err != nil {
			return diag.Errorf("error retrieving agency (%s): %s", agencyId, err)
		}
		agencies = []interface{}{agency}
	} else {
		agencies, err = listV5Agencies(client, d)
		if err != nil {
			return diag.Errorf("error retrieving agencies: %s", err)
		}
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}

	d.SetId(randomId)

	return diag.FromErr(d.Set("agencies", flattenV5Agencies(agencies)))
}

func flattenV5Agencies(agencies []interface{}) []interface{} {
	if len(agencies) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(agencies))
	for _, agency := range agencies {
		result = append(result, map[string]interface{}{
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
		})
	}
	return result
}
