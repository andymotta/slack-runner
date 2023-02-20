#!/bin/bash

DOMAIN="${1}"

QUERY="$(aws es describe-elasticsearch-domain --domain-name "$DOMAIN" --query 'DomainStatus.Processing')"; fi
if [ "$QUERY" == "false" ]
  then
    tput setaf 2; echo "The Elasticsearch cluster is ready"
  else
    tput setaf 1; echo "The Elasticsearch cluster is NOT ready"
fi