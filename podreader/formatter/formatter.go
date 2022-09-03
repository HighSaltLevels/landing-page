package formatter

import (
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"log"
)

/* Parse out and format the json response from the kube api */
func FormatStatus(respBody string) (string, error) {
	var podList corev1.PodList
	err := json.Unmarshal([]byte(respBody), &podList)
	if err != nil {
		log.Printf("Unable to unmarshal response body into a V1Pod. Body was: %s", respBody)
		log.Printf("Error was: %v", err)
		return "Unexpected Error", err
	}

	// Default each length to the length of the title of the column
	longestName, longestStatus, longestPodIp, longestNodeName := 9, 8, 8, 11
	// Assume that order is preserved and is directly correlated to the order of podList
	for _, pod := range podList.Items {
		if len(pod.ObjectMeta.Name) > longestName {
			longestName = len(pod.ObjectMeta.Name)
		}

		if len(pod.Status.Phase) > longestStatus {
			longestStatus = len(pod.Status.Phase)
		}

		if len(pod.Status.PodIP) > longestPodIp {
			longestPodIp = len(pod.Status.PodIP)
		}

		if len(pod.Spec.NodeName) > longestNodeName {
			longestNodeName = len(pod.Spec.NodeName)
		}
	}

	formattedResp := "Pod Name" + createChars(" ", longestName-8) + " | Status" + createChars(" ", longestStatus-6) + " | Pod IP" + createChars(" ", longestPodIp-6) + " | Node Name\n" + createChars("-", longestName) + " + " + createChars("-", longestStatus) + " + " + createChars("-", longestPodIp) + " + " + createChars("-", longestNodeName) + "\n"
	for _, pod := range podList.Items {
		formattedResp = formattedResp + pod.ObjectMeta.Name + createChars(" ", longestName-len(pod.ObjectMeta.Name)) + " | " + string(pod.Status.Phase) + createChars(" ", longestStatus-len(pod.Status.Phase)) + " | " + pod.Status.PodIP + createChars(" ", longestPodIp-len(pod.Status.PodIP)) + " | " + pod.Spec.NodeName + createChars(" ", longestNodeName-len(pod.Spec.NodeName)) + "\n"
	}
	return formattedResp, nil
}

func createChars(char string, numChars int) string {
	chars := ""
	for i := 0; i < numChars; i++ {
		chars += char
	}
	return chars
}
