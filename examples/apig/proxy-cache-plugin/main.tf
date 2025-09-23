data "huaweicloud_availability_zones" "test" {
  count = length(var.availability_zones) == 0 ? 1 : 0
}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_apig_instance" "test" {
  name                  = var.instance_name
  edition               = var.instance_edition
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zones    = length(var.availability_zones) == 0 ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, var.availability_zones_count), null) : var.availability_zones
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_apig_plugin" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = var.plugin_name
  type        = "proxy_cache"
  description = var.plugin_description

  content = jsonencode({
    cache_key = {
      system_params = [],
      parameters = [
        "custom_param"
      ],
      headers = []
    },
    cache_http_status_and_ttl = [
      {
        http_status = [
          202,
          203
        ],
        ttl = 5
      }
    ],
    client_cache_control = {
      mode  = "off",
      datas = []
    },
    cacheable_headers = [
      "X-Custom-Header"
    ]
  })
}
