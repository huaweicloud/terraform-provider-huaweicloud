package drs

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

// @API DRS GET /v3/{project_id}/node-type
func DataSourceNodeTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodeTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
			},
			"multi_write": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"node_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNodeTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listNodeTypesHttpUrl := "v3/{project_id}/node-type"
	listNodeTypesPath := client.Endpoint + listNodeTypesHttpUrl
	listNodeTypesPath = strings.ReplaceAll(listNodeTypesPath, "{project_id}", client.ProjectID)
	listNodeTypesPath += buildListNodeTypesQueryParams(d)
	listNodeTypesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	listNodeTypesResp, err := client.Request("GET", listNodeTypesPath, &listNodeTypesOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	listNodeTypesRespBody, err := utils.FlattenResponse(listNodeTypesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("node_types", utils.PathSearch("node_types[?is_sellout == `false`].node_type", listNodeTypesRespBody, make([]interface{}, 0))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListNodeTypesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("engine_type"); ok {
		res = fmt.Sprintf("%s&engine_type=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&db_use_type=%v", res, v)
	}
	if v, ok := d.GetOk("direction"); ok {
		res = fmt.Sprintf("%s&job_direction=%v", res, v)
	}
	if v, ok := d.GetOk("amulti_write"); ok {
		res = fmt.Sprintf("%s&is_multi_write=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
