package rocketmq

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// @API RocketMQ GET /v2/{engine}/products
func DataSourceRocketMQFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRocketMQFlavorsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"arch_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of CPU architecture.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of availability zone names.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the billing mode of the flavor.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the flavor.`,
			},
			"storage_spec_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the disk IO encoding.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the flavor.`,
			},
			"versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the list of flavor versions.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        rocketmqFlavorSchema(),
				Description: `Indicates the list of flavors.`,
			},
		},
	}
}

func rocketmqFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the flavor.`,
			},
			"arch_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the list of the types of CPU architecture.`,
			},
			"charging_modes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the list of the billing modes.`,
			},
			"ios": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        flavorIosSchema(),
				Description: `Indicates the list of disk IO types.`,
			},
			"support_features": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        supportFeatureSchema(),
				Description: `Indicates the list of features supported by the current specification.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the flavor.`,
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        propertySchema(),
				Description: `Indicates the list of the properties of the current specification.`,
			},
			"vm_specification": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the underlying VM specification.`,
			},
		},
	}
	return &sc
}

func flavorIosSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"storage_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the disk IO encoding.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the disk type.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the list of availability zone names.`,
			},
			"unavailability_zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the list of unavailability zone names.`,
			},
		},
	}
	return &sc
}

func supportFeatureSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the function name.`,
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        supportFeaturePropertySchema(),
				Description: `Indicates the list of the function property details.`,
			},
		},
	}
	return &sc
}

func supportFeaturePropertySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_task": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of tasks for the dump function.`,
			},
			"min_task": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the minimum number of tasks for the dump function.`,
			},
			"max_node": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of nodes for the dump function.`,
			},
			"min_node": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the minimum number of nodes for the dump function.`,
			},
		},
	}
	return &sc
}

func propertySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_broker": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of brokers.`,
			},
			"min_broker": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the minimum number of brokers.`,
			},
			"max_bandwidth_per_broker": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum bandwidth per broker.`,
			},
			"max_consumer_per_broker": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of consumers per broker.`,
			},
			"max_partition_per_broker": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of partitions per broker.`,
			},
			"max_tps_per_broker": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum TPS per broker.`,
			},
			"max_storage_per_node": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum storage per node. The unit is GB.`,
			},
			"min_storage_per_node": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the minimum storage per node. The unit is GB.`,
			},
			"flavor_alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the alias of the flavor.`,
			},
		},
	}
	return &sc
}

func dataSourceRocketMQFlavorsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return kafka.DataSourceFlavorsRead(ctx, d, meta, "rocketmq")
}
