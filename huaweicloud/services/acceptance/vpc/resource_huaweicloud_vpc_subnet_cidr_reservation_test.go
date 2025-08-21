package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getReservationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v3 client: %s", err)
	}

	httpUrl := fmt.Sprintf("vpc/virsubnet-cidr-reservations/%s", state.Primary.ID)
	getPath := client.ResourceBaseURL() + httpUrl

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func TestAccVpcSubnetCidrReservation_basic(t *testing.T) {
	var reservation interface{}

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_subnet_cidr_reservation.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&reservation,
		getReservationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetCidrReservation_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "mask", "26"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testAccVpcSubnetCidrReservation_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVpcSubnetCidrReservation_withCidr(t *testing.T) {
	var reservation interface{}

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_subnet_cidr_reservation.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&reservation,
		getReservationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetCidrReservation_withCidr(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.64/26"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "mask"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVpcSubnetCidrReservation_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}
`, name)
}

func testAccVpcSubnetCidrReservation_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_subnet_cidr_reservation" "test" {
  subnet_id   = huaweicloud_vpc_subnet.test.id
  ip_version  = 4
  mask        = 26
  name        = "%[2]s"
  description = "created by terraform"
}
`, testAccVpcSubnetCidrReservation_base(name), name)
}

func testAccVpcSubnetCidrReservation_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_subnet_cidr_reservation" "test" {
  subnet_id   = huaweicloud_vpc_subnet.test.id
  ip_version  = 4
  mask        = 26
  name        = "%[2]s-update"
  description = "updated by terraform"
}
`, testAccVpcSubnetCidrReservation_base(name), name)
}

func testAccVpcSubnetCidrReservation_withCidr(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_subnet_cidr_reservation" "test" {
  subnet_id   = huaweicloud_vpc_subnet.test.id
  ip_version  = 4
  cidr        = "192.168.0.64/26"
  name        = "%[2]s"
  description = "created by terraform"
}
`, testAccVpcSubnetCidrReservation_base(name), name)
}
