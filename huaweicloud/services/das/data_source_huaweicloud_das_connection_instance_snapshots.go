package das

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DAS GET /v3/{project_id}/connections/{connection_id}/instance/query-snapshots
func DataSourceConnectionInstanceSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionInstanceSnapshotsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the connection instance snapshots are located.",
			},

			// Required parameters.
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database user ID.",
			},
			"module": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The lock snapshot type.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The start time of the query range, in RFC3339 format.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The end time of the query range, in RFC3339 format.",
			},

			// Attributes.
			"snapshots": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of lock snapshots.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The snapshot ID.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The snapshot status.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The snapshot creation time, in RFC3339 format.",
						},
						"find_lock": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether a lock was found.",
						},
					},
				},
			},
		},
	}
}

func dataSourceConnectionInstanceSnapshotsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	snapshots, err := listConnectionInstanceSnapshots(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS connection instance snapshots: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("snapshots", flattenConnectionInstanceSnapshots(snapshots)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildConnectionInstanceSnapshotsQueryParams(d *schema.ResourceData, curPage, perPage int) string {
	startAt := utils.ConvertTimeStrToNanoTimestamp(d.Get("start_time").(string))
	endAt := utils.ConvertTimeStrToNanoTimestamp(d.Get("end_time").(string))

	return fmt.Sprintf("?module=%d&start_at=%d&end_at=%d&per_page=%d&cur_page=%d",
		d.Get("module").(int), startAt, endAt, perPage, curPage)
}

func listConnectionInstanceSnapshots(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/connections/{connection_id}/instance/query-snapshots"
		perPage = 100
		curPage = 1
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{connection_id}", d.Get("user_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
	}

	for {
		listPathWithPage := listPath + buildConnectionInstanceSnapshotsQueryParams(d, curPage, perPage)

		requestResp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, items...)

		if len(items) < perPage {
			break
		}
		curPage++
	}

	return result, nil
}

func flattenConnectionInstanceSnapshots(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":        utils.PathSearch("id", item, nil),
			"status":    utils.PathSearch("status", item, nil),
			"find_lock": utils.PathSearch("find_lock", item, nil),
			"created_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_at", item, float64(0)).(float64))/1000, false),
		})
	}

	return result
}
