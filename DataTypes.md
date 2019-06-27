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
 pr -w $COLUMNS -m -t test_suite-v2.txt test_suite-v3.txt
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
pr -w $COLUMNS -m -t chart-v2.txt chart-v3.txt
Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`										Metadata *Metadata
Templates []*Template `protobuf:"bytes,2,rep,name=templates,proto3" json:"templates,omitempty"`										Templates []*File
Dependencies []*Chart `protobuf:"bytes,3,rep,name=dependencies,proto3" json:"dependencies,omitempty"`									dependencies []*Chart
Values *Config `protobuf:"bytes,4,opt,name=values,proto3" json:"values,omitempty"`											Values map[string]interface{}
Files []*any.Any `protobuf:"bytes,5,rep,name=files,proto3" json:"files,omitempty"`											Files []*File
																					Schema []byte
																					Lock *Lock
																					parent *Chart
```
