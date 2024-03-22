package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePropagations_basic(t *testing.T) {
	var (
		dcName   = "data.huaweicloud_er_propagations.test"
		name     = acceptance.RandomAccResourceName()
		dc       = acceptance.InitDataSourceCheck(dcName)
		bgpAsNum = acctest.RandIntRange(64512, 65534)

		byInstanceId   = "data.huaweicloud_er_propagations.not_found_instance_id"
		dcByInstanceId = acceptance.InitDataSourceCheck(byInstanceId)

		byRouteTableId   = "data.huaweicloud_er_propagations.not_found_route_table_id"
		dcByRouteTableId = acceptance.InitDataSourceCheck(byRouteTableId)

		byAttachmentId   = "data.huaweicloud_er_propagations.filter_by_attachment_id"
		dcByAttachmentId = acceptance.InitDataSourceCheck(byAttachmentId)

		byAttachmentType   = "data.huaweicloud_er_propagations.filter_by_attachment_type"
		dcByAttachmentType = acceptance.InitDataSourceCheck(byAttachmentType)

		byStatus   = "data.huaweicloud_er_propagations.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePropagations_basic(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					dcByInstanceId.CheckResourceExists(),
					resource.TestCheckOutput("instance_id_not_found", "true"),

					dcByRouteTableId.CheckResourceExists(),
					resource.TestCheckOutput("route_table_id_not_found", "true"),
					resource.TestCheckResourceAttrSet(byAttachmentId, "propagations.#"),
					resource.TestCheckResourceAttrSet(byAttachmentId, "propagations.0.resource_id"),
					resource.TestCheckResourceAttrSet(byAttachmentId, "propagations.0.created_at"),
					resource.TestCheckResourceAttrSet(byAttachmentId, "propagations.0.updated_at"),

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
		},
	})
}

func testAccDatasourcePropagations_basic(name string, bgpAsNum int) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_propagations" "not_found_instance_id" {
  depends_on = [
    huaweicloud_er_propagation.test,
  ]

  instance_id    = "%[2]s"
  route_table_id = huaweicloud_er_route_table.test.id
}
  
output "instance_id_not_found" {
  value = length(data.huaweicloud_er_propagations.not_found_instance_id.propagations) == 0
}
  
data "huaweicloud_er_propagations" "not_found_route_table_id" {
  depends_on = [
    huaweicloud_er_propagation.test,
  ]

  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = "%[2]s"
}
  
output "route_table_id_not_found" {
  value = length(data.huaweicloud_er_propagations.not_found_route_table_id.propagations) == 0
}

data "huaweicloud_er_propagations" "test" {
  depends_on = [
    huaweicloud_er_propagation.test,
  ]

  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
}

locals {
  attachment_id = data.huaweicloud_er_propagations.test.propagations[0].attachment_id
}

data "huaweicloud_er_propagations" "filter_by_attachment_id" {
  depends_on = [
    huaweicloud_er_propagation.test,
  ]

  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
  attachment_id  = local.attachment_id
}

locals {
 attachment_id_filter_result = [
    for v in data.huaweicloud_er_propagations.filter_by_attachment_id.propagations[*].attachment_id : 
    v == local.attachment_id
  ]
}

output "is_attachment_id_filter_useful" {
  value = alltrue(local.attachment_id_filter_result) && length(local.attachment_id_filter_result) > 0
}

locals {
  attachment_type = data.huaweicloud_er_propagations.test.propagations[0].attachment_type
}

data "huaweicloud_er_propagations" "filter_by_attachment_type" {
  depends_on = [
    huaweicloud_er_propagation.test,
  ]

  instance_id     = huaweicloud_er_instance.test.id
  route_table_id  = huaweicloud_er_route_table.test.id
  attachment_type = local.attachment_type
}

locals {
  attachment_type_filter_result = [
    for v in data.huaweicloud_er_propagations.filter_by_attachment_type.propagations[*].attachment_type : 
    v == local.attachment_type
  ]
}
   
output "is_attachment_type_filter_useful" {
  value = alltrue(local.attachment_type_filter_result) && length(local.attachment_type_filter_result) > 0
}

locals {
  status = data.huaweicloud_er_propagations.test.propagations[0].status
}
  
data "huaweicloud_er_propagations" "filter_by_status" {
  depends_on = [
    huaweicloud_er_propagation.test,
  ]

  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
  status         = local.status
}
  
locals {
  status_filter_result = [
    for v in data.huaweicloud_er_propagations.filter_by_status.propagations[*].status : v == local.status
  ]
}
   
output "is_status_filter_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, testAccPropagation_basic(name, bgpAsNum), randUUID)
}
