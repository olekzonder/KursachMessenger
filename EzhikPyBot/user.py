import json
import requests


def login(url, email, passwd):
    url += "/api/user/login"
    data = {"email": email, "password": passwd}
    data = json.dumps(data)
    r = requests.post(url, data)
    acc = r.json()
    if r.status_code == 403:
        if acc['message'] == 'INCORRECT_PASS':
            return 0
        if acc['message'] == "USER_NOT_FOUND":
            return 1
    if r.status_code != 500:
        return acc['account']['token']
    else:
        return -1
