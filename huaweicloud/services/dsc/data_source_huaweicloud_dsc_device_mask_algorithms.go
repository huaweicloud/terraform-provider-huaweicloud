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

// @API DSC GET /v1/{project_id}/devices/{device_id}/mask-algorithms
func DataSourceDscDeviceMaskAlgorithms() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscDeviceMaskAlgorithmsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mask_algorithms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscDeviceMaskAlgorithmSchema(),
			},
		},
	}
}

func dscDeviceMaskAlgorithmSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDscDeviceMaskAlgorithmsQueryParams(limit, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func dataSourceDscDeviceMaskAlgorithmsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/devices/{device_id}/mask-algorithms"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Get("device_id").(string))

	for {
		currentPath := requestPath + buildDscDeviceMaskAlgorithmsQueryParams(limit, offset)
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC device mask algorithms: %s", err)
		}

		requestRespBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		maskAlgorithmsList := utils.PathSearch("mask_algorithms", requestRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, maskAlgorithmsList...)
		if len(maskAlgorithmsList) < limit {
			break
		}
		offset += len(maskAlgorithmsList)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("mask_algorithms", flattenDscDeviceMaskAlgorithms(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscDeviceMaskAlgorithms(maskAlgorithmsList []interface{}) []interface{} {
	if len(maskAlgorithmsList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(maskAlgorithmsList))
	for _, v := range maskAlgorithmsList {
		rst = append(rst, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"name": utils.PathSearch("name", v, nil),
		})
	}

	return rst
}
