from locust import HttpUser, task, between

class MyUser(HttpUser):
    wait_time = between(1, 2)  # Wait between 1 and 2 seconds between each task

    @task
    def get_users(self):
        user_id = 11   
        self.client.get(f"/gateway/get_users/{user_id}")

 
