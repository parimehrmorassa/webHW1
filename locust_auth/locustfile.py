from locust import User, task, between, events
import grpc
import auth_pb2 as auth_pb2
import auth_pb2_grpc as auth_pb2_grpc
import req_DH_params_pb2 as dh_pb2
import req_DH_params_pb2_grpc as dh_pb2_grpc

class MyUser(User):
    wait_time = between(1, 5)
    response1 = None

    def on_start(self):
        channel = grpc.insecure_channel("localhost:50052")
        self.client = auth_pb2_grpc.MyServiceStub(channel)

        channel2 = grpc.insecure_channel("localhost:50054")
        self.client2 = dh_pb2_grpc.DHParamsServiceStub(channel2)

    @task
    def call_service_1(self):
        request = auth_pb2.MyRequest(
            message_id=4,
            nonce="abc12311111111111111"
        )
        response1 = self.client.ProcessRequest(request)

        events.request.fire(
            request_type="call_service_1",
            name="call_service_1",
            response_time=0,
            response_length=0,
        )

        self.response1 = response1

    @task
    def call_service_2(self):
        if self.response1 is not None:
            private_key = 3
            A = int(self.response1.g) ** private_key % int(self.response1.p)
            request = dh_pb2.DHParamsRequest(
                nonce=self.response1.nonce,
                server_nonce=self.response1.server_nonce,
                message_id=8,
                a=str(A)
            )
            response = self.client2.ProcessRequest(request)
            shared_key = int(response.b) ** private_key % int(self.response1.p)

            events.request.fire(
                request_type="call_service_2",
                name="call_service_2",
                response_time=0,
                response_length=0,
                response=shared_key,
            )
