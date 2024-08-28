package ccm

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CCM POST /v1/private-certificate-authorities/{ca_id}/restore
func ResourcePrivateCaRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateCaRestoreCreate,
		ReadContext:   resourcePrivateCaRestoreRead,
		DeleteContext: resourcePrivateCaRestoreDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"ca_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `ID of the CA you want to restore.`,
			},
		},
	}
}

func resourcePrivateCaRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/private-certificate-authorities/{ca_id}/restore"
		caID    = d.Get("ca_id").(string)
	)
	client, err := cfg.NewServiceClient("ccm", region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	restorePath := client.Endpoint + httpUrl
	restorePath = strings.ReplaceAll(restorePath, "{ca_id}", caID)
	restoreOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("POST", restorePath, &restoreOpt)
	if err != nil {
		return diag.Errorf("error restoring CCM private CA: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourcePrivateCaRestoreRead(ctx, d, meta)
}

func resourcePrivateCaRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePrivateCaRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
