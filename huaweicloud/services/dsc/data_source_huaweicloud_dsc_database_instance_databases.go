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

// @API DSC GET /v1/{project_id}/asset-center/database/instances/{instance_id}/databases
func DataSourceDscDatabaseInstanceDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscDatabaseInstanceDatabasesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the database instance ID.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the database type.",
			},
			"databases": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The database information list.",
				Elem:        dscDatabaseInstanceDatabaseSchema(),
			},
		},
	}
}

func dscDatabaseInstanceDatabaseSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database ID.",
			},
			"db_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The database port.",
			},
			"db_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database name.",
			},
			"asset_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset name.",
			},
			"authorized": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the database is authorized.",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the database is the default database.",
			},
		},
	}
}

func buildDscDatabaseInstanceDatabasesQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}
	return queryParams
}

func dataSourceDscDatabaseInstanceDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "dsc"
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v1/{project_id}/asset-center/database/instances/{instance_id}/databases"
		offset     = 0
		limit      = 1000
		result     = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildDscDatabaseInstanceDatabasesQueryParams(d, limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC database instance databases: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		databasesResp := utils.PathSearch("databases", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, databasesResp...)

		if len(databasesResp) < limit {
			break
		}

		offset += len(databasesResp)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("databases", flattenDscDatabaseInstanceDatabases(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscDatabaseInstanceDatabases(databases []interface{}) []interface{} {
	if len(databases) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(databases))
	for _, v := range databases {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"db_port":    utils.PathSearch("db_port", v, nil),
			"db_name":    utils.PathSearch("db_name", v, nil),
			"asset_name": utils.PathSearch("asset_name", v, nil),
			"authorized": utils.PathSearch("authorized", v, nil),
			"default":    utils.PathSearch("default", v, nil),
		})
	}

	return rst
}
