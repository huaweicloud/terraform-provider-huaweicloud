package coc

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

var changeDeleteNonUpdatableParams = []string{"ticket_type", "ticket_id"}

// @API COC DELETE /v1/{ticket_type}/tickets/{ticket_id}
func ResourceChangeDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChangeDeleteCreate,
		ReadContext:   resourceChangeDeleteRead,
		UpdateContext: resourceChangeDeleteUpdate,
		DeleteContext: resourceChangeDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(changeDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"ticket_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ticket_id": {
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

func resourceChangeDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	ticketID := d.Get("ticket_id").(string)
	httpUrl := "v1/{ticket_type}/tickets/{ticket_id}"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{ticket_type}", d.Get("ticket_type").(string))
	createPath = strings.ReplaceAll(createPath, "{ticket_id}", ticketID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error deleting COC change ticket: %s", err)
	}

	d.SetId(ticketID)

	return nil
}

func resourceChangeDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceChangeDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceChangeDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting change delete resource is not supported. The change delete resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
