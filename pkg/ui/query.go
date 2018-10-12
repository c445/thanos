package ui

import (
	"html/template"
	"net/http"
	"path"

	"time"

	"os"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/common/route"
	"github.com/prometheus/common/version"
)

var localhostRepresentations = []string{"127.0.0.1", "localhost"}

type Query struct {
	*BaseUI

	flagsMap map[string]string

	cwd   string
	birth time.Time
	now   func() model.Time
}

type thanosVersion struct {
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	Branch    string `json:"branch"`
	BuildUser string `json:"buildUser"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}

func NewQueryUI(logger log.Logger, flagsMap map[string]string) *Query {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "<error retrieving current working directory>"
	}
	return &Query{
		BaseUI:   NewBaseUI(logger, "query_menu.html", template.FuncMap{}),
		flagsMap: flagsMap,
		cwd:      cwd,
		birth:    time.Now(),
		now:      model.Now,
	}
}

// Register registers new GET routes for subpages and retirects from / to /graph.
func (q *Query) Register(r *route.Router) {
	instrf := prometheus.InstrumentHandlerFunc

	r.Get("/", instrf("root", q.root))
	r.Get("/graph", instrf("graph", q.graph))
	r.Get("/status", instrf("status", q.status))
	r.Get("/flags", instrf("flags", q.flags))

	r.Get("/static/*filepath", instrf("static", q.serveStaticAsset))
	// TODO(bplotka): Consider adding more Thanos related data e.g:
	// - what store nodes we see currently
	// - what sidecars we see currently
}

// root redirects / requests to /graph, taking into account the path prefix value
func (q *Query) root(w http.ResponseWriter, r *http.Request) {
	prefix := GetWebPrefix(q.logger, q.flagsMap, r)

	http.Redirect(w, r, path.Join(prefix, "/graph"), http.StatusFound)
}

func (q *Query) graph(w http.ResponseWriter, r *http.Request) {
	prefix := GetWebPrefix(q.logger, q.flagsMap, r)

	q.executeTemplate(w, r, "graph.html", prefix, nil)
}

func (q *Query) status(w http.ResponseWriter, r *http.Request) {
	prefix := GetWebPrefix(q.logger, q.flagsMap, r)

	q.executeTemplate(w, r, "status.html", prefix, struct {
		Birth   time.Time
		CWD     string
		Version thanosVersion
	}{
		Birth: q.birth,
		CWD:   q.cwd,
		Version: thanosVersion{
			Version:   version.Version,
			Revision:  version.Revision,
			Branch:    version.Branch,
			BuildUser: version.BuildUser,
			BuildDate: version.BuildDate,
			GoVersion: version.GoVersion,
		},
	})
}

func (q *Query) flags(w http.ResponseWriter, r *http.Request) {
	prefix := GetWebPrefix(q.logger, q.flagsMap, r)

	q.executeTemplate(w, r, "flags.html", prefix, q.flagsMap)
}
