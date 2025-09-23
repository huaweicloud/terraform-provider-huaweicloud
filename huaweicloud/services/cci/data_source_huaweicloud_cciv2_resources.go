package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI GET /apis/cci/v2
func DataSourceV2Resources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2ResourcesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"categories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespaced": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"short_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"singular_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_version_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"verbs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceV2ResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	listResourcesHttpUrl := "apis/cci/v2"
	listResourcesPath := client.Endpoint + listResourcesHttpUrl
	listResourcesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listResourcesResp, err := client.Request("GET", listResourcesPath, &listResourcesOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	listResourcesRespBody, err := utils.FlattenResponse(listResourcesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resources := utils.PathSearch("resources", listResourcesRespBody, make([]interface{}, 0)).([]interface{})

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", flattenResources(resources)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResources(resources []interface{}) []interface{} {
	if len(resources) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resources))
	for _, v := range resources {
		rst = append(rst, map[string]interface{}{
			"categories":           utils.PathSearch("categories", v, nil),
			"group":                utils.PathSearch("group", v, nil),
			"kind":                 utils.PathSearch("kind", v, nil),
			"name":                 utils.PathSearch("name", v, nil),
			"namespaced":           utils.PathSearch("namespaced", v, nil),
			"short_names":          utils.PathSearch("shortNames", v, nil),
			"singular_name":        utils.PathSearch("singularName", v, nil),
			"storage_version_hash": utils.PathSearch("storageVersionHash", v, nil),
			"verbs":                utils.PathSearch("verbs", v, nil),
			"version":              utils.PathSearch("version", v, nil),
		})
	}
	return rst
}
