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

// @API CFW GET /v1/{project_id}/schedules
func DataSourceSchedules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSchedulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ref_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"periodic": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataRecordsPeriodicElem(),
						},
						"absolute": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataRecordsAbsoluteElem(),
						},
					},
				},
			},
		},
	}
}

func dataRecordsAbsoluteElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataRecordsPeriodicElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"week_mask": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"start_week": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"end_week": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildQuerySchedulesQueryParams(d *schema.ResourceData) string {
	queryParam := fmt.Sprintf("?object_id=%s&limit=1024", d.Get("object_id").(string))
	if v, ok := d.GetOk("name"); ok {
		queryParam += fmt.Sprintf("&name=%s", v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		queryParam += fmt.Sprintf("&desc=%s", v.(string))
	}

	return queryParam
}

func dataSourceSchedulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "cfw"
		httpUrl    = "v1/{project_id}/schedules"
		offset     = 0
		allResults = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildQuerySchedulesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving CFW schedules: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		records := utils.PathSearch("data.records", respBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		allResults = append(allResults, records...)
		offset += len(records)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("records", flattenSchedulesResponse(allResults)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSchedulesResponse(respArray []interface{}) []map[string]interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"schedule_id": utils.PathSearch("schedule_id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"ref_count":   utils.PathSearch("ref_count", v, nil),
			"periodic":    flattenPeriodicResponse(utils.PathSearch("periodic", v, make([]interface{}, 0)).([]interface{})),
			"absolute":    flattenAbsoluteResponse(utils.PathSearch("absolute", v, nil)),
		})
	}

	return rst
}

func flattenAbsoluteResponse(respMap interface{}) []map[string]interface{} {
	if respMap == nil {
		return nil
	}

	rstMap := map[string]interface{}{
		"start_time": utils.PathSearch("start_time", respMap, nil),
		"end_time":   utils.PathSearch("end_time", respMap, nil),
	}

	return []map[string]interface{}{rstMap}
}

func flattenPeriodicResponse(respArray []interface{}) []map[string]interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"type":       utils.PathSearch("type", v, nil),
			"start_time": utils.PathSearch("start_time", v, nil),
			"end_time":   utils.PathSearch("end_time", v, nil),
			"week_mask":  utils.PathSearch("week_mask", v, nil),
			"start_week": utils.PathSearch("start_week", v, nil),
			"end_week":   utils.PathSearch("end_week", v, nil),
		})
	}

	return rst
}
