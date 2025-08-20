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

// @API WAF GET /v1/{project_id}/waf/bundle
func DataSourceUserBundle() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserBundleRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString, // JSON format
				Computed: true,
			},
			"premium_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"premium_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"premium_host": {
				Type:     schema.TypeString, // JSON format
				Computed: true,
			},
			"options": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
			"rule": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"upgrade": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"feature": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
		},
	}
}

func dataSourceUserBundleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/bundle"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		getPath = fmt.Sprintf("%s?enterprise_project_id=%s", getPath, epsId)
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving user bundle information: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("type", utils.PathSearch("type", respBody, nil)),
		d.Set("host", utils.JsonToString(utils.PathSearch("host", respBody, nil))),
		d.Set("premium_name", utils.PathSearch("premium_name", respBody, nil)),
		d.Set("premium_type", utils.PathSearch("premium_type", respBody, nil)),
		d.Set("premium_host", utils.JsonToString(utils.PathSearch("premium_host", respBody, nil))),
		d.Set("options", utils.PathSearch("options", respBody, nil)),
		d.Set("rule", utils.PathSearch("rule", respBody, nil)),
		d.Set("upgrade", utils.PathSearch("upgrade", respBody, nil)),
		d.Set("feature", utils.PathSearch("feature", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
