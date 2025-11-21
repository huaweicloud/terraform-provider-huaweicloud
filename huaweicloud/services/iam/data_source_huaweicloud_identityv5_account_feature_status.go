package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/features/{feature_name}
func DataSourceIdentityV5AccountFeatureStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5AccountFeatureStatusRead,

		Schema: map[string]*schema.Schema{
			"feature_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"feature_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityV5AccountFeatureStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	featureName := d.Get("feature_name").(string)
	path := client.Endpoint + "v5/features/{feature_name}"
	path = strings.ReplaceAll(path, "{feature_name}", featureName)
	reqOpt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	r, err := client.Request("GET", path, reqOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving account status")
	}
	resp, err := utils.FlattenResponse(r)
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	mErr := multierror.Append(nil,
		d.Set("feature_status", utils.PathSearch("feature_status", resp, "").(string)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
