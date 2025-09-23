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
		auth interface{}

		name = acceptance.RandomAccResourceNameWithDash()

		withHost    = "huaweicloud_servicestage_repo_token_authorization.with_host"
		withoutHost = "huaweicloud_servicestage_repo_token_authorization.without_host"

		rcWithHost    = acceptance.InitResourceCheck(withHost, &auth, getAuthResourceFunc)
		rcWithoutHost = acceptance.InitResourceCheck(withoutHost, &auth, getAuthResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRepoTokenAuth(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithHost.CheckResourceDestroy(),
			rcWithoutHost.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccRepoTokenAuth_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithHost.CheckResourceExists(),
					resource.TestCheckResourceAttr(withHost, "name", name+"-with-host"),
					resource.TestCheckResourceAttr(withHost, "type", "github"),
					resource.TestCheckResourceAttr(withHost, "host", "https://api.github.com"),
					rcWithoutHost.CheckResourceExists(),
					resource.TestCheckResourceAttr(withoutHost, "name", name+"-without-host"),
					resource.TestCheckResourceAttr(withoutHost, "type", "github"),
					resource.TestCheckNoResourceAttr(withoutHost, "host"),
				),
			},
			{
				ResourceName:      withHost,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"token",
					"host",
				},
			},
			{
				ResourceName:      withoutHost,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"token",
				},
			},
		},
	})
}

func testAccRepoTokenAuth_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_servicestage_repo_token_authorization" "with_host" {
  type  = "github"
  name  = "%[1]s-with-host"
  token = "%[2]s"
  host  = "https://api.github.com"
}

resource "huaweicloud_servicestage_repo_token_authorization" "without_host" {
  type  = "github"
  name  = "%[1]s-without-host"
  token = "%[2]s"
}
`, rName, acceptance.HW_GITHUB_PERSONAL_TOKEN)
}
