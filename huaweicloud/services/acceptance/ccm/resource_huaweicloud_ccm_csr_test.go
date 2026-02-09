package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ccm"
)

func getCCMCsrResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("scm", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCM client: %s", err)
	}

	return ccm.GetCsr(client, state.Primary.ID)
}

func TestAccCCMCsr_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_ccm_csr.test"
		name  = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCCMCsrResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCCMCsr_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "domain_name", "test.com"),
					resource.TestCheckResourceAttr(rName, "private_key_algo", "RSA_2048"),
					resource.TestCheckResourceAttr(rName, "usage", "PERSONAL"),
					resource.TestCheckResourceAttrSet(rName, "csr"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				Config: testCCMCsr_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
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

func testCCMCsr_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_csr" "test" {
  name             = "%s"
  domain_name      = "test.com"
  private_key_algo = "RSA_2048"
  usage            = "PERSONAL"
}
`, name)
}

func testCCMCsr_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_csr" "test" {
  name             = "%s_update"
  domain_name      = "test.com"
  private_key_algo = "RSA_2048"
  usage            = "PERSONAL"
}
`, name)
}
