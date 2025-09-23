package servicestage

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/servicestage"
)

func getV3ReleasedRuntimeStacksFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("servicestage", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage client: %s", err)
	}

	runtimeStackIdsLen, err := strconv.Atoi(state.Primary.Attributes["runtime_stack_ids.#"])
	if err != nil {
		return nil, fmt.Errorf("failed to convert the %v value from string to int", state.Primary.Attributes["runtime_stack_ids.#"])
	}

	runtimeStackIds := make([]interface{}, 0, runtimeStackIdsLen)
	for i := 0; i < runtimeStackIdsLen; i++ {
		runtimeStackIds = append(runtimeStackIds, state.Primary.Attributes[fmt.Sprintf("runtime_stack_ids.%d", i)])
	}

	return servicestage.FilterV3ReleasedRuntimeStacks(client, runtimeStackIds)
}

func TestAccV3RuntimeStackBatchRelease_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_servicestagev3_runtime_stack_batch_release.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getV3ReleasedRuntimeStacksFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// At least one of ZIP file must be provided.
			acceptance.TestAccPreCheckServiceStageZipStorageURLs(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3RuntimeStackBatchRelease_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack_ids.#", "2"),
				),
			},
			{
				Config: testAccV3RuntimeStackBatchRelease_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack_ids.#", "2"),
				),
			},
		},
	})
}

func testAccV3RuntimeStackBatchRelease_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_servicestagev3_runtime_stack" "test" {
  count = 3

  name        = format("%[1]s_%%d", count.index)
  deploy_mode = "virtualmachine"
  type        = "Java"
  version     = format("1.0.%%d", count.index)
  spec        = jsonencode({
    "parameters": {
      "jdk_url": try(element(split(",", "%[2]s"), 0), "")
    }
  })
}
`, name, acceptance.HW_SERVICESTAGE_ZIP_STORAGE_URLS)
}

func testAccV3RuntimeStackBatchRelease_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_runtime_stack_batch_release" "test" {
  runtime_stack_ids = slice(huaweicloud_servicestagev3_runtime_stack.test[*].id, 0, 2)
}
`, testAccV3RuntimeStackBatchRelease_base(name))
}

func testAccV3RuntimeStackBatchRelease_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_runtime_stack_batch_release" "test" {
  runtime_stack_ids = slice(huaweicloud_servicestagev3_runtime_stack.test[*].id, 1, 3)
}
`, testAccV3RuntimeStackBatchRelease_base(name))
}
