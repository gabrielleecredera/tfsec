package documentdb

import (
	"github.com/aquasecurity/defsec/provider/aws/documentdb"
	"github.com/aquasecurity/defsec/types"
	"github.com/aquasecurity/tfsec/internal/pkg/block"
)

func Adapt(modules block.Modules) documentdb.DocumentDB {
	return documentdb.DocumentDB{
		Clusters: adaptClusters(modules),
	}
}

func adaptClusters(modules block.Modules) []documentdb.Cluster {
	var clusters []documentdb.Cluster
	for _, module := range modules {
		for _, resource := range module.GetResourcesByType("aws_docdb_cluster") {
			clusters = append(clusters, adaptCluster(resource, module))
		}
	}
	return clusters
}

func adaptCluster(resource *block.Block, module *block.Module) documentdb.Cluster {
	identifierAttr := resource.GetAttribute("cluster_identifier")
	identifierVal := identifierAttr.AsStringValueOrDefault("", resource)

	var enabledLogExports []types.StringValue
	var instances []documentdb.Instance

	enabledLogExportsAttr := resource.GetAttribute("enabled_cloudwatch_logs_exports")
	enabledLogExportsList := enabledLogExportsAttr.ValueAsStrings()
	for _, logExport := range enabledLogExportsList {
		enabledLogExports = append(enabledLogExports, types.String(logExport, enabledLogExportsAttr.Metadata()))
	}

	instancesRes := module.GetReferencingResources(resource, "aws_docdb_cluster_instance", "cluster_identifier")
	for _, instanceRes := range instancesRes {
		keyIDAttr := instanceRes.GetAttribute("kms_key_id")
		keyIDVal := keyIDAttr.AsStringValueOrDefault("", instanceRes)

		instances = append(instances, documentdb.Instance{
			Metadata: instanceRes.Metadata(),
			KMSKeyID: keyIDVal,
		})
	}

	storageEncryptedAttr := resource.GetAttribute("storage_encrypted")
	storageEncryptedVal := storageEncryptedAttr.AsBoolValueOrDefault(false, resource)

	KMSKeyIDAttr := resource.GetAttribute("kms_key_id")
	KMSKeyIDVal := KMSKeyIDAttr.AsStringValueOrDefault("", resource)

	return documentdb.Cluster{
		Metadata:          resource.Metadata(),
		Identifier:        identifierVal,
		EnabledLogExports: enabledLogExports,
		Instances:         instances,
		StorageEncrypted:  storageEncryptedVal,
		KMSKeyID:          KMSKeyIDVal,
	}
}
