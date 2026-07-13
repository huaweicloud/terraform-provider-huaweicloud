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

// @API DSC GET /v2/{project_id}/sec-ops/situation-dashboard/metadata
func DataSourceDscShowMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscShowMetadataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"column_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"file_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sensitive_column_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sensitive_file_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"table_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDscShowMetadataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v2/{project_id}/sec-ops/situation-dashboard/metadata"
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
		return diag.Errorf("error retrieving DSC show metadata: %s", err)
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
		d.Set("column_num", utils.PathSearch("column_num", respBody, nil)),
		d.Set("file_num", utils.PathSearch("file_num", respBody, nil)),
		d.Set("sensitive_column_num", utils.PathSearch("sensitive_column_num", respBody, nil)),
		d.Set("sensitive_file_num", utils.PathSearch("sensitive_file_num", respBody, nil)),
		d.Set("table_num", utils.PathSearch("table_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
