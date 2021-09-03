package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccVpcV1_basic(t *testing.T) {
	var vpc vpcs.Vpc

	rName := acceptance.RandomAccResourceName()
	rNameUpdate := rName + "_updated"
	resourceName := "huaweicloud_vpc.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "status", "OK"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcV1_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_updated"),
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

func TestAccVpcV1_WithEpsId(t *testing.T) {
	var vpc vpcs.Vpc

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_epsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "status", "OK"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

// TestAccVpcV1_WithCustomRegion this case will run a test for resource-level region. Before run this case,
// you shoule set `HW_CUSTOM_REGION_NAME` in your system and it should be different from `HW_REGION_NAME`.
func TestAccVpcV1_WithCustomRegion(t *testing.T) {

	vpcName1 := fmt.Sprintf("test_vpc_region_%s", acctest.RandString(5))
	vpcName2 := fmt.Sprintf("test_vpc_region_%s", acctest.RandString(5))

	resName1 := "huaweicloud_vpc.test1"
	resName2 := "huaweicloud_vpc.test2"

	var vpc1, vpc2 vpcs.Vpc

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPrecheckCustomRegion(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_WithCustomRegion(vpcName1, vpcName2, acceptance.HW_CUSTOM_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomRegionVpcV1Exists(resName1, &vpc1, acceptance.HW_REGION_NAME),
					testAccCheckCustomRegionVpcV1Exists(resName2, &vpc2, acceptance.HW_CUSTOM_REGION_NAME),
				),
			},
		},
	})
}

func testAccCheckVpcV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	vpcClient, err := config.NetworkingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating huaweicloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc" {
			continue
		}

		_, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Vpc still exists")
		}
	}

	return nil
}

func testAccCheckCustomRegionVpcV1Exists(name string, vpc *vpcs.Vpc, region string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmtp.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		vpcClient, err := config.NetworkingV1Client(region)
		if err != nil {
			return fmtp.Errorf("Error creating huaweicloud vpc client: %s", err)
		}

		found, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("vpc not found")
		}

		*vpc = *found
		return nil
	}
}

func testAccCheckVpcV1Exists(n string, vpc *vpcs.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		vpcClient, err := config.NetworkingV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating huaweicloud vpc client: %s", err)
		}

		found, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("vpc not found")
		}

		*vpc = *found

		return nil
	}
}

func testAccVpcV1_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccVpcV1_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr ="192.168.0.0/16"

  tags = {
    foo1 = "bar"
    key  = "value_updated"
  }
}
`, rName)
}

func testAccVpcV1_epsId(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name                  = "%s"
  cidr                  = "192.168.0.0/16"
  enterprise_project_id = "%s"
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID)
}

func testAccVpcV1_WithCustomRegion(name1, name2, region string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test1" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc" "test2" {   
  region = "%s"
  name   = "%s"
  cidr   = "192.168.0.0/16"
}
`, name1, region, name2)
}
