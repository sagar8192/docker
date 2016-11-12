package scribe

import (
	// "fmt"
	"net"
	"time"

	"github.com/docker/docker/daemon/logger"

	"github.com/Sirupsen/logrus"
	"github.com/samuel/go-thrift/examples/scribe"
    	"github.com/samuel/go-thrift/thrift"
	// "golang.org/x/net/context"
)

const (
	name = "scribe"
	ip = "127.0.0.1"
	port = "1463"
	streamName = "sagarp-testing-docker-logging"
)

/*
var (
	streamName   string
)
*/

func init() {

	if err := logger.RegisterLogDriver(name, New); err != nil {
		logrus.Fatal(err)
	}

	// if err := logger.RegisterLogOptValidator(name, ValidateLogOpts); err != nil {
	//	logrus.Fatal(err)
	// }
}

type scribelogs struct {
	client    scribe.ScribeClient
	container *containerInfo
}

type containerInfo struct {
	Name      string            `json:"name,omitempty"`
	ID        string            `json:"id,omitempty"`
	ImageName string            `json:"imageName,omitempty"`
	ImageID   string            `json:"imageId,omitempty"`
	Created   time.Time         `json:"created,omitempty"`
	Command   string            `json:"command,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Creates a new logger for scribe
func New(ctx logger.Context) (logger.Logger, error) {

	// Initialize scribe client
	conn, _ := net.Dial("tcp", ip+":"+port)
    	t := thrift.NewTransport(thrift.NewFramedReadWriteCloser(conn, 0), thrift.BinaryProtocol)
    	c := thrift.NewClient(t, false)

	l := &scribelogs{
		client: scribe.ScribeClient{Client: c},
		container: &containerInfo{
			Name:      ctx.ContainerName,
			ID:        ctx.ContainerID,
			ImageName: ctx.ContainerImageName,
			ImageID:   ctx.ContainerImageID,
			Created:   ctx.ContainerCreated,
			Metadata:  ctx.ExtraAttributes(nil),
		},
	}
	
	return l, nil
}


/*
// ValidateLogOpts validates the opts passed to the scribelogs driver. Currently, the scribelogs
// driver doesn't take any arguments.
func ValidateLogOpts(cfg map[string]string) error {
	for k := range cfg {
		switch k {
		case streamName, ip, port:
		default:
			return fmt.Errorf("%q is not a valid option for the scribelogs driver", k)
		}
	}
	return nil
}
*/

func (l *scribelogs) Log(m *logger.Message) error {
	_, err := l.client.Log([]*scribe.LogEntry{{streamName, string(m.Line)}})
	return err
}

func (l *scribelogs) Close() error {
	return nil
}

func (l *scribelogs) Name() string {
	return name
}
