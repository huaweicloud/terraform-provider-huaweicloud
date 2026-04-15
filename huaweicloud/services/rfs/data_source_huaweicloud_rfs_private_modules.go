package rfs

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

// @API RFS GET /v1/private-modules
func DataSourcePrivateModules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateModulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"module_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"module_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"module_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildListPrivateModulesQueryParams(d *schema.ResourceData, marker string) string {
	rst := ""

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%s", v.(string))
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%s", v.(string))
	}

	if marker != "" {
		rst += fmt.Sprintf("&marker=%s", marker)
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourcePrivateModulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/private-modules"
		allModules = make([]interface{}, 0)
		nextMarker string
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildListPrivateModulesQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving private modules: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		modules := utils.PathSearch("modules", respBody, make([]interface{}, 0)).([]interface{})
		if len(modules) == 0 {
			break
		}

		allModules = append(allModules, modules...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("modules", flattenPrivateModules(allModules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPrivateModules(modules []interface{}) []interface{} {
	if len(modules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(modules))
	for _, v := range modules {
		rst = append(rst, map[string]interface{}{
			"module_name":        utils.PathSearch("module_name", v, nil),
			"module_id":          utils.PathSearch("module_id", v, nil),
			"module_description": utils.PathSearch("module_description", v, nil),
			"create_time":        utils.PathSearch("create_time", v, nil),
			"update_time":        utils.PathSearch("update_time", v, nil),
		})
	}
	return rst
}
