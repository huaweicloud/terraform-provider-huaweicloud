package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

// @API IAM GET /v5/account-summary
func DataSourceIamIdentityV5AccountSummary() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIamIdentityV5AccountSummaryRead,

		Schema: map[string]*schema.Schema{
			"agencies_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"attached_policies_per_agency_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"attached_policies_per_group_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"attached_policies_per_user_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"groups_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policies": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policies_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy_size_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"root_user_mfa_enabled": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"users_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"versions_per_policy_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"agencies": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

type IdentityV5AccountSummaryDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newIdentityV5AccountSummaryDSWrapper(d *schema.ResourceData, meta interface{}) *IdentityV5AccountSummaryDSWrapper {
	return &IdentityV5AccountSummaryDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceIamIdentityV5AccountSummaryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newIdentityV5AccountSummaryDSWrapper(d, meta)
	getAccSumV5Rst, err := wrapper.GetAccountSummaryV5()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.getAccountSummaryV5ToSchema(getAccSumV5Rst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func (w *IdentityV5AccountSummaryDSWrapper) GetAccountSummaryV5() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "iam")
	if err != nil {
		return nil, err
	}

	uri := "/v5/account-summary"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *IdentityV5AccountSummaryDSWrapper) getAccountSummaryV5ToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("agencies_quota", body.Get("agencies_quota").Value()),
		d.Set("attached_policies_per_agency_quota", body.Get("attached_policies_per_agency_quota").Value()),
		d.Set("attached_policies_per_group_quota", body.Get("attached_policies_per_group_quota").Value()),
		d.Set("attached_policies_per_user_quota", body.Get("attached_policies_per_user_quota").Value()),
		d.Set("groups_quota", body.Get("groups_quota").Value()),
		d.Set("policies", body.Get("policies").Value()),
		d.Set("policies_quota", body.Get("policies_quota").Value()),
		d.Set("groups", body.Get("groups").Value()),
		d.Set("policy_size_quota", body.Get("policy_size_quota").Value()),
		d.Set("root_user_mfa_enabled", body.Get("root_user_mfa_enabled").Value()),
		d.Set("users", body.Get("users").Value()),
		d.Set("users_quota", body.Get("users_quota").Value()),
		d.Set("versions_per_policy_quota", body.Get("versions_per_policy_quota").Value()),
		d.Set("agencies", body.Get("agencies").Value()),
	)
	return mErr.ErrorOrNil()
}
