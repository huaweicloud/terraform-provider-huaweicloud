package dds

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

// @API DDS GET /v3/{project_id}/instances/{instance_id}/migrate/az
func DataSourceDDSMigrateAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDDSMigrateAvailabilityZonesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceDDSMigrateAvailabilityZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	listAZsHttpUrl := "v3/{project_id}/instances/{instance_id}/migrate/az"
	listAZsPath := client.Endpoint + listAZsHttpUrl
	listAZsPath = strings.ReplaceAll(listAZsPath, "{project_id}", client.ProjectID)
	listAZsPath = strings.ReplaceAll(listAZsPath, "{instance_id}", d.Get("instance_id").(string))
	listAZsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	listAZsResp, err := client.Request("GET", listAZsPath, &listAZsOpt)
	if err != nil {
		return diag.Errorf("error retrieving availability zones list: %s", err)
	}
	listAZsRespBody, err := utils.FlattenResponse(listAZsResp)
	if err != nil {
		return diag.Errorf("error flatten response: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("names", utils.PathSearch("az_list[?status == 'ENABLED'].code | sort(@)", listAZsRespBody, make([]interface{}, 0))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
