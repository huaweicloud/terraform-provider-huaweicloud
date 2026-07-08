package gaussdb

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

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/session/statistic/wait-event
func DataSourceGaussdbRtsWaitingEventStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussdbRtsWaitingEventStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"wait_events_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussdbRtsWaitingEventStatisticsSchema(),
			},
		},
	}
}

func gaussdbRtsWaitingEventStatisticsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussdbRtsWaitingEventStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/session/statistic/wait-event"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	}

	offset := 0
	res := make([]interface{}, 0)
	for {
		getPathWithPage := fmt.Sprintf("%s?offset=%d&limit=%d", getPath, offset, 100)
		getResp, err := client.Request("GET", getPathWithPage, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB session wait events: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		waitEvents := flattenGetGaussdbRtsWaitingEventStatisticsBody(getRespBody)
		if len(waitEvents) == 0 {
			break
		}
		res = append(res, waitEvents...)
		offset += 100
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("wait_events_info", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetGaussdbRtsWaitingEventStatisticsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("wait_event_info", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"node_name":  utils.PathSearch("node_name", v, nil),
			"event_name": utils.PathSearch("event_name", v, nil),
			"count":      utils.PathSearch("count", v, nil),
		})
	}
	return res
}
