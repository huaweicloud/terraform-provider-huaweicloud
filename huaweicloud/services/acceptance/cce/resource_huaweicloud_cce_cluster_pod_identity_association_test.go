package cce

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPodIdentityAssociationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cce", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE client: %s", err)
	}

	getHttpUrl := "api/v3/projects/{project_id}/clusters/{cluster_id}/pod-identity-associations/{association_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", state.Primary.Attributes["cluster_id"])
	getPath = strings.ReplaceAll(getPath, "{association_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE client: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}
	associationId := utils.PathSearch("uid", respBody, "")
	if associationId == "" {
		return nil, golangsdk.ErrDefault404{}
	}
	return respBody, nil
}

func TestAccCcePodIdentityAssociation_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_cce_cluster_pod_identity_association.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPodIdentityAssociationResourceFunc,
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcePodIdentityAssociation_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "namespace", "default"),
					resource.TestCheckResourceAttr(rName, "service_account", "default"),
					resource.TestCheckResourceAttr(rName, "agency_name", "CCENodeAgency"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
				),
			},
			{
				Config: testAccCcePodIdentityAssociation_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "namespace", "default"),
					resource.TestCheckResourceAttr(rName, "service_account", "default"),
					resource.TestCheckResourceAttr(rName, "agency_name", "CCEServiceAgency"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCcePodIdentityAssociationImportStateIdFunc(rName),
			},
		},
	})
}

func testAccCcePodIdentityAssociation_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cce_cluster_pod_identity_association" "test" {
  cluster_id      = "%[1]s"
  namespace       = "default"
  service_account = "default"
  agency_name     = "CCENodeAgency"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, acceptance.HW_CCE_CLUSTER_ID)
}

func testAccCcePodIdentityAssociation_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_cce_cluster_pod_identity_association" "test" {
  cluster_id      = "%[1]s"
  namespace       = "default"
  service_account = "default"
  agency_name     = "CCEServiceAgency"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, acceptance.HW_CCE_CLUSTER_ID)
}

func testAccCcePodIdentityAssociationImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}
		clusterID := rs.Primary.Attributes["cluster_id"]
		if clusterID == "" {
			return "", errors.New("attribute 'cluster_id' is not set in the resource")
		}
		return fmt.Sprintf("%s/%s", clusterID, rs.Primary.ID), nil
	}
}
