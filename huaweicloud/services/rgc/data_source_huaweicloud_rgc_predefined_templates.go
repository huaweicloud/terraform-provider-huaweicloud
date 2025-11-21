package rgc

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

// @API RGC GET /v1/rgc/predefined-templates
func DataSourcePredefinedTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePredefinedTemplatesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     PredefinedTemplatesSchema(),
			},
		},
	}
}

func PredefinedTemplatesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourcePredefinedTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getPredefinedTemplatesProduct = "rgc"
	getPredefinedTemplatesClient, err := cfg.NewServiceClient(getPredefinedTemplatesProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getPredefinedTemplatesRespBody, err := getPredefinedTemplates(getPredefinedTemplatesClient)

	if err != nil {
		return diag.Errorf("error retrieving RGC pre-defined templates: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("templates", utils.PathSearch("templates", getPredefinedTemplatesRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getPredefinedTemplates(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getPredefinedTemplatesHttpUrl = "v1/rgc/predefined-templates"
	)
	getPredefinedTemplatesHttpPath := client.Endpoint + getPredefinedTemplatesHttpUrl

	getPredefinedTemplatesHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPredefinedTemplatesHttpResp, err := client.Request("GET", getPredefinedTemplatesHttpPath, &getPredefinedTemplatesHttpOpt)
	if err != nil {
		return nil, err
	}
	getPredefinedTemplatesRespBody, err := utils.FlattenResponse(getPredefinedTemplatesHttpResp)
	if err != nil {
		return nil, err
	}
	return getPredefinedTemplatesRespBody, nil
}
