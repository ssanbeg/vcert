package policy

const (
	JsonExtention        = ".json"
	YamlExtention        = ".yaml"
	RootPath             = "\\VED\\Policy\\"
	PolicyClass          = "Policy"
	PolicyAttributeClass = "X509 Certificate"

	//tpp policy attributes
	TppContact                    = "Contact"
	TppApprover                   = "Approver"
	TppCertificateAuthority       = "Certificate Authority"
	TppProhibitWildcard           = "Prohibit Wildcard"
	TppDomainSuffixWhitelist      = "Domain Suffix Whitelist"
	TppOrganization               = "Organization"
	TppOrganizationalUnit         = "Organizational Unit"
	TppCity                       = "City"
	TppState                      = "State"
	TppCountry                    = "Country"
	TppKeyAlgorithm               = "Key Algorithm"
	TppKeyBitStrength             = "Key Bit Strength"
	TppEllipticCurve              = "Elliptic Curve"
	ServiceGenerated              = "Manual Csr"
	TppProhibitedSANTypes         = "Prohibited SAN Types"
	TppAllowPrivateKeyReuse       = "Allow Private Key Reuse"
	TppWantRenewal                = "Want Renewal"
	TppDnsAllowed                 = "DNS"
	TppIpAllowed                  = "IP"
	TppEmailAllowed               = "Email"
	TppUriAllowed                 = "URI"
	TppUpnAllowed                 = "UPN"
	AllowAll                      = ".*"
	UserProvided                  = "UserProvided"
	DefaultCA                     = "BUILTIN\\Built-In CA\\Default Product"
	TppManagementType             = "Management Type"
	TppManagementTypeEnrollment   = "Enrollment"
	TppManagementTypeProvisioning = "Provisioning"
	CloudEntrustCA                = "ENTRUST"
	CloudDigicertCA               = "DIGICERT"
	CloudRequesterName            = "Venafi Cloud Service"
	CloudRequesterEmail           = "no-reply@venafi.cloud"
	CloudRequesterPhone           = "801-555-0123"
)
