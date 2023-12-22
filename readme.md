Switcher
========

**3.08.2020**

Переключение провайдеров на Mikrotik. Нужно подлючиться к нему по SSH
и выполнить на нем подготовленный скрипт.

Нужные для работы файлы:

    VIDEODIR
    ├── switcher.exe
    ├── switcher.conf
    ├── *.dsa
    ├── *.dsa.pub (нужен для подключения к mikrotik)
    └── static
        ├── favicon.ico
        ├── etelecom_logo.png
        ├── sumtel_logo.png
        ├── axios.min.js
        ├── bootstrap.min.css
        ├── vue-spinner.min.js
        └── vue.min.js
    
server
------

Перешел на [Fiber](https://github.com/gofiber/fiber) и на `mod` для
управления зависимостями.

```bash   
go mod init switcher
go mod tidy
```

Сборка на базе `Makefile`.

`index.html` - из сгенерированного статического файла. Остальная статика из
каталога. В настоящее время не нашел вариант работы со `static assets` для `fiber` 

Логгер [zerolog](https://github.com/rs/zerolog)

Для windows собираю сервис, который принимает параметры для регистрации и запуска.
Сервис нормально работает под системной учеткой, устанавливать сетевую учетку не потребовалось.

Первую версию запускал при помощи NSSM.

ssh connect to mikrotik
-----------------------

Старые версии Mikrotik работали с DSA ключами. После обновления OpenSSH до версии 7.0 стали работать с RSA ключами.

    # generate key
    ssh-keygen -f mikrotik.dsa -t dsa
    ssh-keygen -f mikrotik.rsa -t rsa

On Mikrotik

    /user add name=switcherUser group=switcher  disabled=no
    # switcher group write, read
    /user ssh-keys import public-key-file=mikrotik.rsa.pub user=switcherUser
    # check connect
    ssh -i mikrotik.rsa switcherUser@192.168.1.202 “/system resource print”

Чтобы заработала DSA на последней OSX with OpenSSH 7 нужно разрешить
использовать

    sudo echo PubkeyAcceptedKeyTypes +ssh-dss >> ~/.ssh/config
    chmod 600 ~/.ssh/config

Дополнительно пришлось разрешить подключаться на 22 порт из локальной
сети минуя цепочку fail2ban. В приложении не храню SSH Session - 
каждый запрос отдельное подключение!

mikrotik scripts
----------------

script          | Result
--------------- | ---------------------------------
get_version     | return model and RouterOS version
get_provider    | return 1 or 2 - default provider
switch_provider | return 1 or 2 or "error"

Передать параметр в скрипт нельзя, для переключения один
скрипт, который возвращает новый текущий провайдер.

client
------

[CSS botstrap 4](https://getbootstrap.com/)

[VueJS](https://vuejs.org/) версия v2.5.17 обновил 8.10.2017

[vue-spinner](https://github.com/greyby/vue-spinner)

webpack не используется - локальное приложение прямо на странице.

[AJAX library Axios](https://github.com/axios/axios) заменил на обертку для `window.fetch()`

Контроль и управление DUNE
--------------------------

Задача отображения состояния и вывод из спящего режима.

dune_ip_control_overview.pdf - описание API.

[Dune Remote Control](http://dune-hd.com/support/rc/)

Добавил дополнительные кнопки для отображения статуса Dune HD.

Динамически подгружаем их имена и IP адреса.

History
-------

**Февраль 2018** старт проекта

На базе [iris](https://github.com/kataras/iris) bootstrap and vue

**Октябрь 2018** старт проекта

Добавление отображения состояния Dune и возможность их пробуждения 

**Июль 2020** обновление версий библиотек

**Август 2020** переход на fiber, добавление windows сервиса

1. Сделал другие проекты и перевел этот тоже.
2. Обновление фронтенда

**Декабрь 2021** переход на fiber v2





