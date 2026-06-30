package eip

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP POST /v3/internet-bandwidth/resource-instances/filter
func DataSourceGlobalInternetBandwidthsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalInternetBandwidthsByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"resources": resourcesSchema(),
		},
	}
}

func resourcesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"resource_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"resource_detail": {
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
			},
		},
	}
}
func buildListInternetBandwidthsByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray, ok := d.Get("tags").([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))

	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		tag := map[string]interface{}{
			"key":   rawMap["key"],
			"value": utils.ValueIgnoreEmpty(rawMap["value"]),
		}

		rst = append(rst, tag)
	}

	return map[string]interface{}{
		"tags": rst,
	}
}

func dataSourceGlobalInternetBandwidthsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v3/internet-bandwidth/resource-instances/filter"
		allResources = make([]interface{}, 0)
		offset       = 0
	)

	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildListInternetBandwidthsByTagsBodyParams(d)),
	}

	for {
		requestPathWithOffset := buildRequestPathWithOffset(requestPath, offset)

		resp, err := client.Request("POST", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving global internet bandwidths by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		if len(resources) == 0 {
			break
		}

		allResources = append(allResources, resources...)
		offset += len(resources)
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid.String())
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("resources", flattenInternetBandwidths(allResources)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInternetBandwidths(resources []interface{}) []interface{} {
	if len(resources) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resources))
	for _, v := range resources {
		rst = append(rst, map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_detail": utils.PathSearch("resource_detail", v, nil),
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"tags":            flattenTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}

func flattenTags(tags []interface{}) []interface{} {
	if len(tags) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(tags))
	for _, v := range tags {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}

func buildRequestPathWithOffset(requestPath string, offset int) string {
	if offset == 0 {
		return requestPath
	}
	return fmt.Sprintf("%s?offset=%d", requestPath, offset)
}
