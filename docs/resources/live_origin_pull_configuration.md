---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_origin_pull_configuration"
description: |-
  Manages a Live origin pull configuration resource within HuaweiCloud.
---

# huaweicloud_live_origin_pull_configuration

Manages a Live origin pull configuration resource within HuaweiCloud.

-> This resource is an operational resource, and destroying it will not change the current origin pull configuration.

## Example Usage

### Create a customer source site with a domain name format for the source site address

```hcl
variable "domain_name" {}
variable "sources" {}
variable "schema" {}

resource "huaweicloud_live_origin_pull_configuration" "test" {
  domain_name = var.domain_name
  source_type = "domain"
  sources     = var.sources
  scheme      = var.scheme
}
```

### Create a customer source site with an IP format source address

```hcl
variable "domain_name" {}
variable "sources_ip" {}
variable "schema" {}

resource "huaweicloud_live_origin_pull_configuration" "test" {
  domain_name = var.domain_name
  source_type = "ipaddr"
  sources_ip  = var.sources_ip
  scheme      = var.scheme
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the streaming domain name.
  Changing this parameter will create a new resource.

* `source_type` - (Required, String) Specifies the type of return to the source.
  The valid values are as follows:
  + **domain**: Return to the source customer's source site, where the source site address is in domain name format.
  + **ipaddr**: Return to the source customer's source site, where the source site address is in IP format.
  + **huawei**: Return to the Huawei origin site, the default value after domain creation.

* `sources` - (Optional, List) Specifies the list of domain names for returning to the source.
  + When `source_type` is set to **domain**, this parameter is required and can be configured with up to `10` values.
    When multiple domain names are configured, if the source return fails, it will be rotated according to the
    configuration order.
  + When `source_type` is set to **ipaddr**, a maximum of `1` value can be configured. If configured, the httpflv HOST
    header should be filled with this domain name and the RTMP tcurl field should be filled with this domain name when
    returning to the source. Otherwise, the current IP address should be used as the HOST.

* `sources_ip` - (Optional, List) Specifies the list of IP addresses for returning to the source.
  Up to `10` can be configured.
  + When `source_type` is set to **ipaddr**, this parameter is required. When configuring multiple IPs, if the return to
    the source fails, it will be rotated according to the configuration order.

* `source_port` - (Optional, Int) Specifies the return port to the source.

* `scheme` - (Optional, String) Specifies the source return protocol.
  This parameter is required when the `source_type` is not **huawei**.
  The valid values are as follows:
  + **http**
  + **rtmp**

* `additional_args` - (Optional, Map) Specifies the parameters carried in the URL when returning to the source client's
  website.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

The resource can be imported using `domain_name`, e.g.

```bash
$ terraform import huaweicloud_live_origin_pull_configuration.test <domain_name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `additional_args`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_live_origin_pull_configuration" "test" { 
  ...

  lifecycle {
    ignore_changes = [
      additional_args,
    ]
  }
}
```
