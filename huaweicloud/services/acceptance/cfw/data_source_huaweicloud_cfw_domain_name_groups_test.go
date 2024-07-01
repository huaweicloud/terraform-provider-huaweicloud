package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwDomainNameGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_domain_name_groups.filter_by_id"
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
				Config: testDataSourceCfwDomainNameGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.domain_names.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.domain_names.0.domain_address_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.rules.0.name"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_config_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCfwDomainNameGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  id            = huaweicloud_cfw_domain_name_group.g2.id
  name          = huaweicloud_cfw_domain_name_group.g2.name
  type          = tostring(huaweicloud_cfw_domain_name_group.g1.type)
  config_status = "3"
}

data "huaweicloud_cfw_domain_name_groups" "test" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id

  depends_on = [
    huaweicloud_cfw_protection_rule.test,
  ]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_cfw_domain_name_groups.test.records) >= 2
}

data "huaweicloud_cfw_domain_name_groups" "filter_by_id" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  group_id       = local.id

  depends_on = [
    huaweicloud_cfw_protection_rule.test,
  ]
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_cfw_domain_name_groups.filter_by_id.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_domain_name_groups.filter_by_id.records[*] : v.group_id == local.id]
  )
}

data "huaweicloud_cfw_domain_name_groups" "filter_by_name" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name           = local.name

  depends_on = [
    huaweicloud_cfw_protection_rule.test,
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_cfw_domain_name_groups.filter_by_name.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_domain_name_groups.filter_by_name.records[*] : v.name == local.name]
  )
}

data "huaweicloud_cfw_domain_name_groups" "filter_by_type" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  type           = local.type

  depends_on = [
    huaweicloud_cfw_protection_rule.test,
  ]
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_cfw_domain_name_groups.filter_by_type.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_domain_name_groups.filter_by_type.records[*] : v.type == local.type]
  )
}

data "huaweicloud_cfw_domain_name_groups" "filter_by_config_status" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  config_status  = local.config_status

  depends_on = [
    huaweicloud_cfw_protection_rule.test,
  ]
}

output "is_config_status_filter_useful" {
  value = length(data.huaweicloud_cfw_domain_name_groups.filter_by_config_status.records) >= 1 && alltrue([
    for v in data.huaweicloud_cfw_domain_name_groups.filter_by_config_status.records[*] : 
      v.config_status == local.config_status
  ])
}
`, testDataSourceCfwDomainNameGroups_base(name), acceptance.HW_CFW_INSTANCE_ID)
}

func testDataSourceCfwDomainNameGroups_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_domain_name_group" "g1" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name           = "%[3]s1"
  type           = 0
  description    = "created by terraform"
	  
  domain_names {
   domain_name = "www.test1.com"
   description = "test domain 1"
  }

  domain_names {
    domain_name = "www.test2.com"
    description = "test domain 2"
  }
}

resource "huaweicloud_cfw_domain_name_group" "g2" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name           = "%[3]s2"
  type           = 1
  description    = "created by terraform"
		
  domain_names {
    domain_name = "www.test3.com"
    description = "test domain 3"
  }

  domain_names {
   domain_name = "www.test4.com"
   description = "test domain 4"
  }
}

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%[3]s"
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
    type            = 6
    domain_set_id   = huaweicloud_cfw_domain_name_group.g2.id
    domain_set_name = huaweicloud_cfw_domain_name_group.g2.name
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
    key = "value"
  }

  depends_on = [
    huaweicloud_cfw_domain_name_group.g1,
    huaweicloud_cfw_domain_name_group.g2,
  ] 
}
`, testAccDatasourceFirewalls_basic(), acceptance.HW_CFW_INSTANCE_ID, name)
}
