// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RAM
// ---------------------------------------------------------------

package ram

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RAM GET /v1/permissions
func DataSourceRAMPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceRAMPermissionsRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource type of RAM permission.`,
			},
			"permission_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the permission.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of RAM permission.`,
			},
			"permissions": {
				Type:        schema.TypeList,
				Elem:        permissionsSchema(),
				Computed:    true,
				Description: `Indicates the list of the RAM permissions`,
			},
		},
	}
}

func permissionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the id of RAM permission.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of RAM permission.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the resource type of RAM permission.`,
			},
			"is_resource_type_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the RAM permission resource type is default.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the RAM permission create time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the RAM permission last update time.`,
			},
			"permission_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the URN for the permission.`,
			},
			"permission_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the permission type.`,
			},
			"default_version": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the current version is the default version.`,
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the version of the permission.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the permission.`,
			},
		},
	}
	return &sc
}

func buildGetRAMPermissionsQueryParams(d *schema.ResourceData, nextMarker string) string {
	queryParam := "?limit=2000"
	if nextMarker != "" {
		queryParam = fmt.Sprintf("%s&marker=%s", queryParam, nextMarker)
	}

	if v, ok := d.GetOk("resource_type"); ok {
		queryParam = fmt.Sprintf("%s&resource_type=%v", queryParam, v)
	}

	if v, ok := d.GetOk("permission_type"); ok {
		queryParam = fmt.Sprintf("%s&permission_type=%v", queryParam, v)
	}
	return queryParam
}

func resourceRAMPermissionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		mErr             *multierror.Error
		httpUrl          = "v1/permissions"
		product          = "ram"
		nextMarker       string
		totalPermissions []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParam := requestPath + buildGetRAMPermissionsQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParam, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RAM permissions, %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		permissions := utils.PathSearch("permissions", respBody, make([]interface{}, 0)).([]interface{})
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
		d.Set("permissions", flattenGetPermissionsResponseBody(totalPermissions, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetPermissionsResponseBody(totalPermissions []interface{}, d *schema.ResourceData) []interface{} {
	if len(totalPermissions) == 0 {
		return nil
	}

	name := d.Get("name").(string)
	rst := make([]interface{}, 0, len(totalPermissions))
	for _, v := range totalPermissions {
		permissionName := utils.PathSearch("name", v, "").(string)
		if name != "" && name != permissionName {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"id":                       utils.PathSearch("id", v, nil),
			"name":                     permissionName,
			"resource_type":            utils.PathSearch("resource_type", v, nil),
			"is_resource_type_default": utils.PathSearch("is_resource_type_default", v, nil),
			"created_at":               utils.PathSearch("created_at", v, nil),
			"updated_at":               utils.PathSearch("updated_at", v, nil),
			"permission_urn":           utils.PathSearch("permission_urn", v, nil),
			"permission_type":          utils.PathSearch("permission_type", v, nil),
			"default_version":          utils.PathSearch("default_version", v, nil),
			"version":                  utils.PathSearch("version", v, nil),
			"status":                   utils.PathSearch("status", v, nil),
		})
	}
	return rst
}
