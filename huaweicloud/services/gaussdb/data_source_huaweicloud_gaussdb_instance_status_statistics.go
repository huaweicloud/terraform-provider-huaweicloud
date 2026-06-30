package gaussdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/instances-statistics
func DataSourceInstanceStatusStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceStatusStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instances_statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     instanceStatusStatisticsSchema(),
			},
		},
	}
}

func instanceStatusStatisticsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
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

func dataSourceInstanceStatusStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/instances-statistics"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.Errorf("error querying GaussDB instance status statistics: %s", err)
	}

	resp, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instances_statistics", flattenInstanceStatusStatistics(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
func flattenInstanceStatusStatistics(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("instances_statistics", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"status": utils.PathSearch("status", v, nil),
			"count":  utils.PathSearch("count", v, nil),
		})
	}
	return res
}
