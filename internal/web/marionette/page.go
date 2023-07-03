package marionette

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"git.sr.ht/~michl/quickbeam/internal/web"
	"github.com/njasm/marionette_client"
	"github.com/shibukawa/configdir"
)

type Page struct {
	FirefoxPath string
	Headless bool
	client *marionette_client.Client
	pageName string
	firefox Firefox
}

func responseDecoder(response *marionette_client.Response) *json.Decoder {
	value := bytes.NewBuffer([]byte(response.Value))
	return json.NewDecoder(value)
}

func responseGetString(response *marionette_client.Response, key string) (value string, ok bool) {
	decoder := responseDecoder(response)
	var m map[string]interface{}
	err := decoder.Decode(m)
	if err != nil {
		return "", false
	}
	v, ok := m[key]
	if !ok {
		return "", false
	}
	value, ok = v.(string)
	return
}

func NewPage() *Page {
	return &Page{}
}

func (p *Page) activate() (err error) {
	err = p.client.SwitchToWindow(p.pageName)
	if err != nil {
		log.Printf("marionette.Page.activate: %v\n", err)
	}
	return
}

func (p *Page) Start() (err error) {
	if p.pageName != "" && p.client != nil {
		log.Println("marionette.Page.Start: already started")
		return nil
	}
	if p.client == nil {
		connected := false
		start := time.Now()
		for time.Since(start) < 30*time.Second {
			p.client = marionette_client.NewClient()
			err := p.client.Connect("127.0.0.1", 2828)
			if err != nil {
				time.Sleep(10 * time.Millisecond)
				continue
			} else {
				connected = true
				break
			}
		}
		if !connected {
			return errors.New("Could not connect to Firefox Marionette")
		}
		p.client.NewSession("", nil)
	}
	if false {
		r, err := p.client.NewWindow(true, "tab", false)
		if err != nil {
			return err
		}

		var ok bool
		p.pageName, ok = responseGetString(r, "name")
		if !ok {
			log.Println("marionette.Page.Start: no name attribute of type string in return value")
			return err
		}
	} else {
		var err error
		p.pageName, err = p.client.GetWindowHandle()
		if err != nil {
			log.Println("marionette.Page.Start: could not get window handle")
		}
	}
	return nil
}

func (p *Page) KillBrowser() {
	if p.firefox.Process != nil {
		p.firefox.Process.Kill()
	}
}

func (p *Page) Close() {
	if p.firefox.Process != nil {
		p.firefox.Process.Kill()
	}
}

func (p *Page) Navigate(url string) (err error) {
	err = p.activate()
	if err != nil {
		return
	}
	_, err = p.client.Navigate(url)
	if err != nil {
		log.Printf("marionette.Page.Navigate: %v", err)
		return
	}
	return
}

func (p *Page) Back() {
	p.activate()
	_ = p.client.Back()
}

func (p *Page) Forward() {
	p.activate()
	_ = p.client.Forward()
}

func (p *Page) Root() web.Noder {
	root, _ := p.client.FindElement(marionette_client.By(marionette_client.CSS_SELECTOR), "body")
	return &Node{root}
}

func (p *Page) Running() bool {
	return true
}

func shell(command string) string {
	parts := strings.Fields(command)
	head := parts[0]
	parts = parts[1:]
	out, err := exec.Command(head, parts...).CombinedOutput()
	if err != nil {
		errorMessage := fmt.Sprintf(
			"shell cmd `%s` failed with: %s", command, string(out))
		log.Println(errorMessage)
	}
	return strings.TrimSpace(string(out))
}

func whichFirefox() string {
	return shell("which firefox")
}

func getConfigNamespace() string {
	return "quickbeam"
}

// Gets a cross-platform path to store a Browsh-specific Firefox profile
func getFirefoxProfilePath() string {
	configDirs := configdir.New(getConfigNamespace(), "firefox_profile")
	folders := configDirs.QueryFolders(configdir.Global)
	folders[0].MkdirAll()
	return folders[0].Path
}

type Firefox struct {
	Process *os.Process
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

func (p *Page) StartBrowser() (err error) {
	if p.FirefoxPath == "" {
		p.FirefoxPath = whichFirefox()
	}
	firefox, err := startFirefox(p.FirefoxPath, p.Headless)
	if err != nil {
		return err
	}
	p.firefox = firefox
	return nil
}

func startFirefox(path string, headless bool) (firefox Firefox, err error) {
	f := Firefox{}
	args := []string{"--marionette"}
	if headless {

		args = append(args, "--headless")
	}
	profilePath := getFirefoxProfilePath()
	log.Println("Using profile at: " + profilePath)
	args = append(args, "--profile", profilePath)
	firefoxCmd := exec.Command(path, args...)
	f.Stdout, err = firefoxCmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("startFirefox/stdout-pipe: %w", err)
		return
	}
	f.Stderr, err = firefoxCmd.StderrPipe()
	if err != nil {
		err = fmt.Errorf("startFirefox/stderr-pipe: %w", err)
		return
	}
	if err = firefoxCmd.Start(); err != nil {
		err = fmt.Errorf("startFirefox start: %w", err)
		return
	}
	f.Process = firefoxCmd.Process
	return f, nil
}

func (p *Page) Execute(js string) (string, error) {
	args := []interface{}{}
	r, err := p.client.ExecuteScript(js, args, 10000, false)
	if err != nil {
		return "", err
	}
	return r.Value, nil
}
