package ram

import (
	"context"
	"errors"
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
func ResourceRAMResourceSharePermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRAMResourceSharePermissionCreate,
		UpdateContext: resourceRAMResourceSharePermissionUpdate,
		ReadContext:   resourceRAMResourceSharePermissionRead,
		DeleteContext: resourceRAMResourceSharePermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRAMResourceSharePermissionImportState,
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

func resourceRAMResourceSharePermissionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	resourceShareId := d.Get("resource_share_id").(string)
	permissionId := d.Get("permission_id").(string)
	var (
		createResourceSharePermissionHttpUrl = "v1/resource-shares/{resource_share_id}/associate-permission"
		createResourceSharePermissionProduct = "ram"
	)
	ramClient, err := cfg.NewServiceClient(createResourceSharePermissionProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	createResourceSharePermissionPath := ramClient.Endpoint + createResourceSharePermissionHttpUrl
	createResourceSharePermissionPath = strings.ReplaceAll(createResourceSharePermissionPath, "{resource_share_id}", resourceShareId)
	createResourceSharePermissionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResourceSharePermissionOpt.JSONBody = utils.RemoveNil(buildCreateResourceSharePermissionBodyParams(d))
	_, err = ramClient.Request("POST", createResourceSharePermissionPath, &createResourceSharePermissionOpt)
	if err != nil {
		return diag.Errorf("error creating RAM share permission: %s", err)
	}

	d.SetId(resourceShareId + "/" + permissionId)
	return resourceRAMResourceSharePermissionRead(ctx, d, meta)
}

func resourceRAMResourceSharePermissionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRAMResourceSharePermissionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	resourceShareId := d.Get("resource_share_id").(string)
	permissionId := d.Get("permission_id").(string)

	var (
		getRAMAssociatedPermissionHttpUrl = "v1/resource-shares/{resource_share_id}/associated-permissions"
		getRAMAssociatedPermissionProduct = "ram"
	)
	getRAMAssociatedPermissionClient, err := cfg.NewServiceClient(getRAMAssociatedPermissionProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	getRAMAssociatedPermissionPath := getRAMAssociatedPermissionClient.Endpoint + getRAMAssociatedPermissionHttpUrl
	getRAMAssociatedPermissionPath = strings.ReplaceAll(getRAMAssociatedPermissionPath, "{resource_share_id}", resourceShareId)
	getRAMShareOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var associatedPermission interface{}
	var marker string
	var queryPath string

	for {
		queryPath = getRAMAssociatedPermissionPath + buildGetAssociatedPermissionQueryParams(marker)
		getRAMAssociatedPermissionResp, err := getRAMAssociatedPermissionClient.Request("GET", queryPath, &getRAMShareOpt)
		if err != nil {
			return diag.Errorf("error retrieving associated permissions, error: %s", err)
		}

		getRAMAssociatedPermissionRespBody, err := utils.FlattenResponse(getRAMAssociatedPermissionResp)
		if err != nil {
			return diag.FromErr(err)
		}

		associatedPermission = utils.PathSearch(
			fmt.Sprintf("associated_permissions[?permission_id=='%s'&&status=='associated']|[0]", permissionId),
			getRAMAssociatedPermissionRespBody,
			nil,
		)
		if associatedPermission != nil {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", getRAMAssociatedPermissionRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	if associatedPermission == nil {
		return common.CheckDeletedDiag(d, err, "error retrieving associated permissions")
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

func resourceRAMResourceSharePermissionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	resourceShareId := d.Get("resource_share_id").(string)
	var (
		deleteResourceSharePermissionHttpUrl = "v1/resource-shares/{resource_share_id}/disassociate-permission"
		deleteResourceSharePermissionProduct = "ram"
	)
	ramClient, err := cfg.NewServiceClient(deleteResourceSharePermissionProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	deleteResourceSharePermissionPath := ramClient.Endpoint + deleteResourceSharePermissionHttpUrl
	deleteResourceSharePermissionPath = strings.ReplaceAll(deleteResourceSharePermissionPath, "{resource_share_id}", resourceShareId)
	deleteResourceSharePermissionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteResourceSharePermissionOpt.JSONBody = utils.RemoveNil(buildDeleteResourceSharePermissionBodyParams(d))
	_, err = ramClient.Request("POST", deleteResourceSharePermissionPath, &deleteResourceSharePermissionOpt)
	if err != nil {
		return diag.Errorf("error delete RAM resource share permission: %s", err)
	}

	return nil
}

func buildCreateResourceSharePermissionBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := make(map[string]interface{})
	params["permission_id"] = d.Get("permission_id").(string)

	if v, ok := d.GetOk("replace"); ok {
		params["replace"] = v
	}

	return params
}

func buildDeleteResourceSharePermissionBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := make(map[string]interface{})
	params["permission_id"] = d.Get("permission_id").(string)

	return params
}

func buildGetAssociatedPermissionQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func resourceRAMResourceSharePermissionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <resource_share_id>/<permission_id>")
	}

	mErr := multierror.Append(nil,
		d.Set("resource_share_id", parts[0]),
		d.Set("permission_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
