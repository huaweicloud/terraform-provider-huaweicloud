package cfw

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var exportLogsNonUpdatableParams = []string{"fw_instance_id", "start_time", "end_time", "log_type", "type",
	"filters", "filters.*.field", "filters.*.operator", "filters.*.values", "time_zone", "export_file_name"}

// @API CFW POST /v1/{project_id}/cfw/{fw_instance_id}/logs/export
func ResourceExportLogs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExportLogsCreate,
		ReadContext:   resourceExportLogsRead,
		UpdateContext: resourceExportLogsUpdate,
		DeleteContext: resourceExportLogsDelete,

		CustomizeDiff: config.FlexibleForceNew(exportLogsNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"export_file_name": {
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

func buildExportLogsFiltersBodyParams(d *schema.ResourceData) []map[string]interface{} {
	rawFilters, ok := d.Get("filters").([]interface{})
	if !ok || len(rawFilters) == 0 {
		return nil
	}

	filters := make([]map[string]interface{}, 0, len(rawFilters))
	for _, raw := range rawFilters {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		var values []string
		if v, ok := m["values"].([]interface{}); ok {
			values = utils.ExpandToStringList(v)
		}

		filters = append(filters, map[string]interface{}{
			"field":    m["field"],
			"operator": m["operator"],
			"values":   utils.ValueIgnoreEmpty(values),
		})
	}

	return filters
}

func buildExportLogsBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"start_time": d.Get("start_time"),
		"end_time":   d.Get("end_time"),
		"log_type":   d.Get("log_type"),
		"type":       d.Get("type"),
		"time_zone":  utils.ValueIgnoreEmpty(d.Get("time_zone")),
		"filters":    buildExportLogsFiltersBodyParams(d),
	}

	return body
}

func resourceExportLogsCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/{project_id}/cfw/{fw_instance_id}/logs/export"
		fwInstanceID = d.Get("fw_instance_id").(string)
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{fw_instance_id}", fwInstanceID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildExportLogsBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error exporting CFW logs: %s", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.Errorf("error reading response body: %s", err)
	}

	var outputFile string
	if v, ok := d.GetOk("export_file_name"); ok {
		outputFile = v.(string)
		if !strings.HasSuffix(outputFile, ".csv") {
			outputFile += ".csv"
		}
	} else {
		outputFile = fmt.Sprintf("cfw-%s-%s-log-%d.csv",
			d.Get("log_type").(string), d.Get("type").(string), d.Get("start_time").(int))
	}

	if err := os.WriteFile(outputFile, bodyBytes, 0600); err != nil {
		return diag.Errorf("failed to write firewall logs to (%s): %s", outputFile, err)
	}

	d.SetId(fwInstanceID)

	return nil
}

func resourceExportLogsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceExportLogsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceExportLogsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to export firewall logs. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
