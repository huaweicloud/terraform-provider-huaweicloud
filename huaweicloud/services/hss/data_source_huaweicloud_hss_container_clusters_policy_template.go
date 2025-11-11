package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/container/clusters/protection-policy-templates/{policy_template_id}
func DataSourceContainerClustersPolicyTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerClustersPolicyTemplateRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy_template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"constraint_template": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceContainerClustersPolicyTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/container/clusters/protection-policy-templates/{policy_template_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_template_id}", d.Get("policy_template_id").(string))
	if epsId != "" {
		requestPath += fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS container clusters protection policy template: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	templateId := utils.PathSearch("id", respBody, "").(string)
	if templateId == "" {
		return diag.Errorf("error retrieving HSS container clusters protection policy template: ID is not found in API response")
	}

	d.SetId(templateId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("template_name", utils.PathSearch("template_name", respBody, nil)),
		d.Set("template_type", utils.PathSearch("template_type", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("target_kind", utils.PathSearch("target_kind", respBody, nil)),
		d.Set("tag", utils.PathSearch("tag", respBody, nil)),
		d.Set("level", utils.PathSearch("level", respBody, nil)),
		d.Set("constraint_template", utils.PathSearch("constraint_template", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
