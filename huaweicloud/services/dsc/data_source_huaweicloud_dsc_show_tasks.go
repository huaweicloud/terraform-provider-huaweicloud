package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v2/{project_id}/sec-ops/situation-dashboard/task-statistics
func DataSourceDscShowTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscShowTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mask_task_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"scan_task_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"watermark_embed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"watermark_extract_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"watermark_task_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDscShowTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v2/{project_id}/sec-ops/situation-dashboard/task-statistics"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving DSC show tasks: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("mask_task_num", utils.PathSearch("mask_task_num", respBody, nil)),
		d.Set("scan_task_num", utils.PathSearch("scan_task_num", respBody, nil)),
		d.Set("watermark_embed_num", utils.PathSearch("watermark_embed_num", respBody, nil)),
		d.Set("watermark_extract_num", utils.PathSearch("watermark_extract_num", respBody, nil)),
		d.Set("watermark_task_num", utils.PathSearch("watermark_task_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
