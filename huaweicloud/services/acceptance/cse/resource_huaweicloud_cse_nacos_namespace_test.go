package cse

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cse"
)

func getNacosNamespaceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cse", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSE client: %s", err)
	}

	return cse.GetNacosNamespaceById(client, state.Primary.Attributes["engine_id"],
		state.Primary.Attributes["enterprise_project_id"], state.Primary.ID)
}

func TestAccNacosNamespace_basic(t *testing.T) {
	var (
		obj interface{}

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		randUUID, _ = uuid.NewRandom()

		resourceName = "huaweicloud_cse_nacos_namespace.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getNacosNamespaceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSENacosMicroserviceEngineID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNacosNamespace_basic(randUUID.String(), name),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf(`unable to create the namespace because the Nacos engine \(%s\) does not exist`, randUUID.String())),
			},
			{
				Config: testAccNacosNamespace_basic(acceptance.HW_CSE_NACOS_MICROSERVICE_ENGINE_ID, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccNacosNamespace_basic(acceptance.HW_CSE_NACOS_MICROSERVICE_ENGINE_ID, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNacosNamespaceImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccNacosNamespaceImportStateIdFunc(resName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var engineId, namespaceId, enterpriseProjectId string
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", resName)
		}

		engineId = rs.Primary.Attributes["engine_id"]
		namespaceId = rs.Primary.ID
		enterpriseProjectId = rs.Primary.Attributes["enterprise_project_id"]

		if engineId == "" || namespaceId == "" || enterpriseProjectId == "" {
			return "", fmt.Errorf("missing some attributes, want '<engine_id>/<id>/<enterprise_project_id>', but got '%s/%s/%s'",
				engineId, namespaceId, enterpriseProjectId)
		}

		return fmt.Sprintf("%s/%s/%s", engineId, namespaceId, enterpriseProjectId), nil
	}
}

func testAccNacosNamespace_basic(engineId, name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  default = "%[1]s"
}

resource "huaweicloud_cse_nacos_namespace" "test" {
  engine_id             = "%[2]s"
  name                  = "%[3]s"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, engineId, name)
}
