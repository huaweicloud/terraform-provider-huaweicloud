package fgs

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/application"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API FunctionGraph GET /v2/{project_id}/fgs/application/templates
func DataSourceFunctionGraphApplicationTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFunctionGraphApplicationTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"runtime": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceFunctionGraphApplicationTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.FgsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	opts := application.ListOpts{
		Runtime: d.Get("runtime").(string),
	}
	templateList, err := application.ListTemplates(client, opts)
	if err != nil {
		return diag.Errorf("error retrieving application templates: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("templates", flattenListTemplates(filterListTemplateByCategory(templateList, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListTemplates(templates []application.Template) []map[string]interface{} {
	if len(templates) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(templates))
	for i, template := range templates {
		result[i] = map[string]interface{}{
			"id":          template.ID,
			"name":        template.Name,
			"runtime":     template.Runtime,
			"category":    template.Category,
			"type":        template.Type,
			"description": template.Description,
		}
	}
	return result
}

func filterListTemplateByCategory(all []application.Template, d *schema.ResourceData) []application.Template {
	if len(all) == 0 {
		return nil
	}
	rst := make([]application.Template, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("category"); ok && v.Category != param {
			continue
		}
		rst = append(rst, v)
	}
	return rst
}
