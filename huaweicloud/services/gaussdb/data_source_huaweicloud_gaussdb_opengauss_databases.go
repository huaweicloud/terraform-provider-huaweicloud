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

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/databases
func DataSourceOpenGaussDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenGaussDatabasesRead,

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
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"compatibility_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOpenGaussDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/databases"
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

	for {
		listPath := listBasePath + buildOpenGaussDatabasesQueryParams(offset)
		listResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return diag.FromErr(err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}
		databases := flattenListOpenGaussDatabasesResponseBody(listRespBody)
		res = append(res, databases...)
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
		d.Set("databases", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildOpenGaussDatabasesQueryParams(page int) string {
	return fmt.Sprintf("?limit=100&offset=%v", page)
}

func flattenListOpenGaussDatabasesResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":               utils.PathSearch("name", v, nil),
			"owner":              utils.PathSearch("owner", v, nil),
			"character_set":      utils.PathSearch("character_set", v, nil),
			"lc_collate":         utils.PathSearch("collate_set", v, nil),
			"size":               utils.PathSearch("size", v, nil),
			"compatibility_type": utils.PathSearch("compatibility_type", v, nil),
		})
	}
	return rst
}
