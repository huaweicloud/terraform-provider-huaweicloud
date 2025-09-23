package organizations

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations POST /v1/received-handshakes/{handshake_id}/decline
func ResourceAccountInviteDecliner() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountInviteDeclinerCreate,
		ReadContext:   resourceAccountInviteDeclinerRead,
		DeleteContext: resourceAccountInviteDeclinerDelete,

		Schema: map[string]*schema.Schema{
			"invitation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the unique ID of an invitation (handshake).`,
			},
		},
	}
}

func resourceAccountInviteDeclinerCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/received-handshakes/{handshake_id}/decline"
		product = "organizations"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{handshake_id}", d.Get("invitation_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Organizations account invite decliner: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	handshakeId := utils.PathSearch("handshake.id", createRespBody, "").(string)
	if handshakeId == "" {
		return diag.Errorf("unable to find the handshake ID of the Organizations account invite accepter from the API response")
	}
	d.SetId(handshakeId)

	return nil
}

func resourceAccountInviteDeclinerRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAccountInviteDeclinerDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting Organizations account invite decliner resource is not supported. The Organizations account " +
		"invite decliner resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
