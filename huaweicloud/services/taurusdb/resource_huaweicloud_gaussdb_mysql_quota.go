package taurusdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL POST /v3/{project_id}/quotas
// @API GaussDBforMySQL GET /v3/{project_id}/quotas
func ResourceGaussDBMysqlQuota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBMysqlQuotaCreateOrUpdate,
		UpdateContext: resourceGaussDBMysqlQuotaCreateOrUpdate,
		ReadContext:   resourceGaussDBMysqlQuotaRead,
		DeleteContext: resourceGaussDBMysqlQuotaDelete,
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
				Default:  -1,
			},
			"vcpus_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"ram_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"enterprise_project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_instance_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"availability_vcpus_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"availability_ram_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceGaussDBMysqlQuotaCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/quotas"
		product = "gaussdb"
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
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBMysqlQuotaBodyParams(d))

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error updating GaussDB MySQL quota: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(d.Get("enterprise_project_id").(string))
	}

	return resourceGaussDBMysqlQuotaRead(ctx, d, meta)
}

func buildCreateGaussDBMysqlQuotaBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"enterprise_project_id": d.Get("enterprise_project_id"),
		"instance_quota":        d.Get("instance_quota"),
		"vcpus_quota":           d.Get("vcpus_quota"),
		"ram_quota":             d.Get("ram_quota"),
	}
	bodyParams := map[string]interface{}{
		"quota_list": []map[string]interface{}{params},
	}
	return bodyParams
}

func resourceGaussDBMysqlQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/quotas"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listMysqlDatabasesResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.FromErr(err)
	}

	listRespJson, err := json.Marshal(listMysqlDatabasesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}
	expression := fmt.Sprintf("quota_list[?enterprise_project_id=='%s']|[0]", d.Id())
	quota := utils.PathSearch(expression, listRespBody, nil)
	if quota == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", quota, nil)),
		d.Set("enterprise_project_name", utils.PathSearch("enterprise_project_name", quota, nil)),
		d.Set("instance_quota", utils.PathSearch("instance_quota", quota, nil)),
		d.Set("vcpus_quota", utils.PathSearch("vcpus_quota", quota, nil)),
		d.Set("ram_quota", utils.PathSearch("ram_quota", quota, nil)),
		d.Set("availability_instance_quota", utils.PathSearch("availability_instance_quota", quota, nil)),
		d.Set("availability_vcpus_quota", utils.PathSearch("availability_vcpus_quota", quota, nil)),
		d.Set("availability_ram_quota", utils.PathSearch("availability_ram_quota", quota, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGaussDBMysqlQuotaDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB MySQL quota resource is not supported. The GaussDB MySQL quota resource is only " +
		"removed from the state, the GaussDB MySQL quota remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
