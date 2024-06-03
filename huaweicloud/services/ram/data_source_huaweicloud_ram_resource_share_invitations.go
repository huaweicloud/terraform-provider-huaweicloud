package ram

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const pageLimitCount = 200

// @API RAM POST /v1/resource-share-invitations/search
func DataSourceResourceShareInvitations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceShareInvitationssRead,
		Schema: map[string]*schema.Schema{
			"resource_share_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_share_invitation_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_share_invitations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceShareInvitationsSchema(),
			},
		},
	}
}

func resourceShareInvitationsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_share_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_share_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"receiver_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sender_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceResourceShareInvitationssRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                                = meta.(*config.Config)
		region                             = cfg.GetRegion(d)
		mErr                               *multierror.Error
		getResourceShareInvitationsHttpUrl = "v1/resource-share-invitations/search"
		getResourceShareInvitationsProduct = "ram"
	)
	getResourceShareInvitationsClient, err := cfg.NewServiceClient(getResourceShareInvitationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	getResourceShareInvitationsPath := getResourceShareInvitationsClient.Endpoint + getResourceShareInvitationsHttpUrl
	getResourceShareInvitationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var nextMarker string
	allResourceShareInvitations := make([]interface{}, 0)
	getResourceShareInvitationsJSONBody := utils.RemoveNil(buildGetResourceShareInvitationsBodyParams(d))
	for {
		getResourceShareInvitationsOpt.JSONBody = getResourceShareInvitationsJSONBody
		getResourceShareInvitationsResp, err := getResourceShareInvitationsClient.Request("POST",
			getResourceShareInvitationsPath, &getResourceShareInvitationsOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving RAM resource share invitations")
		}

		getResourceShareInvitationsRespBody, err := utils.FlattenResponse(getResourceShareInvitationsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		resourceShareInvitations := utils.PathSearch("resource_share_invitations", getResourceShareInvitationsRespBody,
			make([]interface{}, 0)).([]interface{})

		if len(resourceShareInvitations) > 0 {
			allResourceShareInvitations = append(allResourceShareInvitations, resourceShareInvitations...)
		}

		nextMarker = utils.PathSearch("page_info.next_marker", getResourceShareInvitationsRespBody, "").(string)
		if nextMarker == "" {
			break
		}
		getResourceShareInvitationsJSONBody["marker"] = nextMarker
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("resource_share_invitations", flattenGetResourceShareInvitationsResponseBody(allResourceShareInvitations)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetResourceShareInvitationsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"resource_share_ids":            utils.ValueIgnoreEmpty(d.Get("resource_share_ids")),
		"resource_share_invitation_ids": utils.ValueIgnoreEmpty(d.Get("resource_share_invitation_ids")),
		"status":                        utils.ValueIgnoreEmpty(d.Get("status")),
		"limit":                         pageLimitCount,
	}
	return bodyParams
}

func flattenGetResourceShareInvitationsResponseBody(curArray []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("resource_share_invitation_id", v, nil),
			"resource_share_id":   utils.PathSearch("resource_share_id", v, nil),
			"resource_share_name": utils.PathSearch("resource_share_name", v, nil),
			"receiver_account_id": utils.PathSearch("receiver_account_id", v, nil),
			"sender_account_id":   utils.PathSearch("sender_account_id", v, nil),
			"status":              utils.PathSearch("status", v, nil),
			"created_at":          utils.PathSearch("created_at", v, nil),
			"updated_at":          utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}
