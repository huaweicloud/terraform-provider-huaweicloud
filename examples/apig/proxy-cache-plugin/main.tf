data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

resource "huaweicloud_apig_instance" "test" {
  name                  = var.instance_name
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = var.enterprise_project_id

  available_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}

resource "huaweicloud_apig_plugin" "proxy_cache" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = var.plugin_name
  description = "Created by terraform script"
  type        = "proxy_cache"
  content     = jsonencode({
    "cache_key": {
      "system_params": [],
       "parameters": [
        "custom_param"
      ],
      "headers": []
    },
    "cache_http_status_and_ttl": [
      {
        "http_status": [
          202,
          203
        ],
        "ttl": 5
      }
    ],
    "client_cache_control": {
      "mode": "off",
      "datas": []
    },
    "cacheable_headers": [
      "X-Custom-Header"
    ]
  })
}
