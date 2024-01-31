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
				Type:     schema.TypeString,
				Optional: true,
				Description: `Specifies the resource type of RAM permission. Valid values are **vpc:subnets**, 
**dns:zone** and **dns:resolverRule**.`,
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
		},
	}
	return &sc
}

func resourceRAMPermissionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getRAMPermissionsHttpUrl = "v1/permissions"
		getRAMPermissionsProduct = "ram"
	)
	getRAMPermissionsClient, err := cfg.NewServiceClient(getRAMPermissionsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM Client: %s", err)
	}

	getRAMPermissionsPath := getRAMPermissionsClient.Endpoint + getRAMPermissionsHttpUrl
	getRAMPermissionsPath += buildGetRAMPermissionsQueryParams(d)

	getRAMPermissionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getRAMPermissionsResp, err := getRAMPermissionsClient.Request("GET", getRAMPermissionsPath,
		&getRAMPermissionsOpt)

	if err != nil {
		return diag.Errorf("error retrieving RAM permissions, %s", err)
	}

	getRAMPermissionsRespBody, err := utils.FlattenResponse(getRAMPermissionsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("permissions", flattenGetPermissionsResponseBody(getRAMPermissionsRespBody, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetPermissionsResponseBody(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	name := d.Get("name").(string)

	curJson := utils.PathSearch("permissions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		permissionName := utils.PathSearch("name", v, "")
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
		})
	}
	return rst
}

// buildGetRAMPermissionsQueryParams use the max limit number.
// Paging is not currently implemented
func buildGetRAMPermissionsQueryParams(d *schema.ResourceData) string {
	res := "?limit=2000&marker=1"
	if v, ok := d.GetOk("resource_type"); ok {
		res = fmt.Sprintf("%s&resource_type=%v", res, v)
	}
	return res
}
