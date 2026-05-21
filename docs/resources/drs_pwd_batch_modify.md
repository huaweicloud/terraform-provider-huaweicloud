---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_pwd_batch_modify"
description: |-
  Manages a resource to batch modify DRS jobs password within HuaweiCloud.
---

# huaweicloud_drs_pwd_batch_modify

Manages a resource to batch modify DRS jobs password within HuaweiCloud.

-> 1. This resource is a one-time action resource used to batch modify DRS jobs password. Deleting this
  resource will not restore the modified password or undo the modify action, but will only remove the resource information
  from the tf state file.
  <br/>2. Only tasks in the following status can be used: STARTJOBING, STARTJOB_FAILED, FULL_TRANSFER_STARTED,
  FULL_TRANSFER_FAILED, FULL_TRANSFER_COMPLETE, INCRE_TRANSFER_STARTED, INCRE_TRANSFER_FAILED, PAUSING.

## Example Usage

```hcl
variable "jobs" { 
  type = list(object({
    job_id         = string
    db_password    = string
    end_point_type = string
  }))
  sensitive = true
}

resource "huaweicloud_drs_pwd_batch_modify" "test" {
  dynamic "jobs" {
    for_each = var.jobs
    
    content {
      job_id         = jobs.value.job_id
      db_password    = jobs.value.db_password
      end_point_type = jobs.value.end_point_type
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `jobs` - (Required, List, NonUpdatable) Specifies the list of jobs to batch modify database password.
  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `job_id` - (Required, String, NonUpdatable) Specifies the job ID.

* `db_password` - (Required, String, NonUpdatable) Specifies the database password.

* `end_point_type` - (Required, String, NonUpdatable) Specifies the endpoint type.
  The valid values are as follows:
  + **so**: Source database.
  + **ta**: Target database.

* `kerberos` - (Optional, List, NonUpdatable) Specifies the Kerberos authentication information.
  The [kerberos](#kerberos_struct) structure is documented below.

<a name="kerberos_struct"></a>
The `kerberos` block supports:

* `krb5_conf_file` - (Optional, String, NonUpdatable) Specifies the Kerberos configuration file.

* `key_tab_file` - (Optional, String, NonUpdatable) Specifies the keytab file.

* `domain_name` - (Optional, String, NonUpdatable) Specifies the domain name.

* `user_principal` - (Optional, String, NonUpdatable) Specifies the Kerberos user principal.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `results` - The results of batch modifying tasks.
  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `id` - The job ID.

* `status` - The operation status. The valid values are as follows:
  + **success**: Success.
  + **failed**: Failed.

* `end_point_type` - The endpoint type.
  The valid values are as follows:
  + **so**: Source database.
  + **ta**: Target database.
