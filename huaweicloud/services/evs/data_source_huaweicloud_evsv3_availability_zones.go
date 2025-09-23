package evs

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceEvsV3AvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvsV3AvailabilityZonesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of availability zones.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_available": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the availability zone is available.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of availability zone.`,
						},
					},
				},
			},
		},
	}
}

type AvailabilityZonesDSWrapperV3 struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newAvailabilityZonesDSWrapperV3(d *schema.ResourceData, meta interface{}) *AvailabilityZonesDSWrapperV3 {
	return &AvailabilityZonesDSWrapperV3{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceEvsV3AvailabilityZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newAvailabilityZonesDSWrapperV3(d, meta)
	cinLisAvaZonRst, err := wrapper.CinderListAvailabilityZones()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.cinderListAvailabilityZonesToSchema(cinLisAvaZonRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API EVS GET /v3/{project_id}/os-availability-zone
func (w *AvailabilityZonesDSWrapperV3) CinderListAvailabilityZones() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "evs")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/os-availability-zone"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OkCode(200).
		Request().
		Result()
}

func (w *AvailabilityZonesDSWrapperV3) cinderListAvailabilityZonesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("availability_zones", schemas.SliceToList(body.Get("availabilityZoneInfo"),
			func(avaZones gjson.Result) any {
				return map[string]any{
					"is_available": avaZones.Get("zoneState.available").Value(),
					"name":         avaZones.Get("zoneName").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
