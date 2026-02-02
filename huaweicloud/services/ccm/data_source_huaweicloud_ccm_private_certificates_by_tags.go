package ccm

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

// @API CCM POST /v1/private-certificates/resource-instances/filter
func DataSourcePrivateCertificatesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourcePrivateCertificatesByTagsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     privateCertificatesByTagsTagsSchema(),
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     privateCertificatesByTagsMatchesSchema(),
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     privateCertificatesByTagsResourcesSchema(),
			},
		},
	}
}

func privateCertificatesByTagsMatchesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func privateCertificatesByTagsTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"values": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func privateCertificatesByTagsResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     privateCertificatesByTagsResourcesTagsSchema(),
			},
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// The API documentation does not provide the data structure for this field,
			// so it is temporarily treated as a JSON string.
			"resource_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func privateCertificatesByTagsResourcesTagsSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func buildPrivateCertificatesByTagsTagsBodyParams(rawArray []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, raw := range rawArray {
		rawMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":    rawMap["key"],
			"values": rawMap["values"],
		})
	}

	return rst
}

func buildPrivateCertificatesByTagsMatchesBodyParams(rawArray []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, raw := range rawArray {
		rawMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":   rawMap["key"],
			"value": rawMap["value"],
		})
	}

	return rst
}

func buildPrivateCertificatesByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"tags":    buildPrivateCertificatesByTagsTagsBodyParams(d.Get("tags").([]interface{})),
		"matches": buildPrivateCertificatesByTagsMatchesBodyParams(d.Get("matches").([]interface{})),
		"limit":   50,
	}
}

func datasourcePrivateCertificatesByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		httpUrl = "v1/private-certificates/resource-instances/filter"
		product = "ccm"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	// The pagination parameter for this API is invalid. Specifying an offset value will retrieve duplicate data,
	// so pagination is not currently supported.
	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildPrivateCertificatesByTagsBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CCM private certificates by tags: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("total_count", int(utils.PathSearch("total_count", respBody, float64(0)).(float64))),
		d.Set("resources", flattenPrivateCertificatesByTagsResources(resources)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPrivateCertificatesByTagsResources(respArray []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		tagsResp := utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})
		resourceDetailResp := utils.JsonToString(utils.PathSearch("resource_detail", v, nil))

		rst = append(rst, map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"tags":            flattenPrivateCertificatesByTagsTagsResources(tagsResp),
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_detail": resourceDetailResp,
		})
	}

	return rst
}

func flattenPrivateCertificatesByTagsTagsResources(respArray []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return rst
}
