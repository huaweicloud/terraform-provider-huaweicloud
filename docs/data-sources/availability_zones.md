# huaweicloud\_availability\_zones

Use this data source to get a list of availability zones from HuaweiCloud

## Example Usage

```hcl
data "huaweicloud_availability_zones" "zones" {}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the available zones. If omitted, the provider-level region will be used.

* `state` - (Optional) The `state` of the availability zones to match, default ("available").


## Attributes Reference

`id` is set to hash of the returned zone list. In addition, the following attributes
are exported:

* `names` - The names of the availability zones, ordered alphanumerically, that match the queried `state`
