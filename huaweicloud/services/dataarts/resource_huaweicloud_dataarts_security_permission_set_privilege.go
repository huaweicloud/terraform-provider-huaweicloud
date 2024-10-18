package dataarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var PermissionSetPrivilegeResourceNotFoundCodes = []string{
	"DLS.6036",
	"DLS.3027",
}

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

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityPermissionSetPrivilegeImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission_set_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"datasource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"actions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"database_name", "table_name", "column_name"},
			},
			"database_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"database_url"},
			},
			"table_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"database_url"},
			},
			"column_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"database_url"},
			},
			"schema_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSecurityPermissionSetPrivilegeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/security/permission-sets/{permission_set_id}/permissions"
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	creatPath := client.Endpoint + httpUrl
	creatPath = strings.ReplaceAll(creatPath, "{project_id}", client.ProjectID)
	creatPath = strings.ReplaceAll(creatPath, "{permission_set_id}", d.Get("permission_set_id").(string))
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         utils.RemoveNil(buildCreatePermissionSetPrivilegeBodyParams(d)),
	}
	resp, err := client.Request("POST", creatPath, &opts)
	if err != nil {
		return diag.Errorf("error creating DataArts Security permission set privilege: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	privilegeId := utils.PathSearch("id", respBody, "").(string)
	if privilegeId == "" {
		return diag.Errorf("unable to find the privilege ID of the DataArts Security permission set from the API response")
	}

	d.SetId(privilegeId)
	return resourceSecurityPermissionSetPrivilegeRead(ctx, d, meta)
}

func buildCreatePermissionSetPrivilegeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"datasource_type":    d.Get("datasource_type"),
		"permission_type":    d.Get("type"),
		"permission_actions": utils.ExpandToStringList(d.Get("actions").(*schema.Set).List()),
		"cluster_name":       d.Get("cluster_name"),
		"cluster_id":         utils.ValueIgnoreEmpty(d.Get("cluster_id")),
		"dw_id":              utils.ValueIgnoreEmpty(d.Get("connection_id")),
		"url":                utils.ValueIgnoreEmpty(d.Get("database_url")),
		"database_name":      utils.ValueIgnoreEmpty(d.Get("database_name")),
		"table_name":         utils.ValueIgnoreEmpty(d.Get("table_name")),
		"column_name":        utils.ValueIgnoreEmpty(d.Get("column_name")),
		"schema_name":        utils.ValueIgnoreEmpty(d.Get("schema_name")),
	}
	return bodyParams
}

// GetPrivilegeById is a method used to query permission configuration using a specified ID.
func GetPrivilegeById(client *golangsdk.ServiceClient, workspaceId, permissionSetId, id string) (interface{}, error) {
	httpUrl := "v1/{project_id}/security/permission-sets/{permission_set_id}/permissions?limit=100"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{permission_set_id}", permissionSetId)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceId},
	}

	var currentTotal int
	for {
		path := fmt.Sprintf("%s&offset=%v", getPath, currentTotal)
		resp, err := client.Request("GET", path, &opts)
		if err != nil {
			return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", PermissionSetPrivilegeResourceNotFoundCodes...)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		totalNum := utils.PathSearch("total", respBody, 0)
		if privilegeInfo := utils.PathSearch(fmt.Sprintf("permissions|[?id=='%s']|[0]", id), respBody, nil); privilegeInfo != nil {
			return privilegeInfo, nil
		}

		privileges := utils.PathSearch("permissions", respBody, make([]interface{}, 0)).([]interface{})
		currentTotal += len(privileges)
		// The type of `total` is float64.
		if float64(currentTotal) == totalNum {
			break
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func resourceSecurityPermissionSetPrivilegeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := GetPrivilegeById(client, d.Get("workspace_id").(string), d.Get("permission_set_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", PermissionSetPrivilegeResourceNotFoundCodes...),
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

func resourceSecurityPermissionSetPrivilegeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	httpUrl := "v1/{project_id}/security/permission-sets/{permission_set_id}/permissions/{permission_id}"
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{permission_set_id}", d.Get("permission_set_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{permission_id}", d.Id())

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         utils.RemoveNil(buildUpdatePermissionSetPrivilegeBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &opts)
	if err != nil {
		return diag.Errorf("error updating DataArts Security permission set privilege: %s", err)
	}

	return resourceSecurityPermissionSetPrivilegeRead(ctx, d, meta)
}

func buildUpdatePermissionSetPrivilegeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"dw_id":              utils.ValueIgnoreEmpty(d.Get("connection_id")),
		"permission_actions": utils.ExpandToStringList(d.Get("actions").(*schema.Set).List()),
	}
	return bodyParams
}

func resourceSecurityPermissionSetPrivilegeDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/security/permission-sets/{permission_set_id}/permissions/batch-delete"
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{permission_set_id}", d.Get("permission_set_id").(string))
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody: map[string]interface{}{
			"ids": []string{d.Id()},
		},
		OkCodes: []int{204},
	}

	_, err = client.Request("POST", deletePath, &opts)
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
