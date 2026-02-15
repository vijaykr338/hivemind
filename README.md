# Hivemind

Hivemind is a distributed peer-to-peer compute scheduler that aggregates multiple machines into a single virtual execution cluster.

It allows users to submit a computational job which is automatically split into smaller tasks, distributed across available worker nodes, executed in parallel, and reassembled into a final result.

The system continuously monitors worker liveness using heartbeats, dynamically balances load across heterogeneous machines, and recovers from node failures by reassigning unfinished work. The goal of the project is to demonstrate real-world distributed system behavior including concurrency, fault tolerance, and scheduling rather than simple remote execution.

---

## Features

- Dynamic worker discovery and liveness monitoring
- Pull-based distributed task scheduling
- Parallel execution across multiple machines
- Automatic load balancing
- Failure detection and task reassignment
- Self-healing workers with automatic reconnection
- gRPC-based communication
- Designed for containerized sandbox execution (planned)

---

## Architecture Overview

The system consists of three components.

### Coordinator

The central scheduler responsible for:

- tracking active workers
- assigning tasks
- monitoring heartbeats
- reassigning failed tasks
- aggregating results

### Worker

A lightweight agent running on any machine that:

- registers with the coordinator
- sends periodic heartbeats
- requests tasks
- executes computations
- returns results
- automatically reconnects on failure

### Client

Submits jobs to the system and receives final results.

---

## How It Works

1. Workers start and register with the coordinator
2. The coordinator tracks worker liveness using heartbeats
3. A job is submitted and split into smaller tasks
4. Workers pull tasks from the coordinator
5. Each worker executes its assigned task
6. Results are returned and merged into a final output
7. If a worker disconnects, unfinished tasks are reassigned

---

## Tech Stack

- Go
- gRPC / Protocol Buffers
- Concurrency via goroutines and channels
- Planned: Redis for persistent task state
- Planned: Docker sandboxed execution

---

## Running the Project

### Start Coordinator

```bash
go run coordinator/main.go
