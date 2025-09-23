package identitycenter

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter GET /v1/catalog/applications
func DataSourceIdentityCenterCatalogApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCenterCatalogApplicationsRead,

		Schema: map[string]*schema.Schema{
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"display_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"icon": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"application_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildListCatalogApplicationsParams(marker string) string {
	res := "?limit=100"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func listCatalogApplications(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/catalog/applications"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl

	queryParams := buildListCatalogApplicationsParams("")

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestResp, err := client.Request("GET", listPath+queryParams, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Identity Center applications in catalog: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		applications := utils.PathSearch("applications", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, applications...)

		marker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}

		queryParams = buildListCatalogApplicationsParams(marker)
	}
	return result, nil
}

func dataSourceIdentityCenterCatalogApplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	applications, err := listCatalogApplications(client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("applications", flattenCatalogApplications(applications)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCatalogApplications(applications []interface{}) []interface{} {
	if len(applications) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(applications))
	for _, application := range applications {
		result = append(result, map[string]interface{}{
			"application_id":   utils.PathSearch("application_id", application, nil),
			"application_type": utils.PathSearch("application_type", application, nil),
			"display":          flattenDisplay(utils.PathSearch("display", application, nil)),
		})
	}
	return result
}

func flattenDisplay(display interface{}) []map[string]interface{} {
	if display == nil || len(display.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"display_name": utils.PathSearch("display_name", display, nil),
			"description":  utils.PathSearch("description", display, nil),
			"icon":         utils.PathSearch("icon", display, nil),
		},
	}
}
