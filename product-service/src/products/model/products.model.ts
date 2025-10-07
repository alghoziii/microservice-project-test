export interface ProductRequest {
  name: string;
  price: number;
  qty: number;
}

export interface ProductResponse {
  id: number;          
  name: string;
  price: number;
  qty: number;
  createdAt?: Date;     
}

export interface ProductEvent {
  id: number;
  name: string;
  price: number;
  qty: number;
  createdAt: Date;
}