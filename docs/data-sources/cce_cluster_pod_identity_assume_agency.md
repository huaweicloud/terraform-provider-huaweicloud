---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_pod_identity_assume_agency"
description: -|
  Use this data source to get the IAM agency credential of a pod-identity association within a HuaweiCloud CCE cluster.
---

# huaweicloud_cce_cluster_pod_identity_assume_agency

Use this data source to get the IAM agency credential of a pod-identity association within a HuaweiCloud CCE cluster.

This data source is used to obtain temporary security credentials by providing a ServiceAccount token
associated with a pod-identity binding. It allows pods to access Huawei Cloud services using IAM permissions
granted to the associated agency.

## Example Usage

```hcl
variable "cluster_id" {}
variable "service_account_token" {}

data "huaweicloud_cce_cluster_pod_identity_assume_agency" "test" {
  cluster_id = var.cluster_id
  token      = var.service_account_token
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID where the pod-identity association exists.

* `token` - (Required, String) Specifies the ServiceAccount token of the pod-identity association.
  This token is automatically mounted into the pod when pod-identity is enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `assumed_agency` - The agency metadata information corresponding to the credential.
  The [assumed_agency](#assumed_agency_attr) structure is documented below.

* `audience` - The audience attribute passed in when the credential is issued.
  For pod-identity scenario, this value is fixed as **service.cce.pods**.

* `credentials` - The temporary security credential.
  The [credentials](#credentials_attr) structure is documented below.

* `pod_identity_association_id` - The ID of the pod identity association.

* `subject` - The subject information of the token.
  The [subject](#subject_attr) structure is documented below.

<a name="assumed_agency_attr"></a>
The `assumed_agency` block supports:

* `urn` - The unique ID of an agency, in the format of `sts::{account_id}::assumed-agency:{agency_name}/{agency_session_name}`

* `id` - The agency ID, in the format of `{agency_id}:{agency_session_name}`

<a name="credentials_attr"></a>
The `credentials` block supports:

* `access_key_id` - The access key ID of the temporary security credential.

* `secret_access_key` - The secret access key of the temporary security credential.

* `security_token` - The security token of the temporary security credential.

* `expiration` - The expiration time of the temporary security credential in RFC3339 format.

<a name="subject_attr"></a>
The `subject` block supports:

* `namespace` - The namespace of the service account.

* `service_account` - The name of the service account.
