#!/bin/sh

minikube image build -t konstellation/github-trigger-mock . -f Dockerfile -p kai-local
