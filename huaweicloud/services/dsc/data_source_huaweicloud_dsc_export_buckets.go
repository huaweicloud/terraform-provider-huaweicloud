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

// @API DSC GET /v1/{project_id}/scan-jobs/buckets-for-export
func DataSourceDscExportBuckets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscExportBucketsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"buckets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscExportBucketsSchema(),
			},
		},
	}
}

func dscExportBucketsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"asset_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bind_task": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bucket_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_audit_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceDscExportBucketsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/scan-jobs/buckets-for-export"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC export buckets: %s", err)
	}

	requestRespBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	bucketsList := utils.PathSearch("buckets", requestRespBody, make([]interface{}, 0)).([]interface{})

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("buckets", flattenDscExportBuckets(bucketsList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscExportBuckets(bucketsList []interface{}) []interface{} {
	if len(bucketsList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(bucketsList))
	for _, v := range bucketsList {
		rst = append(rst, map[string]interface{}{
			"asset_name":          utils.PathSearch("asset_name", v, nil),
			"bind_task":           utils.PathSearch("bind_task", v, nil),
			"bucket_location":     utils.PathSearch("bucket_location", v, nil),
			"bucket_name":         utils.PathSearch("bucket_name", v, nil),
			"bucket_policy":       utils.PathSearch("bucket_policy", v, nil),
			"create_time":         utils.PathSearch("create_time", v, nil),
			"enable_audit_status": utils.PathSearch("enable_audit_status", v, nil),
			"id":                  utils.PathSearch("id", v, nil),
			"is_deleted":          utils.PathSearch("is_deleted", v, nil),
		})
	}

	return rst
}
