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

// @API CBH GET /v2/{project_id}/cbs/instance/ecs-quota
func DataSourceInstanceEcsQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceEcsQuotaRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the availability zone name.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the specification code of the CBH instance to be created.`,
			},
			"status_v6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the CBH instance specification resources which support Ipv6.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the CBH instance specification resources.`,
			},
		},
	}
}

func dataSourceInstanceEcsQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		mErr       *multierror.Error
		getHttpUrl = "v2/{project_id}/cbs/instance/ecs-quota"
		product    = "cbh"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?availability_zone=%s&resource_spec_code=%s", getPath,
		d.Get("availability_zone").(string),
		d.Get("resource_spec_code").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving CBH instance ECS quota: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("status_v6", utils.PathSearch("status_v6", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
