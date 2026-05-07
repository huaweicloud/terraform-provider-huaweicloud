---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eip_publicip_instances_count"
description: |-
  Use this dataSource to get the number of public IP instances.
---

# huaweicloud_eip_publicip_instances_count

Use this dataSource to get the number of public IP instances.

## Example Usage

```hcl
data "huaweicloud_eip_publicip_instances_count" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instance_num` - The number of public IP instances.
