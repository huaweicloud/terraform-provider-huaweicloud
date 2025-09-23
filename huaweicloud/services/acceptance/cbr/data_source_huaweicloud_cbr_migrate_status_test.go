package cbr

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func checkRegionIDs(resourceName string, expectMultiple bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		regionIDs := make(map[string]struct{})
		count, err := strconv.Atoi(rs.Primary.Attributes["project_status.#"])
		if err != nil {
			return fmt.Errorf("failed to parse project_status count: %v", err)
		}

		for i := 0; i < count; i++ {
			key := fmt.Sprintf("project_status.%d.region_id", i)
			regionID := rs.Primary.Attributes[key]
			regionIDs[regionID] = struct{}{}
		}

		if expectMultiple {
			if len(regionIDs) <= 1 {
				return fmt.Errorf("expected multiple region_ids, but found %d unique values", len(regionIDs))
			}
		} else {
			if len(regionIDs) > 1 {
				return fmt.Errorf("expected single region_id, but found %d unique values", len(regionIDs))
			}
		}
		return nil
	}
}

func TestAccDataSourceMigrateStatus_basic(t *testing.T) {
	var (
		basic        = "data.huaweicloud_cbr_migrate_status.basic"
		dcBasic      = acceptance.InitDataSourceCheck(basic)
		allRegions   = "data.huaweicloud_cbr_migrate_status.all_regions"
		dcAllRegions = acceptance.InitDataSourceCheck(allRegions)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMigrateStatus_basic,
				Check: resource.ComposeTestCheckFunc(
					dcBasic.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(basic, "id"),
					resource.TestCheckResourceAttrSet(basic, "region"),
					resource.TestCheckResourceAttrSet(basic, "status"),
					resource.TestMatchResourceAttr(basic, "status", regexp.MustCompile(`^(migrating|success|failed)$`)),
					resource.TestCheckResourceAttrSet(basic, "project_status.#"),
					resource.TestCheckResourceAttrSet(basic, "project_status.0.status"),
					resource.TestCheckResourceAttrSet(basic, "project_status.0.project_id"),
					resource.TestCheckResourceAttrSet(basic, "project_status.0.project_name"),
					resource.TestCheckResourceAttrSet(basic, "project_status.0.region_id"),
					resource.TestCheckResourceAttrSet(basic, "project_status.0.progress"),
					checkRegionIDs(basic, false),
				),
			},
			{
				Config: testAccDataSourceMigrateStatus_withAllRegions,
				Check: resource.ComposeTestCheckFunc(
					dcAllRegions.CheckResourceExists(),
					resource.TestMatchResourceAttr(allRegions, "status", regexp.MustCompile(`^(migrating|success|failed)$`)),
					resource.TestCheckResourceAttrSet(allRegions, "project_status.#"),
					resource.TestCheckResourceAttrSet(allRegions, "id"),
					checkRegionIDs(allRegions, true),
				),
			},
		},
	})
}

const testAccDataSourceMigrateStatus_basic = `
data "huaweicloud_cbr_migrate_status" "basic" {}
`

const testAccDataSourceMigrateStatus_withAllRegions = `
data "huaweicloud_cbr_migrate_status" "all_regions" {
  all_regions = true
}
`
