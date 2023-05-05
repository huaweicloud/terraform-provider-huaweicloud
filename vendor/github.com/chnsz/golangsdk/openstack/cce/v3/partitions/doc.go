// Package partitions is a new feature introduced by CCE Turbo remote distributed cluster management.
// Currently, partitions are divided into three categories: Center, HomeZone, and IES.
//   - Default: Indicates the availability zone resources of the HuaweiCloud data center
//   - HomeZone: Indicates the edge computing area in the user's home, which is not currently supported
//   - IES: IES is the edge station of HuaweiCloud, which deploys user-specific cloud resources to the
//          local computer zone
// The concept of partitions can support CCE Turbo to manage distributed nodes, which has
// strong applicability in the field of edge computing.
package partitions
