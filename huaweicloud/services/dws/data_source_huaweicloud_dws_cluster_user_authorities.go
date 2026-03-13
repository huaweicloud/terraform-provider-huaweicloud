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

// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/db-manager/users/{name}/authority
func DataSourceClusterUserAuthorities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterUserAuthoritiesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the cluster user authorities are located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the DWS cluster to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user name or role name to be queried.`,
			},

			// Attributes.
			"authorities": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of user or role authorities.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The authority type.`,
						},
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The database name.`,
						},
						"schema_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The schema name.`,
						},
						"object_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The object name.`,
						},
						"all_object": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether all objects are effective.`,
						},
						"future": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether future objects are effective.`,
						},
						"future_object_owners": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The owners of future objects.`,
						},
						"column_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of column names.`,
						},
						"privileges": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        clusterUserAuthorityPrivilegeSchema(),
							Description: `The privilege list under this authority record.`,
						},
					},
				},
			},
		},
	}
}

func clusterUserAuthorityPrivilegeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"permission": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The privilege name.`,
			},
			"grant_with": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the grant option is included.`,
			},
		},
	}
}

func listClusterUserAuthorities(client *golangsdk.ServiceClient, clusterId, name string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/clusters/{cluster_id}/db-manager/users/{name}/authority?limit={limit}"
		limit   = 1000
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", clusterId)
	listPath = strings.ReplaceAll(listPath, "{name}", name)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		authorities := utils.PathSearch("authority_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, authorities...)
		if len(authorities) < limit {
			break
		}
		offset += len(authorities)
	}

	return result, nil
}

func flattenClusterUserAuthorityPrivileges(privileges []interface{}) []map[string]interface{} {
	if len(privileges) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(privileges))
	for _, privilege := range privileges {
		result = append(result, map[string]interface{}{
			"permission": utils.PathSearch("permission", privilege, nil),
			"grant_with": utils.PathSearch("grant_with", privilege, false),
		})
	}

	return result
}

func flattenClusterUserAuthorities(authorities []interface{}) []map[string]interface{} {
	if len(authorities) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(authorities))
	for _, authority := range authorities {
		result = append(result, map[string]interface{}{
			"type":                 utils.PathSearch("type", authority, nil),
			"database":             utils.PathSearch("database", authority, nil),
			"schema_name":          utils.PathSearch("schema", authority, nil),
			"object_name":          utils.PathSearch("obj_name", authority, nil),
			"all_object":           utils.PathSearch("all_object", authority, false),
			"future":               utils.PathSearch("future", authority, false),
			"future_object_owners": utils.PathSearch("future_object_owners", authority, nil),
			"column_names":         utils.PathSearch("column_name", authority, make([]interface{}, 0)),
			"privileges": flattenClusterUserAuthorityPrivileges(utils.PathSearch("privileges", authority,
				make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceClusterUserAuthoritiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
		name      = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	authorities, err := listClusterUserAuthorities(client, clusterId, name)
	if err != nil {
		return diag.Errorf("error querying cluster user authorities: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("authorities", flattenClusterUserAuthorities(authorities)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
