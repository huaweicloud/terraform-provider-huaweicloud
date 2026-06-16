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

// @API EIP POST /v3/global-eip/resource-instances/filter
func DataSourceGlobalEipsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEipsByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Required: true,
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
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// In the API documentation, it is of type `Object`,
						// but here it has been changed to type `string`.
						"resource_detail": {
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
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildGlobalEipsByTagsRequestPath(requestPath string, offset int) string {
	if offset == 0 {
		return requestPath
	}

	return fmt.Sprintf("%s?offset=%d", requestPath, offset)
}

func buildGlobalEipsByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	tags := d.Get("tags").([]interface{})
	res := make([]map[string]interface{}, 0, len(tags))
	for _, tagRaw := range tags {
		tag, ok := tagRaw.(map[string]interface{})
		if !ok {
			continue
		}

		res = append(res, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.ValueIgnoreEmpty(utils.PathSearch("value", tag, nil)),
		})
	}

	return map[string]interface{}{
		"tags": res,
	}
}

func dataSourceGlobalEipsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "geip"
		httpUrl = "v3/global-eip/resource-instances/filter"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildGlobalEipsByTagsBodyParams(d)),
	}

	for {
		currentRequestPath := buildGlobalEipsByTagsRequestPath(requestPath, offset)
		resp, err := client.Request("POST", currentRequestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving global EIPs by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		resourcesResp := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		if len(resourcesResp) == 0 {
			break
		}

		result = append(result, resourcesResp...)

		offset += len(resourcesResp)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", flattenGlobalEipsByTags(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalEipsByTags(resourcesResp []interface{}) []interface{} {
	if len(resourcesResp) == 0 {
		return nil
	}

	rst := make([]interface{}, len(resourcesResp))
	for i, v := range resourcesResp {
		rst[i] = map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_detail": utils.JsonToString(utils.PathSearch("resource_detail", v, nil)),
			"tags": flattenGlobalEipsByTagsTagsResp(
				utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
			"resource_name": utils.PathSearch("resource_name", v, nil),
		}
	}

	return rst
}

func flattenGlobalEipsByTagsTagsResp(tagsResp []interface{}) []interface{} {
	if len(tagsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, len(tagsResp))
	for i, tag := range tagsResp {
		rst[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		}
	}

	return rst
}
