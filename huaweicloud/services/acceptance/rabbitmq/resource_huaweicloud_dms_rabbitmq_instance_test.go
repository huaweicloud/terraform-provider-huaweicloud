package rabbitmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dms/v2/rabbitmq/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getDmsRabbitMqInstanceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID).Extract()
}

func TestAccDmsRabbitmqInstances_newFormat_cluster(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_newFormat_cluster(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "maintain_begin"),
					resource.TestCheckResourceAttrSet(resourceName, "maintain_end"),
					resource.TestCheckResourceAttr(resourceName, "disk_encrypted_enable", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "disk_encrypted_key", "huaweicloud_kms_key.test", "id"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccDmsRabbitmqInstance_newFormat_cluster_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "10:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "used_storage_space",
				},
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_newFormat_single(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_newFormat_single(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "10:00:00"),
				),
			},
			{
				Config: testAccDmsRabbitmqInstance_newFormat_single_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "02:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "06:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "used_storage_space",
				},
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_prePaid(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := rName + "update"
	resourceName := "huaweicloud_dms_rabbitmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "disk_encrypted_enable", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "disk_encrypted_key", "huaweicloud_kms_key.test", "id"),
				),
			},
			{
				Config: testAccDmsRabbitmqInstance_prePaid_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"auto_renew",
					"period",
					"period_unit",
				},
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_compatible(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_compatible(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "used_storage_space",
				},
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_publicID(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_publicID(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip_id", "huaweicloud_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "disk_encrypted_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "disk_encrypted_key", ""),
				),
			},
			{
				Config: testAccDmsRabbitmqInstance_publicID_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip_id", "huaweicloud_vpc_eip.test.1", "id"),
				),
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_amqp_single(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_amqp_single(rName, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "AMQP-0-9-1"),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "true"),
				),
			},
			{
				Config: testAccDmsRabbitmqInstance_amqp_single(rName, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "AMQP-0-9-1"),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "false"),
				),
			},
		},
	})
}

// available_zones are deprecated
func testAccDmsRabbitmqInstance_compatible(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_az" "test" {}

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type = "single"
}

locals {
  flavor = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0]
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test"
  
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  available_zones   = [data.huaweicloud_dms_az.test.id]

  flavor_id         = local.flavor.id
  engine_version    = "3.8.35"
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  access_user       = "user"
  password          = "Rabbitmqtest@123"
  broker_num        = 1

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDmsRabbitmqInstance_newFormat_cluster(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type = "cluster"
}

locals {
  flavor = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0]
}

resource "huaweicloud_kms_key" "test" {
  key_alias     = "%[2]s"
  key_algorithm = "AES_256"
  key_usage     = "ENCRYPT_DECRYPT"
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name        = "%[2]s"
  description = "rabbitmq test"
  
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id             = local.flavor.id
  engine_version        = "3.8.35"
  storage_space         = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  storage_spec_code     = local.flavor.ios[0].storage_spec_code
  broker_num            = 3
  access_user           = "user"
  password              = "Rabbitmqtest@123"
  disk_encrypted_enable = true
  disk_encrypted_key    = huaweicloud_kms_key.test.id

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDmsRabbitmqInstance_newFormat_cluster_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type = "cluster"
}

locals {
  flavor    = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0]
  newFlavor = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[1]
}

resource "huaweicloud_kms_key" "test" {
  key_alias     = "%[2]s"
  key_algorithm = "AES_256"
  key_usage     = "ENCRYPT_DECRYPT"
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name        = "%[3]s"
  description = "rabbitmq test update"
  
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id             = local.newFlavor.id
  engine_version        = "3.8.35"
  storage_space         = 1000
  storage_spec_code     = local.flavor.ios[0].storage_spec_code
  broker_num            = 5
  access_user           = "user"
  password              = "Rabbitmqtest@123"
  maintain_begin        = "06:00:00"
  maintain_end          = "10:00:00"
  disk_encrypted_enable = true
  disk_encrypted_key    = huaweicloud_kms_key.test.id

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}`, common.TestBaseNetwork(rName), rName, updateName)
}

func testAccDmsRabbitmqInstance_newFormat_single(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type = "single"
}

locals {
  flavor = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0]
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test"
  
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id         = local.flavor.id
  engine_version    = "3.8.35"
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  access_user       = "user"
  password          = "Rabbitmqtest@123"
  broker_num        = 1
  maintain_begin    = "06:00"
  maintain_end      = "10:00"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDmsRabbitmqInstance_newFormat_single_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type = "single"
}

locals {
  flavor    = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0]
  newFlavor = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[1]  
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test update"
  
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id         = local.newFlavor.id
  engine_version    = "3.8.35"
  storage_space     = 600
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  access_user       = "user"
  password          = "Rabbitmqtest@123"
  broker_num        = 1
  maintain_begin    = "02:00"
  maintain_end      = "06:00"

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}`, common.TestBaseNetwork(rName), updateName)
}

func testAccDmsRabbitmqInstance_prePaid(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type = "cluster"
}

locals {
  flavor = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0]
}

resource "huaweicloud_kms_key" "test" {
  key_alias     = "%[2]s"
  key_algorithm = "AES_256"
  key_usage     = "ENCRYPT_DECRYPT"
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name        = "%[2]s"
  description = "rabbitmq test"
  
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id         = local.flavor.id
  engine_version    = "3.8.35"
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  broker_num        = 3

  access_user           = "user"
  password              = "Rabbitmqtest@123"
  disk_encrypted_enable = true
  disk_encrypted_key    = huaweicloud_kms_key.test.id

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDmsRabbitmqInstance_prePaid_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type = "cluster"
}

locals {
  flavor = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0]
}

resource "huaweicloud_kms_key" "test" {
  key_alias     = "%[2]s"
  key_algorithm = "AES_256"
  key_usage     = "ENCRYPT_DECRYPT"
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name        = "%[3]s"
  description = "rabbitmq test update"
  
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id         = local.flavor.id
  engine_version    = "3.8.35"
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  broker_num        = 3

  access_user           = "user"
  password              = "Rabbitmqtest@123"
  disk_encrypted_enable = true
  disk_encrypted_key    = huaweicloud_kms_key.test.id

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}`, common.TestBaseNetwork(rName), rName, updateName)
}

func testEip_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  count = 2

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s_${count.index}"
    size        = 5
    charge_mode = "traffic"
  }
}
`, name)
}

func testAccDmsRabbitmqInstance_publicID(rName string) string {
	return fmt.Sprintf(`
%s

%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type = "cluster"
}

locals {
  flavor = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0]
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test"
  
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id         = local.flavor.id
  engine_version    = "3.8.35"
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  broker_num        = 3
  access_user       = "user"
  password          = "Rabbitmqtest@123"
  public_ip_id      = huaweicloud_vpc_eip.test[0].id
}`, common.TestBaseNetwork(rName), testEip_base(), rName)
}

func testAccDmsRabbitmqInstance_publicID_update(rName string) string {
	return fmt.Sprintf(`
%s

%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type = "cluster"
}

locals {
  flavor = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0]
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test"
  
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id         = local.flavor.id
  engine_version    = "3.8.35"
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  broker_num        = 3
  access_user       = "user"
  password          = "Rabbitmqtest@123"
  public_ip_id      = huaweicloud_vpc_eip.test[1].id
}`, common.TestBaseNetwork(rName), testEip_base(), rName)
}

func testAccDmsRabbitmqInstance_amqp_single(rName string, enableAcl bool) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type = "single.professional"
}

locals {
  flavor = data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0]
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name              = "%[2]s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id         = local.flavor.id
  engine_version    = "AMQP-0-9-1"
  storage_space     = local.flavor.properties[0].min_storage_per_node
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  enable_acl        = "%[3]t"
}`, common.TestBaseNetwork(rName), rName, enableAcl)
}
