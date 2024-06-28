package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssociations_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		baseConfig = testAccDataSourceAssociations_base(name)

		all = "data.huaweicloud_er_associations.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byAttachmentId   = "data.huaweicloud_er_associations.filter_by_attachment_id"
		dcByAttachmentId = acceptance.InitDataSourceCheck(byAttachmentId)

		byAttachmentType   = "data.huaweicloud_er_associations.filter_by_attachment_type"
		dcByAttachmentType = acceptance.InitDataSourceCheck(byAttachmentType)

		byStatus   = "data.huaweicloud_er_associations.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byNotFoundInstanceId   = "data.huaweicloud_er_associations.instance_id_not_found"
		dcByNotFoundInstanceId = acceptance.InitDataSourceCheck(byNotFoundInstanceId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssociations_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "associations.#"),
					resource.TestCheckResourceAttrSet(all, "associations.0.resource_id"),
					resource.TestMatchResourceAttr(all, "associations.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "associations.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Check whether filter parameter 'status' is effective.
					dcByAttachmentId.CheckResourceExists(),
					resource.TestCheckOutput("is_attachment_id_filter_useful", "true"),
					// Check whether filter parameter 'attachment_type' is effective.
					dcByAttachmentType.CheckResourceExists(),
					resource.TestCheckOutput("is_attachment_type_filter_useful", "true"),
					// Check whether filter parameter 'status' is effective.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
			{
				// Checks whether the resource returns the expected empty list when the instance ID does not exist.
				Config: testAccDataSourceAssociations_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dcByNotFoundInstanceId.CheckResourceExists(),
					// If the instance ID does not exist, the data source will not report the error.
					// Just return an empty list.
					resource.TestCheckResourceAttr(byNotFoundInstanceId, "associations.#", "0"),
				),
			},
			{
				// Checks whether the resource returns the expected error when the route table ID does not exist.
				// Please ensure that the test account has 'ER FullAccess' permission for version 5.0.
				Config: testAccDataSourceAssociations_basic_step3(baseConfig),
				// If the routing table ID does not exist, the data source will report an error: 'route table {uuid} not found'.
				ExpectError: regexp.MustCompile(`route table [a-f0-9-]+ not found`),
			},
		},
	})
}

func testAccDataSourceAssociations_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_association" "test" {
  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
  attachment_id  = huaweicloud_er_vpc_attachment.test.id
}
`, testAccAssociation_base(name))
}

func testAccDataSourceAssociations_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_associations" "test" {
  depends_on = [
    huaweicloud_er_association.test,
  ]

  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
}

locals {
  attachment_id = data.huaweicloud_er_associations.test.associations[0].attachment_id
}

data "huaweicloud_er_associations" "filter_by_attachment_id" {
  depends_on = [
    huaweicloud_er_association.test,
  ]

  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
  attachment_id  = local.attachment_id
}

locals {
 attachment_id_filter_result = [
    for v in data.huaweicloud_er_associations.filter_by_attachment_id.associations[*].attachment_id : 
    v == local.attachment_id
  ]
}

output "is_attachment_id_filter_useful" {
  value = alltrue(local.attachment_id_filter_result) && length(local.attachment_id_filter_result) > 0
}

locals {
  attachment_type = data.huaweicloud_er_associations.test.associations[0].attachment_type
}

data "huaweicloud_er_associations" "filter_by_attachment_type" {
  depends_on = [
    huaweicloud_er_association.test,
  ]

  instance_id     = huaweicloud_er_instance.test.id
  route_table_id  = huaweicloud_er_route_table.test.id
  attachment_type = local.attachment_type
}

locals {
  attachment_type_filter_result = [
    for v in data.huaweicloud_er_associations.filter_by_attachment_type.associations[*].attachment_type : 
    v == local.attachment_type
  ]
}
   
output "is_attachment_type_filter_useful" {
  value = alltrue(local.attachment_type_filter_result) && length(local.attachment_type_filter_result) > 0
}

locals {
  status = data.huaweicloud_er_associations.test.associations[0].status
}
  
data "huaweicloud_er_associations" "filter_by_status" {
  depends_on = [
    huaweicloud_er_association.test,
  ]

  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
  status         = local.status
}
  
locals {
  status_filter_result = [
    for v in data.huaweicloud_er_associations.filter_by_status.associations[*].status : v == local.status
  ]
}
   
output "is_status_filter_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, baseConfig)
}

func testAccDataSourceAssociations_basic_step2(baseConfig string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_associations" "instance_id_not_found" {
  depends_on = [
    huaweicloud_er_association.test,
  ]

  instance_id    = "%[2]s"
  route_table_id = huaweicloud_er_route_table.test.id
}
`, baseConfig, randUUID)
}

func testAccDataSourceAssociations_basic_step3(baseConfig string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_associations" "route_table_id_not_found" {
  depends_on = [
    huaweicloud_er_association.test,
  ]

  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = "%[2]s"
}
`, baseConfig, randUUID)
}
