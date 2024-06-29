package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceAttachments_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_er_attachments.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byAttachmentId   = "data.huaweicloud_er_attachments.filter_by_attachment_id"
		dcByAttachmentId = acceptance.InitDataSourceCheck(byAttachmentId)

		byType   = "data.huaweicloud_er_attachments.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byName   = "data.huaweicloud_er_attachments.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byResourceId   = "data.huaweicloud_er_attachments.filter_by_resource_id"
		dcByResourceId = acceptance.InitDataSourceCheck(byResourceId)

		byStatus   = "data.huaweicloud_er_attachments.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byTags   = "data.huaweicloud_er_attachments.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAttachments_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "attachments.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Check whether filter parameter 'attachment_id' is effective.
					dcByAttachmentId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byAttachmentId, "attachments.#", "1"),
					resource.TestCheckResourceAttrPair(byAttachmentId, "attachments.0.id",
						"huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttrPair(byAttachmentId, "attachments.0.name",
						"huaweicloud_er_vpc_attachment.test", "name"),
					resource.TestCheckResourceAttrPair(byAttachmentId, "attachments.0.description",
						"huaweicloud_er_vpc_attachment.test", "description"),
					resource.TestCheckResourceAttrSet(byAttachmentId, "attachments.0.status"),
					resource.TestCheckResourceAttrSet(byAttachmentId, "attachments.0.associated"),
					resource.TestCheckResourceAttrPair(byAttachmentId, "attachments.0.resource_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestMatchResourceAttr(byAttachmentId, "attachments.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byAttachmentId, "attachments.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttr(byAttachmentId, "attachments.0.tags.%", "2"),
					resource.TestCheckResourceAttr(byAttachmentId, "attachments.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(byAttachmentId, "attachments.0.tags.key", "value"),
					resource.TestCheckResourceAttr(byAttachmentId, "attachments.0.type", "vpc"),
					resource.TestCheckOutput("is_attachment_id_filter_useful", "true"),
					// Check whether filter parameter 'type' is effective.
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					// Check whether filter parameter 'name' is effective.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Check whether filter parameter 'resource_id' is effective.
					dcByResourceId.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_id_filter_useful", "true"),
					// Check whether filter parameter 'status' is effective.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Check whether filter parameter 'tags' is effective.
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAttachments_base() string {
	var (
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)
	)

	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

%[1]s

resource "huaweicloud_er_instance" "test" {
  availability_zones    = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)
  name                  = "%[2]s"
  asn                   = %[3]d
  enterprise_project_id = "0"
}

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  name = "%[2]s"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_er_route_table" "test" {
  instance_id = huaweicloud_er_instance.test.id
  name        = "%[2]s"
}

resource "huaweicloud_er_static_route" "test" {
  route_table_id = huaweicloud_er_route_table.test.id
  destination    = huaweicloud_vpc.test.cidr
  attachment_id  = huaweicloud_er_vpc_attachment.test.id
}
`, common.TestVpc(name), name, bgpAsNum)
}

func testAccDataSourceAttachments_basic_step1() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_attachments" "test" {
  depends_on = [huaweicloud_er_vpc_attachment.test]

  instance_id = huaweicloud_er_instance.test.id
}

# Filter by attachment ID
locals {
  attachment_id = huaweicloud_er_vpc_attachment.test.id
}

data "huaweicloud_er_attachments" "filter_by_attachment_id" {
  depends_on = [huaweicloud_er_static_route.test]

  instance_id   = huaweicloud_er_instance.test.id
  attachment_id = local.attachment_id
}

locals {
  attachment_id_filter_result = [
    for v in data.huaweicloud_er_attachments.filter_by_attachment_id.attachments[*].id : v == local.attachment_id
  ]
}

output "is_attachment_id_filter_useful" {
  value = length(local.attachment_id_filter_result) > 0 && alltrue(local.attachment_id_filter_result)
}

# Filter by type
locals {
  # There is no attribute names 'type'.
  attachment_type = "vpc"
}

data "huaweicloud_er_attachments" "filter_by_type" {
  depends_on = [huaweicloud_er_vpc_attachment.test]

  instance_id = huaweicloud_er_instance.test.id
  type        = local.attachment_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_er_attachments.filter_by_type.attachments[*].type : v == local.attachment_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by name
locals {
  attachment_name = huaweicloud_er_vpc_attachment.test.name
}

data "huaweicloud_er_attachments" "filter_by_name" {
  # The behavior of parameter 'name' is 'Required', means this parameter does not have 'Know After Apply' behavior.
  depends_on = [huaweicloud_er_vpc_attachment.test]

  instance_id = huaweicloud_er_instance.test.id
  name        = local.attachment_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_er_attachments.filter_by_name.attachments[*].name : v == local.attachment_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by resource ID
locals {
  resource_id = huaweicloud_vpc.test.id
}

data "huaweicloud_er_attachments" "filter_by_resource_id" {
  # The behavior of parameter 'type' is 'Required', means this parameter does not have 'Know After Apply' behavior.
  depends_on = [huaweicloud_er_vpc_attachment.test]

  instance_id = huaweicloud_er_instance.test.id
  resource_id = local.resource_id
}

locals {
  resource_id_filter_result = [
    for v in data.huaweicloud_er_attachments.filter_by_resource_id.attachments[*].resource_id : v == local.resource_id
  ]
}

output "is_resource_id_filter_useful" {
  value = length(local.resource_id_filter_result) > 0 && alltrue(local.resource_id_filter_result)
}

# Filter by status
locals {
  attachment_status = huaweicloud_er_vpc_attachment.test.status
}

data "huaweicloud_er_attachments" "filter_by_status" {
  instance_id = huaweicloud_er_instance.test.id
  status      = local.attachment_status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_er_attachments.filter_by_status.attachments[*].status : v == local.attachment_status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by tags
locals {
  attachment_tags = huaweicloud_er_vpc_attachment.test.tags
}

data "huaweicloud_er_attachments" "filter_by_tags" {
  depends_on = [huaweicloud_er_vpc_attachment.test]

  instance_id = huaweicloud_er_instance.test.id
  tags        = local.attachment_tags
}

locals {
  tags_filter_result = [
    for v in data.huaweicloud_er_attachments.filter_by_tags.attachments[*].tags : length(v) == length(local.attachment_tags) &&
    length(v) == length(merge(v, local.attachment_tags))
  ]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}
`, testAccDataSourceAttachments_base())
}
