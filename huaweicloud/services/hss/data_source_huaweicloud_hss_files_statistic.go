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

// @API HSS GET /v5/{project_id}/files/statistic
func DataSourceFilesStatistic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFilesStatisticRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"change_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"change_file_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"change_registry_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"modify_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"add_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"delete_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildFilesStatisticQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	if v, ok := d.GetOk("begin_time"); ok {
		queryParams = fmt.Sprintf("%s&begin_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceFilesStatisticRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/files/statistic"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildFilesStatisticQueryParams(d, epsId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the server files statistic: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("host_num", utils.PathSearch("host_num", respBody, nil)),
		d.Set("change_total_num", utils.PathSearch("change_total_num", respBody, nil)),
		d.Set("change_file_num", utils.PathSearch("change_file_num", respBody, nil)),
		d.Set("change_registry_num", utils.PathSearch("change_registry_num", respBody, nil)),
		d.Set("modify_num", utils.PathSearch("modify_num", respBody, nil)),
		d.Set("add_num", utils.PathSearch("add_num", respBody, nil)),
		d.Set("delete_num", utils.PathSearch("delete_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
