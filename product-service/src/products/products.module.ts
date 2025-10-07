import { Module } from '@nestjs/common';
import { ProductsController } from './controller/products.controller';
import { ProductsService } from './service/products.service';
import { KafkaProducerService } from '../common/kafka/kafka-producer.service';
import { KafkaConsumerService } from '../common/kafka/kafka-consumer.service';
import { PrismaService } from '../common/prisma/prisma.service';

@Module({
  imports: [],
  controllers: [ProductsController],
  providers: [
    ProductsService,
    KafkaProducerService, 
    KafkaConsumerService,
    PrismaService,
  ],
  exports: [ProductsService],
})
export class ProductsModule {}