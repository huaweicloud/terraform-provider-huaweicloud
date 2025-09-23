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

func getResourceAppImageServerFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}

	return workspace.GetAppImageServerById(client, state.Primary.ID)
}

func getAcceptanceEpsId() string {
	if acceptance.HW_ENTERPRISE_PROJECT_ID_TEST == "" {
		return "0"
	}
	return acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
}

func TestAccResourceAppImageServer_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		resourceName = "huaweicloud_workspace_app_image_server.test"
		serverGroup  interface{}
		rc           = acceptance.InitResourceCheck(
			resourceName,
			&serverGroup,
			getResourceAppImageServerFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
			acceptance.TestAccPreCheckWorkspaceAppImageSpecCode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceAppImageServer_basic(name, name, "Created by script"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"data.huaweicloud_workspace_service.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"data.huaweicloud_workspace_service.test", "network_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "image_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID),
					resource.TestCheckResourceAttr(resourceName, "image_type", "gold"),
					resource.TestCheckResourceAttr(resourceName, "image_source_product_id",
						acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID),
					resource.TestCheckResourceAttr(resourceName, "spec_code",
						acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_SPEC_CODE),
					resource.TestCheckResourceAttr(resourceName, "authorize_accounts.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "authorize_accounts.0.account",
						"huaweicloud_workspace_user.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "authorize_accounts.0.type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "is_vdi", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", getAcceptanceEpsId()),
					resource.TestCheckResourceAttr(resourceName, "is_delete_associated_resources", "true"),
				),
			},
			{
				Config: testResourceAppImageServer_basic(name, updateName, ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"flavor_id",
					"vpc_id",
					"subnet_id",
					"root_volume",
					"image_source_product_id",
					"is_vdi",
					"availability_zone",
					"scheduler_hints",
					"tags",
					"is_delete_associated_resources",
				},
			},
		},
	})
}

func testResourceAppImageServer_basic(name, updateName string, description string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_user" "test" {
  name                   = "%[1]s"
  email                  = "%[1]s@example.com"
  password_never_expires = false
  disabled               = false
}

resource "huaweicloud_workspace_app_image_server" "test" {
  name                           = "%[2]s"
  flavor_id                      = "%[3]s"
  vpc_id                         = try(data.huaweicloud_workspace_service.test.vpc_id, "")
  subnet_id                      = try(data.huaweicloud_workspace_service.test.network_ids[0], "")
  image_id                       = "%[4]s"
  image_type                     = "gold"
  image_source_product_id        = "%[5]s"
  spec_code                      = "%[6]s"
  is_vdi                         = true
  availability_zone              = data.huaweicloud_availability_zones.test.names[0]
  description                    = "%[7]s"
  enterprise_project_id          = "%[8]s"
  is_delete_associated_resources = true

  # Currently only one user can be set.
  authorize_accounts {
    account = huaweicloud_workspace_user.test.name
    type    = "USER"
  }

  root_volume {
    type = "SAS"
    size = 80
  }

  tags = {
    foo = "bar"
  }
}
`, name,
		updateName,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_SPEC_CODE,
		description,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// Before running this test, please enable a service that connects to LocalAD and the corresponding OU is created.
func TestAccResourceAppImageServer_withAD(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_image_server.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()

		serverGroup interface{}
		rc          = acceptance.InitResourceCheck(
			resourceName,
			&serverGroup,
			getResourceAppImageServerFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
			acceptance.TestAccPreCheckWorkspaceAppImageSpecCode(t)
			acceptance.TestAccPreCheckWorkspaceOUName(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceAppImageServer_withAD(name, "Created by script"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"data.huaweicloud_workspace_service.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"data.huaweicloud_workspace_service.test", "network_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "image_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID),
					resource.TestCheckResourceAttr(resourceName, "image_type", "gold"),
					resource.TestCheckResourceAttr(resourceName, "image_source_product_id",
						acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID),
					resource.TestCheckResourceAttr(resourceName, "spec_code",
						acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_SPEC_CODE),
					resource.TestCheckResourceAttr(resourceName, "authorize_accounts.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "authorize_accounts.0.account",
						"data.huaweicloud_workspace_service.test", "ad_domain.0.admin_account"),
					resource.TestCheckResourceAttr(resourceName, "authorize_accounts.0.type", "USER"),
					resource.TestCheckResourceAttrPair(resourceName, "authorize_accounts.0.domain",
						"data.huaweicloud_workspace_service.test", "ad_domain.0.name"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "is_vdi", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "ou_name", acceptance.HW_WORKSPACE_OU_NAME),
					resource.TestCheckResourceAttr(resourceName, "extra_session_type", "CPU"),
					resource.TestCheckResourceAttr(resourceName, "extra_session_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "route_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_policy.0.max_session", "3"),
					resource.TestCheckResourceAttr(resourceName, "route_policy.0.cpu_threshold", "80"),
					resource.TestCheckResourceAttr(resourceName, "route_policy.0.mem_threshold", "80"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", getAcceptanceEpsId()),
					resource.TestCheckResourceAttr(resourceName, "is_delete_associated_resources", "true"),
				),
			},
			{
				Config: testResourceAppImageServer_withAD(updateName, ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"flavor_id",
					"vpc_id",
					"subnet_id",
					"root_volume",
					"image_source_product_id",
					"is_vdi",
					"availability_zone",
					"ou_name",
					"extra_session_type",
					"extra_session_size",
					"route_policy",
					"scheduler_hints",
					"tags",
					"enterprise_project_id",
					"is_delete_associated_resources",
				},
			},
		},
	})
}

func testResourceAppImageServer_withAD(name, description string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_image_server" "test" {
  name                    = "%[1]s"
  flavor_id               = "%[2]s"
  vpc_id                  = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id               = data.huaweicloud_workspace_service.test.network_ids[0]
  image_id                = "%[3]s"
  image_type              = "gold"
  image_source_product_id = "%[4]s"
  spec_code               = "%[5]s"

  # Currently only one user can be set.
  authorize_accounts {
    account = data.huaweicloud_workspace_service.test.ad_domain[0].admin_account
    type    = "USER"
    domain  = data.huaweicloud_workspace_service.test.ad_domain[0].name
  }

  root_volume {
    type = "SAS"
    size = 80
  }

  is_vdi             = false
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  description        = "%[6]s"
  ou_name            = "%[7]s"
  extra_session_type = "CPU"
  extra_session_size = 2

  route_policy {
    max_session   = 3
    cpu_threshold = 80
    mem_threshold = 80
  }

  tags = {
    foo = "bar"
  }

  enterprise_project_id          = "%[8]s"
  is_delete_associated_resources = true
}
`,
		name,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_SPEC_CODE,
		description,
		acceptance.HW_WORKSPACE_OU_NAME,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
