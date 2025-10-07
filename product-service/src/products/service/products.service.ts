import { Injectable, Inject, Logger } from '@nestjs/common';
import { Product } from '@prisma/client';
import { Redis } from 'ioredis';
import { PrismaService } from '../../common/prisma/prisma.service';
import { KafkaProducerService } from '../../common/kafka/kafka-producer.service';
import { ProductRequest } from '../dto/product.validation';

@Injectable()
export class ProductsService {
  private readonly logger = new Logger(ProductsService.name);

  constructor(
    @Inject('REDIS_CLIENT') private readonly redis: Redis,
    private readonly prisma: PrismaService,
    private readonly kafka: KafkaProducerService,
  ) {}

  // 🔹 CREATE PRODUCT
  async create(data: ProductRequest): Promise<Product> {
    try {
      const product = await this.prisma.product.create({ data });
      await this.kafka.sendMessage('product.created', product);
      await this.invalidateCache(product.id);
      return product;
    } catch (e) {
      this.logger.error('Failed to create product', e);
      throw e;
    }
  }

  // 🔹 FIND PRODUCT BY ID (with cache)
  async findOne(id: number): Promise<Product | null> {
    const key = `product:${id}`;

    // Try cache first
    const cached = await this.redis.get(key);
    if (cached) return JSON.parse(cached);

    // Fetch from DB
    const product = await this.prisma.product.findUnique({ where: { id } });
    if (product) await this.redis.set(key, JSON.stringify(product), 'EX', 300);

    return product;
  }

  // 🔹 FIND ALL PRODUCTS
  async findAll(): Promise<Product[]> {
    return this.prisma.product.findMany();
  }

  // 🔹 REDUCE STOCK
  async reduceQuantity(productId: number, qty: number): Promise<void> {
    await this.prisma.product.update({
      where: { id: productId },
      data: { qty: { decrement: qty } },
    });
    await this.invalidateCache(productId);
  }

  // 🔹 Helper invalidate cache
  private async invalidateCache(id: number) {
    await this.redis.del(`product:${id}`);
  }
}
