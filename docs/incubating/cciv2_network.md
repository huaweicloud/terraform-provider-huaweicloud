---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_network"
description: |-
  Manages a CCI Network resource within HuaweiCloud.
---

# huaweicloud_cciv2_network

Manages a CCI Network resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}
variable "name" {}
variable "project_id" {}
variable "domain_id" {}
variable "subnet_id" {}
variable "security_group_ids" {}

resource "huaweicloud_cciv2_network" "test" {
  namespace = var.namespace
  name      = var.name

  annotations = {
    "yangtse.io/project-id"                 = var.project_id,
    "yangtse.io/domain-id"                  = var.domain_id,
    "yangtse.io/warm-pool-size"             = "10",
    "yangtse.io/warm-pool-recycle-interval" = "2",
  }

  subnets {
    subnet_id = var.subnet_id
  }

  security_group_ids = var.security_group_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the CCI network.

* `namespace` - (Required, String, NonUpdatable) Specifies the namespace.

* `annotations` - (Optional, Map) Specifies the annotations of the CCI network.

* `ip_families` - (Optional, List, NonUpdatable) Specifies the IP families of the CCI network.
  When IPV6 is enabled, the value can be **["IPv4", "IPv6"]**.

* `security_group_ids` - (Optional, List) Specifies the security group IDs of the CCI network.

* `subnets` - (Optional, List, NonUpdatable) Specifies the subnets of the CCI network.
  The [subnets](#block_subnets) structure is documented below.

<a name="block_subnets"></a>
The `subnets` block supports:

* `subnet_id` - (Optional, String, NonUpdatable) Specifies the subnet ID of the CCI network.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `api_version` - The API version of the CCI network.

* `kind` - The kind of the CCI network.

* `creation_timestamp` - The creation timestamp of the namespace.

* `finalizers` - The finalizers of the namespace.

* `resource_version` - The resource version of the namespace.

* `status` - The status of the namespace.
  The [status](#attrblock_status) structure is documented below.

* `uid` - The uid of the namespace.

<a name="attrblock_status"></a>
The `status` block supports:

* `conditions` - The conditions of the CCI network.
  The [conditions](#attrblock_status_conditions) structure is documented below.

* `status` - The status of the CCI network.

* `subnet_attrs` - The subnet attributes of the CCI network.
  The [subnet_attrs](#attrblock_status_subnet_attrs) structure is documented below.

<a name="attrblock_status_conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time of the CCI network conditions.

* `message` - The message of the CCI network conditions.

* `reason` - The reason of the CCI network conditions.

* `status` - The status of the CCI network conditions.

* `type` - The type of the CCI network conditions.

<a name="attrblock_status_subnet_attrs"></a>
The `subnet_attrs` block supports:

* `network_id` - The ID of the CCI network.

* `subnet_v4_id` - The subnet IPv4 ID of the CCI network.

* `subnet_v6_id` - The subnet IPv6 ID of the CCI network.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The CCI Network can be imported using `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cciv2_network.test <namespace>/<name>
```
