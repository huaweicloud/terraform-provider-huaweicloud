package ga

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

// @API GA POST /v1/{resource_type}/resource-instances/filter
func DataSourceGaResourcesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaResourceByTagsRead,

		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource type.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the list of tags.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the key of the tag.`,
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the list of the tag values.`,
						},
					},
				},
			},
			"matches": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the list of matches.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the key for matching a resource instance.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the value for matching a resource instance.`,
						},
					},
				},
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of target resources that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource.`,
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the resource.`,
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of tags associated with the resource.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The key of the tag.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The value of the tag.`,
									},
								},
							},
						},
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total count of the resources.`,
			},
		},
	}
}

func dataSourceGaResourceByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "ga"
		httpUrl    = "v1/{resource_type}/resource-instances/filter"
		result     = make([]interface{}, 0)
		limit      = 1000
		offset     = 0
		totalCount = float64(0)
		mErr       *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{resource_type}", d.Get("resource_type").(string))

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildGaResourceByTagsBodyParams(d)),
	}

	for {
		currentListPath := listPath + fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
		listResp, err := client.Request("POST", currentListPath, &reqOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalCount = utils.PathSearch("total_count", listRespBody, float64(0)).(float64)
		resources := utils.PathSearch("resources", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(resources) == 0 {
			break
		}

		result = append(result, resources...)

		offset += len(resources)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("resources", flattenGaResourceByTagsRespBody(result)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGaResourceByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"tags":    buildGaResourceByTagsTagsBodyParams(d),
		"matches": buildGaResourceByTagsMatchesBodyParams(d),
	}

	return bodyParams
}

func buildGaResourceByTagsTagsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	tags, ok := d.GetOk("tags")
	if !ok {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(tags.([]interface{})))

	for i, tag := range tags.([]interface{}) {
		bodyParams[i] = map[string]interface{}{
			"key":    utils.PathSearch("key", tag, nil),
			"values": utils.PathSearch("values", tag, nil),
		}
	}

	return bodyParams
}

func buildGaResourceByTagsMatchesBodyParams(d *schema.ResourceData) []map[string]interface{} {
	matches, ok := d.GetOk("matches")
	if !ok {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(matches.([]interface{})))

	for i, match := range matches.([]interface{}) {
		bodyParams[i] = map[string]interface{}{
			"key":    utils.PathSearch("key", match, nil),
			"values": utils.PathSearch("values", match, nil),
		}
	}

	return bodyParams
}

func flattenGaResourceByTagsRespBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	resources := make([]interface{}, len(resp))
	for i, v := range resp {
		tags := utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})
		tagsList := make([]interface{}, 0, len(tags))
		for _, tag := range tags {
			tagsList = append(tagsList, map[string]interface{}{
				"key":   utils.PathSearch("key", tag, nil),
				"value": utils.PathSearch("value", tag, nil),
			})
		}
		resources[i] = map[string]interface{}{
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"tags":          tagsList,
		}
	}
	return resources
}
