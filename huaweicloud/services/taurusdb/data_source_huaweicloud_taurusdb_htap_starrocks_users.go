package taurusdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/users
func DataSourceTaurusDBHtapStarrocksUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapStarrocksUsersRead,

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
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksUsersUserDetailsSchema(),
			},
		},
	}
}

func starrocksUsersUserDetailsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_bases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dml": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ddl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBHtapStarrocksUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/users"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		queryPath := buildStarrocksUsersQueryParams(listPath, d, limit, offset)
		resp, err := client.Request("GET", queryPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB HTAP StarRocks users: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		userDetails := utils.PathSearch("user_details", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, userDetails...)
		if len(userDetails) < limit {
			break
		}

		offset += len(userDetails)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("user_details", flattenStarrocksUsersUserDetails(result)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildStarrocksUsersQueryParams(baseUrl string, d *schema.ResourceData, limit, offset int) string {
	rst := fmt.Sprintf("%s?limit=%d&offset=%d", baseUrl, limit, offset)
	if v, ok := d.GetOk("user_name"); ok {
		rst += fmt.Sprintf("&user_name=%s", v.(string))
	}
	return rst
}

func flattenStarrocksUsersUserDetails(resp interface{}) []interface{} {
	curArray := resp.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		dataBases := utils.PathSearch("data_bases", v, make([]interface{}, 0)).([]interface{})
		dataBasesResult := make([]string, 0, len(dataBases))
		for _, db := range dataBases {
			dataBasesResult = append(dataBasesResult, db.(string))
		}

		res = append(res, map[string]interface{}{
			"user_name":  utils.PathSearch("user_name", v, nil),
			"data_bases": dataBasesResult,
			"dml":        utils.PathSearch("dml", v, nil),
			"ddl":        utils.PathSearch("ddl", v, nil),
		})
	}
	return res
}
