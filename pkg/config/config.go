package config

import (
	"fmt"
	"io/ioutil"
	"k8s.io/api/core/v1"
	"os"
)

type Config struct {
	Version  string
	Metadata Metadata
	Spec     Spec
}

const (
	PersistenceStorageClassType = "storageClass"
	PersistenceNfsType          = "nfs"
	PersistenceHostPathType     = "hostPath"
)

func (c *Config) GetStorageClassName() string {
	return c.Spec.Persistence.StorageClassName
}

func (c *Config) IgnorePv() bool {
	if c.GetStorageClassName() == "" {
		return false
	}
	// todo :  get storage class and get nfs server how to do?
	// from now just support nfs
	if c.Spec.Persistence.Nfs.Server == "" {
		return true
	}
	return false
}

func (c *Config) GetResource(key string) *Resource {
	if c == nil {
		return nil
	}
	if c.Spec.Resources == nil {
		return nil
	}
	if val, ok := c.Spec.Resources[key]; ok {
		return val
	}
	return nil
}

func (c *Config) GetHelmValuesTpl(key string) ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	dir := c.Spec.HelmConfig.Values.Dir
	if dir == "" {
		dir = "values"
	}

	valuesFilepath := fmt.Sprintf("%s/%s.yaml", dir, key)

	_, err := os.Stat(valuesFilepath)
	if err == nil {
		return ioutil.ReadFile(valuesFilepath)
	}
	if os.IsNotExist(err) {
		return nil, nil
	}
	return nil, err
}

type Metadata struct {
	Name      string
	Namespace string
}

type Spec struct {
	Persistence Persistence
	Resources   map[string]*Resource
	HelmConfig  HelmConfig `yaml:"helm"`
}

type HelmConfig struct {
	Values ValuesConfig `yaml:"values"`
}

type ValuesConfig struct {
	Dir string `yaml:"dir"`
}

type Persistence struct {
	Nfs              `yaml:"nfs"`
	HostPath         `yaml:"hostPath"`
	StorageClassName string `yaml:"storageClassName"`
	Type             string
	AccessModes      []v1.PersistentVolumeAccessMode `yaml:"accessModes"`
}

type Nfs struct {
	Server   string
	RootPath string `yaml:"rootPath"`
}

type HostPath struct {
	RootPath string `yaml:"rootPath"`
	Path     string `yaml:"path"`
}

type Resource struct {
	Host        string
	Port        int32
	Username    string
	Password    string
	Schema      string
	Domain      string
	External    bool
	Url         string
	Persistence *Persistence `yaml:"persistence"`
}

func (p *Persistence) GetStorageType() string {
	if p.StorageClassName != "" {
		p.Type = "storageClass"
		return PersistenceStorageClassType
	}
	if p.Nfs.Server != "" {
		p.Type = "nfs"
		return PersistenceNfsType
	}
	if p.HostPath.RootPath != "" || p.HostPath.Path != "" {
		p.Type = "hostPath"
		return PersistenceHostPathType
	}
	return ""
}

func (p *Persistence) GetPersistentVolumeSource(subPath string) v1.PersistentVolumeSource {
	storageType := p.GetStorageType()
	if storageType == PersistenceNfsType {
		return p.prepareNfsPVS(subPath)
	}
	if storageType == PersistenceHostPathType {
		return p.prepareHostPathPVS(subPath)
	}
	return v1.PersistentVolumeSource{}
}

func (p *Persistence) prepareNfsPVS(subPath string) v1.PersistentVolumeSource {
	pvs := v1.PersistentVolumeSource{
		NFS: &v1.NFSVolumeSource{
			Server:   p.Server,
			Path:     fmt.Sprintf("%s/%s", p.Nfs.RootPath, subPath),
			ReadOnly: false,
		},
	}
	return pvs
}

func (p *Persistence) prepareHostPathPVS(subPath string) v1.PersistentVolumeSource {

	path := fmt.Sprintf("%s/%s", p.HostPath.RootPath, subPath)
	if p.HostPath.Path != "" {
		path = p.HostPath.Path
	}
	pvs := v1.PersistentVolumeSource{
		HostPath: &v1.HostPathVolumeSource{
			Path: path,
		},
	}
	return pvs
}
