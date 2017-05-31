package hostedzones

// HostedZones holds the Hosted Zone IDs for creating guest cluster recordsets.
type HostedZones struct {
	// API is the Hosted Zone ID for the API recordset.
	API string
	// Etcd is the Hosted Zone ID for the Etcd recordset.
	Etcd string
	// Ingress is the Hosted Zone ID for the Ingress recordset.
	Ingress string
}
