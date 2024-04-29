---
subcategory: "Intelligent EdgeCloud (IEC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iec_sites"
description: ""
---

# huaweicloud_iec_sites

Use this data source to get the available of HuaweiCloud IEC sites.

## Example Usage

### Basic IEC Sites

```hcl
data "huaweicloud_iec_sites" "iec_sites" {}
```

## Argument Reference

The following arguments are supported:

* `area` - (Optional, String) Specifies the area of the IEC sites located.

* `province` - (Optional, String) Specifies the province of the IEC sites located.

* `city` - (Optional, String) Specifies the city of the IEC sites located.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `sites` - An array of one or more IEC service sites. The sites object structure is documented below.

The `sites` block supports:

* `id` - The ID of the IEC service site.
* `name` - The name of the IEC service site.
* `area` - The area of the IEC service site located.
* `province` - The province of the IEC service site located.
* `city` - The city of the IEC service site located.
* `status` - The status of the IEC service site.

* `lines` - An array of one or more EIP lines. The object structure is documented below.
  + `id` - The ID of the EIP line.
  + `name` - The name of the EIP line.
  + `operator` - The operator information of the EIP line.
  + `ip_version` - The supported IP version.
