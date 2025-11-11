package cbh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceInstancesByTags_basic(t *testing.T) {
	var (
		rName        = "data.huaweicloud_cbh_instances_by_tags.test"
		rTagsName    = "data.huaweicloud_cbh_instances_by_tags.test_tags"
		rMatchesName = "data.huaweicloud_cbh_instances_by_tags.test_matches"
		dc           = acceptance.InitDataSourceCheck(rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Prepare CBH instances in the test environment in advance.
			acceptance.TestAccPreCheckCbhInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceInstancesByTags_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "total_count"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.alter_permit"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.az_info.#"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.bastion_version"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.created_time"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.instance_id"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.name"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.network.#"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.period_num"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.resource_info.#"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.server_id"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.status_info.#"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.update"),
					resource.TestCheckResourceAttrSet(rName, "resources.0.resource_detail.0.upgrade_time"),

					resource.TestCheckResourceAttr(rTagsName, "resources.#", "0"),
					resource.TestCheckResourceAttr(rMatchesName, "resources.#", "0"),
				),
			},
		},
	})
}

const testAccDatasourceInstancesByTags_basic = `
data "huaweicloud_cbh_instances_by_tags" "test" {}

data "huaweicloud_cbh_instances_by_tags" "test_tags" {
  tags {
    key    = "non-exist-key"
    values = ["non-exist-value"]
  }
}

data "huaweicloud_cbh_instances_by_tags" "test_matches" {
  matches {
    key   = "resource_name"
    value = "non-exist-value"
  }
}
`
