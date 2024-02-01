package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/db_user/detail
func DataSourceRdsSQLServerAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsSQLServerAccountsRead,
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
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Elem:     sqlServerUsersSchema(),
				Computed: true,
			},
		},
	}
}

func sqlServerUsersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceRdsSQLServerAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listSQLServerAccountsHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		listSQLServerAccountsProduct = "rds"
	)
	listSQLServerAccountsClient, err := cfg.NewServiceClient(listSQLServerAccountsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listSQLServerAccountsPath := listSQLServerAccountsClient.Endpoint + listSQLServerAccountsHttpUrl
	listSQLServerAccountsPath = strings.ReplaceAll(listSQLServerAccountsPath, "{project_id}",
		listSQLServerAccountsClient.ProjectID)
	listSQLServerAccountsPath = strings.ReplaceAll(listSQLServerAccountsPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))

	listSQLServerAccountsResp, err := pagination.ListAllItems(
		listSQLServerAccountsClient,
		"page",
		listSQLServerAccountsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving RDS SQLServer accounts %S", err)
	}

	listSQLServerAccountsRespJson, err := json.Marshal(listSQLServerAccountsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listSQLServerAccountsRespBody interface{}
	err = json.Unmarshal(listSQLServerAccountsRespJson, &listSQLServerAccountsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("users", flattenListSQLServerAccountsBody(filterListSQLServerAccounts(d, listSQLServerAccountsRespBody))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListSQLServerAccountsBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"state": utils.PathSearch("state", v, nil),
		})
	}
	return rst
}

func filterListSQLServerAccounts(d *schema.ResourceData, resp interface{}) []interface{} {
	usersJson := utils.PathSearch("users", resp, make([]interface{}, 0))
	userArray := usersJson.([]interface{})
	if len(userArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0)

	rawName, rawNameOK := d.GetOk("user_name")
	rawState, rawStateOk := d.GetOk("state")

	for _, user := range userArray {
		name := utils.PathSearch("name", user, nil)
		state := utils.PathSearch("state", user, nil)
		if rawNameOK && rawName != name {
			continue
		}
		if rawStateOk && rawState != state {
			continue
		}
		result = append(result, user)
	}

	return result
}
