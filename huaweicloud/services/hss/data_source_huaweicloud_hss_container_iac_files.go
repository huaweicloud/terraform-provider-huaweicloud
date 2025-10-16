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

// @API HSS GET /v5/{project_id}/container/iac/files
func DataSourceContainerIacFiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerIacFilesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scan_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"risky": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cicd_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cicd_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"risky": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"first_scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"upload_file_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cicd_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cicd_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scan_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerIacFilesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"
	queryParams = fmt.Sprintf("%s&scan_type=%v", queryParams, d.Get("scan_type"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("file_id"); ok {
		queryParams = fmt.Sprintf("%s&file_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("file_name"); ok {
		queryParams = fmt.Sprintf("%s&file_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("file_type"); ok {
		queryParams = fmt.Sprintf("%s&file_type=%v", queryParams, v)
	}

	queryParams = fmt.Sprintf("%s&risky=%v", queryParams, d.Get("risky").(bool))

	if v, ok := d.GetOk("cicd_id"); ok {
		queryParams = fmt.Sprintf("%s&cicd_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cicd_name"); ok {
		queryParams = fmt.Sprintf("%s&cicd_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceContainerIacFilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/container/iac/files"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerIacFilesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		getResp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS container IAC files: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenContainerIacFiles(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerIacFiles(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"file_id":          utils.PathSearch("file_id", v, nil),
			"file_name":        utils.PathSearch("file_name", v, nil),
			"file_type":        utils.PathSearch("file_type", v, nil),
			"risky":            utils.PathSearch("risky", v, nil),
			"risk_num":         utils.PathSearch("risk_num", v, nil),
			"first_scan_time":  utils.PathSearch("first_scan_time", v, nil),
			"last_scan_time":   utils.PathSearch("last_scan_time", v, nil),
			"upload_file_time": utils.PathSearch("upload_file_time", v, nil),
			"cicd_id":          utils.PathSearch("cicd_id", v, nil),
			"cicd_name":        utils.PathSearch("cicd_name", v, nil),
			"scan_type":        utils.PathSearch("scan_type", v, nil),
		})
	}

	return rst
}
