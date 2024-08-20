package handlers

import (
	"errors"
	"strconv"

	corev1 "k8s.io/api/core/v1"
)

func getEnvVar(envs []corev1.EnvVar, key string) (string, error) {
	for _, env := range envs {
		if env.Name == key {
			return env.Value, nil
		}
	}
	return "", errors.New("env not found")
}

// getEnvPPriority retrieves the priority from the environment variables
func getEnvPPriority(envs []corev1.EnvVar) int {
	env, err := getEnvVar(envs, "PPRIORITY")
	if err != nil {
		return -1
	}

	value, _ := strconv.Atoi(env)
	return value
}

func getEnvPort(envs []corev1.EnvVar, annotation string) int32 {

	env, err := getEnvVar(envs, annotation)
	if err != nil {
		return 8080
	}

	basePort64, err := strconv.ParseInt(env, 10, 32)
	if err != nil {
		return 8080
	} else {
		return int32(basePort64)
	}

}

func getEnvSharedVolume(envs []corev1.EnvVar, annotation map[string]string) bool {
	env, err := getEnvVar(envs, annotation["k8sidecar.shared-volume"])
	if err != nil {
		return true
	}
	res, err := strconv.ParseBool(env)

	if err != nil {
		return true
	}

	return res
}

// setEnvVar sets or appends a new environment variable
func setEnvVar(envs *[]corev1.EnvVar, key, value string) {
	for i, env := range *envs {
		if env.Name == key {
			(*envs)[i].Value = value
			return
		}
	}

	*envs = append(*envs, corev1.EnvVar{Name: key, Value: value})
}

// addVolumeIfNotExist adds a shared volume if it doesn't exist
func addVolumeMountIfNotExist(volumes *[]corev1.VolumeMount) {
	for _, volume := range *volumes {
		if volume.Name == "shared-volume" {
			return
		}
	}

	*volumes = append(*volumes, corev1.VolumeMount{Name: "shared-volume", MountPath: "/shared"})
}

func existVolume(volumes []corev1.Volume) bool {
	for _, volume := range volumes {
		if volume.Name == "shared-volume" {
			return true
		}
	}
	return false
}
