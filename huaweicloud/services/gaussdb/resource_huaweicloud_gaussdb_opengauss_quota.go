package gaussdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB PUT /v3/{project_id}/enterprise-projects/quotas
// @API GaussDB GET /v3/{project_id}/enterprise-projects/quotas
func ResourceOpenGaussQuota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussQuotaCreateOrUpdate,
		UpdateContext: resourceOpenGaussQuotaCreateOrUpdate,
		ReadContext:   resourceOpenGaussQuotaRead,
		DeleteContext: resourceOpenGaussQuotaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"vcpus_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ram_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"volume_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vcpus_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceOpenGaussQuotaCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/enterprise-projects/quotas"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildOpenGaussQuotaBodyParams(d))

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error updating GaussDB OpenGauss quota: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(d.Get("enterprise_project_id").(string))
	}

	return resourceOpenGaussQuotaRead(ctx, d, meta)
}

func buildOpenGaussQuotaBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"enterprise_projects_id": d.Get("enterprise_project_id"),
		"instance_quota":         d.Get("instance_quota"),
		"vcpus_quota":            d.Get("vcpus_quota"),
		"ram_quota":              d.Get("ram_quota"),
		"volume_quota":           d.Get("volume_quota"),
	}
	bodyParams := map[string]interface{}{
		"eps_quotas": []map[string]interface{}{params},
	}
	return bodyParams
}

func resourceOpenGaussQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/enterprise-projects/quotas"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB OpenGauss quota")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	expression := fmt.Sprintf("eps_quotas[?enterprise_project_id=='%s']|[0]", d.Id())
	quota := utils.PathSearch(expression, getRespBody, nil)
	if quota == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", quota, nil)),
		d.Set("enterprise_project_name", utils.PathSearch("enterprise_project_name", quota, nil)),
		d.Set("instance_quota", utils.PathSearch("instance_eps_quota", quota, nil)),
		d.Set("vcpus_quota", utils.PathSearch("vcpus_eps_quota", quota, nil)),
		d.Set("ram_quota", utils.PathSearch("ram_eps_quota", quota, nil)),
		d.Set("volume_quota", utils.PathSearch("volume_eps_quota", quota, nil)),
		d.Set("instance_used", utils.PathSearch("instance_used", quota, nil)),
		d.Set("vcpus_used", utils.PathSearch("vcpus_used", quota, nil)),
		d.Set("ram_used", utils.PathSearch("ram_used", quota, nil)),
		d.Set("volume_used", utils.PathSearch("volume_used", quota, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceOpenGaussQuotaDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB OpenGauss quota resource is not supported. The GaussDB OpenGauss quota resource is " +
		"only removed from the state, the GaussDB OpenGauss quota remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
