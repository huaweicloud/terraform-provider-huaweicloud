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

// @API DNS GET /v2/resolver/queryloggingconfig
func DataSourceResolverAccessLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResolverAccessLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the resolver access logs are located.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the VPC to be queried.`,
			},
			"access_logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resolver access log.`,
						},
						"lts_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log group.`,
						},
						"lts_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log stream.`,
						},
						"vpc_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of VPC IDs associated with the resolver access log.`,
						},
					},
				},
				Description: `The list of resolver access logs that match the filter parameters.`,
			},
		},
	}
}

func buildListResolverAccessLogsQueryParams(d *schema.ResourceData, limit int) string {
	queryParam := fmt.Sprintf("?limit=%d", limit)
	if v, ok := d.GetOk("vpc_id"); ok {
		queryParam = fmt.Sprintf("%s&vpc_id=%v", queryParam, v)
	}

	return queryParam
}

func listResolverAccessLogs(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/resolver/queryloggingconfig"
		result  = make([]interface{}, 0)
		limit   = 500
		marker  = ""
	)

	listPath := client.Endpoint + httpUrl
	listPath += buildListResolverAccessLogsQueryParams(d, limit)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		accessLogs := utils.PathSearch("resolver_query_log_configs", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, accessLogs...)
		if len(accessLogs) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func dataSourceResolverAccessLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	accessLogs, err := listResolverAccessLogs(client, d)
	if err != nil {
		return diag.Errorf("error querying resolver access logs: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("access_logs", flattenResolverAccessLogs(accessLogs)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResolverAccessLogs(accessLogs []interface{}) []map[string]interface{} {
	if len(accessLogs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(accessLogs))
	for _, item := range accessLogs {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", item, nil),
			"lts_group_id": utils.PathSearch("lts_group_id", item, nil),
			"lts_topic_id": utils.PathSearch("lts_topic_id", item, nil),
			"vpc_ids":      utils.PathSearch("vpc_ids", item, nil),
		})
	}

	return result
}
