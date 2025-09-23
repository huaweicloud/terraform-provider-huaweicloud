package rms

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

// @API CONFIG POST /v1/resource-manager/{resource_type}/resource-instances/filter
func DataSourceResourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource type.`,
			},
			"without_any_tag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies if the resource has no tags.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tags.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the tag key.`,
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the tag values.`,
						},
					},
				},
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The resource list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The tags.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The tag key.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The tag value.`,
									},
								},
							},
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource name.`,
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource ID.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceResourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rms", region)

	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	instances, err := getResourceInstances(client, d)
	if err != nil {
		return diag.Errorf("error retrieving RMS Resource Instances: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("resources", flattenInstances(instances)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getResourceInstances(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v1/resource-manager/{resource_type}/resource-instances/filter"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{resource_type}", d.Get("resource_type").(string))
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildResourceInstancesBodyParams(d)),
	}
	rst := make([]interface{}, 0)

	offset := 0
	for {
		path := fmt.Sprintf("%s?limit=100&offset=%d", path, offset)
		resp, err := client.Request("POST", path, &opt)

		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		curArray := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		rst = append(rst, curArray...)

		offset += 100
		total := utils.PathSearch("total_count", respBody, float64(0))
		if int(total.(float64)) <= offset {
			break
		}
	}
	return rst, nil
}

func buildResourceInstancesBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"tags": buildTagsBodyParams(d.Get("tags").([]interface{})),
	}

	if withoutAnyTag := d.Get("without_any_tag").(bool); withoutAnyTag {
		body["without_any_tag"] = withoutAnyTag
	}

	return body
}

func buildTagsBodyParams(params []interface{}) []interface{} {
	if len(params) == 0 {
		return nil
	}

	tags := make([]interface{}, 0, len(params))
	for _, param := range params {
		paramMap := param.(map[string]interface{})
		tags = append(tags, map[string]interface{}{
			"key":    paramMap["key"],
			"values": utils.ExpandToStringList(paramMap["values"].([]interface{})),
		})
	}
	return tags
}

func flattenInstances(instances []interface{}) []interface{} {
	result := make([]interface{}, len(instances))
	for i, instance := range instances {
		tags := utils.PathSearch("tags", instance, make([]interface{}, 0)).([]interface{})
		tagsList := make([]interface{}, 0, len(tags))
		for _, tag := range tags {
			tagsList = append(tagsList, map[string]interface{}{
				"key":   utils.PathSearch("key", tag, nil),
				"value": utils.PathSearch("value", tag, nil),
			})
		}
		result[i] = map[string]interface{}{
			"tags":          tagsList,
			"resource_name": utils.PathSearch("resource_name", instance, nil),
			"resource_id":   utils.PathSearch("resource_id", instance, nil),
		}
	}
	return result
}
