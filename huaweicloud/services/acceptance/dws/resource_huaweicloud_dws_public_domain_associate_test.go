package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getPublicDomainNameAssociatedFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	return dws.GetDomainNameInfoByClusterId(client, state.Primary.Attributes["cluster_id"])
}

func TestAccPublicDomainNameAssociate_basic(t *testing.T) {
	var (
		obj        interface{}
		rName      = "huaweicloud_dws_public_domain_associate.test"
		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			rName,
			&obj,
			getPublicDomainNameAssociatedFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPublicDomainNameAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "domain_name", name),
				),
			},
			{
				Config: testPublicDomainNameAssociate_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", updateName),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testPublicDomainNameAssociateImportState(rName),
				ImportStateVerifyIgnore: []string{"ttl"},
			},
		},
	})
}

func testPublicDomainNameAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_public_domain_associate" "test" {
  cluster_id  = "%[1]s"
  domain_name = "%[2]s"
  ttl         = 400
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}

func testPublicDomainNameAssociate_basic_step2(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_public_domain_associate" "test" {
  cluster_id  = "%[1]s"
  domain_name = "%[2]s"
  ttl         = 1000
}
`, acceptance.HW_DWS_CLUSTER_ID, updateName)
}

func testPublicDomainNameAssociateImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		domainName := rs.Primary.ID
		if clusterId == "" || domainName == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<cluster_id>/<domain_name>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}

		return fmt.Sprintf("%s/%s", clusterId, domainName), nil
	}
}
