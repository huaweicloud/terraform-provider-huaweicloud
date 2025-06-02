package rds

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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/db-user/database
func DataSourceRdsMysqlAuthorizedDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsMysqlAuthorizedDatabasesRead,
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
			"user_name": {
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
						"readonly": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceRdsMysqlAuthorizedDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/db_user/database"
	basePath := client.Endpoint + httpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	basePath = strings.ReplaceAll(basePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	limit := 100
	page := 1
	userName := d.Get("user_name").(string)
	var databases []interface{}

	for {
		url := fmt.Sprintf("%s?user-name=%s&page=%d&limit=%d", basePath, userName, page, limit)

		resp, err := client.Request("GET", url, &golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		})
		if err != nil {
			return diag.Errorf("error retrieving authorized RDS databases: %s", err)
		}

		body, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		pageDatabases := flattenAuthorizedDatabasesBody(body)

		if len(pageDatabases) == 0 {
			break
		}
		databases = append(databases, pageDatabases...)

		page++
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("databases", databases),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAuthorizedDatabasesBody(resp interface{}) []interface{} {
	databasesJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	databasesArray, ok := databasesJson.([]interface{})
	if !ok || len(databasesArray) == 0 {
		return nil
	}
	res := make([]interface{}, 0, len(databasesArray))
	for _, v := range databasesArray {
		res = append(res, map[string]interface{}{
			"name":     utils.PathSearch("name", v, nil),
			"readonly": utils.PathSearch("readonly", v, nil),
		})
	}
	return res
}
