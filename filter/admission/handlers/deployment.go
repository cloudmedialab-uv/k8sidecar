package handlers

import (
	"encoding/json"
	"filter/admission/config"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/wI2L/jsondiff"

	filterv1 "filter/pkg/apis/filtercontroller/v1"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeploymentHandler(w http.ResponseWriter, r *http.Request) {
	config := config.NewIntance()

	log.Println("Received a request to handle deployments")

	admissionReview := admissionv1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		http.Error(w, "Failed to decode admission review request", http.StatusBadRequest)
		return
	}

	deployment := appsv1.Deployment{}
	err = json.Unmarshal(admissionReview.Request.Object.Raw, &deployment)
	if err != nil {
		http.Error(w, "Failed to unmarshal deployment object", http.StatusBadRequest)
		return
	}

	mDeployment := deployment.DeepCopy()

	if deployment.Labels[config.Get("LABEL_KEY")] == config.Get("LABEL_VALUE") {
		log.Println("adding sidecar container")

		jsonStr := config.Get("FILTERS")
		var objs []filterv1.Sidecar

		err := json.Unmarshal([]byte(jsonStr), &objs)
		if err != nil {
			log.Print("Error: ", err)
			return
		}

		baseContainers := mDeployment.Spec.Template.Spec.Containers
		baseContainer := &baseContainers[len(baseContainers)-1].Env
		basePortTag, ok := mDeployment.Annotations["k8sidecar.port"]
		if !ok {
			basePortTag = "PORT"
		}

		basePort := getEnvPort(*baseContainer, basePortTag)
		addVolume := getEnvSharedVolume(*baseContainer, mDeployment.Annotations)

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
			return getEnvPPriority(baseContainers[i].Env) > getEnvPPriority(baseContainers[j].Env)
		})

		for i := range baseContainers {
			pport := basePort + int32(i)

			if i == len(baseContainers)-1 {
				print("ADDING TO baseContainer port: " + strconv.Itoa(int(pport)))
				setEnvVar(baseContainer, basePortTag, strconv.Itoa(int(pport)))
			} else {
				setEnvVar(&baseContainers[i].Env, "PPORT", strconv.Itoa(int(pport)))
			}

			if addVolume {
				addVolumeIfNotExist(&baseContainers[i].VolumeMounts)
			}

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

		mDeployment.Spec.Template.Spec.Containers = baseContainers

		if addVolume {
			mDeployment.Spec.Template.Spec.Volumes = []corev1.Volume{
				{
					Name: "shared-volume",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							Medium: corev1.StorageMedium(config.Get("MEMORY")),
						},
					},
				},
			}
		}

		log.Println(mDeployment.Spec.Template.Spec.Containers)
	}

	patch, err := jsondiff.Compare(deployment, mDeployment)
	if err != nil {
		http.Error(w, "Failed to patch Deployment object", http.StatusBadRequest)
		return
	}

	patchb, err := json.Marshal(patch)
	if err != nil {
		http.Error(w, "Failed to patch Deployment object", http.StatusBadRequest)
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
