package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getAppCustomStoragePolicyFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}
	return workspace.GetAppCustomStoragePolicy(client)
}

func TestAccAppStoragePolicy_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_storage_policy.test"

		policy interface{}
		rc     = acceptance.InitResourceCheck(resourceName, &policy, getAppCustomStoragePolicyFunc)
	)

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppStoragePolicy_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "server_actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server_actions.0", "GetObject"),
					resource.TestCheckResourceAttr(resourceName, "client_actions.#", "0"),
				),
			},
			{
				Config: testAccAppStoragePolicy_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "server_actions.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "client_actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_actions.0", "GetObject"),
				),
			},
			{
				Config: testAccAppStoragePolicy_basic_step3,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "server_actions.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "client_actions.#", "2"),
				),
			},
			{
				Config: testAccAppStoragePolicy_basic_step4,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "server_actions.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "client_actions.#", "3"),
				),
			},
			// Import by ID.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Import by other characters.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppCustomImportStateIdFunc(resourceName, "NA"),
			},
		},
	})
}

func testAccAppCustomImportStateIdFunc(resourceName, customId string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) is not found in the tfstate", resourceName)
		}
		return customId, nil
	}
}

const testAccAppStoragePolicy_basic_step1 string = `
resource "huaweicloud_workspace_app_storage_policy" "test" {
  server_actions = ["GetObject"]
}
`

const testAccAppStoragePolicy_basic_step2 string = `
resource "huaweicloud_workspace_app_storage_policy" "test" {
  server_actions = ["GetObject", "PutObject", "DeleteObject"]
  client_actions = ["GetObject"]
}
`

const testAccAppStoragePolicy_basic_step3 string = `
resource "huaweicloud_workspace_app_storage_policy" "test" {
  server_actions = ["GetObject", "PutObject", "DeleteObject"]
  client_actions = ["PutObject", "DeleteObject"]
}
`

const testAccAppStoragePolicy_basic_step4 string = `
resource "huaweicloud_workspace_app_storage_policy" "test" {
  server_actions = ["GetObject", "PutObject", "DeleteObject"]
  client_actions = ["GetObject", "PutObject", "DeleteObject"]
}
`
