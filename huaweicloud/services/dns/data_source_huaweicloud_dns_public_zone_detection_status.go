package dns

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

// @API DNS GET /v2.1/zones/{zone_id}/detection
func DataSourcePublicZoneDetectionStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicZoneDetectionStatusRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the public zone.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the record set to be detected.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the record set to be detected.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the domain detection.`,
			},
		},
	}
}

func buildPublicZoneDetectionStatusQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?domain_name=%v", d.Get("domain_name"))

	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourcePublicZoneDetectionStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v2.1/zones/{zone_id}/detection"
		zoneId  = d.Get("zone_id").(string)
	)
	client, err := cfg.NewServiceClient("dns", "")
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{zone_id}", zoneId)
	getPath += buildPublicZoneDetectionStatusQueryParams(d)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
	requestResp, err := client.Request("GET", getPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error querying public zone (%s) detection status: %s", zoneId, err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
