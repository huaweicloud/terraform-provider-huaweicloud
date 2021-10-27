package cce

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/cce/v1/namespaces"
)

func getNamespaceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CceV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud CCE v1 client: %s", err)
	}
	resp, err := namespaces.Get(c, state.Primary.Attributes["cluster_id"],
		state.Primary.Attributes["name"]).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("Unable to find the namespace (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccCCENamespaceV1_basic(t *testing.T) {
	var namespace namespaces.Namespace
	resourceName := "huaweicloud_cce_namespace.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&namespace,
		getNamespaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCCENamespaceV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "cluster_id",
						"${huaweicloud_cce_cluster.test.id}"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENamespaceImportStateIdFunc(randName),
			},
		},
	})
}

func TestAccCCENamespaceV1_generateName(t *testing.T) {
	var namespace namespaces.Namespace
	resourceName := "huaweicloud_cce_namespace.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&namespace,
		getNamespaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCCENamespaceV1_generateName(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "cluster_id",
						"${huaweicloud_cce_cluster.test.id}"),
					resource.TestCheckResourceAttr(resourceName, "prefix", randName),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(fmt.Sprintf(`^%s[a-z0-9-]*`, randName))),
				),
			},
		},
	})
}

func testAccCCENamespaceImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_cce_cluster" {
				clusterID = rs.Primary.ID
			}
		}
		if clusterID == "" || name == "" {
			return "", fmtp.Errorf("resource not found: %s/%s", clusterID, name)
		}
		return fmt.Sprintf("%s/%s", clusterID, name), nil

	}
}

func testAccCCENamespaceV1_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_namespace" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  name       = "%s"
}
`, testAccCceCluster_config(rName), rName)
}

func testAccCCENamespaceV1_generateName(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_namespace" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  prefix     = "%s"
}
`, testAccCceCluster_config(rName), rName)
}
