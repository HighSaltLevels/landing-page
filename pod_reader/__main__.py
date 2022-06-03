"""
    Kubernetes Pod Reader service. This service should
    run only with a JWT associated with a serviceaccount
    that can only watch, get, and list pods
"""
from http import HTTPStatus
import json
import os

from flask import Flask, request
import requests

APP = Flask(__name__)
RESP_HEADERS = {"Content-Type": "text/plain", "Access-Control-Allow-Origin": "*"}
SERVICE_ACCOUNT = "/var/run/secrets/kubernetes.io/serviceaccount"

KUBE_URL = os.getenv("KUBERNETES_SERVICE_HOST", "localhost")
KUBE_PORT = os.getenv("KUBERNETES_SERVICE_PORT", "6443")
KUBE_CA_CERT = f"{SERVICE_ACCOUNT}/ca.crt"

with open(f"{SERVICE_ACCOUNT}/token") as token:
    SERVICE_ACCOUNT_TOKEN = token.read()

REQ_HEADERS = {"Authorization": f"Bearer {SERVICE_ACCOUNT_TOKEN}"}

class ListPodsError(Exception):
    """Error raised if there is an issue listing pods"""


def list_pods(namespace="default"):
    """Do an HTTP Get on pods in the specified namespace"""
    url = f"https://{KUBE_URL}:{KUBE_PORT}/api/v1/namespaces/{namespace}/pods"

    # Make 5 attempts
    for _ in range(5):
        try:
            resp = requests.get(
                url, headers=REQ_HEADERS, verify=KUBE_CA_CERT, timeout=5
            )
            resp.raise_for_status()

            data = resp.json()
            if not data.get("items"):
                raise ListPodsError(f"No pods found in namespace {namespace}.")

            return data

        except requests.RequestException as error:
            print(f"Retrying due to {error}.", flush=True)

    raise ListPodsError("Unable to fetch data.")


def parse_list_response(list_resp):
    """
    Parse out and format the json response from the kube api
    Example response:

    Pod Name                       | Status    | Pod IP         | Node Name
    -------------------------------+-----------+----------------+------------
    landing-page-76b8b9677f-nmddz  | Running   | 10.144.420.69  | salt-work1
    """

    response_message = (
        "Pod Name                      | Status    | Pod IP          | Node Name\n"
    )
    response_message += (
        (30 * "-") + "+" + (11 * "-") + "+" + (17 * "-") + "+" + (12 * "-")
    )
    for pod in list_resp.get("items"):
        pod_name = pod.get("metadata", {}).get("name", "Not Found")
        status = pod.get("status", {}).get("phase", "Not Found")
        pod_ip = pod.get("status", {}).get("podIP", "Not Found")
        node_name = pod.get("spec", {}).get("nodeName", "Not Found")

        response_message += f"\n{pod_name:30}| {status:10}| {pod_ip:16}| {node_name:11}"

    return response_message


@APP.after_request
def get_pod_status(response):
    """after_request handler for doing dynamic lookups on any namespace"""
    # Automatically set response to 200 with unknown error message
    response.status_code = HTTPStatus.OK
    response.headers = RESP_HEADERS
    response.data = json.dumps({"error": "Unknown error occured."})

    # Assume requested namespace is the full uri except for the leading slash
    namespace = request.path[1:]

    try:
        list_resp = list_pods(namespace)
        response.data = parse_list_response(list_resp)

    except ListPodsError as error:
        response.data = json.dumps({"error": str(error)})

    return response


APP.run("0.0.0.0", 42069)
