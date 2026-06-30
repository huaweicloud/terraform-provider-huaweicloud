package drs

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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/dirty-data
func DataSourceDrsDirtyData() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsDirtyDataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dirty_data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schema_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_sql": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildDirtyDataQueryParams(d *schema.ResourceData, offset int) string {
	queryParams := "?limit=1000"
	if v, ok := d.GetOk("begin_time"); ok {
		queryParams += fmt.Sprintf("&begin_time=%s", v.(string))
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams += fmt.Sprintf("&end_time=%s", v.(string))
	}
	if offset > 0 {
		queryParams += fmt.Sprintf("&offset=%d", offset)
	}

	return queryParams
}

func dataSourceDrsDirtyDataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/jobs/{job_id}/dirty-data"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		requestPathWithQuery := requestPath + buildDirtyDataQueryParams(d, offset)
		resp, err := client.Request("GET", requestPathWithQuery, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS dirty data: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dirtyDataList := utils.PathSearch("dirty_data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dirtyDataList) == 0 {
			break
		}

		result = append(result, dirtyDataList...)
		offset += len(dirtyDataList)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("dirty_data_list", flattenDirtyDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDirtyDataList(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, item := range respArray {
		result = append(result, map[string]interface{}{
			"db_name":     utils.PathSearch("db_name", item, nil),
			"schema_name": utils.PathSearch("schema_name", item, nil),
			"table_name":  utils.PathSearch("table_name", item, nil),
			"error_sql":   utils.PathSearch("error_sql", item, nil),
			"error_time":  utils.PathSearch("error_time", item, nil),
			"error_msg":   utils.PathSearch("error_msg", item, nil),
		})
	}

	return result
}
