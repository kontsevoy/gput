#!/bin/bash

API_KEY=$(cat api-key)
USERNAME=$(cat username)

# Authenticate via Cloud Files API. This is what we're posting:

# HTTP POST:
#{
#    "auth": {
#        "RAX-KSKEY:apiKeyCredentials": {
#            "apiKey": "f0a0d5090040f77777777777751537c7",
#            "username": "xxxxxxxxxx"
#        }
#    }
#}
curl -s https://identity.api.rackspacecloud.com/v2.0/tokens  \
 -X POST \
 -d "{\"auth\":{\"RAX-KSKEY:apiKeyCredentials\": {\"username\":\"$USERNAME\",\"apiKey\":\"$API_KEY\"}}}" \
 -H "Content-type: application/json" | python -m json.tool | grep '"id":\s"\w\{32\}'

# the output can be seen in auth-output.txt, 
# but here're relevant parts:

# "access": {
#    "token": 
#        {
#            "RAX-AUTH:authenticatedBy": ["APIKEY"],
#            "expires": "2014-01-09T15:08:53.645-06:00",
#            "id": "123abc5368901230076b12357897898"
#        },
#     "serviceCatalog": [
#        {
#            "endpoints": [
#                {
#                    "internalURL": "https://snet-storage101.dfw1.clouddrive.com/v1/MossoCloudFS_7584136c-9514-49b5-b568-393bccbf68f8",
#                    "publicURL": "https://storage101.dfw1.clouddrive.com/v1/MossoCloudFS_7584136c-9514-49b5-b568-393bccbf68f8",
#                    "region": "DFW",
#                    "tenantId": "MossoCloudFS_7584136c-9514-49b5-b568-393bccbf68f8"
#                },
#            ],
#            "name": "cloudFiles",
#            "type": "object-store"
#       }]
# }

#curl -i -X PUT "https://storage101.dfw1.clouddrive.com/v1/MossoCloudFS_7584136c-9514-49b5-b568-393bccbf68f8/ev-public/hello.jpg" /
#-H "X-Auth-Token: yourAuthToken" \
#-H "Content-Type: image/jpeg" \
#-H "Content-Length: 0"
