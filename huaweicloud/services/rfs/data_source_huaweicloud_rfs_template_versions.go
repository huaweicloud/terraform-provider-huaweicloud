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

// @API RFS GET /v1/{project_id}/templates/{template_name}/versions
func DataSourceRfsTemplateVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRfsTemplateVersionRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_name": {
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

func buildTemplateVersionsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if templateId, ok := d.GetOk("template_id"); ok {
		queryParams = fmt.Sprintf("&template_id=%s", templateId)
	}

	return queryParams
}

func dataSourceRfsTemplateVersionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "rfs"
		httpUrl      = "v1/{project_id}/templates/{template_name}/versions"
		templateName = d.Get("template_name").(string)
		allVersions  = make([]interface{}, 0)
		limit        = 1000
		marker       = ""
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_name}", templateName)
	requestPath = fmt.Sprintf("%s?limit=%d%s", requestPath, limit, buildTemplateVersionsQueryParams(d))

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
			"Content-Type":      "application/json",
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := requestPath
		if marker != "" {
			requestPathWithMarker = fmt.Sprintf("%s&marker=%s", requestPathWithMarker, marker)
		}

		resp, err := client.Request("GET", requestPathWithMarker, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RFS template versions: %s", err)
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

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("versions", flattenRfsTemplateVersions(allVersions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRfsTemplateVersions(versions []interface{}) []interface{} {
	if len(versions) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(versions))
	for _, version := range versions {
		rst = append(rst, map[string]interface{}{
			"version_id":          utils.PathSearch("version_id", version, nil),
			"template_id":         utils.PathSearch("template_id", version, nil),
			"template_name":       utils.PathSearch("template_name", version, nil),
			"version_description": utils.PathSearch("version_description", version, nil),
			"create_time":         utils.PathSearch("create_time", version, nil),
		})
	}
	return rst
}
