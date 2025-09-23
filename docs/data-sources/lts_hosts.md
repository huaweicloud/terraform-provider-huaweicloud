---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_hosts"
description: |-
  Use this data source to get the list of the hosts within HuaweiCloud.
---

# huaweicloud_lts_hosts

Use this data source to get the list of the hosts within HuaweiCloud.

## Example Usage

### Query all hosts

```hcl
data "huaweicloud_lts_hosts" "test" {}
```

### Query the hosts by the specified host IDs

```hcl
variables "host_ids" {
  type = list(string)
}

data "huaweicloud_lts_hosts" "test" {
  host_id_list = var.host_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id_list` - (Optional, List) Specifies the list of the host IDs.

* `filter` - (Optional, List) Specifies filtering parameter to query hosts.

  The [filter](#host_filter_struct) structure is documented below.

<a name="host_filter_struct"></a>
The `filter` block supports:

* `host_name_list` - (Optional, List) Specifies the list of the host names.

* `host_ip_list` - (Optional, List) Specifies the list of the host IPs.

* `host_status` - (Optional, String) Specifies the status of the host.  
  The valid values are as follows:
  + **uninstall**
  + **running**
  + **offline**
  + **error**
  + **plugin error**
  + **install-fail**
  + **upgrade failed**
  + **upgrade-fail**
  + **authentication error**

* `host_version` - (Optional, String) Specifies the version of the host.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `hosts` - All hosts that match the filter parameters.

  The [hosts](#lts_hosts_struct) structure is documented below.

<a name="lts_hosts_struct"></a>
The `hosts` block supports:

* `host_id` - The ID of the host.

* `host_name` - The name of the host.

* `host_type` - The type of the host.
  + **linux**
  + **windows**

* `host_ip` - The IP of the host.

* `host_version` - The version of the host.

* `host_status` - The status of the host.

* `updated_at` - The latest update time of the host, in RFC3339 format.
