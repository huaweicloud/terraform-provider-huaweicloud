package rabbitmq

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/kafka"
)

type ChargingMode string

var (
	ChargingModePrePaid  ChargingMode = "prePaid"
	ChargingModePostPaid ChargingMode = "postPaid"

	ChargingModesMap = map[string]ChargingMode{
		"hourly":  ChargingModePrePaid,
		"monthly": ChargingModePostPaid,
	}
)

// @API RabbitMQ GET /v2/{engine}/products
func DataSourceRabbitMQFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRabbitMQFlavorsRead,

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

func dataSourceRabbitMQFlavorsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return kafka.DataSourceFlavorsRead(ctx, d, meta, "rabbitmq")
}
