package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmnLogtanks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_smn_logtanks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSmnLogtanks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "logtanks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "logtanks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "logtanks.0.log_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "logtanks.0.log_stream_id"),
					resource.TestCheckResourceAttrSet(dataSource, "logtanks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "logtanks.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceSmnLogtanks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_smn_logtanks" "test" {
  depends_on = [huaweicloud_smn_logtank.test]

  topic_urn = huaweicloud_smn_topic.test.id
}
`, testAccSMNV2LogtankConfig_buildData(name))
}

func testAccSMNV2LogtankConfig_buildData(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}

resource "huaweicloud_smn_logtank" "test" {
  topic_urn     = huaweicloud_smn_topic.test.topic_urn
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
}
`, rName)
}
