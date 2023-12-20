package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccAttachmentsDataSource_basic(t *testing.T) {
	var (
		dName = "data.huaweicloud_er_attachments.filter_by_name"
		name  = acceptance.RandomAccResourceName()

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAttachmentsDataSource_filterByName(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccAttachmentsDataSource_base(name string) string {
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_er_instance" "test" {
  availability_zones    = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
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
`, common.TestVpc(name), name, bgpAsNum)
}

func testAccAttachmentsDataSource_filterByName(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_attachments" "filter_by_name" {
  // The behavior of parameter 'name' is 'Required', means this parameter does not have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]

  instance_id = huaweicloud_er_instance.test.id
  name        = huaweicloud_er_vpc_attachment.test.name
}

data "huaweicloud_er_attachments" "not_found" {
  // Since a specified name is used, there is no dependency relationship with resource attachment, and the dependency
  // needs to be manually set.
  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]

  instance_id = huaweicloud_er_instance.test.id
  name        = "resource_not_found"
}

locals {
  filter_result = [for v in data.huaweicloud_er_attachments.filter_by_name.attachments[*].id : v == huaweicloud_er_vpc_attachment.test.id]
}

output "is_name_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_er_attachments.not_found.attachments) == 0
}
`, testAccAttachmentsDataSource_base(name))
}

func TestAccAttachmentsDataSource_filterById(t *testing.T) {
	var (
		dName = "data.huaweicloud_er_attachments.filter_by_id"
		name  = acceptance.RandomAccResourceName()

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAttachmentsDataSource_filterById(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccAttachmentsDataSource_filterById(name string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_attachments" "filter_by_id" {
  instance_id   = huaweicloud_er_instance.test.id
  attachment_id = huaweicloud_er_vpc_attachment.test.id
}

data "huaweicloud_er_attachments" "not_found" {
  // Since a random ID is used, there is no dependency relationship with resource attachment, and the dependency needs
  // to be manually set.
  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]

  instance_id   = huaweicloud_er_instance.test.id
  attachment_id = "%[2]s"
}

locals {
  filter_result = [for v in data.huaweicloud_er_attachments.filter_by_id.attachments[*].id : v == huaweicloud_er_vpc_attachment.test.id]
}

output "is_id_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_er_attachments.not_found.attachments) == 0
}
`, testAccAttachmentsDataSource_base(name), randUUID)
}

func TestAccAttachmentsDataSource_filterByType(t *testing.T) {
	var (
		dName = "data.huaweicloud_er_attachments.filter_by_type"
		name  = acceptance.RandomAccResourceName()

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAttachmentsDataSource_filterByType(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccAttachmentsDataSource_filterByType(name string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_attachments" "filter_by_type" {
  // Since a specified type is used, there is no dependency relationship with resource attachment, and the dependency
  // needs to be manually set.
  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]

  instance_id = huaweicloud_er_instance.test.id
  type        = "vpc"
}

data "huaweicloud_er_attachments" "not_found" {
  // Since a specified type is used, there is no dependency relationship with resource attachment, and the dependency
  // needs to be manually set.
  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]

  instance_id = huaweicloud_er_instance.test.id
  type        = "vgw"
}

locals {
  filter_result = [for v in data.huaweicloud_er_attachments.filter_by_type.attachments[*].id : v == huaweicloud_er_vpc_attachment.test.id]
}

output "is_type_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_er_attachments.not_found.attachments) == 0
}
`, testAccAttachmentsDataSource_base(name), randUUID)
}

func TestAccAttachmentsDataSource_filterByStatus(t *testing.T) {
	var (
		dName = "data.huaweicloud_er_attachments.filter_by_status"
		name  = acceptance.RandomAccResourceName()

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAttachmentsDataSource_filterByStatus(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccAttachmentsDataSource_filterByStatus(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_attachments" "filter_by_status" {
  instance_id = huaweicloud_er_instance.test.id
  status      = huaweicloud_er_vpc_attachment.test.status
}

data "huaweicloud_er_attachments" "not_found" {
  // Since a specified status is used, there is no dependency relationship with resource attachment, and the dependency needs
  // to be manually set.
  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]

  instance_id   = huaweicloud_er_instance.test.id
  status        = "failed"
}

locals {
  filter_result = [for v in data.huaweicloud_er_attachments.filter_by_status.attachments[*].id : v == huaweicloud_er_vpc_attachment.test.id]
}

output "is_status_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_er_attachments.not_found.attachments) == 0
}
`, testAccAttachmentsDataSource_base(name))
}

func TestAccAttachmentsDataSource_filterByTags(t *testing.T) {
	var (
		dName = "data.huaweicloud_er_attachments.filter_by_tags"
		name  = acceptance.RandomAccResourceName()

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAttachmentsDataSource_filterByTags(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_is_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccAttachmentsDataSource_filterByTags(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_attachments" "filter_by_tags" {
  // Since a specified key/value pair is used, there is no dependency relationship with resource attachment, and the
  // dependency needs to be manually set.
  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]

  instance_id = huaweicloud_er_instance.test.id

  tags = {
    foo = "bar"
  }
}

data "huaweicloud_er_attachments" "not_found" {
  // Since a specified key/value pair is used, there is no dependency relationship with resource attachment, and the
  // dependency needs to be manually set.
  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]

  instance_id = huaweicloud_er_instance.test.id

  tags = {
    owner = "terraform"
  }
}

locals {
  filter_result = [for v in data.huaweicloud_er_attachments.filter_by_tags.attachments[*].id : v == huaweicloud_er_vpc_attachment.test.id]
}

output "is_tags_filter_is_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_er_attachments.not_found.attachments) == 0
}
`, testAccAttachmentsDataSource_base(name))
}

func TestAccAttachmentsDataSource_filterByResourceId(t *testing.T) {
	var (
		dName = "data.huaweicloud_er_attachments.filter_by_resource_id"
		name  = acceptance.RandomAccResourceName()

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAttachmentsDataSource_filterByResourceId(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccAttachmentsDataSource_filterByResourceId(name string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_attachments" "filter_by_resource_id" {
  // Since a specified resource ID is used, there is no dependency relationship with resource attachment, and the
  // dependency needs to be manually set.
  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]

  instance_id = huaweicloud_er_instance.test.id
  resource_id = huaweicloud_vpc.test.id
}

data "huaweicloud_er_attachments" "not_found" {
  // Since a specified resource ID is used, there is no dependency relationship with resource attachment, and the
  // dependency needs to be manually set.
  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]

  instance_id = huaweicloud_er_instance.test.id
  resource_id = "%[2]s"
}

locals {
  filter_result = [for v in data.huaweicloud_er_attachments.filter_by_resource_id.attachments[*].id : v == huaweicloud_er_vpc_attachment.test.id]
}

output "is_resource_id_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_er_attachments.not_found.attachments) == 0
}
`, testAccAttachmentsDataSource_base(name), randUUID)
}
