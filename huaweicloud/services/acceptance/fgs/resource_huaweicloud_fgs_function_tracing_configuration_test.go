package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
)

func getFunctionTracingConfigurationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("fgs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}

	return fgs.GetFunctionTracingConfiguration(client, state.Primary.Attributes["function_urn"])
}

func TestAccFunctionTracingConfiguration_basic(t *testing.T) {
	var (
		obj interface{}

		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_fgs_function_tracing_configuration.test"

		rc = acceptance.InitResourceCheck(resourceName, &obj, getFunctionTracingConfigurationFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "3.0.0",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionTracingConfiguration_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "tracing_ak", acceptance.HW_ACCESS_KEY),
					resource.TestCheckResourceAttr(resourceName, "tracing_sk", acceptance.HW_SECRET_KEY),
					resource.TestCheckResourceAttrSet(resourceName, "function_urn"),
				),
			},
			{
				Config: testAccFunctionTracingConfiguration_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "tracing_ak", "random_string.access_key", "result"),
					resource.TestCheckResourceAttrPair(resourceName, "tracing_sk", "random_string.secret_key", "result"),
					resource.TestCheckResourceAttrSet(resourceName, "function_urn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"tracing_sk",
				},
			},
		},
	})
}

func testAccFunctionTracingConfiguration_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name                   = "%[1]s"
  runtime                = "Java11"
  app                    = "default"
  handler                = "com.huawei.demo.TriggerTests.apigTest"
  memory_size            = 512
  timeout                = 15
  code_type              = "zip"
  code_filename          = "java-demo.zip"
  agency                 = "%[2]s"
  functiongraph_version  = "v2"
  enable_class_isolation = true
}
`, name, acceptance.HW_FGS_AGENCY_NAME)
}

func testAccFunctionTracingConfiguration_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function_tracing_configuration" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  tracing_ak   = "%[2]s"
  tracing_sk   = "%[3]s"
}
`, testAccFunctionTracingConfiguration_base(name), acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}

func testAccFunctionTracingConfiguration_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "random_string" "access_key" {
  length  = 20
  special = false
}

resource "random_string" "secret_key" {
  length  = 40
  special = false
}

resource "huaweicloud_fgs_function_tracing_configuration" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  tracing_ak   = random_string.access_key.result
  tracing_sk   = random_string.secret_key.result
}
`, testAccFunctionTracingConfiguration_base(name))
}
