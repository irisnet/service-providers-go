# service-providers-go
Service Providers Implemented in Golang

## 1 Config

- Configuration parameter:
| name | description |
| :-: | :-: |
| chain_id | Chain id |
| node_rpc_addr | Node URL |
| node_grpc_addr | Node GRPC address |
| key_path | Key path |
| key_name | Key name |
| fee | Transaction fee |
| key_algorithm | Key algorithm |

- Example
```yaml
chain_id: iris
node_rpc_addr: http://localhost:26657
node_grpc_addr: http://localhost:9090
key_path: .keys
key_name: node0
fee: 4uiris
key_algorithm: sm2
```

## 2 Key management

  - Commond to key management
    | commond | description |
    | :-: | :-: |
    | add | New-build key |
    | show | Show information of key |
    | import | Import key |
      
    - You need to put the exported information into a file .keys, and specify the path of the file in config.yaml.

      ### 1.1 Export node0

        ```shell
        iris keys export node0 --home /home/sunny/iris/node0/ iriscli
        ```

      ### 1.2 Import node0

        ```shell
        {{service_name}}-sp keys import node0
        ```

## 3  The files that need to be modified are on the floder random-seed/random-seed and token-price/token-price.

## 4 Run docker

    - build
    
        ```shell
        docker build -t {{service_name}} .
        ```
    
    -run

        ```shell
        docker run {{service_name}} start
        ```
