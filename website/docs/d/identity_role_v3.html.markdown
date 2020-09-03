---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_role_v3"
sidebar_current: "docs-huaweicloud-datasource-identity-role-v3"
description: |-
  Get information on an HuaweiCloud Role.
---

# huaweicloud\_identity\_role\_v3

Use this data source to get the ID of an HuaweiCloud role.

The Role in Terraform is the same as Policy on console. however,
The policy name is the display name of Role, the Role name cannot
be found on Console. please refer to the following table to configuration
Role:

Role Name | Policy Name
---- | ---
readonly | Tenant Guest
tms_adm | TMS Administrator
cce_adm | CCE Administrator
dcs_admin | DCS Administrator
dis_adm | DIS Administrator
system_all_6 | VPC Viewer
rds_adm | RDS Administrator
system_all_1001 | Full Access
system_all_3 | EVS Viewer
te_agency | Agent Operator
dms_adm | DMS Administrator
ces_adm | CES Administrator
rts_adm | RTS Administrator
system_all_5 | VPC Admin
dns_adm | DNS Administrator
server_adm | Server Administrator
sdrs_adm | SDRS Administrator
system_all_0 | ECS Admin
wks_adm | Workspace Administrator
te_admin | Tenant Administrator
sfs_adm | SFS Administrator
vpc_netadm | VPC Administrator
css_adm | CSS Administrator
as_adm | AutoScaling Administrator
csbs_adm | CSBS Administrator
secu_admin | Security Administrator
system_all_2 | ECS Viewer
dws_adm | DWS Administrator
mobs_adm | MaaS OBS  Administrator
vbs_adm | VBS Administrator
ddos_adm | Anti-DDoS Administrator
system_all_4 | EVS Admin
system_all_1 | ECS User
dws_db_acc | DWS Database Access
kms_adm | KMS Administrator
mrs_adm | MRS Administrator
nat_adm | NAT Gateway Administrator
dds_adm | DDS Administrator
ims_adm | IMS Administrator
smn_adm | SMN Administrator
plas_adm | Config Plas Connector
elb_adm | ELB Administrator


```hcl
data "huaweicloud_identity_role_v3" "auth_admin" {
  name = "secu_admin"
}
```

## Argument Reference

* `name` - (Required) The name of the role.

* `domain_id` - (Optional) The domain the role belongs to.

## Attributes Reference

`id` is set to the ID of the found role. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `domain_id` - See Argument Reference above.
* `region` - See Argument Reference above.
