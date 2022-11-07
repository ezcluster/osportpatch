
## Links

[GO SDK](https://github.com/gophercloud/gophercloud)

[Terraform provider](https://github.com/terraform-provider-openstack/terraform-provider-openstack)

[openstack CLI](https://docs.openstack.org/python-openstackclient/latest/cli/index.html)


## Release

https://unix.stackexchange.com/questions/155046/determine-if-git-working-directory-is-clean-from-a-script

git tag v0.1.0

To delete:
git tag -d v0.1.0
git push --delete origin v0.1.0

export GITHUB_TOKEN=....

goreleaser  release --rm-dist

If git status is not clean:

goreleaser release --rm-dist --skip-validate
