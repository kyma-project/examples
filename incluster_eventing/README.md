# Asynchronous communication between Functions

## Overview

This example provides a very simple scenario of asynchronous communication between two Functions, where: 
- The first Function accepts the incoming traffic via HTTP, sanitizes the payload, and publishes the content as an in-cluster event via [Kyma Eventing](https://kyma-project.io/docs/kyma/latest/01-overview/main-areas/eventing/).
- The second Function is a message receiver. It subscribes to the given event type and stores the payload.

## Prerequisites

* [Kyma CLI](https://github.com/kyma-project/cli)
* Kyma installed locally or on a cluster

## Installation

You can find all installation steps in the [Set asynchronous communication between Functions](https://github.com/kyma-project/kyma/blob/b783d9e6dffc47c0e3c31923aff62371b0a46779/docs/03-tutorials/00-serverless/svls-11-set-asynchronous-connection-of-functions.md) tutorial.