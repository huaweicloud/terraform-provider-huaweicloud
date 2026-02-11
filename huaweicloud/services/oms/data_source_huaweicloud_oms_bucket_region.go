package oms

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API OMS POST /v2/{project_id}/objectstorage/buckets/regions
func DataSourceObjectstorageBucketRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObjectstorageBucketRegionRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ak": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"sk": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"json_auth_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildObjectstorageBucketRegionParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"cloud_type":        d.Get("cloud_type"),
		"bucket_name":       d.Get("bucket_name"),
		"ak":                utils.ValueIgnoreEmpty(d.Get("ak")),
		"sk":                utils.ValueIgnoreEmpty(d.Get("sk")),
		"json_auth_file":    utils.ValueIgnoreEmpty(d.Get("json_auth_file")),
		"connection_string": utils.ValueIgnoreEmpty(d.Get("connection_string")),
		"app_id":            utils.ValueIgnoreEmpty(d.Get("app_id")),
	}
}

func dataSourceObjectstorageBucketRegionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/objectstorage/buckets/regions"
	)

	client, err := cfg.NewServiceClient("oms", region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildObjectstorageBucketRegionParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving the region corresponding to the bucket: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("region_id", utils.PathSearch("id", respBody, nil)),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("support", utils.PathSearch("support", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
