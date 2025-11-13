data "huaweicloud_availability_zones" "test" {
  count = length(var.availability_zones) < var.availability_zones_count ? 1 : 0
}

resource "huaweicloud_vpc" "test" {
  count = var.vpc_id == "" && var.subnet_id == "" ? 1 : 0

  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  count = var.subnet_id == "" ? 1 : 0

  vpc_id            = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  name              = var.subnet_name
  cidr              = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 0)
  gateway_ip        = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : var.subnet_cidr != "" ? cidrhost(var.subnet_cidr, 1) : cidrhost(cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 0), 1)
  availability_zone = length(var.availability_zones) > 0 ? try(var.availability_zones[0], null) : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_apig_instance" "test" {
  name                  = var.instance_name
  edition               = var.instance_edition
  vpc_id                = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  subnet_id             = var.subnet_id != "" ? var.subnet_id : huaweicloud_vpc_subnet.test[0].id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zones    = length(var.availability_zones) < var.availability_zones_count ? concat(var.availability_zones, try(slice([for v in try(data.huaweicloud_availability_zones.test[0].names, []) : v if !contains(var.availability_zones, v)], 0, var.availability_zones_count - length(var.availability_zones)), [])) : var.availability_zones
  enterprise_project_id = var.enterprise_project_id

  lifecycle {
    ignore_changes = [
      availability_zones
    ]
  }
}

resource "huaweicloud_apig_plugin" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = var.plugin_name
  type        = "proxy_cache"
  description = var.plugin_description

  content = jsonencode({
    cache_key                 = {
      system_params = [],
      parameters    = [
        "custom_param"
      ],
      headers       = []
    },
    cache_http_status_and_ttl = [
      {
        http_status = [
          202,
          203
        ],
        ttl         = 5
      }
    ],
    client_cache_control      = {
      mode  = "off",
      datas = []
    },
    cacheable_headers         = [
      "X-Custom-Header"
    ]
  })
}
