package gaussdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3.1/{project_id}/instances/{instance_id}/database-volume
func DataSourceDatabaseStorageUsage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseStorageUsageRead,

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
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table_space_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     databaseStorageUsageSchema(),
			},
		},
	}
}

func databaseStorageUsageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_space_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDatabaseStorageUsageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3.1/{project_id}/instances/{instance_id}/database-volume"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildDatabaseStorageUsageQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GaussDB database storage usage: %s", err)
	}
	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("database_volumes", flattenDatabaseStorageUsageBody(listRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDatabaseStorageUsageQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("database_name"); ok {
		res += fmt.Sprintf("&database_name=%v", v.(string))
	}
	if v, ok := d.GetOk("table_space_name"); ok {
		res += fmt.Sprintf("&table_space_name=%v", v.(string))
	}
	if v, ok := d.GetOk("user_name"); ok {
		res += fmt.Sprintf("&user_name=%v", v.(string))
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenDatabaseStorageUsageBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("database_volumes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"database_name":    utils.PathSearch("database_name", v, nil),
			"table_space_name": utils.PathSearch("table_space_name", v, nil),
			"user_name":        utils.PathSearch("user_name", v, nil),
			"database_size":    utils.PathSearch("database_size", v, nil),
		})
	}
	return rst
}
