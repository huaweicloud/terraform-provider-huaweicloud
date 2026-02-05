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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/database/db_user
func DataSourceRdsMysqlDatabasePrivileges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsMysqlDatabasePrivilegesRead,
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
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"readonly": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     usersSchema(),
			},
		},
	}
}

func usersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"readonly": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceRdsMysqlDatabasePrivilegesRead(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listMysqlDatabasePrivilegesHttpUrl = "v3/{project_id}/instances/{instance_id}/database/db_user"
		listMysqlDatabasePrivilegesProduct = "rds"
	)
	listMysqlDatabasePrivilegesClient, err := cfg.NewServiceClient(listMysqlDatabasePrivilegesProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listMysqlDatabasePrivilegesPath := listMysqlDatabasePrivilegesClient.Endpoint + listMysqlDatabasePrivilegesHttpUrl
	listMysqlDatabasePrivilegesPath = strings.ReplaceAll(listMysqlDatabasePrivilegesPath, "{project_id}",
		listMysqlDatabasePrivilegesClient.ProjectID)
	listMysqlDatabasePrivilegesPath = strings.ReplaceAll(listMysqlDatabasePrivilegesPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))

	listMysqlDatabasePrivilegesQueryParams := buildGetMysqlDatabasePrivilegesQueryParams(
		fmt.Sprintf("%v", d.Get("db_name")))
	listMysqlDatabasePrivilegesPath += listMysqlDatabasePrivilegesQueryParams

	listMysqlDatabasePrivilegesResp, err := pagination.ListAllItems(
		listMysqlDatabasePrivilegesClient,
		"page",
		listMysqlDatabasePrivilegesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving RDS MySQL database privileges, %s", err)
	}

	listMysqlDatabasePrivilegesRespJson, err := json.Marshal(listMysqlDatabasePrivilegesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listMysqlDatabasePrivilegesRespBody interface{}
	err = json.Unmarshal(listMysqlDatabasePrivilegesRespJson, &listMysqlDatabasePrivilegesRespBody)
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
		d.Set("users", flattenListDatabasePrivilegesBody(
			filterListDatabasePrivileges(d, listMysqlDatabasePrivilegesRespBody))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetMysqlDatabasePrivilegesQueryParams(dbName string) string {
	return fmt.Sprintf("?db-name=%s&page=1&limit=100", dbName)
}

func flattenListDatabasePrivilegesBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":     utils.PathSearch("name", v, nil),
			"readonly": utils.PathSearch("readonly", v, nil),
		})
	}
	return rst
}

func filterListDatabasePrivileges(d *schema.ResourceData, resp interface{}) []interface{} {
	usersJson := utils.PathSearch("users", resp, make([]interface{}, 0))
	userArray := usersJson.([]interface{})
	if len(userArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(userArray))

	rawUserName, rawUserNameOK := d.GetOk("user_name")
	rawReadonly := d.Get("readonly")
	for _, user := range userArray {
		name := utils.PathSearch("name", user, nil)
		readonly := utils.PathSearch("readonly", user, nil)
		if rawUserNameOK && rawUserName != name {
			continue
		}
		if rawReadonly != readonly {
			continue
		}
		result = append(result, user)
	}

	return result
}
