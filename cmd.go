package gcmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/PeterYangs/tools"
	"io"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

type Cmd struct {
	Command        string //命令行
	Cmd            *exec.Cmd
	CtxCancel      context.CancelFunc
	Ctx            context.Context
	OutPutBuf      *bufio.Reader   //命令输出流
	ErrPutBuf      *bufio.Reader   //错误输出流
	wait           *sync.WaitGroup //等待getOut处理完毕
	convertUtf8    bool            //是否输出转码成utf8
	customOutChan  chan []byte     //自定义命令输出管道
	customErrChan  chan []byte
	exportOut      bool //是否使用了自定义命令输出管道
	exportErr      bool
	waitCustomChan bool //是否等待自定义输出管道处理完毕
	customChanWait *sync.WaitGroup
	throwPanic     bool          //是否出错抛出异常
	outPut         *bytes.Buffer //输出返回值
	isOutPut       bool          //是否返回输出
}

func Command(command string) *Cmd {

	cmd := &Cmd{
		wait:           &sync.WaitGroup{},
		customOutChan:  make(chan []byte, 10),
		customErrChan:  make(chan []byte, 10),
		customChanWait: &sync.WaitGroup{},
		outPut:         bytes.NewBuffer(nil),
	}

	cmd.Command = command

	ctx, cancel := context.WithCancel(context.Background())

	cmd.CtxCancel = cancel

	cmd.Ctx = ctx

	return cmd
}

func dealCmd(c *Cmd) error {

	sysType := runtime.GOOS

	// linux/mac
	if sysType == "linux" || sysType == "darwin" {

		c.Cmd = exec.CommandContext(c.Ctx, "bash", "-c", c.Command)
	}

	// windows
	if sysType == "windows" {

		c.Cmd = exec.CommandContext(c.Ctx, "cmd", "/c", c.Command)

	}

	outPipe, err := c.Cmd.StdoutPipe()

	if err != nil {

		return err
	}

	c.OutPutBuf = bufio.NewReader(outPipe)

	errPipe, err := c.Cmd.StderrPipe()

	if err != nil {

		//fmt.Println(err)

		return err
	}
	c.ErrPutBuf = bufio.NewReader(errPipe)

	return nil
}

// SetTimeout cxt timeout
func (c *Cmd) SetTimeout(timeout time.Duration) *Cmd {

	cxt, cancel := context.WithTimeout(c.Ctx, timeout)

	c.Ctx = cxt

	c.CtxCancel = cancel

	return c

}

func (c *Cmd) OutPut() *Cmd {

	c.isOutPut = true

	return c
}

// ThrowPanic 出错抛出异常
func (c *Cmd) ThrowPanic() *Cmd {

	c.throwPanic = true

	return c
}

func (c *Cmd) WaitCustomChan() *Cmd {

	c.waitCustomChan = true

	return c
}

// Start run command
func (c *Cmd) Start() ([]byte, error) {

	//defer close(c.customOutChan)
	//defer close(c.customErrChan)

	defer func(cc *Cmd) {

		defer func() {

			if r := recover(); r != nil {
				//fmt.Printf("捕获到的错误：%s\n", r)
			}
		}()

		close(cc.customOutChan)
		close(cc.customErrChan)
		cc.outPut.Reset()

	}(c)

	var err error

	err = dealCmd(c)

	if err != nil {

		//fmt.Println(err)

		if c.throwPanic {

			panic(err)
		}

		return c.outPut.Bytes(), err
	}

	err = c.Cmd.Start()

	if err != nil {

		//fmt.Println(err)
		if c.throwPanic {

			panic(err)
		}

		return c.outPut.Bytes(), err
	}

	c.wait.Add(2)
	go getOut(c.OutPutBuf, 1, c)
	go getOut(c.ErrPutBuf, 2, c)

	//wait for command
	if err := c.Cmd.Wait(); err != nil {

		fmt.Println("wait err:", err)

		//panic()

		close(c.customOutChan)
		close(c.customErrChan)

		if c.throwPanic {

			panic(err)
		}

		return c.outPut.Bytes(), err

	}

	//wait for out done
	c.wait.Wait()

	if c.waitCustomChan {

		c.customChanWait.Wait()
	}

	return c.outPut.Bytes(), nil
}

func (c *Cmd) ConvertUtf8() *Cmd {

	c.convertUtf8 = true

	return c
}

// Done customchan done
func (c *Cmd) Done() {

	c.customChanWait.Done()
}

func (c *Cmd) GetCustomOutChan() chan []byte {

	c.exportOut = true

	return c.customOutChan
}

func (c *Cmd) GetCustomErrChan() chan []byte {

	c.exportErr = true

	return c.customErrChan
}

func getOut(outputBuf *bufio.Reader, types int, c *Cmd) {

	defer c.wait.Done()

	if types == 1 && c.exportOut && c.waitCustomChan {

		c.customChanWait.Add(1)
	}

	if types == 2 && c.exportErr && c.waitCustomChan {

		c.customChanWait.Add(1)
	}

	buf := make([]byte, 1024)

	for {

		n, err := outputBuf.Read(buf)

		if err != nil {

			if err == io.EOF {

				return
			}

			fmt.Println(err)
		}

		var out []byte
		if c.convertUtf8 {

			out = tools.ConvertToByte(string(buf[:n]), "gbk", "utf8")

		} else {

			out = buf[:n]
		}

		isOutToBash := true

		if types == 1 && c.exportOut {

			isOutToBash = false
			c.customOutChan <- out

		}

		if types == 2 && c.exportErr {
			isOutToBash = false
			c.customErrChan <- out
		}

		if isOutToBash && !c.isOutPut {

			fmt.Print(string(out))

		}

		if c.isOutPut {

			//fmt.Println(111)

			c.outPut.Write(out)
		}

	}

}
