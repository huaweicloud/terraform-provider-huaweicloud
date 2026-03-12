---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_smn_topics"
description: |-
  Use this data source to query the SMN topics available for alarms.
---

# huaweicloud_css_smn_topics

Use this data source to query the SMN topics available for alarms.

## Example Usage

```hcl
variable "domain_id" {}

data "huaweicloud_css_smn_topics" "test" {
  domain_id = var.domain_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_id` - (Required, String) Specifies the domain account ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `topics_name` - The SMN topic name list.
