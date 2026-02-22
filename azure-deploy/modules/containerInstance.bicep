param location string
param containerGroupName string
param containerUMIResourceId string
param imageName string
param bloodhoundInstanceDomain string
param azureTenantId string

@secure()
param bloodhoundTokenId string

@secure()
param bloodhoundToken string

var config = {
  app: ''
  auth: ''
  batchsize: 100
  config: '/home/nonroot/.config/azurehound/config.json'
  instance: 'https://${bloodhoundInstanceDomain}/'
  json: false
  'managed-identity': true
  maxconnsperhost: 20
  maxidleconnsperhost: 20
  region: 'cloud'
  streamcount: 25
  tenant: azureTenantId
  token: bloodhoundToken
  tokenid: bloodhoundTokenId
  verbosity: 0
}

resource containerGroup 'Microsoft.ContainerInstance/containerGroups@2023-05-01' = {
  name: containerGroupName
  location: location
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${containerUMIResourceId}': {}
    }
  }
  properties: {
    containers: [
      {
        name: 'azurehound'
        properties: {
          image: imageName
          command: [
            '/azurehound'
            'start'
          ]
          volumeMounts: [
            {
              name: 'config-volume'
              mountPath: '/home/nonroot/.config/azurehound'
            }
          ]
          resources: {
            requests: {
              cpu: '1'
              memoryInGB: '1'
            }
          }
        }
      }
    ]
    volumes: [
      {
        name: 'config-volume'
        secret: {
          'config.json': base64(string(config))
        }
      }
    ]
    osType: 'Linux'
    restartPolicy: 'Never'
  }
}
