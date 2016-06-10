package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/alcortesm/tgpl-exercises/ch01/e12/lissajous"
)

var errHelp = errors.New("")

func main() {
	http.HandleFunc("/", lissajousGif)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func dumpRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

func lissajousGif(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}

	conf, err := formToConf(r.Form)
	if err != nil {
		if err == errHelp {
			fmt.Fprintf(w, help)
			return
		}
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}

	lissajous.Gif(w, conf)
}

func formToConf(forms url.Values) (*lissajous.Conf, error) {
	if len(forms) == 0 {
		return nil, errHelp
	}

	conf := lissajous.DefaultConf()
	var err error

	for k, v := range forms {
		if len(v) != 1 {
			return nil, fmt.Errorf(
				"bad number of arguments to %q form: expected 1, found %d",
				k, len(v))
		}
		switch k {
		case "cycles":
			conf.Cycles, err = strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf(
					"bad cycle value, an int was expected but %s was found",
					v[0])
			}
		case "res":
			conf.Res, err = strconv.ParseFloat(v[0], 64)
			if err != nil {
				return nil, fmt.Errorf(
					"bad res value, a float was expected but %s was found",
					v[0])
			}
		case "side":
			conf.Side, err = strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf(
					"bad side value, an int was expected but %s was found",
					v[0])
			}
		case "nframes":
			conf.NFrames, err = strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf(
					"bad nframes value, an int was expected but %s was found",
					v[0])
			}
		case "delay":
			conf.Delay, err = strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf(
					"bad delay value, an int was expected but %s was found",
					v[0])
			}
		case "phaseInc":
			conf.PhaseInc, err = strconv.ParseFloat(v[0], 64)
			if err != nil {
				return nil, fmt.Errorf(
					"bad phaseInc value, an int was expected but %s was found",
					v[0])
			}
		case "FreqDiff":
			conf.FreqDiff, err = strconv.ParseFloat(v[0], 64)
			if err != nil {
				return nil, fmt.Errorf(
					"bad phaseInc value, an int was expected but %s was found",
					v[0])
			}
		}
	}

	return conf, nil
}

const help = `<html>
<body>
<h1>What is this?</h1>

<p>Foo bar.</p>

<h1>Usage</h1>

Accepted forms:
<ul>
<li>cycles   = <int>: number of complete x oscillator revolutions</li>
<li>res      = <float>:  angular resolution</li>
<li>side     = <int>:    image canvas side in pixels [0..side]</li>
<li>nframes  = <int>:     number of animation frames</li>
<li>delay    = <int>:      delay between frames in 10ms units</li>
<li>phaseInc = <float>:    how much phase to increment in each frame</li>
<li>freqDiff = <float>:    frequency difference between x and y</li>
</ul>


<h1>Examples</h1>
<ul>
	<li>
<a href="http://localhost:8000/?cycle=4&freqDiff=2.3&phaseInc=0.1">http://localhost:8000/?cycle=4&freqDiff=2.3&phaseInc=0.1</a>
	</li>
</ul>

</body>
</html>`
