package ecs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEcsSupplyRecommendations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_supply_recommendations.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsSupplyRecommendations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "supply_recommendations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "supply_recommendations.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "supply_recommendations.0.score"),
					resource.TestCheckOutput("architecture_type_filter_is_useful", "true"),
					resource.TestCheckOutput("vcpu_count_filter_is_useful", "true"),
					resource.TestCheckOutput("memory_mb_filter_is_useful", "true"),
					resource.TestCheckOutput("cpu_manufacturers_filter_is_useful", "true"),
					resource.TestCheckOutput("memory_gb_per_vcpu_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_generations_filter_is_useful", "true"),
					resource.TestCheckOutput("flavor_ids_filter_is_useful", "true"),
					resource.TestCheckOutput("locations_filter_is_useful", "true"),
					resource.TestCheckOutput("option_filter_is_useful", "true"),
					resource.TestCheckOutput("strategy_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEcsSupplyRecommendations_basic() string {
	return `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {}

data "huaweicloud_compute_supply_recommendations" "test" {
  flavor_ids = [data.huaweicloud_compute_flavors.test.flavors[0].id]
}

data "huaweicloud_compute_supply_recommendations" "architecture_type_filter" {
  flavor_constraint {
    architecture_type = ["arm64"]
  }
}
output "architecture_type_filter_is_useful" {
  value = length(data.huaweicloud_compute_supply_recommendations.architecture_type_filter.supply_recommendations) > 0
}

data "huaweicloud_compute_supply_recommendations" "vcpu_count_filter" {
  flavor_constraint {
    flavor_requirements {
      vcpu_count {
	    max = 1000000
	    min = -1
	  }
    }
  }
}
output "vcpu_count_filter_is_useful" {
  value = length(data.huaweicloud_compute_supply_recommendations.vcpu_count_filter.supply_recommendations) > 0
}

data "huaweicloud_compute_supply_recommendations" "memory_mb_filter" {
  flavor_constraint {
    flavor_requirements {
      memory_mb {
	    max = 1000000
	    min = -1
	  }
    }
  }
}
output "memory_mb_filter_is_useful" {
  value = length(data.huaweicloud_compute_supply_recommendations.memory_mb_filter.supply_recommendations) > 0
}

data "huaweicloud_compute_supply_recommendations" "cpu_manufacturers_filter" {
  flavor_constraint {
    flavor_requirements {
      cpu_manufacturers = ["INTEL"]
    }
  }
}
output "cpu_manufacturers_filter_is_useful" {
  value = length(data.huaweicloud_compute_supply_recommendations.cpu_manufacturers_filter.supply_recommendations) > 0
}

data "huaweicloud_compute_supply_recommendations" "memory_gb_per_vcpu_filter" {
  flavor_constraint {
    flavor_requirements {
      memory_gb_per_vcpu {
	    max = 1000000
	    min = -1
	  }
    }
  }
}
output "memory_gb_per_vcpu_filter_is_useful" {
  value = length(data.huaweicloud_compute_supply_recommendations.memory_gb_per_vcpu_filter.supply_recommendations) > 0
}

data "huaweicloud_compute_supply_recommendations" "instance_generations_filter" {
  flavor_constraint {
    flavor_requirements {
      instance_generations = ["CURRENT"]
    }
  }
}
output "instance_generations_filter_is_useful" {
  value = length(data.huaweicloud_compute_supply_recommendations.instance_generations_filter.supply_recommendations) > 0
}

data "huaweicloud_compute_supply_recommendations" "flavor_ids_filter" {
  flavor_ids = [data.huaweicloud_compute_flavors.test.flavors[1].id]
}
output "flavor_ids_filter_is_useful" {
  value = length(data.huaweicloud_compute_supply_recommendations.flavor_ids_filter.supply_recommendations) > 0
}

data "huaweicloud_compute_supply_recommendations" "locations_filter" {
  flavor_ids = [data.huaweicloud_compute_flavors.test.flavors[0].id]

  locations {
    region_id            = data.huaweicloud_availability_zones.test.region
    availability_zone_id = data.huaweicloud_availability_zones.test.names[0]
  }
}
output "locations_filter_is_useful" {
  value = length(data.huaweicloud_compute_supply_recommendations.locations_filter.supply_recommendations) > 0
}

data "huaweicloud_compute_supply_recommendations" "option_filter" {
  flavor_ids = [data.huaweicloud_compute_flavors.test.flavors[0].id]

  option {
    result_granularity = "BY_REGION"
    enable_spot        = "true"
  }
}
output "option_filter_is_useful" {
  value = length(data.huaweicloud_compute_supply_recommendations.option_filter.supply_recommendations) > 0
}

data "huaweicloud_compute_supply_recommendations" "strategy_filter" {
  flavor_ids = [data.huaweicloud_compute_flavors.test.flavors[0].id]

  option {
    result_granularity = "BY_REGION"
    enable_spot        = "true"
  }

  strategy = "COST"
}
output "strategy_filter_is_useful" {
  value = length(data.huaweicloud_compute_supply_recommendations.strategy_filter.supply_recommendations) > 0
}
`
}
