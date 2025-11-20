package identitycenter

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/users/batch-query
func DataSourceIdentityCenterBatchQueryUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceIdentityCenterBatchQueryUsersRead,

		Description: "schema: Internal",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identity_store_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the identity store that associated with IAM Identity Center.`,
			},
			"user_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"users": {
				Type:        schema.TypeList,
				Elem:        identityCenterUsersUserSchema(),
				Computed:    true,
				Description: `Indicates the list of IdentityCenter user.`,
			},
		},
	}
}

func resourceIdentityCenterBatchQueryUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getIdentityCenterUsersHttpUrl = "v1/identity-stores/{identity_store_id}/users/batch-query"
		getIdentityCenterUsersProduct = "identitystore"
	)
	getIdentityCenterUsersClient, err := cfg.NewServiceClient(getIdentityCenterUsersProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	getIdentityCenterUsersPath := getIdentityCenterUsersClient.Endpoint + getIdentityCenterUsersHttpUrl
	getIdentityCenterUsersPath = strings.ReplaceAll(getIdentityCenterUsersPath, "{identity_store_id}",
		d.Get("identity_store_id").(string))

	getIdentityCenterUsersOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getIdentityCenterUsersOpt.JSONBody = utils.RemoveNil(buildDescribeUsersBodyParam(d))

	getIdentityCenterUsersResp, err := getIdentityCenterUsersClient.Request("POST",
		getIdentityCenterUsersPath, &getIdentityCenterUsersOpt)

	if err != nil {
		return diag.FromErr(err)
	}

	getIdentityCenterUsersRespBody, err := utils.FlattenResponse(getIdentityCenterUsersResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("users", flattenGetIdentityCenterUsersResponseBodyUser(d, getIdentityCenterUsersRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDescribeUsersBodyParam(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_ids": d.Get("user_ids").([]interface{}),
	}
	return bodyParams
}
