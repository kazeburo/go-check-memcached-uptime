package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type memcachedSetting struct {
	Host    string  `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Port    string  `short:"p" long:"port" default:"11211" description:"Port"`
	Timeout float64 `short:"t" long:"timeout" default:"10" description:"Seconds before connection times out"`
}

type uptimeOpts struct {
	memcachedSetting
	Crit int64 `short:"c" long:"critical" description:"critical if uptime seconds is less than this number"`
	Warn int64 `short:"w" long:"warning" description:"warning if uptime seconds is less than this number"`
}

func uptime2str(uptime int64) string {
	day := uptime / 86400
	hour := (uptime % 86400) / 3600
	min := ((uptime % 86400) % 3600) / 60
	sec := ((uptime % 86400) % 3600) % 60
	return fmt.Sprintf("%d days, %02d:%02d:%02d", day, hour, min, sec)
}

func main() {
	ckr := checkUptime()
	ckr.Name = "memcached Uptime"
	ckr.Exit()
}

func write(conn net.Conn, content []byte, timeout float64) error {
	if timeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
	_, err := conn.Write(content)
	return err
}

func slurp(conn net.Conn, timeout float64) ([]byte, error) {
	buf := []byte{}
	readLimit := 32 * 1024
	if timeout > 0 {
		conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
	for {
		tmpBuf := make([]byte, readLimit)
		i, err := conn.Read(tmpBuf)
		if i > 0 {
			buf = append(buf, tmpBuf[:i]...)
		}
		if err == io.EOF || i < readLimit {
			return buf, nil
		}
		if err != nil {
			return buf, err
		}
	}
	return buf, nil
}

func retrieve_uptime(conn net.Conn, timeout float64) (int64, error) {
	buf, err := slurp(conn, timeout)
	if err != nil {
		return 0, err
	}

	for _, b := range bytes.Split(buf, []byte("\n")) {
		if match := regexp.MustCompile(`^STAT uptime (\d+)`).FindStringSubmatch(string(b)); match != nil {
			i, err := strconv.ParseInt(match[1], 0, 64)
			if err != nil {
				return 0, err
			}
			return i, nil
		}
	}

	return 0, fmt.Errorf("uptime not found")
}

func checkUptime() *checkers.Checker {
	opts := uptimeOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	_, err := psr.Parse()
	if err != nil {
		os.Exit(1)
	}

	address := fmt.Sprintf("%s:%s", opts.Host, opts.Port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return checkers.Critical(err.Error())
	}
	defer conn.Close()

	err = write(conn, []byte("stats\r\n"), opts.Timeout)
	if err != nil {
		return checkers.Critical(err.Error())
	}

	Uptime, err := retrieve_uptime(conn, opts.Timeout)
	if err != nil {
		return checkers.Critical(err.Error())
	}

	if opts.Crit > 0 && Uptime < opts.Crit {
		return checkers.Critical(fmt.Sprintf("up %s < %s", uptime2str(Uptime), uptime2str(opts.Crit)))
	} else if opts.Warn > 0 && Uptime < opts.Warn {
		return checkers.Warning(fmt.Sprintf("up %s < %s", uptime2str(Uptime), uptime2str(opts.Warn)))
	}
	return checkers.Ok(fmt.Sprintf("up %s", uptime2str(Uptime)))

}
