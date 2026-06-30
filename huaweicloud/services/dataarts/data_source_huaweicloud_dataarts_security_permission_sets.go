package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/security/permission-sets
func DataSourceSecurityPermissionSets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityPermissionSetsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the permission sets are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the permission sets belong.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the permission set to be queried.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The parent ID of the permission set to be queried.`,
			},
			"type_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type filter of the permission sets to be queried.`,
			},
			"manager_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The manager ID of the permission set to be queried.`,
			},
			"manager_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The manager name of the permission set to be queried.`,
			},
			"manager_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The manager type of the permission set to be queried.`,
			},
			"datasource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The datasource type of the permission set to be queried.`,
			},
			"sync_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sync status of the permission set to be queried.`,
			},

			// Attributes.
			"permission_sets": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        securityPermissionSetSchema(),
				Description: `The list of permission sets that match the filter parameters.`,
			},
		},
	}
}

func securityPermissionSetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the permission set.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent ID of the permission set.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the permission set.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the permission set.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the permission set.`,
			},
			"managed_cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The managed cluster ID of the permission set.`,
			},
			"managed_cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The managed cluster name of the permission set.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID of the permission set.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain ID of the permission set.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance ID of the permission set.`,
			},
			"manager_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The manager ID of the permission set.`,
			},
			"manager_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The manager name of the permission set.`,
			},
			"manager_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The manager type of the permission set.`,
			},
			"datasource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The datasource type of the permission set.`,
			},
			"sync_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sync status of the permission set.`,
			},
			"sync_msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sync message of the permission set.`,
			},
			"sync_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sync time of the permission set, in RFC3339 format.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the permission set, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the permission set, in RFC3339 format.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the permission set.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the permission set.`,
			},
		},
	}
}

func buildSecurityPermissionSetsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("parent_id"); ok {
		res = fmt.Sprintf("%s&parent_id=%v", res, v)
	}
	if v, ok := d.GetOk("type_filter"); ok {
		res = fmt.Sprintf("%s&type_filter=%v", res, v)
	}
	if v, ok := d.GetOk("manager_id"); ok {
		res = fmt.Sprintf("%s&manager_id=%v", res, v)
	}
	if v, ok := d.GetOk("manager_name"); ok {
		res = fmt.Sprintf("%s&manager_name=%v", res, v)
	}
	if v, ok := d.GetOk("manager_type"); ok {
		res = fmt.Sprintf("%s&manager_type=%v", res, v)
	}
	if v, ok := d.GetOk("datasource_type"); ok {
		res = fmt.Sprintf("%s&datasource_type=%v", res, v)
	}
	if v, ok := d.GetOk("sync_status"); ok {
		res = fmt.Sprintf("%s&sync_status=%v", res, v)
	}

	return res
}

func listSecurityPermissionSets(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl     = "v1/{project_id}/security/permission-sets?limit={limit}"
		limit       = 100
		offset      = 0
		workspaceId = d.Get("workspace_id").(string)
		result      = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildSecurityPermissionSetsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	for {
		listPathWithOffset := listPathWithLimit + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		permissionSets := utils.PathSearch("permission_sets", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, permissionSets...)
		if len(permissionSets) < limit {
			break
		}
		offset += len(permissionSets)
	}

	return result, nil
}

func flattenSecurityPermissionSets(permissionSets []interface{}) []map[string]interface{} {
	if len(permissionSets) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(permissionSets))
	for _, permissionSet := range permissionSets {
		result = append(result, map[string]interface{}{
			"id":                   utils.PathSearch("id", permissionSet, nil),
			"parent_id":            utils.PathSearch("parent_id", permissionSet, nil),
			"name":                 utils.PathSearch("name", permissionSet, nil),
			"description":          utils.PathSearch("description", permissionSet, nil),
			"type":                 utils.PathSearch("type", permissionSet, nil),
			"managed_cluster_id":   utils.PathSearch("managed_cluster_id", permissionSet, nil),
			"managed_cluster_name": utils.PathSearch("managed_cluster_name", permissionSet, nil),
			"project_id":           utils.PathSearch("project_id", permissionSet, nil),
			"domain_id":            utils.PathSearch("domain_id", permissionSet, nil),
			"instance_id":          utils.PathSearch("instance_id", permissionSet, nil),
			"manager_id":           utils.PathSearch("manager_id", permissionSet, nil),
			"manager_name":         utils.PathSearch("manager_name", permissionSet, nil),
			"manager_type":         utils.PathSearch("manager_type", permissionSet, nil),
			"datasource_type":      utils.PathSearch("datasource_type", permissionSet, nil),
			"sync_status":          utils.PathSearch("sync_status", permissionSet, nil),
			"sync_msg":             utils.PathSearch("sync_msg", permissionSet, nil),
			"sync_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("sync_time",
				permissionSet, float64(0)).(float64))/1000, false),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				permissionSet, float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time",
				permissionSet, float64(0)).(float64))/1000, false),
			"created_by": utils.PathSearch("created_by", permissionSet, nil),
			"updated_by": utils.PathSearch("updated_by", permissionSet, nil),
		})
	}

	return result
}

func dataSourceSecurityPermissionSetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	permissionSets, err := listSecurityPermissionSets(client, d)
	if err != nil {
		return diag.Errorf("error querying permission sets: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("permission_sets", flattenSecurityPermissionSets(permissionSets)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
