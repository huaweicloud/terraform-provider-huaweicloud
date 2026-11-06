package tms

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TMS GET /v2/tms/associated-resources/status
func DataSourceTmsAssociatedResourceStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmsAssociatedResourceStatusRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the associated resource tag inheritance.",
			},
		},
	}
}

func dataSourceTmsAssociatedResourceStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	d.SetId("tms-associated-resource-status")

	return diag.FromErr(mErr.ErrorOrNil())
}
