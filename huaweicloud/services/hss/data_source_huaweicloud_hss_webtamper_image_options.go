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

// @API HSS GET /v5/{project_id}/webtamper/image-options
func DataSourceWebTamperImageOptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebTamperImageOptionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"registry_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"registry_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
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
				Elem:     dataSourceWebTamperImageOptionsDataListSchema(),
			},
		},
	}
}

func dataSourceWebTamperImageOptionsDataListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_full_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_version_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"image_namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registry_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registry_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildWebTamperImageOptionsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"
	queryParams = fmt.Sprintf("%s&image_type=%v", queryParams, d.Get("image_type"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("registry_type"); ok {
		queryParams = fmt.Sprintf("%s&registry_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("image_namespace"); ok {
		queryParams = fmt.Sprintf("%s&image_namespace=%v", queryParams, v)
	}
	if v, ok := d.GetOk("registry_name"); ok {
		queryParams = fmt.Sprintf("%s&registry_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("image_name"); ok {
		queryParams = fmt.Sprintf("%s&image_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceWebTamperImageOptionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
		httpUrl  = "v5/{project_id}/webtamper/image-options"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildWebTamperImageOptionsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS web tamper image options: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		dataListResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataListResp) == 0 {
			break
		}

		result = append(result, dataListResp...)
		offset += len(dataListResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenWebTamperImageOptionsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWebTamperImageOptionsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"image_name":      utils.PathSearch("image_name", v, nil),
			"image_full_name": utils.PathSearch("image_full_name", v, nil),
			"image_id":        utils.PathSearch("image_id", v, nil),
			"image_version_list": utils.ExpandToStringList(
				utils.PathSearch("image_version_list", v, make([]interface{}, 0)).([]interface{})),
			"image_namespace": utils.PathSearch("image_namespace", v, nil),
			"registry_name":   utils.PathSearch("registry_name", v, nil),
			"registry_type":   utils.PathSearch("registry_type", v, nil),
		})
	}

	return rst
}
