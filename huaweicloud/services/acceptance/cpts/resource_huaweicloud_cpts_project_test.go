package cpts

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cpts/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getProjectResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcCptsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CPTS v1 client: %s", err)
	}

	id, err := strconv.ParseInt(state.Primary.ID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("the project ID must be integer: %s", err)
	}

	request := &model.ShowProjectRequest{
		TestSuiteId: int32(id),
	}

	return client.ShowProject(request)
}

func TestAccProject_basic(t *testing.T) {
	var obj model.CreateProjectResponse

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cpts_project.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getProjectResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testProject_basic(rName, "created by acc test"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testProject_basic(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
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

func testProject_basic(rName, desc string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cpts_project" "test" {
  name        = "%s"
  description = "%s"
}
`, rName, desc)
}
