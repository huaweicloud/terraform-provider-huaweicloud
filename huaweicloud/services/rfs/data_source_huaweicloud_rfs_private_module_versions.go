package rfs

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

// @API RFS GET /v1/private-modules/{module_name}/versions
func DataSourcePrivateModuleVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateModuleVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"module_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"module_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"versions": {
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
						"module_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildListPrivateModuleVersionsQueryParams(d *schema.ResourceData, marker string) string {
	rst := ""

	if v, ok := d.GetOk("module_id"); ok {
		rst += fmt.Sprintf("&module_id=%s", v.(string))
	}

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

func dataSourcePrivateModuleVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/private-modules/{module_name}/versions"
		allVersions = make([]interface{}, 0)
		nextMarker  string
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
	requestPath = strings.ReplaceAll(requestPath, "{module_name}", d.Get("module_name").(string))
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildListPrivateModuleVersionsQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RFS private module versions: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		versions := utils.PathSearch("versions", respBody, make([]interface{}, 0)).([]interface{})
		if len(versions) == 0 {
			break
		}

		allVersions = append(allVersions, versions...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("versions", flattenPrivateModuleVersions(allVersions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPrivateModuleVersions(versions []interface{}) []interface{} {
	if len(versions) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(versions))
	for _, v := range versions {
		rst = append(rst, map[string]interface{}{
			"module_name":         utils.PathSearch("module_name", v, nil),
			"module_id":           utils.PathSearch("module_id", v, nil),
			"module_version":      utils.PathSearch("module_version", v, nil),
			"version_description": utils.PathSearch("version_description", v, nil),
			"create_time":         utils.PathSearch("create_time", v, nil),
		})
	}
	return rst
}
