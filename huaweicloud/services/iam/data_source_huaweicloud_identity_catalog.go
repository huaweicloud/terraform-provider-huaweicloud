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

// @API IAM GET /v3/auth/catalog
func DataSourceIdentityCatalog() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCatalogRead,

		Schema: map[string]*schema.Schema{
			"project_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"catalog": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"interface": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityCatalogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient := common.NewCustomClient(true, "https://iam.{region_id}.myhuaweicloud.com")
	getIdentityCatalogBasePath := iamClient.ResourceBase + "v3/auth/catalog"
	getIdentityCatalogBasePath = strings.ReplaceAll(getIdentityCatalogBasePath, "{region_id}", cfg.GetRegion(d))
	getIdentityCatalogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Auth-Token": d.Get("project_token").(string),
		},
	}
	response, err := iamClient.Request("GET", getIdentityCatalogBasePath, &getIdentityCatalogOpt)
	if err != nil {
		return diag.Errorf("error getting identity catalog: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	catalog := flattenCatalogList(utils.PathSearch("catalog", respBody, make([]interface{}, 0)).([]interface{}))
	if e := d.Set("catalog", catalog); e != nil {
		return diag.FromErr(e)
	}
	return nil
}

func flattenCatalogList(catalogBody []interface{}) []map[string]interface{} {
	catalog := make([]map[string]interface{}, len(catalogBody))
	for i, c := range catalogBody {
		catalog[i] = map[string]interface{}{
			"endpoints": flattenEndpointsLists(utils.PathSearch("endpoints", c, make([]interface{}, 0)).([]interface{})),
			"name":      utils.PathSearch("name", c, ""),
			"id":        utils.PathSearch("id", c, ""),
			"type":      utils.PathSearch("type", c, ""),
		}
	}
	return catalog
}

func flattenEndpointsLists(endpointsBody []interface{}) []map[string]interface{} {
	endpoints := make([]map[string]interface{}, len(endpointsBody))
	for i, endpoint := range endpointsBody {
		endpoints[i] = map[string]interface{}{
			"region":    utils.PathSearch("region", endpoint, ""),
			"interface": utils.PathSearch("interface", endpoint, ""),
			"id":        utils.PathSearch("id", endpoint, ""),
			"region_id": utils.PathSearch("region_id", endpoint, ""),
			"url":       utils.PathSearch("url", endpoint, ""),
		}
	}
	return endpoints
}
