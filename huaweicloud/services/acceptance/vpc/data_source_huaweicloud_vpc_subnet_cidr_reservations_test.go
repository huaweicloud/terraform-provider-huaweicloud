package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcSubnetCidrReservations_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_vpc_subnet_cidr_reservations.basic"
	dataSource2 := "data.huaweicloud_vpc_subnet_cidr_reservations.filter_by_id"
	dataSource3 := "data.huaweicloud_vpc_subnet_cidr_reservations.filter_by_subnet_id"
	dataSource4 := "data.huaweicloud_vpc_subnet_cidr_reservations.filter_by_cidr"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcSubnetCidrReservations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_subnet_id_filter_useful", "true"),
					resource.TestCheckOutput("is_cidr_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcSubnetCidrReservations_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnet_cidr_reservations" "basic" {
  depends_on = [huaweicloud_vpc_subnet_cidr_reservation.test]
}

data "huaweicloud_vpc_subnet_cidr_reservations" "filter_by_id" {
  reservation_id = [huaweicloud_vpc_subnet_cidr_reservation.test.id]

  depends_on = [huaweicloud_vpc_subnet_cidr_reservation.test]
}

data "huaweicloud_vpc_subnet_cidr_reservations" "filter_by_subnet_id" {
  subnet_id = [huaweicloud_vpc_subnet.test.id]

  depends_on = [huaweicloud_vpc_subnet_cidr_reservation.test]
}

data "huaweicloud_vpc_subnet_cidr_reservations" "filter_by_cidr" {
  cidr = ["192.168.0.64/26"]

  depends_on = [huaweicloud_vpc_subnet_cidr_reservation.test]
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_vpc_subnet_cidr_reservations.filter_by_id.reservations[*].id :
    v == huaweicloud_vpc_subnet_cidr_reservation.test.id
  ]
  subnet_id_filter_result = [
    for v in data.huaweicloud_vpc_subnet_cidr_reservations.filter_by_subnet_id.reservations[*].subnet_id :
    v == huaweicloud_vpc_subnet_cidr_reservation.test.subnet_id
  ]
  cidr_filter_result = [
    for v in data.huaweicloud_vpc_subnet_cidr_reservations.filter_by_cidr.reservations[*].cidr :
    v == "192.168.0.64/26"
  ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_vpc_subnet_cidr_reservations.basic.reservations) > 0
}

output "is_id_filter_useful" {
  value = alltrue(local.id_filter_result) && length(local.id_filter_result) > 0
}

output "is_subnet_id_filter_useful" {
  value = alltrue(local.subnet_id_filter_result) && length(local.subnet_id_filter_result) > 0
}

output "is_cidr_filter_useful" {
  value = alltrue(local.cidr_filter_result) && length(local.cidr_filter_result) > 0
}
`, testAccVpcSubnetCidrReservation_withCidr(name))
}
