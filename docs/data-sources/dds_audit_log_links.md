---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_audit_log_links"
description: |-
  Use this data source to get the list of DDS instance audit log links.
---

# huaweicloud_dds_audit_log_links

Use this data source to get the list of DDS instance audit log links.

## Example Usage

```hcl
variable "instance_id" {}
variable "audit_log_ids" {}

data "huaweicloud_dds_audit_log_links" "test" {
  instance_id = var.instance_id
  ids         = var.audit_log_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `ids` - (Required, List) Specifies the list of audit log ids. A maximum of 50 audit log IDs are allowed in the list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `links` - Indicates the list of audit log download links. The validity period is 5 minutes.
