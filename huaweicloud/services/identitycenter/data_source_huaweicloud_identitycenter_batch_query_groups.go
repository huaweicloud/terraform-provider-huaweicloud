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

// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/groups/batch-query
func DataSourceIdentityCenterBatchQueryGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceIdentityCenterBatchQueryGroupsRead,

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
			"group_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:        schema.TypeList,
				Elem:        identityCenterGroupsGroupSchema(),
				Computed:    true,
				Description: `Indicates the list of IdentityCenter group.`,
			},
		},
	}
}

func resourceIdentityCenterBatchQueryGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getIdentityCenterGroupsHttpUrl = "v1/identity-stores/{identity_store_id}/groups/batch-query"
		getIdentityCenterGroupsProduct = "identitystore"
	)
	getIdentityCenterGroupsClient, err := cfg.NewServiceClient(getIdentityCenterGroupsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	getIdentityCenterGroupsPath := getIdentityCenterGroupsClient.Endpoint + getIdentityCenterGroupsHttpUrl
	getIdentityCenterGroupsPath = strings.ReplaceAll(getIdentityCenterGroupsPath, "{identity_store_id}",
		d.Get("identity_store_id").(string))

	getIdentityCenterGroupsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getIdentityCenterGroupsOpt.JSONBody = utils.RemoveNil(buildDescribeGroupsBodyParam(d))

	getIdentityCenterGroupsResp, err := getIdentityCenterGroupsClient.Request("POST",
		getIdentityCenterGroupsPath, &getIdentityCenterGroupsOpt)

	if err != nil {
		return diag.FromErr(err)
	}

	getIdentityCenterGroupsRespBody, err := utils.FlattenResponse(getIdentityCenterGroupsResp)
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
		d.Set("groups", flattenGetIdentityCenterGroupsResponseBodyGroup(getIdentityCenterGroupsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDescribeGroupsBodyParam(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_ids": d.Get("group_ids").([]interface{}),
	}
	return bodyParams
}
