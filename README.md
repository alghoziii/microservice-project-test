# Fullstack Developer Test Challenge

Microservices demo with **NestJS (Product Service)** and **Golang (Order Service)**, using **MySQL**, **Redis**, and **Kafka** (event-driven).

---

## Table of Contents
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Default Ports](#default-ports)
- [Environment Variables](#environment-variables)
- [API Reference](#api-reference)
- [Postman](#postman)
- [Performance (k6)](#performance-k6)
- [Project Structure](#project-structure)
- [Troubleshooting](#troubleshooting)

---

## Architecture

```mermaid
flowchart LR
  subgraph Product Service (NestJS)
    PAPI[REST API]
    PDB[(MySQL: products)]
    PREDIS[(Redis)]
  end

  subgraph Order Service (Go)
    OAPI[REST API]
    ODB[(MySQL: orders)]
    OREDIS[(Redis)]
  end

  K[(Kafka)]
  Z[(Zookeeper)]

  PAPI-- GET/POST -->PDB
  PAPI-- cache -->PREDIS
  OAPI-- GET/POST -->ODB
  OAPI-- cache -->OREDIS

  OAPI-- publish order.created -->K
  K-- consume order.created -->PAPI
  PAPI-- publish product.created -->K

  Z---K
