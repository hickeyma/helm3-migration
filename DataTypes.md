Diff of datatypes from v2 to v3

Release:

```console
$ pr -w $COLUMNS -m -t release-v2.txt release-v3.txt
Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`					       Name string `json:"name,omitempty"`
Info *Info `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`					       Info *Info `json:"info,omitempty"`
Chart *chart.Chart `protobuf:"bytes,3,opt,name=chart,proto3" json:"chart,omitempty"`				       Chart *chart.Chart `json:"chart,omitempty"`
Config *chart.Config `protobuf:"bytes,4,opt,name=config,proto3" json:"config,omitempty"`			       Config map[string]interface{} `json:"config,omitempty"`
Manifest string `protobuf:"bytes,5,opt,name=manifest,proto3" json:"manifest,omitempty"`				       Manifest string `json:"manifest,omitempty"`
Hooks []*Hook `protobuf:"bytes,6,rep,name=hooks,proto3" json:"hooks,omitempty"`					       Hooks []*Hook `json:"hooks,omitempty"`
Version int32 `protobuf:"varint,7,opt,name=version,proto3" json:"version,omitempty"`				       Version int `json:"version,omitempty"`
Namespace string `protobuf:"bytes,8,opt,name=namespace,proto3" json:"namespace,omitempty"`			       Namespace string `json:"namespace,omitempty"`
```

Info:

```console
$ pr -w $COLUMNS -m -t info-v2.txt info-v3.txt
FirstDeployed *timestamp.Timestamp `protobuf:"bytes,2,opt,name=first_deployed,json=firstDeployed,proto3" json:"first_deployed,omitempty"`  FirstDeployed time.Time `json:"first_deployed,omitempty"`
LastDeployed  *timestamp.Timestamp `protobuf:"bytes,3,opt,name=last_deployed,json=lastDeployed,proto3" json:"last_deployed,omitempty"`	   LastDeployed time.Time `json:"last_deployed,omitempty"`
Deleted *timestamp.Timestamp `protobuf:"bytes,4,opt,name=deleted,proto3" json:"deleted,omitempty"`					   Deleted time.Time `json:"deleted,omitempty"`
Description string `protobuf:"bytes,5,opt,name=Description,proto3" json:"Description,omitempty"`					   Description string `json:"Description,omitempty"`
Status *Status `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`							   Status Status `json:"status,omitempty"`
																	   Resources string `json:"resources,omitempty"`
																	   Notes string `json:"notes,omitempty"`
																	   LastTestSuiteRun *TestSuite `json:"last_test_suite_run,omitempty"`
```

Status:

```console
$ pr -w $COLUMNS -m -t status-v2.txt status-v3.txt
Code Status_Code `protobuf:"varint,1,opt,name=code,proto3,enum=hapi.release.Status_Code" json:"code,omitempty"`							       string
Resources string `protobuf:"bytes,3,opt,name=resources,proto3" json:"resources,omitempty"`									       
Notes string `protobuf:"bytes,4,opt,name=notes,proto3" json:"notes,omitempty"`											       
LastTestSuiteRun *TestSuite `protobuf:"bytes,5,opt,name=last_test_suite_run,json=lastTestSuiteRun,proto3" json:"last_test_suite_run,omitempty"`
```

Status_Code:

```console
$ cat status_code-v2.txt
int32
```

TestSuite:

```console
$ pr -w $COLUMNS -m -t test_suite-v2.txt test_suite-v3.txt
StartedAt *timestamp.Timestamp `protobuf:"bytes,1,opt,name=started_at,json=startedAt,proto3" json:"started_at,omitempty"`					       StartedAt time.Time `json:"started_at,omitempty"`
CompletedAt *timestamp.Timestamp `protobuf:"bytes,2,opt,name=completed_at,json=completedAt,proto3" json:"completed_at,omitempty"`				       CompletedAt time.Time `json:"completed_at,omitempty"`
Results []*TestRun `protobuf:"bytes,3,rep,name=results,proto3" json:"results,omitempty"`									       Results []*TestRun `json:"results,omitempty"`
```

TestRun:

```console
$ pr -w $COLUMNS -m -t test_run-v2.txt test_run-v3.txt
Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`											       Name string `json:"name,omitempty"`
Status TestRun_Status `protobuf:"varint,2,opt,name=status,proto3,enum=hapi.release.TestRun_Status" json:"status,omitempty"`					       Status TestRunStatus `json:"status,omitempty"`
Info string `protobuf:"bytes,3,opt,name=info,proto3" json:"info,omitempty"`											       Info string `json:"info,omitempty"`
StartedAt *timestamp.Timestamp `protobuf:"bytes,4,opt,name=started_at,json=startedAt,proto3" json:"started_at,omitempty"`					       StartedAt   time.Time `json:"started_at,omitempty"`
CompletedAt *timestamp.Timestamp `protobuf:"bytes,5,opt,name=completed_at,json=completedAt,proto3" json:"completed_at,omitempty"`				       CompletedAt time.Time `json:"completed_at,omitempty"`
```

TestRunStatus:

```console
$ pr -w $COLUMNS -m -t test_run_status-v2.txt test_run_status-v3.txt
 int32																					string
```

Chart:

```console
$ pr -w $COLUMNS -m -t chart-v2.txt chart-v3.txt
Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`										Metadata *Metadata
Templates []*Template `protobuf:"bytes,2,rep,name=templates,proto3" json:"templates,omitempty"`										Templates []*File
Dependencies []*Chart `protobuf:"bytes,3,rep,name=dependencies,proto3" json:"dependencies,omitempty"`									dependencies []*Chart
Values *Config `protobuf:"bytes,4,opt,name=values,proto3" json:"values,omitempty"`											Values map[string]interface{}
Files []*any.Any `protobuf:"bytes,5,rep,name=files,proto3" json:"files,omitempty"`											Files []*File
																					Schema []byte
																					Lock *Lock
																					parent *Chart
```

Metadata:

```console
$ pr -w $COLUMNS -m -t metadata-v2.txt metadata-v3.txt
Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`												Name string `json:"name,omitempty"`
Home string `protobuf:"bytes,2,opt,name=home,proto3" json:"home,omitempty"`												Home string `json:"home,omitempty"`
Sources []string `protobuf:"bytes,3,rep,name=sources,proto3" json:"sources,omitempty"`											Sources []string `json:"sources,omitempty"`
Version string `protobuf:"bytes,4,opt,name=version,proto3" json:"version,omitempty"`											Version string `json:"version,omitempty"`
Description string `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`									Description string `json:"description,omitempty"`
Keywords []string `protobuf:"bytes,6,rep,name=keywords,proto3" json:"keywords,omitempty"`										Keywords []string `json:"keywords,omitempty"`
Maintainers []*Maintainer `protobuf:"bytes,7,rep,name=maintainers,proto3" json:"maintainers,omitempty"`									Maintainers []*Maintainer `json:"maintainers,omitempty"`
Icon string `protobuf:"bytes,9,opt,name=icon,proto3" json:"icon,omitempty"`												Icon string `json:"icon,omitempty"`
ApiVersion string `protobuf:"bytes,10,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`										APIVersion string `json:"apiVersion,omitempty"`
Condition string `protobuf:"bytes,11,opt,name=condition,proto3" json:"condition,omitempty"`										Condition string `json:"condition,omitempty"`
Tags string `protobuf:"bytes,12,opt,name=tags,proto3" json:"tags,omitempty"`												Tags string `json:"tags,omitempty"`
AppVersion string `protobuf:"bytes,13,opt,name=appVersion,proto3" json:"appVersion,omitempty"`										AppVersion string `json:"appVersion,omitempty"`
Deprecated bool `protobuf:"varint,14,opt,name=deprecated,proto3" json:"deprecated,omitempty"`										Deprecated bool `json:"deprecated,omitempty"`
Annotations map[string]string `protobuf:"bytes,16,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"by Annotations map[string]string `json:"annotations,omitempty"`
KubeVersion string `protobuf:"bytes,17,opt,name=kubeVersion,proto3" json:"kubeVersion,omitempty"`									KubeVersion string `json:"kubeVersion,omitempty"`
Engine string `protobuf:"bytes,8,opt,name=engine,proto3" json:"engine,omitempty"`											
TillerVersion string `protobuf:"bytes,15,opt,name=tillerVersion,proto3" json:"tillerVersion,omitempty"`									
																					Dependencies []*Dependency `json:"dependencies,omitempty"`
																					Type string `json:"type,omitempty"`
															
```

Maintainer:

```console
$ pr -w $COLUMNS -m -t maintainer-v2.txt maintainer-v3.txt
Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`												Name string `json:"name,omitempty"`
Email string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`												Email string `json:"email,omitempty"`
Url string `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`												URL string `json:"url,omitempty"`
```

Dependency (v3 only):

```console
$ cat dependency-v3.txt
Name string `json:"name"`
Version string `json:"version,omitempty"`
Repository string `json:"repository"`
Condition string `json:"condition,omitempty"`
Tags []string `json:"tags,omitempty"`
Enabled bool `json:"enabled,omitempty"`
ImportValues []interface{} `json:"import-values,omitempty"`
Alias string `json:"alias,omitempty"`
```

Lock (v3 only):

```console
$ cat lock-v3.txt
Generated time.Time `json:"generated"`
Digest string `json:"digest"`
Dependencies []*Dependency `json:"dependencies"`
```

File (v3 only):

```console
$ cat file-v3.txt
Name string
Data []byte
```

Hook:

```console
$ pr -w $COLUMNS -m -t hook-v2.txt hook-v3.txt
Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`												Name string `json:"name,omitempty"`
Kind string `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`												Kind string `json:"kind,omitempty"`
Path string `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`												Path string `json:"path,omitempty"`
Manifest string `protobuf:"bytes,4,opt,name=manifest,proto3" json:"manifest,omitempty"`											Manifest string `json:"manifest,omitempty"`
Events []Hook_Event `protobuf:"varint,5,rep,packed,name=events,proto3,enum=hapi.release.Hook_Event" json:"events,omitempty"`						Events []HookEvent `json:"events,omitempty"`
LastRun *timestamp.Timestamp `protobuf:"bytes,6,opt,name=last_run,json=lastRun,proto3" json:"last_run,omitempty"`							LastRun time.Time `json:"last_run,omitempty"`
Weight int32 `protobuf:"varint,7,opt,name=weight,proto3" json:"weight,omitempty"`											Weight int `json:"weight,omitempty"`
DeletePolicies []Hook_DeletePolicy `protobuf:"varint,8,rep,packed,name=delete_policies,json=deletePolicies,proto3,enum=hapi.release.Hook_DeletePolicy" json:"delete_pol DeletePolicies []HookDeletePolicy `json:"delete_policies,omitempty"`
```

HookEvent:

```console
$ pr -w $COLUMNS -m -t hook_event-v2.txt hook_event-v3.txt
int32																					string
```

HookDeletePolicy:

```console
$ pr -w $COLUMNS -m -t hook_delete_policy-v2.txt hook_delete_policy-v3.txt
int32																					string
```
