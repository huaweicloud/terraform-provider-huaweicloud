package dsc

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

// @API DSC GET /v1/{project_id}/metadata/catalog/databases
func DataSourceCatalogInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCatalogInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the catalog instances are located.`,
			},

			// Optional parameters.
			"label_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"type_id"},
				Description:  `The ID of the group label to which the instance belongs.`,
			},
			"type_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the data type of the database instance.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the database instance.`,
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The address of the database instance.`,
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The access user of the database instance.`,
			},
			"col_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The key used to sort the instances.`,
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sorting method for query results.`,
			},

			// Attributes.
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of database instances.`,
				Elem:        catalogInstanceInfoSchema(),
			},
		},
	}
}

func catalogInstanceInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the database instance.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the database instance.`,
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The address of the database instance.`,
			},
			"db_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of databases that belong to the instance.`,
				Elem:        catalogInstanceDbInfoSchema(),
			},
			"sensitive_col_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of sensitive columns in the instance.`,
			},
			"sensitive_db_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of sensitive databases in the instance.`,
			},
			"sensitive_table_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of sensitive tables in the instance.`,
			},
			"user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The access user of the database instance.`,
			},
		},
	}
}

func catalogInstanceDbInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the database.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the database.`,
			},
			"db_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the database.`,
			},
			"asset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the asset to which the database belongs.`,
			},
			"classifications": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The classification list of the database.`,
			},
			"latest_scan_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest scan time of the database, in RFC3339 format.`,
			},
			"sensitive_level_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the sensitive level.`,
			},
			"color_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The color number corresponding to the database sensitive level.`,
			},
			"sensitive_table_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of sensitive tables in the database.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The tag list of the database.`,
			},
			"total_table_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of tables in the database.`,
			},
		},
	}
}

func buildCatalogInstancesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("label_id"); ok {
		res = fmt.Sprintf("%s&label_id=%v", res, v)
	}

	if v, ok := d.GetOk("type_id"); ok {
		res = fmt.Sprintf("%s&type_id=%v", res, v)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		res = fmt.Sprintf("%s&instance_name=%v", res, v)
	}

	if v, ok := d.GetOk("address"); ok {
		res = fmt.Sprintf("%s&address=%v", res, v)
	}

	if v, ok := d.GetOk("user"); ok {
		res = fmt.Sprintf("%s&user=%v", res, v)
	}

	if v, ok := d.GetOk("col_id"); ok {
		res = fmt.Sprintf("%s&col_id=%v", res, v)
	}

	if v, ok := d.GetOk("sort"); ok {
		res = fmt.Sprintf("%s&sort=%v", res, v)
	}

	return res
}

func listCatalogInstances(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/metadata/catalog/databases"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += fmt.Sprintf("?limit=%v", limit) + buildCatalogInstancesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%v", offset)
		resp, err := client.Request("GET", listPathWithOffset, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		instances := utils.PathSearch("instance_infos", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, instances...)
		if len(instances) < limit {
			break
		}

		offset += len(instances)
	}

	return result, nil
}

func dataSourceCatalogInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	instances, err := listCatalogInstances(client, d)
	if err != nil {
		return diag.Errorf("error retrieving DSC catalog instances: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", flattenCatalogInstances(instances)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCatalogInstances(instances []interface{}) []interface{} {
	if len(instances) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(instances))
	for _, v := range instances {
		rst = append(rst, map[string]interface{}{
			"instance_id":         utils.PathSearch("instance_id", v, nil),
			"instance_name":       utils.PathSearch("instance_name", v, nil),
			"address":             utils.PathSearch("address", v, nil),
			"db_infos":            flattenCatalogInstanceDatabases(utils.PathSearch("db_infos", v, make([]interface{}, 0)).([]interface{})),
			"sensitive_col_num":   utils.PathSearch("sensitive_col_num", v, nil),
			"sensitive_db_num":    utils.PathSearch("sensitive_db_num", v, nil),
			"sensitive_table_num": utils.PathSearch("sensitive_table_num", v, nil),
			"user":                utils.PathSearch("user", v, nil),
		})
	}

	return rst
}

func flattenCatalogInstanceDatabases(dbs []interface{}) []interface{} {
	if len(dbs) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dbs))
	for _, v := range dbs {
		rst = append(rst, map[string]interface{}{
			"db_id":           utils.PathSearch("dbId", v, nil),
			"db_name":         utils.PathSearch("db_name", v, nil),
			"db_type":         utils.PathSearch("db_type", v, nil),
			"asset_id":        utils.PathSearch("asset_id", v, nil),
			"classifications": utils.PathSearch("classifications", v, nil),
			"latest_scan_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("latest_scan_time",
				v, float64(0)).(float64))/1000, false),
			"sensitive_level_name":  utils.PathSearch("sensitive_level_name", v, nil),
			"color_number":          utils.PathSearch("color_number", v, nil),
			"sensitive_table_count": utils.PathSearch("sensitive_table_count", v, nil),
			"tags":                  utils.PathSearch("tags", v, nil),
			"total_table_count":     utils.PathSearch("total_table_count", v, nil),
		})
	}

	return rst
}
