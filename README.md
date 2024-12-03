> ### Updated!!!
>
> This extension was based on archived and deprecated https://github.com/grafana/xk6-amqp.  
> But all methods/objects was refactored.  
> I've used https://github.com/mostafa/xk6-kafka as an example of approved extension.
>
> USE AT YOUR OWN RISK!

<br />

# xk6-amqp-custom

A k6 extension for publishing and consuming messages from queues and exchanges.
This project utilizes [AMQP 0.9.1](https://www.rabbitmq.com/tutorials/amqp-concepts.html), the most common AMQP protocol in use today.

> :warning: This project is not compatible with [AMQP 1.0](http://docs.oasis-open.org/amqp/core/v1.0/os/amqp-core-overview-v1.0-os.html).
> A list of AMQP 1.0 brokers and other AMQP 1.0 resources may be found at [github.com/xinchen10/awesome-amqp](https://github.com/xinchen10/awesome-amqp).

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download [xk6](https://github.com/grafana/xk6):
  ```bash
  $ go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. [Build the k6 binary](https://github.com/grafana/xk6#command-usage):
  ```bash
  $ xk6 build --with github.com/ptsypyshev/xk6-amqp-custom@latest
  ```

## Development
To make development a little smoother, use the `Makefile` in the root folder. The default target will format your code, run tests, and create a `k6` binary with your local code rather than from GitHub.

```shell
git clone git@github.com:ptsypyshev/xk6-amqp-custom.git
cd xk6-amqp-custom
make
```

## Example

```javascript
import {Publisher, Consumer} from "k6/x/amqp"; // import Amqp extension

const url = "amqp://guest:guest@localhost:5672/"
const exchangeName = 'K6 exchange'
const queueName = 'K6 queue'
const routingKey = 'rkey123'

const exchange = {
  name: exchangeName,
  kind: "direct",
  durable: true,
};

const queue = {
  name: queueName,
  routing_key: routingKey,
  durable: true,
};

const publisher = new Publisher({
  connection_url: url,
  exchange: exchange,
  queue: queue,
});

const consumer = new Consumer({
  connection_url: url,
  exchange: exchange,
  queue: queue,
});

// init
export let options = {
  vus: 5,
  iterations: 10,
  // duration: '20s',
};

export default function () {
  console.log("Publisher and Consumer are ready")

  // Publish messages
  const publish = function(mark) {
    publisher.publish({
      exchange  : exchangeName,
      routing_key: routingKey,
      content_type: "text/plain",
      body: "Ping from k6 -> " + mark
    })
  }
  for (let i = 65; i <= 90; i++) {
    publish(String.fromCharCode(i))
}

  // Consume messages
  let result = consumer.consume({
    read_timeout: '3s',
    consume_limit: 26,
  })

  result.forEach(msg => {
    console.log("msg: " + msg.body)
  });  
}
```

Result output:

```plain
$ ./k6 run ./examples/publish-listen.js
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: ./examples/publish-listen.js
        output: -

     scenarios: (100.00%) 1 scenario, 5 max VUs, 10m30s max duration (incl. graceful stop):
              * default: 10 iterations shared among 5 VUs (maxDuration: 10m0s, gracefulStop: 30s)

INFO[0000] Publisher and Consumer are ready              source=console
INFO[0000] Publisher and Consumer are ready              source=console
INFO[0000] Publisher and Consumer are ready              source=console
INFO[0000] Publisher and Consumer are ready              source=console
INFO[0000] Publisher and Consumer are ready              source=console
INFO[0000] msg: Ping from k6 -> X                        source=console

... some the same console output ...

INFO[0000] msg: Ping from k6 -> F                        source=console

     consumer_message_count....: 156 49.882524/s
     consumer_queue_load.......: 0   min=0       max=90  
     data_received.............: 0 B 0 B/s
     data_sent.................: 0 B 0 B/s
     iteration_duration........: avg=1.26s min=62.09ms med=69.95ms max=3.06s p(90)=3.05s p(95)=3.06s
     iterations................: 10  3.197598/s
     publisher_message_count...: 260 83.13754/s
     publisher_queue_load......: 0   min=0       max=1539
     vus.......................: 4   min=4       max=4   
     vus_max...................: 5   min=5       max=5   


running (00m03.1s), 0/5 VUs, 10 complete and 0 interrupted iterations
default ✓ [======================================] 5 VUs  00m03.1s/10m0s  10/10 shared iters

```

Inspect examples folder for more details.

# Testing Locally

This repository includes a [docker-compose.yml](./docker-compose.yml) file that starts RabbitMQ with Management Plugin for testing the extension locally.

> :warning: This environment is intended for testing only and should not be used for production purposes.

1. Start the docker compose environment.
   ```bash
   docker compose up -d
   ```
   Output should appear similar to the following:
   ```shell
   ✔ Network xk6-amqp_default       Created               ...    0.0s
   ✔ Container xk6-amqp-rabbitmq-1  Started               ...    0.2s
   ```
2. Use your [custom k6 binary](#build) to run a k6 test script connecting to your RabbitMQ server started in the previous step.
   ```bash
   ./k6 run ./examples/publish-listen.js
   ```
3. Use the RabbitMQ admin console by accessing [http://localhost:15672/](http://localhost:15672/), then login using `guest` for both the Username and Password.
   This will allow you to monitor activity within your messaging server.
