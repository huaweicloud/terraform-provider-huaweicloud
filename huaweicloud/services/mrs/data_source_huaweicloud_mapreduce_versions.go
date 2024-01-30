// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product MRS
// ---------------------------------------------------------------

package mrs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API MRS GET /v2/{project_id}/metadata/versions
func DataSourceMrsVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceMrsVersionsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"versions": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of available cluster versions",
			},
		},
	}
}

func resourceMrsVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getVersionsHttpUrl = "v2/{project_id}/metadata/versions"
		getVersionsProduct = "mrs"
	)
	getVersionsClient, err := cfg.NewServiceClient(getVersionsProduct, region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	getVersionsPath := getVersionsClient.Endpoint + getVersionsHttpUrl
	getVersionsPath = strings.ReplaceAll(getVersionsPath, "{project_id}", getVersionsClient.ProjectID)

	getVersionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getVersionsResp, err := getVersionsClient.Request("GET", getVersionsPath, &getVersionsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving MrsVersions")
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

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("versions", utils.PathSearch("cluster_versions", getVersionsRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
