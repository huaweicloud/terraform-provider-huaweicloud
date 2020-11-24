---
subcategory: "Distributed Cache Service"
---

# huaweicloud\_dcs\_az

Use this data source to get the ID of an available Huaweicloud dcs az.
This is an alternative to `huaweicloud_dcs_az_v1`

## Example Usage

```hcl

data "huaweicloud_dcs_az" "az1" {
  name = "AZ1"
  port = "8004"
  code = "cn-north-1a"
}
```

## Argument Reference

For details, See [Querying AZ Information](https://support.huaweicloud.com/en-us/api-dcs/dcs-api-0312039.html).

* `region` - (Optional) The region in which to obtain the dcs az. If omitted, the provider-level region will be used.

* `name` - (Optional) Indicates the name of an AZ.

* `code` - (Optional) Indicates the code of an AZ.

* `port` - (Optional) Indicates the port number of an AZ.


## Attributes Reference

`id` is set to the ID of the found az. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `code` - See Argument Reference above.
* `port` - See Argument Reference above.
