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

// @API CFW GET /v1/{project_id}/eip/auto-protect-status/{object_id}
func DataSourceEipAutoProtectStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEipAutoProtectStatusRead,

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
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"available_eip_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"beyond_max_count": {
							Type:     schema.TypeBool,
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
						"eip_un_protected": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"object_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildEipAutoProtectStatusQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceEipAutoProtectStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		objectID = d.Get("object_id").(string)
		epsId    = cfg.GetEnterpriseProjectID(d)
		httpUrl  = "v1/{project_id}/eip/auto-protect-status/{object_id}"
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{object_id}", objectID)
	requestPath += buildEipAutoProtectStatusQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW EIP auto protect status: %s", err)
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
		d.Set("data", flattenEipAutoProtectStatusData(utils.PathSearch("data", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEipAutoProtectStatusData(data interface{}) []interface{} {
	if data == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"available_eip_count": utils.PathSearch("available_eip_count", data, nil),
			"beyond_max_count":    utils.PathSearch("beyond_max_count", data, nil),
			"eip_protected_self":  utils.PathSearch("eip_protected_self", data, nil),
			"eip_total":           utils.PathSearch("eip_total", data, nil),
			"eip_un_protected":    utils.PathSearch("eip_un_protected", data, nil),
			"object_id":           utils.PathSearch("object_id", data, nil),
			"status":              utils.PathSearch("status", data, nil),
		},
	}
}
