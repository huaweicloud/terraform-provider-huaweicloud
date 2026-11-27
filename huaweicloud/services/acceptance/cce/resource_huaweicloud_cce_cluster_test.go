package cce

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccCluster_basic(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Shanghai"),
					resource.TestCheckResourceAttr(resourceName, "agency_name", "CCEServiceAgency"),
					resource.TestCheckResourceAttr(resourceName, "description", "the description of cluster"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate_users.0.client_certificate_data",
					"certificate_users.0.client_key_data",
					"kube_config_raw",
				},
			},
			{
				Config: testAccCluster_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "agency_name", "CCEAutoClusterAgency"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
				),
			},
		},
	})
}

func TestAccCluster_prePaid(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_prePaid(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccCluster_prePaid(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func TestAccCluster_withEip(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_withEip(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
				),
			},
			{
				Config: testAccCluster_withEipUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"eip", "certificate_users.0.client_key_data", "certificate_users.0.client_certificate_data", "kube_config_raw",
				},
			},
		},
	})
}

func TestAccCluster_withEpsId(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testAccCluster_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccCluster_turbo(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_turbo(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "eni"),
					resource.TestCheckOutput("is_eni_subnet_id_different", "false"),
				),
			},
			{
				Config: testAccCluster_turbo(rName, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "eni"),
					resource.TestCheckOutput("is_eni_subnet_id_different", "false"),
				),
			},
		},
	})
}

// clusterDataPlaneV2Annotation is the cluster annotation that CCE writes only when the data plane v2 of the
// cluster is enabled. The cluster GET response carries neither the configurationsOverride field nor a
// first-class data plane v2 field, so this annotation is the only way to verify the effective state.
const clusterDataPlaneV2Annotation = "cluster.install.addons/yangtse-cilium-install-success"

func testAccCheckClusterDataPlaneV2Enabled(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		cceClient, err := cfg.CceV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CCE v3 client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		if _, ok := found.Metadata.Annotations[clusterDataPlaneV2Annotation]; !ok {
			return fmt.Errorf("the cluster %s does not carry the data plane v2 annotation %q",
				rs.Primary.ID, clusterDataPlaneV2Annotation)
		}
		return nil
	}
}

// testAccCheckClusterRecordID stores the cluster ID of the current step, so that a later step can assert that the
// cluster was destroyed and recreated instead of updated in place.
func testAccCheckClusterRecordID(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}
		*id = rs.Primary.ID
		return nil
	}
}

// testAccCheckClusterRecreated asserts that the cluster ID changed after an attribute that is ForceNew was
// modified, which is the observable evidence of a replacement: the previous cluster is destroyed and a brand-new
// one is created.
//
// This assertion is what makes the case sensitive to the replacement. A cluster updated in place keeps its ID, and
// an attribute that loses ForceNew while the update function ignores it keeps it as well (no request is sent, the
// new value still lands in the state), so an assertion on the attribute value alone cannot tell the difference.
func testAccCheckClusterRecreated(n string, previousID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}
		if *previousID == "" {
			return fmt.Errorf("the previous cluster ID of %s was not recorded", n)
		}
		if rs.Primary.ID == *previousID {
			return fmt.Errorf("the cluster %s was not replaced, its ID did not change after the data_plane_v2 "+
				"switch", rs.Primary.ID)
		}
		return nil
	}
}

func TestAccCluster_dataPlaneV2Turbo(t *testing.T) {
	var cluster clusters.Clusters
	var previousClusterID string

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				// A-02: an eni cluster created through the first-class attribute runs the data plane v2.
				// The cluster version must satisfy the server-side gate (>= 1.34.3-r10).
				// Cost: this step creates a dedicated cluster, and the next step recreates it, so the whole test
				// consumes two cluster creations.
				Config: testAccCluster_dataPlaneV2Turbo(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					testAccCheckClusterRecordID(resourceName, &previousClusterID),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "eni"),
					resource.TestCheckResourceAttr(resourceName, "data_plane_v2", "true"),
					testAccCheckClusterDataPlaneV2Enabled(resourceName),
				),
			},
			{
				// A-03: the attribute is ForceNew, switching it destroys and recreates the cluster instead of
				// updating the data plane in place. The ID change is asserted first, so that a regression of the
				// ForceNew behaviour fails here instead of being masked by the attribute value landing in the
				// state anyway (the update function does not handle the attribute at all).
				Config: testAccCluster_dataPlaneV2Turbo(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterRecreated(resourceName, &previousClusterID),
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "data_plane_v2", "false"),
				),
			},
		},
	})
}

func TestAccCluster_dataPlaneV2VpcRouter(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				// A-01: a vpc-router cluster created through the first-class attribute runs the data plane v2.
				// The cluster version must satisfy the commercial vpc-router gate.
				// Cost: this step creates a dedicated cluster (one cluster creation).
				Config: testAccCluster_dataPlaneV2VpcRouter(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "vpc-router"),
					resource.TestCheckResourceAttr(resourceName, "data_plane_v2", "true"),
					testAccCheckClusterDataPlaneV2Enabled(resourceName),
				),
			},
		},
	})
}

func TestAccCluster_dataPlaneV2UnsupportedNetwork(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				// A-05: the provider does not duplicate the server-side network mode check, the server error is
				// surfaced as is.
				Config:      testAccCluster_dataPlaneV2OverlayL2(rName),
				ExpectError: regexp.MustCompile("does not support DataPlane V2"),
			},
		},
	})
}

func TestAccCluster_componentConfigurationsCompat(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				// A-04: the users that declare the component configurations only keep the previous behavior,
				// the new attribute is not configured at all.
				// The framework asserts an empty plan right after the apply of every step, so a request body that
				// gained a key, or a state that drifted away from the configuration, fails this step as a
				// non-empty plan — the "no diff, behaviour unchanged" evidence of this case.
				// Input difference from delta-design.md §6.2: the design sketched a kube-apiserver package, the
				// eni package is used here instead, so that the end-to-end path of U-02 (the switch declared by
				// hand as a JSON boolean, left untouched by the provider) is covered as well.
				// Cost: this step creates a dedicated cluster (one cluster creation).
				Config: testAccCluster_componentConfigurations(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "eni"),
					resource.TestCheckResourceAttr(resourceName, "data_plane_v2", "false"),
					resource.TestCheckResourceAttr(resourceName, "component_configurations.#", "1"),
				),
			},
		},
	})
}

func TestAccCluster_hibernate(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
				),
			},
			{
				Config: testAccCluster_hibernate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Hibernation"),
				),
			},
			{
				Config: testAccCluster_awake(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
				),
			},
		},
	})
}

func TestAccCluster_multiContainerNetworkCidrs(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_multiContainerNetworkCidrs(rName, "172.16.0.0/24,172.16.1.0/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "vpc-router"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "container_network_cidr", "172.16.0.0/24,172.16.1.0/24"),
				),
			},
			{
				Config: testAccCluster_multiContainerNetworkCidrs(rName, "172.16.0.0/24,172.16.1.0/24,172.16.2.0/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "vpc-router"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "container_network_cidr", "172.16.0.0/24,172.16.1.0/24,172.16.2.0/24"),
				),
			},
		},
	})
}

func TestAccCluster_secGroup(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_secGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test1", "id"),
				),
			},
			{
				Config: testAccCluster_secGroup_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test2", "id"),
				),
			},
		},
	})
}

func TestAccCluster_resize(t *testing.T) {
	var cluster clusters.Clusters

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_resize(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
				),
			},
			{
				Config: testAccCluster_resize_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.medium"),
				),
			},
		},
	})
}

func TestAccCluster_resizePeriod(t *testing.T) {
	var cluster clusters.Clusters

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_resizePeriod(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
				),
			},
			{
				Config: testAccCluster_resize_updatePeriod(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.medium"),
				),
			},
		},
	})
}

func testAccCheckClusterDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	cceClient, err := cfg.CceV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating CCE v3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cce_cluster" {
			continue
		}

		_, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("cluster still exists")
		}
	}

	return nil
}

func testAccCheckClusterExists(n string, cluster *clusters.Clusters) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		cceClient, err := cfg.CceV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CCE v3 client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func testAccCluster_prePaid(rName string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "%[3]v"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName, isAutoRenew)
}

func testAccCluster_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  timezone               = "Asia/Shanghai"
  description            = "the description of cluster"
  agency_name            = "CCEServiceAgency"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  timezone               = "Asia/Shanghai"
  description            = ""
  agency_name            = "CCEAutoClusterAgency"

  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_withEip(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  eip                    = huaweicloud_vpc_eip.test.address
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_withEipUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_vpc_eip" "update" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  eip                    = huaweicloud_vpc_eip.update.address
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  enterprise_project_id  = "%s"
}

`, common.TestVpc(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCluster_turbo(rName string, eniNum int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_subnet" "eni_test" {
  count      = %[3]d

  name       = "%[2]s-eni-${count.index}"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index + 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index + 1), 1)
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "eni"
  eni_subnet_id          = join(",", huaweicloud_vpc_subnet.eni_test[*].ipv4_subnet_id)
}

output "is_eni_subnet_id_different" {
  value = length(setsubtract(split(",", huaweicloud_cce_cluster.test.eni_subnet_id),
  huaweicloud_vpc_subnet.eni_test[*].ipv4_subnet_id)) != 0
}
`, common.TestVpc(rName), rName, eniNum)
}

func testAccCluster_dataPlaneV2Turbo(rName string, dataPlaneV2 bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_subnet" "eni_test" {
  name       = "%[2]s-eni"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 1), 1)
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "eni"
  eni_subnet_id          = huaweicloud_vpc_subnet.eni_test.ipv4_subnet_id
  data_plane_v2          = %[3]v
}
`, common.TestVpc(rName), rName, dataPlaneV2)
}

func testAccCluster_dataPlaneV2VpcRouter(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "vpc-router"
  data_plane_v2          = true
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_dataPlaneV2OverlayL2(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  data_plane_v2          = true
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_componentConfigurations(rName string) string {
	// The eni package carries the switch as a JSON boolean, which is the type the server asserts on.
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "eni"

  component_configurations {
    name           = "eni"
    configurations = jsonencode([
      {
        name  = "dataplane-v2"
        value = true
      }
    ])
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_hibernate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  hibernate              = true
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_awake(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  hibernate              = false
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_multiContainerNetworkCidrs(rName, containerNetworkCidr string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "vpc-router"
  container_network_cidr = "%s"
  service_network_cidr   = "10.248.0.0/16"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName, containerNetworkCidr)
}

func testAccCluster_secGroup(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup" "test1" {
  name = "%[2]s-secgroup-1"
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform"
  security_group_id      = huaweicloud_networking_secgroup.test1.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_secGroup_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup" "test1" {
  name = "%[2]s-secgroup-1"
}

resource "huaweicloud_networking_secgroup" "test2" {
  name = "%[2]s-secgroup-2"
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform update"
  security_group_id      = huaweicloud_networking_secgroup.test2.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_resize(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_resize_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.medium"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform update"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_resizePeriod(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_resize_updatePeriod(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.medium"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform update"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}
