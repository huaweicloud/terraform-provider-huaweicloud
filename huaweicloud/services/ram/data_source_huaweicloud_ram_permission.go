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
func DataSourcePermission() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePermissionRead,
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

func dataSourcePermissionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getPermissionProduct = "ram"
	getPermissionClient, err := cfg.NewServiceClient(getPermissionProduct, region)
	if err != nil {
		return diag.Errorf("Error creating RAM client: %s", err)
	}

	getPermissionRespBody, err := getPermission(getPermissionClient, d)

	if err != nil {
		return diag.Errorf("error retrieving RAM permission: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("permission", []interface{}{utils.PathSearch("permission", getPermissionRespBody, nil)}),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getPermission(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	permissionId := d.Get("permission_id").(string)

	var (
		getPermissionHttpUrl = "v1/permissions/{permission_id}"
	)
	getPermissionHttpPath := client.Endpoint + getPermissionHttpUrl
	getPermissionHttpPath = strings.ReplaceAll(getPermissionHttpPath, "{permission_id}", permissionId)
	getPermissionHttpPath += buildGetPermissionQueryParams(d)

	getPermissionHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getPermissionHttpResp, err := client.Request("GET", getPermissionHttpPath, &getPermissionHttpOpt)
	if err != nil {
		return nil, err
	}
	getPermissionRespBody, err := utils.FlattenResponse(getPermissionHttpResp)
	if err != nil {
		return nil, err
	}
	return getPermissionRespBody, nil
}

func buildGetPermissionQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("permission_version"); ok {
		res = fmt.Sprintf("%s&permission_version=%v", res, v)
	}

	if res != "" {
		return "?" + res[1:]
	}

	return res
}
