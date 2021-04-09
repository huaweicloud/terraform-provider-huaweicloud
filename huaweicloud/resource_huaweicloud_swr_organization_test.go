package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/swr/v2/namespaces"
)

func TestAccSWROrganization_basic(t *testing.T) {
	var org namespaces.Namespace

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_swr_organization.test"
	loginServer := fmt.Sprintf("swr.%s.myhuaweicloud.com", HW_REGION_NAME)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSWROrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSWROrganization_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSWROrganizationExists(resourceName, &org),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "permission", "Manage"),
					resource.TestCheckResourceAttr(resourceName, "login_server", loginServer),
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

func testAccCheckSWROrganizationDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	swrClient, err := config.SwrV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SWR client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_swr_organization" {
			continue
		}

		_, err := namespaces.Get(swrClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SWR organization still exists")
		}
	}

	return nil
}

func testAccCheckSWROrganizationExists(n string, org *namespaces.Namespace) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		swrClient, err := config.SwrV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud SWR client: %s", err)
		}

		found, err := namespaces.Get(swrClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("SWR organization not found")
		}

		*org = *found

		return nil
	}
}

func testAccSWROrganization_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_swr_organization" "test" {
  name = "%s"
}
`, rName)
}
