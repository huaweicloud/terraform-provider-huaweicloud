package cdn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/configuration/templates
func DataSourceDomainTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainTemplatesRead,

		Schema: map[string]*schema.Schema{
			// Attributes.
			"templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        domainTemplateSchema(),
				Description: `The list of domain templates.`,
			},
		},
	}
}

func domainTemplateSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the domain template.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the domain template.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the domain template.`,
			},
			"configs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The configuration of the domain template, in JSON format.`,
			},
			"type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The type of the domain template. Valid values are 1 (system preset template) and 2 (tenant custom template).`,
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The account ID.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the domain template.`,
			},
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The modification time of the domain template.`,
			},
		},
	}
}

func flattenDomainTemplates(templates []interface{}) []map[string]interface{} {
	if len(templates) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(templates))
	for _, template := range templates {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("tml_id", template, nil),
			"name":        utils.PathSearch("tml_name", template, nil),
			"description": utils.PathSearch("remark", template, nil),
			"configs":     utils.JsonToString(utils.PathSearch("configs", template, nil)),
			"type":        utils.PathSearch("type", template, nil),
			"account_id":  utils.PathSearch("account_id", template, nil),
			"create_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", template,
				float64(0)).(float64))/1000, false),
			"modify_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("modify_time", template,
				float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceDomainTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg = meta.(*config.Config)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	templates, err := listDomainTemplates(client)
	if err != nil {
		return diag.Errorf("error querying CDN domain templates: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("templates", flattenDomainTemplates(templates)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
