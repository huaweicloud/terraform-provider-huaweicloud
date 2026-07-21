package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	permissionSetPrivilegeResourceNotFoundCodes = []string{
		"DLS.6036",
		"DLS.3027",
	}
	permissionSetPrivilegeNonUpdatableParams = []string{
		"workspace_id",
		"permission_set_id",
		"datasource_type",
		"type",
		"cluster_name",
		"cluster_id",
		"database_url",
		"database_name",
		"table_name",
		"column_name",
		"schema_name",
	}
)

// @API DataArtsStudio POST /v1/{project_id}/security/permission-sets/{permission_set_id}/permissions
// @API DataArtsStudio GET /v1/{project_id}/security/permission-sets/{permission_set_id}/permissions
// @API DataArtsStudio PUT /v1/{project_id}/security/permission-sets/{permission_set_id}/permissions/{permission_id}
// @API DataArtsStudio POST /v1/{project_id}/security/permission-sets/{permission_set_id}/permissions/batch-delete
func ResourceSecurityPermissionSetPrivilege() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPermissionSetPrivilegeCreate,
		ReadContext:   resourceSecurityPermissionSetPrivilegeRead,
		UpdateContext: resourceSecurityPermissionSetPrivilegeUpdate,
		DeleteContext: resourceSecurityPermissionSetPrivilegeDelete,

		CustomizeDiff: config.FlexibleForceNew(permissionSetPrivilegeNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityPermissionSetPrivilegeImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which to the permission set to be granted privilege is located.`,
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
				Description: `The ID of the permission set to be granted.`,
			},
			"datasource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the granted data source.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the permission to be configured.`,
			},
			"actions": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of granted permissions.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The cluster name corresponding to the granted data source.`,
			},

			// Optional parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The cluster ID corresponding to the granted data source.`,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The data connection ID corresponding to the granted data source.`,
			},
			"database_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"database_name", "table_name", "column_name"},
				Description:   `The URL of the database corresponding to the granted data source.`,
			},
			"database_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"database_url"},
				Description:   `The name of the database corresponding to the granted data source.`,
			},
			"table_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"database_url"},
				Description:   `The name of the data table corresponding to the granted data source.`,
			},
			"column_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"database_url"},
				Description:   `The name of the column corresponding to the granted data source.`,
			},
			"schema_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The schema name corresponding to the DWS data source.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current synchronization status of the resource.`,
			},
			"sync_msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status information of the resource.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildSecurityMoreHeaders(workspaceId string) map[string]string {
	moreHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	if workspaceId != "" {
		moreHeaders["workspace"] = workspaceId
	}

	return moreHeaders
}

func associatePrivilegeToPermissionSet(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v1/{project_id}/security/permission-sets/{permission_set_id}/permissions"
	associatePath := client.Endpoint + httpUrl
	associatePath = strings.ReplaceAll(associatePath, "{project_id}", client.ProjectID)
	associatePath = strings.ReplaceAll(associatePath, "{permission_set_id}", d.Get("permission_set_id").(string))

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody:         utils.RemoveNil(buildCreatePermissionSetPrivilegeBodyParams(d)),
	}

	requestResp, err := client.Request("POST", associatePath, &opts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceSecurityPermissionSetPrivilegeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := associatePrivilegeToPermissionSet(client, d)
	if err != nil {
		return diag.Errorf("error creating DataArts Security permission set privilege: %s", err)
	}
	privilegeId := utils.PathSearch("id", respBody, "").(string)
	if privilegeId == "" {
		return diag.Errorf("unable to find the privilege ID of the DataArts Security permission set from the API response")
	}
	d.SetId(privilegeId)

	return resourceSecurityPermissionSetPrivilegeRead(ctx, d, meta)
}

func buildCreatePermissionSetPrivilegeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"datasource_type":    d.Get("datasource_type"),
		"permission_type":    d.Get("type"),
		"permission_actions": d.Get("actions").(*schema.Set).List(),
		"cluster_name":       d.Get("cluster_name"),
		// Optional parameters.
		"cluster_id":    utils.ValueIgnoreEmpty(d.Get("cluster_id")),
		"dw_id":         utils.ValueIgnoreEmpty(d.Get("connection_id")),
		"url":           utils.ValueIgnoreEmpty(d.Get("database_url")),
		"database_name": utils.ValueIgnoreEmpty(d.Get("database_name")),
		"table_name":    utils.ValueIgnoreEmpty(d.Get("table_name")),
		"column_name":   utils.ValueIgnoreEmpty(d.Get("column_name")),
		"schema_name":   utils.ValueIgnoreEmpty(d.Get("schema_name")),
	}
}

func listPermissionSetAssociatedPrivileges(client *golangsdk.ServiceClient, workspaceId, permissionSetId string,
	queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/security/permission-sets/{permission_set_id}/permissions?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{permission_set_id}", permissionSetId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(queryParams) > 0 {
		listPath = fmt.Sprintf("%s&%s", listPath, queryParams[0])
	}

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &opts)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}
		privileges := utils.PathSearch("permissions", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, privileges...)
		if len(privileges) < limit {
			break
		}
		offset += len(privileges)
	}
	return result, nil
}

// GetPrivilegeById is a method used to query permission configuration using a specified ID.
func GetPrivilegeById(client *golangsdk.ServiceClient, workspaceId, permissionSetId, privilegeId string) (interface{}, error) {
	privileges, err := listPermissionSetAssociatedPrivileges(client, workspaceId, permissionSetId)
	if err != nil {
		return nil, err
	}

	privilege := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", privilegeId), privileges, nil)
	if privilege == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/security/permission-sets/{permission_set_id}/permissions",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the privilege (%s) does not exist", privilegeId)),
			},
		}
	}
	return privilege, nil
}

func resourceSecurityPermissionSetPrivilegeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := GetPrivilegeById(client, d.Get("workspace_id").(string), d.Get("permission_set_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", permissionSetPrivilegeResourceNotFoundCodes...),
			"DataArts Security permission set privilege")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("datasource_type", utils.PathSearch("datasource_type", respBody, nil)),
		d.Set("type", utils.PathSearch("permission_type", respBody, nil)),
		d.Set("actions", utils.PathSearch("permission_actions", respBody, nil)),
		d.Set("cluster_name", utils.PathSearch("cluster_name", respBody, nil)),
		d.Set("cluster_id", utils.PathSearch("cluster_id", respBody, nil)),
		d.Set("database_url", utils.PathSearch("url", respBody, nil)),
		d.Set("database_name", utils.PathSearch("database_name", respBody, nil)),
		d.Set("table_name", utils.PathSearch("table_name", respBody, nil)),
		d.Set("column_name", utils.PathSearch("column_name", respBody, nil)),
		d.Set("schema_name", utils.PathSearch("schema_name", respBody, nil)),
		d.Set("status", utils.PathSearch("sync_status", respBody, nil)),
		d.Set("sync_msg", utils.PathSearch("sync_msg", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdatePermissionSetPrivilegeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// Required parameters.
		"permission_actions": utils.ExpandToStringList(d.Get("actions").(*schema.Set).List()),
		// Optional parameters.
		"dw_id": utils.ValueIgnoreEmpty(d.Get("connection_id")),
	}
	return bodyParams
}

func updatePermissionSetPrivilege(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl         = "v1/{project_id}/security/permission-sets/{permission_set_id}/permissions/{permission_id}"
		permissionSetId = d.Get("permission_set_id").(string)
		privilegeId     = d.Id()
		workspaceId     = d.Get("workspace_id").(string)
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{permission_set_id}", permissionSetId)
	updatePath = strings.ReplaceAll(updatePath, "{permission_id}", privilegeId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
		JSONBody:         utils.RemoveNil(buildUpdatePermissionSetPrivilegeBodyParams(d)),
	}

	resp, err := client.Request("PUT", updatePath, &opts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func resourceSecurityPermissionSetPrivilegeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if d.HasChangesExcept("enable_force_new") {
		_, err = updatePermissionSetPrivilege(client, d)
		if err != nil {
			return diag.Errorf("error updating DataArts Security permission set privilege: %s", err)
		}
	}

	return resourceSecurityPermissionSetPrivilegeRead(ctx, d, meta)
}

func deletePermissionSetPrivilege(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl         = "v1/{project_id}/security/permission-sets/{permission_set_id}/permissions/batch-delete"
		permissionSetId = d.Get("permission_set_id").(string)
		privilegeId     = d.Id()
		workspaceId     = d.Get("workspace_id").(string)
	)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{permission_set_id}", permissionSetId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
		JSONBody: map[string]interface{}{
			"ids": []string{privilegeId},
		},
		OkCodes: []int{200, 204},
	}

	_, err := client.Request("POST", deletePath, &opts)
	return err
}

func resourceSecurityPermissionSetPrivilegeDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	err = deletePermissionSetPrivilege(client, d)
	if err != nil {
		return diag.Errorf("error deleting DataArts Security permission set privilege: %s", err)
	}

	return nil
}

func resourceSecurityPermissionSetPrivilegeImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<permission_set_id>/<id>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(
		d.Set("workspace_id", parts[0]),
		d.Set("permission_set_id", parts[1]),
	)
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
