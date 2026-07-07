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

// @API DSC GET /v1/{project_id}/sdg/asset/quota
func DataSourceDscAssetQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscAssetQuotaRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"total_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"use_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDscAssetQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/sdg/asset/quota"
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
		return diag.Errorf("error retrieving DSC asset quota: %s", err)
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
		d.Set("total_size", utils.PathSearch("total_size", respBody, nil)),
		d.Set("use_size", utils.PathSearch("use_size", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
