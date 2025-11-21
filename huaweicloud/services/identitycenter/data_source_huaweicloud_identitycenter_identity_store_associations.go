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

// @API IdentityCenter GET /v1/instances/{instance_id}/identity-store-associations
func DataSourceIdentityCenterIdentityStoreAssociations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCenterIdentityStoreAssociationsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_store_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authentication_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_type": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityCenterIdentityStoreAssociationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/instances/{instance_id}/identity-store-associations"
	)

	client, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listResp, err := client.Request("GET", listPath, &opt)

	if err != nil {
		return diag.FromErr(err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("identity_store_id", utils.PathSearch("identity_store_associations|[0].identity_store_id", listRespBody, nil)),
		d.Set("identity_store_type", utils.PathSearch("identity_store_associations|[0].identity_store_type", listRespBody, nil)),
		d.Set("authentication_type", utils.PathSearch("identity_store_associations|[0].authentication_type", listRespBody, nil)),
		d.Set("provisioning_type", utils.PathSearch("identity_store_associations|[0].provisioning_type", listRespBody, nil)),
		d.Set("status", utils.PathSearch("identity_store_associations|[0].status", listRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
