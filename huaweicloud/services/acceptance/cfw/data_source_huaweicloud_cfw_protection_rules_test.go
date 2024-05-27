package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwProtectionRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_protection_rules.filter_by_id"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwProtectionRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.rule_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.action_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.direction"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.source.0.address"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.destination.0.address"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.service.0.type"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_action_type_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_direction_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCfwProtectionRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  id          = huaweicloud_cfw_protection_rule.r1.id
  name        = huaweicloud_cfw_protection_rule.r1.name
  action_type = tostring(huaweicloud_cfw_protection_rule.r1.action_type)
  status      = tostring(huaweicloud_cfw_protection_rule.r2.status)
  direction   = tostring(huaweicloud_cfw_protection_rule.r2.direction)
  tags        = huaweicloud_cfw_protection_rule.r1.tags
}

data "huaweicloud_cfw_protection_rules" "filter_by_id" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  rule_id   = local.id

  depends_on = [
    huaweicloud_cfw_protection_rule.r1,
    huaweicloud_cfw_protection_rule.r2,
  ]
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_cfw_protection_rules.filter_by_id.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_protection_rules.filter_by_id.records[*] : v.rule_id == local.id]
  )
}

data "huaweicloud_cfw_protection_rules" "filter_by_name" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name      = local.name

  depends_on = [
    huaweicloud_cfw_protection_rule.r1,
    huaweicloud_cfw_protection_rule.r2,
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_cfw_protection_rules.filter_by_name.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_protection_rules.filter_by_name.records[*] : v.name == local.name]
  )
}

data "huaweicloud_cfw_protection_rules" "filter_by_action_type" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  action_type = local.action_type

  depends_on = [
    huaweicloud_cfw_protection_rule.r1,
    huaweicloud_cfw_protection_rule.r2,
  ]
}

output "is_action_type_filter_useful" {
  value = length(data.huaweicloud_cfw_protection_rules.filter_by_action_type.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_protection_rules.filter_by_action_type.records[*] : v.action_type == local.action_type]
  )
}

data "huaweicloud_cfw_protection_rules" "filter_by_status" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  status    = local.status

  depends_on = [
    huaweicloud_cfw_protection_rule.r1,
    huaweicloud_cfw_protection_rule.r2,
  ]
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_cfw_protection_rules.filter_by_status.records) >= 1 && alltrue([
    for v in data.huaweicloud_cfw_protection_rules.filter_by_status.records[*] : 
      v.status == local.status
  ])
}

data "huaweicloud_cfw_protection_rules" "filter_by_direction" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  direction = local.direction

  depends_on = [
    huaweicloud_cfw_protection_rule.r1,
    huaweicloud_cfw_protection_rule.r2,
  ]
}

output "is_direction_filter_useful" {
  value = length(data.huaweicloud_cfw_protection_rules.filter_by_direction.records) >= 1 && alltrue([
    for v in data.huaweicloud_cfw_protection_rules.filter_by_direction.records[*] : 
      v.direction == local.direction
  ])
}

data "huaweicloud_cfw_protection_rules" "filter_by_tags" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  tags      = local.tags

  depends_on = [
    huaweicloud_cfw_protection_rule.r1,
    huaweicloud_cfw_protection_rule.r2,
  ]
}

output "is_tags_filter_useful" {
  value = length(data.huaweicloud_cfw_protection_rules.filter_by_tags.records) >= 1 && alltrue([
    for pr in data.huaweicloud_cfw_protection_rules.filter_by_tags.records : alltrue([
      for k, v in local.tags : pr.tags[k] == v
    ])
  ])
}
`, testDataSourceCfwProtectionRules_base(name))
}

func testDataSourceCfwProtectionRules_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_protection_rule" "r1" {
  name                = "%[2]s1"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1
  direction           = 1
  
  source {
    type    = 0
    address = "1.1.1.1"
  }

  destination {
    type    = 0
    address = "2.2.2.1"
  }

  service {
    type = 2

    custom_service {
      protocol    = 6
      source_port = 80
      dest_port   = 80			
    }

    custom_service {
      protocol    = 6
      source_port = 8080
      dest_port   = 8080
    }
  }

  sequence {
    top = 1
  }

  tags = {
    k2 = "v2"
  }
}

resource "huaweicloud_cfw_protection_rule" "r2" {
  name                = "%[2]s2"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 1
  long_connect_enable = 0
  status              = 0
  direction           = 0
  
  source {
    type    = 0
    address = "1.1.1.2"
  }

  destination {
    type    = 0
    address = "2.2.2.2"
  }

  service {
    type = 2

    custom_service {
      protocol    = 6
      source_port = 80
      dest_port   = 80			
    }

    custom_service {
      protocol    = 6
      source_port = 8080
      dest_port   = 8080
    }
  }

  sequence {
    top = 1
  }

  tags = {
    k1 = "v1"
  }
}
`, testAccDatasourceFirewalls_basic(), name)
}
