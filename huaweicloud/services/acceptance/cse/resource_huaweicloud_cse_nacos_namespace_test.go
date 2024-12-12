package cse

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
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

	return cse.GetNacosNamespaceById(client, state.Primary.Attributes["engine_id"], state.Primary.ID)
}

func TestAccNacosNamespace_basic(t *testing.T) {
	var (
		obj interface{}

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		uuid, _ = uuid.GenerateUUID()

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
				Config:      testAccNacosNamespace_basic(uuid, name),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`unable to create the namespace because the Nacos engine \(%s\) does not exist`, uuid)),
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
		var engineId, namespaceId string
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", resName)
		}

		engineId = rs.Primary.Attributes["engine_id"]
		namespaceId = rs.Primary.ID

		if engineId == "" || namespaceId == "" {
			return "", fmt.Errorf("missing some attributes, want '<engine_id>/<id>', but got '%s/%s'", engineId, namespaceId)
		}

		return fmt.Sprintf("%s/%s", engineId, namespaceId), nil
	}
}

func testAccNacosNamespace_basic(engineId, name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cse_nacos_namespace" "test" {
  engine_id = "%[1]s"
  name      = "%[2]s"
}
`, engineId, name)
}
