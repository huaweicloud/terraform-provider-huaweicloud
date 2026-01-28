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
func DataSourceResourcePermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourcePermissionsRead,
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
				Description: `Indicates the ID of RAM permission.`,
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
				Description: `Indicates the RAM permission creation time.`,
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

func buildGetResourcePermissionsQueryParams(d *schema.ResourceData, nextMarker string) string {
	queryParams := "?limit=2000"
	if nextMarker != "" {
		queryParams = fmt.Sprintf("%s&marker=%s", queryParams, nextMarker)
	}

	if v, ok := d.GetOk("resource_type"); ok {
		queryParams = fmt.Sprintf("%s&resource_type=%v", queryParams, v)
	}

	if v, ok := d.GetOk("permission_type"); ok {
		queryParams = fmt.Sprintf("%s&permission_type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceResourcePermissionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		mErr       *multierror.Error
		httpUrl    = "v1/permissions"
		product    = "ram"
		nextMarker string
		result     []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := requestPath + buildGetResourcePermissionsQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving RAM resource permissions: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		permissionsResp := utils.PathSearch("permissions", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, permissionsResp...)
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

	mErr = multierror.Append(mErr,
		d.Set("permissions", flattenResourcePermissions(filterResourcePermissions(result, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterResourcePermissions(all []interface{}, d *schema.ResourceData) []interface{} {
	name := d.Get("name").(string)
	if name == "" {
		return all
	}

	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if fmt.Sprint(name) != fmt.Sprint(utils.PathSearch("name", v, "").(string)) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenResourcePermissions(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                       utils.PathSearch("id", v, nil),
			"name":                     utils.PathSearch("name", v, nil),
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
