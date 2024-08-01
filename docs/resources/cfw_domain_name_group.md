---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_domain_name_group"
description: |-
  Manages a CFW domain name group resource within HuaweiCloud.
---

# huaweicloud_cfw_domain_name_group

Manages a CFW domain name group resource within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "object_id" {}
variable "name" {}

resource "huaweicloud_cfw_domain_name_group" "test" {
  fw_instance_id = var.fw_instance_id
  object_id      = var.object_id
  name           = var.name
  type           = 0
  description    = "created by terraform"

  domain_names {
    domain_name = "www.cfw-test.com"
    description = "test domain"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, ForceNew) Specifies the firewall instance ID.

  Changing this parameter will create a new resource.

* `object_id` - (Required, String, ForceNew) Specifies the protected object ID.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the domain name group.

* `type` - (Required, Int, ForceNew) Specifies the type of the domain name group.
  The value can be:
  + **0**: means application type;
  + **1**: means network type;

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the domain name group.

* `domain_names` - (Optional, List) Specifies the list of domain names.
  The [domain_names](#DomainNameGroup_DomainNames) structure is documented below.

<a name="DomainNameGroup_DomainNames"></a>
The `domain_names` block supports:

* `domain_name` - (Required, String) Specifies the domain name.

* `description` - (Optional, String) Specifies the description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `config_status` - The config status of the domain name group.

* `ref_count` - The reference count of the domain name group.

* `message` - The exception message of the domain name group.

* `domain_names` - The list of domain names.
  The [domain_names](#DomainNameGroup_DomainNames_Resp) structure is documented below.

<a name="DomainNameGroup_DomainNames_Resp"></a>
The `domain_names` block supports:

* `domain_address_id` - The domain address ID.

* `dns_ips` - The DNS IP list.

## Import

The domainnamegroup can be imported using the `fw_instance_id`, `object_id` and `id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_cfw_domain_name_group.test <fw_instance_id>/<object_id>/<id>
```
