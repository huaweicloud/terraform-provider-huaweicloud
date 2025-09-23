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

// @API HSS GET /v5/{project_id}/asset/common-port-info
func DataSourceAssetPortInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetPortInfoRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description_en": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAssetPortInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v5/{project_id}/asset/common-port-info"
		epsId    = cfg.GetEnterpriseProjectID(d)
		port     = d.Get("port").(int)
		category = d.Get("category").(string)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?port=%d&category=%s", getPath, port, category)
	if epsId != "" {
		getPath = fmt.Sprintf("%s&enterprise_project_id=%s", getPath, epsId)
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the port details information: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("type", utils.PathSearch("type", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("description_en", utils.PathSearch("description_en", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
