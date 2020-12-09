package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cce/v3/addons"
)

func TestAccCCEAddonV3_basic(t *testing.T) {
	var addon addons.Addon

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_addon.test"
	clusterName := "huaweicloud_cce_cluster_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEAddonV3Exists(resourceName, clusterName, &addon),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
				),
			},
		},
	})
}

func testAccCheckCCEAddonV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.CceAddonV3Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE Addon client: %s", err)
	}

	var clusterId string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "huaweicloud_cce_cluster_v3" {
			clusterId = rs.Primary.ID
		}

		if rs.Type != "huaweicloud_cce_addon" {
			continue
		}

		if clusterId != "" {
			_, err := addons.Get(cceClient, rs.Primary.ID, clusterId).Extract()
			if err == nil {
				return fmt.Errorf("addon still exists")
			}
		}
	}
	return nil
}

func testAccCheckCCEAddonV3Exists(n string, cluster string, addon *addons.Addon) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		c, ok := s.RootModule().Resources[cluster]
		if !ok {
			return fmt.Errorf("Cluster not found: %s", c)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		if c.Primary.ID == "" {
			return fmt.Errorf("Cluster id is not set")
		}

		config := testAccProvider.Meta().(*Config)
		cceClient, err := config.CceAddonV3Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud CCE Addon client: %s", err)
		}

		found, err := addons.Get(cceClient, rs.Primary.ID, c.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Addon not found")
		}

		*addon = *found

		return nil
	}
}

func testAccCCEAddonV3_Base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_v3" "test" {
	cluster_id        = huaweicloud_cce_cluster_v3.test.id
	name              = "%s"
	flavor_id         = "s6.large.2"
	availability_zone = data.huaweicloud_availability_zones.test.names[0]
	key_pair          = huaweicloud_compute_keypair_v2.test.name

	root_volume {
	size       = 40
	volumetype = "SSD"
	}
	data_volumes {
	size       = 100
	volumetype = "SSD"
	}
}
`, testAccCCENodeV3_Base(rName), rName)
}

func testAccCCEAddonV3_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_addon" "test" {
    cluster_id = huaweicloud_cce_cluster_v3.test.id
    version = "1.0.3"
	template_name = "metrics-server"
	depends_on = [huaweicloud_cce_node_v3.test]
}
`, testAccCCEAddonV3_Base(rName))
}
