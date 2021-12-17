resource "huaweicloud_vpc_secgroup_v3" "test"  {
  region = "cn-north-4"
  dry_run = false
  name = "aaa"
  description = "123"
  enterprise_project_id = "0"
}
//For updating
resource "huaweicloud_vpc_secgroup_v3" "test"  {
  region = "cn-north-4"
  dry_run = false
  name = "a"
  description = "12"
}