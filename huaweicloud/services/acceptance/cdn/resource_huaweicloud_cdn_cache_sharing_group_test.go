package cdn

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdn"
)

func getCacheSharingGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return nil, fmt.Errorf("error creating CDN client: %s", err)
	}

	return cdn.GetCacheSharingGroupById(client, state.Primary.ID)
}

func TestAccCacheSharingGroup_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_cdn_cache_sharing_group.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCacheSharingGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCdntDomainNames(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCacheSharingGroup_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "primary_domain"),
					resource.TestCheckResourceAttr(rName, "share_cache_records.#", "2"),
					resource.TestMatchResourceAttr(rName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccCacheSharingGroup_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "primary_domain"),
					resource.TestCheckResourceAttr(rName, "share_cache_records.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCacheSharingGroupImportStateWithName(rName),
			},
		},
	})
}

func testCacheSharingGroupImportStateWithName(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		groupName := rs.Primary.Attributes["name"]
		if groupName == "" {
			return "", errors.New("The share cache group name is missing, want '<name>'")
		}
		return groupName, nil
	}
}

func testAccCacheSharingGroup_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_cache_sharing_group" "test" {
  name           = "%[1]s"
  primary_domain = try(element(split(",", "%[2]s"), 0), "")

  dynamic "share_cache_records" {
    for_each = split(",", "%[2]s")

    content {
      domain_name = share_cache_records.value
    }
  }
}
`, name, acceptance.HW_CDN_DOMAIN_NAMES)
}

func testAccCacheSharingGroup_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_cache_sharing_group" "test" {
  name           = "%[1]s"
  primary_domain = try(element(split(",", "%[2]s"), 0), "")

  dynamic "share_cache_records" {
    for_each = slice(split(",", "%[2]s"), 0, 1)

    content {
      domain_name = share_cache_records.value
    }
  }
}
`, name, acceptance.HW_CDN_DOMAIN_NAMES)
}
