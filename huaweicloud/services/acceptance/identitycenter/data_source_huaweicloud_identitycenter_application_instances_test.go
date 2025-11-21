package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceApplicationInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identitycenter_application_instances.test"
	uuid, _ := uuid.GenerateUUID()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckIdentityCenterAccountId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceApplicationInstances_basic(uuid),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "application_instances.0.name", uuid),
					resource.TestCheckResourceAttr(dataSource, "application_instances.0.display_name", "create"),
					resource.TestCheckResourceAttr(dataSource, "application_instances.0.description", "create"),
					resource.TestCheckResourceAttrSet(dataSource, "application_instances.0.response_config"),
					resource.TestCheckResourceAttrSet(dataSource, "application_instances.0.response_schema_config"),
					resource.TestCheckResourceAttr(dataSource,
						"application_instances.0.service_provider_config.0.audience", "https://create.com"),
					resource.TestCheckResourceAttr(dataSource,
						"application_instances.0.service_provider_config.0.consumers.0.location", "https://create.com"),
					resource.TestCheckResourceAttr(dataSource, "application_instances.0.security_config.0.ttl", "P9M"),
					resource.TestCheckResourceAttr(dataSource, "application_instances.0.status", "CREATED"),
				),
			},
		},
	})
}

func testDataSourceApplicationInstances_basic(uuid string) string {
	return fmt.Sprintf(`
%[1]s
 
data "huaweicloud_identitycenter_application_instances" "test"{
  depends_on  = [huaweicloud_identitycenter_application_instance.test]
  instance_id = data.huaweicloud_identitycenter_instance.test.id
}
`, testApplicationInstance_basic(uuid))
}
