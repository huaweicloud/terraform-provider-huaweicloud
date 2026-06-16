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

// @API GaussDB GET /v3/{project_id}/instances/metric-data
func DataSourceGaussDbAllInstancesMetrics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDbAllInstancesMetricsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbAllInstancesMetricsSchema(),
			},
		},
	}
}

func gaussDbAllInstancesMetricsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"solution": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_used_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_total_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_usage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"p80": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"p95": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deadlocks": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"buffer_hit_ratio": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbAllInstancesMetricsNodeSchema(),
			},
		},
	}
}

func gaussDbAllInstancesMetricsNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"component_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGaussDbAllInstancesMetricsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/metric-data"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	limit := 50
	offset := 0
	allInstances := make([]interface{}, 0)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		queryParams := buildGaussDbInstancesQueryParams(d, offset, limit)
		requestPath := getPath + queryParams

		getResp, err := client.Request("GET", requestPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB all instances metrics: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flattening GaussDB all instances metrics response: %s", err)
		}

		instances := flattenGaussDbAllInstancesMetrics(getRespBody)
		if len(instances) == 0 {
			break
		}
		allInstances = append(allInstances, instances...)
		offset += limit
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("instances", allInstances),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGaussDbInstancesQueryParams(d *schema.ResourceData, offset int, limit int) string {
	res := fmt.Sprintf("?offset=%d", offset)
	res = fmt.Sprintf("%s&limit=%d", res, limit)

	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	return res
}

func flattenGaussDbAllInstancesMetrics(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"status":           utils.PathSearch("status", v, nil),
			"mode":             utils.PathSearch("mode", v, nil),
			"engine_name":      utils.PathSearch("engine_name", v, nil),
			"engine_version":   utils.PathSearch("engine_version", v, nil),
			"solution":         utils.PathSearch("solution", v, nil),
			"disk_used_size":   utils.PathSearch("disk_used_size", v, nil),
			"disk_total_size":  utils.PathSearch("disk_total_size", v, nil),
			"disk_usage":       utils.PathSearch("disk_usage", v, nil),
			"p80":              utils.PathSearch("p80", v, nil),
			"p95":              utils.PathSearch("p95", v, nil),
			"deadlocks":        utils.PathSearch("deadlocks", v, nil),
			"buffer_hit_ratio": utils.PathSearch("buffer_hit_ratio", v, nil),
			"nodes":            flattenGaussDbAllInstancesMetricsNodes(v),
		})
	}
	return res
}

func flattenGaussDbAllInstancesMetricsNodes(resp interface{}) []interface{} {
	curJson := utils.PathSearch("nodes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"role":          utils.PathSearch("role", v, nil),
			"component_ids": utils.PathSearch("component_ids", v, nil),
		})
	}
	return res
}
