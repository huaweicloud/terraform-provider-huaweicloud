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

// @API CodeArtsDeploy PUT /v2/applications/{application_id}/environments/{environment_id}/permissions
// @API CodeArtsDeploy GET /v2/applications/{application_id}/environments/{environment_id}/permissions
func ResourceDeployEnvironmentPermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeployEnvironmentPermissionCreateOrUpdate,
		ReadContext:   resourceDeployEnvironmentPermissionRead,
		UpdateContext: resourceDeployEnvironmentPermissionCreateOrUpdate,
		DeleteContext: resourceDeployEnvironmentPermissionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDeployEnvironmentPermissionImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the application ID.`,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the environment ID.`,
			},
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the role ID.`,
			},
			"permission_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the permission name.`,
			},
			"permission_value": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to enable the permission.`,
			},
		},
	}
}

func resourceDeployEnvironmentPermissionCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("codearts_deploy", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	err = modifyDeployEnvironmentPermission(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.IsNewResource() {
		d.SetId(d.Get("application_id").(string) + "/" + d.Get("environment_id").(string) + "/" +
			d.Get("role_id").(string) + "/" + d.Get("permission_name").(string))
	}

	return resourceDeployEnvironmentPermissionRead(ctx, d, meta)
}

func modifyDeployEnvironmentPermission(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v2/applications/{application_id}/environments/{environment_id}/permissions"
	modifyPath := client.Endpoint + httpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{application_id}", d.Get("application_id").(string))
	modifyPath = strings.ReplaceAll(modifyPath, "{environment_id}", d.Get("environment_id").(string))

	modifyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: buildDeployEnvironmentPermissionBodyParams(d),
	}

	_, err := client.Request("PUT", modifyPath, &modifyOpt)
	if err != nil {
		return fmt.Errorf("error modifying CodeArts deploy environment permission: %s", err)
	}

	return nil
}

func buildDeployEnvironmentPermissionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"role_id":          d.Get("role_id"),
		"permission_name":  d.Get("permission_name"),
		"permission_value": d.Get("permission_value"),
	}
}

func resourceDeployEnvironmentPermissionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	permissionMatrix, err := getDeployEnvironmentPermissionMatrix(client, d.Get("application_id").(string), d.Get("environment_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts deploy environment permission")
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

func resourceDeployEnvironmentPermissionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting permission resource is not supported. The resource is only removed from the state," +
		" the environment permission matrix remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceDeployEnvironmentPermissionImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid format specified for import ID, want "+
			"'<application_id>/<environment_id>/<role_id>/<permission_name>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("application_id", parts[0]),
		d.Set("environment_id", parts[1]),
		d.Set("role_id", parts[2]),
		d.Set("permission_name", parts[3]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
