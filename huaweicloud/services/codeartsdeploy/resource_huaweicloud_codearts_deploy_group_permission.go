package codeartsdeploy

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

// @API CodeArtsDeploy PUT /v2/host-groups/{group_id}/permissions
// @API CodeArtsDeploy GET /v2/host-groups/{group_id}/permissions
func ResourceDeployGroupPermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeployGroupPermissionCreateOrUpdate,
		ReadContext:   resourceDeployGroupPermissionRead,
		UpdateContext: resourceDeployGroupPermissionCreateOrUpdate,
		DeleteContext: resourceDeployGroupPermissionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDeployGroupPermissionImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission_value": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceDeployGroupPermissionCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("codearts_deploy", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	err = modifyDeployGroupPermission(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.IsNewResource() {
		d.SetId(d.Get("group_id").(string) + "/" + d.Get("role_id").(string) + "/" + d.Get("permission_name").(string))
	}

	return resourceDeployGroupPermissionRead(ctx, d, meta)
}

func modifyDeployGroupPermission(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v2/host-groups/{group_id}/permissions"
	modifyPath := client.Endpoint + httpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{group_id}", d.Get("group_id").(string))

	modifyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: buildDeployGroupPermissionBodyParams(d),
	}

	_, err := client.Request("PUT", modifyPath, &modifyOpt)
	if err != nil {
		return fmt.Errorf("error modifying CodeArts deploy group permission: %s", err)
	}

	return nil
}

func buildDeployGroupPermissionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"project_id":       d.Get("project_id"),
		"role_id":          d.Get("role_id"),
		"permission_name":  d.Get("permission_name"),
		"permission_value": d.Get("permission_value"),
	}
}

func resourceDeployGroupPermissionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	permissionMatrix, err := getDeployGroupPermissionMatrix(client, d.Get("group_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts deploy group permission")
	}

	roleId := d.Get("role_id").(string)
	permissionName := d.Get("permission_name").(string)
	expression := fmt.Sprintf("[?role_id=='%s']|[0]", roleId)
	role := utils.PathSearch(expression, permissionMatrix, nil)
	if role == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "unable to find role from API response")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("permission_value", utils.PathSearch(permissionName, role, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDeployGroupPermissionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting permission resource is not supported. The resource is only removed from the state," +
		" the group permission matrix remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceDeployGroupPermissionImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<project_id>/<group_id>/<role_id>/<permission_name>',"+
			" but got '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("project_id", parts[0]),
		d.Set("group_id", parts[1]),
		d.Set("role_id", parts[2]),
		d.Set("permission_name", parts[3]),
	)

	d.SetId(parts[1] + "/" + parts[2] + "/" + parts[3])

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
