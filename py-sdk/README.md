# KAI Python SDK

KAI SDK's implementation in Python.


## How it works

The SDK can be used in a KAI workflows for working with Python code

Once the Python SDK is deployed, it connects to NATS and it subscribes permanently to an input subject. Each node knows to which subject it has to subscribe and also to which subject it has to send messages, since the K8s manager tells it with environment variables. It is important to note that the nodes use a queue subscription, which allows load balancing of messages when there are multiple replicas

When a new message is published in the input subject of a node, it passes it down to a handler function, along with a context object formed by variables and useful methods for processing data. This handler is the solution implemented by the client and given in the krt file generated. Once executed, the result will be taken and transformed into a NATS message that will then be published to the next node's subject (indicated by an environment variable). After that, the node ACKs the message manually

## Usage (TODO UPDATE)

## Requirements (TODO UPDATE)

It is necessary to set the following environment variables in order to use the runner:

## Example (TODO UPDATE)

## Development

Install the dependencies 

If you don't have poetry installed (you must have python 3.11 installed in your system):

`python3 -m pip install --user poetry`

## Unit tests

Just run `make pytest` from the root folder
