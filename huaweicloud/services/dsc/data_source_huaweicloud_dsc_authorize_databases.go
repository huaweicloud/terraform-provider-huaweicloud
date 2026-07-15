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

// @API DSC GET /v1/{project_id}/asset-center/database/authorized-databases
func DataSourceDscAuthorizeDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscAuthorizeDatabasesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the instance type.",
			},
			"databases": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The authorized database information list.",
				Elem:        dscAuthorizeDatabaseSchema(),
			},
		},
	}
}

func dscAuthorizeDatabaseSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database ID.",
			},
			"asset_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset name.",
			},
			"auth_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The authorization type.",
			},
			"authorize_fail_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The authorization failure reason.",
			},
			"authorized": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the database is authorized.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The creation time.",
			},
			"db_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database address.",
			},
			"db_authorized": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database authorization status.",
			},
			"db_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database name.",
			},
			"db_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The database port.",
			},
			"db_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database type.",
			},
			"db_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database username.",
			},
			"db_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database version.",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether it is the default database.",
			},
			"ins_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance ID.",
			},
			"ins_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance name.",
			},
			"ins_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance type.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region where the instance is located.",
			},
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The service name.",
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The session ID.",
			},
			"subnet_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The subnet ID list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vpc_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC ID list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildDscAuthorizeDatabasesQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?instance_type=%s&limit=%d&offset=%d", d.Get("instance_type").(string), limit, offset)

	return queryParams
}

func dataSourceDscAuthorizeDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/asset-center/database/authorized-databases"
		offset  = 0
		limit   = 1000
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildDscAuthorizeDatabasesQueryParams(d, limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC authorized databases: %s", err)
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
		d.Set("databases", flattenDscAuthorizeDatabases(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscAuthorizeDatabases(databases []interface{}) []interface{} {
	if len(databases) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(databases))
	for _, v := range databases {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"asset_name":            utils.PathSearch("asset_name", v, nil),
			"auth_type":             utils.PathSearch("auth_type", v, nil),
			"authorize_fail_reason": utils.PathSearch("authorize_fail_reason", v, nil),
			"authorized":            utils.PathSearch("authorized", v, nil),
			"create_time":           utils.PathSearch("create_time", v, nil),
			"db_address":            utils.PathSearch("db_address", v, nil),
			"db_authorized":         utils.PathSearch("db_authorized", v, nil),
			"db_name":               utils.PathSearch("db_name", v, nil),
			"db_port":               utils.PathSearch("db_port", v, nil),
			"db_type":               utils.PathSearch("db_type", v, nil),
			"db_user":               utils.PathSearch("db_user", v, nil),
			"db_version":            utils.PathSearch("db_version", v, nil),
			"default":               utils.PathSearch("default", v, nil),
			"ins_id":                utils.PathSearch("ins_id", v, nil),
			"ins_name":              utils.PathSearch("ins_name", v, nil),
			"ins_type":              utils.PathSearch("ins_type", v, nil),
			"region":                utils.PathSearch("region", v, nil),
			"service_name":          utils.PathSearch("service_name", v, nil),
			"sid":                   utils.PathSearch("sid", v, nil),
			"subnet_ids":            utils.PathSearch("subnet_ids", v, nil),
			"vpc_ids":               utils.PathSearch("vpc_ids", v, nil),
		})
	}

	return rst
}
