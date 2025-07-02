package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dcs/v2/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDcsResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DcsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID)
}

func TestAccDcsInstances_basic(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.id", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "timeout"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "launched_at"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "cache_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "cpu_type"),
					resource.TestCheckResourceAttrSet(resourceName, "replica_count"),
					resource.TestCheckResourceAttrSet(resourceName, "readonly_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "transparent_client_ip_enable"),
					resource.TestCheckResourceAttrSet(resourceName, "sharding_count"),
					resource.TestCheckResourceAttrSet(resourceName, "product_type"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.bandwidth"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_info.0.begin_time", ""),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.current_time"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_info.0.end_time", ""),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.expand_count"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.expand_effect_time"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.expand_interval_time"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.max_expand_count"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_info.0.next_expand_time", ""),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.task_running"),
				),
			},
			{
				Config: testAccDcsV1Instance_updated(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6389"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.begin_at", "01:00-02:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.save_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_at.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.id", "10"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "latency-monitor-threshold"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "120"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "launched_at"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "cache_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "cpu_type"),
					resource.TestCheckResourceAttrSet(resourceName, "replica_count"),
					resource.TestCheckResourceAttrSet(resourceName, "readonly_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "transparent_client_ip_enable"),
					resource.TestCheckResourceAttrSet(resourceName, "sharding_count"),
					resource.TestCheckResourceAttrSet(resourceName, "product_type"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.bandwidth"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_info.0.begin_time", ""),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.current_time"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_info.0.end_time", ""),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.expand_count"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.expand_effect_time"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.expand_interval_time"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.max_expand_count"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_info.0.next_expand_time", ""),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth_info.0.task_running"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at", "parameters",
					"used_memory", "bandwidth_info"},
			},
		},
	})
}

func TestAccDcsInstances_ha_change_capacity(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_ha(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_ha_expand_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				Config: testAccDcsV1Instance_ha_reduce_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_ha_expand_replica(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_ha(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_ha_expand_replica(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_ha_to_proxy(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_ha(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_ha_to_proxy(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_rw_change_capacity(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_rw(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_rw_expand_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "16"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				Config: testAccDcsV1Instance_rw_reduce_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_rw_expand_replica(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_rw(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_rw_expand_replica(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_rw_to_proxy(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_rw(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_rw_to_proxy(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_proxy_change_capacity(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_proxy(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_proxy_expand_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "16"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				Config: testAccDcsV1Instance_proxy_reduce_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_proxy_to_ha(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_proxy(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_proxy_to_ha(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_proxy_to_rw(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_proxy(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_proxy_to_rw(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_cluster_change_capacity(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_cluster(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_cluster_expand_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				Config: testAccDcsV1Instance_cluster_reduce_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_cluster_expand_replica(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_cluster(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_cluster_expand_replica(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_whitelists(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_whitelists(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.group_name", "test-group1"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.0", "192.168.10.100"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.1", "192.168.0.0/24"),
				),
			},
			{
				Config: testAccDcsV1Instance_whitelists_update(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.group_name", "test-group2"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.0", "172.16.10.100"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.1", "172.16.0.0/24"),
				),
			},
		},
	})
}

func TestAccDcsInstances_tiny(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_tiny(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.125"),
				),
			},
		},
	})
}

func TestAccDcsInstances_single(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_single(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "2"),
				),
			},
		},
	})
}

func TestAccDcsInstances_prePaid(t *testing.T) {
	var instance instances.DcsInstance
	var rName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsInstance_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccDcsInstance_prePaid_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "bandwidth_info",
					"used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_ssl(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_ssl(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "6.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.125"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_update_ssl(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "6.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.125"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at", "parameters",
					"bandwidth_info", "used_memory"},
			},
		},
	})
}

func testAccDcsV1Instance_basic(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode       = "ha"
  capacity         = 1
  cpu_architecture = "x86_64"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "22:00:00"
  maintain_end       = "23:00:00"

  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  rename_commands = {
    command  = "command001"
    keys     = "keys001"
    flushall = "flushall001"
    flushdb  = "flushdb001"
    hgetall  = "hgetall001"
  }

  parameters {
    id    = "1"
    name  = "timeout"
    value = "100"
  }

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, instanceName)
}

func testAccDcsV1Instance_updated(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode       = "ha"
  capacity         = 2
  cpu_architecture = "x86_64"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6389
  capacity           = 2
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"

  backup_policy {
    backup_type = "auto"
    begin_at    = "01:00-02:00"
    period_type = "weekly"
    backup_at   = [1, 2, 4]
    save_days   = 2
  }

  rename_commands = {
    command  = "command001"
    keys     = "keys001"
    flushall = "flushall001"
    flushdb  = "flushdb001"
    hgetall  = "hgetall001"
  }

  parameters {
    id    = "10"
    name  = "latency-monitor-threshold"
    value = "120"
  }

  tags = {
    key   = "value_update"
    owner = "terraform_update"
  }
}`, instanceName)
}

func testAccDcsV1Instance_ha(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 1
  name           = "redis.ha.xu1.large.r2.1"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "22:00:00"
  maintain_end       = "23:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_ha_expand_capacity(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 4
  name           = "redis.ha.xu1.large.r2.4"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 4
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_ha_reduce_capacity(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 1
  name           = "redis.ha.xu1.large.r2.1"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_ha_expand_replica(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 1
  name           = "redis.ha.xu1.large.r4.1"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_ha_to_proxy(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 4
  name           = "redis.proxy.xu1.large.4"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 4
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_rw(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 8
  name           = "redis.ha.xu1.large.p2.8"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 8
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "22:00:00"
  maintain_end       = "23:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_rw_expand_capacity(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 16
  name           = "redis.ha.xu1.large.p2.16"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 16
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_rw_reduce_capacity(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 8
  name           = "redis.ha.xu1.large.p2.8"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 8
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_rw_expand_replica(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 8
  name           = "redis.ha.xu1.large.p4.8"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 8
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_rw_to_proxy(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 8
  name           = "redis.proxy.xu1.large.8"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 8
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_proxy(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 4
  name           = "redis.proxy.xu1.large.4"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 4
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "22:00:00"
  maintain_end       = "23:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_proxy_expand_capacity(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 16
  name           = "redis.proxy.xu1.large.s1.16"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 16
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_proxy_reduce_capacity(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 8
  name           = "redis.proxy.xu1.large.s1.8"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test_update"
  engine             = "Redis"
  port               = 6388
  capacity           = 8
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_proxy_to_ha(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 4
  name           = "redis.ha.xu1.large.r2.4"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 4
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_proxy_to_rw(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 8
  name           = "redis.ha.xu1.large.p2.8"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 8
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_cluster(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 4
  name           = "redis.cluster.xu1.large.r2.4"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 4
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "22:00:00"
  maintain_end       = "23:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_cluster_expand_capacity(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 8
  name           = "redis.cluster.xu1.large.r2.s1.8"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 8
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_cluster_reduce_capacity(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 4
  name           = "redis.cluster.xu1.large.r2.s1.4"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 4
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_cluster_expand_replica(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "5.0"
  capacity       = 4
  name           = "redis.cluster.xu1.large.r3.4"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 4
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin     = "06:00:00"
  maintain_end       = "07:00:00"
}`, instanceName)
}

func testAccDcsV1Instance_tiny(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 0.125
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  
  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1]
    save_days   = 1
  }
}`, instanceName)
}

func testAccDcsV1Instance_single(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "single"
  capacity   = 2
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 2
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}`, instanceName)
}

func testAccDcsV1Instance_whitelists(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 2
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 2
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  
  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1]
    save_days   = 1
  }

  whitelists {
    group_name = "test-group1"
    ip_address = ["192.168.10.100", "192.168.0.0/24"]
  }
}`, instanceName)
}

func testAccDcsV1Instance_whitelists_update(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 2
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 2
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  
  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1]
    save_days   = 1
  }

  whitelists {
    group_name = "test-group2"
    ip_address = ["172.16.10.100", "172.16.0.0/24"]
  }
}`, instanceName)
}

func testAccDcsInstance_prePaid(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 0.125
}

resource "huaweicloud_dcs_instance" "test" {
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), [])
  name               = "%s"
  engine             = "Redis"
  engine_version     = "5.0"
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  capacity           = 0.125
  password           = "Huawei_test"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"
}`, instanceName)
}

func testAccDcsInstance_prePaid_update(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode = "ha"
  capacity   = 0.5
}

resource "huaweicloud_dcs_instance" "test" {
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), [])
  name               = "%s"
  engine             = "Redis"
  engine_version     = "5.0"
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  capacity           = 0.5
  password           = "Huawei_test"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}`, instanceName)
}

func testAccDcsV1Instance_ssl(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode     = "ha"
  capacity       = 0.125
  engine_version = "6.0"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%s"
  engine_version     = "6.0"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  ssl_enable         = true
}`, instanceName)
}

func testAccDcsV1Instance_update_ssl(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode     = "ha"
  capacity       = 0.125
  engine_version = "6.0"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%s"
  engine_version     = "6.0"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  ssl_enable         = false
}`, instanceName)
}
