import { Global, Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { APP_FILTER } from '@nestjs/core';
import { PrismaService } from './prisma/prisma.service';
import { RedisModule } from './redis/redis.module';
import { ErrorFilter } from './errorfilter/error.filter';
import { WinstonModule } from 'nest-winston';
import * as winston from 'winston';
import { ValidationService } from './validation/validation.service';

@Global()
@Module({
  imports: [
    WinstonModule.forRoot({
      level: 'debug',
      format: winston.format.json(),
      transports: [new winston.transports.Console()],
    }),
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    RedisModule,
  ],
  providers: [
    PrismaService,
    ValidationService,
    {
      provide: APP_FILTER,
      useClass: ErrorFilter, 
    },
  ],
  exports: [
    PrismaService, 
    ValidationService,
    WinstonModule, 
    RedisModule,  
  ],
})
export class CommonModule {}