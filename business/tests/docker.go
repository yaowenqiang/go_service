package tests
import (
    "testing"
    "bytes"
    "net"
    "os/exec"
    "encoding/json"
)

type Container struct {
    ID string
    Host string // IP:Port
}


func startContainer(t *testing.T, image string, port string, args ...string) *Container {
    arg := []string{"run", "-P", "-d"}
    arg = append(arg, args...)
    arg = append(arg, image)

    cmd := exec.Command("docker", arg...)

    var out bytes.Buffer
    cmd.Stdout = &out

    if err := cmd.Run(); err  != nil {
        t.Fatalf("could not start container %s : %v", image, err)
    }

    id := out.String()[:12]
    cmd = exec.Command("docker", "inspect", id)

    out.Reset()

    cmd.Stdout = &out

    if err := cmd.Run(); err != nil {
        t.Fatalf("could not inspect container %s : %v", image, err)
    }

    var doc []map[string]interface{}
    if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
        t.Fatalf("could not decode  json  %v",  err)
    }

    ip, randPort := extractIPPort(t, doc, port)

    c := Container{
        ID: id,
        Host: net.JoinHostPort(ip, randPort),
    }

    t.Logf("Image:            %s", image)
    t.Logf("ContainerID:      %s", c.ID)
    t.Logf("Host:             %s", c.Host)

    return &c

}

func stopContainer(t *testing.T, id string) {
    if err := exec.Command("docker", "stop", id).Run(); err != nil {
        t.Fatalf("could not stop container %s : %v", id, err)
    }
    t.Log("Stopped", id)

    if err := exec.Command("docker", "rm", id, "-v").Run(); err != nil {
        t.Fatalf("could not remove container %s : %v", id, err)
    }
    t.Log("Removed", id)

}

func dumpContainerLogs(t *testing.T, id string) {
    out ,err := exec.Command("docker", "logs" , id).CombinedOutput()
    if err != nil {
        t.Fatalf("could not log container %s : %v", id, err)
    }
    t.Logf("logs for %s\n%s", id, out)
}

func extractIPPort(t *testing.T, doc []map[string]interface{}, port string) (string, string) {
    nw, exists := doc[0]["NetworkSettings"]
    if !exists {
        t.Fatal("could not get network settings")
    }
    ports, exists := nw.(map[string]interface{})["Ports"]
    if !exists {
        t.Fatal("could not get network port settings")
    }
    tcp, exists := ports.(map[string]interface{})[port + "/tcp"]
    if !exists {
        t.Fatal("could not get network ports/tcp settings")
    }
    list, exists := tcp.([]interface{})
    if !exists {
        t.Fatal("could not get network ports/tcp list settings")
    }
    if len(list) != 1 {
        t.Fatal("could not get network ports/tcp list settings")
    }
    data, exists := list[0].(map[string]interface{})
    if !exists {
        t.Fatal("could not get network ports/tcp list data")
    }

    return data["HostIp"].(string), data["HostPort"].(string)
}
