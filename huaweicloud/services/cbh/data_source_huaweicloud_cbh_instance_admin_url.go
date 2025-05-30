package cbh

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBH GET /v2/{project_id}/cbs/instances/{server_id}/admin-url
func DataSourceInstanceAdminUrl() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceAdminUrlRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CBH instance ID.`,
			},
			"admin_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URL for logging in to a CBH instance as user admin.`,
			},
		},
	}
}

func dataSourceInstanceAdminUrlRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		mErr       *multierror.Error
		getHttpUrl = "v2/{project_id}/cbs/instances/{server_id}/admin-url"
		product    = "cbh"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{server_id}", d.Get("server_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving CBH instance admin url: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("server_id").(string))

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("admin_url", utils.PathSearch("admin_url", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
