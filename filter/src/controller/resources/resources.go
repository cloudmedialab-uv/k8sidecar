package resources

import (
	"context"
	"encoding/json"
	"log"

	v1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"filter/src/controller/gen"
	filterv1 "filter/src/pkg/apis/filtercontroller/v1"
)

// Api wraps a Kubernetes client to manage custom filter resources.
type Api struct {
	client *kubernetes.Clientset
}

// NewApi initializes and returns a new API instance.
func NewApi(client *kubernetes.Clientset) *Api {
	return &Api{
		client: client,
	}
}

// CreateResources creates all required resources for the given filter.
func (api *Api) CreateResources(filters *filterv1.Filter) error {
	// Extract the filter name and marshal its sidecars to JSON.
	name := filters.ObjectMeta.Name
	jFilters, err := json.Marshal(filters.Spec.Sidecars)
	if err != nil {
		log.Printf("Failed to marshal sidecars for filter %s: %v", name, err)
		return err
	}

	// Generate a certificate for the filter.
	cert, key, err := gen.GenCert(name)
	if err != nil {
		return err
	}

	// Create the associated resources, logging and collecting errors as needed.
	if err = api.createDeployment(name, jFilters, cert, key); err != nil {
		log.Printf("Failed to create deployment for filter %s: %v", name, err)
		return err
	}

	if err = api.createService(name); err != nil {
		log.Printf("Failed to create service for filter %s: %v", name, err)
		return err
	}

	if err = api.createKnativeMutatingWebhook(name, cert); err != nil {
		log.Printf("Failed to create Knative mutating webhook for filter %s: %v", name, err)
		return err
	}

	if err = api.createDeploymentMutatingWebhook(name, cert); err != nil {
		log.Printf("Failed to create deployment mutating webhook for filter %s: %v", name, err)
		return err
	}

	return nil
}

// createDeployment creates a Kubernetes deployment for the filter.

func (api *Api) createDeployment(name string, jFilters []byte, cert string, key string) error {
	// Deployment definition
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "filter-deployment-" + name,
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func() *int32 { i := int32(1); return &i }(),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "filter-deployment-" + name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "filter-deployment-" + name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "admssion-server",
							Image: "routerdi1315.uv.es:33443/sidecar/filter/admission:1.7.test",
							Env: []corev1.EnvVar{
								{
									Name:  "FILTERS",
									Value: string(jFilters),
								},
								{
									Name:  "LABEL",
									Value: name,
								},
								{
									Name:  "TLS_CRT",
									Value: cert,
								},
								{
									Name:  "TLS_KEY",
									Value: key,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8443,
								},
							},
						},
					},
				},
			},
		},
	}
	// Create the deployment in the cluster.
	_, err := api.client.AppsV1().Deployments(deployment.Namespace).Create(context.Background(), deployment, metav1.CreateOptions{})

	return err
}

// createService creates a Kubernetes service for the filter.
func (api *Api) createService(name string) error {
	// Service definition for the filter.
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-sidecar-" + name,
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "filter-deployment-" + name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "https",
					Port:       443,
					TargetPort: intstr.FromInt(8443),
				},
			},
		},
	}

	// Attempt to create the service in the cluster.
	_, err := api.client.CoreV1().Services("default").Create(context.Background(), service, metav1.CreateOptions{})
	return err
}

// createKnativeMutatingWebhook creates a Knative mutating webhook for the filter.
func (api *Api) createKnativeMutatingWebhook(name string, cert string) error {

	// Define the mutating webhook configuration for Knative.
	mutatingWebhookConfiguration := &v1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "filter-mutation-knative." + name,
		},
		Webhooks: []v1.MutatingWebhook{
			{
				Name: "filter.knative." + name + ".com",
				Rules: []v1.RuleWithOperations{
					{
						Operations: []v1.OperationType{"CREATE"},
						Rule: v1.Rule{
							APIGroups:   []string{"serving.knative.dev"},
							APIVersions: []string{"v1"},
							Resources:   []string{"services"},
							Scope:       &[]v1.ScopeType{v1.NamespacedScope}[0],
						},
					},
				},
				ClientConfig: v1.WebhookClientConfig{
					Service: &v1.ServiceReference{
						Namespace: "default",
						Name:      "service-sidecar-" + name,
						Path:      &[]string{"/kservice"}[0],
					},
					CABundle: []byte(cert),
				},
				SideEffects:             &[]v1.SideEffectClass{v1.SideEffectClassNone}[0],
				AdmissionReviewVersions: []string{"v1", "v1beta1"},
			},
		},
	}
	// Attempt to create the mutating webhook in the cluster.
	_, err := api.client.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(context.Background(), mutatingWebhookConfiguration, metav1.CreateOptions{})
	return err
}

// createDeploymentMutatingWebhook creates a deployment mutating webhook for the filter.
func (api *Api) createDeploymentMutatingWebhook(name string, cert string) error {

	// Define the mutating webhook configuration for deployments.
	mutatingWebhookConfiguration := &v1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "filter-mutation-deployment-" + name,
		},
		Webhooks: []v1.MutatingWebhook{
			{
				Name: "filter.deployment." + name + ".com",
				Rules: []v1.RuleWithOperations{
					{
						Operations: []v1.OperationType{"CREATE"},
						Rule: v1.Rule{
							APIGroups:   []string{"apps"},
							APIVersions: []string{"v1"},
							Resources:   []string{"deployments"},
							Scope:       &[]v1.ScopeType{v1.NamespacedScope}[0],
						},
					},
				},
				ClientConfig: v1.WebhookClientConfig{
					Service: &v1.ServiceReference{
						Namespace: "default",
						Name:      "service-sidecar-" + name,
						Path:      &[]string{"/deployment"}[0],
					},
					CABundle: []byte(cert),
				},
				SideEffects:             &[]v1.SideEffectClass{v1.SideEffectClassNone}[0],
				AdmissionReviewVersions: []string{"v1", "v1beta1"},
			},
		},
	}
	// Attempt to create the mutating webhook in the cluster.
	_, err := api.client.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(context.Background(), mutatingWebhookConfiguration, metav1.CreateOptions{})
	return err
}

// DeleteResources deletes all resources associated with the given filter name.
func (api *Api) DeleteResources(name string) error {
	// Delete associated resources, logging and collecting errors as needed.
	if err := api.deleteDeployment(name); err != nil {
		log.Printf("Failed to delete deployment for filter %s: %v", name, err)
		return err
	}

	if err := api.deleteService(name); err != nil {
		log.Printf("Failed to delete service for filter %s: %v", name, err)
		return err
	}

	if err := api.deleteKnativeMutatingWebhookConfiguration(name); err != nil {
		log.Printf("Failed to delete Knative mutating webhook configuration for filter %s: %v", name, err)
		return err
	}

	if err := api.deleteDeploymentMutatingWebhookConfiguration(name); err != nil {
		log.Printf("Failed to delete deployment mutating webhook configuration for filter %s: %v", name, err)
		return err
	}

	return nil
}

// deleteDeployment deletes the Kubernetes deployment associated with the filter.
func (api *Api) deleteDeployment(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return api.client.AppsV1().Deployments("default").Delete(context.Background(), "filter-deployment-"+name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

// deleteService deletes the Kubernetes service associated with the filter.
func (api *Api) deleteService(name string) error {
	return api.client.CoreV1().Services("default").Delete(context.Background(), "service-sidecar-"+name, metav1.DeleteOptions{})
}

// deleteKnativeMutatingWebhookConfiguration deletes the Knative mutating webhook configuration associated with the filter.
func (api *Api) deleteKnativeMutatingWebhookConfiguration(name string) error {
	return api.client.AdmissionregistrationV1().MutatingWebhookConfigurations().Delete(context.Background(), "filter-mutation-"+name, metav1.DeleteOptions{})
}

// deleteDeploymentMutatingWebhookConfiguration deletes the deployment mutating webhook configuration associated with the filter.
func (api *Api) deleteDeploymentMutatingWebhookConfiguration(name string) error {
	return api.client.AdmissionregistrationV1().MutatingWebhookConfigurations().Delete(context.Background(), "filter-mutation-"+name, metav1.DeleteOptions{})
}
