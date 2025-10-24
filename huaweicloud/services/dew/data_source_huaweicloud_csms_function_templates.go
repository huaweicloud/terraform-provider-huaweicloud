package dew

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

// @API DEW GET /v1/csms/function-templates
func DataSourceCsmsFunctionTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCsmsFunctionTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"secret_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the secret type.`,
			},
			"secret_sub_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the secret rotation account type.`,
			},
			"engine": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the database type.`,
			},
			"function_templates": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The secret rotation function templates.`,
			},
		},
	}
}

func buildCsmsFunctionTemplatesQueryParams(d *schema.ResourceData) string {
	rst := fmt.Sprintf("?secret_type=%v&secret_sub_type=%v", d.Get("secret_type"), d.Get("secret_sub_type"))

	if engine, ok := d.GetOk("engine"); ok {
		rst += fmt.Sprintf("&engine=%v", engine)
	}
	return rst
}

func dataSourceCsmsFunctionTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		getHttpUrl = "v1/csms/function-templates"
		product    = "kms"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + getHttpUrl
	requestPath += buildCsmsFunctionTemplatesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DEW CSMS function templates: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("function_templates", utils.PathSearch("function_templates", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
