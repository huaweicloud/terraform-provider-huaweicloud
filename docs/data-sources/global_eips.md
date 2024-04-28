---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eips"
description: ""
---

# huaweicloud_global_eips

Use this data source to get the list of global EIPs.

## Example Usage

### Get all global EIPs

```hcl
data "huaweicloud_global_eips" "all" {}
```

### Get specific global EIPs

```hcl
data "huaweicloud_global_eips" "test" {
  status = "inuse"
}
```

## Argument Reference

The following arguments are supported:

* `geip_id` - (Optional, String) Specifies the GEIP ID.

* `internet_bandwidth_id` - (Optional, String) Specifies the global internet bandwidth ID which the GEIP associates to.

* `ip_address` - (Optional, String) Specifies the GEIP IP address.

* `name` - (Optional, String) Specifies the GEIP name.

* `status` - (Optional, String) Specifies the GEIP status.
  
  Valid valus are as follows:
  + **idle**: Not associates with instance.
  + **inuse**: Associates with instance.
  + **pending_create**: Creating.
  + **pending_update**: Updating.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the GEIP.

* `tags` - (Optional, Map) Specifies the tags of the GEIP.

* `associate_instance_id` - (Optional, String) Specifies the instance ID which the GEIP associates to.

* `associate_instance_type` - (Optional, String) Specifies the type of instance which the GEIP associates to.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `global_eips` - The GEIP list.
  The [global_eips](#attrblock--global_eips) structure is documented below.

<a name="attrblock--global_eips"></a>
The `global_eips` block supports:

* `id` - The GEIP ID.

* `access_site` - The access site name.

* `geip_pool_name` - The GEIP pool name.

* `internet_bandwidth_id` - The global internet bandwidth ID which the GEIP associates to.

* `ip_address` - The IP address of GEIP.

* `ip_version` - The IP version of GEIP.

* `isp` - The the internet service provider of the global EIP.

* `name` - The GEIP name.

* `enterprise_project_id` - The enterprise project ID of GEIP.

* `description` - The description of GEIP.

* `tags` - The tags of GEIP.

* `global_connection_bandwidth_id` - The ID of the global connection bandwidth.

* `global_connection_bandwidth_type` - The type of the global connection bandwidth.

* `associate_instance_region` - The region of the associate instance.

* `associate_instance_id` - The ID of the associate instance.

* `associate_instance_type` - The type of the associate instance.

* `frozen` - The global EIP is frozen or not.

* `frozen_info` - The frozen info of the global EIP.

* `polluted` - The global EIP is polluted or not.

* `status` - The status of the global EIP.

* `created_at` - The create time of the global EIP.

* `updated_at` - The update time of the global EIP.
