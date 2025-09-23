package ddm

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

// @API DDM GET /v3/{project_id}/instances/{instance_id}/database-version/available-versions
func DataSourceDdmInstanceAvailableVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdmInstanceAvailableVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DDM instance ID.`,
			},
			"current_favored_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the preferred version of the current series.`,
			},
			"current_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the current version.`,
			},
			"latest_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the latest version.`,
			},
			"previous_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the previous version of the current instance.`,
			},
			"versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the available version.`,
			},
		},
	}
}

func dataSourceDdmInstanceAvailableVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/database-version/available-versions"
		product = "ddm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DDM instance available versions")
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
		d.Set("current_favored_version", utils.PathSearch("current_favored_version", getRespBody, nil)),
		d.Set("current_version", utils.PathSearch("current_version", getRespBody, nil)),
		d.Set("latest_version", utils.PathSearch("latest_version", getRespBody, nil)),
		d.Set("previous_version", utils.PathSearch("previous_version", getRespBody, nil)),
		d.Set("versions", utils.PathSearch("versions", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
