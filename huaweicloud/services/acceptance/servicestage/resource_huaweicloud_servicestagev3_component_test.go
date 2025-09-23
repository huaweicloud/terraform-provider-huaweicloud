package servicestage

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/servicestage"
)

func getV3ComponentFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("servicestage", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage client: %s", err)
	}
	return servicestage.QueryV3Component(client, state.Primary.Attributes["application_id"], state.Primary.ID)
}

func TestAccV3Component_basic(t *testing.T) {
	var (
		component interface{}

		resourceName = "huaweicloud_servicestagev3_component.test"
		rc           = acceptance.InitResourceCheck(resourceName, &component, getV3ComponentFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Make sure the networks of the CCE cluster, CSE engine and ELB loadbalancer are same.
			acceptance.TestAccPreCheckCceClusterId(t) // Make sure at least one of node exist.
			acceptance.TestAccPreCheckCSEMicroserviceEngineID(t)
			acceptance.TestAccPreCheckElbLoadbalancerID(t)
			// Two different JAR packages need to be provided.
			acceptance.TestAccPreCheckServiceStageJarPkgStorageURLs(t, 2)
			acceptance.TestAccPreCheckCertificateBase(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3Component_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "application_id", "huaweicloud_servicestagev3_application.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id", "huaweicloud_servicestagev3_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "limit_cpu", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "limit_memory", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "request_cpu", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "request_memory", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "replica", "2"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Shanghai"),
					resource.TestCheckResourceAttrSet(resourceName, "build"),
					resource.TestCheckResourceAttrSet(resourceName, "source"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.deploy_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.type"),
					resource.TestCheckResourceAttr(resourceName, "refer_resources.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "external_accesses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_accesses.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "external_accesses.0.forward_port", "8000"),
					resource.TestCheckResourceAttrPair(resourceName, "external_accesses.0.address",
						"huaweicloud_elb_certificate.test", "domain"),
					resource.TestCheckResourceAttr(resourceName, "envs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "envs.0.name", "MOCK_ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "envs.0.value", "true"),
					resource.TestCheckResourceAttr(resourceName, "post_start.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.command.0", "test"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.type", "command"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.command.0", "test"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.type", "command"),
					resource.TestCheckResourceAttr(resourceName, "mesher.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "mesher.0.port", "60"),
					resource.TestCheckResourceAttr(resourceName, "storages.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.type", "HostPath"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.name", name),
					resource.TestCheckResourceAttrSet(resourceName, "storages.0.parameters"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.path", "/category"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.sub_path", "sub"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.read_only", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "command"),
					resource.TestCheckResourceAttr(resourceName, "logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.log_path", "/tmp"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.rotate", "Hourly"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.host_path", "/tmp"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.host_extend_path", "PodName"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.path", "/tmp"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.port", "600"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.dimensions", "cpu_usage,mem_usage"),
					resource.TestCheckResourceAttr(resourceName, "affinity.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "anti_affinity.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.type", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.delay", "30"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.port", "800"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.type", "http"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.delay", "30"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.scheme", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.host", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.port", "8000"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.path", "/v1/test"),
					resource.TestCheckResourceAttrSet(resourceName, "update_strategy"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV3Component_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "application_id", "huaweicloud_servicestagev3_application.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id", "huaweicloud_servicestagev3_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "limit_cpu", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "limit_memory", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "request_cpu", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "request_memory", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "replica", "2"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Shanghai"),
					resource.TestCheckResourceAttrSet(resourceName, "build"),
					resource.TestCheckResourceAttrSet(resourceName, "source"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.deploy_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.type"),
					resource.TestCheckResourceAttr(resourceName, "refer_resources.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(resourceName, "tags.new_key", "value"),
					resource.TestCheckResourceAttr(resourceName, "external_accesses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_accesses.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "external_accesses.0.forward_port", "8080"),
					resource.TestCheckResourceAttrPair(resourceName, "external_accesses.0.address",
						"huaweicloud_elb_certificate.test", "domain"),
					resource.TestCheckResourceAttr(resourceName, "envs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "envs.0.name", "MOCK_ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "envs.0.value", "false"),
					resource.TestCheckResourceAttr(resourceName, "post_start.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.command.0", "new test command"),
					resource.TestCheckResourceAttr(resourceName, "post_start.0.type", "command"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.command.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.command.0", "new test command"),
					resource.TestCheckResourceAttr(resourceName, "pre_stop.0.type", "command"),
					resource.TestCheckResourceAttr(resourceName, "mesher.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "mesher.0.port", "100"),
					resource.TestCheckResourceAttr(resourceName, "storages.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.type", "HostPath"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.name", name+"-new"),
					resource.TestCheckResourceAttrSet(resourceName, "storages.0.parameters"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.path", "/category/new"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.sub_path", "sub/new"),
					resource.TestCheckResourceAttr(resourceName, "storages.0.mounts.0.read_only", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "command"),
					resource.TestCheckResourceAttr(resourceName, "logs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.log_path", "/tmp/new"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.rotate", "Daily"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.host_path", "/tmp/new"),
					resource.TestCheckResourceAttr(resourceName, "logs.0.host_extend_path", "PodUID"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.path", "/tmp/new"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.port", "800"),
					resource.TestCheckResourceAttr(resourceName, "custom_metric.0.dimensions", "mem_usage"),
					resource.TestCheckResourceAttr(resourceName, "affinity.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "anti_affinity.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.type", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.delay", "60"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.timeout", "60"),
					resource.TestCheckResourceAttr(resourceName, "liveness_probe.0.port", "900"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.type", "http"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.delay", "60"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.timeout", "60"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.scheme", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.host", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "readiness_probe.0.path", "/v1/test/new"),
					resource.TestCheckResourceAttr(resourceName, "deploy_strategy.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "update_strategy"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3ComponentImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"tags",
					"source_origin",
					"build_origin",
					"deploy_strategy.0.rolling_release_origin",
					"command_origin",
					"tomcat_opts_origin",
					"update_strategy_origin",
				},
			},
		},
	})
}

func testAccV3ComponentImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var applicationId, resourceId string
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of compnent is not found in the tfstate", resourceName)
		}
		applicationId = rs.Primary.Attributes["application_id"]
		resourceId = rs.Primary.ID
		if applicationId == "" || resourceId == "" {
			return "", fmt.Errorf("the component ID is not exist or application ID is missing")
		}
		return fmt.Sprintf("%s/%s", applicationId, resourceId), nil
	}
}

func testAccV3Component_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cce_clusters" "test" {
  cluster_id = "%[1]s"
}

data "huaweicloud_cse_microservice_engines" "test" {}

data "huaweicloud_elb_loadbalancers" "test" {
  loadbalancer_id = "%[2]s"
}

data "huaweicloud_servicestagev3_runtime_stacks" "test" {}

locals {
  java_runtime_stack = try([for o in data.huaweicloud_servicestagev3_runtime_stacks.test.runtime_stacks:
    o if o.type == "Java" && o.deploy_mode == "container"][0], {})
}

resource "huaweicloud_servicestagev3_application" "test" {
  name                  = "%[3]s"
  enterprise_project_id = "0"
}

resource "huaweicloud_servicestagev3_environment" "test" {
  name                  = "%[3]s"
  vpc_id                = try(data.huaweicloud_cce_clusters.test.clusters[0].vpc_id, "")
  enterprise_project_id = "0"
}

resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = huaweicloud_servicestagev3_environment.test.id

  resources {
    id   = try(data.huaweicloud_cce_clusters.test.clusters[0].id, "")
    type = "cce"
  }
  resources {
    id   = "%[4]s"
    type = "cse"
  }
  resources {
    id   = try(data.huaweicloud_elb_loadbalancers.test.loadbalancers[0].id, "")
    type = "elb"
  }
}

resource "huaweicloud_elb_certificate" "test" {
  name        = "test-server"
  domain      = "p2cserver.com"
  type        = "server"
  private_key = "%[5]s"
  certificate = "%[6]s"
}
`, acceptance.HW_CCE_CLUSTER_ID,
		acceptance.HW_ELB_LOADBALANCER_ID,
		name,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ID,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY,
		acceptance.HW_CERTIFICATE_CONTENT,
	)
}

func testAccV3Component_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_component" "test" {
  application_id = huaweicloud_servicestagev3_application.test.id
  environment_id = huaweicloud_servicestagev3_environment.test.id
  name           = "%[2]s"
  version        = "1.0.1"
  description    = "Created by terraform script"
  limit_cpu      = 0.25
  limit_memory   = 0.5
  request_cpu    = 0.25
  request_memory = 0.5
  replica        = 2
  timezone       = "Asia/Shanghai"

  build = jsonencode({
    "parameters" : {
      "cluster_id": try(data.huaweicloud_cce_clusters.test.clusters[0].id, ""),
      "dockerfile_path": "./",
      "build_env_selected": "current"
    }
  })

  source = jsonencode({
    "kind": "package",
    "storage": "obs",
    "url": try(element(split(",", "%[3]s"), 0), "")
  })

  runtime_stack {
    deploy_mode = try(local.java_runtime_stack.deploy_mode, "container")
    name        = try(local.java_runtime_stack.name, "OpenJDK17")
    type        = try(local.java_runtime_stack.type, "Java")
    version     = try(local.java_runtime_stack.version, null)
  }

  refer_resources {
    id   = try(data.huaweicloud_cce_clusters.test.clusters[0].id, "")
    type = "cce"
    parameters = jsonencode({
      "namespace": "default",
      "type": "VirtualMachine"
      "name": format("%%s-first-version", try(data.huaweicloud_cce_clusters.test.clusters[0].name, ""))
    })
  }
  refer_resources {
    id   = try(data.huaweicloud_elb_loadbalancers.test.loadbalancers[0].id, "")
    type = "elb"
    parameters = jsonencode({
      "name": try(data.huaweicloud_elb_loadbalancers.test.loadbalancers[0].name, "")
    })
  }
  refer_resources {
    id   = "%[4]s"
    type = "cse"
  }

  tags = {
    foo = "bar"
    key = "value"
  }

  external_accesses {
    protocol     = "HTTP"
    forward_port = "8000"
    address      = huaweicloud_elb_certificate.test.domain
  }

  envs {
    name  = "MOCK_ENABLED"
    value = "true"
  }

  post_start {
    command = ["test"]
    type    = "command"
  }

  pre_stop {
    command = ["test"]
    type    = "command"
  }

  mesher {
    port = 60
  }

  storages {
    type       = "HostPath"
    name       = "%[2]s"
    parameters = jsonencode({
      "default_mode": 0,
      "path": "/tmp"
    })
    mounts {
      path      = "/category"
      sub_path  = "sub"
      read_only = true
    }
  }

  command = jsonencode({
    "args": ["-a"],
    "command": ["ls"]
  })

  logs {
    log_path         = "/tmp"
    rotate           = "Hourly"
    host_path        = "/tmp"
    host_extend_path = "PodName"
  }

  custom_metric {
    path       = "/tmp"
    port       = 600
    dimensions = "cpu_usage,mem_usage"
  }

  affinity {
    condition = "required"
    kind      = "node"
    match_expressions {
      key       = "affinity1"
      value     = "foo"
      operation = "In"
    }
    weight = 100
  }
  affinity {
    condition = "preferred"
    kind      = "node"
    match_expressions {
      key       = "affinity2"
      value     = "bar"
      operation = "NotIn"
    }
    weight = 1
  }

  anti_affinity {
    condition = "required"
    kind      = "pod"
    match_expressions {
      key       = "anit-affinity1"
      operation = "Exists"
    }
    weight = 100
  }
  anti_affinity {
    condition = "preferred"
    kind      = "pod"
    match_expressions {
      key       = "anti-affinity2"
      operation = "DoesNotExist"
    }
    weight = 1
  }

  liveness_probe {
    type    = "tcp"
    delay   = 30
    timeout = 30
    port    = 800
  }

  readiness_probe {
    type    = "http"
    delay   = 30
    timeout = 30
    scheme  = "HTTPS"
    host    = "127.0.0.1"
    port    = 8000
    path    = "/v1/test"
  }

  update_strategy = jsonencode({
    "type": "RollingUpdate",
    "min_ready_seconds": 20,
    "revision_history_limit": 15,
    "progress_deadline_seconds": 900,
    "termination_period_seconds": 30,
    "max_unavailable": "30%%",
    "max_surge": "30%%"
  })

  depends_on = [
    huaweicloud_elb_certificate.test,
    huaweicloud_servicestagev3_environment_associate.test,
  ]
}
`, testAccV3Component_base(name),
		name,
		acceptance.HW_SERVICESTAGE_JAR_PKG_STORAGE_URLS,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ID)
}

func testAccV3Component_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_component" "test" {
  application_id = huaweicloud_servicestagev3_application.test.id
  environment_id = huaweicloud_servicestagev3_environment.test.id
  name           = "%[2]s"
  version        = "1.0.2"
  description    = "Updated by terraform script"
  limit_cpu      = 0.25
  limit_memory   = 0.5
  request_cpu    = 0.25
  request_memory = 0.5
  replica        = 2
  timezone       = "Asia/Shanghai"

  build = jsonencode({
    "parameters": {
      "environment_id": huaweicloud_servicestagev3_environment.test.id,
      "cluster_namespace": "default",
      "use_public_cluster": false,
      "cluster_id": try(data.huaweicloud_cce_clusters.test.clusters[0].id, ""),
      "dockerfile_path": "./",
      "build_env_selected": "current"
    }
  })

  source = jsonencode({
    "kind": "package",
    "storage": "obs"
    "url": try(element(split(",", "%[3]s"), 1), "")
  })

  runtime_stack {
    deploy_mode = try(local.java_runtime_stack.deploy_mode, "container")
    name        = try(local.java_runtime_stack.name, "OpenJDK17")
    type        = try(local.java_runtime_stack.type, "Java")
    version     = try(local.java_runtime_stack.version, null)
  }

  refer_resources {
    id         = try(data.huaweicloud_cce_clusters.test.clusters[0].id, "")
    type       = "cce"
    parameters = jsonencode({
      "namespace": "default",
      "type": "VirtualMachine"
      "name": try(data.huaweicloud_cce_clusters.test.clusters[0].name, "")
    })
  }
  refer_resources {
    id         = try(data.huaweicloud_elb_loadbalancers.test.loadbalancers[0].id, "")
    type       = "elb"
    parameters = jsonencode({
      "name": try(data.huaweicloud_elb_loadbalancers.test.loadbalancers[0].name, "")
    })
  }
  refer_resources {
    id   = "%[4]s"
    type = "cse"
  }

  tags = {
    foo     = "baar"
    new_key = "value"
  }

  external_accesses {
    protocol     = "HTTP"
    forward_port = "8080"
    address      = huaweicloud_elb_certificate.test.domain
  }

  envs {
    name  = "MOCK_ENABLED"
    value = "false"
  }

  post_start {
    command = ["new test command"]
    type    = "command"
  }

  pre_stop {
    command = ["new test command"]
    type    = "command"
  }

  mesher {
    port = 100
  }

  storages {
    type       = "HostPath"
    name       = "%[2]s-new"
    parameters = jsonencode({
      "default_mode": 0,
      "path": "/tmp/new"
    })
    mounts {
      path      = "/category/new"
      sub_path  = "sub/new"
      read_only = true
    }
  }

  command = jsonencode({
    "args": ["-l"],
    "command": ["ls"]
  })

  logs {
    log_path         = "/tmp/new"
    rotate           = "Daily"
    host_path        = "/tmp/new"
    host_extend_path = "PodUID"
  }

  custom_metric {
    path       = "/tmp/new"
    port       = 800
    dimensions = "mem_usage"
  }

  affinity {
    condition = "required"
    kind      = "node"
    match_expressions {
      key       = "new_affinity1"
      value     = "1"
      operation = "Gt"
    }
    weight = 100
  }
  affinity {
    condition = "preferred"
    kind      = "node"
    match_expressions {
      key       = "new_affinity2"
      value     = "100"
      operation = "Lt"
    }
    weight = 1
  }

  anti_affinity {
    condition = "required"
    kind      = "pod"
    match_expressions {
      key       = "new_anit-affinity1"
      operation = "Exists"
    }
    weight = 100
  }
  anti_affinity {
    condition = "preferred"
    kind      = "pod"
    match_expressions {
      key       = "new_anti-affinity2"
      operation = "DoesNotExist"
    }
    weight = 1
  }

  liveness_probe {
    type    = "tcp"
    delay   = 60
    timeout = 60
    port    = 900
  }

  readiness_probe {
    type    = "http"
    delay   = 60
    timeout = 60
    scheme  = "HTTP"
    host    = "192.168.0.1"
    port    = 8080
    path    = "/v1/test/new"
  }

  deploy_strategy {
    type            = "RollingRelease"
    rolling_release = jsonencode({
      "batches": 1
    })
  }
  update_strategy = jsonencode({
    "max_surge": "30%%",
    "max_unavailable": 2,
    "progress_deadline_seconds": 900,
    "revision_history_limit": 20,
    "termination_period_seconds": 30,
    "type": "RollingUpdate"
  })

  depends_on = [
    huaweicloud_elb_certificate.test,
    huaweicloud_servicestagev3_environment_associate.test,
  ]
}
`, testAccV3Component_base(name),
		name,
		acceptance.HW_SERVICESTAGE_JAR_PKG_STORAGE_URLS,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ID)
}

func TestAccV3Component_yaml(t *testing.T) {
	var (
		component interface{}

		resourceName = "huaweicloud_servicestagev3_component.test"
		rc           = acceptance.InitResourceCheck(resourceName, &component, getV3ComponentFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Make sure at least one of node exist.
			acceptance.TestAccPreCheckCceClusterId(t)
			// At least one of JAR package must be provided.
			acceptance.TestAccPreCheckServiceStageJarPkgStorageURLs(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3Component_yaml_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "application_id", "huaweicloud_servicestagev3_application.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id", "huaweicloud_servicestagev3_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttr(resourceName, "config_mode", "yaml"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "limit_cpu", "0"),
					resource.TestCheckResourceAttr(resourceName, "limit_memory", "0"),
					resource.TestCheckResourceAttr(resourceName, "request_cpu", "0"),
					resource.TestCheckResourceAttr(resourceName, "request_memory", "0"),
					resource.TestCheckResourceAttr(resourceName, "replica", "2"),
					resource.TestCheckResourceAttr(resourceName, "build", ""),
					resource.TestCheckResourceAttrSet(resourceName, "source"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.deploy_mode", "container"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.name", "Docker"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.type", "Docker"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "refer_resources.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "update_strategy"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3ComponentImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"workload_content",
					"tags",
					"source_origin",
					"build_origin",
					"deploy_strategy.0.rolling_release_origin",
					"command_origin",
					"tomcat_opts_origin",
					"update_strategy_origin",
				},
			},
		},
	})
}

func testAccV3Component_yaml_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cce_clusters" "test" {
  cluster_id = "%[1]s"
}

data "huaweicloud_servicestagev3_runtime_stacks" "test" {}

locals {
  docker_runtime_stack = try([for o in data.huaweicloud_servicestagev3_runtime_stacks.test.runtime_stacks:
    o if o.type == "Docker" && o.deploy_mode == "container"][0], {})
}

resource "huaweicloud_servicestagev3_application" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "0"
}

resource "huaweicloud_servicestagev3_environment" "test" {
  name                  = "%[2]s"
  vpc_id                = try(data.huaweicloud_cce_clusters.test.clusters[0].vpc_id, "")
  enterprise_project_id = "0"
}

resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = huaweicloud_servicestagev3_environment.test.id

  resources {
    id   = try(data.huaweicloud_cce_clusters.test.clusters[0].id, "")
    type = "cce"
  }
}
`, acceptance.HW_CCE_CLUSTER_ID, name)
}

func testAccV3Component_yaml_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_component" "test" {
  name           = "%[2]s"
  description    = "Created by terraform script"
  version        = "1.0.1"
  environment_id = huaweicloud_servicestagev3_environment.test.id
  application_id = huaweicloud_servicestagev3_application.test.id

  runtime_stack {
    deploy_mode = try(local.docker_runtime_stack.deploy_mode, "container")
    name        = try(local.docker_runtime_stack.name, "Docker")
    type        = try(local.docker_runtime_stack.type, "Docker")
    version     = try(local.docker_runtime_stack.version, "1.0")
  }

  source = jsonencode({
    kind    = "image"
    storage = "swr"
    url     = try(element(split(",", "%[3]s"), 0), "")
  })

  refer_resources {
    type       = "cce"
    id         = try(data.huaweicloud_cce_clusters.test.clusters[0].id, "")
    parameters = jsonencode({
      type      = "VirtualMachine"
      namespace = "default"
    })
  }

  config_mode      = "yaml"
  workload_content = jsonencode({
    apiVersion = "apps/v1"
    kind       = "Deployment"
    metadata   = {
      name      = "%[2]s"
      namespace = "default"
    }
    spec = {
      selector = {}
      template = {
        metadata = {}
        spec = {
          imagePullSecrets = [
            {
              name = "default-secret",
            }
          ]
          terminationGracePeriodSeconds = 30
          volumes                       = []
          restartPolicy                 = "Always"
          dnsPolicy                     = "ClusterFirst"
          containers                    = [
            {
              image           = try(element(split(",", "%[3]s"), 0), "")
              name            = "%[2]s"
              imagePullPolicy = "Always"
              resources       = {
                requests = {
                  cpu    = "0"
                  memory = "0"
                }
                limits = {
                  cpu    = "0"
                  memory = "0"
                }
              }
              ports = [
                {
                  containerPort = 8080,
                  protocol      = "TCP"
                }
              ]
            }
          ]
        }
      }
    }
    strategy = {
      type = "RollingUpdate"
      rollingUpdate = {
        maxSurge       = 0
        maxUnavailable = 1
      }
    }
    replicas = 2
  })
}
`, testAccV3Component_yaml_base(name), name, acceptance.HW_SERVICESTAGE_JAR_PKG_STORAGE_URLS)
}
