# osportpatch

This small utility is intended to manipulate the `allowed_address_pairs` field of the port associated to an openstack VM.

This may be required by some tools implementing a virtual IP system, to ensure fail over between several servers, such as [metallb](https://metallb.universe.tf/) for a kubernetes infrastructure.

## Configuration

All openstack access parameters must be provided as environment variables, as the following:

```
export OS_USERNAME="john"
export OS_PROJECT_NAME="projectX"
export OS_PASSWORD="thepasswordofjohn"
export OS_AUTH_URL="http://keystone.openstack.mycompany.com:5000"
export OS_REGION_NAME="RegionOne"
```

## USAGE:

```
osportpatch [add|remove] <ipaddr> <serverName> [<serverName>...]
```
