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

// DataSourceIdentityRegions
// @API IAM GET /v3/regions
// @API IAM GET /v3/regions/{region_id}
func DataSourceIdentityRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityRegionsRead,

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The id of the region to be queried.`,
			},

			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the region.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the region.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the region.`,
						},
						"link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource link of the region.`,
						},
						"locales": {
							Type:     schema.TypeMap,
							Computed: true,
							Description: `The map of localized region names, where the key is the language code and
                                the value is the region name in that language.`,
							Elem: &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityRegionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	regionId := d.Get("region_id").(string)
	if regionId == "" {
		return listRegions(iamClient, d)
	}
	return showRegion(iamClient, regionId, d)
}

func listRegions(iamClient *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	listRegionsPath := iamClient.Endpoint + "v3/regions"
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	response, err := iamClient.Request("GET", listRegionsPath, &options)
	if err != nil {
		return diag.Errorf("error listRegions: %s", err)
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	regionsBody := utils.PathSearch("regions", respBody, make([]interface{}, 0)).([]interface{})
	regions := make([]interface{}, 0, len(regionsBody))
	for _, region := range regionsBody {
		regions = append(regions, flattenRegion(region))
	}

	if err = d.Set("regions", regions); err != nil {
		return diag.Errorf("error setting regions fields: %s", err)
	}
	return nil
}

func showRegion(iamClient *golangsdk.ServiceClient, regionId string, d *schema.ResourceData) diag.Diagnostics {
	showRegionPath := iamClient.Endpoint + "v3/regions/{region_id}"
	showRegionPath = strings.ReplaceAll(showRegionPath, "{region_id}", regionId)
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	response, err := iamClient.Request("GET", showRegionPath, &options)
	if err != nil {
		return diag.Errorf("error showRegion: %s", err)
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	regionBody := utils.PathSearch("region", respBody, make([]interface{}, 0))
	regions := append(make([]interface{}, 0, 1), flattenRegion(regionBody))
	if err = d.Set("regions", regions); err != nil {
		return diag.Errorf("error setting regions fields: %s", err)
	}
	return nil
}

func flattenRegion(regionModel interface{}) map[string]interface{} {
	if regionModel == nil {
		return nil
	}

	region := map[string]interface{}{
		"id":          utils.PathSearch("id", regionModel, nil),
		"type":        utils.PathSearch("type", regionModel, nil),
		"description": utils.PathSearch("description", regionModel, nil),
		"link":        utils.PathSearch("links.self", regionModel, nil),
		"locales": map[string]string{
			"zh-cn": utils.PathSearch("locales.\"zh-cn\"", regionModel, "").(string),
			"en-us": utils.PathSearch("locales.\"en-us\"", regionModel, "").(string),
			"pt-br": utils.PathSearch("locales.\"pt-br\"", regionModel, "").(string),
			"es-us": utils.PathSearch("locales.\"es-us\"", regionModel, "").(string),
			"es-es": utils.PathSearch("locales.\"es-es\"", regionModel, "").(string),
		},
	}

	return region
}
