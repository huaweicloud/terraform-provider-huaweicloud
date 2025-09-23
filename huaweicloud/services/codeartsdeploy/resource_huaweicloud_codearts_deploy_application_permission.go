package codeartsdeploy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CodeArtsDeploy PUT /v3/applications/permissions
func ResourceDeployApplicationPermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeployApplicationPermissionCreateOrUpdate,
		ReadContext:   resourceDeployApplicationPermissionRead,
		UpdateContext: resourceDeployApplicationPermissionCreateOrUpdate,
		DeleteContext: resourceDeployApplicationPermissionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the project ID for CodeArts service.`,
			},
			"application_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the application IDs.`,
			},
			"roles": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `Specifies the role permissions list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the role ID.`,
						},
						"can_modify": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Specifies whether the role has the editing permission.`,
						},
						"can_disable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Specifies whether the role has the permission to disable application.`,
						},
						"can_delete": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Specifies whether the role has the deletion permission.`,
						},
						"can_view": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Specifies whether the role has the view permission.`,
						},
						"can_execute": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Specifies whether the role has the deployment permission.`,
						},
						"can_copy": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Specifies whether the role has the copy permission.`,
						},
						"can_manage": {
							Type:     schema.TypeBool,
							Required: true,
							Description: `Specifies whether the role has the management permission, including adding, deleting,
		modifying, querying deployment and permission modification.`,
						},
						"can_create_env": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Specifies whether the role has the permission to create an environment.`,
						},
					},
				},
			},
		},
	}
}

func resourceDeployApplicationPermissionCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("codearts_deploy", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	if !d.IsNewResource() {
		timeout = d.Timeout(schema.TimeoutUpdate)
	}

	err = modifyDeployApplicationPermission(ctx, client, d, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.IsNewResource() {
		id, err := uuid.GenerateUUID()
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(id)
	}

	return resourceDeployApplicationPermissionRead(ctx, d, meta)
}

func modifyDeployApplicationPermission(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	httpUrl := "v3/applications/permissions"
	modifyPath := client.Endpoint + httpUrl
	modifyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			// API's language defaults to zh-cn.
			// When response is error, diag.Diagnostics sometimes will cause error:
			// `error while marshaling: string field contains invalid UTF-8`
			"X-Language": "en-us",
		},
		JSONBody: buildDeployApplicationPermissionBodyParams(d),
	}

	err := resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err := client.Request("PUT", modifyPath, &modifyOpt)
		isRetry, err := handleDeployApplicationPermissionLevelOperationError(err)
		if isRetry {
			// lintignore:R018
			time.Sleep(10 * time.Second)
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error updating CodeArts deploy application permission: %s", err)
	}

	return nil
}

func buildDeployApplicationPermissionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"project_id":      d.Get("project_id"),
		"application_ids": d.Get("application_ids").(*schema.Set).List(),
		"roles":           buildDeployApplicationPermissionBodyParamsRoles(d),
	}
}

func buildDeployApplicationPermissionBodyParamsRoles(d *schema.ResourceData) []map[string]interface{} {
	rawArray := d.Get("roles").(*schema.Set).List()
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		if raw, isMap := v.(map[string]interface{}); isMap {
			rst = append(rst, map[string]interface{}{
				"dev_role_id":    raw["role_id"],
				"can_modify":     raw["can_modify"],
				"can_disable":    raw["can_disable"],
				"can_delete":     raw["can_delete"],
				"can_view":       raw["can_view"],
				"can_execute":    raw["can_execute"],
				"can_copy":       raw["can_copy"],
				"can_manage":     raw["can_manage"],
				"can_create_env": raw["can_create_env"],
			})
		}
	}

	return rst
}

func resourceDeployApplicationPermissionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDeployApplicationPermissionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting permission resource is not supported. The resource is only removed from the state," +
		" the application permission matrix remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
