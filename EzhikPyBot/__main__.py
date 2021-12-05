import asyncio
import json
from datetime import datetime
import user
import requests
from pathlib import Path

jsonsettings = Path('settings.json')
if jsonsettings.is_file():
    with open('settings.json') as settings:
        d = json.load(settings)
        IP = d["ip"]
        Port = d["port"]
        email = d["email"]
        passwd = d["passwd"]
    url = "http://" + IP + ":" + Port
    if requests.get(url + '/').status_code != 200:
        print("Не удалось установить связь с сервером...")
        exit(0)
    token = user.login(url, email, passwd)
else:
    print("Файл settings.json отсутствует.")
    IP = input("Введите IP-адрес сервера > ")
    Port = input("Введите порт сервера (по умолчанию 8080) > ")
    if Port == '':
        Port = "8080"
    url = "http://" + IP + ":" + Port
    if requests.get(url + '/').status_code != 200:
        print("Не удалось установить связь с сервером... ")
        exit(0)
    k = input("Хотите войти или зарегистрироваться? (l/r) > ")
    while k.upper() != "L" and k.upper() != "R":
        k = input("Хотите войти или зарегистрироваться? (l/r) > ")
    if k.upper() == "L":
        email = input("Введите адрес электронной почты > ")
        passwd = input("Введите пароль > ")
        token = user.login(url, email, passwd)
        l = input("Хотите записать файл settings.json? (y/n) > ")
        while l.upper() != "Y" and l.upper() != "N":
            l = input("Хотите записать файл settings.json? (y/n) > ")
        if l.upper() == "Y":
            with open('settings.json', 'w') as f:
                json.dump({'ip': IP, 'port': Port, 'email': email, 'passwd': passwd}, f)
    else:
        email = input("Введите email > ")
        passwd = input("Введите пароль (минимум 6 символов) > ")
        while len(passwd) <= 6:
            passwd = input("Введите пароль (минимум 6 символов) > ")
        url_reg = url+"/api/user/create"
        data = json.dumps({"email": email, "password": passwd})
        r = requests.post(url_reg, data)
        if r.status_code == 200:
            token = user.login(url, email, passwd)
        l = input("Хотите записать файл settings.json? (y/n) > ")
        while l.upper() != "Y" and l.upper() != "N":
            l = input("Хотите записать файл settings.json? (y/n) > ")
        if l.upper() == "Y":
            with open('settings.json', 'w') as f:
                json.dump({'ip': IP, 'port': Port, 'email': email, 'passwd': passwd}, f)
        else:
            token = 2


if token == 0:
    print("Неверный пароль!")
    exit(0)
if token == 1:
    print("Пользователь не найден!")
    exit(0)
if token == 2:
    print("Регистрация не была успешной...")
if token == -1:
    print("Неполадки с сервером...")
    exit(0)

header = {"Authorization": "Bearer " + token}  # Заголовок авторизации
r = requests.get(url + '/api/messages/get', headers=header)  # Очистка сообщений
print("Готов к работе!")

payload = {"msg": "БОТ ЁЖИК 1.0 К ВАШИМ УСЛУГАМ."}
payload = json.dumps(payload)
requests.post(url + '/api/messages/send', payload, headers=header)


async def async_read(url, header):
    timer = asyncio.get_event_loop().time()
    last = timer
    interval = 1.5
    while True:
        await asyncio.sleep(interval - int(timer) % interval)
        r = requests.get(url + '/api/messages/get/unread', headers=header)
        msgs = r.json()
        if str(msgs["message"]) != "NO_NEW_MSG":
            for i in range(int(msgs["message"])):
                if msgs[str(i)]["msg"] == "/help":
                    print(msgs[str(i)]["name"] + ": /help")
                    payload = json.dumps({
                        "msg": "\n Список команд: \n /help - вывести это вообщение \n /time- вывести время \n /who - "
                               "вывести ""имя пользователя \n /about - о боте"})
                    r = requests.post(url + '/api/messages/send', payload, headers=header)
                if msgs[str(i)]["msg"] == "/who":
                    print(msgs[str(i)]["name"] + ": /who")
                    payload = json.dumps({"msg": "Ты " + str(msgs[str(i)]["name"])})
                    r = requests.post(url + '/api/messages/send', payload, headers=header)
                if msgs[str(i)]["msg"] == "/time":
                    print(msgs[str(i)]["name"] + ": /time")
                    payload = json.dumps({"msg": str(datetime.now())})
                    r = requests.post(url + '/api/messages/send', payload, headers=header)
                if msgs[str(i)]["msg"] == "/about":
                    print(msgs[str(i)]["name"] + ": /about")
                    payload = json.dumps({"msg": "\n ЁЖИК БОТ v 1.0 \n ПЕРЕДАЮ ПРИВЕТ ВСЕМ, КТО МЕНЯ ЗНАЕТ"})
                    r = requests.post(url + '/api/messages/send', payload, headers=header)


asyncio.run(async_read(url, header))
