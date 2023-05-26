import sys
sys.path.append("../service1/pyhtonPb")

from locust import User, task, between
import grpc
import get_users_pb2
import get_users_pb2_grpc
from locust import TaskSet


class GrpcUser(User):
    wait_time = between(1, 3)

    def on_start(self):
        channel = grpc.insecure_channel("localhost:50051")
        self.client = get_users_pb2_grpc.get_usersStub(channel)

    @task
    def get_user(self):
        # Create a request message
        request = get_users_pb2.GetDataRequest(user_id=1)

        # Make request
        response = self.client.GetData(request)

        if response.message_id == 1:
            print(f"User: {response.return_users}")
        elif response.message_id == 3:
            for user in response.return_users:
                print(f"User: {user}")
    


if __name__ == "__main__":
    GrpcUser().run()
