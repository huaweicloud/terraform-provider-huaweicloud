package swr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/swr/v2/domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceRepositorySharing(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	swrClient, err := conf.SwrV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	return domains.Get(swrClient, state.Primary.Attributes["organization"],
		strings.ReplaceAll(state.Primary.Attributes["repository"], "/", "$"), state.Primary.ID).Extract()
}

func TestAccSWRRepositorySharing_basic(t *testing.T) {
	var domain domains.AccessDomain
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_swr_repository_sharing.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&domain,
		getResourceRepositorySharing,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSWRDomian(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSWRRepositorySharing_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${huaweicloud_swr_organization.test.name}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "repository",
						"${huaweicloud_swr_repository.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "sharing_account", acceptance.HW_SWR_SHARING_ACCOUNT),
					resource.TestCheckResourceAttr(resourceName, "deadline", "forever"),
					resource.TestCheckResourceAttr(resourceName, "permission", "pull"),
				),
			},
			{
				Config: testAccSWRRepositorySharing_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${huaweicloud_swr_organization.test.name}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "repository",
						"${huaweicloud_swr_repository.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "sharing_account", acceptance.HW_SWR_SHARING_ACCOUNT),
					resource.TestCheckResourceAttr(resourceName, "deadline", "2099-12-31"),
					resource.TestCheckResourceAttr(resourceName, "permission", "pull"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSWRRepositorySharingImportStateIdFunc(),
			},
		},
	})
}

func testAccSWRRepositorySharingImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var organization string
		var repositoryID string
		var sharingAccount string
		for _, rs := range s.RootModule().Resources {
			switch rs.Type {
			case "huaweicloud_swr_organization":
				organization = rs.Primary.Attributes["name"]
			case "huaweicloud_swr_repository":
				repositoryID = rs.Primary.ID
			case "huaweicloud_swr_repository_sharing":
				sharingAccount = rs.Primary.ID
			}
		}
		if organization == "" || repositoryID == "" || sharingAccount == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", organization, repositoryID, sharingAccount)
		}
		return fmt.Sprintf("%s/%s/%s", organization, repositoryID, sharingAccount), nil
	}
}

func testAccSWRRepositorySharing_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_repository_sharing" "test" {
  organization    = huaweicloud_swr_organization.test.name
  repository      = huaweicloud_swr_repository.test.name
  sharing_account = "%s"
  permission      = "pull"
  deadline        = "forever"
}
`, testAccSWRRepository_basic(rName), acceptance.HW_SWR_SHARING_ACCOUNT)
}

func testAccSWRRepositorySharing_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_repository_sharing" "test" {
  organization    = huaweicloud_swr_organization.test.name
  repository      = huaweicloud_swr_repository.test.name
  sharing_account = "%s"
  permission      = "pull"
  deadline        = "2099-12-31"
}
`, testAccSWRRepository_basic(rName), acceptance.HW_SWR_SHARING_ACCOUNT)
}
