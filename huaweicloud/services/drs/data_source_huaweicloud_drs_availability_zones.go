package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS POST /v3/{project_id}/available-zone
func DataSourceAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAvailabilityZonesRead,

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
			"node_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"multi_write": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAvailabilityZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v3 client, error: %s", err)
	}

	listAZsHttpUrl := "v3/{project_id}/available-zone"
	listAZsPath := client.Endpoint + listAZsHttpUrl
	listAZsPath = strings.ReplaceAll(listAZsPath, "{project_id}", client.ProjectID)
	listAZsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"engine_type":   d.Get("engine_type"),
			"db_use_type":   d.Get("type"),
			"job_direction": d.Get("direction"),
			"node_type":     d.Get("node_type"),
			"multi_write":   utils.ValueIgnoreEmpty(d.Get("multi_write")),
		}),
	}

	listAZsResp, err := client.Request("POST", listAZsPath, &listAZsOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	listAZsRespBody, err := utils.FlattenResponse(listAZsResp)
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
		d.Set("names", utils.PathSearch("az_infos[?status == 'ENABLED'].code | sort(@)", listAZsRespBody, make([]interface{}, 0))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
