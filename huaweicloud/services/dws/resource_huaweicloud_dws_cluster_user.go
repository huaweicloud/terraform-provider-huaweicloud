package dws

import (
	"context"
	"fmt"
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

var clusterUserNonUpdatableParams = []string{
	"cluster_id",
	"name",
	"type",
	"password",
	"description",
	"password_disable",
	"logical_cluster",
	"grant_list",
}

// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/db-manager/users
// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/db-manager/users/{name}
// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/db-manager/users/{name}/authority
// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/db-manager/users/{name}
// @API DWS DELETE /v1/{project_id}/clusters/{cluster_id}/db-manager/users/{name}
func ResourceClusterUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterUserCreate,
		ReadContext:   resourceClusterUserRead,
		UpdateContext: resourceClusterUserUpdate,
		DeleteContext: resourceClusterUserDelete,

		CustomizeDiff: config.FlexibleForceNew(clusterUserNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceClusterUserImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the cluster user is located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the DWS cluster.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The object type.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the cluster user or role.`,
			},

			// Optional parameters.
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The password of the user.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the user or role.`,
			},
			"password_disable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to disable password authentication.`,
			},
			"logical_cluster": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The logical cluster name.`,
			},
			"grant_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        clusterUserGrantListElemSchema(),
				Description: `The list of grants.`,
			},
			"cascade": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to cascade delete the dependencies when deleting the user or role.`,
			},
			"login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to allow the user to log in.`,
			},
			"create_role": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to grant the permission to create roles.`,
			},
			"create_db": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to grant the permission to create databases.`,
			},
			"system_admin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to grant the system administrator permission.`,
			},
			"audit_admin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to grant the audit administrator permission.`,
			},
			"inherit": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to inherit permissions from roles.`,
			},
			"use_ft": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to grant the external table permission.`,
			},
			"conn_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The maximum number of concurrent connections.`,
			},
			"replication": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to grant the replication permission.`,
			},
			"valid_begin": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The valid begin time.`,
			},
			"valid_until": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The valid until time.`,
			},
			"lock": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether the user is locked.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func clusterUserGrantListElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The object type.`,
			},
			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The database name.`,
			},
			"schema_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The schema name.`,
			},
			"object_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The object name.`,
			},
			"all_object": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether all objects are included.`,
			},
			"future": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to grant privileges on future objects.`,
			},
			"future_object_owners": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The owners of the future objects.`,
			},
			"column_names": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of column names.`,
			},
			"privileges": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        clusterUserAuthorityPrivilegeSchema(),
				Description: `The list of privileges.`,
			},
		},
	}
}

func clusterUserAuthorityPrivilegeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"permission": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The privilege name.`,
			},
			"grant_with": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether the grant option is included.`,
			},
		},
	}
}

func buildClusterUserAuthorityPrivileges(privileges []interface{}) []map[string]interface{} {
	if len(privileges) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(privileges))
	for _, privilege := range privileges {
		result = append(result, map[string]interface{}{
			"permission": utils.PathSearch("permission", privilege, nil),
			"grant_with": utils.PathSearch("grant_with", privilege, nil),
		})
	}

	return result
}

func buildClusterUserGrantList(grants []interface{}) []map[string]interface{} {
	if len(grants) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(grants))
	for _, grant := range grants {
		result = append(result, map[string]interface{}{
			"type":                 utils.PathSearch("type", grant, nil),
			"database":             utils.ValueIgnoreEmpty(utils.PathSearch("database", grant, "")),
			"schema":               utils.ValueIgnoreEmpty(utils.PathSearch("schema_name", grant, "")),
			"obj_name":             utils.ValueIgnoreEmpty(utils.PathSearch("object_name", grant, "")),
			"all_object":           utils.ValueIgnoreEmpty(utils.PathSearch("all_object", grant, nil)),
			"future":               utils.ValueIgnoreEmpty(utils.PathSearch("future", grant, nil)),
			"future_object_owners": utils.ValueIgnoreEmpty(utils.PathSearch("future_object_owners", grant, "")),
			"column_name": utils.ValueIgnoreEmpty(
				utils.ExpandToStringListBySet(utils.PathSearch("column_names", grant, &schema.Set{}).(*schema.Set)),
			),
			"privileges": buildClusterUserAuthorityPrivileges(
				utils.PathSearch("privileges", grant, &schema.Set{}).(*schema.Set).List(),
			),
		})
	}

	return result
}

func buildClusterUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"type":             d.Get("type"),
		"name":             d.Get("name"),
		"password":         utils.ValueIgnoreEmpty(d.Get("password")),
		"desc":             utils.ValueIgnoreEmpty(d.Get("description")),
		"password_disable": utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "password_disable"),
		"logical_cluster":  utils.ValueIgnoreEmpty(d.Get("logical_cluster")),
		"grant_list":       utils.ValueIgnoreEmpty(buildClusterUserGrantList(d.Get("grant_list").(*schema.Set).List())),
		"login":            utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "login"),
		"create_role":      utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "create_role"),
		"create_db":        utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "create_db"),
		"system_admin":     utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "system_admin"),
		"inherit":          utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "inherit"),
		"conn_limit":       utils.ValueIgnoreEmpty(utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "conn_limit")),
	}
}

func buildClusterUserUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"login":       utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "login"),
		"createrole":  utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "create_role"),
		"createdb":    utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "create_db"),
		"systemadmin": utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "system_admin"),
		"auditadmin":  utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "audit_admin"),
		"inherit":     utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "inherit"),
		"useft":       utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "use_ft"),
		"conn_limit":  utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "conn_limit"),
		"replication": utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "replication"),
		"valid_begin": utils.ValueIgnoreEmpty(d.Get("valid_begin")),
		"valid_until": utils.ValueIgnoreEmpty(d.Get("valid_until")),
		"lock":        utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "lock"),
	}
}

// Some parameters (such as createdb, systemadmin, etc.) are not supported by the Create API and must be
// initialized through the Update API after the user is created.
func resourceClusterUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		clusterId     = d.Get("cluster_id").(string)
		name          = d.Get("name").(string)
		createHttpUrl = "v1/{project_id}/clusters/{cluster_id}/db-manager/users"
		updateHttpUrl = "v1/{project_id}/clusters/{cluster_id}/db-manager/users/{name}"
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", clusterId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
		JSONBody:         utils.RemoveNil(buildClusterUserBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating cluster user: %s", err)
	}

	id := fmt.Sprintf("%s/%s/%s", region, clusterId, name)
	d.SetId(id)

	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", clusterId)
	updatePath = strings.ReplaceAll(updatePath, "{name}", name)

	updateBody := utils.RemoveNil(buildClusterUserUpdateBodyParams(d))
	if len(updateBody) < 1 {
		return resourceClusterUserRead(ctx, d, meta)
	}

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
		JSONBody:         updateBody,
	}

	_, err = client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating cluster user after create: %s", err)
	}

	return resourceClusterUserRead(ctx, d, meta)
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

func resourceClusterUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	userDetail, err := GetClusterUser(client, clusterId, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving cluster user")
	}
	authorities, err := listClusterUserAuthorities(client, clusterId, name)
	if err != nil {
		return diag.Errorf("error querying cluster user authorities: %s", err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("description", utils.PathSearch("desc", userDetail, nil)),
		d.Set("logical_cluster", utils.PathSearch("logical_cluster", userDetail, nil)),
		d.Set("login", utils.PathSearch("login", userDetail, nil)),
		d.Set("create_role", utils.PathSearch("createrole", userDetail, nil)),
		d.Set("create_db", utils.PathSearch("createdb", userDetail, nil)),
		d.Set("system_admin", utils.PathSearch("systemadmin", userDetail, nil)),
		d.Set("audit_admin", utils.PathSearch("auditadmin", userDetail, nil)),
		d.Set("inherit", utils.PathSearch("inherit", userDetail, nil)),
		d.Set("use_ft", utils.PathSearch("useft", userDetail, nil)),
		d.Set("conn_limit", utils.PathSearch("conn_limit", userDetail, nil)),
		d.Set("replication", utils.PathSearch("replication", userDetail, nil)),
		d.Set("valid_begin", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("valid_begin", userDetail, float64(0)).(float64))/1000, false)),
		d.Set("valid_until", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("valid_until", userDetail, float64(0)).(float64))/1000, false)),
		d.Set("lock", utils.PathSearch("lock", userDetail, nil)),
		d.Set("grant_list", flattenClusterUserAuthorities(authorities)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateClusterUser(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl   = "v1/{project_id}/clusters/{cluster_id}/db-manager/users/{name}"
		clusterId = d.Get("cluster_id").(string)
		name      = d.Get("name").(string)
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", clusterId)
	updatePath = strings.ReplaceAll(updatePath, "{name}", name)

	updateBody := utils.RemoveNil(buildClusterUserUpdateBodyParams(d))
	if len(updateBody) < 1 {
		return nil
	}

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
		JSONBody:         updateBody,
	}

	_, err := client.Request("POST", updatePath, &updateOpt)
	return err
}

func resourceClusterUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if d.HasChangeExcept("enable_force_new") {
		err = updateClusterUser(client, d)
		if err != nil {
			return diag.Errorf("error updating cluster user: %s", err)
		}
	}

	return resourceClusterUserRead(ctx, d, meta)
}

func resourceClusterUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
		name      = d.Get("name").(string)
		cascade   = d.Get("cascade").(bool)
		httpUrl   = "v1/{project_id}/clusters/{cluster_id}/db-manager/users/{name}"
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", clusterId)
	deletePath = strings.ReplaceAll(deletePath, "{name}", name)
	deletePathWithQuery := fmt.Sprintf("%s?cascade=%t", deletePath, cascade)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	_, err = client.Request("DELETE", deletePathWithQuery, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting cluster user")
	}

	return nil
}

// GetClusterUser is a method used to query the database user or role details from the DWS API.
// It returns the user details as an interface{} and any error encountered.
func GetClusterUser(client *golangsdk.ServiceClient, clusterId, name string) (interface{}, error) {
	httpUrl := "v1/{project_id}/clusters/{cluster_id}/db-manager/users/{name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)
	getPath = strings.ReplaceAll(getPath, "{name}", name)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceClusterUserImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for resource ID, want '<region>/<cluster_id>/<name>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(
		d.Set("region", parts[0]),
		d.Set("cluster_id", parts[1]),
		d.Set("name", parts[2]),
	)
	d.SetId(fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2]))

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
