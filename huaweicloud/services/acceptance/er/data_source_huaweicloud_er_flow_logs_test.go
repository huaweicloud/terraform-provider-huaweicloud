package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceFlowLogs_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		bgpAsNum       = acctest.RandIntRange(64512, 65534)
		dataSourceName = "data.huaweicloud_er_flow_logs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byResourceType   = "data.huaweicloud_er_flow_logs.filter_by_resource_type"
		dcByResourceType = acceptance.InitDataSourceCheck(byResourceType)

		byResourceId   = "data.huaweicloud_er_flow_logs.filter_by_resource_id"
		dcByResourceId = acceptance.InitDataSourceCheck(byResourceId)

		byFlowLogId   = "data.huaweicloud_er_flow_logs.filter_by_flow_log_id"
		dcByFlowLogId = acceptance.InitDataSourceCheck(byFlowLogId)

		byName   = "data.huaweicloud_er_flow_logs.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_er_flow_logs.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byEnabled   = "data.huaweicloud_er_flow_logs.filter_by_enabled"
		dcByEnabled = acceptance.InitDataSourceCheck(byEnabled)

		byLogGroupId   = "data.huaweicloud_er_flow_logs.filter_by_log_group_id"
		dcByLogGroupId = acceptance.InitDataSourceCheck(byLogGroupId)

		byLogStreamId   = "data.huaweicloud_er_flow_logs.filter_by_log_stream_id"
		dcByLogStreamId = acceptance.InitDataSourceCheck(byLogStreamId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceFlowLogs_basic(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("resource_type_filter_is_useful", "true"),

					dcByResourceId.CheckResourceExists(),
					resource.TestCheckOutput("resource_id_filter_is_useful", "true"),

					dcByFlowLogId.CheckResourceExists(),
					resource.TestCheckOutput("flow_log_id_filter_is_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByEnabled.CheckResourceExists(),
					resource.TestCheckOutput("enabled_filter_is_useful", "true"),

					dcByLogGroupId.CheckResourceExists(),
					resource.TestCheckOutput("log_group_id_filter_is_useful", "true"),

					dcByLogStreamId.CheckResourceExists(),
					resource.TestCheckOutput("log_stream_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccFlowLogsDataSource_base(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)
  name               = "%[2]s"
  asn                = %[3]d
}

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 7
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id            = huaweicloud_er_instance.test.id
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  name                   = "%[2]s"
  auto_create_vpc_routes = true

  tags = {
    foo = "bar"
  }
}

resource "huaweicloud_er_flow_log" "test" {
  instance_id    = huaweicloud_er_instance.test.id
  log_store_type = "LTS"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  resource_type  = "attachment"
  resource_id    = huaweicloud_er_vpc_attachment.test.id
  name           = "%[2]s"
  description    = "Create ER flow log"
}
`, common.TestVpc(name), name, bgpAsNum)
}

func testAccDatasourceFlowLogs_basic(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_er_flow_logs" "test" {
  depends_on  = [huaweicloud_er_flow_log.test]
  instance_id = huaweicloud_er_instance.test.id
}

locals {
  resource_type = data.huaweicloud_er_flow_logs.test.flow_logs[0].resource_type
}

data "huaweicloud_er_flow_logs" "filter_by_resource_type" {
  instance_id   = huaweicloud_er_instance.test.id
  resource_type = local.resource_type
}

locals {
  resource_type_filter_result = [
    for v in data.huaweicloud_er_flow_logs.filter_by_resource_type.flow_logs[*].resource_type : 
    v == local.resource_type
  ]
}

output "resource_type_filter_is_useful" {
  value = alltrue(local.resource_type_filter_result) && length(local.resource_type_filter_result) > 0
}

locals {
  resource_id = data.huaweicloud_er_flow_logs.test.flow_logs[0].resource_id
}

data "huaweicloud_er_flow_logs" "filter_by_resource_id" {
  instance_id = huaweicloud_er_instance.test.id
  resource_id = local.resource_id
}

locals {
  resource_id_filter_result = [
    for v in data.huaweicloud_er_flow_logs.filter_by_resource_id.flow_logs[*].resource_id : v == local.resource_id
  ]
}

output "resource_id_filter_is_useful" {
  value = alltrue(local.resource_id_filter_result) && length(local.resource_id_filter_result) > 0
}

locals {
  flow_log_id = data.huaweicloud_er_flow_logs.test.flow_logs[0].id
}

data "huaweicloud_er_flow_logs" "filter_by_flow_log_id" {
  instance_id = huaweicloud_er_instance.test.id
  flow_log_id = local.flow_log_id
}

locals {
  flow_log_id_filter_result = [
    for v in data.huaweicloud_er_flow_logs.filter_by_flow_log_id.flow_logs[*].id : v == local.flow_log_id
  ]
}

output "flow_log_id_filter_is_useful" {
  value = alltrue(local.flow_log_id_filter_result) && length(local.flow_log_id_filter_result) > 0
}

locals {
  name = data.huaweicloud_er_flow_logs.test.flow_logs[0].name
}

data "huaweicloud_er_flow_logs" "filter_by_name" {
  instance_id = huaweicloud_er_instance.test.id
  name        = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_er_flow_logs.filter_by_name.flow_logs[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

locals {
  status = data.huaweicloud_er_flow_logs.test.flow_logs[0].status
}

data "huaweicloud_er_flow_logs" "filter_by_status" {
  instance_id = huaweicloud_er_instance.test.id
  status      = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_er_flow_logs.filter_by_status.flow_logs[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

locals {
  log_group_id = data.huaweicloud_er_flow_logs.test.flow_logs[0].log_group_id
}

data "huaweicloud_er_flow_logs" "filter_by_log_group_id" {
  instance_id  = huaweicloud_er_instance.test.id
  log_group_id = local.log_group_id
}

locals {
  log_group_id_filter_result = [
    for v in data.huaweicloud_er_flow_logs.filter_by_log_group_id.flow_logs[*].log_group_id : v == local.log_group_id
  ]
}

output "log_group_id_filter_is_useful" {
  value = alltrue(local.log_group_id_filter_result) && length(local.log_group_id_filter_result) > 0
}

locals {
  log_stream_id = data.huaweicloud_er_flow_logs.test.flow_logs[0].log_stream_id
}

data "huaweicloud_er_flow_logs" "filter_by_log_stream_id" {
  instance_id   = huaweicloud_er_instance.test.id
  log_stream_id = local.log_stream_id
}

locals {
  log_stream_id_filter_result = [
    for v in data.huaweicloud_er_flow_logs.filter_by_log_stream_id.flow_logs[*].log_stream_id : 
    v == local.log_stream_id
  ]
}

output "log_stream_id_filter_is_useful" {
  value = alltrue(local.log_stream_id_filter_result) && length(local.log_stream_id_filter_result) > 0
}

locals {
  enabled = data.huaweicloud_er_flow_logs.test.flow_logs[0].enabled
}

data "huaweicloud_er_flow_logs" "filter_by_enabled" {
  instance_id = huaweicloud_er_instance.test.id
  enabled     = local.enabled
}

locals {
  enabled_filter_result = [
    for v in data.huaweicloud_er_flow_logs.filter_by_enabled.flow_logs[*].enabled : v == local.enabled
  ]
}

output "enabled_filter_is_useful" {
  value = alltrue(local.enabled_filter_result) && length(local.enabled_filter_result) > 0
}
`, testAccFlowLogsDataSource_base(name, bgpAsNum))
}
