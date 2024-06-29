package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstances_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_er_instances.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byInstanceId   = "data.huaweicloud_er_instances.filter_by_instance_id"
		dcByInstanceId = acceptance.InitDataSourceCheck(byInstanceId)

		byName   = "data.huaweicloud_er_instances.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byEpsId   = "data.huaweicloud_er_instances.filter_by_eps_id"
		dcByEpsId = acceptance.InitDataSourceCheck(byEpsId)

		byOwnedBySelf   = "data.huaweicloud_er_instances.filter_by_owned_by_self"
		dcByOwnedBySelf = acceptance.InitDataSourceCheck(byOwnedBySelf)

		byStatus   = "data.huaweicloud_er_instances.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byTags   = "data.huaweicloud_er_instances.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "instances.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Check whether filter parameter 'instance_id' is effective.
					dcByInstanceId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byInstanceId, "instances.#", "1"),
					resource.TestCheckResourceAttrPair(byInstanceId, "instances.0.id", "huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(byInstanceId, "instances.0.asn", "huaweicloud_er_instance.test", "asn"),
					resource.TestCheckResourceAttrPair(byInstanceId, "instances.0.name", "huaweicloud_er_instance.test", "name"),
					resource.TestCheckResourceAttrPair(byInstanceId, "instances.0.description", "huaweicloud_er_instance.test", "description"),
					resource.TestCheckResourceAttrPair(byInstanceId, "instances.0.enterprise_project_id",
						"huaweicloud_er_instance.test", "enterprise_project_id"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.tags.%", "2"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.tags.key", "value"),
					resource.TestMatchResourceAttr(byInstanceId, "instances.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.enable_default_propagation", "true"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.enable_default_association", "true"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.auto_accept_shared_attachments", "true"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.default_propagation_route_table_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.default_association_route_table_id"),
					resource.TestMatchResourceAttr(byInstanceId, "instances.0.availability_zones.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),
					// Check whether filter parameter 'name' is effective.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Check whether filter parameter 'enterprise_project_id' is effective.
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					// Check whether filter parameter 'owned_by_self' is effective.
					dcByOwnedBySelf.CheckResourceExists(),
					resource.TestCheckOutput("is_owned_by_self_filter_useful", "true"),
					// Check whether filter parameter 'status' is effective.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					// Check whether filter parameter 'tags' is effective.
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceInstances_base() string {
	var (
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)
	)

	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones    = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)
  name                  = "%[1]s"
  asn                   = %[2]d
  description           = "Created by script"
  enterprise_project_id = "0"

  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, bgpAsNum)
}

func testAccDataSourceInstances_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_instances" "test" {
  depends_on = [huaweicloud_er_instance.test]
}

# Filter by instance ID
locals {
  instance_id = huaweicloud_er_instance.test.id
}

data "huaweicloud_er_instances" "filter_by_instance_id" {
  instance_id = local.instance_id
}

locals {
  instance_id_filter_result = [
    for v in data.huaweicloud_er_instances.filter_by_instance_id.instances[*].id : v == local.instance_id
  ]
}

output "is_instance_id_filter_useful" {
  value = length(local.instance_id_filter_result) > 0 && alltrue(local.instance_id_filter_result)
}

# Filter by name
locals {
  instance_name = huaweicloud_er_instance.test.name
}

data "huaweicloud_er_instances" "filter_by_name" {
  depends_on = [huaweicloud_er_instance.test]

  name = local.instance_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_er_instances.filter_by_name.instances[*].name : v == local.instance_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by enterprise project ID
locals {
  eps_id = huaweicloud_er_instance.test.enterprise_project_id
}

data "huaweicloud_er_instances" "filter_by_eps_id" {
  depends_on = [huaweicloud_er_instance.test]
  
  enterprise_project_id = local.eps_id
}

locals {
  eps_id_filter_result = [
    for v in data.huaweicloud_er_instances.filter_by_eps_id.instances[*].enterprise_project_id : v == local.eps_id
  ]
}

output "is_eps_id_filter_useful" {
  value = length(local.eps_id_filter_result) > 0 && alltrue(local.eps_id_filter_result)
}

# Filter by owned_by_self param
data "huaweicloud_er_instances" "filter_by_owned_by_self" {
  depends_on = [huaweicloud_er_instance.test]

  owned_by_self = true
}

output "is_owned_by_self_filter_useful" {
  value = contains(data.huaweicloud_er_instances.filter_by_owned_by_self.instances[*].id, local.instance_id)
}

# Filter by status
locals {
  instance_status = huaweicloud_er_instance.test.status
}

data "huaweicloud_er_instances" "filter_by_status" {
  status = local.instance_status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_er_instances.filter_by_status.instances[*].status : v == local.instance_status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by tags
locals {
  instance_tags = huaweicloud_er_instance.test.tags
}

data "huaweicloud_er_instances" "filter_by_tags" {
  depends_on = [huaweicloud_er_instance.test]

  tags = local.instance_tags
}

locals {
  tags_filter_result = [
    for v in data.huaweicloud_er_instances.filter_by_tags.instances[*].tags : length(v) == length(local.instance_tags) &&
    length(v) == length(merge(v, local.instance_tags))
  ]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}
`, testAccDataSourceInstances_base())
}
