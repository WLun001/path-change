package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func getBranch(ref string) string {
	return strings.Split(ref, "/")[2]
}

func ensureHomeEnv(homepath string) {
	homeenv := os.Getenv("HOME")
	if _, err := os.Stat(filepath.Join(homeenv, ".ssh")); err != nil {
		// There's no $HOME/.ssh directory to access or the user doesn't have permissions
		// to read it, or something else; in any event there's no need to try creating a
		// symlink to it.
		return
	}
	if homeenv != "" {
		ensureHomeEnvSSHLinkedFromPath(homeenv, homepath)
	}
}

func ensureHomeEnvSSHLinkedFromPath(homeenv string, homepath string) {
	if filepath.Clean(homeenv) != filepath.Clean(homepath) {
		homeEnvSSH := filepath.Join(homeenv, ".ssh")
		homePathSSH := filepath.Join(homepath, ".ssh")
		if _, err := os.Stat(homePathSSH); os.IsNotExist(err) {
			if err := os.Symlink(homeEnvSSH, homePathSSH); err != nil {
				// Only do a warning, in case we don't have a real home
				// directory writable in our image
				log.Printf("Unexpected error: creating symlink: %v\n", err)
			}
		}
	}
}

// Canonical updates the map keys to use the Canonical name
func Canonical(h map[string][]string) http.Header {
	c := map[string][]string{}
	for k, v := range h {
		c[http.CanonicalHeaderKey(k)] = v
	}
	return http.Header(c)
}
