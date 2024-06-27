package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssociations_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		bgpAsNum   = acctest.RandIntRange(64512, 65534)
		baseConfig = testAccAssociation_basic(name, bgpAsNum)

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

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAssociations_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "associations.#"),
					resource.TestCheckResourceAttrSet(all, "associations.0.resource_id"),
					resource.TestCheckResourceAttrSet(all, "associations.0.created_at"),
					resource.TestCheckResourceAttrSet(all, "associations.0.updated_at"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),

					dcByAttachmentId.CheckResourceExists(),
					resource.TestCheckOutput("is_attachment_id_filter_useful", "true"),

					dcByAttachmentType.CheckResourceExists(),
					resource.TestCheckOutput("is_attachment_type_filter_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
			// If the instance ID does not exist, the data source will not report the error.
			// Just return an empty list.
			{
				Config: testAccDatasourceAssociations_instanceIdNotFound(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dcByNotFoundInstanceId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byNotFoundInstanceId, "associations.#", "0"),
				),
			},
			// If the routing table ID does not exist, the data source will report an error: 'route table {uuid} not found'.
			{
				Config:      testAccDatasourceAssociations_routeTableIdNotFound(baseConfig),
				ExpectError: regexp.MustCompile(`route table [a-f0-9-]+ not found`),
			},
		},
	})
}

func testAccDatasourceAssociations_basic(baseConfig string) string {
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

func testAccDatasourceAssociations_instanceIdNotFound(baseConfig string) string {
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

func testAccDatasourceAssociations_routeTableIdNotFound(baseConfig string) string {
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
