---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_custom_lines"
description: |-
  Use this data source to get the list of DNS custom lines within HuaweiCloud.
---

# huaweicloud_dns_custom_lines

Use this data source to get the list of DNS custom lines within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dns_custom_lines" "test" {}
```

## Argument Reference

The following arguments are supported:

* `line_id` - (Optional, String) Specifies the ID of the custom line. Fuzzy search is supported.

* `name` - (Optional, String) Specifies the name of the custom line. Fuzzy search is supported.

* `ip` - (Optional, String) Specifies the IP address used to query custom line which is in the IP address range.

* `status` - (Optional, String) Specifies the status of the custom line.  
  The valid values are as follows:
  + **ACTIVE**
  + **FREEZE**
  + **DISABLE**
  + **ERROR**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `lines` - All custom lines that match the filter parameters.

  The [lines](#lines_struct) structure is documented below.

<a name="lines_struct"></a>
The `lines` block supports:

* `id` - The ID of the custom line.

* `name` - The name of the custom line.

* `ip_segments` - The IP address range of the custom line.

* `status` - The current status of the custom line.

* `description` - The description of the custom line.

* `created_at` - The creation time of the custom line, in RFC339 format.

* `updated_at` - The latest update time of the custom line, in RFC339 format.
