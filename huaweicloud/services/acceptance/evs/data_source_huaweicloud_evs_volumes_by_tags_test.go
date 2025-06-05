package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVolumesByTags_basic(t *testing.T) {
	var (
		filter_by_tags_all    = "data.huaweicloud_evs_volumes_by_tags.filter_by_tags"
		dc_tags               = acceptance.InitDataSourceCheck(filter_by_tags_all)
		filter_by_matches_all = "data.huaweicloud_evs_volumes_by_tags.filter_by_matches"
		dc_matches            = acceptance.InitDataSourceCheck(filter_by_matches_all)
		testName              = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVolumesByTags_basic1(testName),
				Check: resource.ComposeTestCheckFunc(
					dc_tags.CheckResourceExists(),
					resource.TestCheckResourceAttr(filter_by_tags_all, "resources.0.resource_name", fmt.Sprintf("%s_0", testName)),
					resource.TestCheckResourceAttrSet(filter_by_tags_all, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(filter_by_tags_all, "resources.0.resource_detail.0.id"),
					resource.TestCheckResourceAttrSet(filter_by_tags_all, "resources.0.resource_detail.0.name"),
					resource.TestCheckResourceAttrSet(filter_by_tags_all, "resources.0.resource_detail.0.status"),
					resource.TestCheckResourceAttrSet(filter_by_tags_all, "resources.0.resource_detail.0.size"),
					resource.TestCheckResourceAttrSet(filter_by_tags_all, "resources.0.resource_detail.0.volume_type"),
					resource.TestCheckResourceAttrSet(filter_by_tags_all, "resources.0.resource_detail.0.availability_zone"),
					resource.TestCheckResourceAttrSet(filter_by_tags_all, "resources.0.resource_detail.0.created_at"),
					resource.TestCheckResourceAttrSet(filter_by_tags_all, "resources.0.resource_detail.0.updated_at"),
					resource.TestCheckResourceAttrSet(filter_by_tags_all, "resources.0.resource_detail.0.bootable"),
					resource.TestCheckResourceAttrSet(filter_by_tags_all, "resources.0.resource_detail.0.multiattach"),
					resource.TestCheckResourceAttr(filter_by_tags_all, "resources.0.tags.0.key", "foo1"),
					resource.TestCheckResourceAttr(filter_by_tags_all, "resources.0.tags.0.value", "bar1"),
					resource.TestCheckResourceAttr(filter_by_tags_all, "resources.0.tags.1.key", "foo0"),
					resource.TestCheckResourceAttr(filter_by_tags_all, "resources.0.tags.1.value", "bar0"),
				),
			},
			{
				Config: testDataSourceVolumesByTags_basic2(testName),
				Check: resource.ComposeTestCheckFunc(
					dc_matches.CheckResourceExists(),
					resource.TestCheckResourceAttr(filter_by_matches_all, "resources.0.resource_name", fmt.Sprintf("%s_0", testName)),
					resource.TestCheckResourceAttrSet(filter_by_matches_all, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(filter_by_matches_all, "resources.0.resource_detail.0.id"),
					resource.TestCheckResourceAttrSet(filter_by_matches_all, "resources.0.resource_detail.0.name"),
					resource.TestCheckResourceAttrSet(filter_by_matches_all, "resources.0.resource_detail.0.status"),
					resource.TestCheckResourceAttrSet(filter_by_matches_all, "resources.0.resource_detail.0.size"),
					resource.TestCheckResourceAttrSet(filter_by_matches_all, "resources.0.resource_detail.0.volume_type"),
					resource.TestCheckResourceAttrSet(filter_by_matches_all, "resources.0.resource_detail.0.availability_zone"),
					resource.TestCheckResourceAttrSet(filter_by_matches_all, "resources.0.resource_detail.0.created_at"),
					resource.TestCheckResourceAttrSet(filter_by_matches_all, "resources.0.resource_detail.0.updated_at"),
					resource.TestCheckResourceAttrSet(filter_by_matches_all, "resources.0.resource_detail.0.bootable"),
					resource.TestCheckResourceAttrSet(filter_by_matches_all, "resources.0.resource_detail.0.multiattach"),
					resource.TestCheckResourceAttr(filter_by_matches_all, "resources.0.tags.0.key", "foo1"),
					resource.TestCheckResourceAttr(filter_by_matches_all, "resources.0.tags.0.value", "bar1"),
					resource.TestCheckResourceAttr(filter_by_matches_all, "resources.0.tags.1.key", "foo0"),
					resource.TestCheckResourceAttr(filter_by_matches_all, "resources.0.tags.1.value", "bar0"),
				),
			},
		},
	},
	)
}

func testDataSourceVolumesByTags_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  count = 2
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = format("%[1]s_%%d", count.index)
  size              = 100
  volume_type       = "GPSSD"
  device_type       = "VBD"
  multiattach       = false

  tags = {
    "foo0" = "bar0"
    "foo1" = "bar1"
  }
}
`, name)
}

func testDataSourceVolumesByTags_basic1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_evs_volumes_by_tags" "filter_by_tags" {
  depends_on = [huaweicloud_evs_volume.test]
  action     = "filter"

  tags {
    key    = "foo0"
    values = ["bar0"]
  }

  tags {
    key    = "foo1"
    values = ["bar1"]
  }
}`, testDataSourceVolumesByTags_base(name))
}

func testDataSourceVolumesByTags_basic2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_evs_volumes_by_tags" "filter_by_matches" {
  depends_on = [huaweicloud_evs_volume.test]
  action     = "filter"

  tags {
    key    = "foo0"
    values = ["bar0"]
  }

  tags {
    key    = "foo1"
    values = ["bar1"]
  }

  matches {
    key    = "resource_name"
    value  = format("%%s_0", "%[2]s")
  }
}`, testDataSourceVolumesByTags_base(name), name)
}
