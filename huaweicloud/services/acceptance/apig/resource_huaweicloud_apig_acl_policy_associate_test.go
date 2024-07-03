package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/acls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getAclPolicyAssociateFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	opt := acls.ListBindOpts{
		InstanceId: state.Primary.Attributes["instance_id"],
		PolicyId:   state.Primary.Attributes["policy_id"],
	}
	resp, err := acls.ListBind(c, opt)
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, err
}

func TestAccAclPolicyAssociate_basic(t *testing.T) {
	var (
		apiDetails []acls.AclBindApiInfo

		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_apig_acl_policy_associate.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&apiDetails,
		getAclPolicyAssociateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckApigChannelRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAclPolicyAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_apig_acl_policy.test", "id"),
					resource.TestCheckResourceAttr(rName, "publish_ids.#", "1"),
				),
			},
			{
				Config: testAccAclPolicyAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_apig_acl_policy.test", "id"),
					resource.TestCheckResourceAttr(rName, "publish_ids.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAclPolicyAssociateImportStateFunc(rName),
			},
		},
	})
}

func testAccAclPolicyAssociateImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["policy_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<policy_id>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["policy_id"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["policy_id"]), nil
	}
}

func testAccAclPolicyAssociate_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[3]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = "%[4]s"
  }
}

resource "huaweicloud_apig_group" "test" {
  instance_id = local.instance_id
  name        = "%[3]s"
}

resource "huaweicloud_apig_channel" "test" {
  instance_id      = local.instance_id
  name             = "%[3]s"
  port             = 8000
  balance_strategy = 2
  member_type      = "ecs"
  type             = 2

  health_check {
    protocol           = "HTTPS"
    threshold_normal   = 10  # maximum value
    threshold_abnormal = 10  # maximum value
    interval           = 300 # maximum value
    timeout            = 30  # maximum value
    path               = "/"
    method             = "HEAD"
    port               = 8080
    http_codes         = "201,202,303-404"
  }

  member {
    id   = huaweicloud_compute_instance.test.id
    name = huaweicloud_compute_instance.test.name
  }
}

resource "huaweicloud_apig_api" "test" {
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[3]s"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/user_info/{user_age}"
  security_authentication = "APP"
  matching                = "Exact"
  success_response        = "Success response"
  failure_response        = "Failed response"
  description             = "Created by script"

  request_params {
    name     = "user_age"
    type     = "NUMBER"
    location = "PATH"
    required = true
    maximum  = 200
    minimum  = 0
  }
  
  backend_params {
    type     = "REQUEST"
    name     = "userAge"
    location = "PATH"
    value    = "user_age"
  }

  web {
    path             = "/getUserAge/{userAge}"
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
  }

  web_policy {
    name             = "%[3]s_web"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    path             = "/getUserAge/{userAge}"
    timeout          = 30000
    vpc_channel_id   = huaweicloud_apig_channel.test.id

    backend_params {
      type     = "REQUEST"
      name     = "userAge"
      location = "PATH"
      value    = "user_age"
    }

    conditions {
      source     = "param"
      param_name = "user_age"
      type       = "Equal"
      value      = "28"
    }
  }
}

resource "huaweicloud_apig_environment" "test" {
  count = 2

  instance_id = local.instance_id
  name        = "%[3]s_${count.index}"
}

resource "huaweicloud_apig_api_publishment" "test" {
  count = 2

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test[count.index].id
}

resource "huaweicloud_apig_acl_policy" "test" {
  instance_id = local.instance_id
  name        = "%[3]s"
  type        = "PERMIT"
  entity_type = "IP"
  value       = "10.201.33.4,10.30.2.15"
}
`, common.TestBaseComputeResources(name),
		acceptance.HW_APIG_DEDICATED_INSTANCE_ID,
		name,
		acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID)
}

func testAccAclPolicyAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_acl_policy_associate" "test" {
  instance_id = local.instance_id
  policy_id   = huaweicloud_apig_acl_policy.test.id

  publish_ids = [
    huaweicloud_apig_api_publishment.test[0].publish_id
  ]
}
`, testAccAclPolicyAssociate_base(name))
}

func testAccAclPolicyAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_acl_policy_associate" "test" {
  instance_id = local.instance_id
  policy_id   = huaweicloud_apig_acl_policy.test.id

  publish_ids = [
    huaweicloud_apig_api_publishment.test[1].publish_id
  ]
}
`, testAccAclPolicyAssociate_base(name))
}
