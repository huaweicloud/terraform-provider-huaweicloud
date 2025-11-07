package fgs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph GET /v2/{project_id}/fgs/service-trusted-agencies
func DataSourceServiceTrustedAgencies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceTrustedAgenciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the service trusted agencies are located.`,
			},
			"agencies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the trusted agency.`,
						},
						// "null" means it never expires
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The expiration time of the trusted agency.`,
						},
					},
				},
				Description: `The list of service trusted agencies.`,
			},
		},
	}
}

func getServiceTrustedAgencies(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/fgs/service-trusted-agencies"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("[]", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func dataSourceServiceTrustedAgenciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	agencies, err := getServiceTrustedAgencies(client)
	if err != nil {
		return diag.Errorf("error querying service trusted agencies: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("agencies", flattenServiceTrustedAgencies(agencies)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenServiceTrustedAgencies(agencies []interface{}) []map[string]interface{} {
	if len(agencies) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(agencies))
	for _, agency := range agencies {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", agency, nil),
			"expire_time": utils.PathSearch("expire_time", agency, nil),
		})
	}
	return result
}
