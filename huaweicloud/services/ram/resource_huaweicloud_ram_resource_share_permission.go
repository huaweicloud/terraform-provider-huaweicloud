package ram

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

var resourceSharePermissionNonUpdatableParams = []string{"resource_share_id", "permission_id", "replace"}

// @API RAM POST /v1/resource-shares/{resource_share_id}/associate-permission
// @API RAM POST /v1/resource-shares/{resource_share_id}/disassociate-permission
// @API RAM GET /v1/resource-shares/{resource_share_id}/associated-permissions
func ResourceSharePermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharePermissionCreate,
		UpdateContext: resourceSharePermissionUpdate,
		ReadContext:   resourceSharePermissionRead,
		DeleteContext: resourceSharePermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSharePermissionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(resourceSharePermissionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"resource_share_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permission_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replace": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"permission_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildCreateResourceSharePermissionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"permission_id": d.Get("permission_id"),
		"replace":       d.Get("replace"),
	}
}

func resourceSharePermissionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		resourceShareId = d.Get("resource_share_id").(string)
		permissionId    = d.Get("permission_id").(string)
		httpUrl         = "v1/resource-shares/{resource_share_id}/associate-permission"
		product         = "ram"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{resource_share_id}", resourceShareId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateResourceSharePermissionBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error binding RAM shared resource permission: %s", err)
	}

	d.SetId(resourceShareId + "/" + permissionId)

	return resourceSharePermissionRead(ctx, d, meta)
}

func buildGetAssociatedPermissionsQueryParams(marker string) string {
	// The default value of limit is `2000`
	queryParams := "?limit=2000"
	if marker != "" {
		queryParams = fmt.Sprintf("%s&marker=%v", queryParams, marker)
	}

	return queryParams
}

func resourceSharePermissionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                      = meta.(*config.Config)
		region                   = cfg.GetRegion(d)
		mErr                     *multierror.Error
		resourceShareId          = d.Get("resource_share_id").(string)
		permissionId             = d.Get("permission_id").(string)
		httpUrl                  = "v1/resource-shares/{resource_share_id}/associated-permissions"
		product                  = "ram"
		allAssociatedPermissions = make([]interface{}, 0)
		marker                   string
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{resource_share_id}", resourceShareId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := requestPath + buildGetAssociatedPermissionsQueryParams(marker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving RAM shared resource permissions: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		permissionsResp := utils.PathSearch("associated_permissions", respBody, make([]interface{}, 0)).([]interface{})
		allAssociatedPermissions = append(allAssociatedPermissions, permissionsResp...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	jsonPath := fmt.Sprintf("[?permission_id=='%s'&&status=='associated']|[0]", permissionId)
	associatedPermission := utils.PathSearch(jsonPath, allAssociatedPermissions, nil)
	if associatedPermission == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(mErr,
		d.Set("permission_id", utils.PathSearch("permission_id", associatedPermission, nil)),
		d.Set("permission_name", utils.PathSearch("permission_id", associatedPermission, nil)),
		d.Set("resource_type", utils.PathSearch("resource_type", associatedPermission, nil)),
		d.Set("status", utils.PathSearch("status", associatedPermission, nil)),
		d.Set("created_at", utils.PathSearch("created_at", associatedPermission, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", associatedPermission, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSharePermissionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func buildDeleteResourceSharePermissionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"permission_id": d.Get("permission_id"),
	}
}

func resourceSharePermissionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		resourceShareId = d.Get("resource_share_id").(string)
		httpUrl         = "v1/resource-shares/{resource_share_id}/disassociate-permission"
		product         = "ram"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{resource_share_id}", resourceShareId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteResourceSharePermissionBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error unbinding RAM shared resource permission: %s", err)
	}

	return nil
}

func resourceSharePermissionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID,"+
			" must be '<resource_share_id>/<permission_id>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("resource_share_id", parts[0]),
		d.Set("permission_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
