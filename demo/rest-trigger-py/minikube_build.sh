#!/bin/sh

minikube image build -t konstellation/demo-rest-trigger-py . -f Dockerfile -p kai-local
