package rds

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dbKernelUpdateNonUpdatableParams = []string{"instance_id"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/db-upgrade
func ResourceRdsUpgradingMinorVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsUpgradingMinorVersionCreate,
		ReadContext:   resourceRdsUpgradingMinorVersionRead,
		UpdateContext: resourceRdsUpgradingMinorVersionUpdate,
		DeleteContext: resourceRdsUpgradingMinorVersionDelete,

		CustomizeDiff: config.FlexibleForceNew(dbKernelUpdateNonUpdatableParams),

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
			"is_delayed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceRdsUpgradingMinorVersionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	const (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db-upgrade"
		product = "rds"
	)

	instanceID := d.Get("instance_id").(string)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	url := client.Endpoint + httpUrl
	url = strings.ReplaceAll(url, "{project_id}", client.ProjectID)
	url = strings.ReplaceAll(url, "{instance_id}", instanceID)

	body := map[string]interface{}{}
	if v, ok := d.GetOk("is_delayed"); ok {
		body["is_delayed"] = v.(bool)
	}

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	opts.JSONBody = utils.RemoveNil(body)

	resp, err := client.Request("POST", url, &opts)
	if err != nil {
		return diag.Errorf("error upgrading kernel minor version for instance (%s): %s", instanceID, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing upgrade kernel response: %s", err)
	}

	jobID := utils.PathSearch("job_id", respBody, nil)
	if jobID == nil {
		return diag.Errorf("job_id not found in response")
	}

	d.SetId(instanceID)
	_ = d.Set("job_id", jobID)

	return nil
}

func resourceRdsUpgradingMinorVersionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsUpgradingMinorVersionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsUpgradingMinorVersionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting kernel upgrade resource is not supported. This resource is only removed from the state.",
		},
	}
}
