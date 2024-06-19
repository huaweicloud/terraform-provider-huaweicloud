package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccLogstashConnectivity_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_logstash_connectivity.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLogstashConnectivity_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccLogstashConnectivity_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_connectivity" "test" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  address_and_ports {
    address = split(":",huaweicloud_css_logstash_cluster.test.endpoint)[0]
    port    = parseint(split(":",huaweicloud_css_logstash_cluster.test.endpoint)[1], 10)
  }
}
`, testAccLogstashCluster_basic(rName, 1, "bar"))
}
