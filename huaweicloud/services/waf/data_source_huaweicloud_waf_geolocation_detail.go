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

// @API WAF GET /v1/{project_id}/waf/tag/geoip/map
func DataSourceGeolocationDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeolocationDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"continent": {
				Type:     schema.TypeString, // JSON format
				Computed: true,
			},
			"geomap": {
				Type:     schema.TypeString, // JSON format
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGeolocationDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/waf/tag/geoip/map"
		geoLanguage = d.Get("lang").(string)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	if geoLanguage != "" {
		getPath = fmt.Sprintf("%s?lang=%s", getPath, geoLanguage)
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving geolocation detail: %s", err)
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
		d.Set("continent", utils.JsonToString(utils.PathSearch("continent", respBody, nil))),
		d.Set("geomap", utils.JsonToString(utils.PathSearch("geomap", respBody, nil))),
		d.Set("locale", utils.PathSearch("locale", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
