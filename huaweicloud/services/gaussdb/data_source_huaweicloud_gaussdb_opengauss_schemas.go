package gaussdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/schemas
func DataSourceOpenGaussSchemas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenGaussSchemasRead,

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
			"database_schemas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schema_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOpenGaussSchemasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/schemas"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listBasePath := client.Endpoint + httpUrl
	listBasePath = strings.ReplaceAll(listBasePath, "{project_id}", client.ProjectID)
	listBasePath = strings.ReplaceAll(listBasePath, "{instance_id}", d.Get("instance_id").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	var res []interface{}
	dnName := d.Get("db_name").(string)

	for {
		listPath := listBasePath + buildOpenGaussSchemasQueryParams(dnName, offset)
		listResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return diag.FromErr(err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}
		schemas := flattenListOpenGaussSchemasResponseBody(listRespBody)
		res = append(res, schemas...)
		totalCount := utils.PathSearch("total_count", listRespBody, float64(0)).(float64)
		if int(totalCount) <= (offset+1)*100 {
			break
		}
		offset++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("database_schemas", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildOpenGaussSchemasQueryParams(dbName string, page int) string {
	return fmt.Sprintf("?db_name=%s&limit=100&offset=%v", dbName, page)
}

func flattenListOpenGaussSchemasResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("database_schemas", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"schema_name": utils.PathSearch("schema_name", v, nil),
			"owner":       utils.PathSearch("owner", v, nil),
		})
	}
	return rst
}
