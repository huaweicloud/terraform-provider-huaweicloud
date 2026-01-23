package eps

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EPS GET /v1.0/enterprise-projects/migrate-record/list
func DataSourceMigrateRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMigrateRecordRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildMigrateRecordSchema(),
			},
		},
	}
}

func buildMigrateRecordSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"associated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// The API documentation describes this field `event_time` as a string type, but it is actually a float type.
			"event_time": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operate_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initiate_ep_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initiate_ep_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin_ep_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin_ep_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_ep_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_ep_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// The API documentation describes this field `exist_failed` as a string type, but it is actually a boolean type.
			"exist_failed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildMigrateRecordQueryParams(d *schema.ResourceData, startTime, endTime, offset string) string {
	rst := ""
	if v, ok := d.GetOk("resource_id"); ok {
		rst += fmt.Sprintf("&resource_id=%s", v)
	}

	if startTime != "" {
		rst += fmt.Sprintf("&start_time=%s", startTime)
	}

	if endTime != "" {
		rst += fmt.Sprintf("&end_time=%s", endTime)
	}

	if offset != "" {
		rst += fmt.Sprintf("&offset=%s", offset)
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

// Convert the time string to the corresponding timestamp (in millisecond).
func convertTimeQueryParam(time string) (string, error) {
	if time == "" {
		return "", nil
	}

	timestamp, err := utils.FormatUTCTimeStamp(time)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(timestamp*1000, 10), nil
}

func dataSourceMigrateRecordRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1.0/enterprise-projects/migrate-record/list"
		allResults = make([]interface{}, 0)
		offset     = ""
	)
	client, err := cfg.NewServiceClient("eps", region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}

	startTime, err := convertTimeQueryParam(d.Get("start_time").(string))
	if err != nil {
		return diag.Errorf("error converting start_time: %s", err)
	}

	endTime, err := convertTimeQueryParam(d.Get("end_time").(string))
	if err != nil {
		return diag.Errorf("error converting end_time: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildMigrateRecordQueryParams(d, startTime, endTime, offset)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving EPS migrate records: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		records := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		allResults = append(allResults, records...)
		offset = utils.PathSearch("offset", respBody, "").(string)
		if offset == "" {
			break
		}
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("records", flattenListMigrateRecordResponseBody(allResults)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListMigrateRecordResponseBody(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"associated":       utils.PathSearch("associated", v, nil),
			"code":             utils.PathSearch("code", v, nil),
			"message":          utils.PathSearch("message", v, nil),
			"project_id":       utils.PathSearch("project_id", v, nil),
			"project_name":     utils.PathSearch("project_name", v, nil),
			"region_id":        utils.PathSearch("region_id", v, nil),
			"event_time":       utils.PathSearch("event_time", v, nil),
			"user_name":        utils.PathSearch("user_name", v, nil),
			"operate_type":     utils.PathSearch("operate_type", v, nil),
			"resource_id":      utils.PathSearch("resource_id", v, nil),
			"resource_name":    utils.PathSearch("resource_name", v, nil),
			"resource_type":    utils.PathSearch("resource_type", v, nil),
			"initiate_ep_id":   utils.PathSearch("initiate_ep_id", v, nil),
			"initiate_ep_name": utils.PathSearch("initiate_ep_name", v, nil),
			"origin_ep_id":     utils.PathSearch("origin_ep_id", v, nil),
			"origin_ep_name":   utils.PathSearch("origin_ep_name", v, nil),
			"target_ep_id":     utils.PathSearch("target_ep_id", v, nil),
			"target_ep_name":   utils.PathSearch("target_ep_name", v, nil),
			"exist_failed":     utils.PathSearch("exist_failed", v, nil),
		})
	}
	return rst
}
