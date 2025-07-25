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

// @API HSS GET /v5/{project_id}/antivirus/statistic
func DataSourceAntivirusStatistic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAntivirusStatisticRead,

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
			"total_malware_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"malware_host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_task_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"scanning_task_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"latest_scan_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"scan_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAntivirusStatisticRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/antivirus/statistic"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		requestPath += fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS antivirus statistic: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_malware_num", utils.PathSearch("total_malware_num", respBody, nil)),
		d.Set("malware_host_num", utils.PathSearch("malware_host_num", respBody, nil)),
		d.Set("total_task_num", utils.PathSearch("total_task_num", respBody, nil)),
		d.Set("scanning_task_num", utils.PathSearch("scanning_task_num", respBody, nil)),
		d.Set("latest_scan_time", utils.PathSearch("latest_scan_time", respBody, nil)),
		d.Set("scan_type", utils.PathSearch("scan_type", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
