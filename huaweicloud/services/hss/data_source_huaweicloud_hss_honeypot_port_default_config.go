package hss

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

// @API HSS GET /v5/{project_id}/honeypot-port/default-config
func DataSourceHoneypotPortDefaultConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHoneypotPortDefaultConfigRead,

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
			"auto_bind": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			// The field name in the API document is `windows_policy`,
			// but the actual returned name is `windows_policy_id`.
			"windows_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// The field name in the API document is `linux_policy`,
			// but the actual returned name is `linux_policy_id`.
			"linux_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildHoneypotPortDefaultConfigQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	return ""
}

func dataSourceHoneypotPortDefaultConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/honeypot-port/default-config"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildHoneypotPortDefaultConfigQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS honeypot port default config: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("auto_bind", utils.PathSearch("auto_bind", respBody, nil)),
		d.Set("windows_policy_id", utils.PathSearch("windows_policy_id", respBody, nil)),
		d.Set("linux_policy_id", utils.PathSearch("linux_policy_id", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
