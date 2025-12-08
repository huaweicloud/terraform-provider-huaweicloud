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

// @API HSS GET /v5/{project_id}/image/{image_id}/files-statistics
func DataSourceImageFilesStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageFilesStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_files_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_files_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildImageFilesStatisticsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?image_type=%s", d.Get("image_type"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	if v, ok := d.GetOk("namespace"); ok {
		queryParams = fmt.Sprintf("%s&namespace=%v", queryParams, v)
	}

	if v, ok := d.GetOk("image_name"); ok {
		queryParams = fmt.Sprintf("%s&image_name=%v", queryParams, v)
	}

	if v, ok := d.GetOk("tag_name"); ok {
		queryParams = fmt.Sprintf("%s&tag_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceImageFilesStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/image/{image_id}/files-statistics"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	imageId := d.Get("image_id").(string)
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{image_id}", imageId)
	requestPath += buildImageFilesStatisticsQueryParams(d, epsId)

	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS image files statistics: %s", err)
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
		d.Set("total_files_num", utils.PathSearch("total_files_num", respBody, nil)),
		d.Set("total_files_size", utils.PathSearch("total_files_size", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
