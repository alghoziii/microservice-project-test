import { Controller, Post, Get, Body, Param, ParseIntPipe, Logger, Patch } from '@nestjs/common';
import { ProductsService } from '../service/products.service';
import { ProductValidation, ProductRequest } from '../dto/product.validation';

@Controller('products')
export class ProductsController {
  private readonly logger = new Logger(ProductsController.name);

  constructor(private readonly productsService: ProductsService) {}

  @Post()
  async create(@Body() body: any) {
    this.logger.log('POST /products - Creating product');
    
    // Validasi dengan Zod
    const productRequest: ProductRequest = ProductValidation.PRODUCT.parse(body);
    
    return this.productsService.create(productRequest);
  }

  @Get(':id')
  findOne(@Param('id', ParseIntPipe) id: number) {
    this.logger.log(`GET /products/${id} - Finding product`);
    return this.productsService.findOne(id);
  }

  @Get()
  findAll() {
    this.logger.log('GET /products - Finding all products');
    return this.productsService.findAll();
  }

  @Patch(':id/reduce-quantity')
  async reduceQuantity(
    @Param('id', ParseIntPipe) id: number,
    @Body() body: { quantity: number }
  ) {
    this.logger.log(`PATCH /products/${id}/reduce-quantity - Reducing quantity by ${body.quantity}`);
    
    await this.productsService.reduceQuantity(id, body.quantity);
    return { message: 'Quantity reduced successfully' };
  }
}