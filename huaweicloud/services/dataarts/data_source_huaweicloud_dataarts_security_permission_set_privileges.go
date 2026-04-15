package dataarts

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/security/permission-sets/{permission_set_id}/permissions
func DataSourceSecurityPermissionSetPrivileges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityPermissionSetPrivilegesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the permission set privileges are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the permission set belongs.`,
			},
			"permission_set_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the permission set to which the granted privileges belong.`,
			},

			// Optional parameters.
			"privilege_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The privilege type used to filter privileges.`,
			},
			"privilege_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The privilege action used to filter privileges.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The cluster ID used to filter privileges.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The cluster name used to filter privileges.`,
			},
			"datasource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The data source type used to filter privileges.`,
			},
			"database_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The database name used to filter privileges.`,
			},
			"table_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The table name used to filter privileges.`,
			},
			"column_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The column name used to filter privileges.`,
			},
			"sync_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The synchronization status used to filter privileges.`,
			},

			// Attributes.
			"privileges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of privileges that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the granted privilege.`,
						},
						"permission_set_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the permission set to which the granted privilege belongs.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the instance to which the granted privilege belongs.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the permission to be configured.`,
						},
						"actions": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `The list of granted privileges.`,
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the cluster to which the granted privilege belongs.`,
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the cluster to which the granted privilege belongs.`,
						},
						"datasource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the granted data source.`,
						},
						"database_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the database to which the granted privilege belongs.`,
						},
						"schema_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the schema to which the granted privilege belongs.`,
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the namespace to which the granted privilege belongs.`,
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the table to which the granted privilege belongs.`,
						},
						"column_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the column to which the granted privilege belongs.`,
						},
						"row_level_security": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The row level security of the granted privilege.`,
						},
						"sync_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The synchronization status of the granted privilege.`,
						},
						"sync_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The synchronization message of the granted privilege.`,
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URL path name of the granted privilege.`,
						},
					},
				},
			},
		},
	}
}

func buildSecurityPermissionSetPrivilegesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("privilege_type"); ok {
		res = fmt.Sprintf("%s&permission_type=%v", res, v)
	}
	if v, ok := d.GetOk("privilege_action"); ok {
		res = fmt.Sprintf("%s&permission_action=%v", res, v)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		res = fmt.Sprintf("%s&cluster_id=%v", res, v)
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		res = fmt.Sprintf("%s&cluster_name=%v", res, v)
	}
	if v, ok := d.GetOk("datasource_type"); ok {
		res = fmt.Sprintf("%s&datasource_type=%v", res, v)
	}
	if v, ok := d.GetOk("database_name"); ok {
		res = fmt.Sprintf("%s&database_name=%v", res, v)
	}
	if v, ok := d.GetOk("table_name"); ok {
		res = fmt.Sprintf("%s&table_name=%v", res, v)
	}
	if v, ok := d.GetOk("column_name"); ok {
		res = fmt.Sprintf("%s&column_name=%v", res, v)
	}
	if v, ok := d.GetOk("sync_status"); ok {
		res = fmt.Sprintf("%s&sync_status=%v", res, v)
	}

	if len(res) < 1 {
		return res
	}
	return res[1:]
}

func flattenSecurityPermissionSetPrivileges(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"id":                 utils.PathSearch("id", item, nil),
			"permission_set_id":  utils.PathSearch("permission_set_id", item, nil),
			"instance_id":        utils.PathSearch("instance_id", item, nil),
			"type":               utils.PathSearch("permission_type", item, nil),
			"actions":            utils.PathSearch("permission_actions", item, nil),
			"cluster_id":         utils.PathSearch("cluster_id", item, nil),
			"cluster_name":       utils.PathSearch("cluster_name", item, nil),
			"datasource_type":    utils.PathSearch("datasource_type", item, nil),
			"database_name":      utils.PathSearch("database_name", item, nil),
			"schema_name":        utils.PathSearch("schema_name", item, nil),
			"namespace":          utils.PathSearch("namespace", item, nil),
			"table_name":         utils.PathSearch("table_name", item, nil),
			"column_name":        utils.PathSearch("column_name", item, nil),
			"row_level_security": utils.PathSearch("row_level_security", item, nil),
			"sync_status":        utils.PathSearch("sync_status", item, nil),
			"sync_msg":           utils.PathSearch("sync_msg", item, nil),
			"url":                utils.PathSearch("url", item, nil),
		}
	}
	return result
}

func dataSourceSecurityPermissionSetPrivilegesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		workspaceId     = d.Get("workspace_id").(string)
		permissionSetId = d.Get("permission_set_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	privileges, err := listPermissionSetAssociatedPrivileges(client, workspaceId, permissionSetId,
		buildSecurityPermissionSetPrivilegesQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying DataArts Security permission set privileges: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("privileges", flattenSecurityPermissionSetPrivileges(privileges)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
