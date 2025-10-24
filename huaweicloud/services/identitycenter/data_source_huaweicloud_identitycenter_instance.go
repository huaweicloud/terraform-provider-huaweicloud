// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IdentityCenter
// ---------------------------------------------------------------

package identitycenter

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter GET /v1/instances
func DataSourceIdentityCenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceIdentityCenterRead,

		Schema: map[string]*schema.Schema{
			"identity_store_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the identity store that associated with IAM Identity Center.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the uniform resource name of the instance.`,
			},
			"alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the alias of the instance.`,
			},
		},
	}
}

func resourceIdentityCenterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getIdentityCenterInstance: Query Identity Center instance
	var (
		getIdentityCenterInstanceHttpUrl = "v1/instances"
		getIdentityCenterInstanceProduct = "identitycenter"
	)
	getIdentityCenterInstanceClient, err := cfg.NewServiceClient(getIdentityCenterInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	getIdentityCenterInstancePath := getIdentityCenterInstanceClient.Endpoint + getIdentityCenterInstanceHttpUrl

	getIdentityCenterInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getIdentityCenterInstanceResp, err := getIdentityCenterInstanceClient.Request("GET",
		getIdentityCenterInstancePath, &getIdentityCenterInstanceOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center Instance")
	}

	getIdentityCenterInstanceRespBody, err := utils.FlattenResponse(getIdentityCenterInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.PathSearch("instances|[0].instance_id",
		getIdentityCenterInstanceRespBody, "").(string))

	mErr = multierror.Append(
		mErr,
		d.Set("identity_store_id", utils.PathSearch("instances|[0].identity_store_id",
			getIdentityCenterInstanceRespBody, nil)),
		d.Set("urn", utils.PathSearch("instances|[0].instance_urn",
			getIdentityCenterInstanceRespBody, nil)),
		d.Set("alias", utils.PathSearch("instances|[0].alias",
			getIdentityCenterInstanceRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
