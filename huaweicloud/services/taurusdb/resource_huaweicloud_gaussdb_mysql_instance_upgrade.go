package taurusdb

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/db-upgrade
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
func ResourceGaussDBMysqlUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBMysqlUpgradeCreate,
		ReadContext:   resourceGaussDBMysqlUpgradeRead,
		DeleteContext: resourceGaussDBMysqlUpgradeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delay": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGaussDBMysqlUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db-upgrade"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBMysqlUpgradeBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB MySQL instance upgrade: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error upgrading GaussDB MySQL instance: job_id is not found in API response")
	}
	err = checkGaussDBMySQLJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for updating GaussDB MySQL instance(%s) job to complete: %s", instanceId,
			err)
	}

	d.SetId(instanceId)

	return nil
}

func buildCreateGaussDBMysqlUpgradeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"delay": d.Get("delay"),
	}
	return bodyParams
}

func resourceGaussDBMysqlUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDBMysqlUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB MySQL instance upgrade resource is not supported. The upgrade resource is only " +
		"removed from the state, the GaussDB MySQL instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
