package ga

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA GET /v1/pops
func DataSourceGaPops() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaPopsRead,

		Schema: map[string]*schema.Schema{
			"pops": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The access point ID.`,
						},
					},
				},
			},
		},
	}
}

func listpops(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/pops"
		result  = make([]interface{}, 0)
		limit   = 500
		marker  = ""
	)

	listPath := client.Endpoint + httpUrl
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	reqOpt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, reqOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		pops := utils.PathSearch("pops", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, pops...)
		if len(pops) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func dataSourceGaPopsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("ga", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	pops, err := listpops(client)
	if err != nil {
		return diag.Errorf("error listing GA Pops: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}

	d.SetId(generateUUID)

	return diag.FromErr(d.Set("pops", flattenPops(pops)))
}

func flattenPops(pops []interface{}) []interface{} {
	if len(pops) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(pops))
	for _, pop := range pops {
		result = append(result, map[string]interface{}{
			"id": utils.PathSearch("id", pop, nil),
		})
	}

	return result
}
