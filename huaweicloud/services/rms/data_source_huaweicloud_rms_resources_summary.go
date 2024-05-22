package rms

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/all-resources/summary
// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/tracked-resources/summary

func DataSourceResourcesSummary() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourcesSummaryRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tracked": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_deleted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"resources_summary": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"regions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"count": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceResourcesSummaryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rms", region)
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}

	getQueryParams := buildResourcesSummaryQueryParams(d)

	getResourcesSummaryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var resp *http.Response
	if d.Get("tracked").(bool) || d.Get("resource_deleted").(bool) {
		getTrackedResourcesSummaryHttpUrl := "v1/resource-manager/domains/{domain_id}/tracked-resources/summary"
		getTrackedResourcesSummaryPath := client.Endpoint + getTrackedResourcesSummaryHttpUrl
		getTrackedResourcesSummaryPath = strings.ReplaceAll(getTrackedResourcesSummaryPath, "{domain_id}", cfg.DomainID)
		getTrackedResourcesSummaryPath += getQueryParams

		resp, err = client.Request("GET", getTrackedResourcesSummaryPath, &getResourcesSummaryOpt)
		if err != nil {
			return diag.Errorf("error retrieving RMS tracked resources summary: %s", err)
		}
	} else {
		getResourcesSummaryHttpUrl := "v1/resource-manager/domains/{domain_id}/all-resources/summary"
		getResourcesSummaryPath := client.Endpoint + getResourcesSummaryHttpUrl
		getResourcesSummaryPath = strings.ReplaceAll(getResourcesSummaryPath, "{domain_id}", cfg.DomainID)
		getResourcesSummaryPath += getQueryParams

		resp, err = client.Request("GET", getResourcesSummaryPath, &getResourcesSummaryOpt)
		if err != nil {
			return diag.Errorf("error retrieving RMS resources summary: %s", err)
		}
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		nil,
		d.Set("resources_summary", flattenResourcesSummary(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildResourcesSummaryQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("region_id"); ok {
		res = fmt.Sprintf("%s&region_id=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&ep_id=%v", res, v)
	}
	if v, ok := d.GetOk("project_id"); ok {
		res = fmt.Sprintf("%s&project_id=%v", res, v)
	}
	if v, ok := d.GetOk("resource_deleted"); ok {
		res = fmt.Sprintf("%s&resource_deleted=%v", res, v)
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := v.(map[string]interface{})
		for k, val := range tagsMap {
			tagsString := fmt.Sprintf(`%s%%3D%s`, k, val)
			res = fmt.Sprintf("%s&tags=%v", res, tagsString)
		}
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func flattenResourcesSummary(resourcesSummaryRaw interface{}) []map[string]interface{} {
	if resourcesSummaryRaw == nil {
		return nil
	}

	resourcesSummary := resourcesSummaryRaw.([]interface{})
	res := make([]map[string]interface{}, len(resourcesSummary))
	for i, v := range resourcesSummary {
		summary := v.(map[string]interface{})
		res[i] = map[string]interface{}{
			"service": summary["provider"],
			"types":   flattenTypes(summary["types"]),
		}
	}

	return res
}

func flattenTypes(typesRaw interface{}) []map[string]interface{} {
	if typesRaw == nil {
		return nil
	}

	types := typesRaw.([]interface{})
	res := make([]map[string]interface{}, len(types))
	for i, v := range types {
		t := v.(map[string]interface{})
		res[i] = map[string]interface{}{
			"type":    t["type"],
			"regions": flattenRegions(t["regions"]),
		}
	}

	return res
}

func flattenRegions(regionsRaw interface{}) []map[string]interface{} {
	if regionsRaw == nil {
		return nil
	}

	regions := regionsRaw.([]interface{})
	res := make([]map[string]interface{}, len(regions))
	for i, v := range regions {
		region := v.(map[string]interface{})
		res[i] = map[string]interface{}{
			"region_id": region["region_id"],
			"count":     region["count"],
		}
	}

	return res
}
