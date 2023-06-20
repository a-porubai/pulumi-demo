package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Mocks struct {
	ResourceOutputs map[string]func(outputs resource.PropertyMap) resource.PropertyMap
	CallArgs        map[string]func(args pulumi.MockCallArgs) (resource.PropertyMap, error)
}

func (m *Mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	outputs := args.Inputs

	if f, ok := m.ResourceOutputs[args.TypeToken+"::"+args.Name]; ok {
		outputs = f(outputs)
	}

	return args.Name + "_id", outputs, nil
}

func (m *Mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	if f, ok := m.CallArgs[args.Token]; ok {
		return f(args)
	}

	return args.Args, nil
}

func WithMocksAndConfig(project, stack string, config map[string]string, mocks pulumi.MockResourceMonitor) pulumi.RunOption {
	return func(info *pulumi.RunInfo) {
		info.Project, info.Stack, info.Mocks, info.Config = project, stack, mocks, config
	}
}
