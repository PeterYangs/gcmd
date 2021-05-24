# gcmd

This package is simple to run cmd.

### installation
```
go get github.com/PeterYangs/gcmd
```

### Quick start
```go
package main

import "github.com/PeterYangs/gcmd"

func main() {

	gcmd.Command("php index.php").Start()

}

```

php script is
```php
<?php

$index=0;

while (true){

    if($index>=10){

        throw new Exception("error here!");

    }

    echo 'echo success'.PHP_EOL;

    $index++;

    sleep(1);

}

```

Console output
```bash
echo success
echo success
echo success
echo success
echo success
echo success
echo success
echo success
echo success
echo success

Fatal error: Uncaught Exception: error here! in D:\goDemo\cmd\index.php:9
Stack trace:
#0 {main}
  thrown in D:\goDemo\cmd\index.php on line 9
wait err: exit status 255
```

### Custom output
```go
package main

import (
	"fmt"
	"github.com/PeterYangs/gcmd"
)

func main() {

	command := gcmd.Command("php index.php").WaitCustomChan()

	outChan := command.GetCustomOutChan()
	errChan := command.GetCustomErrChan()

	go func() {

		for t := range outChan {

			fmt.Println(string(t))
		}

		fmt.Println("out chan close")

		command.Done()

	}()

	go func() {

		for t := range errChan {

			fmt.Println("err:", string(t))
		}

		fmt.Println("err chan close")

		command.Done()

	}()

	command.Start()

}

```
### Throw Panic
```go
package main

import (
	"github.com/PeterYangs/gcmd"
)

func main() {

	gcmd.Command("php index.php").ThrowPanic().Start()

}

```

### Return Output
```go
package main

import (
	"fmt"
	"github.com/PeterYangs/gcmd"
)

func main() {

	out, err := gcmd.Command("php index.php").OutPut().Start()

	if err != nil {

		fmt.Println(err)

		return
	}

	fmt.Println(string(out))
}

```
