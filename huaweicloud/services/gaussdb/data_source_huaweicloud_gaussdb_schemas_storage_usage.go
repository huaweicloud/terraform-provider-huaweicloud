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

// @API GaussDB GET /v3.1/{project_id}/instances/{instance_id}/schema-volume
func DataSourceSchemasStorageUsage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSchemasStorageUsageRead,

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
				Required: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schema_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemasStorageUsageSchema(),
			},
		},
	}
}

func schemasStorageUsageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"schema_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSchemasStorageUsageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3.1/{project_id}/instances/{instance_id}/schema-volume"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildSchemasStorageUsageQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GaussDB schemas storage usage: %s", err)
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
		d.Set("schema_volumes", flattenSchemasStorageUsageBody(listRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildSchemasStorageUsageQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?database_name=%v", d.Get("database_name").(string))

	if v, ok := d.GetOk("schema_name"); ok {
		res += fmt.Sprintf("&schema_name=%v", v.(string))
	}
	return res
}

func flattenSchemasStorageUsageBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("schema_volumes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"schema_size": utils.PathSearch("schema_size", v, nil),
			"table_count": utils.PathSearch("table_count", v, nil),
			"user_name":   utils.PathSearch("user_name", v, nil),
			"schema_name": utils.PathSearch("schema_name", v, nil),
		})
	}
	return rst
}
