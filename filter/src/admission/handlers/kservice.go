package handlers

import (
	"encoding/json"
	"filter/src/admission/config"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/wI2L/jsondiff"

	filterv1 "filter/src/pkg/apis/filtercontroller/v1"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func KserviceHandler(w http.ResponseWriter, r *http.Request) {

	admissionReview := admissionv1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		http.Error(w, "Failed to decode admission review request", http.StatusBadRequest)
		return
	}

	knativeService := servingv1.Service{}
	err = json.Unmarshal(admissionReview.Request.Object.Raw, &knativeService)
	if err != nil {
		http.Error(w, "Failed to unmarshal service object", http.StatusBadRequest)
		return
	}

	mknativeService := knativeService.DeepCopy()

	if knativeService.Labels[config.LabelKey] == config.LabelValue {
		log.Println("adding sidecar container")

		jsonStr := os.Getenv("FILTERS")
		var objs []filterv1.Object

		err := json.Unmarshal([]byte(jsonStr), &objs)
		if err != nil {
			log.Print("Error: ", err)
			return
		}

		basePort := int32(8080)
		baseContainers := mknativeService.Spec.Template.Spec.Containers

		for _, obj := range objs {
			c := corev1.Container{
				Name:  obj.Name,
				Image: obj.Image,
				Env: append(obj.Env, corev1.EnvVar{
					Name:  "PPRIORITY",
					Value: strconv.Itoa(int(obj.Priority)),
				}),
			}
			baseContainers = append(baseContainers, c)
		}

		sort.Slice(baseContainers, func(i, j int) bool {
			return getPPriority(baseContainers[i].Env) > getPPriority(baseContainers[j].Env)
		})

		for i := range baseContainers {
			pport := basePort + int32(i) - int32(len(baseContainers))
			setEnvVar(&baseContainers[i].Env, "PPORT", strconv.Itoa(int(pport)))
			addVolumeIfNotExist(&baseContainers[i].VolumeMounts)

			if i == 0 {
				baseContainers[0].Ports = []corev1.ContainerPort{
					{
						Name:          "http1",
						Protocol:      corev1.ProtocolTCP,
						ContainerPort: pport,
					},
				}
			}
		}

		mknativeService.Spec.Template.Spec.Containers = baseContainers

		mknativeService.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: "shared-data",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		}

		log.Println(mknativeService.Spec.Template.Spec.Containers)
	}

	patch, err := jsondiff.Compare(knativeService, mknativeService)
	if err != nil {
		http.Error(w, "Failed to patch Pod object", http.StatusBadRequest)
		return
	}

	patchb, err := json.Marshal(patch)
	if err != nil {
		http.Error(w, "Failed to patch Pod object", http.StatusBadRequest)
		return
	}

	admissionReview.Response = &admissionv1.AdmissionResponse{
		UID:     admissionReview.Request.UID,
		Allowed: true,
		Result: &metav1.Status{
			Message: "Sidecar injection successful",
		},
		Patch: patchb,
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(admissionReview)
	if err != nil {
		http.Error(w, "Failed to encode admission review response", http.StatusInternalServerError)
		return
	}
}

func getPPriority(envs []corev1.EnvVar) int {
	for _, env := range envs {
		if env.Name == "PPRIORITY" {
			value, _ := strconv.Atoi(env.Value)
			return value
		}
	}
	return 0
}

func setEnvVar(envs *[]corev1.EnvVar, key, value string) {
	for i, env := range *envs {
		if env.Name == key {
			(*envs)[i].Value = value
			return
		}
	}

	*envs = append(*envs, corev1.EnvVar{Name: key, Value: value})
}

func addVolumeIfNotExist(volumes *[]corev1.VolumeMount) {
	for _, volume := range *volumes {
		if volume.Name == "shared-data" {
			return
		}
	}

	*volumes = append(*volumes, corev1.VolumeMount{Name: "shared-data", MountPath: "/shared"})
}
