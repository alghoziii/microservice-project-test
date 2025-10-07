import { Module } from '@nestjs/common';
import { CommonModule } from './common/common.module';
import { ProductsModule } from './products/products.module';

@Module({
  imports: [CommonModule, ProductsModule],
})
export class AppModule {}