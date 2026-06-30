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

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/databases
func DataSourceTaurusDBHtapStarrocksDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapStarrocksDatabasesRead,

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
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTaurusDBHtapStarrocksDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		offset = 0
		result = make([]interface{}, 0)
	)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/databases"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		queryPath := buildStarRocksDatabasesQueryParams(listPath, d, offset)
		resp, err := client.Request("GET", queryPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB Htap StarRocks databases: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		databases := utils.PathSearch("databases", respBody, make([]interface{}, 0)).([]interface{})
		if len(databases) == 0 {
			break
		}

		result = append(result, databases...)

		// API returns total_count always is 0, so we check if the returned count is less than limit
		if len(databases) < 100 {
			break
		}

		offset += len(databases)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("databases", flattenHtapStarrocksDatabasesBody(result)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildStarRocksDatabasesQueryParams(baseUrl string, d *schema.ResourceData, offset int) string {
	rst := fmt.Sprintf("%s?limit=100&offset=%d", baseUrl, offset)

	if v, ok := d.GetOk("database_name"); ok {
		rst += fmt.Sprintf("&database_name=%s", v.(string))
	}

	return rst
}
func flattenHtapStarrocksDatabasesBody(resp interface{}) []interface{} {
	curArray := resp.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, v.(string))
	}
	return res
}
