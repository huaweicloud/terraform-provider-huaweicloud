package dli

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

// The error code corresponding to ErrCodeConnNotFound is an important sign that the related resource does not exist.
// When the object is enhanced connection and the connection not exist, the API return this error:
// + {"error_code": "DLI.0001", "error_msg": "Connection xxx is not exist"}
const ErrCodeConnNotFound string = "DLI.0001"

// @API DLI PUT /v1.0/{project_id}/authorization
// @API DLI GET /v2.0/{project_id}/datasource/enhanced-connections/{connection_id}/privileges
func ResourceDatasourceConnectionPrivilege() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatasourceConnectionPrivilegeCreate,
		UpdateContext: resourceDatasourceConnectionPrivilegeUpdate,
		ReadContext:   resourceDatasourceConnectionPrivilegeRead,
		DeleteContext: resourceDatasourceConnectionPrivilegeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDatasourceConnectionPrivilegeImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the connection to be granted.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the granted project.`,
			},
			"privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of permissions granted to the connection.`,
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					if _, ok := d.GetOk("privileges"); !ok {
						// If the value of the permission is nil, the Computed behavior will prevent changes.
						// If the value of the permission is not nil and the length of the array is zero, need a
						// DiffSuppress function to prevent changes.
						return true
					}
					return false
				},
			},
		},
	}
}

func buildModifyDatasourceConnectionPrivilegesBodyParams(d *schema.ResourceData, action string) map[string]interface{} {
	// Default values of the privileges parameter is 'BIND_QUEUE'.
	privileges := []string{"BIND_QUEUE"}
	if d.Get("privileges").(*schema.Set).Len() > 0 {
		privileges = utils.ExpandToStringListBySet(d.Get("privileges").(*schema.Set))
	}

	return map[string]interface{}{
		"projectId": d.Get("project_id"),
		"action":    action,
		"privileges": []map[string]interface{}{
			{
				"object":     fmt.Sprintf("edsconnections.%v", d.Get("connection_id")),
				"privileges": privileges,
			},
		},
	}
}

func modifyDatasourceConnectionPrivileges(client *golangsdk.ServiceClient, d *schema.ResourceData, action string) error {
	httpUrl := "v1.0/{project_id}/authorization"

	modifyPath := client.Endpoint + httpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	opts.JSONBody = utils.RemoveNil(buildModifyDatasourceConnectionPrivilegesBodyParams(d, action))
	requestResp, err := client.Request("PUT", modifyPath, &opts)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}
	if !utils.PathSearch("is_success", respBody, "").(bool) {
		return fmt.Errorf("unable to %s the privileges: %s", action,
			utils.PathSearch("message", respBody, "Message Not Found"))
	}
	return nil
}

func resourceDatasourceConnectionPrivilegeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	err = modifyDatasourceConnectionPrivileges(client, d, "grant")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%v/%v", d.Get("connection_id"), d.Get("project_id")))

	return resourceDatasourceConnectionPrivilegeRead(ctx, d, meta)
}

func resourceDatasourceConnectionPrivilegeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{connection_id}/privileges"
	)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connection_id}", d.Get("connection_id").(string))

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", ErrCodeConnNotFound),
			"error retrieving privileges")
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return diag.Errorf("unable to query the privileges: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	privilege := utils.PathSearch(fmt.Sprintf("privileges[?project_id=='%v']|[0]", d.Get("project_id")), respBody, nil)
	if privilege == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving privileges")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("connection_id", utils.PathSearch("connection_id", respBody, nil)),
		d.Set("project_id", utils.PathSearch("project_id", privilege, nil)),
		d.Set("privileges", utils.PathSearch("privileges", privilege, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDatasourceConnectionPrivilegeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	err = modifyDatasourceConnectionPrivileges(client, d, "update")
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDatasourceConnectionPrivilegeRead(ctx, d, meta)
}

func resourceDatasourceConnectionPrivilegeDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	err = modifyDatasourceConnectionPrivileges(client, d, "revoke")
	return diag.FromErr(err)
}

func resourceDatasourceConnectionPrivilegeImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	var (
		importId = d.Id()
		parts    = strings.Split(importId, "/")
	)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid resource ID format for privilege management, want '<connection_id>/<project_id>', but got '%s'", importId)
	}
	mErr := multierror.Append(
		d.Set("connection_id", parts[0]),
		d.Set("project_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
