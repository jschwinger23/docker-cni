package cni

import (
	"strings"

	"github.com/projecteru2/docker-cni/config"
	"github.com/projecteru2/docker-cni/oci"
)

func (h *CNIHandler) HandleCreate(conf config.Config, containerMeta *oci.ContainerMeta) (err error) {
	if err = h.AddCNIStartHook(conf, containerMeta); err != nil {
		return
	}
	if err = h.AddCNIStopHook(conf, containerMeta); err != nil {
		return
	}
	return containerMeta.Save()
}

func (h *CNIHandler) AddCNIStartHook(conf config.Config, containerMeta *oci.ContainerMeta) (err error) {
	env := []string{}
	cniArgs := []string{}
	if containerMeta.RequiresSpecificIPPool() {
		cniArgs = append(cniArgs, "IPPOOL="+containerMeta.SpecificIPPool())
	}
	if containerMeta.RequiresSpecificIP() {
		cniArgs = append(cniArgs, "IP="+containerMeta.SpecificIP())
	}
	env = append(env, "CNI_ARGS="+strings.Join(cniArgs, ";"))
	containerMeta.AppendHook("prestart",
		conf.BinPathname,
		[]string{conf.BinPathname, "cni", "--config", conf.Filename, "--command", "add"}, // args
		env, // envs
	)
	return
}

func (h *CNIHandler) AddCNIStopHook(conf config.Config, containerMeta *oci.ContainerMeta) (err error) {
	containerMeta.AppendHook("poststop",
		conf.BinPathname,
		[]string{conf.BinPathname, "cni", "--config", conf.Filename, "--command", "del"}, // args
		nil, // env
	)
	return
}
