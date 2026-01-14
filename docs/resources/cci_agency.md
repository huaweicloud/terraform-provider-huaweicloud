---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cci_agency"
description: |-
  Manages a CCI agency resource within HuaweiCloud.
---

# huaweicloud_cci_agency

Manages a CCI agency resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_cci_agency" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to the domain ID.

* `name` - Indicates the name of the agency.

* `domain_id` - Indicates the domain ID.

* `trust_domain_id` - Indicates the trust domain ID.

* `trust_domain_name` - Indicates the trust domain name.

* `description` - Indicates the description of the agency.

* `duration` - Indicates the duration of the agency.

* `create_time` - Indicates the creation time of the agency.

* `need_update` - Indicates whether the agency need update.

## Import

The CCI agency can be imported using the `id`, e.g.:

```bash
$ terraform import huaweicloud_cci_agency.test <id>
```
