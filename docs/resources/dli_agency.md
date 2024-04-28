---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_agency"
description: ""
---

# huaweicloud_dli_agency

Assign Agency Permissions of DLI service within HuaweiCloud.  
Once service authorization has succeeded, an agency named **dli_admin_agency** on IAM will be created.
You can only create one this resource.

## Example Usage

```hcl
  resource "huaweicloud_dli_agency" "test" {
    roles = [
      "obs_adm",
      "dis_adm",
      "ctable_adm",
      "vpc_netadm",
      "smn_adm",
      "te_admin",
      "iam_readOnlyAccess"
    ]
  }
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `roles` - (Required, List) The list of roles.  
  The options are as follows:
   + **te_admin**: Tenant Administrator permissions are required to access data from OBS to execute Flink jobs on DLI,
     for example, obtaining OBS/DWS data sources, log dump (including bucket authorization), checkpoint enabling,
     and job import and export. Due to cloud service cache differences, operations require about 60 minutes to take effect.
   + **dis_adm**: DIS Administrator permissions are required to use DIS data as the data source of DLI Flink jobs.
     Due to cloud service cache differences, operations require about 30 minutes to take effect.
   + **ctable_adm**: CloudTable Administrator permissions are required to use CloudTable data as the data source of
     DLI Flink jobs.
   + **vpc_netadm**: VPC Administrator permissions are required to use the VPC, subnet, route, VPC peering connection,
     and port for DLI datasource connections.
   + **smn_adm**: SMN Administrator permissions are required to receive notifications when a DLI job fails.
   + **obs_adm**: OBS Administrator permissions are required to use the Object Storage service.
   + **iam_readOnlyAccess**: IAM ReadOnlyAccess permissions are required to obtain information about the
     IAM users who have never logged in to DLI.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which value is **dli_admin_agency**.

* `version` - Agency version information.

## Import

The agency can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dli_agency.test dli_admin_agency
```
