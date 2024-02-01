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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/database/detail
func DataSourceSQLServerDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceSQLServerDatabasesRead,
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"character_set": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sqlserverDatabasesSchema(),
			},
		},
	}
}

func sqlserverDatabasesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"character_set": {
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

func resourceSQLServerDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listSQLServerDatabasesHttpUrl = "v3/{project_id}/instances/{instance_id}/database/detail?page=1&limit=100"
		listSQLServerDatabasesProduct = "rds"
	)
	listSQLServerDatabasesClient, err := cfg.NewServiceClient(listSQLServerDatabasesProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listSQLServerDatabasesPath := listSQLServerDatabasesClient.Endpoint + listSQLServerDatabasesHttpUrl
	listSQLServerDatabasesPath = strings.ReplaceAll(listSQLServerDatabasesPath, "{project_id}",
		listSQLServerDatabasesClient.ProjectID)
	listSQLServerDatabasesPath = strings.ReplaceAll(listSQLServerDatabasesPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))

	listSQLServerDatabasesResp, err := pagination.ListAllItems(
		listSQLServerDatabasesClient,
		"page",
		listSQLServerDatabasesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving RDS SQLServer databases, %s", err)
	}

	listSQLServerDatabasesRespJson, err := json.Marshal(listSQLServerDatabasesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listSQLServerDatabasesRespBody interface{}
	err = json.Unmarshal(listSQLServerDatabasesRespJson, &listSQLServerDatabasesRespBody)
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
		d.Set("databases", flattenListSQLServerDatabasesBody(
			filterListSQLServerDatabases(d, listSQLServerDatabasesRespBody))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListSQLServerDatabasesBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":          utils.PathSearch("name", v, nil),
			"character_set": utils.PathSearch("character_set", v, nil),
			"state":         utils.PathSearch("state", v, nil),
		})
	}
	return rst
}

func filterListSQLServerDatabases(d *schema.ResourceData, resp interface{}) []interface{} {
	databaseJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	databaseArray := databaseJson.([]interface{})
	if len(databaseArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(databaseArray))

	rawName, rawNameOK := d.GetOk("name")
	rawState, rawStateOK := d.GetOk("state")
	rawCharacterSet, rawCharacterSetOK := d.GetOk("character_set")
	for _, database := range databaseArray {
		name := utils.PathSearch("name", database, nil)
		state := utils.PathSearch("state", database, nil)
		characterSet := utils.PathSearch("character_set", database, nil)
		if rawNameOK && rawName != name {
			continue
		}
		if rawStateOK && rawState != state {
			continue
		}
		if rawCharacterSetOK && rawCharacterSet != characterSet {
			continue
		}
		result = append(result, database)
	}

	return result
}
