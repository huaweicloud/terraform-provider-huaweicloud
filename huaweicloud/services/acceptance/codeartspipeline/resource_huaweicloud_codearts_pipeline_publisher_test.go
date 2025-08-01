package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/codeartspipeline"
)

func getPipelinePublisherResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	return codeartspipeline.GetPipelinePublisher(client, cfg.DomainID, state.Primary.ID)
}

func TestAccPipelinePublisher_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_publisher.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelinePublisherResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelinePublisher_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "en_name", name+"-en"),
					resource.TestCheckResourceAttr(rName, "description", "desc"),
					resource.TestCheckResourceAttr(rName, "website", "https://github.com"),
					resource.TestCheckResourceAttr(rName, "support_url", "terraform"),
					resource.TestCheckResourceAttr(rName, "source_url", "source_address"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testPipelinePublisher_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_pipeline_publisher" "test" {
  name        = "%[1]s"
  en_name     = "%[1]s-en"
  description = "desc"
  website     = "https://github.com"
  support_url = "terraform"
  source_url  = "source_address"
}
`, name)
}
