package admin

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"

	"<projectrepo>/pkg/model/core"
	"github.com/gin-gonic/gin"
	mailjet "github.com/mailjet/mailjet-apiv3-go/v3"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func NotifyPageFailure(c *gin.Context, err any, user *core.User) {
	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))

	decodedStack := strings.Split(ClonedCustomRecovery(c, err), "\r\n")

	var userInfo string = "Anonymous"

	if user != nil {
		userInfo = fmt.Sprintf("email (%s), id (%s)", user.Email, user.ID)
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "<projectemail>",
				Name:  "Your <projectname>",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: NotifyAddress,
				},
			},
			Subject: "Panic on the page",
			TextPart: fmt.Sprintf(`
			Hi!

			Panic on the page

			* User: %s
			* Request data:

			%s
			`, userInfo, strings.Join(decodedStack, "\r\n")),
			HTMLPart: fmt.Sprintf(`
			<p>Hi!</p>

			<p>Panic on the page:</p>

			<ul>
			<li>user: %s</li>
			<li>Request data: <br /><pre>%s</pre></li>
			</ul>`, userInfo, strings.Join(decodedStack, "\r\n")),
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)

	if err != nil {
		log.Fatal(err)
	}
}

func ClonedCustomRecovery(c *gin.Context, err any) string {
	// Check for a broken connection, as it is not really a
	// condition that warrants a panic stack trace.
	var brokenPipe bool
	if ne, ok := err.(*net.OpError); ok {
		var se *os.SyscallError
		if errors.As(ne, &se) {
			seStr := strings.ToLower(se.Error())
			if strings.Contains(seStr, "broken pipe") ||
				strings.Contains(seStr, "connection reset by peer") {
				brokenPipe = true
			}
		}
	}

	stack := stack(3)
	httpRequest, _ := httputil.DumpRequest(c.Request, false)
	headers := strings.Split(string(httpRequest), "\r\n")
	for idx, header := range headers {
		current := strings.Split(header, ":")
		if current[0] == "Cookie" && !gin.IsDebugging() {
			headers[idx] = current[0] + ": <hidden>"
		}

		if current[0] == "Authorization" {
			headers[idx] = current[0] + ": *"
		}
	}
	headersToStr := strings.Join(headers, "\r\n")
	if brokenPipe {
		return fmt.Sprintf("%s\n%s", err, headersToStr)
	} else {
		return fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n%s",
			timeFormat(time.Now()), headersToStr, err, stack)
	}
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contain dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// timeFormat returns a customized time string for logger.
func timeFormat(t time.Time) string {
	return t.Format("2006/01/02 - 15:04:05")
}
