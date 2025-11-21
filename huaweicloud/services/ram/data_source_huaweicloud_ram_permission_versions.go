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

// @API RAM GET /v1/permissions/{permission_id}/versions
func DataSourcePermissionVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePermissionVersionsRead,
		Schema: map[string]*schema.Schema{
			"permission_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     permissionVersionsSchema(),
			},
		},
	}
}

func permissionVersionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_version": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_resource_type_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission_urn": {
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
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourcePermissionVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	permissionId := d.Get("permission_id").(string)
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listPermissionVersionsHttpUrl := "v1/permissions/{permission_id}/versions"
	listPermissionVersionsProduct := "ram"
	listPermissionVersionsClient, err := cfg.NewServiceClient(listPermissionVersionsProduct, region)
	if err != nil {
		return diag.Errorf("error creating ram client: %s", err)
	}

	listPermissionVersionsPath := listPermissionVersionsClient.Endpoint + listPermissionVersionsHttpUrl
	listPermissionVersionsPath = strings.ReplaceAll(listPermissionVersionsPath, "{permission_id}", permissionId)

	listPermissionVersionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var permissionVersions []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listPermissionVersionsPath + buildListPermissionVersionsQueryParams(marker)
		listPermissionVersionsResp, err := listPermissionVersionsClient.Request("GET", queryPath, &listPermissionVersionsOpt)
		if err != nil {
			return diag.Errorf("error retrieving RAM permission versions: %s", err)
		}

		listPermissionVersionsRespBody, err := utils.FlattenResponse(listPermissionVersionsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePagePermissionVersions := FlattenPermissionVersionsResp(listPermissionVersionsRespBody)
		permissionVersions = append(permissionVersions, onePagePermissionVersions...)
		marker = utils.PathSearch("page_info.next_marker", listPermissionVersionsRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("permissions", permissionVersions),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListPermissionVersionsQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func FlattenPermissionVersionsResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("permissions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"created_at":               utils.PathSearch("created_at", v, nil),
			"default_version":          utils.PathSearch("default_version", v, nil),
			"id":                       utils.PathSearch("id", v, nil),
			"is_resource_type_default": utils.PathSearch("is_resource_type_default", v, nil),
			"name":                     utils.PathSearch("name", v, nil),
			"permission_type":          utils.PathSearch("permission_type", v, nil),
			"permission_urn":           utils.PathSearch("permission_urn", v, nil),
			"resource_type":            utils.PathSearch("resource_type", v, nil),
			"status":                   utils.PathSearch("status", v, nil),
			"updated_at":               utils.PathSearch("updated_at", v, nil),
			"version":                  utils.PathSearch("version", v, nil),
		})
	}
	return rst
}
