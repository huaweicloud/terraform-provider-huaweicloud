package drs

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS GET /v5/{project_id}/agency/permissions
func DataSourceDrsAgencyPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsAgencyPermissionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_non_dbs": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"common_permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"engine_permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildAgencyPermissionsQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?is_non_dbs=%t", d.Get("is_non_dbs").(bool))

	if v, ok := d.GetOk("source_type"); ok {
		queryParams += fmt.Sprintf("&source_type=%s", v.(string))
	}
	if v, ok := d.GetOk("target_type"); ok {
		queryParams += fmt.Sprintf("&target_type=%s", v.(string))
	}

	return queryParams
}

func dataSourceDrsAgencyPermissionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/agency/permissions"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAgencyPermissionsQueryParams(d)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS agency permissions: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("common_permissions", utils.PathSearch("common_permissions", respBody, nil)),
		d.Set("engine_permissions", utils.PathSearch("engine_permissions", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
