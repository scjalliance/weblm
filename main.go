package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	shellwords "github.com/mattn/go-shellwords"
)

func main() {
	httpHostPortDefault := "0.0.0.0:7259"
	lmutilPathDefault := "/opt/lmutil"
	lmstatArgDefault := "lmstat -a"

	if os.Getenv("WEBLM_HOSTPORT") != "" {
		httpHostPortDefault = os.Getenv("WEBLM_HOSTPORT")
	}
	httpHostPort := flag.String("http", httpHostPortDefault, "http host:port value (or set env var WEBLM_HOSTPORT)")

	if os.Getenv("WEBLM_LMUTIL") != "" {
		lmutilPathDefault = os.Getenv("WEBLM_LMUTIL")
	}
	lmutilPath := flag.String("lmutil", lmutilPathDefault, "lmutil binary path (or set env var WEBLM_LMUTIL)")

	if os.Getenv("WEBLM_LMSTATARG") != "" {
		lmstatArgDefault = os.Getenv("WEBLM_LMSTATARG")
	}
	lmstatArg := flag.String("lmstat", lmstatArgDefault, "lmstat arguments (or set env var WEBLM_LMSTATARG)")

	flag.Parse()

	lmstatArgs, err := shellwords.Parse(*lmstatArg)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/lmstat", func(w http.ResponseWriter, r *http.Request) { lmstatHandler(w, r, *lmutilPath, lmstatArgs) })
	http.ListenAndServe(*httpHostPort, nil)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Nothing at this URL.", http.StatusNotFound)
}

func lmstatHandler(w http.ResponseWriter, r *http.Request, lmutilPath string, lmstatArgs []string) {
	c := exec.Command(lmutilPath, lmstatArgs...)
	out, err := c.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprint(w, string(out))
}
