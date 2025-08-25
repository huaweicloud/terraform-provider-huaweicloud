---
subcategory: "Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_antiddos_open_protection"
description: |-
  Use this resource to open Anti-DDoS protection for an EIP.
---

# huaweicloud_antiddos_open_protection

Use this resource to open Anti-DDoS protection for an EIP.

-> This resource is a one-time action resource for opening Anti-DDoS protection. Deleting this resource will not
disable the protection, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "eip_id" {
  type = string
}
variable "app_type_id" {
  type = number
}
variable "cleaning_access_pos_id" {
  type = number
}
variable "enable_l7" {
  type = bool
}
variable "http_request_pos_id" {
  type = number
}
variable "traffic_pos_id" {
  type = number
}
variable "antiddos_config_id" {
  type = string
}

resource "huaweicloud_antiddos_open_protection" "test" {
  floating_ip_id         = var.eip_id
  app_type_id            = var.app_type_id
  cleaning_access_pos_id = var.cleaning_access_pos_id
  enable_l7              = var.enable_l7
  http_request_pos_id    = var.http_request_pos_id
  traffic_pos_id         = var.traffic_pos_id
  antiddos_config_id     = var.antiddos_config_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `floating_ip_id` - (Required, String, NonUpdatable) Specifies the ID of the EIP to enable Anti-DDoS protection for.

* `app_type_id` - (Required, Int, NonUpdatable) Specifies the application type ID. Only `0` is supported.

* `cleaning_access_pos_id` - (Required, Int, NonUpdatable) Specifies the cleaning access position ID.
  The value can be:
  + `1`: 10M.
  + `2`: 30M.
  + `3`: 50M.
  + `4`: 70M.
  + `5`: 100M.
  + `6`: 150M.
  + `7`: 200M.
  + `8`: 250M.
  + `9`: 300M.
  + `10`: 500M.
  + `11`: 800M.
  + `88`: 1000M.
  + `99`: Default protection.

* `enable_l7` - (Required, Bool, NonUpdatable) Specifies whether to enable L7 protection. Only **false** is supported.

* `http_request_pos_id` - (Required, Int, NonUpdatable) Specifies the HTTP request position ID. Only `1` is supported.

* `traffic_pos_id` - (Required, Int, NonUpdatable) Specifies the traffic position ID.
  This field can be configured with the same value as `cleaning_access_pos_id`.

* `antiddos_config_id` - (Optional, String, NonUpdatable) Specifies the Anti-DDoS configuration ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
