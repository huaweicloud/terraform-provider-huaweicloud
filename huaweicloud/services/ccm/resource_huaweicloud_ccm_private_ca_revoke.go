package ccm

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCM POST /v1/private-certificate-authorities/{ca_id}/revoke
// The resource is a one-time action resource.
func ResourcePrivateCaRevoke() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateCaRevokeCreate,
		ReadContext:   resourcePrivateCaRevokeRead,
		DeleteContext: resourcePrivateCaRevokeDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"ca_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reason": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func buildPrivateCaRevokeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"reason": utils.ValueIgnoreEmpty(d.Get("reason")),
	}
	return bodyParams
}

func resourcePrivateCaRevokeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/private-certificate-authorities/{ca_id}/revoke"
		product = "ccm"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{ca_id}", d.Get("ca_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: buildPrivateCaRevokeBodyParams(d),
	}
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error revoking CCM private CA: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourcePrivateCaRevokeRead(ctx, d, meta)
}

func resourcePrivateCaRevokeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePrivateCaRevokeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
