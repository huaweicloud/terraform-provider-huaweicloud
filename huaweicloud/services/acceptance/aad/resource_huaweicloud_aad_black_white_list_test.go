package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aad"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getBlackWhiteListFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "aad"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating AAD client: %s", err)
	}

	return aad.ReadBlackWhiteList(client, state.Primary.Attributes["type"], state.Primary.Attributes["instance_id"])
}

func TestAccBlackWhiteList_whitelist(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_aad_black_white_list.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getBlackWhiteListFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a valid AAD instance ID and config it to the environment variable.
			acceptance.TestAccPreCheckAadInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBlackWhiteList_whitelist(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_AAD_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "type", "white"),
					resource.TestCheckResourceAttr(resourceName, "ips.#", "2"),
				),
			},
			{
				Config: testAccBlackWhiteList_whitelist_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_AAD_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "type", "white"),
					resource.TestCheckResourceAttr(resourceName, "ips.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ips.0", "11.1.2.116"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testBlackWhiteListImportState(resourceName),
			},
		},
	})
}

func testAccBlackWhiteList_whitelist() string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_black_white_list" "test" {
  instance_id = "%s"
  type        = "white"
  ips         = ["11.1.2.114", "11.1.2.115"]
}
`, acceptance.HW_AAD_INSTANCE_ID)
}

func testAccBlackWhiteList_whitelist_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_black_white_list" "test" {
  instance_id = "%s"
  type        = "white"
  ips         = ["11.1.2.116"]
}
`, acceptance.HW_AAD_INSTANCE_ID)
}

func TestAccBlackWhiteList_blacklist(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_aad_black_white_list.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getBlackWhiteListFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a valid AAD instance ID and config it to the environment variable.
			acceptance.TestAccPreCheckAadInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBlackWhiteList_blacklist(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_AAD_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "type", "black"),
					resource.TestCheckResourceAttr(resourceName, "ips.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ips.0", "12.1.2.117"),
				),
			},
			{
				Config: testAccBlackWhiteList_blacklist_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_AAD_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "type", "black"),
					resource.TestCheckResourceAttr(resourceName, "ips.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testBlackWhiteListImportState(resourceName),
			},
		},
	})
}

func testAccBlackWhiteList_blacklist() string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_black_white_list" "test" {
  instance_id = "%s"
  type        = "black"
  ips         = ["12.1.2.117"]
}
`, acceptance.HW_AAD_INSTANCE_ID)
}

func testAccBlackWhiteList_blacklist_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_black_white_list" "test" {
  instance_id = "%s"
  type        = "black"
  ips         = ["12.1.2.118", "12.1.2.119"]
}
`, acceptance.HW_AAD_INSTANCE_ID)
}

func testBlackWhiteListImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		if instanceID == "" {
			return "", fmt.Errorf("attribute (instance_id) of resource (%s) not found", name)
		}

		typeParam := rs.Primary.Attributes["type"]
		if typeParam == "" {
			return "", fmt.Errorf("attribute (type) of resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", instanceID, typeParam), nil
	}
}
