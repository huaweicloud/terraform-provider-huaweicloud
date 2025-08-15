package workspace

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

// @API Workspace GET /v1/{project_id}/check/quota
func DataSourceAppServerQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppServerQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the Workspace APP server quotas are located.`,
			},
			"product_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the product to be queried.`,
			},
			"subscription_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The number of server instances to be queried.`,
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The disk size of the single server instance to be queried.`,
			},
			"disk_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The number of disks for the single server instance to be queried.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the flavor to be queried.`,
			},
			"is_period": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the instance is prepaid.`,
			},
			"deh_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the dedicated host.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The of the cloud dedicated distributed storage pool.`,
			},
			"is_enough": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the quota is sufficient.`,
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The quota resource type.`,
						},
						"remainder": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The remaining quota.`,
						},
						"need": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The required quota.`,
						},
					},
				},
				Description: `The list of the quotas that match the filter parameters.`,
			},
		},
	}
}

func buildAppServerQuotasQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?product_id=%v&subscription_num=%v&disk_size=%v&disk_num=%v",
		d.Get("product_id"),
		d.Get("subscription_num"),
		d.Get("disk_size"),
		d.Get("disk_num"),
	)

	if v, ok := d.GetOk("flavor_id"); ok {
		res = fmt.Sprintf("%s&flavor_id=%v", res, v)
	}

	if v := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "is_period"); v != nil {
		res = fmt.Sprintf("%s&is_period=%v", res, v)
	}

	if v, ok := d.GetOk("deh_id"); ok {
		res = fmt.Sprintf("%s&deh_id=%v", res, v)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		res = fmt.Sprintf("%s&cluster_id=%v", res, v)
	}

	return res
}

func dataSourceAppServerQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/check/quota"
		opt     = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}
	)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildAppServerQuotasQueryParams(d)

	resp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.FromErr(err)
	}

	respBodoy, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error querying Workspace APP server quotas: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("is_enough", utils.PathSearch("is_enough", respBodoy, nil)),
		d.Set("quotas", flattenAppServerQuotas(utils.PathSearch("quota_remainder",
			respBodoy, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAppServerQuotas(quotas []interface{}) []map[string]interface{} {
	if len(quotas) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(quotas))
	for i, item := range quotas {
		result[i] = map[string]interface{}{
			"type":      utils.PathSearch("type", item, nil),
			"remainder": utils.PathSearch("remainder", item, nil),
			"need":      utils.PathSearch("need", item, nil),
		}
	}

	return result
}
