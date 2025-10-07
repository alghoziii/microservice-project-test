import { Test } from '@nestjs/testing';
import { ProductsService } from './products.service';
import { PrismaService } from '../../common/prisma/prisma.service';
import { KafkaProducerService } from '../../common/kafka/kafka-producer.service';

describe('ProductsService', () => {
  let service: ProductsService;

  const redisMock = {
    get: jest.fn(),
    set: jest.fn(),
    del: jest.fn(),
  };

  const prismaMock = {
    product: {
      create: jest.fn(),
      findUnique: jest.fn(),
      update: jest.fn(),
    },
  };

  const kafkaMock = {
    sendMessage: jest.fn(),
  };

  beforeEach(async () => {
    const module = await Test.createTestingModule({
      providers: [
        ProductsService,
        { provide: 'REDIS_CLIENT', useValue: redisMock },
        { provide: PrismaService, useValue: prismaMock },
        { provide: KafkaProducerService, useValue: kafkaMock },
      ],
    }).compile();

    service = module.get(ProductsService);
    jest.clearAllMocks();
  });

  it('create() menyimpan produk & kirim event Kafka', async () => {
    const data = { name: 'Test Product', price: 10000, qty: 5 };
    const created = { id: 1, ...data };

    prismaMock.product.create.mockResolvedValue(created);

    const result = await service.create(data as any);

    expect(prismaMock.product.create).toHaveBeenCalledWith({ data });
    expect(kafkaMock.sendMessage).toHaveBeenCalledWith('product.created', created);
    expect(redisMock.del).toHaveBeenCalledWith('product:1');
    expect(result).toEqual(created);
  });

  it('findOne() cache HIT (tidak panggil DB)', async () => {
    const cached = { id: 9, name: 'Cached Product' };
    redisMock.get.mockResolvedValue(JSON.stringify(cached));

    const result = await service.findOne(9);

    expect(redisMock.get).toHaveBeenCalledWith('product:9');
    expect(result).toEqual(cached);
    expect(prismaMock.product.findUnique).not.toHaveBeenCalled();
  });

  it('reduceQuantity() update stok & hapus cache', async () => {
    prismaMock.product.update.mockResolvedValue(true);

    await service.reduceQuantity(7, 2);

    expect(prismaMock.product.update).toHaveBeenCalledWith({
      where: { id: 7 },
      data: { qty: { decrement: 2 } },
    });
    expect(redisMock.del).toHaveBeenCalledWith('product:7');
  });
});
