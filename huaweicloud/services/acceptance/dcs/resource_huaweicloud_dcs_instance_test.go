package dcs

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDcsResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DCS instance: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccDcsInstances_basic(t *testing.T) {
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "1"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(rName, "parameters.0.id", "2"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "maxmemory-policy"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "volatile-lfu"),
					resource.TestCheckResourceAttr(rName, "big_key_enable_auto_scan", "true"),
					resource.TestCheckResourceAttr(rName, "big_key_schedule_at.0", "10:00"),
					resource.TestCheckResourceAttr(rName, "hot_key_enable_auto_scan", "false"),
					resource.TestCheckResourceAttr(rName, "hot_key_schedule_at.0", "13:00"),
					resource.TestCheckResourceAttr(rName, "expire_key_enable_auto_scan", "true"),
					resource.TestCheckResourceAttr(rName, "expire_key_interval", "20"),
					resource.TestCheckResourceAttr(rName, "expire_key_timeout", "100"),
					resource.TestCheckResourceAttr(rName, "expire_key_scan_keys_count", "20000"),
					resource.TestCheckResourceAttr(rName, "transparent_client_ip_enable", "false"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "launched_at"),
					resource.TestCheckResourceAttrSet(rName, "subnet_cidr"),
					resource.TestCheckResourceAttrSet(rName, "cache_mode"),
					resource.TestCheckResourceAttrSet(rName, "cpu_type"),
					resource.TestCheckResourceAttrSet(rName, "replica_count"),
					resource.TestCheckResourceAttrSet(rName, "readonly_domain_name"),
					resource.TestCheckResourceAttrSet(rName, "sharding_count"),
					resource.TestCheckResourceAttrSet(rName, "product_type"),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.bandwidth"),
					resource.TestCheckResourceAttr(rName, "bandwidth_info.0.begin_time", ""),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.current_time"),
					resource.TestCheckResourceAttr(rName, "bandwidth_info.0.end_time", ""),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.expand_count"),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.expand_effect_time"),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.expand_interval_time"),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.max_expand_count"),
					resource.TestCheckResourceAttr(rName, "bandwidth_info.0.next_expand_time", ""),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.task_running"),
					resource.TestCheckResourceAttrSet(rName, "big_key_updated_at"),
					resource.TestCheckResourceAttrSet(rName, "hot_key_updated_at"),
					resource.TestCheckResourceAttrSet(rName, "expire_key_first_scan_at"),
					resource.TestCheckResourceAttrSet(rName, "expire_key_updated_at"),
				),
			},
			{
				Config: testAccDcsV1Instance_updated(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6389"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "2"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
					resource.TestCheckResourceAttr(rName, "backup_policy.0.begin_at", "01:00-02:00"),
					resource.TestCheckResourceAttr(rName, "backup_policy.0.save_days", "2"),
					resource.TestCheckResourceAttr(rName, "backup_policy.0.backup_at.#", "3"),
					resource.TestCheckResourceAttr(rName, "parameters.0.id", "2"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "maxmemory-policy"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "allkeys-lfu"),
					resource.TestCheckResourceAttr(rName, "big_key_enable_auto_scan", "false"),
					resource.TestCheckResourceAttr(rName, "big_key_schedule_at.0", "17:00"),
					resource.TestCheckResourceAttr(rName, "hot_key_enable_auto_scan", "true"),
					resource.TestCheckResourceAttr(rName, "hot_key_schedule_at.0", "20:00"),
					resource.TestCheckResourceAttr(rName, "expire_key_enable_auto_scan", "false"),
					resource.TestCheckResourceAttr(rName, "transparent_client_ip_enable", "true"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "launched_at"),
					resource.TestCheckResourceAttrSet(rName, "subnet_cidr"),
					resource.TestCheckResourceAttrSet(rName, "cache_mode"),
					resource.TestCheckResourceAttrSet(rName, "cpu_type"),
					resource.TestCheckResourceAttrSet(rName, "replica_count"),
					resource.TestCheckResourceAttrSet(rName, "readonly_domain_name"),
					resource.TestCheckResourceAttrSet(rName, "sharding_count"),
					resource.TestCheckResourceAttrSet(rName, "product_type"),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.bandwidth"),
					resource.TestCheckResourceAttr(rName, "bandwidth_info.0.begin_time", ""),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.current_time"),
					resource.TestCheckResourceAttr(rName, "bandwidth_info.0.end_time", ""),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.expand_count"),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.expand_effect_time"),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.expand_interval_time"),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.max_expand_count"),
					resource.TestCheckResourceAttr(rName, "bandwidth_info.0.next_expand_time", ""),
					resource.TestCheckResourceAttrSet(rName, "bandwidth_info.0.task_running"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "1"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_ha_expand_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "4"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				Config: testAccDcsV1Instance_ha_reduce_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "1"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "1"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_ha_expand_replica(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "1"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "1"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_ha_to_proxy(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "4"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "8"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_rw_expand_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "16"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				Config: testAccDcsV1Instance_rw_reduce_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "8"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "8"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_rw_expand_replica(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "8"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "8"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_rw_to_proxy(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "8"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "4"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_proxy_expand_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "16"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				Config: testAccDcsV1Instance_proxy_reduce_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "8"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "4"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_proxy_to_ha(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "4"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "4"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_proxy_to_rw(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "8"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "4"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_cluster_expand_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "8"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				Config: testAccDcsV1Instance_cluster_reduce_capacity(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "4"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "4"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "22:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "23:00:00"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_cluster_expand_replica(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "port", "6388"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(rName, "capacity", "4"),
					resource.TestCheckResourceAttr(rName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(rName, "maintain_end", "07:00:00"),
				),
			},
			{
				ResourceName:      rName,
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
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "whitelist_enable", "true"),
					resource.TestCheckResourceAttr(rName, "whitelists.#", "1"),
					resource.TestCheckResourceAttr(rName, "whitelists.0.group_name", "test-group1"),
					resource.TestCheckResourceAttr(rName, "whitelists.0.ip_address.0", "192.168.10.100"),
					resource.TestCheckResourceAttr(rName, "whitelists.0.ip_address.1", "192.168.0.0/24"),
				),
			},
			{
				Config: testAccDcsV1Instance_whitelists_update(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "whitelist_enable", "true"),
					resource.TestCheckResourceAttr(rName, "whitelists.#", "1"),
					resource.TestCheckResourceAttr(rName, "whitelists.0.group_name", "test-group2"),
					resource.TestCheckResourceAttr(rName, "whitelists.0.ip_address.0", "172.16.10.100"),
					resource.TestCheckResourceAttr(rName, "whitelists.0.ip_address.1", "172.16.0.0/24"),
				),
			},
		},
	})
}

func TestAccDcsInstances_single(t *testing.T) {
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(rName, "capacity", "2"),
				),
			},
		},
	})
}

func TestAccDcsInstances_prePaid(t *testing.T) {
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
				Config: testAccDcsInstance_prePaid(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccDcsInstance_prePaid_update(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "auto_renew", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "bandwidth_info",
					"used_memory"},
			},
		},
	})
}

func TestAccDcsInstances_ssl(t *testing.T) {
	var instance interface{}

	var instanceName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
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
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "6.0"),
					resource.TestCheckResourceAttr(rName, "capacity", "2"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(rName, "ssl_enable", "true"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "port"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_update_ssl(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", instanceName),
					resource.TestCheckResourceAttr(rName, "engine", "Redis"),
					resource.TestCheckResourceAttr(rName, "engine_version", "6.0"),
					resource.TestCheckResourceAttr(rName, "capacity", "2"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dcs_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(rName, "ssl_enable", "false"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "port"),
					resource.TestCheckResourceAttrSet(rName, "domain_name"),
				),
			},
			{
				ResourceName:      rName,
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
	firstScanTime := time.Now().UTC().Add(1 * time.Hour)
	firstScanTimeString := firstScanTime.Format("2006-01-02T15:04:05.000z")
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

resource "huaweicloud_dcs_instance" "test" {
  name                         = "%[1]s"
  engine_version               = "5.0"
  password                     = "Huawei_test"
  engine                       = "Redis"
  port                         = 6388
  capacity                     = 1
  vpc_id                       = data.huaweicloud_vpc.test.id
  subnet_id                    = data.huaweicloud_vpc_subnet.test.id
  availability_zones           = [data.huaweicloud_availability_zones.test.names[0]]
  flavor                       = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin               = "22:00:00"
  maintain_end                 = "23:00:00"
  transparent_client_ip_enable = false

  big_key_enable_auto_scan    = true
  big_key_schedule_at         = ["10:00"]
  hot_key_enable_auto_scan    = false
  hot_key_schedule_at         = ["13:00"]
  expire_key_enable_auto_scan = true
  expire_key_first_scan_at    = "%[2]s"
  expire_key_interval         = 20
  expire_key_timeout          = 100
  expire_key_scan_keys_count  = 20000

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
    id    = "2"
    name  = "maxmemory-policy"
    value = "volatile-lfu"
  }

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, instanceName, firstScanTimeString)
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

resource "huaweicloud_dcs_instance" "test" {
  name                         = "%[1]s"
  engine_version               = "5.0"
  password                     = "Huawei_test"
  engine                       = "Redis"
  port                         = 6389
  capacity                     = 2
  vpc_id                       = data.huaweicloud_vpc.test.id
  subnet_id                    = data.huaweicloud_vpc_subnet.test.id
  availability_zones           = [data.huaweicloud_availability_zones.test.names[0]]
  flavor                       = data.huaweicloud_dcs_flavors.test.flavors[0].name
  maintain_begin               = "06:00:00"
  maintain_end                 = "07:00:00"
  transparent_client_ip_enable = true

  big_key_enable_auto_scan    = false
  big_key_schedule_at         = ["17:00"]
  hot_key_enable_auto_scan    = true
  hot_key_schedule_at         = ["20:00"]
  expire_key_enable_auto_scan = false

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
    id    = "2"
    name  = "maxmemory-policy"
    value = "allkeys-lfu"
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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

resource "huaweicloud_dcs_instance" "test" {
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
  capacity   = 2
}

resource "huaweicloud_dcs_instance" "test" {
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), [])
  name               = "%s"
  engine             = "Redis"
  engine_version     = "5.0"
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  capacity           = 2
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
  capacity   = 4
}

resource "huaweicloud_dcs_instance" "test" {
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), [])
  name               = "%s"
  engine             = "Redis"
  engine_version     = "5.0"
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  capacity           = 4
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
  capacity       = 2
  engine_version = "6.0"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%s"
  engine_version     = "6.0"
  engine             = "Redis"
  capacity           = 2
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
  capacity       = 2
  engine_version = "6.0"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%s"
  engine_version     = "6.0"
  engine             = "Redis"
  capacity           = 2
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
  ssl_enable         = false
}`, instanceName)
}
