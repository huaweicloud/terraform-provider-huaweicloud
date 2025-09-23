package sms

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var sourceServerCommandResultReportNonUpdatableParams = []string{"server_id", "command_name", "result", "result_detail"}

// @API SMS POST /v3/sources/{server_id}/command_result
func ResourceSourceServerCommandResultReport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSourceServerCommandResultReportCreate,
		ReadContext:   resourceSourceServerCommandResultReportRead,
		UpdateContext: resourceSourceServerCommandResultReportUpdate,
		DeleteContext: resourceSourceServerCommandResultReportDelete,

		CustomizeDiff: config.FlexibleForceNew(sourceServerCommandResultReportNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"command_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"result": {
				Type:     schema.TypeString,
				Required: true,
			},
			"result_detail": {
				Type:     schema.TypeString,
				Required: true,
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

func buildSourceServerCommandResultReportCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"command_name":  d.Get("command_name"),
		"result":        d.Get("result"),
		"result_detail": parseJson(d.Get("result_detail").(string)),
	}

	return bodyParams
}

func parseJson(v string) interface{} {
	if v == "" {
		return nil
	}

	var data interface{}
	err := json.Unmarshal([]byte(v), &data)
	if err != nil {
		log.Printf("[WARN] Unable to parse JSON: %s", err)
		return v
	}

	return data
}

func resourceSourceServerCommandResultReportCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	createHttpUrl := "v3/sources/{server_id}/command_result"
	serverID := d.Get("server_id").(string)
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{server_id}", serverID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSourceServerCommandResultReportCreateOpts(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error operating the SMS source server command result report (%s): %s", serverID, err)
	}

	d.SetId(serverID)

	return nil
}

func resourceSourceServerCommandResultReportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSourceServerCommandResultReportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSourceServerCommandResultReportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting source server command result report resource is not supported. The source server command result" +
		" report resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
