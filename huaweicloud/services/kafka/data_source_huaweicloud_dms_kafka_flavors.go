package kafka

import (
	"context"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/dms/v2/products"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ChargingMode string

var (
	ChargingModePrePaid  ChargingMode = "prePaid"
	ChargingModePostPaid ChargingMode = "postPaid"

	ChargingModesMap map[string]ChargingMode = map[string]ChargingMode{
		"hourly":  ChargingModePrePaid,
		"monthly": ChargingModePostPaid,
	}
)

// @API Kafka GET /v2/{engine}/products
func DataSourceKafkaFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKafkaFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// We call the product as flavor.
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"arch_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ChargingModePrePaid), string(ChargingModePostPaid),
				}, false),
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"arch_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"charging_modes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ios": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage_spec_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"availability_zones": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"unavailability_zones": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"support_features": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"properties": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max_task": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"min_task": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"max_node": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"min_node": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_broker": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"min_broker": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_bandwidth_per_broker": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_consumer_per_broker": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_partition_per_broker": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_tps_per_broker": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_storage_per_node": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"min_storage_per_node": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"flavor_alias": {
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

type ios []products.IOEntity

func (s ios) Len() int {
	return len(s)
}

func (s ios) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ios) Less(i, j int) bool {
	return strings.ToLower(s[i].IoSpec) < strings.ToLower(s[j].IoSpec)
}

func filterFlavors(d *schema.ResourceData, flavorList []products.Product) ([]products.Product, []string) {
	if len(flavorList) < 1 {
		return nil, nil
	}

	result := make([]products.Product, 0, len(flavorList))
	ids := make([]string, 0, len(flavorList))

	t, tOk := d.GetOk("type")
	cm, cmOk := d.GetOk("charging_mode")
	at, atOk := d.GetOk("arch_type")
	sc, scOk := d.GetOk("storage_spec_code")
	azs := d.Get("availability_zones").([]interface{})

	for i, flavor := range flavorList {
		if tOk && flavor.Type != t.(string) {
			continue
		}
		if cmOk && !utils.StrSliceContains(flavor.ChargingModes, parseChargingMode(cm.(string))) {
			continue
		}
		if atOk && !utils.StrSliceContains(flavor.ArchTypes, at.(string)) {
			continue
		}
		validIOs := make([]products.IOEntity, 0, len(flavor.IOs))
		for _, io := range flavor.IOs {
			if scOk && io.IoSpec != sc.(string) {
				continue
			}
			if utils.StrSliceContainsAnother(io.AvailableZones, utils.ExpandToStringList(azs)) {
				validIOs = append(validIOs, io)
			}
		}
		if len(validIOs) < 1 {
			continue
		}
		sort.Sort(ios(validIOs))
		// The element "flavor" is just a copy.
		flavorList[i].IOs = validIOs
		result = append(result, flavorList[i])
		ids = append(ids, flavor.ProductId)
	}

	return result, ids
}

func flattenIOs(ios []products.IOEntity) []map[string]interface{} {
	if len(ios) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(ios))

	for i, io := range ios {
		result[i] = map[string]interface{}{
			"storage_spec_code":    io.IoSpec,
			"type":                 io.Type,
			"availability_zones":   io.AvailableZones,
			"unavailability_zones": io.UnavailableZones,
		}
	}

	log.Printf("[DEBUG] The result of IO list is: %#v", result)
	return result
}

func convertStringNumberIgnoreErr(strNum string) int {
	result, _ := strconv.Atoi(strNum)
	return result
}

func flattenSupportFeatures(features []products.SupportFeatureEntity) []map[string]interface{} {
	if len(features) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(features))

	for i, feature := range features {
		result[i] = map[string]interface{}{
			"name": feature.Name,
			"properties": []map[string]interface{}{
				{
					"max_task": convertStringNumberIgnoreErr(feature.Properties.MaxTask),
					"min_task": convertStringNumberIgnoreErr(feature.Properties.MinTask),
					"max_node": convertStringNumberIgnoreErr(feature.Properties.MaxNode),
					"min_node": convertStringNumberIgnoreErr(feature.Properties.MinNode),
				},
			},
		}
	}

	log.Printf("[DEBUG] The support features result is: %#v", result)
	return result
}

func flattenProperties(properties products.PropertiesEntity) []map[string]interface{} {
	if reflect.DeepEqual(properties, products.PropertiesEntity{}) {
		return nil
	}

	result := []map[string]interface{}{
		{
			"max_partition_per_broker": convertStringNumberIgnoreErr(properties.MaxPartitionPerBroker),
			"max_broker":               convertStringNumberIgnoreErr(properties.MaxBroker),
			"max_storage_per_node":     convertStringNumberIgnoreErr(properties.MaxStoragePerNode),
			"max_consumer_per_broker":  convertStringNumberIgnoreErr(properties.MaxConsumerPerBroker),
			"min_broker":               convertStringNumberIgnoreErr(properties.MinBroker),
			"max_bandwidth_per_broker": convertStringNumberIgnoreErr(properties.MaxBandwidthPerBroker),
			"min_storage_per_node":     convertStringNumberIgnoreErr(properties.MinStoragePerNode),
			"max_tps_per_broker":       convertStringNumberIgnoreErr(properties.MaxTpsPerBroker),
			"flavor_alias":             properties.ProductAlias,
		},
	}

	log.Printf("[DEBUG] The properties result is: %#v", result)
	return result
}

func parseChargingMode(chargingMode string) string {
	if chargingMode == string(ChargingModePostPaid) {
		return "hourly"
	}
	return "monthly"
}

func parseChargingModesResp(chargingModes []string) []interface{} {
	result := make([]interface{}, len(chargingModes))
	for i, val := range chargingModes {
		if cm, ok := ChargingModesMap[val]; ok {
			result[i] = cm
		} else {
			result[i] = val
		}
	}
	return result
}

func flattenFlavors(flavorList []products.Product) []map[string]interface{} {
	if len(flavorList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(flavorList))

	for i, val := range flavorList {
		result[i] = map[string]interface{}{
			"id":               val.ProductId,
			"type":             val.Type,
			"vm_specification": val.EcsFlavorId,
			"arch_types":       val.ArchTypes,
			"charging_modes":   parseChargingModesResp(val.ChargingModes),
			"ios":              flattenIOs(val.IOs),
			"support_features": flattenSupportFeatures(val.SupportFeatures),
			"properties":       flattenProperties(val.Properties),
		}
	}

	log.Printf("[DEBUG] The result of DMS flavor list is: %#v", result)
	return result
}

func dataSourceKafkaFlavorsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return DataSourceFlavorsRead(ctx, d, meta, "kafka")
}

func DataSourceFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}, engine string) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DmsV2Client(region)
	if err != nil {
		return diag.Errorf("error getting DMS v2 client: %v", err)
	}

	opt := products.ListOpts{
		ProductId: d.Get("flavor_id").(string),
	}
	resp, err := products.List(client, engine, opt)
	if err != nil {
		return diag.Errorf("error getting %s flavor list: %v", engine, err)
	}
	log.Printf("[DEBUG] The response of DMS %s flavor list request is: %#v", engine, resp)

	filtered, ids := filterFlavors(d, resp.Products)
	log.Printf("[DEBUG] The filtered list of DMS %s flavors is: %#v", engine, filtered)
	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("versions", resp.Versions),
		d.Set("flavors", flattenFlavors(filtered)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
