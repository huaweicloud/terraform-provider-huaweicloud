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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/database/detail
func DataSourceRdsMysqlDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsMysqlDatabasesRead,
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
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     databasesSchema(),
			},
		},
	}
}

func databasesSchema() *schema.Resource {
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
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceRdsMysqlDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listMysqlDatabasesHttpUrl = "v3/{project_id}/instances/{instance_id}/database/detail?page=1&limit=100"
		listMysqlDatabasesProduct = "rds"
	)
	listMysqlDatabasesClient, err := cfg.NewServiceClient(listMysqlDatabasesProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listMysqlDatabasesPath := listMysqlDatabasesClient.Endpoint + listMysqlDatabasesHttpUrl
	listMysqlDatabasesPath = strings.ReplaceAll(listMysqlDatabasesPath, "{project_id}", listMysqlDatabasesClient.ProjectID)
	listMysqlDatabasesPath = strings.ReplaceAll(listMysqlDatabasesPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	listMysqlDatabasesResp, err := pagination.ListAllItems(
		listMysqlDatabasesClient,
		"page",
		listMysqlDatabasesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS Mysql databases")
	}

	listMysqlDatabasesRespJson, err := json.Marshal(listMysqlDatabasesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listMysqlDatabasesRespBody interface{}
	err = json.Unmarshal(listMysqlDatabasesRespJson, &listMysqlDatabasesRespBody)
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
		d.Set("databases", flattenListDatabasesBody(filterListDatabases(d, listMysqlDatabasesRespBody))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListDatabasesBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":          utils.PathSearch("name", v, nil),
			"character_set": utils.PathSearch("character_set", v, nil),
			"description":   utils.PathSearch("comment", v, nil),
		})
	}
	return rst
}

func filterListDatabases(d *schema.ResourceData, resp interface{}) []interface{} {
	databaseJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	databaseArray := databaseJson.([]interface{})
	if len(databaseArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(databaseArray))

	rawName, rawNameOK := d.GetOk("name")
	rawCharacterSet, rawCharacterSetOK := d.GetOk("character_set")
	for _, database := range databaseArray {
		name := utils.PathSearch("name", database, nil)
		characterSet := utils.PathSearch("character_set", database, nil)
		if rawNameOK && rawName != name {
			continue
		}
		if rawCharacterSetOK && rawCharacterSet != characterSet {
			continue
		}
		result = append(result, database)
	}

	return result
}
