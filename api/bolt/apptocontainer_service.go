package bolt//click2cloud-apptocontainer
import (
	"github.com/akash-Click2cloud/dolphin/api"
	"os/exec"
	"bytes"
	 "strconv"
	"github.com/akash-Click2cloud/dolphin/api/file"

)
var Output       string
type AppToContainerService struct {
	store *Store
}


// CreateEndpoint assign an ID to a new endpoint and saves it.

func (service *AppToContainerService) BuildAppToContainer(atoc *dockm.AToC) ( error , string ) {
	var TLSCertPath="C:\\data"+file.TLSStorePath+strconv.Itoa(atoc.EndPointId)+"cert.pem"
	var TLSCaPath  ="C:\\data"+file.TLSStorePath+strconv.Itoa(atoc.EndPointId)+"ca.pem"
	var TLkeyPath="C:\\data"+file.TLSStorePath+strconv.Itoa(atoc.EndPointId)+"key.pem"
	var DockerURL="tcp://"+atoc.EndPointUrl+":2376"
	var TLS="true"
	command:="s2i"
	comarg := []string{"build",atoc.GitUrl,atoc.BaseImage,atoc.ImageName,"--ca",TLSCaPath,"--cert",TLSCertPath,"--key",TLkeyPath,"--tls",TLS, "--url",DockerURL}
	cmd := exec.Command(command,comarg...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	stdout:=outb.String()
	stderr:=errb.String()
	if err != nil {
		Output=stdout+stderr
	}
	Output=stdout+stderr
	return err ,Output
}
//func (service *AppToContainerService) OutputData () string{
//         return  Output
//	}
