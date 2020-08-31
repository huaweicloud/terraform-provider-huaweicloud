---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_az"
sidebar_current: "docs-huaweicloud-datasource-dms-az"
description: |-
  Get information on an HuaweiCloud dms az.
---

# huaweicloud\_dms\_az

Use this data source to get the ID of an available HuaweiCloud dms az.
This is an alternative to `huaweicloud_dms_az_v1`

## Example Usage

```hcl

data "huaweicloud_dms_az" "az1" {
  name = "可用区1"
  port = "8002"
  code = "cn-north-1a"
}
```

## Argument Reference

* `name` - (Required) Indicates the name of an AZ.

* `code` - (Optional) Indicates the code of an AZ.

* `port` - (Required) Indicates the port number of an AZ.


## Attributes Reference

`id` is set to the ID of the found az. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `code` - See Argument Reference above.
* `port` - See Argument Reference above.
