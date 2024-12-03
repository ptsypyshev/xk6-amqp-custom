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
