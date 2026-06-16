package taurusdb

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

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/database-parameters
func DataSourceTaurusDBHtapStarrocksReplicationDBParameters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapStarrocksReplicationDBParametersRead,

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
			"add_task_scenario": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"main_task_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksReplicationDbParametersSchema(),
			},
		},
	}
}

func starrocksReplicationDbParametersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"param_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_modifiable": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBHtapStarrocksReplicationDBParametersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		limit  = 100
		offset = 0
		result = make([]interface{}, 0)
	)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/database-parameters"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		queryPath := buildStarrocksReplicationDbParametersQueryParams(listPath, d, limit, offset)
		resp, err := client.Request("GET", queryPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB HTAP StarRocks replication database parameters: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dbParameters := utils.PathSearch("db_parameters", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, dbParameters...)
		if len(dbParameters) < limit {
			break
		}

		offset += len(dbParameters)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("db_parameters", flattenStarrocksReplicationDbParameters(result)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildStarrocksReplicationDbParametersQueryParams(baseUrl string, d *schema.ResourceData, limit, offset int) string {
	rst := fmt.Sprintf("%s?limit=%d&offset=%d", baseUrl, limit, offset)
	if v, ok := d.GetOk("add_task_scenario"); ok {
		rst += fmt.Sprintf("&add_task_scenario=%s", v.(string))
	}
	if v, ok := d.GetOk("main_task_name"); ok {
		rst += fmt.Sprintf("&main_task_name=%s", v.(string))
	}
	return rst
}

func flattenStarrocksReplicationDbParameters(resp interface{}) []interface{} {
	curArray := resp.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"param_name":    utils.PathSearch("param_name", v, nil),
			"data_type":     utils.PathSearch("data_type", v, nil),
			"default_value": utils.PathSearch("default_value", v, nil),
			"value_range":   utils.PathSearch("value_range", v, nil),
			"description":   utils.PathSearch("description", v, nil),
			"is_modifiable": utils.PathSearch("is_modifiable", v, nil),
		})
	}
	return res
}
