#!/bin/sh

minikube image build -t konstellation/demo-nats-trigger-py . -f Dockerfile -p kai-local
