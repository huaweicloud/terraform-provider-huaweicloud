package gaussdb

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/sync-limit-task
func ResourceOpenGaussSyncSqlThrottlingTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussSyncSqlThrottlingTaskCreate,
		ReadContext:   resourceOpenGaussSyncSqlThrottlingTaskRead,
		DeleteContext: resourceOpenGaussSyncSqlThrottlingTaskDelete,

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
		},
	}
}

func resourceOpenGaussSyncSqlThrottlingTaskCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/sync-limit-task"
		product = "opengauss"
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

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB OpenGauss(%s) sync SQL throttling task: %s", instanceId, err)
	}

	d.SetId(instanceId)

	return nil
}

func resourceOpenGaussSyncSqlThrottlingTaskRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOpenGaussSyncSqlThrottlingTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB OpenGauss sync SQL throttling task resource is not supported. The GaussDB OpenGauss" +
		" sync SQL throttling task resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
