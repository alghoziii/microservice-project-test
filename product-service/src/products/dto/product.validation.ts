import { z } from 'zod';

export const ProductValidation = {
  PRODUCT: z.object({
    name: z.string().min(3, 'Name must be at least 3 characters'),
    price: z.number().positive('Price must be positive'),
    qty: z.number().int('Quantity must be integer').min(0, 'Quantity cannot be negative'),
  }),
};

export type ProductRequest = z.infer<typeof ProductValidation.PRODUCT>;