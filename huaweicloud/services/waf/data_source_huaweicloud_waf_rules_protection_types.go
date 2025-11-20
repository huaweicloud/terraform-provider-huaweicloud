package waf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF GET /v1/{project_id}/waf/rules/protection-types
func DataSourceRulesProtectionTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRulesProtectionTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildRulesProtectionTypesQueryParams(offset int) string {
	if offset > 0 {
		return fmt.Sprintf("?offset=%d", offset)
	}

	return ""
}

func dataSourceRulesProtectionTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/rules/protection-types"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	for {
		currentPath := requestPath + buildRulesProtectionTypesQueryParams(offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving rules protection types: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(items) == 0 {
			break
		}

		result = append(result, items...)
		offset += len(items)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("items", utils.ExpandToStringList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
