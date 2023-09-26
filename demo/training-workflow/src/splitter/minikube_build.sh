#!/bin/sh

minikube image build -t konstellation/splitter . -f Dockerfile -p kai-local
