package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sdrs/v1/protectiongroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getProtectionGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SDRS Client: %s", err)
	}
	return protectiongroups.Get(client, state.Primary.ID).Extract()
}

// Lack of testing for `enable`, will test it in resource replication pair
func TestAccProtectionGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_sdrs_protection_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProtectionGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testProtectionGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttrPair(rName, "source_availability_zone", "data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "target_availability_zone", "data.huaweicloud_availability_zones.test", "names.1"),
					resource.TestCheckResourceAttrPair(rName, "domain_id", "data.huaweicloud_sdrs_domain.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "source_vpc_id", "huaweicloud_vpc.test", "id"),
				),
			},
			{
				Config: testProtectionGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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

func testProtectionGroup_base(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_sdrs_domain" "test" {}
data "huaweicloud_availability_zones" "test" {}
`, common.TestVpc(name))
}

func testProtectionGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sdrs_protection_group" "test" {
  name                     = "%[2]s"
  source_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  target_availability_zone = data.huaweicloud_availability_zones.test.names[1]
  domain_id                = data.huaweicloud_sdrs_domain.test.id
  source_vpc_id            = huaweicloud_vpc.test.id
  description              = "test description"
}
`, testProtectionGroup_base(name), name)
}

func testProtectionGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sdrs_protection_group" "test" {
  name                     = "%[2]s_update"
  source_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  target_availability_zone = data.huaweicloud_availability_zones.test.names[1]
  domain_id                = data.huaweicloud_sdrs_domain.test.id
  source_vpc_id            = huaweicloud_vpc.test.id
  description              = "test description"
}
`, testProtectionGroup_base(name), name)
}
