---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_system_lines"
description: |-
  Use this data source to get the list of DNS system lines within HuaweiCloud.
---

# huaweicloud_dns_system_lines

Use this data source to get the list of DNS system lines within HuaweiCloud.

## Example Usage

### Query system lines (Chinese)

```hcl
data "huaweicloud_dns_system_lines" "test" {}
```

### Query system lines (English)

```hcl
data "huaweicloud_dns_system_lines" "test" {
  locale = "en-us"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the system lines are located.  
  If omitted, the default value is the provider-level region.

* `locale` - (Optional, String) Specifies the display language.  
  The valid values are as follows:
  + **zh-cn**
  + **en-us**
  + **es-us**
  + **pt-br**

  If omitted, the default value is **zh-cn**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `lines` - The list of system lines that match the filter parameters.  
  The [lines](#system_lines) structure is documented below.

<a name="system_lines"></a>
The `lines` block supports:

* `id` - The ID of the system line.

* `name` - The name of the system line.

* `father_id` - The ID of the parent line.
