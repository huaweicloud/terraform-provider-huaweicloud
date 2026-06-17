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
			"databases": {
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
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		userName   = d.Get("user_name").(string)
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	users, err := QueryHtapStarrocksUsers(client, instanceId, userName)
	if err != nil {
		return diag.Errorf("error querying TaurusDB users: %s", err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("user_details", flattenStarrocksUsersUserDetails(users)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func QueryHtapStarrocksUsers(client *golangsdk.ServiceClient, instanceId, userName string) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/starrocks/users"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		queryPath := buildStarrocksUsersQueryParams(listPath, userName, limit, offset)
		resp, err := client.Request("GET", queryPath, &listOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		userDetails := utils.PathSearch("user_details", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, userDetails...)
		if len(userDetails) < limit {
			break
		}

		offset += len(userDetails)
	}
	return result, nil
}

func buildStarrocksUsersQueryParams(baseUrl, userName string, limit, offset int) string {
	rst := fmt.Sprintf("%s?limit=%d&offset=%d", baseUrl, limit, offset)
	if userName != "" {
		rst += fmt.Sprintf("&user_name=%s", userName)
	}
	return rst
}

func flattenStarrocksUsersUserDetails(users []interface{}) []interface{} {
	res := make([]interface{}, 0, len(users))
	for _, v := range users {
		dataBases := utils.PathSearch("databases", v, make([]interface{}, 0)).([]interface{})
		res = append(res, map[string]interface{}{
			"user_name": utils.PathSearch("user_name", v, nil),
			"databases": utils.ExpandToStringList(dataBases),
			"dml":       utils.PathSearch("dml", v, nil),
			"ddl":       utils.PathSearch("ddl", v, nil),
		})
	}
	return res
}
