package ram

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RAM GET /v1/resource-shares/{resource_share_id}/associated-permissions
func DataSourceAssociatedPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssociatedPermissionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"resource_share_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the resource share.`,
			},
			"permission_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the RAM managed permission.`,
			},
			"associated_permissions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of RAM managed permissions associated with the resource share.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the permission was last updated.`,
						},
						"permission_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The permission ID.`,
						},
						"permission_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the RAM managed permission.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource type to which the permission applies.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the permission.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the permission was created.`,
						},
					},
				},
			},
		},
	}
}

// buildAssociatedPermissionPathWithQueryParams The default limit value for paging query is `200`, so the limit value is
// not configured here.
func buildAssociatedPermissionPathWithQueryParams(d *schema.ResourceData, path, nextMarker string) string {
	queryParam := ""
	if permissionName := d.Get("permission_name").(string); permissionName != "" {
		queryParam = fmt.Sprintf("%s&permission_name=%s", queryParam, permissionName)
	}

	if nextMarker != "" {
		queryParam = fmt.Sprintf("%s&marker=%s", queryParam, nextMarker)
	}

	if queryParam == "" {
		return path
	}
	return fmt.Sprintf("%s?%s", path, queryParam[1:])
}

func dataSourceAssociatedPermissionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		mErr             *multierror.Error
		nextMarker       string
		httpUrl          = "v1/resource-shares/{resource_share_id}/associated-permissions"
		product          = "ram"
		totalPermissions []interface{}
		resourceShareID  = d.Get("resource_share_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{resource_share_id}", resourceShareID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParam := buildAssociatedPermissionPathWithQueryParams(d, requestPath, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParam, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RAM resource share associated permissions: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		permissions := utils.PathSearch("associated_permissions", respBody, make([]interface{}, 0)).([]interface{})
		if len(permissions) > 0 {
			totalPermissions = append(totalPermissions, permissions...)
		}

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("associated_permissions", flattenResourceShareAssociatedPermissions(totalPermissions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResourceShareAssociatedPermissions(permissions []interface{}) []interface{} {
	if len(permissions) == 0 {
		return nil
	}

	rst := make([]interface{}, len(permissions))
	for i, v := range permissions {
		rst[i] = map[string]interface{}{
			"updated_at":      utils.PathSearch("updated_at", v, nil),
			"permission_id":   utils.PathSearch("permission_id", v, nil),
			"permission_name": utils.PathSearch("permission_name", v, nil),
			"resource_type":   utils.PathSearch("resource_type", v, nil),
			"status":          utils.PathSearch("status", v, nil),
			"created_at":      utils.PathSearch("created_at", v, nil),
		}
	}
	return rst
}
