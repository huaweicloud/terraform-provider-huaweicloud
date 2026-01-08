package dns

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

// @API DNS GET /v2.1/system-lines
func DataSourceSystemLines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSystemLinesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the system lines are located.`,
			},
			"locale": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The display language.`,
			},
			"lines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the system line.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the system line.`,
						},
						"father_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the parent line.`,
						},
					},
				},
				Description: `The list of system lines that match the filter parameters.`,
			},
		},
	}
}

func buildListSystemLinesQueryParams(d *schema.ResourceData, limit int) string {
	queryParams := fmt.Sprintf("?limit=%d", limit)

	if v, ok := d.GetOk("locale"); ok {
		queryParams = fmt.Sprintf("%s&locale=%v", queryParams, v)
	}

	return queryParams
}

func listSystemLines(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2.1/system-lines"
		limit   = 1000
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath += buildListSystemLinesQueryParams(d, limit)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		lines := utils.PathSearch("lines", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, lines...)
		if len(lines) < limit {
			break
		}

		offset += len(lines)
	}

	return result, nil
}

func dataSourceSystemLinesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	lines, err := listSystemLines(client, d)
	if err != nil {
		return diag.Errorf("error querying system lines: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("lines", flattenSystemLines(lines)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSystemLines(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":        utils.PathSearch("id", item, nil),
			"name":      utils.PathSearch("name", item, nil),
			"father_id": utils.PathSearch("father_id", item, nil),
		})
	}

	return result
}
