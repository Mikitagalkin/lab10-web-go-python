"""
gRPC-клиент для Go-сервера.

Задание 7 повышенной сложности: gRPC.
"""

import grpc

import calculator_pb2
import calculator_pb2_grpc


def run():
    channel = grpc.insecure_channel("localhost:50051")
    stub = calculator_pb2_grpc.CalculatorStub(channel)

    # Add
    print("=== Add ===")
    for a, b in [(2, 3), (10, 20), (100, 200)]:
        req = calculator_pb2.AddRequest(a=a, b=b)
        resp = stub.Add(req)
        print(f"{a} + {b} = {resp.result}")

    # Multiply
    print("\n=== Multiply ===")
    for a, b in [(2, 3), (10, 20), (7, 8)]:
        req = calculator_pb2.MultiplyRequest(a=a, b=b)
        resp = stub.Multiply(req)
        print(f"{a} * {b} = {resp.result}")


if __name__ == "__main__":
    run()
