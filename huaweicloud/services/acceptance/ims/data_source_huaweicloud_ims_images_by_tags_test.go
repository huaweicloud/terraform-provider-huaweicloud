package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceIMSImagesByTags_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_ims_images_by_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceIMSImagesByTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.%"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_ims_images_by_tags.total_count_test",
						"total_count"),

					resource.TestCheckOutput("without_any_tag_filter_is_useful", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_any_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_any_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testIMSImagesByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_compute_instance" "test" {
  name       = "%[2]s"
  image_name = "Ubuntu 18.04 server 64bit"
  flavor_id  = data.huaweicloud_compute_flavors.test.ids[0]

  security_group_ids = [
    huaweicloud_networking_secgroup.test.id
  ]

  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_ims_ecs_system_image" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_compute_instance.test.id
  description = "created by Terraform AccTest"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceIMSImagesByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ims_images_by_tags" "test" {
  depends_on = [huaweicloud_ims_ecs_system_image.test]

  action = "filter"
}

data "huaweicloud_ims_images_by_tags" "without_any_tag_filter" {
  depends_on = [huaweicloud_ims_ecs_system_image.test]

  action          = "filter"
  without_any_tag = true
}

output "without_any_tag_filter_is_useful" {
  value = length(data.huaweicloud_ims_images_by_tags.without_any_tag_filter.resources) == 0
}

data "huaweicloud_ims_images_by_tags" "matches_filter" {
  depends_on = [huaweicloud_ims_ecs_system_image.test]

  action = "filter"

  matches {
    key   = "resource_name"
    value = "%[2]s"
  }
}

output "matches_filter_is_useful" {
  value = length(data.huaweicloud_ims_images_by_tags.matches_filter.resources) > 0
}

data "huaweicloud_ims_images_by_tags" "tags_filter" {
  depends_on = [huaweicloud_ims_ecs_system_image.test]

  action = "filter"

  tags {
    key    = "foo"
    values = ["bar"]
  }
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_ims_images_by_tags.tags_filter.resources) > 0
}

data "huaweicloud_ims_images_by_tags" "not_tags_filter" {
  depends_on = [huaweicloud_ims_ecs_system_image.test]

  action = "filter"

  not_tags {
    key    = "no_foo"
    values = ["no_bar"]
  }
}

output "not_tags_filter_is_useful" {
  value = length(data.huaweicloud_ims_images_by_tags.not_tags_filter.resources) > 0
}

data "huaweicloud_ims_images_by_tags" "tags_any_filter" {
  depends_on = [huaweicloud_ims_ecs_system_image.test]

  action = "filter"

  tags_any {
    key    = "foo"
    values = ["bar"]
  }
}

output "tags_any_filter_is_useful" {
  value = length(data.huaweicloud_ims_images_by_tags.tags_any_filter.resources) > 0
}

data "huaweicloud_ims_images_by_tags" "not_tags_any_filter" {
  depends_on = [huaweicloud_ims_ecs_system_image.test]

  action = "filter"

  not_tags {
    key    = "no_foo"
    values = ["no_bar"]
  }
}

output "not_tags_any_filter_is_useful" {
  value = length(data.huaweicloud_ims_images_by_tags.not_tags_any_filter.resources) > 0
}

data "huaweicloud_ims_images_by_tags" "total_count_test" {
  depends_on = [huaweicloud_ims_ecs_system_image.test]

  action = "count"
}
`, testIMSImagesByTags_basic(name), name)
}
