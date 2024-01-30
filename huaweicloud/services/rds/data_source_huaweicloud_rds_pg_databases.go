package rds

import (
	"context"
	"encoding/json"
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
func DataSourcePgDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePgDatabasesRead,
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
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"character_set": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lc_collate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"databases": {
				Type:     schema.TypeList,
				Elem:     pgDatabasesSchema(),
				Computed: true,
			},
		},
	}
}

func pgDatabasesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"character_set": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lc_collate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
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

func dataSourcePgDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listPgDatabasesHttpUrl = "v3/{project_id}/instances/{instance_id}/database/detail?page=1&limit=100"
		listPgDatabasesProduct = "rds"
	)

	listPgDatabasesClient, err := cfg.NewServiceClient(listPgDatabasesProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listPgDatabasesPath := listPgDatabasesClient.Endpoint + listPgDatabasesHttpUrl
	listPgDatabasesPath = strings.ReplaceAll(listPgDatabasesPath, "{project_id}", listPgDatabasesClient.ProjectID)
	listPgDatabasesPath = strings.ReplaceAll(listPgDatabasesPath, "{instance_id}", d.Get("instance_id").(string))

	listPgDatabasesResp, err := pagination.ListAllItems(
		listPgDatabasesClient,
		"page",
		listPgDatabasesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS PostgreSQL databases")
	}

	listPgDatabasesRespJson, err := json.Marshal(listPgDatabasesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listPgDatabasesRespBody interface{}
	err = json.Unmarshal(listPgDatabasesRespJson, &listPgDatabasesRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("databases", flattenPgDatabasesBody(filterPgDatabases(d, listPgDatabasesRespBody))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPgDatabasesBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":          utils.PathSearch("name", v, nil),
			"owner":         utils.PathSearch("owner", v, nil),
			"character_set": utils.PathSearch("character_set", v, nil),
			"lc_collate":    utils.PathSearch("collate_set", v, nil),
			"size":          utils.PathSearch("size", v, 0),
			"description":   utils.PathSearch("comment", v, nil),
		})
	}
	return rst
}

func filterPgDatabases(d *schema.ResourceData, resp interface{}) []interface{} {
	databasesJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	databaseArray := databasesJson.([]interface{})
	if len(databaseArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0)

	rawName, rawNameOK := d.GetOk("name")
	rawOwner, rawOwnerOK := d.GetOk("owner")
	rawCharacterSet, rawCharacterSetOK := d.GetOk("character_set")
	rawLcCollate, rawLcCollateOK := d.GetOk("lc_collate")
	rawSize, rawSizeOK := d.GetOk("size")
	for _, database := range databaseArray {
		name := utils.PathSearch("name", database, nil)
		owner := utils.PathSearch("owner", database, nil)
		characterSet := utils.PathSearch("character_set", database, nil)
		lcCollate := utils.PathSearch("collate_set", database, nil)
		size := utils.PathSearch("size", database, 0)
		if rawNameOK && rawName != name {
			continue
		}
		if rawOwnerOK && rawOwner != owner {
			continue
		}
		if rawCharacterSetOK && rawCharacterSet != characterSet {
			continue
		}
		if rawLcCollateOK && rawLcCollate != lcCollate {
			continue
		}
		if rawSizeOK && rawSize != size {
			continue
		}
		result = append(result, database)
	}
	return result
}
