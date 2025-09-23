package servicestage

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/servicestage"
)

func getV3EnvFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("servicestage", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage client: %s", err)
	}
	return servicestage.QueryV3Environment(client, state.Primary.ID)
}

func TestAccV3Environment_basic(t *testing.T) {
	var (
		env interface{}

		resourceName = "huaweicloud_servicestagev3_environment.test"
		rc           = acceptance.InitResourceCheck(resourceName, &env, getV3EnvFunc)

		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()
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
				Config: testAccV3Environment_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform test"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "creator"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV3Environment_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "baar"),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV3Environment_basic_step3(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
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

func testAccV3Environment_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_environment" "test" {
  name                  = "%[2]s"
  description           = "Created by terraform test"
  vpc_id                = huaweicloud_vpc.test.id
  enterprise_project_id = "0"

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, common.TestVpc(name), name)
}

func testAccV3Environment_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_environment" "test" {
  name                  = "%[2]s"
  vpc_id                = huaweicloud_vpc.test.id
  enterprise_project_id = "%[3]s"

  tags = {
    foo = "baar"
  }
}
`, common.TestVpc(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccV3Environment_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_environment" "test" {
  name                  = "%[2]s"
  vpc_id                = huaweicloud_vpc.test.id
  enterprise_project_id = "%[3]s"
}
`, common.TestVpc(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
