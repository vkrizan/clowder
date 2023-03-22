package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/RedHatInsights/rhc-osdk-utils/utils"
)

type mutantPod struct {
	Client   client.Client
	Recorder record.EventRecorder
	decoder  *admission.Decoder
}

func (p *mutantPod) Handle(_ context.Context, req admission.Request) admission.Response {
	pod := &core.Pod{}

	err := p.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if v, ok := pod.GetAnnotations()["clowder/authsidecar-enabled"]; ok && v == "true" {
		ridx := -1
		for idx, container := range pod.Spec.Containers {
			if container.Name == "crcauth" {
				ridx = idx
				break
			}
		}

		port, ok := pod.GetAnnotations()["clowder/authsidecar-port"]
		if !ok {
			return admission.Errored(http.StatusBadRequest, fmt.Errorf("pod does not specify authsidecar port"))
		}

		config, ok := pod.GetAnnotations()["clowder/authsidecar-config"]
		if !ok {
			return admission.Errored(http.StatusBadRequest, fmt.Errorf("pod does not specify authsidecar config"))
		}

		image, ok := pod.GetAnnotations()["clowder/authsidecar-image"]
		if !ok {
			return admission.Errored(http.StatusBadRequest, fmt.Errorf("pod does not specify authsidecar image"))
		}

		probeHandler := core.ProbeHandler{
			TCPSocket: &core.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 8080,
				},
			},
		}

		livenessProbe := core.Probe{
			ProbeHandler:        probeHandler,
			InitialDelaySeconds: 10,
			TimeoutSeconds:      2,
			PeriodSeconds:       10,
			SuccessThreshold:    1,
			FailureThreshold:    3,
		}
		readinessProbe := core.Probe{
			ProbeHandler:        probeHandler,
			InitialDelaySeconds: 15,
			TimeoutSeconds:      2,
			PeriodSeconds:       10,
			SuccessThreshold:    1,
			FailureThreshold:    3,
		}

		container := core.Container{
			Name:           "crcauth",
			Image:          image,
			LivenessProbe:  &livenessProbe,
			ReadinessProbe: &readinessProbe,
			Env: []core.EnvVar{
				{
					Name:  "CADDY_PORT",
					Value: port,
				},
				{
					Name: "CADDY_BOP_URL",
					ValueFrom: &core.EnvVarSource{
						SecretKeyRef: &core.SecretKeySelector{
							LocalObjectReference: core.LocalObjectReference{
								Name: config,
							},
							Optional: utils.BoolPtr(false),
							Key:      "bopurl",
						},
					},
				},
				{
					Name: "CADDY_KEYCLOAK_URL",
					ValueFrom: &core.EnvVarSource{
						SecretKeyRef: &core.SecretKeySelector{
							LocalObjectReference: core.LocalObjectReference{
								Name: config,
							},
							Optional: utils.BoolPtr(false),
							Key:      "keycloakurl",
						},
					},
				},
				{
					Name: "CADDY_WHITELIST",
					ValueFrom: &core.EnvVarSource{
						SecretKeyRef: &core.SecretKeySelector{
							LocalObjectReference: core.LocalObjectReference{
								Name: config,
							},
							Optional: utils.BoolPtr(false),
							Key:      "whitelist",
						},
					},
				},
			},
			TerminationMessagePath:   "/dev/termination-log",
			TerminationMessagePolicy: core.TerminationMessageReadFile,
			ImagePullPolicy:          core.PullIfNotPresent,
			Ports: []core.ContainerPort{{
				Name:          "auth",
				ContainerPort: 8080,
				Protocol:      "TCP",
			}},
			Resources: core.ResourceRequirements{
				Limits: core.ResourceList{
					"memory": resource.MustParse("200Mi"),
					"cpu":    resource.MustParse("100m"),
				},
				Requests: core.ResourceList{
					"memory": resource.MustParse("100Mi"),
					"cpu":    resource.MustParse("50m"),
				},
			},
		}

		if ridx == -1 {
			pod.Spec.Containers = append(pod.Spec.Containers, container)
		} else {
			pod.Spec.Containers[ridx] = container
		}
	}

	marshaledObj, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledObj)
}

func (p *mutantPod) InjectDecoder(d *admission.Decoder) error {
	p.decoder = d
	return nil
}
