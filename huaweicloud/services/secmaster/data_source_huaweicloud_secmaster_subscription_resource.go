package secmaster

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

// @API Secmaster GET /v2/{project_id}/workspaces/{workspace_id}/siem/subscription/resource
func DataSourceSecmasterSubscriptionResource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterSubscriptionResourceRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sku": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sku_attribute": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upper_limit": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"step": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"used_amount": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"unused_amount": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"index_storage_upper_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"index_shards_upper_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"index_shards_unused": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"partitions_unused": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"partition_upper_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildSubscriptionResourceQueryParams(d *schema.ResourceData) string {
	if sku, ok := d.GetOk("sku"); ok {
		return "?sku=" + sku.(string)
	}
	return ""
}

func dataSourceSecmasterSubscriptionResourceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/workspaces/{workspace_id}/siem/subscription/resource"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath += buildSubscriptionResourceQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster subscription resource: %s", err)
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

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("sku_attribute", utils.PathSearch("sku", respBody, nil)),
		d.Set("upper_limit", utils.PathSearch("upper_limit", respBody, nil)),
		d.Set("unit", utils.PathSearch("unit", respBody, nil)),
		d.Set("step", utils.PathSearch("step", respBody, nil)),
		d.Set("used_amount", utils.PathSearch("used_amount", respBody, nil)),
		d.Set("unused_amount", utils.PathSearch("unused_amount", respBody, nil)),
		d.Set("version", utils.PathSearch("version", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
		d.Set("index_storage_upper_limit", utils.PathSearch("index_storage_upper_limit", respBody, nil)),
		d.Set("index_shards_upper_limit", utils.PathSearch("index_shards_upper_limit", respBody, nil)),
		d.Set("index_shards_unused", utils.PathSearch("index_shards_unused", respBody, nil)),
		d.Set("partitions_unused", utils.PathSearch("partitions_unused", respBody, nil)),
		d.Set("partition_upper_limit", utils.PathSearch("partition_upper_limit", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
