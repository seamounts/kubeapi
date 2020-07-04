/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"github.com/seamounts/kubeapi/pkg/model/config"
	"github.com/seamounts/kubeapi/pkg/plugin"
)

const (
	pluginName    = "go" + plugin.DefaultNameQualifier
	pluginVersion = "v1.0.0"
)

var supportedProjectVersions = []string{config.Version1}

var (
	_ plugin.Base                  = Plugin{}
	_ plugin.InitPluginGetter      = Plugin{}
	_ plugin.CreateAPIPluginGetter = Plugin{}
)

type Plugin struct {
	initPlugin
	createAPIPlugin
}

func (Plugin) Name() string                           { return pluginName }
func (Plugin) Version() string                        { return pluginVersion }
func (Plugin) SupportedProjectVersions() []string     { return supportedProjectVersions }
func (p Plugin) GetInitPlugin() plugin.Init           { return &p.initPlugin }
func (p Plugin) GetCreateAPIPlugin() plugin.CreateAPI { return &p.createAPIPlugin }
