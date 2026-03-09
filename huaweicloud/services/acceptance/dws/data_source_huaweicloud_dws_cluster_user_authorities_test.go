package dws

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataClusterUserAuthorities_basic(t *testing.T) {
	dcName := "data.huaweicloud_dws_cluster_user_authorities.test"
	dc := acceptance.InitDataSourceCheck(dcName)
	userName := strings.Split(acceptance.HW_DWS_ASSOCIATE_USER_NAMES, ",")[0]

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
			acceptance.TestAccPreCheckDwsClusterUserNames(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataClusterUserAuthorities_clusterNotFound(userName),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testAccDataClusterUserAuthorities_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(dcName, "name", userName),
					resource.TestMatchResourceAttr(dcName, "authorities.#", regexp.MustCompile(`^([0-9]|[1-9][0-9]*)$`)),
					resource.TestCheckResourceAttrSet(dcName, "authorities.0.database"),
					resource.TestCheckResourceAttrSet(dcName, "authorities.0.all_object"),
					resource.TestCheckResourceAttrSet(dcName, "authorities.0.future"),
					resource.TestMatchResourceAttr(dcName, "authorities.0.privileges.#", regexp.MustCompile(`^([0-9]|[1-9][0-9]*)$`)),
					resource.TestCheckResourceAttrSet(dcName, "authorities.0.privileges.0.permission"),
					resource.TestCheckResourceAttrSet(dcName, "authorities.0.privileges.0.grant_with"),
				),
			},
		},
	})
}

func testAccDataClusterUserAuthorities_clusterNotFound(userName string) string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_user_authorities" "test" {
  cluster_id = "%s"
  name       = "%s"
}
`, randUUID, userName)
}

func testAccDataClusterUserAuthorities_basic(userName string) string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_user_authorities" "test" {
  cluster_id = "%s"
  name       = "%s"
}
`, acceptance.HW_DWS_CLUSTER_ID, userName)
}
