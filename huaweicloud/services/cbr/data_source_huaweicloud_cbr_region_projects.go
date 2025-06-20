package cbr

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBR GET /v3/region-projects
func DataSourceCbrRegionProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCbrRegionProjectsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataCbrRegionProjectSchema(),
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"self": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataCbrRegionProjectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_domain": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"self": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCbrRegionProjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + "v3/region-projects"
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error querying CBR region projects: %s", err)
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	projects := utils.PathSearch("projects", respBody, []interface{}{}).([]interface{})

	uuidStr, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuidStr)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("projects", flattenCbrRegionProjects(projects)),
		d.Set("links", flattenLinks(utils.PathSearch("links", respBody, nil))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCbrRegionProjects(projects []interface{}) []map[string]interface{} {
	if len(projects) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(projects))
	for i, p := range projects {
		result[i] = map[string]interface{}{
			"domain_id":   utils.PathSearch("domain_id", p, nil),
			"is_domain":   utils.PathSearch("is_domain", p, nil),
			"parent_id":   utils.PathSearch("parent_id", p, nil),
			"name":        utils.PathSearch("name", p, nil),
			"description": utils.PathSearch("description", p, nil),
			"id":          utils.PathSearch("id", p, nil),
			"enabled":     utils.PathSearch("enabled", p, nil),
			"links":       flattenLinks(utils.PathSearch("links", p, nil)),
		}
	}
	return result
}

func flattenLinks(raw interface{}) []map[string]interface{} {
	if raw == nil {
		return nil
	}
	rstMap := map[string]interface{}{
		"self": utils.PathSearch("self", raw, nil),
	}

	return []map[string]interface{}{rstMap}
}
