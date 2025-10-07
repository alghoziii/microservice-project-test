import { Injectable, OnModuleInit, Logger } from '@nestjs/common';
import { Kafka, Consumer } from 'kafkajs';
import { ProductsService } from '../../products/service/products.service';

@Injectable()
export class KafkaConsumerService implements OnModuleInit {
  private readonly logger = new Logger(KafkaConsumerService.name);
  private consumer?: Consumer;
  private enabled = false;

  constructor(private readonly productsService: ProductsService) {
    const broker = process.env.KAFKA_BROKER?.trim();       // ← tidak ada default
    if (!broker) {
      this.logger.warn('Kafka disabled: KAFKA_BROKER is empty');
      return; // keep disabled
    }

    const kafka = new Kafka({
      clientId: 'product-service',
      brokers: [broker],
    });
    this.consumer = kafka.consumer({ groupId: 'product-service-group' });
    this.enabled = true;
  }

  async onModuleInit() {
    if (!this.enabled || !this.consumer) return;

    try {
      await this.consumer.connect();
      await this.consumer.subscribe({ topic: 'order.created', fromBeginning: true });

      await this.consumer.run({
        eachMessage: async ({ message }) => {
          if (!message?.value) return;
          try {
            const evt = JSON.parse(message.value.toString());
            this.logger.log(`📨 Received order.created: ${JSON.stringify(evt)}`);
            await this.productsService.reduceQuantity(evt.productId, evt.quantity);
          } catch (e) {
            this.logger.warn(`process message failed (ignored): ${e}`);
          }
        },
      });

      this.logger.log('✅ Kafka consumer started for order.created');
    } catch (e) {
      this.enabled = false;
      this.logger.warn(`Kafka disabled (consumer connect failed): ${e}`);
    }
  }
}
