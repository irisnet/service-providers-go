# random-seed

## 1 Config

- Configuration parameter
  
    | name | description |
    | :---: | :---: |
    | chain_id | Chain id |
    | node_rpc_addr | Node URL |
    | node_grpc_addr | Node GRPC address |
    | key_path | Key path |
    | key_name | Key name |
    | fee | Transaction fee |
    | key_algorithm | Key algorithm |

- Example
    ```yaml
    chain_id: irishub
    node_rpc_addr: http://localhost:26657
    node_grpc_addr: http://localhost:9090
    key_path: .keys
    key_name: node0
    fee: 0.4iris
    key_algorithm: secp256k1
    ```

## 2 Key management

  - Commond to key management
    
    | commond | description |
    | :---: | :---: |
    | add | Add a new key |
    | show | Show information of key |
    | import | Import key |
      
- You need to put the exported information into a file .keys, and specify the path of the file in config.yaml.

  ### 2.1 Export node0

    ```shell
    iris keys export node0 --home /home/testnet/node0/iriscli
    ```

  ### 2.2 Import node0

    ```shell
    random-sp keys import node0
    ```

## 3 Modify config.yaml

## 4 Run docker

- build
  
    ```shell
    docker build -t random .
    ```
    
  -run

    ```shell
    docker run -it random
    echo ${your_password} | random-sp start
    ```