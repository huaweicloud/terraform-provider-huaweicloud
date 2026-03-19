package cfw

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

// @API CFW GET /v1/{project_id}/eip-count/{object_id}
func DataSourceEipCount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEipCountRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"eip_protected": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"eip_protected_self": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"eip_total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildQueryEipCountQueryParams(cfg *config.Config, d *schema.ResourceData) string {
	queryParam := ""

	if v := cfg.GetEnterpriseProjectID(d); v != "" {
		queryParam += fmt.Sprintf("&enterprise_project_id=%s", v)
	}

	if v, ok := d.GetOk("fw_instance_id"); ok {
		queryParam += fmt.Sprintf("&fw_instance_id=%s", v.(string))
	}

	if len(queryParam) > 0 {
		queryParam = "?" + queryParam[1:]
	}

	return queryParam
}

func dataSourceEipCountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/eip-count/{object_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{object_id}", d.Get("object_id").(string))
	requestPath += buildQueryEipCountQueryParams(cfg, d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW EIP count: %s", err)
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

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("eip_protected", utils.PathSearch("data.eip_protected", respBody, nil)),
		d.Set("eip_protected_self", utils.PathSearch("data.eip_protected_self", respBody, nil)),
		d.Set("eip_total", utils.PathSearch("data.eip_total", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
