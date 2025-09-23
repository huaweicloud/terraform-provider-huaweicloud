package dew

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

// @API DEW GET /v1.0/{project_id}/kms/user-quotas
func DataSourceKMSQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKMSQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the quotas.",
									},
									"used": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of quotas used.",
									},
									"quota": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The total number of quotas.",
									},
								},
							},
							Description: "The list of the resource quotas.",
						},
					},
				},
				Description: "The quota details.",
			},
		},
	}
}

func dataSourceKMSQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		mErr   *multierror.Error
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	httpUrl := "v1.0/{project_id}/kms/user-quotas"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.Errorf("error retrieving KMS quotas: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("quotas", flattenQuotasResponseBody(respBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenQuotasResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	resources := utils.PathSearch("quotas.resources", resp, make([]interface{}, 0)).([]interface{})

	return []interface{}{
		map[string]interface{}{
			"resources": flattenQuotasResponseBodyResources(resources),
		},
	}
}

func flattenQuotasResponseBodyResources(resourcesResp []interface{}) []interface{} {
	if len(resourcesResp) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(resourcesResp))
	for _, v := range resourcesResp {
		rst = append(rst, map[string]interface{}{
			"type":  utils.PathSearch("type", v, nil),
			"used":  utils.PathSearch("used", v, nil),
			"quota": utils.PathSearch("quota", v, nil),
		})
	}
	return rst
}
