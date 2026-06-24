package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceInstancesByTags_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_instances_by_tags.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		name       = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceInstancesByTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.tags.#"),

					resource.TestCheckOutput("results_is_not_empty", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceInstancesByTags_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "test_1234"
  mode              = "Cluster"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "16"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceInstancesByTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_geminidb_instances_by_tags" "test" {
  action = "filter"

  depends_on = [huaweicloud_geminidb_instance.test]
}

data "huaweicloud_geminidb_instances_by_tags" "filter_by_count" {
  action = "count"

  depends_on = [huaweicloud_geminidb_instance.test]
}

data "huaweicloud_geminidb_instances_by_tags" "filter_by_matches" {
  action = "filter"

  matches {
    key   = "instance_id"
    value = huaweicloud_geminidb_instance.test.id
  }
}

data "huaweicloud_geminidb_instances_by_tags" "filter_by_tags" {
  action = "filter"

  tags {
    key    = data.huaweicloud_geminidb_instances_by_tags.test.instances.0.tags.0.key
    values = [data.huaweicloud_geminidb_instances_by_tags.test.instances.0.tags.0.value]
  }
}


output "results_is_not_empty" {
  value = data.huaweicloud_geminidb_instances_by_tags.filter_by_count.total_count > 0
}

output "matches_filter_is_useful" {
  value = length(data.huaweicloud_geminidb_instances_by_tags.filter_by_matches.instances) > 0
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_geminidb_instances_by_tags.filter_by_tags.instances) > 0
}
`, testDataSourceInstancesByTags_base(name))
}
