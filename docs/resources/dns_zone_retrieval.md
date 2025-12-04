---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_zone_retrieval"
description: |-
  Manages a DNS zone retrieval resource within HuaweiCloud.
---

# huaweicloud_dns_zone_retrieval

Manages a DNS zone retrieval resource within HuaweiCloud.

## Example Usage

```hcl
variable "zone_name" {}

resource "huaweicloud_dns_zone_retrieval" "test" {
  zone_name = var.zone_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `zone_name` - (Required, String, NonUpdatable) Specifies the zone name.
  Note the `.` at the end of the name is available.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `retrieval_id` - Indicates the retrieval ID

* `record` - Indicates the record detail.
  The [record](#attrblock--record) structure is documented below.

* `status` - Indicates the status.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the last update time.

<a name="attrblock--record"></a>
The `record` block supports:

* `host` - Indicates the record host.

* `value` - Indicates the record value.

## Import

The zone retrieval can be imported using `zone_name`, e.g.

```bash
$ terraform import huaweicloud_dns_zone_retrieval.test <zone_name>
```
