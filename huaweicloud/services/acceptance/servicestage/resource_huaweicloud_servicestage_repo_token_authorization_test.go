package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/servicestage/v1/repositories"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getAuthResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ServiceStageV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage v1 client: %s", err)
	}

	resp, err := repositories.List(c)
	if err != nil {
		return resp, err
	}
	for _, v := range resp {
		if v.Name == state.Primary.ID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("error getting ServiceStage authorization (%s)", state.Primary.ID)
}

func TestAccRepoTokenAuth_basic(t *testing.T) {
	var (
		auth         repositories.Authorization
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_servicestage_repo_token_authorization.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&auth,
		getAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRepoTokenAuth(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRepoTokenAuth_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "github"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"token",
					"host",
				},
			},
		},
	})
}

func testAccRepoTokenAuth_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_servicestage_repo_token_authorization" "test" {
  type  = "github"
  name  = "%s"
  host  = "%s"
  token = "%s"
}
`, rName, acceptance.HW_GITHUB_REPO_HOST, acceptance.HW_GITHUB_PERSONAL_TOKEN)
}
