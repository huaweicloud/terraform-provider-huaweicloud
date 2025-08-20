package eg

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

// @API EG GET /v1/{project_id}/target-catalogs
func DataSourceEventTargetCatalogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventTargetCatalogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the event target catalogs are located.`,
			},
			"fuzzy_label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The label of the event target catalog to be queried.`,
			},
			"support_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The support type list of event targets to be queried.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sort order for querying event target catalogs.`,
			},
			"catalogs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the event target catalog.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the event target catalog.`,
						},
						"label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The display name of the event target catalog.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the event target catalog.`,
						},
						"support_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The support type list of event target catalog.`,
						},
						"provider_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The provider type of the event target catalog.`,
						},
						"parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the target parameter.`,
									},
									"label": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The display name of the target parameter.`,
									},
									"metadata": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The metadata of the target parameter, in JSON format.`,
									},
								},
							},
							Description: `The parameters of the event target catalog.`,
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the event target catalog, in UTC format.`,
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the event target catalog, in UTC format.`,
						},
					},
				},
				Description: `All event target catalogs that match the filter parameters.`,
			},
		},
	}
}

func buildEventTargetCatalogsQueryParams(d *schema.ResourceData) string {
	res := ""

	if fuzzyLabel, ok := d.GetOk("fuzzy_label"); ok {
		res = fmt.Sprintf("%s&fuzzy_label=%v", res, fuzzyLabel)
	}

	if supportTypes, ok := d.GetOk("support_types"); ok {
		for _, supportType := range supportTypes.([]interface{}) {
			res = fmt.Sprintf("%s&support_types=%v", res, supportType)
		}
	}

	if sort, ok := d.GetOk("sort"); ok {
		res = fmt.Sprintf("%s&sort=%v", res, sort)
	}

	return res
}

func queryEventTargetCatalogs(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/target-catalogs"
		offset  = 0
		limit   = 500
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += fmt.Sprintf("?limit=%d", limit)
	listPath += buildEventTargetCatalogsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		eventTargetCatalogs := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, eventTargetCatalogs...)
		if len(eventTargetCatalogs) < limit {
			break
		}

		offset += len(eventTargetCatalogs)
	}

	return result, nil
}

func flattenEventTargetCatalogParameters(parameters []interface{}) []map[string]interface{} {
	if len(parameters) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(parameters))
	for _, parameter := range parameters {
		result = append(result, map[string]interface{}{
			"name":     utils.PathSearch("name", parameter, nil),
			"label":    utils.PathSearch("label", parameter, nil),
			"metadata": utils.JsonToString(utils.PathSearch("metadata", parameter, nil)),
		})
	}

	return result
}

func flattenDataEventTargetCatalogs(catalogs []interface{}) []map[string]interface{} {
	if len(catalogs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(catalogs))
	for _, catalog := range catalogs {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", catalog, nil),
			"name":          utils.PathSearch("name", catalog, nil),
			"label":         utils.PathSearch("label", catalog, nil),
			"description":   utils.PathSearch("description", catalog, nil),
			"support_types": utils.PathSearch("support_types", catalog, nil),
			"provider_type": utils.PathSearch("provider_type", catalog, nil),
			"parameters": flattenEventTargetCatalogParameters(utils.PathSearch("parameters",
				catalog, make([]interface{}, 0)).([]interface{})),
			"created_time": utils.PathSearch("created_time", catalog, nil),
			"updated_time": utils.PathSearch("updated_time", catalog, nil),
		})
	}

	return result
}

func dataSourceEventTargetCatalogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	catalogs, err := queryEventTargetCatalogs(client, d)
	if err != nil {
		return diag.Errorf("error querying event target catalogs: %s", err)
	}

	randomID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("catalogs", flattenDataEventTargetCatalogs(catalogs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
