package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/abbot/go-http-auth"
	"golang.org/x/net/webdav"
)

var (
	flPort     = flag.Int("port", 9999, "server port")
	flCert     = flag.String("cert", "", "server SSL cert (both -cert and -key must be present to use SSL). See `go run $(go env GOROOT)/src/crypto/tls/generate_cert.go -h` to generate development cert/key")
	flKey      = flag.String("key", "", "server SSL key")
	flHtpasswd = flag.String("htpasswd", "", "htpasswd file for auth (must be present to use auth) See htpasswd(1) to create this file.")
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("One argument required. Please provide path to serve (you can use a special keyword of 'mem' to serve an in-memory filesystem)")
	}

	var fs webdav.FileSystem
	if flag.Args()[0] == "mem" {
		fs = webdav.NewMemFS()
	} else {
		fs = NewPassThroughFS(flag.Args()[0])
	}
	log.SetFlags(0)
	h := &webdav.Handler{
		FileSystem: fs,
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			t := time.Now()
			tStamp := fmt.Sprintf("%d.%9.9d", t.Unix(), t.Nanosecond())
			switch r.Method {
			case "COPY", "MOVE":
				dst := ""
				if u, err := url.Parse(r.Header.Get("Destination")); err == nil {
					dst = u.Path
				}
				o := r.Header.Get("Overwrite")
				log.Printf("%-21s%-25s%-10s%-30s%-30so=%-2s%v", tStamp, r.RemoteAddr, r.Method, r.URL.Path, dst, o, err)
			default:
				log.Printf("%-21s%-25s%-10s%-30s%v", tStamp, r.RemoteAddr, r.Method, r.URL.Path, err)
			}
		},
	}
	if *flHtpasswd != "" {
		secret := auth.HtpasswdFileProvider(*flHtpasswd)
		authenticator := auth.NewBasicAuthenticator("", secret)
		authHandlerFunc := func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
			h.ServeHTTP(w, &r.Request)
		}
		http.HandleFunc("/", authenticator.Wrap(authHandlerFunc))

	} else {
		log.Println("WARNING: connections are not authenticated. STRONGLY consider using -htpasswd.")
		http.Handle("/", h)
	}
	addr := fmt.Sprintf(":%d", *flPort)
	if *flCert != "" && *flKey != "" {
		log.Printf("Serving HTTPS:// %v", addr)
		log.Fatal(http.ListenAndServeTLS(addr, *flCert, *flKey, nil))
	} else {
		log.Println("WARNING: connections are not encrypted. STRONGLY consider using -cert/-key.")
		log.Printf("Serving HTTP:// %v", addr)
		log.Fatal(http.ListenAndServe(addr, nil))
	}
}

func NewPassThroughFS(path string) webdav.FileSystem {
	return &passThroughFS{root: path}
}

type passThroughFS struct {
	root string
}

func (ptfs *passThroughFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	// TODO(vbatts) check for escaping the root directory
	return os.Mkdir(filepath.Join(ptfs.root, name), perm)
}
func (ptfs *passThroughFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	// TODO(vbatts) check for escaping the root directory
	return os.OpenFile(filepath.Join(ptfs.root, name), flag, perm)
}
func (ptfs *passThroughFS) RemoveAll(ctx context.Context, name string) error {
	// TODO(vbatts) check for escaping the root directory
	return os.RemoveAll(filepath.Join(ptfs.root, name))
}
func (ptfs *passThroughFS) Rename(ctx context.Context, oldName, newName string) error {
	// TODO(vbatts) check for escaping the root directory
	return os.Rename(filepath.Join(ptfs.root, oldName), filepath.Join(ptfs.root, newName))
}
func (ptfs *passThroughFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	// TODO(vbatts) check for escaping the root directory
	return os.Stat(filepath.Join(ptfs.root, name))
}
