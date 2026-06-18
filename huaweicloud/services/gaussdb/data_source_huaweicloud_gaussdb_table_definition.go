package gaussdb

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

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/list-table-definition
func DataSourceTableDefinition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTableDefinitionRead,

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
			"table_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table_definitions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     tableDefinitionSchema(),
			},
		},
	}
}

func tableDefinitionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table_definition": {
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

func dataSourceTableDefinitionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/list-table-definition"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath += buildTableDefinitionQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB table definitions: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
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
		d.Set("table_definitions", flattenTableDefinitionBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildTableDefinitionQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?database_name=%v", d.Get("database_name").(string))
	res = fmt.Sprintf("%s&table_name=%v", res, d.Get("table_name").(string))

	if v, ok := d.GetOk("schema_name"); ok {
		res = fmt.Sprintf("%s&schema_name=%v", res, v.(string))
	}
	return res
}

func flattenTableDefinitionBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("table_definitions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"table_definition": utils.PathSearch("table_definition", v, nil),
			"schema_name":      utils.PathSearch("schema_name", v, nil),
		})
	}
	return rst
}
