#!/bin/sh

minikube image build -t konstellation/training-py . -f Dockerfile -p kai-local
