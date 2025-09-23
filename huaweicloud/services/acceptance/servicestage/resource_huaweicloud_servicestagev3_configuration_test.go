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

func getV3Configuration(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("servicestage", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage client: %s", err)
	}
	return servicestage.GetV3ConfigurationFile(client, state.Primary.ID)
}

func TestAccV3Configuration_basic(t *testing.T) {
	var (
		configuration interface{}
		resourceName  = "huaweicloud_servicestagev3_configuration.test"
		rc            = acceptance.InitResourceCheck(resourceName, &configuration, getV3Configuration)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3Configuration_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "config_group_id", "huaweicloud_servicestagev3_configuration_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "properties"),
					resource.TestCheckResourceAttr(resourceName, "content", "testkey = testvalue"),
					resource.TestCheckResourceAttr(resourceName, "sensitive", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by TF script"),
					resource.TestCheckResourceAttr(resourceName, "components.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "creator"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV3Configuration_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "yaml"),
					resource.TestCheckResourceAttr(resourceName, "sensitive", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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

func testAccV3Configuration_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_configuration" "test" {
  config_group_id = huaweicloud_servicestagev3_configuration_group.test.id
  name            = "%[2]s"
  type            = "properties"
  content         = "testkey = testvalue"
  sensitive       = true
  description     = "Created by TF script"
}
`, testAccV3ConfigurationGroup_basic(name), name)
}

func testAccV3Configuration_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_configuration" "test" {
  config_group_id = huaweicloud_servicestagev3_configuration_group.test.id
  name            = "%[2]s"
  type            = "yaml"
  content         = <<EOF
spring:
  application:
    name: "service"
  cloud:
    servicecomb:
      service:
        name: service
        version: $${CAS_INSTANCE_VERSION}
        application: $${CAS_APPLICATION_NAME}
      discovery:
        address: $${PAAS_CSE_SC_ENDPOINT}
        healthCheckInterval: 10
        pollInterval: 15000
        waitTimeForShutDownInMillis: 15000
      config:
        serverAddr: $${PAAS_CSE_CC_ENDPOINT}
        serverType: kie
        fileSource: governance.yaml,application.yaml
EOF
}
`, testAccV3ConfigurationGroup_basic(name), name)
}
