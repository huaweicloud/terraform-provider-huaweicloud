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

func getComponentResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	return coc.GetComponent(client, state.Primary.Attributes["application_id"], state.Primary.ID)
}

func TestAccResourceComponent_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	nameUpdate := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_component.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getComponentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComponent_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "code"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComponent_updated(nameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdate),
					resource.TestCheckResourceAttrSet(resourceName, "code"),
				),
			},
		},
	})
}

func testAccComponent_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_component" "test" {
  application_id = huaweicloud_coc_application.test.id
  name           = "%[2]s"
}
`, testAccApplication_basic(name), name)
}

func testAccComponent_updated(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_component" "test" {
  application_id = huaweicloud_coc_application.test.id
  name           = "%[2]s"
}
`, testAccApplication_basic(name), name)
}
