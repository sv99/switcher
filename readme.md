Switcher
========

Переключение провайдеров на Mikrotik. Нужно подлючиться к нему по SSH
и выполнить на нем подготовленный скрипт. 

Нужные для работы файлы:

    VIDEODIR
    ├── switcher.exe
    ├── switcher.conf
    ├── *.dsa
    ├── *.dsa.pub (нужен для загрузки на mikrotik)
    ├── static
    │   ├── favicon.ico
    │   ├── etelecom_logo.png
    │   ├── sumtel_logo.png
    │   └── *.js
    ├── favicon.ico
    └── index.html
    
dependencies using dep
----------------------

Не хранит историю git - только актуалные файлы. В результате получаем очень
компактный размер папки vendor.
   
    brew install dep
    dep init
    # -v show extended log
    dep ensure -v

ssh connect to mikrotik
-----------------------

Mikrotik разрешает коннектиться только по DSA ключам.

    # generate key
    ssh-keygen -f mikrotik.dsa -t dsa

On Mikrotik

    /user add name=switcherUser group=switcher  disabled=no
    # switcher group write, read
    /user ssh-keys import public-key-file=mikrotik.dsa.pub user=switcherUser
    # check connect
    ssh -i mikrotik.dsa scriptUser@192.168.1.202 “/system resource print”

Чтобы заработала DSA на последней OSX with OpenSSH 7 нужно разрешить
использовать

    sudo echo PubkeyAcceptedKeyTypes +ssh-dss >> ~/.ssh/config
    chmod 600 ~/.ssh/config

Дополнительно пришлось разрешить подключаться на 22 порт из локальной
сети минуя цепочку fail2ban. В приложении не храню SSH Session - 
каждый запрос отдельное подключение!

scripts
-------

script          | Result
--------------- | ---------------------------------
get_version     | return model and RouterOS version
get_provider    | return 1 or 2 - default provider
switch_provider | return 1 or 2 or "error"

Передать параметр в скрипт нельзя, для переключения один
скрип, который возвращает новый текущий провайдер.

client
------

[CSS botstrap 4](https://getbootstrap.com/)

[VueJS](https://vuejs.org/)

[AJAX library Axios](https://github.com/axios/axios)

Одна страница webpack не используется.