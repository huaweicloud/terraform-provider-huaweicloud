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

// @API RAM GET /v1/resource-types
func DataSourceResourceTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceTypesRead,
		Schema: map[string]*schema.Schema{
			"resource_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ResourceTypesSchema(),
			},
		},
	}
}

func ResourceTypesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceResourceTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listResourceTypesHttpUrl := "v1/resource-types"
	listResourceTypesProduct := "ram"
	listResourceTypesClient, err := cfg.NewServiceClient(listResourceTypesProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	listResourceTypesPath := listResourceTypesClient.Endpoint + listResourceTypesHttpUrl

	listResourceTypesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var resourceTypes []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listResourceTypesPath + buildListResourceTypesQueryParams(marker)
		listResourceTypesResp, err := listResourceTypesClient.Request("GET", queryPath, &listResourceTypesOpt)
		if err != nil {
			return diag.Errorf("error retrieving RAM resource types: %s", err)
		}

		listResourceTypesRespBody, err := utils.FlattenResponse(listResourceTypesResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageResourceTypes := flattenResourceTypesResp(listResourceTypesRespBody)
		resourceTypes = append(resourceTypes, onePageResourceTypes...)
		marker = utils.PathSearch("page_info.next_marker", listResourceTypesRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("resource_types", resourceTypes),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListResourceTypesQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func flattenResourceTypesResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("resource_types", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"region_id":     utils.PathSearch("region_id", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
		})
	}
	return rst
}
