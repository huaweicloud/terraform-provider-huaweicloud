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

// @API RAM GET /v1/permissions/{permission_id}
func DataSourceResourcePermission() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourcePermissionRead,
		Schema: map[string]*schema.Schema{
			"permission_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permission_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permission": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     permissionSchema(),
			},
		},
	}
}

func permissionSchema() *schema.Resource {
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
			"content": {
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

func buildGetResourcePermissionQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("permission_version"); ok {
		return fmt.Sprintf("?permission_version=%v", v)
	}

	return ""
}

func dataSourceResourcePermissionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/permissions/{permission_id}"
		product = "ram"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{permission_id}", d.Get("permission_id").(string))
	requestPath += buildGetResourcePermissionQueryParams(d)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving RAM resource permission: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("permission", flattenPermission(utils.PathSearch("permission", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPermission(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":                       utils.PathSearch("id", resp, nil),
			"name":                     utils.PathSearch("name", resp, nil),
			"resource_type":            utils.PathSearch("resource_type", resp, nil),
			"content":                  utils.PathSearch("content", resp, nil),
			"is_resource_type_default": utils.PathSearch("is_resource_type_default", resp, nil),
			"created_at":               utils.PathSearch("created_at", resp, nil),
			"updated_at":               utils.PathSearch("updated_at", resp, nil),
			"permission_urn":           utils.PathSearch("permission_urn", resp, nil),
			"permission_type":          utils.PathSearch("permission_type", resp, nil),
			"default_version":          utils.PathSearch("default_version", resp, nil),
			"version":                  utils.PathSearch("version", resp, nil),
			"status":                   utils.PathSearch("status", resp, nil),
		},
	}
}
