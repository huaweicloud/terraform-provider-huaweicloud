package cbh

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

// @API CBH GET /v2/{project_id}/cbs/instance/get-om-url
func DataSourceInstanceOmUrl() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceOmUrlRead,

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
				Description: `Specifies the ID of the CBH instance.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the IP address of the managed host.`,
			},
			"host_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the account name of the managed host.`,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The state of getting the OM URL.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description when failed to get the OM URL.`,
			},
			"login_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OM login URL.`,
			},
		},
	}
}

func buildInstanceOmUrlQueryParams(d *schema.ResourceData) string {
	var (
		serverId        = d.Get("server_id").(string)
		ipAddress       = d.Get("ip_address").(string)
		hostAccountName = d.Get("host_account_name").(string)
	)

	return fmt.Sprintf("?server_id=%s&ip_address=%s&host_account_name=%s", serverId, ipAddress, hostAccountName)
}

func dataSourceInstanceOmUrlRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/cbs/instance/get-om-url"
		product = "cbh"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildInstanceOmUrlQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CBH instance OM URL: %s", err)
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
		d.Set("state", utils.PathSearch("state", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("login_url", utils.PathSearch("login_url", respBody, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CBH instance OM URL fields: %s", err)
	}

	return nil
}
