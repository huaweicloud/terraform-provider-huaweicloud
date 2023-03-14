package ddm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdmSchemas_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	instanceName := strings.ReplaceAll(name, "_", "-")
	rName := "data.huaweicloud_ddm_schemas.test"
	dbPwd := "test_1234"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmSchemas_basic(instanceName, name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "schemas.0.name", name),
					resource.TestCheckResourceAttr(rName, "schemas.0.status", "RUNNING"),
					resource.TestCheckResourceAttr(rName, "schemas.0.shard_mode", "single"),
					resource.TestCheckResourceAttr(rName, "schemas.0.shard_number", "1"),
					resource.TestCheckResourceAttr(rName, "schemas.0.data_nodes.#", "1"),
				),
			},
		},
	})
}

func testAccDatasourceDdmSchemas_basic(instanceName, name, dbPwd string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ddm_schemas" "test" {
  instance_id = huaweicloud_ddm_instance.test.id
  name        = huaweicloud_ddm_schema.test.name
}
`, testDdmSchema_basic(instanceName, name, dbPwd))
}
