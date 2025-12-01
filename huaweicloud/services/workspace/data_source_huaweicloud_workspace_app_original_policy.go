package workspace

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

// @API Workspace GET /v1/{project_id}/policy-groups/actions/list-original-policy
func DataSourceAppOriginalPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppOriginalPolicyRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the original policy is located.`,
			},
			"policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The original policy configuration, in JSON format.`,
			},
		},
	}
}

func dataSourceAppOriginalPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/policy-groups/actions/list-original-policy"
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return diag.FromErr(err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.Errorf("error querying original policy information: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("policy", utils.JsonToString(utils.PathSearch("policies", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
