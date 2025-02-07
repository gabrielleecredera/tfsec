package elasticache

import (
	"github.com/aquasecurity/defsec/provider/aws/elasticache"
	"github.com/aquasecurity/tfsec/internal/pkg/block"
)

func Adapt(modules block.Modules) elasticache.ElastiCache {
	return elasticache.ElastiCache{
		Clusters:          adaptClusters(modules),
		ReplicationGroups: adaptReplicationGroups(modules),
		SecurityGroups:    adaptSecurityGroups(modules),
	}
}
func adaptClusters(modules block.Modules) []elasticache.Cluster {
	var clusters []elasticache.Cluster
	for _, module := range modules {
		for _, resource := range module.GetResourcesByType("aws_elasticache_cluster") {
			clusters = append(clusters, adaptCluster(resource))
		}
	}
	return clusters
}

func adaptReplicationGroups(modules block.Modules) []elasticache.ReplicationGroup {
	var replicationGroups []elasticache.ReplicationGroup
	for _, module := range modules {
		for _, resource := range module.GetResourcesByType("aws_elasticache_replication_group") {
			replicationGroups = append(replicationGroups, adaptReplicationGroup(resource))
		}
	}
	return replicationGroups
}

func adaptSecurityGroups(modules block.Modules) []elasticache.SecurityGroup {
	var securityGroups []elasticache.SecurityGroup
	for _, module := range modules {
		for _, resource := range module.GetResourcesByType("aws_elasticache_security_group") {
			securityGroups = append(securityGroups, adaptSecurityGroup(resource))
		}
	}
	return securityGroups
}

func adaptCluster(resource *block.Block) elasticache.Cluster {
	engineAttr := resource.GetAttribute("engine")
	engineVal := engineAttr.AsStringValueOrDefault("", resource)

	nodeTypeAttr := resource.GetAttribute("node_type")
	nodeTypeVal := nodeTypeAttr.AsStringValueOrDefault("", resource)

	snapshotRetentionAttr := resource.GetAttribute("snapshot_retention_limit")
	snapshotRetentionVal := snapshotRetentionAttr.AsIntValueOrDefault(0, resource)

	return elasticache.Cluster{
		Metadata:               resource.Metadata(),
		Engine:                 engineVal,
		NodeType:               nodeTypeVal,
		SnapshotRetentionLimit: snapshotRetentionVal,
	}
}

func adaptReplicationGroup(resource *block.Block) elasticache.ReplicationGroup {
	transitEncryptionAttr := resource.GetAttribute("transit_encryption_enabled")
	transitEncryptionVal := transitEncryptionAttr.AsBoolValueOrDefault(false, resource)

	atRestEncryptionAttr := resource.GetAttribute("at_rest_encryption_enabled")
	atRestEncryptionVal := atRestEncryptionAttr.AsBoolValueOrDefault(false, resource)

	return elasticache.ReplicationGroup{
		Metadata:                 resource.Metadata(),
		TransitEncryptionEnabled: transitEncryptionVal,
		AtRestEncryptionEnabled:  atRestEncryptionVal,
	}
}

func adaptSecurityGroup(resource *block.Block) elasticache.SecurityGroup {
	descriptionAttr := resource.GetAttribute("description")
	descriptionVal := descriptionAttr.AsStringValueOrDefault("Managed by Terraform", resource)

	return elasticache.SecurityGroup{
		Metadata:    resource.Metadata(),
		Description: descriptionVal,
	}
}
