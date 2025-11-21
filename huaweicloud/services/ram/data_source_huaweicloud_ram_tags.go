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

// @API RAM GET /v1/resource-shares/tags
func DataSourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTagsRead,
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     tagsSchema(),
			},
		},
	}
}

func tagsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

	return &sc
}

func dataSourceTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listTagsHttpUrl := "v1/resource-shares/tags"
	listTagsProduct := "ram"
	listTagsClient, err := cfg.NewServiceClient(listTagsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	listTagsPath := listTagsClient.Endpoint + listTagsHttpUrl

	listTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var tags []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listTagsPath + buildListTagsQueryParams(marker)
		listTagsResp, err := listTagsClient.Request("GET", queryPath, &listTagsOpt)
		if err != nil {
			return diag.Errorf("error retrieving RAM tags: %s", err)
		}

		listTagsRespBody, err := utils.FlattenResponse(listTagsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageTags := FlattenTagsResp(listTagsRespBody)
		tags = append(tags, onePageTags...)
		marker = utils.PathSearch("page_info.next_marker", listTagsRespBody, "").(string)
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
		d.Set("tags", tags),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListTagsQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func FlattenTagsResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("tags", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":    utils.PathSearch("key", v, nil),
			"values": utils.PathSearch("values", v, nil),
		})
	}
	return rst
}
