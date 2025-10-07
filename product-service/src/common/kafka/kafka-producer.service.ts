import { Injectable, OnModuleInit, OnModuleDestroy, Logger } from '@nestjs/common';
import { Kafka, Producer } from 'kafkajs';

@Injectable()
export class KafkaProducerService implements OnModuleInit, OnModuleDestroy {
  private readonly logger = new Logger(KafkaProducerService.name);
  private producer?: Producer;
  private enabled = false;

  constructor() {
    const broker = process.env.KAFKA_BROKER?.trim();       
    if (!broker) {
      this.logger.warn('Kafka disabled: KAFKA_BROKER is empty');
      return; // keep disabled
    }

    const kafka = new Kafka({
      clientId: 'product-service',
      brokers: [broker],
    });

    this.producer = kafka.producer();
    this.enabled = true;
  }

  async onModuleInit() {
    if (!this.enabled || !this.producer) return;
    try {
      await this.producer.connect();
      this.logger.log('✅ Kafka producer connected');
    } catch (e) {
      this.enabled = false;
      this.logger.warn(`Kafka disabled (producer connect failed): ${e}`);
    }
  }

  async onModuleDestroy() {
    if (!this.enabled || !this.producer) return;
    try {
      await this.producer.disconnect();
      this.logger.log('✅ Kafka producer disconnected');
    } catch {
      /* ignore */
    }
  }

  async sendMessage(topic: string, message: any) {
    if (!this.enabled || !this.producer) {
      this.logger.debug(`Kafka disabled, skip send to ${topic}`);
      return;
    }
    try {
      await this.producer.send({
        topic,
        messages: [{ value: JSON.stringify(message) }],
      });
      this.logger.log(`📤 Message sent to ${topic}`);
    } catch (e) {
      this.logger.warn(`Kafka send failed (ignored): ${e}`);
    }
  }
}
