scripts
=======

get_version
-----------
```lua
﻿:local boardname [/system resource get board-name]
:local version [/system resource get version]

:put "Mikrotik: $boardname RouterOS: $version"
```

get_provider
------------

debug version - работаем с глобальной переменной

```lua
﻿:global provider

:if ( $provider=nil ) do={ :set provider 1; }

:put $provider
```

worked version

```lua
﻿# find distance for isp1 by comment
﻿:local distance

:foreach counter in=[/ip route find where comment="isp1"] do={
 :set $distance [/ip route get $counter distance]
}

#:put $distance
:if ( $distance != 1) do={ :put 2 } else={ :put 1 } 
```

switch_provider
---------------

debug version

```lua
﻿:do {
:global provider

:if ( $provider != "1" ) do={ :set provider 1 } \
else={ :set provider 2 }

:delay 1
:put $provider
} on-error={ :put "error"}
```

worked version
 
Для переключения провайдеров меняем distance route (ищем нужный по соментарию).
Для переключения текущих соединений на 1 секунду отключаем интерфейс
текущего провайдера (ищем нужный по имени интерфейса).

```lua
﻿# switch distance for isp1 find by comment
# switch interface find by name
:do {

:local distance
:local newdistance
:local currentisp

:foreach counter in=[/ip route find where comment="isp1"] do={
 :set $distance [/ip route get $counter distance]
}

:if ( $distance != 1) do={
  :set newdistance 1
  :set currentisp "isp2"
} else={
  :set newdistance 4
  :set currentisp "isp1"
}

# set new distance
:foreach counter in=[/ip route find where comment="isp1"] do={
 /ip route set $counter distance $newdistance
}

# switch on-off old default interface for reset connection
:foreach counter in=[/interface ethernet find where name=$currentisp] do={
 /interface ethernet set $counter disabled=yes
}
:delay 1
:foreach counter in=[/interface ethernet find where name=$currentisp] do={
 /interface ethernet set $counter disabled=no
}

#:put $currentisp
#:put $newdistance
:if ( $newdistance != 1 ) do={ :put 2 } else={ :put 1 } 
} on-error={ :put "error" }
```