package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTopics_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dms_kafka_topics.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dms_kafka_topics.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTopics_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "topics.#", regexp.MustCompile(`^[1-9]\d*$`)),
					// Filter by 'name' parameter.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestMatchResourceAttr(byName, "max_partitions", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(byName, "topic_max_partitions", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet(byName, "topics.0.name"),
					resource.TestMatchResourceAttr(byName, "topics.0.partitions", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(byName, "topics.0.replicas", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(byName, "topics.0.aging_time", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet(byName, "topics.0.sync_replication"),
					resource.TestCheckResourceAttrSet(byName, "topics.0.sync_flushing"),
					resource.TestCheckResourceAttrSet(byName, "topics.0.description"),
					resource.TestCheckResourceAttrSet(byName, "topics.0.configs.#"),
					resource.TestCheckResourceAttrSet(byName, "topics.0.configs.0.name"),
					resource.TestCheckResourceAttrSet(byName, "topics.0.configs.0.value"),
					resource.TestCheckResourceAttrSet(byName, "topics.0.policies_only"),
					resource.TestCheckResourceAttrSet(byName, "topics.0.type"),
					resource.TestMatchResourceAttr(byName, "topics.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourceTopics_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id      = "%[1]s"
  name             = "%[2]s"
  partitions       = 4
  aging_time       = 36
  description      = "Created by Terraform script"
  sync_replication = true
  sync_flushing    = true
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func testDataSourceTopics_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_dms_kafka_topics" "all" {
  depends_on = [huaweicloud_dms_kafka_topic.test]

  instance_id = "%[2]s"
}

# Filter by 'name' parameter.
locals {
  topic_name = huaweicloud_dms_kafka_topic.test.name
}

data "huaweicloud_dms_kafka_topics" "filter_by_name" {
  depends_on = [huaweicloud_dms_kafka_topic.test]

  instance_id = "%[2]s"
  name        = local.topic_name
}

locals {
  name_filter_result = [for v in data.huaweicloud_dms_kafka_topics.filter_by_name.topics[*].name : strcontains(v, local.topic_name)]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}
`, testAccDataSourceTopics_basic(name), acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}
