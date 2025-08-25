package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/coc"
)

func getApplicationResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	return coc.GetApplication(client, state.Primary.ID)
}

func TestAccResourceApplication_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_application.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getApplicationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApplication_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "is_collection", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "code"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApplication_updated(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", updateName),
					resource.TestCheckResourceAttr(resourceName, "is_collection", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "code"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
		},
	})
}

func testAccApplication_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_application" "test" {
  name          = "%[1]s"
  description   = "%[1]s"
  is_collection = true
}
`, name)
}

func testAccApplication_updated(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_application" "test" {
  name          = "%[1]s"
  description   = "%[1]s"
  is_collection = false
}
`, name)
}
