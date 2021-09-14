package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/addons"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCCEAddonV3_basic(t *testing.T) {
	var addon addons.Addon

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_addon.test"
	clusterName := "huaweicloud_cce_cluster.test"

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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCEAddonImportStateIdFunc(),
			},
		},
	})
}

func TestAccCCEAddonV3_values(t *testing.T) {
	var addon addons.Addon

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_addon.test"
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckProjectID(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3_values(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEAddonV3Exists(resourceName, clusterName, &addon),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
				),
			},
		},
	})
}

func testAccCheckCCEAddonV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	cceClient, err := config.CceAddonV3Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE Addon client: %s", err)
	}

	var clusterId string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "huaweicloud_cce_cluster" {
			clusterId = rs.Primary.ID
		}

		if rs.Type != "huaweicloud_cce_addon" {
			continue
		}

		if clusterId != "" {
			_, err := addons.Get(cceClient, rs.Primary.ID, clusterId).Extract()
			if err == nil {
				return fmtp.Errorf("addon still exists")
			}
		}
	}
	return nil
}

func testAccCheckCCEAddonV3Exists(n string, cluster string, addon *addons.Addon) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}
		c, ok := s.RootModule().Resources[cluster]
		if !ok {
			return fmtp.Errorf("Cluster not found: %s", c)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}
		if c.Primary.ID == "" {
			return fmtp.Errorf("Cluster id is not set")
		}

		config := testAccProvider.Meta().(*config.Config)
		cceClient, err := config.CceAddonV3Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud CCE Addon client: %s", err)
		}

		found, err := addons.Get(cceClient, rs.Primary.ID, c.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmtp.Errorf("Addon not found")
		}

		*addon = *found

		return nil
	}
}

func testAccCCEAddonImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterID string
		var addonID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_cce_cluster" {
				clusterID = rs.Primary.ID
			} else if rs.Type == "huaweicloud_cce_addon" {
				addonID = rs.Primary.ID
			}
		}
		if clusterID == "" || addonID == "" {
			return "", fmtp.Errorf("resource not found: %s/%s", clusterID, addonID)
		}
		return fmt.Sprintf("%s/%s", clusterID, addonID), nil
	}
}

func testAccCCEAddonV3_Base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = "%s"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  key_pair          = huaweicloud_compute_keypair.test.name

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
  cluster_id    = huaweicloud_cce_cluster.test.id
  version       = "1.1.0"
  template_name = "metrics-server"
  depends_on    = [huaweicloud_cce_node.test]
}
`, testAccCCEAddonV3_Base(rName))
}

func testAccCCEAddonV3_values(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_addon_template" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  name       = "autoscaler"
  version    = "1.19.6"
}

resource "huaweicloud_cce_addon" "test" {
  cluster_id    = huaweicloud_cce_cluster.test.id
  template_name = "autoscaler"
  version       = "1.19.6"

  values {
    basic  = jsondecode(data.huaweicloud_cce_addon_template.test.spec).basic
    custom = merge(
      jsondecode(data.huaweicloud_cce_addon_template.test.spec).parameters.custom,
      {
        cluster_id = huaweicloud_cce_cluster.test.id
        tenant_id  = "%s"
      }
    )
    flavor_json = jsonencode(jsondecode(data.huaweicloud_cce_addon_template.test.spec).parameters.flavor2)
  }
  
  depends_on = [huaweicloud_cce_node.test]
}
`, testAccCCEAddonV3_Base(rName), HW_PROJECT_ID)
}
