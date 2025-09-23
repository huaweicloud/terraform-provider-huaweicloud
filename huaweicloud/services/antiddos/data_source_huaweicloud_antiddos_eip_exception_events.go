package antiddos

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

// @API ANTI-DDOS GET /v1/{project_id}/antiddos/{floating_ip_id}/logs
func DataSourceEipExceptionEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEipExceptionEventsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"floating_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the EIP.`,
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the EIP address.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sort direction.`,
			},
			"logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of exception logs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The start time of the exception event.`,
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The end time of the exception event.`,
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The protection status.`,
						},
						"trigger_bps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The traffic when the exception event is triggered.`,
						},
						"trigger_pps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The packet rate when the exception event is triggered.`,
						},
						"trigger_http_pps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The HTTP request rate when the exception event is triggered.`,
						},
					},
				},
			},
		},
	}
}

func buildEipExceptionEventsQueryParams(d *schema.ResourceData) string {
	rst := ""

	if v, ok := d.GetOk("ip"); ok {
		rst += fmt.Sprintf("&ip=%s", v.(string))
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%s", v.(string))
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

// When the limit and offset fields are not included, the API defaults to querying the full data.
// Therefore, no paging logic is done here.
func dataSourceEipExceptionEventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/antiddos/{floating_ip_id}/logs"
		product = "anti-ddos"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{floating_ip_id}", d.Get("floating_ip_id").(string))
	requestPath += buildEipExceptionEventsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving Anti-DDoS exception events: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("logs", flattenExceptionEvents(utils.PathSearch("logs", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenExceptionEvents(logsRaw []interface{}) []map[string]interface{} {
	if len(logsRaw) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(logsRaw))
	for _, log := range logsRaw {
		result = append(result, map[string]interface{}{
			"start_time":       utils.PathSearch("start_time", log, nil),
			"end_time":         utils.PathSearch("end_time", log, nil),
			"status":           utils.PathSearch("status", log, nil),
			"trigger_bps":      utils.PathSearch("trigger_bps", log, nil),
			"trigger_pps":      utils.PathSearch("trigger_pps", log, nil),
			"trigger_http_pps": utils.PathSearch("trigger_http_pps", log, nil),
		})
	}

	return result
}
