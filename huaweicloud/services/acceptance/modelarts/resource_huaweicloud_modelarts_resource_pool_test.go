package modelarts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getModelartsResourcePoolResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getModelartsResourcePoolHttpUrl = "v2/{project_id}/pools/{id}"
		getModelartsResourcePoolProduct = "modelarts"
	)
	getModelartsResourcePoolClient, err := cfg.NewServiceClient(getModelartsResourcePoolProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	getModelartsResourcePoolPath := getModelartsResourcePoolClient.Endpoint + getModelartsResourcePoolHttpUrl
	getModelartsResourcePoolPath = strings.ReplaceAll(getModelartsResourcePoolPath, "{project_id}", getModelartsResourcePoolClient.ProjectID)
	getModelartsResourcePoolPath = strings.ReplaceAll(getModelartsResourcePoolPath, "{id}", state.Primary.ID)

	getModelartsResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getModelartsResourcePoolResp, err := getModelartsResourcePoolClient.Request("GET", getModelartsResourcePoolPath,
		&getModelartsResourcePoolOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Modelarts resource pool: %s", err)
	}

	getModelartsResourcePoolRespBody, err := utils.FlattenResponse(getModelartsResourcePoolResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Modelarts resource pool: %s", err)
	}

	return getModelartsResourcePoolRespBody, nil
}

func TestAccModelartsResourcePool_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_modelarts_resource_pool.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getModelartsResourcePoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testModelartsResourcePool_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttr(rName, "scope.#", "3"),
					resource.TestCheckResourceAttrPair(rName, "network_id",
						"huaweicloud_modelarts_network.test", "id"),
					resource.TestCheckResourceAttr(rName, "resources.0.flavor_id", "modelarts.vm.cpu.8ud"),
					resource.TestCheckResourceAttr(rName, "resources.0.count", "1"),
				),
			},
			{
				Config: testModelartsResourcePool_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo update"),
					resource.TestCheckResourceAttr(rName, "scope.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "network_id",
						"huaweicloud_modelarts_network.test", "id"),
					resource.TestCheckResourceAttr(rName, "resources.0.flavor_id", "modelarts.vm.cpu.8ud"),
					resource.TestCheckResourceAttr(rName, "resources.0.count", "1"),
					resource.TestCheckResourceAttr(rName, "resources.1.flavor_id", "modelarts.vm.cpu.16u64g.d"),
					resource.TestCheckResourceAttr(rName, "resources.1.count", "1"),
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

func TestAccModelartsResourcePool_privilege_pool(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	prefix := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_modelarts_resource_pool.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getModelartsResourcePoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelartsUserLoginPassword(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testModelartsResourcePool_privilege_pool(name, prefix),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "prefix", prefix),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"huaweicloud_cce_cluster.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"huaweicloud_cce_cluster.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.provider_id",
						"huaweicloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.name",
						"huaweicloud_cce_cluster.test", "name"),
					resource.TestCheckResourceAttr(rName, "resources.0.flavor_id", "modelarts.vm.cpu.8ud"),
					resource.TestCheckResourceAttr(rName, "resources.0.count", "1"),
					resource.TestCheckResourceAttr(rName, "resources.0.node_pool", fmt.Sprintf("%s-node-pool", name)),
					resource.TestCheckResourceAttr(rName, "resources.0.max_count", "1"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.vpc_id",
						"huaweicloud_cce_cluster.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.subnet_id",
						"huaweicloud_cce_cluster.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.security_group_ids.0",
						"huaweicloud_cce_cluster.test", "security_group_id"),
					resource.TestCheckResourceAttr(rName, "resources.0.azs.0.count", "1"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.azs.0.az",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(rName, "resources.0.labels.aaa", "value_aaa"),
					resource.TestCheckResourceAttr(rName, "resources.0.labels.bbb", "value_bbb"),
					resource.TestCheckResourceAttr(rName, "resources.0.tags.key", "terraform"),
					resource.TestCheckResourceAttr(rName, "resources.0.tags.owner", "value"),
					resource.TestCheckResourceAttr(rName, "resources.0.taints.0.key", "node_1_key_1"),
					resource.TestCheckResourceAttr(rName, "resources.0.taints.0.value", "node_1_value_1"),
					resource.TestCheckResourceAttr(rName, "resources.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(rName, "resources.0.taints.1.key", "node_1_key_2"),
					resource.TestCheckResourceAttr(rName, "resources.0.taints.1.value", "node_1_value_2"),
					resource.TestCheckResourceAttr(rName, "resources.0.taints.1.effect", "PreferNoSchedule"),
					resource.TestCheckResourceAttr(rName, "resources.0.taints.2.key", "node_1_key_3"),
					resource.TestCheckResourceAttr(rName, "resources.0.taints.2.value", "node_1_value_3"),
					resource.TestCheckResourceAttr(rName, "resources.0.taints.2.effect", "NoExecute"),
					resource.TestCheckResourceAttr(rName, "resources.0.post_install", "test_post_install"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "resource_pool_id"),
				),
			},
			{
				Config: testModelartsResourcePool_privilege_pool_update(name, prefix),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "prefix", prefix),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo update"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"huaweicloud_cce_cluster.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"huaweicloud_cce_cluster.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.provider_id",
						"huaweicloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.name",
						"huaweicloud_cce_cluster.test", "name"),
					resource.TestCheckResourceAttr(rName, "resources.0.flavor_id", "modelarts.vm.cpu.8ud"),
					resource.TestCheckResourceAttr(rName, "resources.0.count", "1"),
					resource.TestCheckResourceAttr(rName, "resources.0.node_pool", fmt.Sprintf("%s-node-pool", name)),
					resource.TestCheckResourceAttr(rName, "resources.0.max_count", "1"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.vpc_id",
						"huaweicloud_cce_cluster.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.subnet_id",
						"huaweicloud_cce_cluster.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.security_group_ids.0",
						"huaweicloud_cce_cluster.test", "security_group_id"),
					resource.TestCheckResourceAttr(rName, "resources.0.azs.0.count", "1"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.azs.0.az",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(rName, "resources.0.taints.#", "0"),

					resource.TestCheckResourceAttr(rName, "resources.1.flavor_id", "modelarts.vm.cpu.16u64g.d"),
					resource.TestCheckResourceAttr(rName, "resources.1.count", "1"),
					resource.TestCheckResourceAttr(rName, "resources.1.node_pool", fmt.Sprintf("%s-node-pool-2", name)),
					resource.TestCheckResourceAttr(rName, "resources.1.max_count", "1"),
					resource.TestCheckResourceAttrPair(rName, "resources.1.vpc_id",
						"huaweicloud_cce_cluster.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(rName, "resources.1.subnet_id",
						"huaweicloud_cce_cluster.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(rName, "resources.1.security_group_ids.0",
						"huaweicloud_cce_cluster.test", "security_group_id"),
					resource.TestCheckResourceAttr(rName, "resources.1.azs.0.count", "1"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.azs.0.az",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(rName, "resources.1.labels.aaa_2", "value_aaa_2"),
					resource.TestCheckResourceAttr(rName, "resources.1.labels.bbb_2", "value_bbb_2"),
					resource.TestCheckResourceAttr(rName, "resources.1.tags.key_2", "terraform_2"),
					resource.TestCheckResourceAttr(rName, "resources.1.tags.owner_2", "value_2"),
					resource.TestCheckResourceAttr(rName, "resources.1.taints.0.key", "node_2_key_1"),
					resource.TestCheckResourceAttr(rName, "resources.1.taints.0.value", "node_2_value_1"),
					resource.TestCheckResourceAttr(rName, "resources.1.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(rName, "resources.1.taints.1.key", "node_2_key_2"),
					resource.TestCheckResourceAttr(rName, "resources.1.taints.1.value", "node_2_value_2"),
					resource.TestCheckResourceAttr(rName, "resources.1.taints.1.effect", "PreferNoSchedule"),
					resource.TestCheckResourceAttr(rName, "resources.1.taints.2.key", "node_2_key_3"),
					resource.TestCheckResourceAttr(rName, "resources.1.taints.2.value", "node_2_value_3"),
					resource.TestCheckResourceAttr(rName, "resources.1.taints.2.effect", "NoExecute"),
					resource.TestCheckResourceAttr(rName, "resources.1.post_install", "test_post_install"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period_unit", "period", "auto_renew", "user_login"},
			},
		},
	})
}

func testModelartsResourcePool_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = "%s"
  description = "This is a demo"
  scope       = ["Train", "Infer", "Notebook"]
  network_id  = huaweicloud_modelarts_network.test.id

  resources {
    flavor_id = "modelarts.vm.cpu.8ud"
    count     = 1
  }
}
`, testModelartsResourcePool_base(name), name)
}

func testModelartsResourcePool_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = "%s"
  description = "This is a demo update"
  scope       = ["Infer", "Train"]
  network_id  = huaweicloud_modelarts_network.test.id

  resources {
    flavor_id = "modelarts.vm.cpu.8ud"
    count     = 1
  }

  resources {
    flavor_id = "modelarts.vm.cpu.16u64g.d"
    count     = 1
  }
}
`, testModelartsResourcePool_base(name), name)
}

func testModelartsResourcePool_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_network" "test" {
  name = "%s"
  cidr = "172.16.0.0/12"
}`, name)
}

func testModelartsResourcePool_privilege_pool(name, prefix string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = "%[2]s"
  prefix      = "%[3]s"
  description = "This is a demo"
  vpc_id      = huaweicloud_cce_cluster.test.vpc_id
  subnet_id   = huaweicloud_cce_cluster.test.subnet_id

  clusters {
    provider_id = huaweicloud_cce_cluster.test.id
  }

  user_login {
    password = "%[4]s"
  }

  resources {
    flavor_id          = "modelarts.vm.cpu.8ud"
    count              = 1
    node_pool          = "%[2]s-node-pool"
    max_count          = "1"
    vpc_id             = huaweicloud_cce_cluster.test.vpc_id
    subnet_id          = huaweicloud_cce_cluster.test.subnet_id
    security_group_ids = [huaweicloud_cce_cluster.test.security_group_id]

    azs {
      count = 1
      az    = data.huaweicloud_availability_zones.test.names[0]
    }

    labels = {
      aaa = "value_aaa"
      bbb = "value_bbb"
    }

    tags = {
      key   = "terraform"
      owner = "value"
    }

    taints {
      key    = "node_1_key_1"
      value  = "node_1_value_1"
      effect = "NoSchedule"
    }
    taints {
      key    = "node_1_key_2"
      value  = "node_1_value_2"
      effect = "PreferNoSchedule"
    }
    taints {
      key    = "node_1_key_3"
      value  = "node_1_value_3"
      effect = "NoExecute"
    }

    post_install = "test_post_install"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true
}
`, testModelartsResourcePool_privilege_pool_base(name), name, prefix, acceptance.HW_MODELARTS_USER_LOGIN_PASSWORD)
}

func testModelartsResourcePool_privilege_pool_update(name, prefix string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = "%[2]s"
  prefix      = "%[3]s"
  description = "This is a demo update"
  vpc_id      = huaweicloud_cce_cluster.test.vpc_id
  subnet_id   = huaweicloud_cce_cluster.test.subnet_id

  clusters {
    provider_id = huaweicloud_cce_cluster.test.id
  }

   user_login {
    password = "%[4]s"
  }

  resources {
    flavor_id          = "modelarts.vm.cpu.8ud"
    count              = 1
    node_pool          = "%[2]s-node-pool"
    max_count          = "1"
    vpc_id             = huaweicloud_cce_cluster.test.vpc_id
    subnet_id          = huaweicloud_cce_cluster.test.subnet_id
    security_group_ids = [huaweicloud_cce_cluster.test.security_group_id]

    azs {
      count = 1
      az    = data.huaweicloud_availability_zones.test.names[0]
    }
  }

  resources {
    flavor_id          = "modelarts.vm.cpu.16u64g.d"
    count              = 1
    node_pool          = "%[2]s-node-pool-2"
    max_count          = "1"
    vpc_id             = huaweicloud_cce_cluster.test.vpc_id
    subnet_id          = huaweicloud_cce_cluster.test.subnet_id
    security_group_ids = [huaweicloud_cce_cluster.test.security_group_id]

    azs {
      count = 1
      az    = data.huaweicloud_availability_zones.test.names[0]
    }

    labels = {
      aaa_2 = "value_aaa_2"
      bbb_2 = "value_bbb_2"
    }

    tags = {
      key_2   = "terraform_2"
      owner_2 = "value_2"
    }

    taints {
      key    = "node_2_key_1"
      value  = "node_2_value_1"
      effect = "NoSchedule"
    }
    taints {
      key    = "node_2_key_2"
      value  = "node_2_value_2"
      effect = "PreferNoSchedule"
    }
    taints {
      key    = "node_2_key_3"
      value  = "node_2_value_3"
      effect = "NoExecute"
    }
    post_install = "test_post_install"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true
}
`, testModelartsResourcePool_privilege_pool_base(name), name, prefix, acceptance.HW_MODELARTS_USER_LOGIN_PASSWORD)
}

func testModelartsResourcePool_privilege_pool_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  cluster_version        = "v1.25"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "true"
}
`, common.TestVpc(name), name)
}
