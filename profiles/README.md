## Анализ использования памяти

Симуляция нагрузки выполнялась с помощью: 

```hey -n 5000 -c 1 \                                                  
  -m POST \
  -H "Content-Type: application/json" \
  -d '[{"id":"a","type":"gauge","value":1.0},{"id":"b","type":"counter","delta":1}]' \
  http://localhost:8080/updates
```
Анализ профиля показал, что значительная часть памяти выделялась при
создании gzip-компрессоров (`compress/flate.NewWriter`), которые
создавались на каждый HTTP-запрос.

Был переработан GzipMiddleware с использованием
`sync.Pool` для переиспользования `gzip.Writer` и `gzip.Reader`,
что позволило уменьшить количество аллокаций и нагрузку на планировщик Go.

Сравнение профилей:

``` Showing nodes accounting for -234.12kB, 4.24% of 5515.24kB total ```