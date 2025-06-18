package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI GET /apis/cci/v2/imagesnapshots
// @API CCI GET /apis/cci/v2/imagesnapshots/{name}
func DataSourceV2ImageSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2ImageSnapshotsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"finalizers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"building_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_create_eip": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"auto_create_eip_attribute": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bandwidth_charge_mode": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"bandwidth_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"bandwidth_size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"ip_version": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"eip_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"image_snapshot_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"images": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"registries": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_pull_secret": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"insecure_skip_verify": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"plain_http": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"server": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ttl_days_after_created": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expire_date_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"images": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"digest": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"image": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"size_bytes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"last_updated_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"phase": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshot_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceV2ImageSnapshotsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	results := make([]interface{}, 0)
	if name, ok := d.GetOk("name"); ok {
		resp, err := GetImageSnapshot(client, name.(string))
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); !ok {
				return diag.Errorf("error getting the image snapshots from the server: %s", err)
			}
		}
		if resp != nil {
			results = append(results, resp)
		}
	} else {
		resp, err := listImageSnapshots(client)
		if err != nil {
			return diag.Errorf("error getting the image snapshots from the server: %s", err)
		}
		results = utils.PathSearch("items", resp, make([]interface{}, 0)).([]interface{})
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("image_snapshots", flattenImageSnapshots(results)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImageSnapshots(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, v := range results {
		rst = append(rst, map[string]interface{}{
			"name":                   utils.PathSearch("metadata.name", v, nil),
			"annotations":            utils.PathSearch("metadata.annotations", v, nil),
			"labels":                 utils.PathSearch("metadata.labels", v, nil),
			"finalizers":             utils.PathSearch("spec.finalizers", v, nil),
			"resource_version":       utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                    utils.PathSearch("metadata.uid", v, nil),
			"building_config":        flattenImageSnapshotSpecBuildingConfig(v),
			"image_snapshot_size":    utils.PathSearch("spec.imageSnapshotSize", v, nil),
			"images":                 flattenImageSnapshotSpecImages(v),
			"registries":             flattenImageSnapshotSpecRegistries(v),
			"ttl_days_after_created": utils.PathSearch("spec.ttlDaysAfterCreated", v, nil),
			"status":                 flattenImageSnapshotStatus(v),
		})
	}
	return rst
}

func listImageSnapshots(client *golangsdk.ServiceClient) (interface{}, error) {
	listImageSnapshotsHttpUrl := "apis/cci/v2/imagesnapshots"
	listImageSnapshotsPath := client.Endpoint + listImageSnapshotsHttpUrl
	listImageSnapshotsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listImageSnapshotssResp, err := client.Request("GET", listImageSnapshotsPath, &listImageSnapshotsOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(listImageSnapshotssResp)
}
