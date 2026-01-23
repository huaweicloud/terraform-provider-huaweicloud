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
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_resource_type_default": {
				Type:     schema.TypeBool,
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
			"permission_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_version": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourcePermissionVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		permissionId = d.Get("permission_id").(string)
		httpUrl      = "v1/permissions/{permission_id}/versions"
		product      = "ram"
		marker       string
		result       []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ram client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{permission_id}", permissionId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := requestPath + buildListPermissionVersionsQueryParams(marker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving RAM permission versions: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		permissionsResp := utils.PathSearch("permissions", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, permissionsResp...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
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
		d.Set("permissions", flattenPermissionVersionsResp(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListPermissionVersionsQueryParams(marker string) string {
	// The default value of limit is `2000`
	res := "?limit=2000"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func flattenPermissionVersionsResp(resp []interface{}) []interface{} {
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
