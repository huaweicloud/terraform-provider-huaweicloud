package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceArchitectureSubjects_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_architecture_subjects.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_architecture_subjects.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_dataarts_architecture_subjects.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byCreateBy   = "data.huaweicloud_dataarts_architecture_subjects.filter_by_create_by"
		dcByCreateBy = acceptance.InitDataSourceCheck(byCreateBy)

		byLevel   = "data.huaweicloud_dataarts_architecture_subjects.filter_by_level"
		dcByLevel = acceptance.InitDataSourceCheck(byLevel)

		byWithRelation   = "data.huaweicloud_dataarts_architecture_subjects.filter_by_with_relation"
		dcByWithRelation = acceptance.InitDataSourceCheck(byWithRelation)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArchitectureSubjects_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "subjects.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(all, "subjects.0.id"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.name"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.name_en"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.description"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.qualified_name"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.code"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.status"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.data_owner_list"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.path"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.level"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.ordinal"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.owner"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.from_public"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.created_by"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.updated_by"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.created_at"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.updated_at"),
					resource.TestCheckResourceAttrSet(all, "subjects.0.children_num"),
					resource.TestMatchResourceAttr(all, "subjects.0.relations.#", regexp.MustCompile(`^[0-9]+$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					dcByCreateBy.CheckResourceExists(),
					resource.TestCheckOutput("is_create_by_filter_useful", "true"),
					dcByLevel.CheckResourceExists(),
					resource.TestCheckOutput("is_level_filter_useful", "true"),
					dcByWithRelation.CheckResourceExists(),
					resource.TestCheckOutput("is_with_relation_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceArchitectureSubjects_basic_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_subject" "test" {
  count = 3

  workspace_id = "%[1]s"
  name         = format("%[2]s-%%d", count.index)
  code         = format("%[2]s_code-%%d", count.index)
  owner        = "terraform_test"
  level        = 1
  description  = "Created by terraform script"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccDataSourceArchitectureSubjects_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all subjects without any filter.
data "huaweicloud_dataarts_architecture_subjects" "all" {
  workspace_id = "%[2]s"

  depends_on = [huaweicloud_dataarts_architecture_subject.test]
}

# Filter subjects by name.
locals {
  name = huaweicloud_dataarts_architecture_subject.test[0].name
}

data "huaweicloud_dataarts_architecture_subjects" "filter_by_name" {
  workspace_id = "%[2]s"
  name         = local.name

  depends_on = [huaweicloud_dataarts_architecture_subject.test]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_subjects.filter_by_name.subjects 
	: v.name == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter subjects by status.
locals {
  status = "DRAFT"
}

data "huaweicloud_dataarts_architecture_subjects" "filter_by_status" {
  workspace_id = "%[2]s"
  status       = local.status

  depends_on = [huaweicloud_dataarts_architecture_subject.test]
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_subjects.filter_by_status.subjects
    : v.status == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter subjects by create_by.
locals {
  owner = data.huaweicloud_dataarts_architecture_subjects.all.subjects[0].owner
}

data "huaweicloud_dataarts_architecture_subjects" "filter_by_create_by" {
  workspace_id = "%[2]s"
  create_by    = local.owner

  depends_on = [huaweicloud_dataarts_architecture_subject.test]
}

locals {
  create_by_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_subjects.filter_by_create_by.subjects
    : v.created_by == local.owner
  ]
}

output "is_create_by_filter_useful" {
  value = length(local.create_by_filter_result) > 0 && alltrue(local.create_by_filter_result)
}

# Filter subjects by level.
data "huaweicloud_dataarts_architecture_subjects" "filter_by_level" {
  workspace_id = "%[2]s"
  level        = 1

  depends_on = [huaweicloud_dataarts_architecture_subject.test]
}

locals {
  level_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_subjects.filter_by_level.subjects
    : v.level == 1
  ]
}

output "is_level_filter_useful" {
  value = length(local.level_filter_result) > 0 && alltrue(local.level_filter_result)
}

# Filter subjects by with_relation.
locals {
  with_relation = true
}

data "huaweicloud_dataarts_architecture_subjects" "filter_by_with_relation" {
  workspace_id   = "%[2]s"
  with_relation  = local.with_relation

  depends_on = [huaweicloud_dataarts_architecture_subject.test]
}

locals {
  with_relation_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_subjects.filter_by_with_relation.subjects
    : length(v.relations) >= 0
  ]
}

output "is_with_relation_filter_useful" {
  value = length(local.with_relation_filter_result) > 0 && alltrue(local.with_relation_filter_result)
}
`, testAccDataSourceArchitectureSubjects_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
