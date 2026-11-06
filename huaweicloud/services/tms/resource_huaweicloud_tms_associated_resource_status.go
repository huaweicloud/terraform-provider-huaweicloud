package tms

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TMS GET /v2/tms/associated-resources/status
// @API TMS PUT /v2/tms/associated-resources/status
const tmsAssociatedResourceStatusID = "tms-associated-resource-status"

func ResourceTmsAssociatedResourceStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTmsAssociatedResourceStatusCreate,
		ReadContext:   resourceTmsAssociatedResourceStatusRead,
		UpdateContext: resourceTmsAssociatedResourceStatusUpdate,
		DeleteContext: resourceTmsAssociatedResourceStatusDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceTmsAssociatedResourceStatusImportState,
		},

		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, false),
				Description:  "Specifies whether the associated resource tag inheritance is enabled.",
			},
		},
	}
}

func resourceTmsAssociatedResourceStatusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("tms", region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	if err := changeAssociatedResourceStatus(client, d); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tmsAssociatedResourceStatusID)

	return resourceTmsAssociatedResourceStatusRead(ctx, d, meta)
}

func resourceTmsAssociatedResourceStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("tms", region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	getPath := client.Endpoint + "v2/tms/associated-resources/status"
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving TMS associated resource status: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	status := utils.PathSearch("status", getRespBody, "").(string)
	if status == "" {
		return diag.Errorf("error retrieving TMS associated resource status: status is not found in API response")
	}

	mErr := multierror.Append(
		d.Set("status", normalizeAssociatedResourceStatus(status)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func normalizeAssociatedResourceStatus(status string) string {
	switch status {
	case "enabled":
		return "enable"
	case "disabled":
		return "disable"
	default:
		return status
	}
}

func resourceTmsAssociatedResourceStatusUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("tms", region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	if err := changeAssociatedResourceStatus(client, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceTmsAssociatedResourceStatusRead(ctx, d, meta)
}

func changeAssociatedResourceStatus(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	changePath := client.Endpoint + "v2/tms/associated-resources/status"
	changeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	changeOpt.JSONBody = utils.RemoveNil(buildChangeAssociatedResourceStatusBodyParams(d))

	changeResp, err := client.Request("PUT", changePath, &changeOpt)
	if err != nil {
		return fmt.Errorf("error changing TMS associated resource status: %s", err)
	}

	changeRespBody, err := utils.FlattenResponse(changeResp)
	if err != nil {
		return err
	}

	errors := utils.PathSearch("errors", changeRespBody, make([]interface{}, 0)).([]interface{})
	if len(errors) > 0 {
		errMsgs := make([]string, 0, len(errors))
		for _, e := range errors {
			code := utils.PathSearch("error_code", e, "").(string)
			if code == "" {
				code = utils.PathSearch("code", e, "").(string)
			}
			message := utils.PathSearch("error_msg", e, "").(string)
			if message == "" {
				message = utils.PathSearch("message", e, "").(string)
			}
			errMsgs = append(errMsgs, fmt.Sprintf("code: %s, message: %s", code, message))
		}
		return fmt.Errorf("failed to change the TMS associated resource status, errors: %s",
			strings.Join(errMsgs, "; "))
	}

	return nil
}

func buildChangeAssociatedResourceStatusBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"status": d.Get("status"),
	}
}

func resourceTmsAssociatedResourceStatusImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	d.SetId(tmsAssociatedResourceStatusID)
	return []*schema.ResourceData{d}, nil
}

func resourceTmsAssociatedResourceStatusDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting the TMS associated resource status is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
