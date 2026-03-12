package cfw

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

// @API CFW GET /v1/{project_id}/ptf/ip-blacklist
func DataSourceIpBlacklist() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIpBlacklistRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"records": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"effect_scope": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"import_status": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"import_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"error_message": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildIpBlacklistQueryParams(limit int, d *schema.ResourceData) string {
	return fmt.Sprintf("?fw_instance_id=%s&limit=%d", d.Get("fw_instance_id").(string), limit)
}

func dataSourceIpBlacklistRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v1/{project_id}/ptf/ip-blacklist"
		recordsResult = make([]interface{}, 0)
		// The maximum limit is `1024`
		limit  = 1024
		offset = 0
		total  float64
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildIpBlacklistQueryParams(limit, d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		requestResp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving CFW IP blacklist: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		total = utils.PathSearch("data.total", respBody, float64(0)).(float64)
		records := utils.PathSearch("data.records", respBody, make([]interface{}, 0)).([]interface{})
		recordsResult = append(recordsResult, records...)
		if len(records) < limit {
			break
		}

		offset += len(records)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenIpBlacklistData(total, recordsResult)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenIpBlacklistData(total float64, records []interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"records": flattenIpBlacklistRecords(records),
			"total":   total,
		},
	}
}

func flattenIpBlacklistRecords(records []interface{}) []interface{} {
	result := make([]interface{}, 0, len(records))
	for _, v := range records {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"effect_scope": utils.ExpandToIntList(
				utils.PathSearch("effect_scope", v, make([]interface{}, 0)).([]interface{})),
			"import_status": utils.PathSearch("import_status", v, nil),
			"import_time":   utils.PathSearch("import_time", v, nil),
			"error_message": utils.PathSearch("error_message", v, nil),
		})
	}

	return result
}
