package taurusdb

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL DELETE /v3/{project_id}/instance/{instance_id}/scheduled-jobs
func ResourceGaussDBScheduledTaskDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBMysqlScheduledTaskDeleteCreate,
		ReadContext:   resourceGaussDBMysqlScheduledTaskDeleteRead,
		DeleteContext: resourceGaussDBMysqlScheduledTaskDeleteDelete,

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
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGaussDBMysqlScheduledTaskDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instance/{instance_id}/scheduled-jobs"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	jobId := d.Get("job_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBMySQLScheduledTaskDeleteBodyParams(jobId))

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error deleting GaussDB MySQL(%s) scheduled task(%s): %s", instanceId, jobId, err)
	}

	d.SetId(jobId)

	return nil
}

func buildCreateGaussDBMySQLScheduledTaskDeleteBodyParams(jobId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"job_id": jobId,
	}
	return bodyParams
}

func resourceGaussDBMysqlScheduledTaskDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDBMysqlScheduledTaskDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB MySQL scheduled task delete resource is not supported. The GaussDB MySQL scheduled " +
		"task delete resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
