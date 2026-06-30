package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataArchitectureDataStandards_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_architecture_data_standards.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byDirectoryId   = "data.huaweicloud_dataarts_architecture_data_standards.filter_by_directory_id"
		dcByDirectoryId = acceptance.InitDataSourceCheck(byDirectoryId)

		byNameCh   = "data.huaweicloud_dataarts_architecture_data_standards.filter_by_name_ch"
		dcByNameCh = acceptance.InitDataSourceCheck(byNameCh)

		byNameEn   = "data.huaweicloud_dataarts_architecture_data_standards.filter_by_name_en"
		dcByNameEn = acceptance.InitDataSourceCheck(byNameEn)

		byBeginTime   = "data.huaweicloud_dataarts_architecture_data_standards.filter_by_begin_time"
		dcByBeginTime = acceptance.InitDataSourceCheck(byBeginTime)

		byEndTime   = "data.huaweicloud_dataarts_architecture_data_standards.filter_by_end_time"
		dcByEndTime = acceptance.InitDataSourceCheck(byEndTime)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataArchitectureDataStandards_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts Architecture data standards"),
			},
			{
				Config: testAccDataArchitectureDataStandards_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "data_standards.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.id"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.directory_id"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.status"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.created_by"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.updated_by"),
					resource.TestMatchResourceAttr(all, "data_standards.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "data_standards.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "data_standards.0.values.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.values.0.fd_name"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.values.0.fd_value"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.values.0.id"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.values.0.fd_id"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.values.0.directory_id"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.values.0.row_id"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.values.0.status"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.values.0.created_by"),
					resource.TestCheckResourceAttrSet(all, "data_standards.0.values.0.updated_by"),
					resource.TestMatchResourceAttr(all, "data_standards.0.values.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "data_standards.0.values.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					// Filter by directory ID.
					dcByDirectoryId.CheckResourceExists(),
					resource.TestCheckOutput("is_directory_id_filter_useful", "true"),

					// Filter by Chinese name.
					dcByNameCh.CheckResourceExists(),
					resource.TestCheckOutput("is_name_ch_filter_useful", "true"),

					// Filter by English code.
					dcByNameEn.CheckResourceExists(),
					resource.TestCheckOutput("is_name_en_filter_useful", "true"),

					// Filter by begin time.
					dcByBeginTime.CheckResourceExists(),
					resource.TestCheckOutput("is_begin_time_filter_useful", "true"),

					// Filter by end time.
					dcByEndTime.CheckResourceExists(),
					resource.TestCheckOutput("is_end_time_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataArchitectureDataStandards_nonExistentWorkspace() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_data_standards" "test" {
  workspace_id = "%[1]s"
}
`, randUUID.String())
}

func testAccDataArchitectureDataStandards_basic_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_directory" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  type         = "STANDARD_ELEMENT"
}

resource "huaweicloud_dataarts_architecture_data_standard" "test" {
  workspace_id = "%[1]s"
  directory_id = huaweicloud_dataarts_architecture_directory.test.id

  values {
    fd_name  = "nameCh"
    fd_value = "%[2]s"
  }

  values {
    fd_name  = "nameEn"
    fd_value = "%[2]s_en"
  }

  values {
    fd_name  = "description"
    fd_value = "Created by acceptance test"
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccDataArchitectureDataStandards_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dataarts_architecture_data_standards" "all" {
  depends_on = [huaweicloud_dataarts_architecture_data_standard.test]

  workspace_id = "%[2]s"
}

# Filter by directory ID.
locals {
  data_standard_directory_id = huaweicloud_dataarts_architecture_data_standard.test.directory_id
}

data "huaweicloud_dataarts_architecture_data_standards" "filter_by_directory_id" {
  depends_on = [huaweicloud_dataarts_architecture_data_standard.test]

  workspace_id = "%[2]s"
  directory_id = local.data_standard_directory_id
}

locals {
  directory_id_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_data_standards.filter_by_directory_id.data_standards :
    v.directory_id == local.data_standard_directory_id
  ]
}

output "is_directory_id_filter_useful" {
  value = length(local.directory_id_filter_result) > 0 && alltrue(local.directory_id_filter_result)
}

# Filter by Chinese name (exact match).
locals {
  data_standard_name_ch = "%[3]s"
}

data "huaweicloud_dataarts_architecture_data_standards" "filter_by_name_ch" {
  depends_on = [huaweicloud_dataarts_architecture_data_standard.test]

  workspace_id = "%[2]s"
  name_ch      = local.data_standard_name_ch
}

output "is_name_ch_filter_useful" {
  value = contains(data.huaweicloud_dataarts_architecture_data_standards.filter_by_name_ch.data_standards[*].id,
    huaweicloud_dataarts_architecture_data_standard.test.id)
}

# Filter by English code (exact match).
locals {
  data_standard_name_en = "%[3]s_en"
}

data "huaweicloud_dataarts_architecture_data_standards" "filter_by_name_en" {
  depends_on = [huaweicloud_dataarts_architecture_data_standard.test]

  workspace_id = "%[2]s"
  name_en      = local.data_standard_name_en
}

output "is_name_en_filter_useful" {
  value = contains(data.huaweicloud_dataarts_architecture_data_standards.filter_by_name_en.data_standards[*].id,
    huaweicloud_dataarts_architecture_data_standard.test.id)
}

# Filter by begin time.
locals {
  data_standard_begin_time = timeadd(huaweicloud_dataarts_architecture_directory.test.created_at, "-10m")
}

data "huaweicloud_dataarts_architecture_data_standards" "filter_by_begin_time" {
  depends_on = [huaweicloud_dataarts_architecture_data_standard.test]

  workspace_id = "%[2]s"
  begin_time   = local.data_standard_begin_time
}

output "is_begin_time_filter_useful" {
  value = contains(data.huaweicloud_dataarts_architecture_data_standards.filter_by_begin_time.data_standards[*].id,
    huaweicloud_dataarts_architecture_data_standard.test.id)
}

# Filter by end time.
locals {
  data_standard_end_time = timeadd(huaweicloud_dataarts_architecture_directory.test.created_at, "10m")
}

data "huaweicloud_dataarts_architecture_data_standards" "filter_by_end_time" {
  depends_on = [huaweicloud_dataarts_architecture_data_standard.test]

  workspace_id = "%[2]s"
  end_time     = local.data_standard_end_time
}

output "is_end_time_filter_useful" {
  value = contains(data.huaweicloud_dataarts_architecture_data_standards.filter_by_end_time.data_standards[*].id,
    huaweicloud_dataarts_architecture_data_standard.test.id)
}
`, testAccDataArchitectureDataStandards_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
