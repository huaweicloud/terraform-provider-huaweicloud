package lakeformation

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

// @API LakeFormation GET /v1/{project_id}/specs
func DataSourceSpecifications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSpecificationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the specifications are located.`,
			},

			// Optional parameters.
			"spec_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The code of the specification to be queried.`,
			},

			// Attributes.
			"specifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The code of the specification.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the specification.`,
						},
						"stride": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The stride of the specification.`,
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unit of the specification.`,
						},
						"min_stride_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The minimum stride number of the specification.`,
						},
						"max_stride_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum value of the specification.`,
						},
						"usage_measure_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The usage measurement unit ID of the specification.`,
						},
						"usage_factor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the specification.`,
						},
						"usage_value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The product ID of the specification.`,
						},
						"free_usage_value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The free usage value of the specification.`,
						},
						"stride_num_whitelist": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The stride number whitelist of the specification.`,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
				Description: `The list of specifications that matched filter parameters.`,
			},
		},
	}
}

func listSpecifications(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/specs"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("spec_codes", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenSpecifications(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"spec_code":            utils.PathSearch("spec_code", item, nil),
			"resource_type":        utils.PathSearch("resource_type", item, nil),
			"stride":               utils.PathSearch("stride", item, nil),
			"unit":                 utils.PathSearch("unit", item, nil),
			"min_stride_num":       utils.PathSearch("min_stride_num", item, nil),
			"max_stride_num":       utils.PathSearch("max_stride_num", item, nil),
			"usage_measure_id":     utils.PathSearch("usage_measure_id", item, nil),
			"usage_factor":         utils.PathSearch("usage_factor", item, nil),
			"usage_value":          utils.PathSearch("usage_value", item, nil),
			"free_usage_value":     utils.PathSearch("free_usage_value", item, nil),
			"stride_num_whitelist": utils.PathSearch("stride_num_whitelist", item, nil),
		})
	}

	return result
}

func filterSpecifications(all []interface{}, d *schema.ResourceData) []interface{} {
	result := make([]interface{}, 0, len(all))

	for _, v := range all {
		if param, ok := d.GetOk("spec_code"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("spec_code", v, nil)) {
			continue
		}

		result = append(result, v)
	}
	return result
}

func dataSourceSpecificationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("lakeformation", region)
	if err != nil {
		return diag.Errorf("error creating LakeFormation client: %s", err)
	}

	specifications, err := listSpecifications(client)
	if err != nil {
		return diag.Errorf("error querying instance specifications: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("specifications", flattenSpecifications(filterSpecifications(specifications, d))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
