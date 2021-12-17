resource "huaweicloud_vpc_address_group" "test" {
	dry_run = false
	name = "test1"
	ip_version = 4
	description  =  "vpc test"
	ip_set	=	[
		"192.168.5.0/24",
		"192.168.9.0/24"
	]
}