# 1. Overview

Helm v3 introduces quite a lot of change in the underlying architecture and plumbing from the previous release, 
Helm v2. One key change is around Release storage. The changes includes the Kubernetes resource for storage and the 
release object metadata contained in the resource. Releases will also be on a per user namespace instead of using the the Tiller namespace (for example, v2 default Tiller namespace `kube-system`).

# 2. Requirement

When Helm v3 is installed in a cluster that is currently managed by a Helm v2 installation, the Helm v3 installation 
should be able to manage the existing v2 Releases.

Note: This proposal covers the migration use case of Helm v3 managing existing Helm v2 releases (i.e. converting v2 releases to v3 releases). Other migration use cases are covered by documentation which is currently a WIP: https://github.com/helm/helm/pull/5582.

# 3. Proposal

1. Get v2 [release history](https://github.com/helm/helm/blob/master/pkg/helm/client.go#L273)
2. For each release version:
   - Map v2 release version object to v3 release object
   - [Add this v3 release object to Helm v3](https://github.com/helm/helm/blob/dev-v3/pkg/action/upgrade.go#L69)
3. [Delete v2 release history](https://github.com/helm/helm/blob/master/pkg/helm/client.go#L122)

Note:
- The namespace is based off each v2 release version namespace
- When adding v3 release version, only the release object is added and not the underlying kubernetes resources
- When deleting v2 release , only the release (and history) is deleted and not the underlying kubernetes resources
- It is assumed that all v2 release versions are incremental and inorder, and that the versions can be added incrementally as well

## 3.1 Implementation

A standalone migration tool that migrates from Helm v2 to Helm v3. (@prydonius 
https://github.com/helm/community/issues/67#issuecomment-448033387) 

The primary function of the tool is to:

- Automatically back up Helm v2 Release and convert them to Helm v3 Release

The suggestion is for a simple, Helm-org supported plugin named `helm 2to3`. The plugin should concentrate at the 
start on its primary function of converting releases from v2 to v3 through the `convert` subcommand. It should be able 
to be extended if need be. (@jdolitsky https://github.com/helm/community/issues/67#issuecomment-448045222)

```console
$ helm 2to3 convert myrelease --dry-run

NOTE: This is in dry-run mode, the following actions will not be executed.
Run without --dry-run to take the actions described below:

Release "myrelease" will be converted from Helm 2 to Helm 3. 
[Helm 3] Release "myrelease" will be created.
[Helm 3] ReleaseVersion "myrelease.v1" will be created.
[Helm 3] ReleaseVersion "myrelease.v2" will be created.
[Helm 3] ReleaseVersion "myrelease.v3" will be created.
[Helm 2] ReleaseVersion "myrelease" will be deleted.

$ helm 2to3 convert myrelease

Release "myrelease" will be converted from Helm 2 to Helm 3. 
[Helm 3] Release "myrelease" will be created.
[Helm 3] ReleaseVersion "myrelease.v1" will be created.
[Helm 3] ReleaseVersion "myrelease.v1" created.
[Helm 3] ReleaseVersion "myrelease.v2" will be created.
[Helm 3] ReleaseVersion "myrelease.v2" created.
[Helm 3] ReleaseVersion "myrelease.v3" will be created.
[Helm 3] ReleaseVersion "myrelease.v3" created.
[Helm 3] Release "myrelease" created.
[Helm 2] Release "myrelease" will be deleted.
[Helm 2] Release "myrelease" deleted.
Release "myrelease" was converted successfully from Helm 2 to Helm 3. 
```

## 3.2 Helm v2 and v3 SDK updates

The following changes are required in the Helm v2 and v3 SDKs:

- When adding a release in v3, be able to set the following to correspond to the v2 details:
  - `modifiedAt`
  - `status`
  - `creationTimestamp`
  - `description`
- A flag (e.g. `--state-only`) in [v3 release upgrade](https://github.com/helm/helm/blob/dev-v3/pkg/action/upgrade.go#L69) which only adds the release object and not the undelying kubernetes resources
- A flag (e.g. `--state-only`) in [v2 release delete](https://github.com/helm/helm/blob/master/pkg/helm/client.go#L122) which only deletes the release object (and history) and not the undelying kubernetes resources
- v3 `helm upgrade` `--install` flag seems to be disabled
- Be able to use Go modules for building the tool as it allows [multiple major versions (helm 2 and 3)](https://github.com/golang/go/wiki/Modules#v2-modules-allow-multiple-major-versions-within-a-single-build) to be used. Need @mattfarina [v3 Go dependency PR](https://github.com/helm/helm/pull/5498) to merge

# 4. Reference

- Migration was raised at *KubeCon/CloudNativeCon Seattle 2018* at the *Helm Deep Dive session*. 
Ref: https://github.com/helm/community/issues/67.
