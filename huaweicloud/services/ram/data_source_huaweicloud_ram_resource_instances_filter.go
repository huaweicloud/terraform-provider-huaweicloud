package ram

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

// @API RAM POST /v1/resource-shares/resource-instances/filter
func DataSourceResourceInstancesFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceInstancesFilterRead,
		Schema: map[string]*schema.Schema{
			"without_any_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceInstancesFilterSchema(),
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceInstancesFilterSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resource_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func buildResourceInstancesFilterQueryParams(offset int) string {
	// The limit default value is `1000`.
	return fmt.Sprintf("?limit=1000&offset=%d", offset)
}

func buildResourceInstancesFilterBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"without_any_tag": d.Get("without_any_tag"),
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsInput := v.([]interface{})
		tags := make([]map[string]interface{}, 0, len(tagsInput))
		for _, item := range tagsInput {
			tag, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			m := map[string]interface{}{
				"key": tag["key"],
			}
			if v, ok := tag["values"]; ok && v != nil {
				m["values"] = utils.ExpandToStringList(v.([]interface{}))
			}

			tags = append(tags, m)
		}

		params["tags"] = tags
	}

	if v, ok := d.GetOk("matches"); ok {
		matchesInput := v.([]interface{})
		matches := make([]map[string]interface{}, 0, len(matchesInput))
		for _, item := range matchesInput {
			match, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			matches = append(matches, map[string]interface{}{
				"key":   match["key"],
				"value": match["value"],
			})
		}

		params["matches"] = matches
	}

	return params
}

func dataSourceResourceInstancesFilterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/resource-shares/resource-instances/filter"
		product    = "ram"
		offset     = 0
		result     = make([]interface{}, 0)
		totalCount float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildResourceInstancesFilterBodyParams(d)),
	}

	for {
		requestPathWithOffset := requestPath + buildResourceInstancesFilterQueryParams(offset)
		resp, err := client.Request("POST", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving RAM resource instances filter: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalCount = utils.PathSearch("total_count", respBody, float64(0)).(float64)
		resourcesResp := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		if len(resourcesResp) == 0 {
			break
		}

		result = append(result, resourcesResp...)
		offset += len(resourcesResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("resources", flattenResourceInstancesFilter(result)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResourceInstancesFilter(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		resourceMap := map[string]interface{}{
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"tags": flattenResourcesTags(
				utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		}

		if detail := utils.PathSearch("resource_detail", v, nil); detail != nil {
			resourceMap["resource_detail"] = fmt.Sprintf("%v", detail)
		}

		rst = append(rst, resourceMap)
	}

	return rst
}

func flattenResourcesTags(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rstMap := map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
		rst = append(rst, rstMap)
	}

	return rst
}
