---
subcategory: "Deprecated"
---

# huaweicloud\_vbs\_backup\_policy

!> **Warning:** It has been deprecated.

The VBS Backup Policy data source provides details about a specific VBS backup policy.
This is an alternative to `huaweicloud_vbs_backup_policy_v2`

## Example Usage

 ```hcl

variable "policy_name" {}
variable "policy_id" {}

data "huaweicloud_vbs_backup_policy" "policies" {
  name = var.policy_name
  id   = var.policy_id
}
 ```

## Argument Reference

The arguments of this data source act as filters for querying the available VBS backup policy.
The given filters must match exactly one VBS backup policy whose data will be exported as attributes.

* `region` - (Optional) The region in which to obtain the VBS backup policy. If omitted, the provider-level region will work as default.

* `id` (Optional) - The ID of the specific VBS backup policy to retrieve.

* `name` (Optional) - The name of the specific VBS backup policy to retrieve.

* `status` (Optional) - The status of the specific VBS backup policy to retrieve. The values can be ON or OFF

**filter_tags** **- (Optional)** Represents the list of tags. Backup policy with these tags will be filtered.

* `key` - (Required) Specifies the tag key. Tag keys must be unique.

* `values` - (Required) Specifies the List of tag values. This list can have a maximum of 10 values and all be unique.



## Attributes Reference

The following attributes are exported:

* `id` - See Argument Reference above.

* `name` - See Argument Reference above.

* `status` - See Argument Reference above.

* `start_time` - Specifies the start time of the backup job.The value is in the HH:mm format.                                                         

* `retain_first_backup` - Specifies whether to retain the first backup in the current month. 

* `rentention_num` - Specifies number of retained backups.

* `frequency` - Specifies the backup interval. The value is in the range of 1 to 14 days.

**tags** - Represents the list of tag details associated with the backup policy.

* `key` - Specifies the tag key. 

* `value` - Specifies the tag value. 
