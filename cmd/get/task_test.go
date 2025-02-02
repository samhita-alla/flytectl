package get

import (
	"os"
	"testing"

	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	resourceListRequestTask *admin.ResourceListRequest
	objectGetRequestTask    *admin.ObjectGetRequest
	namedIDRequestTask      *admin.NamedEntityIdentifierListRequest
	taskListResponse        *admin.TaskList
	argsTask                []string
)

func getTaskSetup() {
	ctx = testutils.Ctx
	cmdCtx = testutils.CmdCtx
	mockClient = testutils.MockClient
	argsTask = []string{"task1"}
	sortedListLiteralType := core.Variable{
		Type: &core.LiteralType{
			Type: &core.LiteralType_CollectionType{
				CollectionType: &core.LiteralType{
					Type: &core.LiteralType_Simple{
						Simple: core.SimpleType_INTEGER,
					},
				},
			},
		},
	}
	variableMap := map[string]*core.Variable{
		"sorted_list1": &sortedListLiteralType,
		"sorted_list2": &sortedListLiteralType,
	}

	task1 := &admin.Task{
		Id: &core.Identifier{
			Name:    "task1",
			Version: "v1",
		},
		Closure: &admin.TaskClosure{
			CreatedAt: &timestamppb.Timestamp{Seconds: 0, Nanos: 0},
			CompiledTask: &core.CompiledTask{
				Template: &core.TaskTemplate{
					Interface: &core.TypedInterface{
						Inputs: &core.VariableMap{
							Variables: variableMap,
						},
					},
				},
			},
		},
	}

	task2 := &admin.Task{
		Id: &core.Identifier{
			Name:    "task1",
			Version: "v2",
		},
		Closure: &admin.TaskClosure{
			CreatedAt: &timestamppb.Timestamp{Seconds: 1, Nanos: 0},
			CompiledTask: &core.CompiledTask{
				Template: &core.TaskTemplate{
					Interface: &core.TypedInterface{
						Inputs: &core.VariableMap{
							Variables: variableMap,
						},
					},
				},
			},
		},
	}

	tasks := []*admin.Task{task2, task1}
	resourceListRequestTask = &admin.ResourceListRequest{
		Id: &admin.NamedEntityIdentifier{
			Project: projectValue,
			Domain:  domainValue,
			Name:    argsTask[0],
		},
		SortBy: &admin.Sort{
			Key:       "created_at",
			Direction: admin.Sort_DESCENDING,
		},
		Limit: 100,
	}

	taskListResponse = &admin.TaskList{
		Tasks: tasks,
	}

	objectGetRequestTask = &admin.ObjectGetRequest{
		Id: &core.Identifier{
			ResourceType: core.ResourceType_TASK,
			Project:      projectValue,
			Domain:       domainValue,
			Name:         argsTask[0],
			Version:      "v2",
		},
	}
	namedIDRequestTask = &admin.NamedEntityIdentifierListRequest{
		Project: projectValue,
		Domain:  domainValue,
		SortBy: &admin.Sort{
			Key:       "name",
			Direction: admin.Sort_ASCENDING,
		},
		Limit: 100,
	}

	var taskEntities []*admin.NamedEntityIdentifier
	idTask1 := &admin.NamedEntityIdentifier{
		Project: projectValue,
		Domain:  domainValue,
		Name:    "task1",
	}
	idTask2 := &admin.NamedEntityIdentifier{
		Project: projectValue,
		Domain:  domainValue,
		Name:    "task2",
	}
	taskEntities = append(taskEntities, idTask1, idTask2)
	namedIdentifierListTask := &admin.NamedEntityIdentifierList{
		Entities: taskEntities,
	}

	mockClient.OnListTasksMatch(ctx, resourceListRequestTask).Return(taskListResponse, nil)
	mockClient.OnGetTaskMatch(ctx, objectGetRequestTask).Return(task2, nil)
	mockClient.OnListTaskIdsMatch(ctx, namedIDRequestTask).Return(namedIdentifierListTask, nil)

	taskConfig.Latest = false
	taskConfig.ExecFile = ""
	taskConfig.Version = ""
}

func TestGetTaskFunc(t *testing.T) {
	setup()
	getTaskSetup()
	err = getTaskFunc(ctx, argsTask, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListTasks", ctx, resourceListRequestTask)
	tearDownAndVerify(t, `[
	{
		"id": {
			"name": "task1",
			"version": "v2"
		},
		"closure": {
			"compiledTask": {
				"template": {
					"interface": {
						"inputs": {
							"variables": {
								"sorted_list1": {
									"type": {
										"collectionType": {
											"simple": "INTEGER"
										}
									}
								},
								"sorted_list2": {
									"type": {
										"collectionType": {
											"simple": "INTEGER"
										}
									}
								}
							}
						}
					}
				}
			},
			"createdAt": "1970-01-01T00:00:01Z"
		}
	},
	{
		"id": {
			"name": "task1",
			"version": "v1"
		},
		"closure": {
			"compiledTask": {
				"template": {
					"interface": {
						"inputs": {
							"variables": {
								"sorted_list1": {
									"type": {
										"collectionType": {
											"simple": "INTEGER"
										}
									}
								},
								"sorted_list2": {
									"type": {
										"collectionType": {
											"simple": "INTEGER"
										}
									}
								}
							}
						}
					}
				}
			},
			"createdAt": "1970-01-01T00:00:00Z"
		}
	}
]`)
}

func TestGetTaskFuncLatest(t *testing.T) {
	setup()
	getTaskSetup()
	taskConfig.Latest = true
	err = getTaskFunc(ctx, argsTask, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListTasks", ctx, resourceListRequestTask)
	tearDownAndVerify(t, `{
	"id": {
		"name": "task1",
		"version": "v2"
	},
	"closure": {
		"compiledTask": {
			"template": {
				"interface": {
					"inputs": {
						"variables": {
							"sorted_list1": {
								"type": {
									"collectionType": {
										"simple": "INTEGER"
									}
								}
							},
							"sorted_list2": {
								"type": {
									"collectionType": {
										"simple": "INTEGER"
									}
								}
							}
						}
					}
				}
			}
		},
		"createdAt": "1970-01-01T00:00:01Z"
	}
}`)
}

func TestGetTaskWithVersion(t *testing.T) {
	setup()
	getTaskSetup()
	taskConfig.Version = "v2"
	objectGetRequestTask.Id.ResourceType = core.ResourceType_TASK
	err = getTaskFunc(ctx, argsTask, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "GetTask", ctx, objectGetRequestTask)
	tearDownAndVerify(t, `{
	"id": {
		"name": "task1",
		"version": "v2"
	},
	"closure": {
		"compiledTask": {
			"template": {
				"interface": {
					"inputs": {
						"variables": {
							"sorted_list1": {
								"type": {
									"collectionType": {
										"simple": "INTEGER"
									}
								}
							},
							"sorted_list2": {
								"type": {
									"collectionType": {
										"simple": "INTEGER"
									}
								}
							}
						}
					}
				}
			}
		},
		"createdAt": "1970-01-01T00:00:01Z"
	}
}`)
}

func TestGetTasks(t *testing.T) {
	setup()
	getTaskSetup()
	argsTask = []string{}
	err = getTaskFunc(ctx, argsTask, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListTaskIds", ctx, namedIDRequest)
	tearDownAndVerify(t, `[
	{
		"project": "dummyProject",
		"domain": "dummyDomain",
		"name": "task1"
	},
	{
		"project": "dummyProject",
		"domain": "dummyDomain",
		"name": "task2"
	}
]`)
}

func TestGetTaskWithExecFile(t *testing.T) {
	setup()
	getTaskSetup()
	taskConfig.Version = "v2"
	taskConfig.ExecFile = testDataFolder + "task_exec_file"
	err = getTaskFunc(ctx, argsTask, cmdCtx)
	os.Remove(taskConfig.ExecFile)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "GetTask", ctx, objectGetRequestTask)
	tearDownAndVerify(t, `{
	"id": {
		"name": "task1",
		"version": "v2"
	},
	"closure": {
		"compiledTask": {
			"template": {
				"interface": {
					"inputs": {
						"variables": {
							"sorted_list1": {
								"type": {
									"collectionType": {
										"simple": "INTEGER"
									}
								}
							},
							"sorted_list2": {
								"type": {
									"collectionType": {
										"simple": "INTEGER"
									}
								}
							}
						}
					}
				}
			}
		},
		"createdAt": "1970-01-01T00:00:01Z"
	}
}`)
}
