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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

// KserviceHandler handles the incoming requests to add a sidecar container to the Knative service.
func KserviceHandler(w http.ResponseWriter, r *http.Request) {

	config := config.NewIntance()

	log.Println("Received a request to handle Kservice")

	// Parse the admission review request from the incoming HTTP body.
	admissionReview := admissionv1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		log.Printf("Failed to decode admission review request: %v", err)
		http.Error(w, "Failed to decode admission review request", http.StatusBadRequest)
		return
	}

	knativeService := servingv1.Service{}
	// Unmarshal the object from the admission request into the knativeService object
	err = json.Unmarshal(admissionReview.Request.Object.Raw, &knativeService)
	if err != nil {
		log.Printf("Failed to unmarshal service object: %v", err)
		http.Error(w, "Failed to unmarshal service object", http.StatusBadRequest)
		return
	}

	mknativeService := knativeService.DeepCopy()

	// Check if the label matches the desired value
	if knativeService.Labels[config.Get("LABEL_KEY")] == config.Get("LABEL_VALUE") {
		log.Println("Label match found, preparing to add sidecar container")

		// Parse FILTERS from the environment
		jsonStr := config.Get("FILTERS")
		var objs []filterv1.Object

		err := json.Unmarshal([]byte(jsonStr), &objs)
		if err != nil {
			log.Printf("Error while unmarshaling FILTERS from environment: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Initialize a base port with the value 8080.
		basePort64, err := strconv.ParseInt(knativeService.Annotations["k8sidecar.port"], 10, 32)

		var basePort int32
		if err != nil {
			basePort = int32(8080)
		} else {
			basePort = int32(basePort64)
		}

		// Get a deep copy of the current containers from the modified Knative service.
		baseContainers := mknativeService.Spec.Template.Spec.Containers

		// Loop through each object (which presumably represents desired containers).
		for _, obj := range objs {
			// Create a new container using the object's name and image, and append a new environment variable for priority.
			c := corev1.Container{
				Name:  obj.Name,
				Image: obj.Image,
				// Append the object's environment variables with an additional one for priority.
				Env: append(obj.Env, corev1.EnvVar{
					Name:  "PPRIORITY",
					Value: strconv.Itoa(int(obj.Priority)),
				}),
			}
			// Append this new container to our list of base containers.
			baseContainers = append(baseContainers, c)
		}

		// Sort the base containers by their priority in descending order.
		sort.Slice(baseContainers, func(i, j int) bool {
			return getPPriority(baseContainers[i].Env) > getPPriority(baseContainers[j].Env)
		})

		// Iterate through each container in the base containers.
		for i := range baseContainers {
			// Calculate a dynamic port based on the container's position.
			pport := basePort + int32(i+1)
			// Set an environment variable for the calculated port.
			setEnvVar(&baseContainers[i].Env, "PPORT", strconv.Itoa(int(pport)))
			// Ensure a shared volume exists for this container.
			addVolumeIfNotExist(&baseContainers[i].VolumeMounts)
			// If this is the first container in the list, set its port details.
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
		// Update the modified Knative service's containers with our base containers.
		mknativeService.Spec.Template.Spec.Containers = baseContainers

		// Set a shared volume to the modified Knative service's volumes.
		mknativeService.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: "shared-data",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{
						Medium: corev1.StorageMedium(config.Get("MEMORY")),
					},
				},
			},
		}

	}

	// Compare original and modified services and generate a patch
	patch, err := jsondiff.Compare(knativeService, mknativeService)
	if err != nil {
		log.Printf("Failed to generate patch: %v", err)
		http.Error(w, "Failed to patch Pod object", http.StatusBadRequest)
		return
	}

	patchb, err := json.Marshal(patch)
	if err != nil {
		log.Printf("Failed to marshal patch into bytes: %v", err)
		http.Error(w, "Failed to patch Pod object", http.StatusBadRequest)
		return
	}

	// Create an admission response
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

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(admissionReview)
	if err != nil {
		log.Printf("Failed to encode admission review response: %v", err)
		http.Error(w, "Failed to encode admission review response", http.StatusInternalServerError)
		return
	}

	log.Println("KserviceHandler completed successfully")
}

// getPPriority retrieves the priority from the environment variables
func getPPriority(envs []corev1.EnvVar) int {
	for _, env := range envs {
		if env.Name == "PPRIORITY" {
			value, _ := strconv.Atoi(env.Value)
			return value
		}
	}
	return -1
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
func addVolumeIfNotExist(volumes *[]corev1.VolumeMount) {
	for _, volume := range *volumes {
		if volume.Name == "shared-data" {
			return
		}
	}

	*volumes = append(*volumes, corev1.VolumeMount{Name: "shared-data", MountPath: "/shared"})
}
