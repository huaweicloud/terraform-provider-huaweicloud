package huaweicloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jen20/awspolicyequivalence"
)

func suppressEquivalentAwsPolicyDiffs(k, old, new string, d *schema.ResourceData) bool {
	equivalent, err := awspolicy.PoliciesAreEquivalent(old, new)
	if err != nil {
		return false
	}

	return equivalent
}

// Suppress all changes?
func suppressDiffAll(k, old, new string, d *schema.ResourceData) bool {
	return true
}

// Suppress changes if we get a computed min_disk_gb if value is unspecified (default 0)
func suppressMinDisk(k, old, new string, d *schema.ResourceData) bool {
	return new == "0" || old == new
}
