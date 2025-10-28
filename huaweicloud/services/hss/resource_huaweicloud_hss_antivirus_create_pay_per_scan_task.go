package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var antivirusCreatePayPerScanTaskNonUpdatableParams = []string{
	"task_name",
	"scan_type",
	"action",
	"host_ids",
	"file_types",
	"scan_dir",
	"ignore_dir",
	"task_id",
	"enterprise_project_id",
}

// @API HSS POST /v5/{project_id}/antivirus/pay-per-scan/tasks
func ResourceAntivirusCreatePayPerScanTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAntivirusCreatePayPerScanTaskCreate,
		ReadContext:   resourceAntivirusCreatePayPerScanTaskRead,
		UpdateContext: resourceAntivirusCreatePayPerScanTaskUpdate,
		DeleteContext: resourceAntivirusCreatePayPerScanTaskDelete,

		CustomizeDiff: config.FlexibleForceNew(antivirusCreatePayPerScanTaskNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scan_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"file_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"scan_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ignore_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildAntivirusCreatePayPerScanTaskQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func buildAntivirusCreatePayPerScanTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"task_name":  d.Get("task_name"),
		"scan_type":  d.Get("scan_type"),
		"action":     d.Get("action"),
		"host_ids":   utils.ExpandToStringList(d.Get("host_ids").([]interface{})),
		"file_types": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("file_types").([]interface{}))),
		"scan_dir":   utils.ValueIgnoreEmpty(d.Get("scan_dir")),
		"ignore_dir": utils.ValueIgnoreEmpty(d.Get("ignore_dir")),
		"task_id":    utils.ValueIgnoreEmpty(d.Get("task_id")),
	}

	return bodyParams
}

func resourceAntivirusCreatePayPerScanTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/antivirus/pay-per-scan/tasks"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAntivirusCreatePayPerScanTaskQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
		JSONBody:         utils.RemoveNil(buildAntivirusCreatePayPerScanTaskBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating HSS antivirus pay-per-scan task: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceAntivirusCreatePayPerScanTaskRead(ctx, d, meta)
}

func resourceAntivirusCreatePayPerScanTaskRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceAntivirusCreatePayPerScanTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceAntivirusCreatePayPerScanTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to create HSS antivirus pay-per-scan task. Deleting
    this resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
