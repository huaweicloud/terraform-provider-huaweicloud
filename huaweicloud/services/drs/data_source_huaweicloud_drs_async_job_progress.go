package drs

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

// @API DRS GET /v5/{project_id}/jobs/{async_job_id}/progress
func DataSourceDrsAsyncJobProgress() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsAsyncJobProgressRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source. If omitted, the provider-level region will be used.",
			},
			"async_job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the unique identifier of the async task.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type of the async task to query.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the async task.",
			},
			"progress": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The execution progress of the async task.",
			},
			"error_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The error code of the async task.",
			},
			"error_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The error message of the async task.",
			},
		},
	}
}

func dataSourceDrsAsyncJobProgressRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/{async_job_id}/progress"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{async_job_id}", d.Get("async_job_id").(string))
	getPath += buildListAsyncJobProgressQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS async job progress: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("progress", utils.PathSearch("progress", getRespBody, nil)),
		d.Set("error_code", utils.PathSearch("error_code", getRespBody, nil)),
		d.Set("error_message", utils.PathSearch("error_message", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListAsyncJobProgressQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
