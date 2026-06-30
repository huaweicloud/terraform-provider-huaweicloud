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

// @API DRS GET /v5/{project_id}/drivers
func DataSourceDrsDrivers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsDriversRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"driver_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"driver_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildDriversQueryParams(d *schema.ResourceData, offset int) string {
	queryParams := fmt.Sprintf("driver_type=%s", d.Get("driver_type").(string))

	if offset > 0 {
		queryParams += fmt.Sprintf("&offset=%d", offset)
	}

	return queryParams
}

func dataSourceDrsDriversRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/drivers"
		result  = make([]interface{}, 0)
		offset  = 0
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		queryParams := buildDriversQueryParams(d, offset)
		currentListPath := fmt.Sprintf("%s?%s", listPath, queryParams)

		listResp, err := client.Request("GET", currentListPath, &reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS drivers data: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		drivers := utils.PathSearch("items", listRespBody, make([]interface{}, 0)).([]interface{})

		if len(drivers) == 0 {
			break
		}

		result = append(result, drivers...)

		offset += len(drivers)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("items", flattenDrivers(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDrivers(driversResp []interface{}) []interface{} {
	if len(driversResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(driversResp))
	for _, v := range driversResp {
		rst = append(rst, map[string]interface{}{
			"driver_name":   utils.PathSearch("driver_name", v, nil),
			"last_modified": utils.PathSearch("last_modified", v, nil),
			"size":          utils.PathSearch("size", v, nil),
		})
	}
	return rst
}
