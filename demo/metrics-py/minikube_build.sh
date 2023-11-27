#!/bin/sh

minikube image build -t konstellation/metrics-py . -f Dockerfile -p kai-local
