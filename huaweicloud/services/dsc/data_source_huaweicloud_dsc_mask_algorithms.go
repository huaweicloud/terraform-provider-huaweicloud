package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v2/{project_id}/mask/algorithms
func DataSourceMaskAlgorithms() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMaskAlgorithmsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"algorithms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     maskAlgorithmsSchema(),
			},
		},
	}
}

func maskAlgorithmsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"algorithm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"algorithm_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"algorithm_name_en": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"processed_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildMaskAlgorithmsQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d", limit)
	// In the API, this field is designated as the page number.
	if offset > 0 {
		queryParams = fmt.Sprintf("%s&offset=%v", queryParams, offset)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams += fmt.Sprintf("&name=%v", v)
	}

	if v, ok := d.GetOk("type"); ok {
		queryParams += fmt.Sprintf("&type=%v", v)
	}

	return queryParams
}

func dataSourceMaskAlgorithmsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v2/{project_id}/mask/algorithms"
		product       = "dsc"
		limit         = 1000
		offset        = 0
		allAlgorithms = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildMaskAlgorithmsQueryParams(d, limit, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC mask algorithms: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		algorithmsResp := utils.PathSearch("algorithms", respBody, make([]interface{}, 0)).([]interface{})
		allAlgorithms = append(allAlgorithms, algorithmsResp...)
		if len(algorithmsResp) < limit {
			break
		}

		offset++
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("algorithms", flattenMaskAlgorithms(allAlgorithms)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMaskAlgorithms(algorithms []interface{}) []interface{} {
	if len(algorithms) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(algorithms))
	for _, v := range algorithms {
		rst = append(rst, map[string]interface{}{
			"algorithm":         utils.PathSearch("algorithm", v, nil),
			"algorithm_id":      utils.PathSearch("algorithm_id", v, nil),
			"algorithm_name":    utils.PathSearch("algorithm_name", v, nil),
			"algorithm_name_en": utils.PathSearch("algorithm_name_en", v, nil),
			"category":          utils.PathSearch("category", v, nil),
			"data":              utils.PathSearch("data", v, nil),
			"parameter":         utils.PathSearch("parameter", v, nil),
			"processed_data":    utils.PathSearch("processed_data", v, nil),
		})
	}

	return rst
}
