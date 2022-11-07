
## Links

[GO SDK](https://github.com/gophercloud/gophercloud)

[Terraform provider](https://github.com/terraform-provider-openstack/terraform-provider-openstack)

[openstack CLI](https://docs.openstack.org/python-openstackclient/latest/cli/index.html)


## Release

https://unix.stackexchange.com/questions/155046/determine-if-git-working-directory-is-clean-from-a-script

git tag v0.1.0
git push origin v0.1.0

To delete:

git tag -d X.X.X

export GITHUB_TOKEN=....

goreleaser --rm-dist release

If git status is not clean:

goreleaser --rm-dist release --skip-validate
