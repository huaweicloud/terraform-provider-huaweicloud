package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBandWidthsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_vpc_bandwidths.all"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBandWidthsDataSource_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "bandwidths.#"),
				),
			},
		},
	})
}

func testAccBandWidthsDataSourceBase(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name                  = "%s"
  size                  = 5
  enterprise_project_id = "%s"
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type = "WHOLE"
    id         = huaweicloud_vpc_bandwidth.test.id
  }
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccBandWidthsDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_bandwidths" "all" {
  depends_on = [huaweicloud_vpc_eip.test]
}
`, testAccBandWidthsDataSourceBase(rName))
}

func TestAccBandWidthsDataSource_filter(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceName := "huaweicloud_vpc_bandwidth.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	eipResourceName := "huaweicloud_vpc_eip.test"

	byEpsId := "data.huaweicloud_vpc_bandwidths.filter_by_eps_id"
	byEpsIdNotFound := "data.huaweicloud_vpc_bandwidths.filter_by_eps_id_not_found"
	dcbyEpsId := acceptance.InitDataSourceCheck(byEpsId)
	dcbyEpsIdNotFound := acceptance.InitDataSourceCheck(byEpsIdNotFound)

	byName := "data.huaweicloud_vpc_bandwidths.filter_by_name"
	byNameNotFound := "data.huaweicloud_vpc_bandwidths.filter_by_name_not_found"
	dcbyName := acceptance.InitDataSourceCheck(byName)
	dcbyNameNotFound := acceptance.InitDataSourceCheck(byNameNotFound)

	bySize := "data.huaweicloud_vpc_bandwidths.filter_by_size"
	bySizeNotFound := "data.huaweicloud_vpc_bandwidths.filter_by_size_not_found"
	dcbySize := acceptance.InitDataSourceCheck(bySize)
	dcbySizeNotFound := acceptance.InitDataSourceCheck(bySizeNotFound)

	byChargeMode := "data.huaweicloud_vpc_bandwidths.filter_by_charge_mode"
	byChargeModeNotFound := "data.huaweicloud_vpc_bandwidths.filter_by_charge_mode_not_found"
	dcbyChargeMode := acceptance.InitDataSourceCheck(byChargeMode)
	dcbyChargeModeNotFound := acceptance.InitDataSourceCheck(byChargeModeNotFound)

	byID := "data.huaweicloud_vpc_bandwidths.filter_by_id"
	byIDNotFound := "data.huaweicloud_vpc_bandwidths.filter_by_id_not_found"
	dcbyID := acceptance.InitDataSourceCheck(byID)
	dcbyIDNotFound := acceptance.InitDataSourceCheck(byIDNotFound)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataArtsStudioWorkspaces_filter(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					dcbyEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					dcbyEpsIdNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful_not_found", "true"),

					dcbyName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcbyNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful_not_found", "true"),

					dcbySize.CheckResourceExists(),
					resource.TestCheckOutput("is_size_filter_useful", "true"),
					dcbySizeNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_size_filter_useful_not_found", "true"),

					dcbyChargeMode.CheckResourceExists(),
					resource.TestCheckOutput("is_charge_mode_filter_useful", "true"),
					dcbyChargeModeNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_charge_mode_filter_useful_not_found", "true"),

					dcbyID.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					dcbyIDNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful_not_found", "true"),
					resource.TestCheckResourceAttrPair(byID, "bandwidths.0.publicips.0.id",
						eipResourceName, "id"),
					resource.TestCheckResourceAttrPair(byID, "bandwidths.0.publicips.0.ip_address",
						eipResourceName, "address"),
				),
			},
		},
	})
}

func testAccDataArtsStudioWorkspaces_filter(rName string) string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
%[1]s

// filter_by_eps_id
data "huaweicloud_vpc_bandwidths" "filter_by_eps_id" {
  depends_on = [huaweicloud_vpc_eip.test]

  enterprise_project_id = "%[2]s"
}

data "huaweicloud_vpc_bandwidths" "filter_by_eps_id_not_found" {
  depends_on = [huaweicloud_vpc_eip.test]

  enterprise_project_id = "%[3]s"
}

locals {
  filter_result_by_eps_id = [for v in data.huaweicloud_vpc_bandwidths.filter_by_eps_id.bandwidths[*].enterprise_project_id : 
    v == huaweicloud_vpc_bandwidth.test.enterprise_project_id]
}

output "is_eps_id_filter_useful" {
  value = length(local.filter_result_by_eps_id) > 0 && alltrue(local.filter_result_by_eps_id) 
}

output "is_eps_id_filter_useful_not_found" {
  value = length(data.huaweicloud_vpc_bandwidths.filter_by_eps_id_not_found.bandwidths) == 0
}

// filter_by_name
data "huaweicloud_vpc_bandwidths" "filter_by_name" {
  depends_on = [huaweicloud_vpc_eip.test]

  name = "%[4]s"
}

data "huaweicloud_vpc_bandwidths" "filter_by_name_not_found" {
  depends_on = [huaweicloud_vpc_eip.test]

  name = "test-not-found-%[4]s"
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_vpc_bandwidths.filter_by_name.bandwidths[*].name :
    v == huaweicloud_vpc_bandwidth.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name)
}

output "is_name_filter_useful_not_found" {
  value = length(data.huaweicloud_vpc_bandwidths.filter_by_name_not_found.bandwidths) == 0
}

// filter_by_size
data "huaweicloud_vpc_bandwidths" "filter_by_size" {
  depends_on = [huaweicloud_vpc_eip.test]

  size = huaweicloud_vpc_bandwidth.test.size
}

data "huaweicloud_vpc_bandwidths" "filter_by_size_not_found" {
  depends_on = [huaweicloud_vpc_eip.test]

  size = -1
}

locals {
  filter_result_by_size = [for v in data.huaweicloud_vpc_bandwidths.filter_by_size.bandwidths[*].size :
    v == huaweicloud_vpc_bandwidth.test.size]
}

output "is_size_filter_useful" {
  value = length(local.filter_result_by_size) > 0 && alltrue(local.filter_result_by_size)
}

output "is_size_filter_useful_not_found" {
  value = length(data.huaweicloud_vpc_bandwidths.filter_by_size_not_found.bandwidths) == 0
}

// filter_by_charge_mode
data "huaweicloud_vpc_bandwidths" "filter_by_charge_mode" {
  depends_on = [huaweicloud_vpc_eip.test]

  charge_mode = huaweicloud_vpc_bandwidth.test.charge_mode
}

data "huaweicloud_vpc_bandwidths" "filter_by_charge_mode_not_found" {
  depends_on = [huaweicloud_vpc_eip.test]

  charge_mode = "95peak_plus"
}

locals {
  filter_result_by_charge_mode = [for v in data.huaweicloud_vpc_bandwidths.filter_by_charge_mode.bandwidths[*].charge_mode :
    v == huaweicloud_vpc_bandwidth.test.charge_mode]
}

output "is_charge_mode_filter_useful" {
  value = length(local.filter_result_by_charge_mode) > 0 && alltrue(local.filter_result_by_charge_mode)
}

output "is_charge_mode_filter_useful_not_found" {
  value = length(data.huaweicloud_vpc_bandwidths.filter_by_charge_mode_not_found.bandwidths) == 0
}

// filter_by_id
data "huaweicloud_vpc_bandwidths" "filter_by_id" {
  depends_on = [huaweicloud_vpc_eip.test]

  bandwidth_id = huaweicloud_vpc_bandwidth.test.id
}

data "huaweicloud_vpc_bandwidths" "filter_by_id_not_found" {
  depends_on = [huaweicloud_vpc_eip.test]

  bandwidth_id = "%[3]s"
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_vpc_bandwidths.filter_by_id.bandwidths) == 1
}

output "is_id_filter_useful_not_found" {
  value = length(data.huaweicloud_vpc_bandwidths.filter_by_id_not_found.bandwidths) == 0
}
`, testAccBandWidthsDataSourceBase(rName), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, randUUID, rName)
}
