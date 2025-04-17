package servicestage

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/servicestage"
)

func getV3RuntimeStackFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("servicestage", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage client: %s", err)
	}
	return servicestage.GetV3RuntimeStackById(client, state.Primary.ID)
}

func TestAccV3RuntimeStack_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_servicestagev3_runtime_stack.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getV3RuntimeStackFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// At least two ZIP files must be provided.
			acceptance.TestAccPreCheckServiceStageZipStorageURLs(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3RuntimeStack_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "deploy_mode", "virtualmachine"),
					resource.TestCheckResourceAttr(resourceName, "type", "Java"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttrSet(resourceName, "spec"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "creator"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(resourceName, "component_count"),
				),
			},
			{
				Config: testAccV3RuntimeStack_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "deploy_mode", "virtualmachine"),
					resource.TestCheckResourceAttr(resourceName, "type", "Java"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttrSet(resourceName, "spec"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "creator"),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(resourceName, "component_count"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"spec_origin",
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3RuntimeStackImportStateFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"spec_origin",
				},
			},
		},
	})
}

func testAccV3RuntimeStackImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("invalid format specified for runtime stack name ('%s')",
				rs.Primary.Attributes["name"])
		}
		return rs.Primary.Attributes["name"], nil
	}
}

func testAccV3RuntimeStack_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_servicestagev3_runtime_stack" "test" {
  name        = "%[1]s"
  deploy_mode = "virtualmachine"
  type        = "Java"
  version     = "1.0.1"
  spec        = jsonencode({
    "parameters": {
      "jdk_url": try(element(split(",", "%[2]s"), 0), "")
    }
  })
  description = "Created by terraform script"
}
`, name, acceptance.HW_SERVICESTAGE_ZIP_STORAGE_URLS)
}

func testAccV3RuntimeStack_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_servicestagev3_runtime_stack" "test" {
  name        = "%[1]s"
  deploy_mode = "virtualmachine"
  type        = "Java"
  version     = "1.0.2"
  spec        = jsonencode({
    "parameters": {
      "jdk_url": try(element(split(",", "%[2]s"), 1), "")
    }
  })
}
`, name, acceptance.HW_SERVICESTAGE_ZIP_STORAGE_URLS)
}
