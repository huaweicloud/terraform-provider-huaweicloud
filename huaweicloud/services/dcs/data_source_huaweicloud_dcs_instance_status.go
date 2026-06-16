package dcs

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

// @API DCS GET /v2/{project_id}/instances/status
func DataSourceDcsInstanceStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsInstanceStatusRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"include_failure": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"paying_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"freezing_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"migrating_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"flushing_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"upgrading_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"restoring_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"extending_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creating_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"running_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"error_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"frozen_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"createfailed_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"restarting_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"redis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsStatusStatisticSchema(),
			},
			"memcached": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsStatusStatisticSchema(),
			},
		},
	}
}

func dcsStatusStatisticSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"paying_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"freezing_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"migrating_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"flushing_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"upgrading_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"restoring_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"extending_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creating_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"running_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"error_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"frozen_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"createfailed_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"restarting_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsInstanceStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	url := "v2/{project_id}/instances/status"
	getPath := client.Endpoint + url
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGetInstanceStatusQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error querying DCS instance status: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	redis := utils.PathSearch("redis", getRespBody, nil)
	memcached := utils.PathSearch("memcached", getRespBody, nil)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("redis", flattenInstanceStatusStatistic(redis)),
		d.Set("memcached", flattenInstanceStatusStatistic(memcached)),
		d.Set("paying_count", utils.PathSearch("paying_count", getRespBody, 0)),
		d.Set("freezing_count", utils.PathSearch("freezing_count", getRespBody, 0)),
		d.Set("migrating_count", utils.PathSearch("migrating_count", getRespBody, 0)),
		d.Set("flushing_count", utils.PathSearch("flushing_count", getRespBody, 0)),
		d.Set("upgrading_count", utils.PathSearch("upgrading_count", getRespBody, 0)),
		d.Set("restoring_count", utils.PathSearch("restoring_count", getRespBody, 0)),
		d.Set("extending_count", utils.PathSearch("extending_count", getRespBody, 0)),
		d.Set("creating_count", utils.PathSearch("creating_count", getRespBody, 0)),
		d.Set("running_count", utils.PathSearch("running_count", getRespBody, 0)),
		d.Set("error_count", utils.PathSearch("error_count", getRespBody, 0)),
		d.Set("frozen_count", utils.PathSearch("frozen_count", getRespBody, 0)),
		d.Set("createfailed_count", utils.PathSearch("createfailed_count", getRespBody, 0)),
		d.Set("restarting_count", utils.PathSearch("restarting_count", getRespBody, 0)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetInstanceStatusQueryParams(d *schema.ResourceData) string {
	if includeFailure, ok := d.GetOk("include_failure"); ok {
		return fmt.Sprintf("?include_failure=%v", includeFailure)
	}

	return ""
}

func flattenInstanceStatusStatistic(statusStatistic interface{}) []interface{} {
	if statusStatistic == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"paying_count":       utils.PathSearch("paying_count", statusStatistic, 0),
			"freezing_count":     utils.PathSearch("freezing_count", statusStatistic, 0),
			"migrating_count":    utils.PathSearch("migrating_count", statusStatistic, 0),
			"flushing_count":     utils.PathSearch("flushing_count", statusStatistic, 0),
			"upgrading_count":    utils.PathSearch("upgrading_count", statusStatistic, 0),
			"restoring_count":    utils.PathSearch("restoring_count", statusStatistic, 0),
			"extending_count":    utils.PathSearch("extending_count", statusStatistic, 0),
			"creating_count":     utils.PathSearch("creating_count", statusStatistic, 0),
			"running_count":      utils.PathSearch("running_count", statusStatistic, 0),
			"error_count":        utils.PathSearch("error_count", statusStatistic, 0),
			"frozen_count":       utils.PathSearch("frozen_count", statusStatistic, 0),
			"createfailed_count": utils.PathSearch("createfailed_count", statusStatistic, 0),
			"restarting_count":   utils.PathSearch("restarting_count", statusStatistic, 0),
		},
	}
	return rst
}
