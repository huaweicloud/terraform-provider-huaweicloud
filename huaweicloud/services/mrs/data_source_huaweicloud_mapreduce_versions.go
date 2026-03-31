package mrs

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

// @API MRS GET /v2/{project_id}/metadata/versions
func DataSourceVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVersionsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the cluster versions are located.`,
			},
			// Attributes.
			"versions": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The list of available cluster versions.`,
			},
		},
	}
}

func dataSourceVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
		getVersionsHttpUrl = "v2/{project_id}/metadata/versions"
	)

	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	getVersionsPath := client.Endpoint + getVersionsHttpUrl
	getVersionsPath = strings.ReplaceAll(getVersionsPath, "{project_id}", client.ProjectID)

	getVersionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getVersionsResp, err := client.Request("GET", getVersionsPath, &getVersionsOpt)
	if err != nil {
		return diag.Errorf("error retrieving cluster available versions: %s", err)
	}

	getVersionsRespBody, err := utils.FlattenResponse(getVersionsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("versions", utils.PathSearch("cluster_versions", getVersionsRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
