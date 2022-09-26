# gc22-concurrent-web-apps-workshop
Code for the GC22 online workshop titled Building Concurrent Web Applications

# To Run Load Tests

```bash
$ go get -u github.com/rakyll/hey
$ hey -n 5 -c 2 -z 3s -m POST -T "application/json" -d ./body.txt http://localhost:3000/orders
```
