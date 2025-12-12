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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func dataSourceResourceInstancesFilterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listResourceInstancesFilterHttpUrl := "v1/resource-shares/resource-instances/filter"
	listResourceInstancesFilterProduct := "ram"
	listResourceInstancesFilterClient, err := cfg.NewServiceClient(listResourceInstancesFilterProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	listResourceInstancesFilterPath := listResourceInstancesFilterClient.Endpoint + listResourceInstancesFilterHttpUrl

	listResourceInstancesFilterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildResourceInstancesFilterBody(d),
	}

	var resourceInstancesFilter []interface{}
	var queryPath string
	var totalCount float64

	limit := 200
	offset := 0

	for {
		queryPath = listResourceInstancesFilterPath + buildResourceInstancesFilterQueryParams(limit, offset)
		listResourceInstancesFilterResp, err := listResourceInstancesFilterClient.Request(
			"POST",
			queryPath,
			&listResourceInstancesFilterOpt,
		)
		if err != nil {
			return diag.Errorf("error retrieving RAM resource instance filter: %s", err)
		}

		listResourceInstancesFilterRespBody, err := utils.FlattenResponse(listResourceInstancesFilterResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageResourceInstancesFilter := flattenResourceInstancesFilterResp(listResourceInstancesFilterRespBody)
		resourceInstancesFilter = append(resourceInstancesFilter, onePageResourceInstancesFilter...)

		offset += limit

		totalCount = utils.PathSearch(
			"total_count",
			listResourceInstancesFilterRespBody,
			0,
		).(float64)

		if len(onePageResourceInstancesFilter) == 0 {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("resources", resourceInstancesFilter),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildResourceInstancesFilterBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{}

	if v, ok := d.GetOk("without_any_tag"); ok {
		params["without_any_tag"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		params["tags"] = v
	}

	if v, ok := d.GetOk("matches"); ok {
		params["matches"] = v
	}

	return params
}

func buildResourceInstancesFilterQueryParams(limit int, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func flattenResourceInstancesFilterResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("resources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rstMap := map[string]interface{}{
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"tags":          utils.PathSearch("tags", v, nil),
		}
		if detail := utils.PathSearch("resource_detail", v, nil); detail != nil {
			rstMap["resource_detail"] = fmt.Sprintf("%v", detail)
		}
		rst = append(rst, rstMap)
	}
	return rst
}
