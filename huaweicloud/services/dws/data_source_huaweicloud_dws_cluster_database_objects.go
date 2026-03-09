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

// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/db-manager/objects
func DataSourceClusterDatabaseObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterDatabaseObjectsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the database objects are located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the database object.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the database object.`,
			},
			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the database.`,
			},
			"schema": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the schema.`,
			},
			"table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the table.`,
			},
			"is_fine_grained_disaster": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Whether fine-grained disaster recovery is enabled.`,
			},

			// Attributes.
			"objects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"obj_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the database object.`,
						},
					},
				},
				Description: `The list of the database objects that matched filter parameters.`,
			},
		},
	}
}

func buildClusterDatabaseObjectsQueryParams(d *schema.ResourceData) string {
	res := ""

	res = fmt.Sprintf("%s&type=%v", res, d.Get("type"))

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("database"); ok {
		res = fmt.Sprintf("%s&database=%v", res, v)
	}
	if v, ok := d.GetOk("schema"); ok {
		res = fmt.Sprintf("%s&schema=%v", res, v)
	}
	if v, ok := d.GetOk("table"); ok {
		res = fmt.Sprintf("%s&table=%v", res, v)
	}
	if v, ok := d.GetOk("is_fine_grained_disaster"); ok {
		res = fmt.Sprintf("%s&is_fine_grained_disaster=%v", res, v)
	}

	return res
}

func listClusterDatabaseObjects(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl   = "v1/{project_id}/clusters/{cluster_id}/db-manager/objects?limit={limit}"
		clusterId = d.Get("cluster_id").(string)
		limit     = 1000
		offset    = 0
		result    = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", clusterId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildClusterDatabaseObjectsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		objects := utils.PathSearch("object_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, objects...)
		if len(objects) < limit {
			break
		}
		offset += len(objects)
	}

	return result, nil
}

func flattenClusterDatabaseObjects(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"obj_name": utils.PathSearch("obj_name", item, nil),
		})
	}
	return result
}

func dataSourceClusterDatabaseObjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	objects, err := listClusterDatabaseObjects(client, d)
	if err != nil {
		return diag.Errorf("error querying cluster database objects: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("objects", flattenClusterDatabaseObjects(objects)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
