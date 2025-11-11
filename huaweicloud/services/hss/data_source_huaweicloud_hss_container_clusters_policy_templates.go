package hss

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

// @API HSS GET /v5/{project_id}/container/clusters/protection-policy-templates
func DataSourceContainerClustersPolicyTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerClustersPolicyTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_kind": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func buildContainerClustersPolicyTemplatesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("template_name"); ok {
		queryParams = fmt.Sprintf("%s&template_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("template_type"); ok {
		queryParams = fmt.Sprintf("%s&template_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("target_kind"); ok {
		queryParams = fmt.Sprintf("%s&target_kind=%v", queryParams, v)
	}
	if v, ok := d.GetOk("tag"); ok {
		queryParams = fmt.Sprintf("%s&tag=%v", queryParams, v)
	}
	if v, ok := d.GetOk("level"); ok {
		queryParams = fmt.Sprintf("%s&level=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceContainerClustersPolicyTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum = 0
		httpUrl  = "v5/{project_id}/container/clusters/protection-policy-templates"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerClustersPolicyTemplatesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS container clusters protection policy templates: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = int(utils.PathSearch("total_num", respBody, float64(0)).(float64))
		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		result = append(result, dataList...)
		offset += len(dataList)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenContainerClustersPolicyTemplatesDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerClustersPolicyTemplatesDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"template_name":       utils.PathSearch("template_name", v, nil),
			"template_type":       utils.PathSearch("template_type", v, nil),
			"description":         utils.PathSearch("description", v, nil),
			"target_kind":         utils.PathSearch("target_kind", v, nil),
			"tag":                 utils.PathSearch("tag", v, nil),
			"level":               utils.PathSearch("level", v, nil),
			"constraint_template": utils.PathSearch("constraint_template", v, nil),
		})
	}

	return rst
}
