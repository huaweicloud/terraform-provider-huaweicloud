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

// @API HSS GET /v5/{project_id}/image/baseline/statistic
func DataSourceImageBaselineStatistic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageBaselineStatisticRead,

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
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_weak_pwd": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pwd_policy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"security_check": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildImageBaselineStatisticQueryParams(d *schema.ResourceData, epsId string) string {
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
	if v, ok := d.GetOk("image_version"); ok {
		queryParams = fmt.Sprintf("%s&image_version=%v", queryParams, v)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		queryParams = fmt.Sprintf("%s&instance_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("image_id"); ok {
		queryParams = fmt.Sprintf("%s&image_id=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceImageBaselineStatisticRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/image/baseline/statistic"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildImageBaselineStatisticQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS image baseline statistic: %s", err)
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
		d.Set("image_weak_pwd", utils.PathSearch("image_weak_pwd", respBody, nil)),
		d.Set("pwd_policy", utils.PathSearch("pwd_policy", respBody, nil)),
		d.Set("security_check", utils.PathSearch("security_check", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
