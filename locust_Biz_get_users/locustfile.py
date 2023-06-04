from locust import User, task, between, events
import grpc
import auth_pb2 as auth_pb2
import auth_pb2_grpc as auth_pb2_grpc
import req_DH_params_pb2 as dh_pb2
import req_DH_params_pb2_grpc as dh_pb2_grpc
import get_users_pb2 as get_user_pb2
import get_users_pb2_grpc as get_user_pb2_grpc
import base64

class MyUser(User):
    wait_time = between(1, 5)
    response1 = None
    shared_key = None
    redis_keyy=None
    def on_start(self):
        channel = grpc.insecure_channel("localhost:50052")
        self.client = auth_pb2_grpc.MyServiceStub(channel)

        channel2 = grpc.insecure_channel("localhost:50054")
        self.client2 = dh_pb2_grpc.DHParamsServiceStub(channel2)

        # Call auth_req_pq
        self.call_service_1()

        # Call auth_DH_params
        if self.response1 is not None:
            self.call_service_2()

    def call_service_1(self):
        request = auth_pb2.MyRequest(
            message_id=4,
            nonce="abc12311111111111111"
        )
        response1 = self.client.ProcessRequest(request)

        events.request.fire(
            request_type="call_service_1",
            name="auth_req_pq",
            response_time=0,
            response_length=0,
        )

        self.response1 = response1

    def call_service_2(self):
        private_key = 3
        A = int(self.response1.g) ** private_key % int(self.response1.p)
        request = dh_pb2.DHParamsRequest(
            nonce=self.response1.nonce,
            server_nonce=self.response1.server_nonce,
            message_id=8,
            a=str(A)
        )
        response        = self.client2.ProcessRequest(request)
        shared_key      = int(response.b) ** private_key % int(self.response1.p)
        self.shared_key = shared_key
        self.redis_keyy = self.response1.nonce+":"+self.response1.server_nonce
        events.request.fire(
            request_type="call_service_2",
            name="auth_DH_params",
            response_time=0,
            response_length=0,
            response=shared_key,
        )

    @task
    def call_get_users(self):

        channel = grpc.insecure_channel("localhost:50051")
        stub = get_user_pb2_grpc.get_usersStub(channel)


        shared_key_bytes = self.shared_key.to_bytes((self.shared_key.bit_length() + 7) // 8, 'big')

        request = get_user_pb2.GetDataRequest(
            user_id=123,
            auth_key=shared_key_bytes,
            message_id=1,
            redis_key=self.redis_keyy
        )

        # Call get_users service
        response = stub.GetData(request)

        events.request.fire(
            request_type="call_get_users",
            name="call_get_users",
            response_time=0,
            response_length=0,
        )
