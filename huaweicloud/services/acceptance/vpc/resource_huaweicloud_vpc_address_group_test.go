package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	vpc_model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getVpcAddressGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcVpcV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Huaweicloud VPC client: %s", err)
	}

	request := &vpc_model.ShowAddressGroupRequest{
		AddressGroupId: state.Primary.ID,
	}

	return client.ShowAddressGroup(request)
}

func TestAccVpcAddressGroup_basic(t *testing.T) {
	var group vpc_model.ShowAddressGroupResponse

	rName := acceptance.RandomAccResourceName()
	rNameUpdate := rName + "_updated"
	resourceName := "huaweicloud_vpc_address_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getVpcAddressGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testVpcAdressGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_capacity", "20"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testVpcAdressGroup_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by acc test"),
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "max_capacity", "10"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccVpcAddressGroup_ipExtraSet(t *testing.T) {
	var group vpc_model.ShowAddressGroupResponse

	rName := acceptance.RandomAccResourceName()
	rNameUpdate := rName + "_updated"
	resourceName := "huaweicloud_vpc_address_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getVpcAddressGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testVpcAdressGroup_ipExtraSet(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "ip_extra_set.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_capacity", "20"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testVpcAdressGroup_ipExtraSetUpdate(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by acc test"),
					resource.TestCheckResourceAttr(resourceName, "ip_extra_set.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "max_capacity", "10"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccVpcAddressGroup_ipv6(t *testing.T) {
	var group vpc_model.ShowAddressGroupResponse

	rName := acceptance.RandomAccResourceName()
	rNameUpdate := rName + "_updated"
	resourceName := "huaweicloud_vpc_address_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getVpcAddressGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testVpcAdressGroup_ipv6(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "6"),
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_capacity", "20"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testVpcAdressGroup_ipv6_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by acc test"),
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_capacity", "10"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccVpcAddressGroup_eps(t *testing.T) {
	var group vpc_model.ShowAddressGroupResponse

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_address_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getVpcAddressGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testVpcAdressGroup_eps(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_capacity", "20"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func testVpcAdressGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_address_group" "test" {
  name        = "%s"
  description = "created by acc test"
  addresses   = [
    "192.168.3.2",
    "192.168.3.20-192.168.3.100"
  ]
}
`, rName)
}

func testVpcAdressGroup_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_address_group" "test" {
  name        = "%s"
  description = "updated by acc test"
  addresses = [
    "192.168.5.0/24",
    "192.168.3.2",
    "192.168.3.20-192.168.3.100"
  ]
  max_capacity = 10
}
`, rName)
}

func testVpcAdressGroup_ipExtraSet(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_address_group" "test" {
  name        = "%s"
  description = "created by acc test"

  ip_extra_set {
    ip      = "192.168.3.2"
    remarks = "terraform test 1"
  }

  ip_extra_set {
    ip      = "192.168.3.20-192.168.3.100"
    remarks = "terraform test 2"
  }
}
`, rName)
}

func testVpcAdressGroup_ipExtraSetUpdate(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_address_group" "test" {
  name        = "%s"
  description = "updated by acc test"

  ip_extra_set {
    ip      = "192.168.3.2"
    remarks = "terraform test 1"
  }

  ip_extra_set {
    ip      = "192.168.5.0/24"
    remarks = "terraform test 2"
  }

  ip_extra_set {
    ip = "192.168.3.20-192.168.3.100"
  }

  max_capacity = 10
}
`, rName)
}

func testVpcAdressGroup_ipv6(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_address_group" "test" {
  name        = "%s"
  description = "created by acc test"
  ip_version  = 6
  addresses   = [
    "2001:db8:a583:6e::/64"
  ]
}
`, rName)
}

func testVpcAdressGroup_ipv6_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_address_group" "test" {
  name        = "%s"
  description = "updated by acc test"
  ip_version  = 6
  addresses = [
    "2001:db8:a583:8e::1-2001:db8:a583:8e::50",
    "2001:db8:a583:6e::/64"
  ]
  max_capacity = 10
}
`, rName)
}

func testVpcAdressGroup_eps(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_address_group" "test" {
  name        = "%s"
  description = "created by acc test"
  addresses = [
    "192.168.3.2",
    "192.168.3.20-192.168.3.100"
  ]
  enterprise_project_id = "%s"
  force_destroy         = true
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
