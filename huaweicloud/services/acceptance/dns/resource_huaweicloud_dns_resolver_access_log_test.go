package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dns"
)

func getResolverAccessLog(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("dns_region", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client: %s", err)
	}

	return dns.GetResolverAccessLog(client, state.Primary.ID)
}

func TestAccResolverAccessLog_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dns_resolver_access_log.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getResolverAccessLog)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverAccessLog_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "lts_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "lts_topic_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "vpc_ids.#", "3"),
				),
			},
			{
				Config: testAccResolverAccessLog_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "vpc_ids.#", "2"),
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

func testAccResolverAccessLog_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}

resource "huaweicloud_vpc" "test" {
  count = 4

  name = "%[1]s_${count.index}"
  cidr = "192.168.0.0/16"
}`, name)
}

func testAccResolverAccessLog_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dns_resolver_access_log" "test" {
  lts_group_id = huaweicloud_lts_group.test.id
  lts_topic_id = huaweicloud_lts_stream.test.id
  vpc_ids      = slice(huaweicloud_vpc.test[*].id, 0, 3)
}`, testAccResolverAccessLog_base(name))
}

func testAccResolverAccessLog_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dns_resolver_access_log" "test" {
  lts_group_id = huaweicloud_lts_group.test.id
  lts_topic_id = huaweicloud_lts_stream.test.id
  vpc_ids      = slice(huaweicloud_vpc.test[*].id, 2, 4)
}`, testAccResolverAccessLog_base(name))
}
