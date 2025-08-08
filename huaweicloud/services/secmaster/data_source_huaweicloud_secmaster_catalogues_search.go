package secmaster

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

// @API Secmaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/catalogues/search
func DataSourceSecmasterCataloguesSearch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterCataloguesSearchRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_catalogue": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"second_catalogue": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"catalogue_status": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"layout_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"publisher_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"analysis_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_catalogue": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"second_catalogue": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"catalogue_status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"catalogue_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"layout_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"layout_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"publisher_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_card_area": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_display": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_landing_page": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_navigation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						// Fields `parent_alias_en` and `parent_alias_zh` are misspelled in the API documentation.
						// The actual names defined in the current schema are valid.
						"parent_alias_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_alias_zh": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"second_alias_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"second_alias_zh": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildSecmasterCataloguesSearchParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"parent_catalogue": utils.ValueIgnoreEmpty(d.Get("parent_catalogue")),
		"second_catalogue": utils.ValueIgnoreEmpty(d.Get("second_catalogue")),
		"catalogue_status": d.Get("catalogue_status"),
		"layout_name":      utils.ValueIgnoreEmpty(d.Get("layout_name")),
		"publisher_name":   utils.ValueIgnoreEmpty(d.Get("publisher_name")),
		"analysis_version": utils.ValueIgnoreEmpty(d.Get("analysis_version")),
	}
}

func dataSourceSecmasterCataloguesSearchRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/catalogues/search"
		offset  = 0
		allData = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildSecmasterCataloguesSearchParams(d)),
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s?offset=%d", requestPath, offset)

		resp, err := client.Request("POST", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster catalogues search: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		data := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(data) == 0 {
			break
		}

		// There is a problem with the API's paging logic. The offset can always be queried for data.
		// It is necessary to compare the `total` value to safely jump out of the loop logic.
		allData = append(allData, data...)
		totalCount := int(utils.PathSearch("total", respBody, float64(0)).(float64))
		if len(allData) >= totalCount {
			break
		}

		offset += len(data)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenCataloguesSearch(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCataloguesSearch(catalogues []interface{}) []interface{} {
	if len(catalogues) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(catalogues))
	for _, v := range catalogues {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"parent_catalogue":  utils.PathSearch("parent_catalogue", v, nil),
			"second_catalogue":  utils.PathSearch("second_catalogue", v, nil),
			"catalogue_status":  utils.PathSearch("catalogue_status", v, nil),
			"catalogue_address": utils.PathSearch("catalogue_address", v, nil),
			"layout_id":         utils.PathSearch("layout_id", v, nil),
			"layout_name":       utils.PathSearch("layout_name", v, nil),
			"publisher_name":    utils.PathSearch("publisher_name", v, nil),
			"is_card_area":      utils.PathSearch("is_card_area", v, nil),
			"is_display":        utils.PathSearch("is_display", v, nil),
			"is_landing_page":   utils.PathSearch("is_landing_page", v, nil),
			"is_navigation":     utils.PathSearch("is_navigation", v, nil),
			"parent_alias_en":   utils.PathSearch("parent_alias_en", v, nil),
			"parent_alias_zh":   utils.PathSearch("parent_alias_zh", v, nil),
			"second_alias_en":   utils.PathSearch("second_alias_en", v, nil),
			"second_alias_zh":   utils.PathSearch("second_alias_zh", v, nil),
		})
	}

	return rst
}
