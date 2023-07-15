package marionette

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"git.sr.ht/~michl/quickbeam/internal/protocol"
	"git.sr.ht/~michl/quickbeam/internal/web"
	"github.com/njasm/marionette_client"
	"github.com/shibukawa/configdir"
)

var (
	defaultFFPrefs = map[string]string{
		"startup.homepage_welcome_url.additional": "''",
		"devtools.errorconsole.enabled":           "true",
		"devtools.chrome.enabled":                 "true",

		// allow microphone and camera
		"media.navigator.permission.disabled": "true",

		// Send Browser Console (different from Devtools console) output to
		// STDOUT.
		"browser.dom.window.dump.enabled": "true",

		// From:
		// http://hg.mozilla.org/mozilla-central/file/1dd81c324ac7/build/automation.py.in//l388
		// Make url-classifier updates so rare that they won"t affect tests.
		"urlclassifier.updateinterval": "172800",
		// Point the url-classifier to a nonexistent local URL for fast failures.
		"browser.safebrowsing.provider.0.gethashURL": "'http://localhost/safebrowsing-dummy/gethash'",
		"browser.safebrowsing.provider.0.keyURL":     "'http://localhost/safebrowsing-dummy/newkey'",
		"browser.safebrowsing.provider.0.updateURL":  "'http://localhost/safebrowsing-dummy/update'",

		// Disable self repair/SHIELD
		"browser.selfsupport.url": "'https://localhost/selfrepair'",
		// Disable Reader Mode UI tour
		"browser.reader.detectedFirstArticle": "true",

		// Set the policy firstURL to an empty string to prevent
		// the privacy info page to be opened on every "web-ext run".
		// (See #1114 for rationale)
		"datareporting.policy.firstRunURL": "''",
	}
)

type Firefox struct {
	FirefoxPath string
	ProfilePath string
	Headless bool
	client *marionette_client.Client
	process *os.Process
	stdout io.ReadCloser
	stderr io.ReadCloser
	emptyTab bool
}

func NewFirefox() *Firefox {
	return &Firefox{}
}

func (f *Firefox) Start() (err error) {
	return start(f, cmdExecute{})
}

// TODO invalidate all pages?
func (f *Firefox) Quit() (err error) {
	client := f.client
	f.client = nil
	process := f.process
	f.process = nil
	if client != nil {
		_, err = client.Quit()
		if err != nil {
			if process != nil {
				err = process.Kill()
				if err != nil {
					return protocol.EnvironmentError(
						"Could not kill Firefox: %v", err)
				}
			}
			return protocol.EnvironmentError(
				"Could not quit Firefox: %v", err)
		}
	}
	return nil
}

func (f *Firefox) Running() bool {
	return f.client != nil && f.process != nil
}


func (f *Firefox) NewPage() (res web.Page, err error) {
	page := Page{client: f.client,}
	if f.emptyTab {
		page.pageName, err = f.client.GetWindowHandle()
	} else {
		r, err := f.client.NewWindow(true, "tab", false)
		if err != nil {
			return nil, err
		}

		var ok bool
		page.pageName, ok = responseGetString(r, "name")
		if !ok {
			log.Println("marionette.Page.Start: no name attribute of type string in return value")
			return nil, err
		}

	}
	return &page, nil
}


func start(f *Firefox, shell cmdExecuter) (err error) {
	if f.process == nil {
		err = startFirefox(f, shell)
	}
	if err != nil {
		return err
	}
	if f.client == nil {
		err = startMarionette(f)
	}
	f.emptyTab = true
	return
}

func startFirefox(f *Firefox, shell cmdExecuter) (err error) {
	if f.FirefoxPath == "" {
		f.FirefoxPath = shell.ExecOrEmpty("which firefox")
	}
	if f.FirefoxPath == "" {
		return protocol.ConfigurationError("Firefox executable not found")
	}

	if f.ProfilePath == "" {
		f.ProfilePath = getFirefoxProfilePath()
	}

	args := []string{"--marionette"}
	if f.Headless {
		args = append(args, "--headless")
	}
	args = append(args, "--profile", f.ProfilePath)

	firefoxCmd := exec.Command(f.FirefoxPath, args...)
	f.stdout, err = firefoxCmd.StdoutPipe()
	if err != nil {
		return protocol.InternalError(
			fmt.Sprintf("Could not connect Firefox' stdout: %v", err))
	}
	f.stderr, err = firefoxCmd.StderrPipe()
	if err != nil {
		return protocol.InternalError(
			fmt.Sprintf("Could not connect Firefox' stderr: %v", err))
	}
	if err = firefoxCmd.Start(); err != nil {
		return protocol.ConfigurationError(
			fmt.Sprintf("Could not start Firefox: %v", err))
	}
	f.process = firefoxCmd.Process
	return nil
}

func startMarionette(f *Firefox) (err error) {
	connected := false
	start := time.Now()
	for time.Since(start) < 30*time.Second {
		f.client = marionette_client.NewClient()
		err := f.client.Connect("127.0.0.1", 2828)
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		} else {
			connected = true
			break
		}
	}
	if !connected {
		return protocol.EnvironmentError("Could not connect to Firefox Marionette")
	}
	f.client.NewSession("", nil)
	for key, value := range defaultFFPrefs {
		setFirefoxPreference(f.client, key, value)
	}
	return
}

type cmdExecuter interface {
	ExecOrEmpty(string) string
}

type cmdExecute struct {
}

func (_ cmdExecute) ExecOrEmpty(command string) string {
	parts := strings.Fields(command)
	head := parts[0]
	parts = parts[1:]
	out, err := exec.Command(head, parts...).CombinedOutput()
	if err != nil {
		errorMessage := fmt.Sprintf(
			"shell cmd `%s` failed with: %s", command, string(out))
		log.Println(errorMessage)
		return ""
	}
	return strings.TrimSpace(string(out))
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

func setFirefoxPreference(client *marionette_client.Client, key string, value string) {
	var script string
	client.SetContext(marionette_client.CHROME)
	script = fmt.Sprintf(`
		Components.utils.import("resource://gre/modules/Preferences.jsm");
		prefs = new Preferences({defaultBranch: "root"});
    prefs.set("%s", %s);`, key, value)
	args := []interface{}{}
	client.ExecuteScript(script, args, 1000, false)
	client.SetContext(marionette_client.CONTENT)
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
