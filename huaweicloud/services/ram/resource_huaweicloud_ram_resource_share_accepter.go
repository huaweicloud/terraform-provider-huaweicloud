package ram

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RAM POST /v1/resource-share-invitations/{resource_share_invitation_id}/accept
// @API RAM POST /v1/resource-share-invitations/{resource_share_invitation_id}/reject
// ResourceShareAccepter is a definition of the one-time action resource that used to manage resource share.
func ResourceShareAccepter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceShareAccepterCreate,
		ReadContext:   resourceShareAccepterRead,
		DeleteContext: resourceShareAccepterDelete,

		Schema: map[string]*schema.Schema{
			"resource_share_invitation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceShareAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                          = meta.(*config.Config)
		region                       = cfg.GetRegion(d)
		resourceShareAccepterProduct = "ram"
	)

	resourceShareAccepterClient, err := cfg.NewServiceClient(resourceShareAccepterProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	action := d.Get("action").(string)
	resourceShareInvitationId := d.Get("resource_share_invitation_id").(string)
	httpUrl := "v1/resource-share-invitations/{resource_share_invitation_id}/{action}"
	httpUrl = strings.ReplaceAll(httpUrl, "{resource_share_invitation_id}", resourceShareInvitationId)
	httpUrl = strings.ReplaceAll(httpUrl, "{action}", action)
	createResourceShareAccepterPath := resourceShareAccepterClient.Endpoint + httpUrl
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResourceShareAccepterResp, err := resourceShareAccepterClient.Request("POST", createResourceShareAccepterPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error creating RAM resource share accepter: %s", err)
	}

	createResourceShareAccepterRespBody, err := utils.FlattenResponse(createResourceShareAccepterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	invitationId := utils.PathSearch("resource_share_invitation.resource_share_invitation_id", createResourceShareAccepterRespBody, "").(string)
	if invitationId == "" {
		return diag.Errorf("unable to find the resource share invitation ID of the RAM share accepter from the API response")
	}

	d.SetId(invitationId)
	return resourceShareAccepterRead(ctx, d, meta)
}

func resourceShareAccepterRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceShareAccepterDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a one-time action resource.
	return nil
}
