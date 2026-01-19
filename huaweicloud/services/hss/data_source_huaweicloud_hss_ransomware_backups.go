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

// @API HSS GET /v5/{project_id}/ransomware/backup/{host_id}
func DataSourceRansomwareBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRansomwareBackupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceRansomwareBackupsDataListSchema(),
			},
		},
	}
}

func dataSourceRansomwareBackupsDataListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// In the API documentation, it is `backup_status`, but it actually returns `status`.
			"backup_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// `create_time` in API documentation as an int type, but it actually returns a string type.
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_images_data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"backup_tag": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildRansomwareBackupsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("last_days"); ok {
		queryParams = fmt.Sprintf("%s&last_days=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceRansomwareBackupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
		httpUrl  = "v5/{project_id}/ransomware/backup/{host_id}"
		hostId   = d.Get("host_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{host_id}", hostId)
	requestPath += buildRansomwareBackupsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS ransomware backups: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
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
		d.Set("data_list", flattenRansomwareBackupsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRansomwareBackupsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"backup_id":     utils.PathSearch("backup_id", v, nil),
			"backup_name":   utils.PathSearch("backup_name", v, nil),
			"backup_status": utils.PathSearch("status", v, nil),
			"create_time":   utils.PathSearch("create_time", v, nil),
			"os_images_data": flattenRansomwareBackupsOsImagesData(
				utils.PathSearch("os_images_data", v, make([]interface{}, 0)).([]interface{})),
			"backup_tag": utils.PathSearch("backup_tag", v, nil),
		})
	}

	return rst
}

func flattenRansomwareBackupsOsImagesData(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"image_id": utils.PathSearch("image_id", v, nil),
		})
	}

	return rst
}
