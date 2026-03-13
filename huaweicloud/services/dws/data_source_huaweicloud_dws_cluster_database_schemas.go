package dws

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/databases/{database_name}/schemas
func DataSourceClusterDatabaseSchemas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterDatabaseSchemasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the database schemas are located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster to be queried.`,
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the database to be queried.`,
			},

			// Optional parameters.
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The field used to sort query results.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The direction of sorting query results.`,
			},
			"keywords": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The keywords used for fuzzy query.`,
			},

			// Attributes.
			"schemas": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        clusterDatabaseSchemasSchema(),
				Description: `The list of the database schemas that matched filter parameters.`,
			},
		},
	}
}

func clusterDatabaseSchemasSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"schema_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the schema.`,
			},
			"database_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the database to which the schema belongs.`,
			},
			"total_value": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total used space value of the schema.`,
			},
			"perm_space": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The space threshold of the schema.`,
			},
			"skew_percent": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The skew percentage of the schema.`,
			},
			"min_value": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum value.`,
			},
			"max_value": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum value.`,
			},
			"min_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The minimum DN node.`,
			},
			"max_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The maximum DN node.`,
			},
			"dn_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of DN nodes.`,
			},
		},
	}
}

func buildClusterDatabaseSchemasQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("sort_key"); ok {
		res = fmt.Sprintf("%s&sort_key=%v", res, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		res = fmt.Sprintf("%s&sort_dir=%v", res, v)
	}
	if v, ok := d.GetOk("keywords"); ok {
		res = fmt.Sprintf("%s&keywords=%v", res, v)
	}

	return res
}

func listClusterDatabaseSchemas(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl      = "v2/{project_id}/clusters/{cluster_id}/databases/{database_name}/schemas?limit={limit}"
		clusterID    = d.Get("cluster_id").(string)
		databaseName = d.Get("database_name").(string)
		limit        = 100
		offset       = 0
		result       = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{cluster_id}", clusterID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{database_name}", databaseName)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildClusterDatabaseSchemasQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPathWithLimit + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		schemas := utils.PathSearch("schemas", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, schemas...)
		if len(schemas) < limit {
			break
		}

		offset += len(schemas)
	}

	return result, nil
}

func flattenClusterDatabaseSchemas(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"schema_name":   utils.PathSearch("schema_name", item, nil),
			"database_name": utils.PathSearch("database_name", item, nil),
			"total_value":   utils.PathSearch("total_value", item, nil),
			"perm_space":    utils.PathSearch("perm_space", item, nil),
			"skew_percent":  utils.PathSearch("skew_percent", item, nil),
			"min_value":     utils.PathSearch("min_value", item, nil),
			"max_value":     utils.PathSearch("max_value", item, nil),
			"min_dn":        utils.PathSearch("min_dn", item, nil),
			"max_dn":        utils.PathSearch("max_dn", item, nil),
			"dn_num":        utils.PathSearch("dn_num", item, nil),
		})
	}

	return result
}

func dataSourceClusterDatabaseSchemasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	schemas, err := listClusterDatabaseSchemas(client, d)
	if err != nil {
		return diag.Errorf("error querying cluster database schemas: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("schemas", flattenClusterDatabaseSchemas(schemas)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
