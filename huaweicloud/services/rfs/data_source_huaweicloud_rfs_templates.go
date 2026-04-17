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

// @API RFS GET /v1/{project_id}/templates
func DataSourceRfsTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRfsTemplateRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_description": {
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
						"latest_version_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"latest_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRfsTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "rfs"
		httpUrl      = "v1/{project_id}/templates"
		allTemplates = make([]interface{}, 0)
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
	requestPath = fmt.Sprintf("%s?limit=%d", requestPath, limit)
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
			return diag.Errorf("error retrieving RFS templates: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		templates := utils.PathSearch("templates", respBody, make([]interface{}, 0)).([]interface{})
		if len(templates) == 0 {
			break
		}

		allTemplates = append(allTemplates, templates...)

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("templates", flattenRfsTemplates(allTemplates)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRfsTemplates(templates []interface{}) []interface{} {
	if len(templates) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(templates))
	for _, template := range templates {
		rst = append(rst, map[string]interface{}{
			"template_id":                utils.PathSearch("template_id", template, nil),
			"template_name":              utils.PathSearch("template_name", template, nil),
			"template_description":       utils.PathSearch("template_description", template, nil),
			"create_time":                utils.PathSearch("create_time", template, nil),
			"update_time":                utils.PathSearch("update_time", template, nil),
			"latest_version_description": utils.PathSearch("latest_version_description", template, nil),
			"latest_version_id":          utils.PathSearch("latest_version_id", template, nil),
		})
	}
	return rst
}
